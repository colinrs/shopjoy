package fulfillment_statistics

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"time"

	appfulfillment "github.com/colinrs/shopjoy/admin/internal/application/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExportFulfillmentStatisticsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
	w      http.ResponseWriter
}

func NewExportFulfillmentStatisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext, w http.ResponseWriter, r *http.Request) ExportFulfillmentStatisticsLogic {
	return ExportFulfillmentStatisticsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		w:      w,
		r:      r,
	}
}

func (l *ExportFulfillmentStatisticsLogic) ExportFulfillmentStatistics(req *types.ExportFulfillmentStatisticsReq) error {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Parse time range
	var startTime, endTime time.Time
	if req.StartDate != "" {
		startTime, err := time.Parse(time.RFC3339, req.StartDate)
		if err != nil {
			return err
		}
		l.ctx = context.WithValue(l.ctx, "startTime", startTime)
	}
	if req.EndDate != "" {
		endTime, err := time.Parse(time.RFC3339, req.EndDate)
		if err != nil {
			return err
		}
		l.ctx = context.WithValue(l.ctx, "endTime", endTime)
	}

	// If no time range specified, calculate based on period
	if startTime.IsZero() || endTime.IsZero() {
		endTime = time.Now().UTC()
		switch req.Period {
		case "daily":
			startTime = endTime.AddDate(0, 0, -1)
		case "weekly":
			startTime = endTime.AddDate(0, 0, -7)
		case "monthly":
			startTime = endTime.AddDate(0, -1, 0)
		default: // "weekly" as default
			startTime = endTime.AddDate(0, 0, -7)
		}
	}

	// Get refund statistics
	refundStats, err := l.svcCtx.RefundApp.GetRefundStatistics(l.ctx, shared.TenantID(tenantID), startTime, endTime)
	if err != nil {
		return err
	}

	// Get fulfillment summary (shipment counts)
	fulfillmentSummary, err := l.svcCtx.OrderFulfillmentApp.GetFulfillmentSummary(l.ctx, shared.TenantID(tenantID))
	if err != nil {
		// Fallback to zeros if error
		fulfillmentSummary = &appfulfillment.FulfillmentSummary{}
	}

	// Set response headers for CSV download
	filename := fmt.Sprintf("fulfillment_statistics_export_%s.csv", time.Now().Format("20060102"))
	l.w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	l.w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// Write UTF-8 BOM for Excel compatibility
	if _, err := l.w.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
		return err
	}

	// Create CSV writer
	writer := csv.NewWriter(l.w)
	defer writer.Flush()

	// Write Overview Section
	if err := writer.Write([]string{"=== Fulfillment Overview ==="}); err != nil {
		return err
	}

	overviewHeader := []string{
		"Total Shipments", "Pending Shipments", "Total Refunds", "Pending Refunds",
		"Refund Rate", "Delivery Success Rate", "Total Amount", "Currency",
	}
	if err := writer.Write(overviewHeader); err != nil {
		return err
	}

	var deliverySuccessRate string
	if fulfillmentSummary.TotalOrders > 0 {
		rate := float64(fulfillmentSummary.Delivered) / float64(fulfillmentSummary.TotalOrders) * 100
		deliverySuccessRate = fmt.Sprintf("%.1f%%", rate)
	} else {
		deliverySuccessRate = "0.0%"
	}

	overviewRow := []string{
		fmt.Sprintf("%d", fulfillmentSummary.TotalOrders),
		fmt.Sprintf("%d", fulfillmentSummary.PendingShipment),
		fmt.Sprintf("%d", refundStats.TotalRefunds),
		fmt.Sprintf("%d", refundStats.PendingCount),
		refundStats.RefundRate + "%",
		deliverySuccessRate,
		refundStats.TotalAmount,
		refundStats.Currency,
	}
	if err := writer.Write(overviewRow); err != nil {
		return err
	}

	// Empty row for separation
	if err := writer.Write([]string{}); err != nil {
		return err
	}

	// Write Fulfillment Summary Section
	if err := writer.Write([]string{"=== Fulfillment Summary ==="}); err != nil {
		return err
	}

	summaryHeader := []string{
		"Pending Shipment", "Partial Shipped", "Shipped", "Delivered",
		"Pending Refund", "Refunding", "Total Orders", "Today Orders",
	}
	if err := writer.Write(summaryHeader); err != nil {
		return err
	}

	summaryRow := []string{
		fmt.Sprintf("%d", fulfillmentSummary.PendingShipment),
		fmt.Sprintf("%d", fulfillmentSummary.PartialShipped),
		fmt.Sprintf("%d", fulfillmentSummary.Shipped),
		fmt.Sprintf("%d", fulfillmentSummary.Delivered),
		fmt.Sprintf("%d", fulfillmentSummary.PendingRefund),
		fmt.Sprintf("%d", fulfillmentSummary.Refunding),
		fmt.Sprintf("%d", fulfillmentSummary.TotalOrders),
		fmt.Sprintf("%d", fulfillmentSummary.TodayOrders),
	}
	if err := writer.Write(summaryRow); err != nil {
		return err
	}

	// Empty row for separation
	if err := writer.Write([]string{}); err != nil {
		return err
	}

	// Write Refund Rate Trend Section
	if len(refundStats.DailyTrend) > 0 {
		if err := writer.Write([]string{"=== Refund Rate Trend ==="}); err != nil {
			return err
		}

		trendHeader := []string{"Date", "Count", "Amount"}
		if err := writer.Write(trendHeader); err != nil {
			return err
		}

		for _, d := range refundStats.DailyTrend {
			trendRow := []string{d.Date, fmt.Sprintf("%d", d.Count), d.Amount}
			if err := writer.Write(trendRow); err != nil {
				return err
			}
		}

		// Empty row for separation
		if err := writer.Write([]string{}); err != nil {
			return err
		}
	}

	// Write Refund Reasons Section
	if len(refundStats.ReasonBreakdown) > 0 {
		if err := writer.Write([]string{"=== Refund Reasons ==="}); err != nil {
			return err
		}

		reasonHeader := []string{"Reason Type", "Reason Name", "Count", "Percentage"}
		if err := writer.Write(reasonHeader); err != nil {
			return err
		}

		for _, r := range refundStats.ReasonBreakdown {
			reasonRow := []string{r.ReasonType, r.ReasonName, fmt.Sprintf("%d", r.Count), r.Percentage + "%"}
			if err := writer.Write(reasonRow); err != nil {
				return err
			}
		}

		// Empty row for separation
		if err := writer.Write([]string{}); err != nil {
			return err
		}
	}

	// Write Top Products Section (Problem Products)
	if len(refundStats.TopProducts) > 0 {
		if err := writer.Write([]string{"=== Top Refunded Products ==="}); err != nil {
			return err
		}

		productHeader := []string{"Product ID", "Product Name", "Refund Count", "Refund Rate"}
		if err := writer.Write(productHeader); err != nil {
			return err
		}

		for _, p := range refundStats.TopProducts {
			productRow := []string{
				fmt.Sprintf("%d", p.ProductID),
				p.ProductName,
				fmt.Sprintf("%d", p.RefundCount),
				p.RefundRate + "%",
			}
			if err := writer.Write(productRow); err != nil {
				return err
			}
		}
	}

	return nil
}
