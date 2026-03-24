package admin_users

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAdminUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteAdminUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteAdminUserLogic {
	return DeleteAdminUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteAdminUserLogic) DeleteAdminUser(req *types.AdminUserIDRequest) error {
	operatorID := contextx.GetCurrentUserID(l.ctx)
	return l.svcCtx.AdminUserService.Delete(l.ctx, operatorID, req.ID)
}
