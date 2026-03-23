package reviews

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteReplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteReplyLogic {
	return DeleteReplyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteReplyLogic) DeleteReply(req *types.DeleteReplyReq) (resp *types.DeleteReplyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
