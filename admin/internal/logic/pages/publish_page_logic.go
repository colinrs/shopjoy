package pages

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/zeromicro/go-zero/core/logx"
)

type PublishPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) PublishPageLogic {
	return PublishPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishPageLogic) PublishPage(req *types.PublishPageRequest) error {

	userID, _ := contextx.GetUserID(l.ctx)
	return l.svcCtx.PageService.PublishPage(l.ctx, req.ID, userID)
}
