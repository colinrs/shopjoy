package dashboard

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
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
func (h *DashboardHelper) GetTenantID() shared.TenantID {
	tenantID, _ := contextx.GetTenantID(h.ctx)
	return shared.TenantID(tenantID)
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
		TodaySales:     formatAmount(todayGMV, "CNY"),
		TodayGrowth:    todayGrowth,
		YesterdaySales: formatAmount(yesterdayGMV, "CNY"),
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
	type statusCount struct {
		Status fulfillment.OrderStatus
		Count  int64
	}

	var results []statusCount
	err := h.svcCtx.DB.Model(&fulfillment.Order{}).
		Where("tenant_id = ?", tenantID).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error
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

	// Query top products by sales quantity from order items
	type productSales struct {
		ProductID   int64
		ProductName string
		Image       string
		Sales       int64
		Revenue     decimal.Decimal
	}

	var results []productSales
	query := h.svcCtx.DB.Table("order_items oi").
		Select("oi.product_id, oi.product_name, oi.image, SUM(oi.quantity) as sales, SUM(oi.total_price) as revenue").
		Joins("JOIN orders o ON o.id = oi.order_id").
		Where("o.tenant_id = ?", tenantID).
		Where("o.status IN ?", []fulfillment.OrderStatus{
			fulfillment.OrderStatusPaid,
			fulfillment.OrderStatusShipped,
			fulfillment.OrderStatusDelivered,
		})

	if !daysAgo.IsZero() {
		query = query.Where("o.paid_at >= ?", daysAgo)
	}

	err := query.Group("oi.product_id, oi.product_name, oi.image").
		Order("sales DESC").
		Limit(limit).
		Scan(&results).Error
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
			Revenue:     formatAmount(r.Revenue, "CNY"),
		})
	}

	return &types.TopProductsResponse{
		List:     list,
		Currency: "CNY",
	}, nil
}

// GetPendingOrders returns recent pending payment orders
func (h *DashboardHelper) GetPendingOrders(tenantID shared.TenantID, limit int) (*types.PendingOrdersResponse, error) {
	var orders []*fulfillment.Order
	err := h.svcCtx.DB.Model(&fulfillment.Order{}).
		Where("tenant_id = ?", tenantID).
		Where("status = ?", fulfillment.OrderStatusPendingPayment).
		Order("created_at DESC").
		Limit(limit).
		Find(&orders).Error
	if err != nil {
		return nil, err
	}

	// Get total count
	var total int64
	h.svcCtx.DB.Model(&fulfillment.Order{}).
		Where("tenant_id = ?", tenantID).
		Where("status = ?", fulfillment.OrderStatusPendingPayment).
		Count(&total)

	list := make([]*types.PendingOrderItem, 0, len(orders))
	for _, o := range orders {
		list = append(list, &types.PendingOrderItem{
			OrderID:    o.ID,
			OrderNo:    o.OrderNo,
			PayAmount:  formatAmount(o.PayAmount, o.Currency),
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
	var recentOrders []*fulfillment.Order
	err := h.svcCtx.DB.Model(&fulfillment.Order{}).
		Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Limit(5).
		Find(&recentOrders).Error
	if err == nil {
		for _, o := range recentOrders {
			activities = append(activities, &types.ActivityItem{
				ID:      o.ID,
				Type:    "order_created",
				Content: fmt.Sprintf("新订单 %s，金额 %s", o.OrderNo, formatAmount(o.PayAmount, o.Currency)),
				Time:    o.Audit.CreatedAt.Format(time.RFC3339),
			})
		}
	}

	// Get recent payments
	var paidOrders []*fulfillment.Order
	err = h.svcCtx.DB.Model(&fulfillment.Order{}).
		Where("tenant_id = ?", tenantID).
		Where("paid_at IS NOT NULL").
		Order("paid_at DESC").
		Limit(5).
		Find(&paidOrders).Error
	if err == nil {
		for _, o := range paidOrders {
			activities = append(activities, &types.ActivityItem{
				ID:      o.ID,
				Type:    "payment_received",
				Content: fmt.Sprintf("订单 %s 已支付 %s", o.OrderNo, formatAmount(o.PayAmount, o.Currency)),
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
	var total decimal.Decimal
	err := h.svcCtx.DB.Model(&fulfillment.Order{}).
		Where("tenant_id = ?", tenantID).
		Where("status IN ?", []fulfillment.OrderStatus{
			fulfillment.OrderStatusPaid,
			fulfillment.OrderStatusShipped,
			fulfillment.OrderStatusDelivered,
		}).
		Where("paid_at >= ? AND paid_at < ?", start, end).
		Select("COALESCE(SUM(pay_amount), 0)").
		Scan(&total).Error
	if err != nil {
		return decimal.Zero, err
	}
	return total, nil
}

func (h *DashboardHelper) getTotalProducts(tenantID shared.TenantID) (int64, error) {
	var count int64
	err := h.svcCtx.DB.Model(&product.Product{}).
		Where("tenant_id = ?", tenantID).
		Where("deleted_at IS NULL").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (h *DashboardHelper) getTotalUsers(tenantID shared.TenantID) (int64, error) {
	var count int64
	err := h.svcCtx.DB.Table("users").
		Where("tenant_id = ?", tenantID).
		Where("deleted_at IS NULL").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (h *DashboardHelper) getNewUsersToday(tenantID shared.TenantID, todayStart time.Time) (int64, error) {
	var count int64
	err := h.svcCtx.DB.Table("users").
		Where("tenant_id = ?", tenantID).
		Where("created_at >= ?", todayStart).
		Where("deleted_at IS NULL").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (h *DashboardHelper) getSalesTrendData(tenantID shared.TenantID, days int) ([]*types.SalesTrendData, error) {
	now := time.Now().UTC()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// Generate date range
	dates := make([]time.Time, days)
	for i := 0; i < days; i++ {
		dates[days-1-i] = todayStart.AddDate(0, 0, -i)
	}

	// Query daily sales
	type dailySales struct {
		Date   time.Time
		Sales  decimal.Decimal
		Orders int64
	}

	var results []dailySales
	err := h.svcCtx.DB.Model(&fulfillment.Order{}).
		Select("DATE(paid_at) as date, SUM(pay_amount) as sales, COUNT(*) as orders").
		Where("tenant_id = ?", tenantID).
		Where("status IN ?", []fulfillment.OrderStatus{
			fulfillment.OrderStatusPaid,
			fulfillment.OrderStatusShipped,
			fulfillment.OrderStatusDelivered,
		}).
		Where("paid_at >= ?", dates[0]).
		Group("DATE(paid_at)").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// Map results by date
	salesMap := make(map[string]dailySales)
	for _, r := range results {
		dateStr := r.Date.Format(time.DateOnly)
		salesMap[dateStr] = r
	}

	// Build response for all dates
	data := make([]*types.SalesTrendData, days)
	for i, d := range dates {
		dateStr := d.Format(time.DateOnly)
		sales := salesMap[dateStr]
		data[i] = &types.SalesTrendData{
			Date:   dateStr,
			Sales:  formatAmount(sales.Sales, "CNY"),
			Orders: sales.Orders,
		}
	}

	return data, nil
}

// formatAmount formats amount to currency string
func formatAmount(amount decimal.Decimal, currency string) string {
	if amount.IsZero() {
		return "0.00"
	}
	return amount.StringFixed(2)
}