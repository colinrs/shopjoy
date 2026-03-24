package admin_users

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAdminUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAdminUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetAdminUserLogic {
	return GetAdminUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAdminUserLogic) GetAdminUser(req *types.AdminUserIDRequest) (resp *types.AdminUserInfo, err error) {
	user, err := l.svcCtx.AdminUserService.GetByID(l.ctx, req.ID)
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
