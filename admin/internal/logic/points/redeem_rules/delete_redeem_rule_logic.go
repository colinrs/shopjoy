package redeem_rules

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRedeemRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteRedeemRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteRedeemRuleLogic {
	return DeleteRedeemRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteRedeemRuleLogic) DeleteRedeemRule(req *types.DeleteRedeemRuleReq) error {

	return l.svcCtx.PointsService.DeleteRedeemRule(l.ctx, req.ID)
}
