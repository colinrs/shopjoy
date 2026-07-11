package admin_users

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"

	"github.com/zeromicro/go-zero/core/logx"
)

type EnableAdminUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEnableAdminUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) EnableAdminUserLogic {
	return EnableAdminUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EnableAdminUserLogic) EnableAdminUser(req *types.AdminUserIDRequest) error {
	operatorID := contextx.GetCurrentUserID(l.ctx)
	return l.svcCtx.AdminUserService.Enable(l.ctx, operatorID, req.ID)
}
