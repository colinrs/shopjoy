package reviews

import (
	"context"
	"time"

	appReview "github.com/colinrs/shopjoy/admin/internal/application/review"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListReviewsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListReviewsLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListReviewsLogic {
	return ListReviewsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListReviewsLogic) ListReviews(req *types.ListReviewsReq) (resp *types.ListReviewsResp, err error) {
	// Get tenantID from context with proper validation
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	listReq := appReview.ListReviewsRequest{
		ProductID: req.ProductID,
		Status:    req.Status,
		HasImage:  req.HasImage,
		Keyword:   req.Keyword,
		Page:      req.Page,
		PageSize:  req.PageSize,
	}

	if req.RatingMin > 0 {
		min := req.RatingMin
		listReq.RatingMin = &min
	}
	if req.RatingMax > 0 {
		max := req.RatingMax
		listReq.RatingMax = &max
	}

	if req.StartTime != "" {
		t, err := time.Parse(time.RFC3339, req.StartTime)
		if err == nil {
			listReq.StartTime = &t
		}
	}
	if req.EndTime != "" {
		t, err := time.Parse(time.RFC3339, req.EndTime)
		if err == nil {
			listReq.EndTime = &t
		}
	}

	listResp, err := l.svcCtx.ReviewService.ListReviews(l.ctx, shared.TenantID(tenantID), listReq)
	if err != nil {
		return nil, err
	}

	list := make([]*types.ReviewListItem, len(listResp.List))
	for i, item := range listResp.List {
		list[i] = &types.ReviewListItem{
			ID:            item.ID,
			OrderID:       item.OrderID,
			ProductID:     item.ProductID,
			ProductName:   item.ProductName,
			SKUCode:       item.SKUCode,
			UserName:      item.UserName,
			IsAnonymous:   item.IsAnonymous,
			IsVerified:    item.IsVerified,
			QualityRating: item.QualityRating,
			ValueRating:   item.ValueRating,
			OverallRating: item.OverallRating,
			Content:       item.Content,
			Images:        item.Images,
			Status:        item.Status,
			IsFeatured:    item.IsFeatured,
			HelpfulCount:  item.HelpfulCount,
			HasReply:      item.HasReply,
			CreatedAt:     item.CreatedAt,
		}
	}

	return &types.ListReviewsResp{
		List:     list,
		Total:    listResp.Total,
		Page:     listResp.Page,
		PageSize: listResp.PageSize,
	}, nil
}
