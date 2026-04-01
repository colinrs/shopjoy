package stats

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPointsTrendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPointsTrendLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetPointsTrendLogic {
	return GetPointsTrendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPointsTrendLogic) GetPointsTrend(req *types.GetPointsTrendReq) (resp *types.PointsTrendResp, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	// Calculate time range based on period
	var startTime, endTime time.Time
	now := time.Now().UTC()
	endTime = now

	switch req.Period {
	case "7d":
		startTime = now.AddDate(0, 0, -7)
	case "30d":
		startTime = now.AddDate(0, 0, -30)
	case "90d":
		startTime = now.AddDate(0, 0, -90)
	case "1y":
		startTime = now.AddDate(-1, 0, 0)
	default:
		startTime = now.AddDate(0, 0, -7)
	}

	granularity := req.Granularity
	if granularity == "" {
		granularity = "daily"
	}

	trendData, err := l.svcCtx.PointsService.GetTrend(l.ctx, shared.TenantID(tenantID), startTime, endTime, granularity)
	if err != nil {
		return nil, err
	}

	dataPoints := make([]types.TrendDataPoint, len(trendData))
	for i, d := range trendData {
		dataPoints[i] = types.TrendDataPoint{
			Date:     d.Date,
			Earned:   d.Earned,
			Redeemed: d.Redeemed,
			Expired:  d.Expired,
		}
	}

	return &types.PointsTrendResp{
		Data:   dataPoints,
		Period: req.Period,
	}, nil
}