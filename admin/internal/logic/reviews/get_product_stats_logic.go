package reviews

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetProductStatsLogic {
	return GetProductStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductStatsLogic) GetProductStats(req *types.ProductStatsReq) (resp *types.ProductStatsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
