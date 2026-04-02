package dashboard

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDashboardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDashboardLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetDashboardLogic {
	return GetDashboardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDashboardLogic) GetDashboard(req *types.GetDashboardRequest) (resp *types.GetDashboardResponse, err error) {
	helper := NewDashboardHelper(l.ctx, l.svcCtx)
	tenantID := helper.GetTenantID()

	// Get overview
	overview, err := helper.GetOverview(tenantID)
	if err != nil {
		l.Logger.Errorf("failed to get overview: %v", err)
		overview = &types.DashboardOverviewResponse{}
	}

	// Get status distribution
	statusDistribution, err := helper.GetOrderStatusDistribution(tenantID)
	if err != nil {
		l.Logger.Errorf("failed to get status distribution: %v", err)
		statusDistribution = &types.OrderStatusDistributionResponse{}
	}

	// Get pending orders
	pendingOrdersResp, err := helper.GetPendingOrders(tenantID, 5)
	if err != nil {
		l.Logger.Errorf("failed to get pending orders: %v", err)
		pendingOrdersResp = &types.PendingOrdersResponse{}
	}

	// Get top products
	topProductsResp, err := helper.GetTopProducts(tenantID, 5, "week")
	if err != nil {
		l.Logger.Errorf("failed to get top products: %v", err)
		topProductsResp = &types.TopProductsResponse{}
	}

	// Get recent activities
	activitiesResp, err := helper.GetRecentActivities(tenantID, 10)
	if err != nil {
		l.Logger.Errorf("failed to get recent activities: %v", err)
		activitiesResp = &types.RecentActivitiesResponse{}
	}

	return &types.GetDashboardResponse{
		Overview:           overview,
		StatusDistribution: statusDistribution,
		PendingOrders:      pendingOrdersResp.List,
		TopProducts:        topProductsResp.List,
		RecentActivities:   activitiesResp.List,
	}, nil
}
