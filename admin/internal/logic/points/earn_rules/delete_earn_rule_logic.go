package earn_rules

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

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
	tenantID, _ := contextx.GetTenantID(l.ctx)
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	return l.svcCtx.PointsService.DeleteEarnRule(l.ctx, shared.TenantID(tenantID), req.ID)
}
