package reviews

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateReplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateReplyLogic {
	return UpdateReplyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateReplyLogic) UpdateReply(req *types.UpdateReplyReq) (resp *types.UpdateReplyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
