package roles

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/role"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateRoleLogic {
	return CreateRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRoleLogic) CreateRole(req *types.CreateRoleRequest) (resp *types.CreateRoleResponse, err error) {
	// Get tenant ID from context
	tenantIDRaw, _ := contextx.GetTenantID(l.ctx)
	tenantID := shared.TenantID(tenantIDRaw)

	// Check if role with same code already exists
	existingRole, err := l.svcCtx.RoleRepo.FindByCode(l.ctx, l.svcCtx.DB, tenantID, req.Code)
	if err == nil && existingRole != nil {
		return nil, code.ErrRoleDuplicate
	}

	// Create new role
	now := time.Now().UTC()
	newRole := &role.Role{
		TenantID:    tenantID,
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      role.StatusEnabled,
		IsSystem:    false,
		Audit: shared.AuditInfo{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// Save role
	if err := l.svcCtx.RoleRepo.Create(l.ctx, l.svcCtx.DB, newRole); err != nil {
		return nil, err
	}

	// Assign permissions if provided
	if len(req.PermissionIDs) > 0 {
		if err := l.svcCtx.PermissionRepo.AssignToRole(l.ctx, l.svcCtx.DB, newRole.ID, req.PermissionIDs); err != nil {
			l.Logger.Errorf("failed to assign permissions: %v", err)
		}
	}

	return &types.CreateRoleResponse{
		ID: newRole.ID,
	}, nil
}