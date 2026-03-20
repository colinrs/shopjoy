package auth

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/application/adminuser"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterTenantAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterTenantAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) RegisterTenantAdminLogic {
	return RegisterTenantAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterTenantAdminLogic) RegisterTenantAdmin(req *types.RegisterTenantAdminRequest) (resp *types.RegisterTenantAdminResponse, err error) {
	registerReq := adminuser.RegisterRequest{
		Email:    req.Email,
		Mobile:   req.Mobile,
		Password: req.Password,
		RealName: req.RealName,
	}

	loginResp, err := l.svcCtx.AdminUserService.RegisterTenantAdmin(l.ctx, registerReq)
	if err != nil {
		return nil, err
	}

	return &types.RegisterTenantAdminResponse{
		AccessToken:  loginResp.AccessToken,
		RefreshToken: loginResp.RefreshToken,
		ExpiresIn:    loginResp.ExpiresIn,
		User: types.AdminUserInfo{
			ID:        loginResp.User.ID,
			TenantID:  loginResp.User.TenantID,
			Username:  loginResp.User.Username,
			Email:     loginResp.User.Email,
			Mobile:    loginResp.User.Mobile,
			RealName:  loginResp.User.RealName,
			Avatar:    loginResp.User.Avatar,
			Type:      loginResp.User.Type,
			TypeText:  loginResp.User.TypeText,
			Status:    loginResp.User.Status,
			CreatedAt: loginResp.User.CreatedAt,
		},
	}, nil
}