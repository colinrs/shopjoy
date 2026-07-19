package promotions

import (
	"context"
	"strconv"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

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

// CreatePromotionRules resolves the owner's kind (PROMOTION or COUPON) by
// reading the unified Promotion record, then forwards the rule set to
// PromotionApp.CreateRules.
func (l *CreatePromotionRulesLogic) CreatePromotionRules(req *types.CreatePromotionRulesReq) (resp *types.CreatePromotionRulesResp, err error) {
	owner, err := l.svcCtx.PromotionApp.Get(l.ctx, req.ID)
	if err != nil {
		return nil, err
	}
	rules := convertRuleReqsToDomain(req.Rules)
	out, err := l.svcCtx.PromotionApp.CreateRules(l.ctx, owner.Kind, req.ID, rules)
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(out))
	for i, r := range out {
		ids[i] = strconv.FormatInt(r.ID, 10)
	}
	return &types.CreatePromotionRulesResp{IDs: ids}, nil
}