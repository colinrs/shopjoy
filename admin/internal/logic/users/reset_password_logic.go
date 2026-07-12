package users

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/domain/user"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) ResetPasswordLogic {
	return ResetPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetPasswordLogic) ResetPassword(req *types.ResetPasswordRequest) (resp *types.ResetPasswordResponse, err error) {
	tempPassword, err := l.svcCtx.UserService.ResetPassword(l.ctx, req.ID)
	if err != nil {
		return nil, err
	}

	recordOperationLog(l.ctx, l.svcCtx, req.ID, user.ActionResetPassword, "")

	return &types.ResetPasswordResponse{
		TemporaryPassword: tempPassword,
	}, nil
}
