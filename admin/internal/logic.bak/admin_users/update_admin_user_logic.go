package admin_users

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/application/adminuser"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAdminUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAdminUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateAdminUserLogic {
	return UpdateAdminUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAdminUserLogic) UpdateAdminUser(req *types.UpdateAdminUserRequest) (resp *types.AdminUserInfo, err error) {
	operatorID := contextx.GetCurrentUserID(l.ctx)

	updateReq := adminuser.UpdateAdminUserRequest{
		ID:       req.ID,
		RealName: req.RealName,
		Avatar:   req.Avatar,
		Mobile:   req.Mobile,
		Email:    req.Email,
	}

	user, err := l.svcCtx.AdminUserService.Update(l.ctx, operatorID, updateReq)
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
