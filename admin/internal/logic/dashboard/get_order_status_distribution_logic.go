package dashboard

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderStatusDistributionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderStatusDistributionLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetOrderStatusDistributionLogic {
	return GetOrderStatusDistributionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderStatusDistributionLogic) GetOrderStatusDistribution(req *types.OrderStatusDistributionRequest) (resp *types.OrderStatusDistributionResponse, err error) {
	helper := NewDashboardHelper(l.ctx, l.svcCtx)
	tenantID := helper.GetTenantID()
	return helper.GetOrderStatusDistribution(tenantID)
}
