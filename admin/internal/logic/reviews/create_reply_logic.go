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

type CreateReplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateReplyLogic {
	return CreateReplyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateReplyLogic) CreateReply(req *types.CreateReplyReq) (resp *types.CreateReplyResp, err error) {
	// Get tenantID and admin info from context
	tenantID, _ := contextx.GetTenantID(l.ctx)
	adminID := contextx.GetCurrentUserID(l.ctx)
	adminName := "Admin" // Default admin name, can be enhanced to get from user service

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	reply, err := l.svcCtx.ReviewService.CreateReply(
		l.ctx,
		shared.TenantID(tenantID),
		adminID,
		adminName,
		req.ID,
		appReview.CreateReplyRequest{Content: req.Content},
	)
	if err != nil {
		return nil, err
	}

	return &types.CreateReplyResp{
		ID:        reply.ID,
		ReviewID:  req.ID,
		Content:   reply.Content,
		AdminID:   adminID,
		AdminName: reply.AdminName,
		CreatedAt: reply.CreatedAt,
	}, nil
}