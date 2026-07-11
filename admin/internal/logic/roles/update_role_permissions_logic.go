package roles

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRolePermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRolePermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateRolePermissionsLogic {
	return UpdateRolePermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRolePermissionsLogic) UpdateRolePermissions(req *types.UpdateRolePermissionsRequest) error {
	// Get tenant ID from context

	// Verify role exists
	r, err := l.svcCtx.RoleRepo.FindByID(l.ctx, l.svcCtx.DB, req.ID)
	if err != nil {
		return err
	}

	// Check if it's a system role
	if r.IsSystem {
		return code.ErrRoleCannotModifySystem
	}

	// Update permissions
	permissionIDs, err := utils.ParseInt64Slice(req.PermissionIDs)
	if err != nil {
		return code.ErrParam
	}
	return l.svcCtx.PermissionRepo.AssignToRole(l.ctx, l.svcCtx.DB, req.ID, permissionIDs)
}
