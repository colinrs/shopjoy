package users

import (
	"context"

	appUser "github.com/colinrs/shopjoy/admin/internal/application/user"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/tenant"

	"github.com/zeromicro/go-zero/core/logx"
)

type SuspendUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSuspendUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) SuspendUserLogic {
	return SuspendUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SuspendUserLogic) SuspendUser(req *types.SuspendUserRequest) (resp *types.GetUserResponse, err error) {
	tenantID, ok := tenant.FromContext(l.ctx)
	if !ok {
		return nil, code.ErrTenantInvalidID
	}

	if err := l.svcCtx.UserService.Suspend(l.ctx, tenantID, req.ID); err != nil {
		return nil, err
	}

	userResp, err := l.svcCtx.UserService.GetByID(l.ctx, tenantID, req.ID)
	if err != nil {
		return nil, err
	}

	return toGetUserResponse(userResp), nil
}

func toGetUserResponse(u *appUser.UserResponse) *types.GetUserResponse {
	return &types.GetUserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Phone:     u.Phone,
		Name:      u.Name,
		Avatar:    u.Avatar,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		LastLogin: u.LastLogin,
	}
}
