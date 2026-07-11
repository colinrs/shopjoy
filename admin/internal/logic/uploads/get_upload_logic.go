// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package uploads

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取图片元数据
func NewGetUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUploadLogic {
	return &GetUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUploadLogic) GetUpload(req *types.GetUploadReq) (resp *types.UploadResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
