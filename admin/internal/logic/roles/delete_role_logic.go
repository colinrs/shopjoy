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

type DeleteRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteRoleLogic {
	return DeleteRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteRoleLogic) DeleteRole(req *types.RoleIDRequest) error {
	// Get tenant ID from context
	tenantIDRaw, _ := contextx.GetTenantID(l.ctx)
	tenantID := shared.TenantID(tenantIDRaw)

	// Get role
	r, err := l.svcCtx.RoleRepo.FindByID(l.ctx, l.svcCtx.DB, tenantID, req.ID)
	if err != nil {
		return err
	}

	// Check if it's a system role
	if r.IsSystem {
		return code.ErrRoleCannotModifySystem
	}

	// TODO: Check if role is in use by any users
	// For now, we'll allow deletion

	// Delete role
	return l.svcCtx.RoleRepo.Delete(l.ctx, l.svcCtx.DB, tenantID, req.ID)
}