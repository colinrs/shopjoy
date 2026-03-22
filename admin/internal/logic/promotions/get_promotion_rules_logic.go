package promotions

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPromotionRulesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPromotionRulesLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetPromotionRulesLogic {
	return GetPromotionRulesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPromotionRulesLogic) GetPromotionRules(req *types.GetPromotionRulesReq) (resp *types.ListPromotionRulesResp, err error) {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Get promotion with rules
	promotionResp, err := l.svcCtx.PromotionApp.GetPromotion(l.ctx, shared.TenantID(tenantID), req.PromotionID)
	if err != nil {
		return nil, err
	}

	// Convert rules
	rules := make([]*types.PromotionRuleResp, 0, len(promotionResp.Rules))
	for _, rule := range promotionResp.Rules {
		rules = append(rules, &types.PromotionRuleResp{
			ID:            rule.ID,
			PromotionID:   rule.PromotionID,
			RuleType:      mapConditionTypeToString(rule.ConditionType),
			Operator:      "gte",
			Value:         formatInt64ToString(rule.ConditionValue),
			DiscountType:  mapActionTypeIntToString(rule.ActionType),
			DiscountValue: formatInt64ToString(rule.ActionValue),
			Priority:      0,
			CreatedAt:     promotionResp.CreatedAt,
			UpdatedAt:     promotionResp.UpdatedAt,
		})
	}

	return &types.ListPromotionRulesResp{
		List:  rules,
		Total: int64(len(rules)),
	}, nil
}