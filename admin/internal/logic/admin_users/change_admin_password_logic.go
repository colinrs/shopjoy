package admin_users

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/application/adminuser"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeAdminPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangeAdminPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) ChangeAdminPasswordLogic {
	return ChangeAdminPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeAdminPasswordLogic) ChangeAdminPassword(req *types.AdminChangePasswordRequest) error {
	userID := contextx.GetCurrentUserID(l.ctx)
	tenantID, _ := contextx.GetTenantID(l.ctx)

	changeReq := adminuser.ChangePasswordRequest{
		UserID:          userID,
		TenantID:        tenantID,
		OldPassword:     req.OldPassword,
		NewPassword:     req.NewPassword,
		ConfirmPassword: req.ConfirmPassword,
	}

	return l.svcCtx.AdminUserService.ChangePassword(l.ctx, changeReq)
}
