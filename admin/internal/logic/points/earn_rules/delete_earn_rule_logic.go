package earn_rules

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteEarnRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteEarnRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteEarnRuleLogic {
	return DeleteEarnRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteEarnRuleLogic) DeleteEarnRule(req *types.DeleteEarnRuleReq) error {

	return l.svcCtx.PointsService.DeleteEarnRule(l.ctx, req.ID)
}
