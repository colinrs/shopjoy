// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package uploads

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadSignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 申请 Cloudinary 签名（前端直传第一步）
func NewUploadSignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadSignLogic {
	return &UploadSignLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadSignLogic) UploadSign(req *types.UploadSignRequest) (resp *types.UploadSignResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
