package admin_users

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/application/adminuser"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAssignRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) AssignRolesLogic {
	return AssignRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AssignRolesLogic) AssignRoles(req *types.AssignRolesRequest) error {
	operatorID := contextx.GetCurrentUserID(l.ctx)

	assignReq := adminuser.AssignRolesRequest{
		AdminUserID: req.ID,
		RoleIDs:     req.RoleIDs,
	}

	return l.svcCtx.AdminUserService.AssignRoles(l.ctx, operatorID, assignReq)
}
