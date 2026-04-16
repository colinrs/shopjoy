package dashboard

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/utils"
	"github.com/shopspring/decimal"

	"github.com/zeromicro/go-zero/core/logx"
)

// DashboardHelper provides shared logic for dashboard endpoints
type DashboardHelper struct {
	Logger logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDashboardHelper(ctx context.Context, svcCtx *svc.ServiceContext) *DashboardHelper {
	return &DashboardHelper{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetTenantID extracts tenant ID from context
func (h *DashboardHelper) GetTenantID() (shared.TenantID, bool) {
	tenantID, ok := contextx.GetTenantID(h.ctx)
	if !ok && !contextx.IsPlatformAdmin(h.ctx) {
		return 0, false
	}
	if contextx.IsPlatformAdmin(h.ctx) {
		tenantID = 0
	}
	return shared.TenantID(tenantID), true
}

// GetOverview retrieves dashboard overview statistics
func (h *DashboardHelper) GetOverview(tenantID shared.TenantID) (*types.DashboardOverviewResponse, error) {
	now := time.Now().UTC()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	yesterdayStart := todayStart.Add(-24 * time.Hour)
	yesterdayEnd := todayStart

	// Get today's orders count
	todayOrders, err := h.svcCtx.OrderRepo.CountTodayOrders(h.ctx, h.svcCtx.DB, tenantID)
	if err != nil {
		h.Logger.Errorf("failed to get today orders: %v", err)
		todayOrders = 0
	}

	// Get today's GMV
	todayGMV, err := h.svcCtx.OrderRepo.SumTodayGMV(h.ctx, h.svcCtx.DB, tenantID)
	if err != nil {
		h.Logger.Errorf("failed to get today GMV: %v", err)
		todayGMV = decimal.Zero
	}

	// Get yesterday's GMV for comparison
	yesterdayGMV, err := h.getYesterdayGMV(tenantID, yesterdayStart, yesterdayEnd)
	if err != nil {
		h.Logger.Errorf("failed to get yesterday GMV: %v", err)
		yesterdayGMV = decimal.Zero
	}

	// Calculate growth percentage
	todayGrowth := "0%"
	if yesterdayGMV.IsPositive() {
		growth := todayGMV.Sub(yesterdayGMV).Div(yesterdayGMV).Mul(decimal.NewFromFloat(100))
		todayGrowth = fmt.Sprintf("%.1f%%", growth.InexactFloat64())
	} else if todayGMV.IsPositive() {
		todayGrowth = "+100%"
	}

	// Get total products
	totalProducts, err := h.getTotalProducts(tenantID)
	if err != nil {
		h.Logger.Errorf("failed to get total products: %v", err)
		totalProducts = 0
	}

	// Get total users
	totalUsers, err := h.getTotalUsers(tenantID)
	if err != nil {
		h.Logger.Errorf("failed to get total users: %v", err)
		totalUsers = 0
	}

	// Get new users today
	newUsersToday, err := h.getNewUsersToday(tenantID, todayStart)
	if err != nil {
		h.Logger.Errorf("failed to get new users today: %v", err)
		newUsersToday = 0
	}

	return &types.DashboardOverviewResponse{
		TodayOrders:    todayOrders,
		TodaySales:     utils.FormatAmountWithCurrency(todayGMV, "CNY"),
		TodayGrowth:    todayGrowth,
		YesterdaySales: utils.FormatAmountWithCurrency(yesterdayGMV, "CNY"),
		TotalProducts:  totalProducts,
		TotalUsers:     totalUsers,
		NewUsersToday:  newUsersToday,
		Currency:       "CNY",
	}, nil
}

// GetSalesTrend retrieves sales trend data
func (h *DashboardHelper) GetSalesTrend(tenantID shared.TenantID, period string) (*types.SalesTrendResponse, error) {
	var days int
	switch period {
	case "week":
		days = 7
	case "month":
		days = 30
	case "year":
		days = 365
	default:
		days = 7
	}

	data, err := h.getSalesTrendData(tenantID, days)
	if err != nil {
		return nil, err
	}

	return &types.SalesTrendResponse{
		Period:   period,
		Data:     data,
		Currency: "CNY",
	}, nil
}

// GetOrderStatusDistribution returns order counts by status
func (h *DashboardHelper) GetOrderStatusDistribution(tenantID shared.TenantID) (*types.OrderStatusDistributionResponse, error) {
	results, err := h.svcCtx.OrderRepo.CountByStatus(h.ctx, h.svcCtx.DB, tenantID)
	if err != nil {
		return nil, err
	}

	// Calculate total
	var total int64
	for _, r := range results {
		total += r.Count
	}

	// Map to response
	statusColors := map[fulfillment.OrderStatus]string{
		fulfillment.OrderStatusPendingPayment: "#E6A23C",
		fulfillment.OrderStatusPaid:           "#409EFF",
		fulfillment.OrderStatusShipped:        "#67C23A",
		fulfillment.OrderStatusDelivered:      "#67C23A",
		fulfillment.OrderStatusCancelled:      "#909399",
		fulfillment.OrderStatusRefunded:       "#F56C6C",
	}

	list := make([]*types.OrderStatusItem, 0, len(results))
	for _, r := range results {
		percentage := "0%"
		if total > 0 {
			percentage = fmt.Sprintf("%.1f%%", float64(r.Count)/float64(total)*100)
		}
		list = append(list, &types.OrderStatusItem{
			Status:     string(r.Status),
			StatusText: r.Status.Text(),
			Count:      r.Count,
			Percentage: percentage,
			Color:      statusColors[r.Status],
		})
	}

	return &types.OrderStatusDistributionResponse{
		List:  list,
		Total: total,
	}, nil
}

// GetTopProducts returns top selling products
func (h *DashboardHelper) GetTopProducts(tenantID shared.TenantID, limit int, period string) (*types.TopProductsResponse, error) {
	var daysAgo time.Time
	switch period {
	case "week":
		daysAgo = time.Now().UTC().AddDate(0, 0, -7)
	case "month":
		daysAgo = time.Now().UTC().AddDate(0, -1, 0)
	case "all":
		daysAgo = time.Time{} // Zero time means all time
	default:
		daysAgo = time.Now().UTC().AddDate(0, 0, -7)
	}

	results, err := h.svcCtx.OrderRepo.FindTopProducts(h.ctx, h.svcCtx.DB, tenantID, daysAgo, limit)
	if err != nil {
		return nil, err
	}

	list := make([]*types.TopProductItem, 0, len(results))
	for _, r := range results {
		list = append(list, &types.TopProductItem{
			ProductID:   r.ProductID,
			ProductName: r.ProductName,
			Image:       r.Image,
			Sales:       r.Sales,
			Revenue:     utils.FormatAmountWithCurrency(r.Revenue, "CNY"),
		})
	}

	return &types.TopProductsResponse{
		List:     list,
		Currency: "CNY",
	}, nil
}

// GetPendingOrders returns recent pending payment orders
func (h *DashboardHelper) GetPendingOrders(tenantID shared.TenantID, limit int) (*types.PendingOrdersResponse, error) {
	orders, err := h.svcCtx.OrderRepo.FindPendingOrders(h.ctx, h.svcCtx.DB, tenantID, limit)
	if err != nil {
		return nil, err
	}

	total, err := h.svcCtx.OrderRepo.CountPendingOrders(h.ctx, h.svcCtx.DB, tenantID)
	if err != nil {
		return nil, err
	}

	list := make([]*types.PendingOrderItem, 0, len(orders))
	for _, o := range orders {
		list = append(list, &types.PendingOrderItem{
			OrderID:    o.ID,
			OrderNo:    o.OrderNo,
			PayAmount:  utils.FormatAmountWithCurrency(o.PayAmount, o.Currency),
			Status:     string(o.Status),
			StatusText: o.Status.Text(),
			CreatedAt:  o.Audit.CreatedAt.Format(time.RFC3339),
		})
	}

	return &types.PendingOrdersResponse{
		List:  list,
		Total: total,
	}, nil
}

// GetRecentActivities returns recent system activities
func (h *DashboardHelper) GetRecentActivities(tenantID shared.TenantID, limit int) (*types.RecentActivitiesResponse, error) {
	activities := make([]*types.ActivityItem, 0, limit)

	// Get recent orders
	recentOrders, err := h.svcCtx.OrderRepo.FindRecentOrders(h.ctx, h.svcCtx.DB, tenantID, 5)
	if err == nil {
		for _, o := range recentOrders {
			activities = append(activities, &types.ActivityItem{
				ID:      o.ID,
				Type:    "order_created",
				Content: fmt.Sprintf("新订单 %s，金额 %s", o.OrderNo, utils.FormatAmountWithCurrency(o.PayAmount, o.Currency)),
				Time:    o.Audit.CreatedAt.Format(time.RFC3339),
			})
		}
	}

	// Get recent payments
	paidOrders, err := h.svcCtx.OrderRepo.FindRecentPaidOrders(h.ctx, h.svcCtx.DB, tenantID, 5)
	if err == nil {
		for _, o := range paidOrders {
			activities = append(activities, &types.ActivityItem{
				ID:      o.ID,
				Type:    "payment_received",
				Content: fmt.Sprintf("订单 %s 已支付 %s", o.OrderNo, utils.FormatAmountWithCurrency(o.PayAmount, o.Currency)),
				Time:    o.PaidAt.Format(time.RFC3339),
			})
		}
	}

	// Sort by time (most recent first) and limit
	sort.Slice(activities, func(i, j int) bool {
		return activities[i].Time > activities[j].Time
	})

	if len(activities) > limit {
		activities = activities[:limit]
	}

	return &types.RecentActivitiesResponse{
		List: activities,
	}, nil
}

// Helper methods

func (h *DashboardHelper) getYesterdayGMV(tenantID shared.TenantID, start, end time.Time) (decimal.Decimal, error) {
	statuses := []fulfillment.OrderStatus{
		fulfillment.OrderStatusPaid,
		fulfillment.OrderStatusShipped,
		fulfillment.OrderStatusDelivered,
	}
	return h.svcCtx.OrderRepo.SumGMVByDateRange(h.ctx, h.svcCtx.DB, tenantID, start, end, statuses)
}

func (h *DashboardHelper) getTotalProducts(tenantID shared.TenantID) (int64, error) {
	return h.svcCtx.ProductRepo.CountTotal(h.ctx, h.svcCtx.DB, tenantID)
}

func (h *DashboardHelper) getTotalUsers(tenantID shared.TenantID) (int64, error) {
	stats, err := h.svcCtx.UserService.GetStats(h.ctx, tenantID)
	if err != nil {
		return 0, err
	}
	return stats.Total, nil
}

func (h *DashboardHelper) getNewUsersToday(tenantID shared.TenantID, todayStart time.Time) (int64, error) {
	stats, err := h.svcCtx.UserService.GetStats(h.ctx, tenantID)
	if err != nil {
		return 0, err
	}
	return stats.NewToday, nil
}

func (h *DashboardHelper) getSalesTrendData(tenantID shared.TenantID, days int) ([]*types.SalesTrendData, error) {
	now := time.Now().UTC()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endDate := todayStart.Add(24 * time.Hour)

	// Generate date range
	dates := make([]time.Time, days)
	for i := 0; i < days; i++ {
		dates[days-1-i] = todayStart.AddDate(0, 0, -i)
	}

	// Query daily sales
	results, err := h.svcCtx.OrderRepo.FindSalesTrend(h.ctx, h.svcCtx.DB, tenantID, dates[0], endDate)
	if err != nil {
		return nil, err
	}

	// Map results by date
	salesMap := make(map[string]*fulfillment.DailySalesTrend)
	for _, r := range results {
		salesMap[r.Date] = r
	}

	// Build response for all dates
	data := make([]*types.SalesTrendData, days)
	for i, d := range dates {
		dateStr := d.Format(time.DateOnly)
		sales := salesMap[dateStr]
		salesAmount := decimal.Zero
		orders := int64(0)
		if sales != nil {
			salesAmount = sales.Sales
			orders = sales.Orders
		}
		data[i] = &types.SalesTrendData{
			Date:   dateStr,
			Sales:  utils.FormatAmountWithCurrency(salesAmount, "CNY"),
			Orders: orders,
		}
	}

	return data, nil
}
