package promotions

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	pkgpromotion "github.com/colinrs/shopjoy/pkg/domain/promotion"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePromotionRulesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePromotionRulesLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreatePromotionRulesLogic {
	return CreatePromotionRulesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePromotionRulesLogic) CreatePromotionRules(req *types.CreatePromotionRulesReq) (resp *types.CreatePromotionRulesResp, err error) {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Get existing promotion
	p, err := l.svcCtx.PromotionApp.GetPromotion(l.ctx, shared.TenantID(tenantID), req.PromotionID)
	if err != nil {
		return nil, err
	}

	// Convert rules
	rules := make([]pkgpromotion.PromotionRule, 0, len(req.Rules))
	ids := make([]int64, 0, len(req.Rules))

	for _, ruleReq := range req.Rules {
		id, err := l.svcCtx.IDGen.NextID(l.ctx)
		if err != nil {
			return nil, err
		}

		rule := pkgpromotion.PromotionRule{
			ID:            id,
			PromotionID:   req.PromotionID,
			ConditionType: mapConditionType(ruleReq.RuleType),
			ActionType:    mapDiscountActionType(ruleReq.DiscountType),
		}

		// Parse condition value
		rule.ConditionValue = parseMoneyToDecimal(ruleReq.Value)

		// Parse action value
		rule.ActionValue = parseMoneyToDecimal(ruleReq.DiscountValue)

		rules = append(rules, rule)
		ids = append(ids, id)
	}

	// Update promotion with new rules (this would require adding a method to the app service)
	// For now, we return the IDs
	_ = p // promotion retrieved for validation

	return &types.CreatePromotionRulesResp{
		IDs: ids,
	}, nil
}
