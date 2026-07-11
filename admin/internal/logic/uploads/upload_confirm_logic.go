// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package uploads

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadConfirmLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 前端直传后回调确认入库
func NewUploadConfirmLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadConfirmLogic {
	return &UploadConfirmLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadConfirmLogic) UploadConfirm(req *types.UploadConfirmRequest) (resp *types.UploadConfirmResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
