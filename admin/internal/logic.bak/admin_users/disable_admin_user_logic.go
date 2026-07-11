package admin_users

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"

	"github.com/zeromicro/go-zero/core/logx"
)

type DisableAdminUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDisableAdminUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) DisableAdminUserLogic {
	return DisableAdminUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DisableAdminUserLogic) DisableAdminUser(req *types.AdminUserIDRequest) error {
	operatorID := contextx.GetCurrentUserID(l.ctx)
	return l.svcCtx.AdminUserService.Disable(l.ctx, operatorID, req.ID)
}
