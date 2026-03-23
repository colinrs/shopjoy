package reviews

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HideReviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHideReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) HideReviewLogic {
	return HideReviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HideReviewLogic) HideReview(req *types.HideReviewReq) (resp *types.HideReviewResp, err error) {
	// todo: add your logic here and delete this line

	return
}
