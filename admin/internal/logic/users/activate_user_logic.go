package users

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/domain/user"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

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
	if err := l.svcCtx.UserService.Activate(l.ctx, req.ID); err != nil {
		return nil, err
	}

	recordOperationLog(l.ctx, l.svcCtx, req.ID, user.ActionActivateUser, "")

	userResp, err := l.svcCtx.UserService.GetByID(l.ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return toGetUserResponse(userResp), nil
}
