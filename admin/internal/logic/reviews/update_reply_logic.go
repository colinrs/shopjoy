package reviews

import (
	"context"

	appReview "github.com/colinrs/shopjoy/admin/internal/application/review"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateReplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateReplyLogic {
	return UpdateReplyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateReplyLogic) UpdateReply(req *types.UpdateReplyReq) (resp *types.UpdateReplyResp, err error) {
	// Get tenantID and admin info from context
	tenantID, _ := contextx.GetTenantID(l.ctx)
	adminID := contextx.GetCurrentUserID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	reply, err := l.svcCtx.ReviewService.UpdateReply(
		l.ctx,
		shared.TenantID(tenantID),
		req.ID,
		appReview.UpdateReplyRequest{Content: req.Content},
	)
	if err != nil {
		return nil, err
	}

	return &types.UpdateReplyResp{
		ID:        reply.ID,
		ReviewID:  req.ID,
		Content:   reply.Content,
		AdminID:   adminID,
		AdminName: reply.AdminName,
		CreatedAt: reply.CreatedAt,
		UpdatedAt: reply.UpdatedAt,
	}, nil
}