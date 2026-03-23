package reviews

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListReviewsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListReviewsLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListReviewsLogic {
	return ListReviewsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListReviewsLogic) ListReviews(req *types.ListReviewsReq) (resp *types.ListReviewsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
