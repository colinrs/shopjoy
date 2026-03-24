package payments

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaymentStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPaymentStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetPaymentStatsLogic {
	return GetPaymentStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPaymentStatsLogic) GetPaymentStats(req *types.GetPaymentStatsReq) (resp *types.PaymentStatsResp, err error) {
	// Get tenant ID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)
	if tenantID == 0 {
		tenantID = 0
	}

	// Get payment stats from service
	stats, err := l.svcCtx.PaymentService.GetPaymentStats(l.ctx, shared.TenantID(tenantID), req.Period)
	if err != nil {
		return nil, err
	}

	// Convert channel distribution
	channelDistribution := make([]types.ChannelDistributionResp, len(stats.ChannelDistribution))
	for i, cd := range stats.ChannelDistribution {
		channelDistribution[i] = types.ChannelDistributionResp{
			Name:    cd.Name,
			Percent: formatPercentInt(cd.Percent),
			Amount:  cd.Amount,
			Count:   cd.Count,
			Color:   cd.Color,
		}
	}

	return &types.PaymentStatsResp{
		TodayReceived:       stats.TodayReceived,
		TodayGrowth:         stats.TodayGrowth,
		PeriodReceived:      stats.PeriodReceived,
		RefundAmount:        stats.RefundAmount,
		RefundRate:          stats.RefundRate,
		Currency:            stats.Currency,
		ChannelDistribution: channelDistribution,
	}, nil
}

// formatPercentInt formats an int percentage to string
func formatPercentInt(p int) string {
	if p == 0 {
		return "0%"
	}
	return formatIntToString(p) + "%"
}

// formatIntToString converts int to string
func formatIntToString(i int) string {
	if i == 0 {
		return "0"
	}
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	var digits []byte
	for i > 0 {
		digits = append([]byte{byte('0' + i%10)}, digits...)
		i /= 10
	}
	if neg {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}