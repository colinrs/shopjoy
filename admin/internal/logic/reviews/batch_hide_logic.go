package reviews

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
