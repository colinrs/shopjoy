package reviews

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchApproveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchApproveLogic(ctx context.Context, svcCtx *svc.ServiceContext) BatchApproveLogic {
	return BatchApproveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchApproveLogic) BatchApprove(req *types.BatchApproveReq) (resp *types.BatchApproveResp, err error) {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	result, err := l.svcCtx.ReviewService.BatchApprove(l.ctx, shared.TenantID(tenantID), req.IDs)
	if err != nil {
		return nil, err
	}

	return &types.BatchApproveResp{
		SuccessCount: result.SuccessCount,
		FailedCount:  result.FailedCount,
		Errors:       result.Errors,
	}, nil
}