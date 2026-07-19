package promotions

import (
	"context"

	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	pkgpromotion "github.com/colinrs/shopjoy/pkg/domain/promotion"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePromotionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePromotionLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdatePromotionLogic {
	return UpdatePromotionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdatePromotion mirrors CreatePromotion: parse times, build a
// PromotionRule from the wire discount fields, then call the unified
// PromotionApp.Update with Rules wrapped in a pointer (so the app
// layer knows to replace — not preserve — the rules).
func (l *UpdatePromotionLogic) UpdatePromotion(req *types.UpdatePromotionReq) (resp *types.PromotionDetailResp, err error) {
	startAt, err := parseTime(req.StartTime)
	if err != nil {
		return nil, err
	}
	endAt, err := parseTime(req.EndTime)
	if err != nil {
		return nil, err
	}

	scope := buildPromotionScope(req.ScopeType, req.ProductIDs, req.CategoryIDs, req.BrandIDs)

	var rulesPtr *[]pkgpromotion.PromotionRule
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
		rules := []pkgpromotion.PromotionRule{rule}
		rulesPtr = &rules
	}

	updateReq := &apppromotion.UpdatePromotionRequest{
		ID:           req.ID,
		Name:         req.Name,
		Description:  req.Description,
		Type:         parsePromotionType(req.Type),
		Currency:     currencyWithDefault(req.Currency),
		Scope:        scope,
		UsageLimit:   req.UsageLimit,
		PerUserLimit: req.PerUserLimit,
		Tags:         req.Tags,
		StartAt:      startAt,
		EndAt:        endAt,
		Rules:        rulesPtr,
		ActorID:      actorID(l.ctx),
	}

	promotionResp, err := l.svcCtx.PromotionApp.Update(l.ctx, updateReq)
	if err != nil {
		return nil, err
	}

	return convertPromotionToDetailResp(promotionResp), nil
}

// actorID extracts the audit user ID from context with a 0 fallback.
func actorID(ctx context.Context) int64 {
	if v := ctx.Value("user_id"); v != nil {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	return 0
}