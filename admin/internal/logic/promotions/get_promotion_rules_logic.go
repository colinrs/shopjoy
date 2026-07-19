package promotions

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

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

// GetPromotionRules resolves the owner's kind (PROMOTION or COUPON) by
// reading the unified Promotion record, then asks PromotionApp.GetRules.
// Wire response uses the unified rule shape (condition_type / condition_value
// / action_type / action_value).
func (l *GetPromotionRulesLogic) GetPromotionRules(req *types.GetPromotionRulesReq) (resp *types.ListPromotionRulesResp, err error) {
	owner, err := l.svcCtx.PromotionApp.Get(l.ctx, req.ID)
	if err != nil {
		return nil, err
	}
	rules, err := l.svcCtx.PromotionApp.GetRules(l.ctx, owner.Kind, req.ID)
	if err != nil {
		return nil, err
	}
	return &types.ListPromotionRulesResp{
		List:  convertRulesPtrToResp(rules),
		Total: int64(len(rules)),
	}, nil
}