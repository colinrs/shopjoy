package promotions

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	pkgpromotion "github.com/colinrs/shopjoy/pkg/domain/promotion"

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

// UpdatePromotionRule mutates a single rule by ID. The app layer's
// UpdateRule preserves OwnerKind / OwnerID — we don't need to fetch
// the parent promotion here.
func (l *UpdatePromotionRuleLogic) UpdatePromotionRule(req *types.UpdatePromotionRuleReq) (resp *types.PromotionRuleResp, err error) {
	rule := &pkgpromotion.PromotionRule{
		ID:            req.ID,
		ConditionType: mapConditionType(req.RuleType),
		ConditionValue: parseMoneyToDecimal(req.Value),
		ActionType:     mapDiscountActionType(req.DiscountType),
		ActionValue:    parseMoneyToDecimal(req.DiscountValue),
	}
	ruleResp, err := l.svcCtx.PromotionApp.UpdateRule(l.ctx, rule)
	if err != nil {
		return nil, err
	}

	return &types.PromotionRuleResp{
		ID:            ruleResp.ID,
		RuleType:      mapConditionTypeToString(ruleResp.ConditionType),
		Operator:      "gte",
		Value:         formatDecimalToString(ruleResp.ConditionValue),
		DiscountType:  mapActionTypeIntToString(ruleResp.ActionType),
		DiscountValue: formatDecimalToString(ruleResp.ActionValue),
		Priority:      ruleResp.SortOrder,
		CreatedAt:     "",
		UpdatedAt:     "",
	}, nil
}