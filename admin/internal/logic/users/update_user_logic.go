package users

import (
	"context"

	appUser "github.com/colinrs/shopjoy/admin/internal/application/user"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/tenant"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateUserLogic {
	return UpdateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLogic) UpdateUser(req *types.UpdateUserRequest) (resp *types.GetUserResponse, err error) {
	tenantID, ok := tenant.FromContext(l.ctx)
	if !ok {
		tenantID = shared.TenantID(1) // 默认租户
	}

	updateReq := appUser.UpdateUserRequest{
		ID:       req.ID,
		TenantID: tenantID,
		Name:     req.Name,
		Avatar:   req.Avatar,
	}

	userResp, err := l.svcCtx.UserService.Update(l.ctx, updateReq)
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
