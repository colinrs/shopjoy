package stats

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPointsStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPointsStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetPointsStatsLogic {
	return GetPointsStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPointsStatsLogic) GetPointsStats(req *types.GetPointsStatsReq) (resp *types.PointsStats, err error) {

	// Calculate time range based on period
	var startTime, endTime *time.Time
	now := time.Now().UTC()
	endTime = &now

	switch req.Period {
	case "7d":
		t := now.AddDate(0, 0, -7)
		startTime = &t
	case "30d":
		t := now.AddDate(0, 0, -30)
		startTime = &t
	case "90d":
		t := now.AddDate(0, 0, -90)
		startTime = &t
	case "1y":
		t := now.AddDate(-1, 0, 0)
		startTime = &t
	default:
		t := now.AddDate(0, 0, -7)
		startTime = &t
	}

	stats, err := l.svcCtx.PointsService.GetStats(l.ctx, startTime, endTime)
	if err != nil {
		return nil, err
	}

	return &types.PointsStats{
		TotalIssued:        stats.TotalIssued,
		TotalRedeemed:      stats.TotalRedeemed,
		TotalExpired:       stats.TotalExpired,
		OutstandingBalance: stats.OutstandingBalance,
		RedemptionRate:     stats.RedemptionRate,
		ActiveUsers:        stats.ActiveUsers,
		PeriodStart:        stats.PeriodStart,
		PeriodEnd:          stats.PeriodEnd,
	}, nil
}
