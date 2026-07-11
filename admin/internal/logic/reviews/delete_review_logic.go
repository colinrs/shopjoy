package reviews

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/review"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteReviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteReviewLogic {
	return DeleteReviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteReviewLogic) DeleteReview(req *types.DeleteReviewReq) (resp *types.DeleteReviewResp, err error) {
	// Get tenantID from context with proper validation

	if err := l.svcCtx.ReviewService.DeleteReview(l.ctx, req.ID); err != nil {
		return nil, err
	}

	return &types.DeleteReviewResp{
		ID:        req.ID,
		Status:    review.StatusDeleted.String(),
		DeletedAt: time.Now().UTC().Format(time.RFC3339),
	}, nil
}
