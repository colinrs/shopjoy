package reviews

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetReviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetReviewLogic {
	return GetReviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetReviewLogic) GetReview(req *types.GetReviewReq) (resp *types.ReviewDetailResp, err error) {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	detail, err := l.svcCtx.ReviewService.GetReview(l.ctx, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}

	resp = &types.ReviewDetailResp{
		ID:            detail.ID,
		TenantID:      detail.TenantID,
		OrderID:       detail.OrderID,
		ProductID:     detail.ProductID,
		ProductName:   detail.ProductName,
		SKUCode:       detail.SKUCode,
		UserID:        detail.UserID,
		UserName:      detail.UserName,
		IsAnonymous:   detail.IsAnonymous,
		IsVerified:    detail.IsVerified,
		QualityRating: detail.QualityRating,
		ValueRating:   detail.ValueRating,
		OverallRating: detail.OverallRating,
		Content:       detail.Content,
		Images:        detail.Images,
		Status:        detail.Status,
		IsFeatured:    detail.IsFeatured,
		HelpfulCount:  detail.HelpfulCount,
		CreatedAt:     detail.CreatedAt,
		UpdatedAt:     detail.UpdatedAt,
	}

	if detail.Reply != nil {
		resp.Reply = &types.ReviewReplyResp{
			ID:        detail.Reply.ID,
			Content:   detail.Reply.Content,
			AdminName: detail.Reply.AdminName,
			CreatedAt: detail.Reply.CreatedAt,
			UpdatedAt: detail.Reply.UpdatedAt,
		}
	}

	return resp, nil
}