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
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
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

	// Calculate ReplyCount and ReplyRate
	replyCount, replyRate := l.calculateReplyStats(shared.TenantID(tenantID), req.ProductID, stats.TotalReviews)

	return &types.ProductStatsResp{
		ProductID:          stats.ProductID,
		TotalReviews:       stats.TotalReviews,
		AverageRating:      formatProductRating(stats.AverageRating),
		QualityAvgRating:   formatProductRating(stats.QualityAvgRating),
		ValueAvgRating:     formatProductRating(stats.ValueAvgRating),
		RatingDistribution: ratingDistribution,
		WithImageCount:     stats.WithImageCount,
		ReplyCount:         replyCount,
		ReplyRate:          replyRate,
	}, nil
}

func (l *GetProductStatsLogic) calculateReplyStats(tenantID shared.TenantID, productID int64, totalReviews int) (replyCount int, replyRate float64) {
	if totalReviews == 0 {
		return 0, 0
	}

	// Get all reviews for this product
	reviews, err := l.svcCtx.ReviewRepo.FindByProductID(l.ctx, l.svcCtx.DB, tenantID, productID)
	if err != nil {
		l.Logger.Errorf("failed to find reviews for product %d: %v", productID, err)
		return 0, 0
	}

	if len(reviews) == 0 {
		return 0, 0
	}

	// Extract review IDs
	reviewIDs := make([]int64, len(reviews))
	for i, r := range reviews {
		reviewIDs[i] = int64(r.ID)
	}

	// Get all replies for these reviews
	replies, err := l.svcCtx.ReplyRepo.FindByReviewIDs(l.ctx, l.svcCtx.DB, reviewIDs)
	if err != nil {
		l.Logger.Errorf("failed to find replies for product %d: %v", productID, err)
		return 0, 0
	}

	// Count unique review IDs that have replies
	replyMap := make(map[int64]bool)
	for _, reply := range replies {
		replyMap[reply.ReviewID] = true
	}
	replyCount = len(replyMap)

	// Calculate reply rate as percentage
	replyRate = float64(replyCount) / float64(totalReviews) * 100

	return replyCount, replyRate
}

func formatProductRating(r float64) string {
	return fmt.Sprintf("%.2f", r)
}
