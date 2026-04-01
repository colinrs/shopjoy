package reviews

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchHideLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchHideLogic(ctx context.Context, svcCtx *svc.ServiceContext) BatchHideLogic {
	return BatchHideLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchHideLogic) BatchHide(req *types.BatchHideReq) (resp *types.BatchHideResp, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	result, err := l.svcCtx.ReviewService.BatchHide(l.ctx, shared.TenantID(tenantID), req.IDs, req.Reason)
	if err != nil {
		return nil, err
	}

	return &types.BatchHideResp{
		SuccessCount: result.SuccessCount,
		FailedCount:  result.FailedCount,
		Errors:       result.Errors,
	}, nil
}