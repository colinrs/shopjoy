package admin_users

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/application/adminuser"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAdminProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAdminProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateAdminProfileLogic {
	return UpdateAdminProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAdminProfileLogic) UpdateAdminProfile(req *types.UpdateProfileRequest) (resp *types.AdminUserInfo, err error) {
	userID := contextx.GetCurrentUserID(l.ctx)
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	updateReq := adminuser.UpdateProfileRequest{
		UserID:   userID,
		TenantID: tenantID,
		RealName: req.RealName,
		Avatar:   req.Avatar,
		Mobile:   req.Mobile,
		Email:    req.Email,
	}

	userResp, err := l.svcCtx.AdminUserService.UpdateProfile(l.ctx, updateReq)
	if err != nil {
		return nil, err
	}

	return &types.AdminUserInfo{
		ID:        userResp.ID,
		TenantID:  userResp.TenantID,
		Username:  userResp.Username,
		Email:     userResp.Email,
		Mobile:    userResp.Mobile,
		RealName:  userResp.RealName,
		Avatar:    userResp.Avatar,
		Type:      userResp.Type,
		TypeText:  userResp.TypeText,
		Status:    userResp.Status,
		CreatedAt: userResp.CreatedAt,
	}, nil
}
