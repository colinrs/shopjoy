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

type ListRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListRolesLogic {
	return ListRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRolesLogic) ListRoles(req *types.ListRolesRequest) (resp *types.ListRolesResponse, err error) {
	// Get tenant ID from context
	tenantIDRaw, _ := contextx.GetTenantID(l.ctx)
	tenantID := shared.TenantID(tenantIDRaw)

	// Build query
	query := role.Query{
		PageQuery: shared.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		Name: req.Name,
		Code: req.Code,
	}

	// Get roles
	roles, total, err := l.svcCtx.RoleRepo.FindList(l.ctx, l.svcCtx.DB, tenantID, query)
	if err != nil {
		return nil, err
	}

	// Convert to response
	list := make([]*types.RoleInfo, 0, len(roles))
	for _, r := range roles {
		list = append(list, &types.RoleInfo{
			ID:          r.ID,
			Name:        r.Name,
			Code:        r.Code,
			Description: r.Description,
			Status:      int8(r.Status),
			StatusText:  getStatusText(r.Status),
			IsSystem:    r.IsSystem,
			CreatedAt:   r.Audit.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   r.Audit.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &types.ListRolesResponse{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func getStatusText(status role.Status) string {
	if status == role.StatusEnabled {
		return "启用"
	}
	return "禁用"
}
