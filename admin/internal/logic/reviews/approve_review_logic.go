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

type ApproveReviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApproveReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) ApproveReviewLogic {
	return ApproveReviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApproveReviewLogic) ApproveReview(req *types.ApproveReviewReq) (resp *types.ApproveReviewResp, err error) {
	// Get tenantID from context with proper validation
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	if err := l.svcCtx.ReviewService.ApproveReview(l.ctx, shared.TenantID(tenantID), req.ID); err != nil {
		return nil, err
	}

	return &types.ApproveReviewResp{
		ID:        req.ID,
		Status:    review.StatusApproved.String(),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}, nil
}