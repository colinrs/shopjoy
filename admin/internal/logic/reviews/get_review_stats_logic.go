package reviews

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetReviewStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetReviewStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetReviewStatsLogic {
	return GetReviewStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetReviewStatsLogic) GetReviewStats(req *types.ReviewStatsResp) (resp *types.ReviewStatsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
