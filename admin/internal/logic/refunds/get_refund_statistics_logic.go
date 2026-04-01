package refunds

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRefundStatisticsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRefundStatisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetRefundStatisticsLogic {
	return GetRefundStatisticsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRefundStatisticsLogic) GetRefundStatistics(req *types.GetRefundStatisticsReq) (resp *types.RefundStatisticsResp, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	// Parse time range
	var startTime, endTime time.Time
	if req.StartTime != "" {
		startTime, err = parseTime(req.StartTime)
		if err != nil {
			return nil, err
		}
	}
	if req.EndTime != "" {
		endTime, err = parseTime(req.EndTime)
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

	stats, err := l.svcCtx.RefundApp.GetRefundStatistics(l.ctx, shared.TenantID(tenantID), startTime, endTime)
	if err != nil {
		return nil, err
	}

	// Convert reason breakdown
	reasonBreakdown := make([]types.RefundReasonStats, len(stats.ReasonBreakdown))
	for i, r := range stats.ReasonBreakdown {
		reasonBreakdown[i] = types.RefundReasonStats{
			ReasonType: r.ReasonType,
			ReasonName: r.ReasonName,
			Count:      r.Count,
			Percentage: r.Percentage,
		}
	}

	// Convert daily trend
	dailyTrend := make([]types.RefundDailyStats, len(stats.DailyTrend))
	for i, d := range stats.DailyTrend {
		dailyTrend[i] = types.RefundDailyStats{
			Date:   d.Date,
			Count:  d.Count,
			Amount: d.Amount,
		}
	}

	// Convert top products
	topProducts := make([]types.RefundProductStats, len(stats.TopProducts))
	for i, p := range stats.TopProducts {
		topProducts[i] = types.RefundProductStats{
			ProductID:   p.ProductID,
			ProductName: p.ProductName,
			RefundCount: p.RefundCount,
			RefundRate:  p.RefundRate,
		}
	}

	return &types.RefundStatisticsResp{
		TotalRefunds:    stats.TotalRefunds,
		TotalAmount:     stats.TotalAmount,
		Currency:        stats.Currency,
		RefundRate:      stats.RefundRate,
		PendingCount:    stats.PendingCount,
		ApprovedCount:   stats.ApprovedCount,
		RejectedCount:   stats.RejectedCount,
		CompletedCount:  stats.CompletedCount,
		ReasonBreakdown: reasonBreakdown,
		DailyTrend:      dailyTrend,
		TopProducts:     topProducts,
	}, nil
}