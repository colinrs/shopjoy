package admin_users

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/application/adminuser"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAdminUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateAdminUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateAdminUserLogic {
	return CreateAdminUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAdminUserLogic) CreateAdminUser(req *types.CreateAdminUserRequest) (resp *types.AdminUserInfo, err error) {
	operatorID := contextx.GetCurrentUserID(l.ctx)

	createReq := adminuser.CreateAdminUserRequest{
		TenantID: req.TenantID,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Mobile:   req.Mobile,
		RealName: req.RealName,
		Avatar:   req.Avatar,
		Type:     req.Type,
		RoleIDs:  req.RoleIDs,
	}

	user, err := l.svcCtx.AdminUserService.Create(l.ctx, operatorID, createReq)
	if err != nil {
		return nil, err
	}

	return &types.AdminUserInfo{
		ID:        user.ID,
		TenantID:  user.TenantID,
		Username:  user.Username,
		Email:     user.Email,
		Mobile:    user.Mobile,
		RealName:  user.RealName,
		Avatar:    user.Avatar,
		Type:      user.Type,
		TypeText:  user.TypeText,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}, nil
}
