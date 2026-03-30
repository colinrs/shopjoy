package reviews

import (
	"context"
	"fmt"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetProductStatsLogic {
	return GetProductStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductStatsLogic) GetProductStats(req *types.ProductStatsReq) (resp *types.ProductStatsResp, err error) {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	stats, err := l.svcCtx.ReviewService.GetProductStats(l.ctx, shared.TenantID(tenantID), req.ProductID)
	if err != nil {
		return nil, err
	}

	// Build rating distribution
	ratingDistribution := &types.RatingDistribution{
		Rating1: stats.Rating1Count,
		Rating2: stats.Rating2Count,
		Rating3: stats.Rating3Count,
		Rating4: stats.Rating4Count,
		Rating5: stats.Rating5Count,
	}

	return &types.ProductStatsResp{
		ProductID:          stats.ProductID,
		TotalReviews:       stats.TotalReviews,
		AverageRating:      formatProductRating(stats.AverageRating),
		QualityAvgRating:   formatProductRating(stats.QualityAvgRating),
		ValueAvgRating:     formatProductRating(stats.ValueAvgRating),
		RatingDistribution: ratingDistribution,
		WithImageCount:     stats.WithImageCount,
		ReplyCount:         0, // TODO: Calculate from replies
		ReplyRate:          0, // TODO: Calculate reply rate
	}, nil
}

func formatProductRating(r float64) string {
	return fmt.Sprintf("%.2f", r)
}