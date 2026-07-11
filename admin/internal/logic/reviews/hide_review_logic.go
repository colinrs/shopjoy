package reviews

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/review"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HideReviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHideReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) HideReviewLogic {
	return HideReviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HideReviewLogic) HideReview(req *types.HideReviewReq) (resp *types.HideReviewResp, err error) {
	// Get tenantID from context with proper validation

	if err := l.svcCtx.ReviewService.HideReview(l.ctx, req.ID, req.Reason); err != nil {
		return nil, err
	}

	return &types.HideReviewResp{
		ID:        req.ID,
		Status:    review.StatusHidden.String(),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}, nil
}
