package promotions

import (
	"context"

	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	pkgpromotion "github.com/colinrs/shopjoy/pkg/domain/promotion"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePromotionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePromotionLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreatePromotionLogic {
	return CreatePromotionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreatePromotion maps the wire request onto the unified
// CreatePromotionRequest. Wire fields now match the unified shape
// (Kind, MarketID, Code, TotalCount, Rules). The wire-level
// ScopeIDs / ExcludeIDs are folded into a single PromotionScope so
// the application layer doesn't need to know about the legacy
// per-type arrays.
func (l *CreatePromotionLogic) CreatePromotion(req *types.CreatePromotionReq) (resp *types.CreatePromotionResp, err error) {
	startAt, err := parseTime(req.StartTime)
	if err != nil {
		return nil, err
	}
	endAt, err := parseTime(req.EndTime)
	if err != nil {
		return nil, err
	}

	kind := parsePromotionKind(req.Kind)
	typeVal := parsePromotionType(req.Type)

	createReq := &apppromotion.CreatePromotionRequest{
		TenantID:     shared.TenantID(l.getTenantID()),
		Kind:         kind,
		Name:         req.Name,
		Description:  req.Description,
		Code:         optionalString(req.Code),
		Type:         typeVal,
		MarketID:     optionalInt64(req.MarketID),
		Currency:     currencyWithDefault(req.Currency),
		Scope:        buildPromotionScopeFromWire(req.ScopeType, req.ScopeIDs, req.ExcludeIDs),
		TotalCount:   optionalInt(req.TotalCount),
		UsageLimit:   req.UsageLimit,
		PerUserLimit: req.PerUserLimit,
		Tags:         req.Tags,
		StartAt:      startAt,
		EndAt:        endAt,
		Rules:        convertRuleReqsToDomainPtr(req.Rules),
		ActorID:      l.getUserID(),
	}

	promotionResp, err := l.svcCtx.PromotionApp.Create(l.ctx, createReq)
	if err != nil {
		return nil, err
	}

	return &types.CreatePromotionResp{
		ID: promotionResp.ID,
	}, nil
}

// getTenantID extracts the tenant ID from the request context. The
// admin API currently runs single-tenant; treat a missing value as 0
// (the app layer's repository filters then scope by the value).
func (l *CreatePromotionLogic) getTenantID() int64 {
	if v := l.ctx.Value("tenant_id"); v != nil {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	return 0
}

// getUserID is the same helper for the actor (audit) field. Missing
// values fall back to 0 which the repo treats as "system".
func (l *CreatePromotionLogic) getUserID() int64 {
	if v := l.ctx.Value("user_id"); v != nil {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	return 0
}

// parsePromotionKind maps the wire Kind ("promotion" / "coupon") onto
// the domain Kind enum. Unknown values default to KindPromotion so
// pre-existing forms keep working.
func parsePromotionKind(s string) pkgpromotion.Kind {
	switch s {
	case "coupon":
		return pkgpromotion.KindCoupon
	case "promotion":
		return pkgpromotion.KindPromotion
	default:
		return pkgpromotion.KindPromotion
	}
}

// optionalString returns a pointer to s, or nil if s is empty.
func optionalString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// optionalInt returns a pointer to v, or nil if v is zero.
func optionalInt(v int) *int {
	if v == 0 {
		return nil
	}
	return &v
}

// buildPromotionScopeFromWire folds the wire-level (scope_type,
// scope_ids, exclude_ids) into the domain's PromotionScope.
func buildPromotionScopeFromWire(scopeType string, scopeIDs, excludeIDs []string) pkgpromotion.PromotionScope {
	return pkgpromotion.PromotionScope{
		Type:       normalizeScopeType(scopeType),
		IDs:        parseInt64Slice(scopeIDs),
		ExcludeIDs: parseInt64Slice(excludeIDs),
	}
}

// normalizeScopeType maps a wire-level scope hint onto the domain's
// ScopeType enum. Empty / unknown values default to storewide.
func normalizeScopeType(scopeType string) pkgpromotion.ScopeType {
	switch scopeType {
	case "products":
		return pkgpromotion.ScopeTypeProducts
	case "categories":
		return pkgpromotion.ScopeTypeCategories
	case "brands":
		return pkgpromotion.ScopeTypeBrands
	case "storewide":
		return pkgpromotion.ScopeTypeStorewide
	default:
		return pkgpromotion.ScopeTypeStorewide
	}
}

// convertRuleReqsToDomainPtr maps the wire rule requests (with the
// unified ConditionType / ActionType / ConditionValue / ActionValue
// shape) onto the domain PromotionRule slice.
func convertRuleReqsToDomainPtr(reqs []*types.PromotionRuleReq) []pkgpromotion.PromotionRule {
	if len(reqs) == 0 {
		return nil
	}
	out := make([]pkgpromotion.PromotionRule, 0, len(reqs))
	for _, r := range reqs {
		if r == nil {
			continue
		}
		out = append(out, pkgpromotion.PromotionRule{
			ConditionType:  mapWireConditionType(r.ConditionType),
			ConditionValue: parseMoneyToDecimal(r.ConditionValue),
			ActionType:     mapWireActionType(r.ActionType),
			ActionValue:    parseMoneyToDecimal(r.ActionValue),
			MaxDiscount:    parseMoneyToDecimal(r.MaxDiscount),
			SortOrder:      r.SortOrder,
		})
	}
	return out
}

// mapWireConditionType maps the wire ConditionType onto the domain
// ConditionType. Wire sends "min_amount" / "min_quantity"; the domain
// uses ConditionMinAmount / ConditionMinQuantity.
func mapWireConditionType(s string) pkgpromotion.ConditionType {
	switch s {
	case "min_quantity":
		return pkgpromotion.ConditionMinQuantity
	case "min_amount":
		return pkgpromotion.ConditionMinAmount
	default:
		return pkgpromotion.ConditionMinAmount
	}
}

// mapWireActionType maps the wire ActionType onto the domain
// ActionType. Wire sends "fixed_amount" / "percentage" /
// "free_shipping"; the domain uses ActionFixedAmount / etc.
func mapWireActionType(s string) pkgpromotion.ActionType {
	switch s {
	case "percentage":
		return pkgpromotion.ActionPercentage
	case "free_shipping":
		return pkgpromotion.ActionFreeShipping
	case "fixed_amount":
		return pkgpromotion.ActionFixedAmount
	default:
		return pkgpromotion.ActionFixedAmount
	}
}
