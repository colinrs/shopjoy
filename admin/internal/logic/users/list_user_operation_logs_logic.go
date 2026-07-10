// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package users

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/domain/user"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/tenant"

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
	tenantID, ok := tenant.FromContext(l.ctx)
	if !ok {
		return nil, code.ErrTenantInvalidID
	}

	query := user.OperationLogQuery{
		Page:     req.Page,
		PageSize: req.PageSize,
		Action:   req.Action,
		Keyword:  req.Keyword,
	}

	out, err := l.svcCtx.OperationLogService.List(l.ctx, l.svcCtx.DB, tenantID, req.ID, query)
	if err != nil {
		return nil, err
	}

	list := make([]*types.UserOperationLog, 0, len(out.List))
	for _, item := range out.List {
		list = append(list, &types.UserOperationLog{
			ID:           item.ID,
			UserID:       item.UserID,
			Action:       item.Action,
			ActionText:   item.ActionText,
			OperatorID:   item.OperatorID,
			OperatorName: item.OperatorName,
			Reason:       item.Reason,
			IPAddress:    item.IPAddress,
			UserAgent:    item.UserAgent,
			CreatedAt:    item.CreatedAt,
		})
	}

	return &types.ListUserOperationLogsResp{
		List:     list,
		Total:    out.Total,
		Page:     out.Page,
		PageSize: out.PageSize,
	}, nil
}