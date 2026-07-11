package admin_users

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetAdminPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetAdminPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) ResetAdminPasswordLogic {
	return ResetAdminPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetAdminPasswordLogic) ResetAdminPassword(req *types.AdminUserIDRequest) (resp *types.ResetAdminPasswordResponse, err error) {
	operatorID := contextx.GetCurrentUserID(l.ctx)

	result, err := l.svcCtx.AdminUserService.ResetPassword(l.ctx, operatorID, req.ID)
	if err != nil {
		return nil, err
	}

	return &types.ResetAdminPasswordResponse{
		TemporaryPassword: result.TemporaryPassword,
	}, nil
}
