package fulfillment_statistics

import (
	"context"
	"fmt"
	"time"

	appfulfillment "github.com/colinrs/shopjoy/admin/internal/application/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFulfillmentStatisticsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFulfillmentStatisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetFulfillmentStatisticsLogic {
	return GetFulfillmentStatisticsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFulfillmentStatisticsLogic) GetFulfillmentStatistics(req *types.GetRefundStatisticsReq) (resp *types.FulfillmentStatisticsResp, err error) {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Parse time range
	var startTime, endTime time.Time
	if req.StartTime != "" {
		startTime, err = time.Parse(time.RFC3339, req.StartTime)
		if err != nil {
			return nil, err
		}
	}
	if req.EndTime != "" {
		endTime, err = time.Parse(time.RFC3339, req.EndTime)
		if err != nil {
			return nil, err
		}
	}

	// If no time range specified, calculate based on period
	if startTime.IsZero() || endTime.IsZero() {
		endTime = time.Now().UTC()
		switch req.Period {
		case "7d":
			startTime = endTime.AddDate(0, 0, -7)
		case "90d":
			startTime = endTime.AddDate(0, 0, -90)
		default: // "30d"
			startTime = endTime.AddDate(0, 0, -30)
		}
	}

	// Get refund statistics
	refundStats, err := l.svcCtx.RefundApp.GetRefundStatistics(l.ctx, shared.TenantID(tenantID), startTime, endTime)
	if err != nil {
		return nil, err
	}

	// Get fulfillment summary (shipment counts)
	fulfillmentSummary, err := l.svcCtx.OrderFulfillmentApp.GetFulfillmentSummary(l.ctx, shared.TenantID(tenantID))
	if err != nil {
		// Fallback to zeros if error
		fulfillmentSummary = &appfulfillment.FulfillmentSummary{}
	}

	// Calculate total shipments and delivery success rate
	pendingShipments := fulfillmentSummary.PendingShipment
	deliveredShipments := fulfillmentSummary.Delivered
	totalShipments := fulfillmentSummary.TotalOrders // Total orders as proxy

	var deliverySuccessRate string
	if totalShipments > 0 {
		rate := float64(deliveredShipments) / float64(totalShipments) * 100
		deliverySuccessRate = fmt.Sprintf("%.1f", rate)
	} else {
		deliverySuccessRate = "0.0"
	}

	// Build overview
	overview := &types.FulfillmentOverview{
		TotalShipments:      totalShipments,
		PendingShipments:    pendingShipments,
		TotalRefunds:        refundStats.TotalRefunds,
		PendingRefunds:      refundStats.PendingCount,
		RefundRate:          refundStats.RefundRate,
		DeliverySuccessRate: deliverySuccessRate,
		RefundAmount:        refundStats.TotalAmount,
		Currency:            refundStats.Currency,
	}

	// Convert refund rate trend
	refundRateTrend := make([]*types.RefundRateTrend, len(refundStats.DailyTrend))
	for i, d := range refundStats.DailyTrend {
		refundRateTrend[i] = &types.RefundRateTrend{
			Date: d.Date,
			Rate: d.Amount,
		}
	}

	// Convert refund reasons
	refundReasons := make([]*types.RefundReasonSummary, len(refundStats.ReasonBreakdown))
	for i, r := range refundStats.ReasonBreakdown {
		refundReasons[i] = &types.RefundReasonSummary{
			ReasonType: r.ReasonType,
			ReasonName: r.ReasonName,
			Count:      r.Count,
			Percentage: r.Percentage,
		}
	}

	// Convert problem products
	problemProducts := make([]*types.ProblemProductStats, len(refundStats.TopProducts))
	for i, p := range refundStats.TopProducts {
		problemProducts[i] = &types.ProblemProductStats{
			ProductID:   p.ProductID,
			ProductName: p.ProductName,
			Image:       "",
			TotalSales:  0,
			RefundCount: p.RefundCount,
			RefundRate:  p.RefundRate,
		}
	}

	// Carrier performance - not available in current implementation
	var carrierPerformance []*types.CarrierPerformanceStats

	return &types.FulfillmentStatisticsResp{
		Overview:           overview,
		RefundRateTrend:     refundRateTrend,
		RefundReasons:       refundReasons,
		ProblemProducts:     problemProducts,
		CarrierPerformance:  carrierPerformance,
	}, nil
}
