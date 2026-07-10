// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package users

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserOperationLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户操作日志
func NewListUserOperationLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserOperationLogsLogic {
	return &ListUserOperationLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUserOperationLogsLogic) ListUserOperationLogs(req *types.ListUserOperationLogsReq) (resp *types.ListUserOperationLogsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
