package reviews

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShowReviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShowReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) ShowReviewLogic {
	return ShowReviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShowReviewLogic) ShowReview(req *types.ShowReviewReq) (resp *types.ShowReviewResp, err error) {
	// todo: add your logic here and delete this line

	return
}
