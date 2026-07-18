package promotions

import (
	"context"

	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/utils"

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
	ruleReqs := make([]apppromotion.CreatePromotionRuleRequest, 0, len(req.Rules))

	for _, ruleReq := range req.Rules {
		ruleReqs = append(ruleReqs, apppromotion.CreatePromotionRuleRequest{
			ConditionType:  mapConditionType(ruleReq.RuleType),
			ConditionValue: parseMoneyToDecimal(ruleReq.Value),
			ActionType:     mapDiscountActionType(ruleReq.DiscountType),
			ActionValue:    parseMoneyToDecimal(ruleReq.DiscountValue),
		})
	}

	ids, err := l.svcCtx.PromotionApp.CreatePromotionRules(l.ctx, req.PromotionID, ruleReqs)
	if err != nil {
		return nil, err
	}

	return &types.CreatePromotionRulesResp{
		IDs: utils.FormatInt64Slice(ids),
	}, nil
}
