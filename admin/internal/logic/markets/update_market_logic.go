package markets

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMarketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateMarketLogic {
	return UpdateMarketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMarketLogic) UpdateMarket(req *types.UpdateMarketReq) (resp *types.MarketResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
