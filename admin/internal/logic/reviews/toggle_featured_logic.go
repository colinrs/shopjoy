package reviews

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type ToggleFeaturedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewToggleFeaturedLogic(ctx context.Context, svcCtx *svc.ServiceContext) ToggleFeaturedLogic {
	return ToggleFeaturedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ToggleFeaturedLogic) ToggleFeatured(req *types.ToggleFeaturedReq) (resp *types.ToggleFeaturedResp, err error) {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	if err := l.svcCtx.ReviewService.ToggleFeatured(l.ctx, shared.TenantID(tenantID), req.ID, req.IsFeatured); err != nil {
		return nil, err
	}

	return &types.ToggleFeaturedResp{
		ID:         req.ID,
		IsFeatured: req.IsFeatured,
		UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
	}, nil
}