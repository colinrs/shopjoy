package users

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/tenant"

	"github.com/zeromicro/go-zero/core/logx"
)

type ActivateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewActivateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) ActivateUserLogic {
	return ActivateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ActivateUserLogic) ActivateUser(req *types.ActivateUserRequest) (resp *types.GetUserResponse, err error) {
	tenantID, ok := tenant.FromContext(l.ctx)
	if !ok {
		return nil, code.ErrTenantInvalidID
	}

	if err := l.svcCtx.UserService.Activate(l.ctx, tenantID, req.ID); err != nil {
		return nil, err
	}

	userResp, err := l.svcCtx.UserService.GetByID(l.ctx, tenantID, req.ID)
	if err != nil {
		return nil, err
	}

	return toGetUserResponse(userResp), nil
}
