package roles

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

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
	tenantIDRaw, _ := contextx.GetTenantID(l.ctx)
	tenantID := shared.TenantID(tenantIDRaw)

	// Verify role exists
	r, err := l.svcCtx.RoleRepo.FindByID(l.ctx, l.svcCtx.DB, tenantID, req.ID)
	if err != nil {
		return err
	}

	// Check if it's a system role
	if r.IsSystem {
		return code.ErrRoleCannotModifySystem
	}

	// Update permissions
	return l.svcCtx.PermissionRepo.AssignToRole(l.ctx, l.svcCtx.DB, req.ID, req.PermissionIDs)
}