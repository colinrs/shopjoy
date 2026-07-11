package pages

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type UnpublishPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnpublishPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) UnpublishPageLogic {
	return UnpublishPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnpublishPageLogic) UnpublishPage(req *types.UnpublishPageRequest) error {

	return l.svcCtx.PageService.UnpublishPage(l.ctx, req.ID)
}
