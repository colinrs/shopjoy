package promotions

import (
	"context"

	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePromotionRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePromotionRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdatePromotionRuleLogic {
	return UpdatePromotionRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePromotionRuleLogic) UpdatePromotionRule(req *types.UpdatePromotionRuleReq) (resp *types.PromotionRuleResp, err error) {
	updateReq := apppromotion.UpdatePromotionRuleRequest{
		ConditionType:  mapConditionType(req.RuleType),
		ConditionValue: parseMoneyToDecimal(req.Value),
		ActionType:     mapDiscountActionType(req.DiscountType),
		ActionValue:    parseMoneyToDecimal(req.DiscountValue),
	}

	ruleResp, err := l.svcCtx.PromotionApp.UpdatePromotionRule(l.ctx, req.ID, updateReq)
	if err != nil {
		return nil, err
	}

	return &types.PromotionRuleResp{
		ID:            ruleResp.ID,
		PromotionID:   ruleResp.PromotionID,
		RuleType:      mapConditionTypeToString(ruleResp.ConditionType),
		Operator:      "gte",
		Value:         ruleResp.ConditionValue.StringFixed(2),
		DiscountType:  mapActionTypeIntToString(ruleResp.ActionType),
		DiscountValue: ruleResp.ActionValue.StringFixed(2),
		Priority:      0,
		CreatedAt:     "",
		UpdatedAt:     "",
	}, nil
}
