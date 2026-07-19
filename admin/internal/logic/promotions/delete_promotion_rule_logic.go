package promotions

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePromotionRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletePromotionRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeletePromotionRuleLogic {
	return DeletePromotionRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DeletePromotionRule removes a single rule by ID via the unified
// app service. No need to look up the parent promotion — DeleteRule
// is owner-kind-agnostic.
func (l *DeletePromotionRuleLogic) DeletePromotionRule(req *types.DeletePromotionRuleReq) (resp *types.CreatePromotionResp, err error) {
	if err := l.svcCtx.PromotionApp.DeleteRule(l.ctx, req.ID); err != nil {
		return nil, err
	}
	return &types.CreatePromotionResp{ID: req.ID}, nil
}
