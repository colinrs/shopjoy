package users

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/tenant"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetUserLogic {
	return GetUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser(req *types.GetUserRequest) (resp *types.GetUserResponse, err error) {
	tenantID, ok := tenant.FromContext(l.ctx)
	if !ok {
		return nil, code.ErrTenantInvalidID
	}

	userResp, err := l.svcCtx.UserService.GetByID(l.ctx, tenantID, req.ID)
	if err != nil {
		return nil, err
	}

	return &types.GetUserResponse{
		ID:        userResp.ID,
		Email:     userResp.Email,
		Phone:     userResp.Phone,
		Name:      userResp.Name,
		Avatar:    userResp.Avatar,
		Status:    userResp.Status,
		CreatedAt: userResp.CreatedAt,
		LastLogin: userResp.LastLogin,
	}, nil
}
