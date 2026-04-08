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

type GetReviewStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetReviewStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetReviewStatsLogic {
	return GetReviewStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetReviewStatsLogic) GetReviewStats(req *types.ReviewStatsReq) (resp *types.ReviewStatsResp, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	stats, err := l.svcCtx.ReviewService.GetStats(l.ctx, shared.TenantID(tenantID))
	if err != nil {
		return nil, err
	}

	return &types.ReviewStatsResp{
		TotalReviews:     stats.TotalReviews,
		PendingReviews:   stats.PendingReviews,
		ApprovedReviews:  stats.ApprovedReviews,
		HiddenReviews:    stats.HiddenReviews,
		AverageRating:    formatRating(stats.AverageRating),
		QualityAvgRating: formatRating(stats.QualityAvgRating),
		ValueAvgRating:   formatRating(stats.ValueAvgRating),
		FiveStarCount:    stats.FiveStarCount,
		FourStarCount:    stats.FourStarCount,
		ThreeStarCount:   stats.ThreeStarCount,
		TwoStarCount:     stats.TwoStarCount,
		OneStarCount:     stats.OneStarCount,
		WithImageCount:   stats.WithImageCount,
		ReplyRate:        stats.ReplyRate,
		FeaturedCount:    stats.FeaturedCount,
	}, nil
}

func formatRating(r float64) string {
	return fmt.Sprintf("%.2f", r)
}
