package dashboard

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPendingOrdersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPendingOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetPendingOrdersLogic {
	return GetPendingOrdersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPendingOrdersLogic) GetPendingOrders(req *types.PendingOrdersRequest) (resp *types.PendingOrdersResponse, err error) {
	helper := NewDashboardHelper(l.ctx, l.svcCtx)
	tenantID, ok := helper.GetTenantID()
	if !ok {
		return nil, code.ErrUnauthorized
	}

	limit := req.Limit
	if limit <= 0 || limit > 50 {
		limit = 5
	}

	return helper.GetPendingOrders(tenantID, limit)
}
