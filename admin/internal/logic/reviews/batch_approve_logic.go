package reviews

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
