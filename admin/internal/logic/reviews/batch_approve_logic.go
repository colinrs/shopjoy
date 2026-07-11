package reviews

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/utils"

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

	ids, err := utils.ParseInt64Slice(req.IDs)
	if err != nil {
		return nil, code.ErrParam
	}

	result, err := l.svcCtx.ReviewService.BatchApprove(l.ctx, ids)
	if err != nil {
		return nil, err
	}

	return &types.BatchApproveResp{
		SuccessCount: result.SuccessCount,
		FailedCount:  result.FailedCount,
		Errors:       result.Errors,
	}, nil
}
