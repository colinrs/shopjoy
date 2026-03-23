package reviews

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
