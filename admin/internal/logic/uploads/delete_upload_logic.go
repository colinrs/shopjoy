package uploads

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUploadLogic {
	return &DeleteUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUploadLogic) DeleteUpload(req *types.DeleteUploadReq) error {
	err := l.svcCtx.Storage.Delete(l.ctx, req.ID)
	if err != nil {
		return code.ErrUploadNotFound
	}
	return nil
}