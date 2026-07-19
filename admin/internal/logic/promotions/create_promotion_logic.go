package promotions

import (
	"context"
	"strings"

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
// CreatePromotionRequest. The wire type still uses the pre-merge
// "PROMOTION only" shape (no Kind field); Task 8 will add it. Until
// then, the create path defaults Kind=KindPromotion so the app layer
// doesn't reject the request.
//
// Discount fields are optional on the wire. When present, they're
// folded into a single PromotionRule so the existing UI keeps working.
func (l *CreatePromotionLogic) CreatePromotion(req *types.CreatePromotionReq) (resp *types.CreatePromotionResp, err error) {
	startAt, err := parseTime(req.StartTime)
	if err != nil {
		return nil, err
	}
	endAt, err := parseTime(req.EndTime)
	if err != nil {
		return nil, err
	}

	kind := pkgpromotion.KindPromotion
	if req.Type == "coupon" {
		// Frontend sometimes signals a coupon via the type field on
		// the legacy wire shape. Defer to KindPromotion otherwise.
		kind = pkgpromotion.KindCoupon
	}

	createReq := &apppromotion.CreatePromotionRequest{
		TenantID:     shared.TenantID(l.getTenantID()),
		Kind:         kind,
		Name:         req.Name,
		Description:  req.Description,
		Type:         parsePromotionType(req.Type),
		Currency:     currencyWithDefault(req.Currency),
		Scope:        buildPromotionScope(req.ScopeType, req.ProductIDs, req.CategoryIDs, req.BrandIDs),
		UsageLimit:   req.UsageLimit,
		PerUserLimit: req.PerUserLimit,
		Tags:         req.Tags,
		StartAt:      startAt,
		EndAt:        endAt,
		ActorID:      l.getUserID(),
	}

	// Build a single rule from the wire discount fields when both
	// type and value are present.
	if req.DiscountType != "" && req.DiscountValue != "" {
		rule := pkgpromotion.PromotionRule{
			ConditionType: pkgpromotion.ConditionMinAmount,
			ActionType:    mapDiscountActionType(req.DiscountType),
		}
		if req.MinOrderAmount != "" {
			rule.ConditionValue = parseMoneyToDecimal(req.MinOrderAmount)
		}
		if req.DiscountValue != "" {
			rule.ActionValue = parseMoneyToDecimal(req.DiscountValue)
		}
		if req.MaxDiscount != "" {
			rule.MaxDiscount = parseMoneyToDecimal(req.MaxDiscount)
		}
		createReq.Rules = []pkgpromotion.PromotionRule{rule}
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

// strings is imported so we keep package layout consistent with
// sibling files; the helper below is also reachable.
var _ = strings.ToLower