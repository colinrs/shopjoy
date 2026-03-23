package reviews

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/review"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

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
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	if err := l.svcCtx.ReviewService.DeleteReview(l.ctx, shared.TenantID(tenantID), req.ID); err != nil {
		return nil, err
	}

	return &types.DeleteReviewResp{
		ID:        req.ID,
		Status:    review.StatusDeleted.String(),
		DeletedAt: time.Now().Format(time.RFC3339),
	}, nil
}