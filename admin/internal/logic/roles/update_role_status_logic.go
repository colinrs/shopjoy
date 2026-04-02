package roles

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRoleStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRoleStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateRoleStatusLogic {
	return UpdateRoleStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRoleStatusLogic) UpdateRoleStatus(req *types.UpdateRoleStatusRequest) (resp *types.RoleInfo, err error) {
	// Get tenant ID from context
	tenantIDRaw, _ := contextx.GetTenantID(l.ctx)
	tenantID := shared.TenantID(tenantIDRaw)

	// Get role
	r, err := l.svcCtx.RoleRepo.FindByID(l.ctx, l.svcCtx.DB, tenantID, req.ID)
	if err != nil {
		return nil, err
	}

	// Check if it's a system role
	if r.IsSystem {
		return nil, code.ErrRoleCannotModifySystem
	}

	// Update status
	if req.Status == 1 {
		r.Enable()
	} else {
		r.Disable()
	}
	r.Audit.UpdatedAt = time.Now().UTC()

	// Save role
	if err := l.svcCtx.RoleRepo.Update(l.ctx, l.svcCtx.DB, r); err != nil {
		return nil, err
	}

	return &types.RoleInfo{
		ID:          r.ID,
		Name:        r.Name,
		Code:        r.Code,
		Description: r.Description,
		Status:      int8(r.Status), // #nosec G115 // status values are small (tinyint range)
		StatusText:  getStatusText(r.Status),
		IsSystem:    r.IsSystem,
		CreatedAt:   r.Audit.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   r.Audit.UpdatedAt.Format(time.RFC3339),
	}, nil
}
