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

// GetPromotionRules looks up the owner Promotion's Kind, then asks
// the unified PromotionApp for its rules. The wire response keeps the
// old shape (rule_type / value / discount_type / discount_value) so
// convertRulesToResp is reused.
func (l *GetPromotionRulesLogic) GetPromotionRules(req *types.GetPromotionRulesReq) (resp *types.ListPromotionRulesResp, err error) {
	owner, err := l.svcCtx.PromotionApp.Get(l.ctx, req.PromotionID)
	if err != nil {
		return nil, err
	}
	rules, err := l.svcCtx.PromotionApp.GetRules(l.ctx, owner.Kind, req.PromotionID)
	if err != nil {
		return nil, err
	}

	// Inject the owner ID into each rule so the wire response can
	// carry a PromotionID back to the form (used by inline-edit
	// controls that round-trip the rule).
	wire := convertRulesToResp(rules)
	for _, r := range wire {
		if r == nil {
			continue
		}
		r.PromotionID = req.PromotionID
	}

	return &types.ListPromotionRulesResp{
		List:  wire,
		Total: int64(len(wire)),
	}, nil
}