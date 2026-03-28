package roles

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/role"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetRoleLogic {
	return GetRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRoleLogic) GetRole(req *types.RoleIDRequest) (resp *types.RoleWithPermissions, err error) {
	// Get tenant ID from context
	tenantIDRaw, _ := contextx.GetTenantID(l.ctx)
	tenantID := shared.TenantID(tenantIDRaw)

	// Get role
	r, err := l.svcCtx.RoleRepo.FindByID(l.ctx, l.svcCtx.DB, tenantID, req.ID)
	if err != nil {
		return nil, err
	}

	// Get role permissions
	permissions, err := l.svcCtx.PermissionRepo.FindByRoleIDs(l.ctx, l.svcCtx.DB, []int64{r.ID})
	if err != nil {
		l.Logger.Errorf("failed to get role permissions: %v", err)
		permissions = []*role.Permission{}
	}

	// Convert permissions
	permList := make([]*types.PermissionInfo, 0, len(permissions))
	for _, p := range permissions {
		permList = append(permList, &types.PermissionInfo{
			ID:       p.ID,
			Name:     p.Name,
			Code:     p.Code,
			Type:     int8(p.Type),
			TypeText: getPermissionTypeText(int8(p.Type)),
			ParentID: p.ParentID,
			Path:     p.Path,
			Icon:     p.Icon,
			Sort:     p.Sort,
		})
	}

	return &types.RoleWithPermissions{
		RoleInfo: types.RoleInfo{
			ID:          r.ID,
			Name:        r.Name,
			Code:        r.Code,
			Description: r.Description,
			Status:      int8(r.Status),
			StatusText:  getStatusText(r.Status),
			IsSystem:    r.IsSystem,
			CreatedAt:   r.Audit.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   r.Audit.UpdatedAt.Format(time.RFC3339),
		},
		Permissions: permList,
	}, nil
}

func getPermissionTypeText(permType int8) string {
	switch permType {
	case 0:
		return "菜单"
	case 1:
		return "按钮"
	case 2:
		return "API"
	default:
		return "未知"
	}
}