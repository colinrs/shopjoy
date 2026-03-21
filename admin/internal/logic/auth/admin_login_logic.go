package auth

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/application/adminuser"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) AdminLoginLogic {
	return AdminLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminLoginLogic) AdminLogin(req *types.AdminLoginRequest) (resp *types.AdminLoginResponse, err error) {
	loginReq := adminuser.LoginRequest{
		Account:  req.Account,
		Password: req.Password,
		IP:       req.IP,
	}

	loginResp, err := l.svcCtx.AdminUserService.Login(l.ctx, loginReq)
	if err != nil {
		return nil, err
	}

	return &types.AdminLoginResponse{
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
