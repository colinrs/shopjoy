package reviews

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteReviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteReviewLogic {
	return DeleteReviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteReviewLogic) DeleteReview(req *types.DeleteReviewReq) (resp *types.DeleteReviewResp, err error) {
	// todo: add your logic here and delete this line

	return
}
