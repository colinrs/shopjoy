package markets

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMarketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateMarketLogic {
	return CreateMarketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateMarketLogic) CreateMarket(req *types.CreateMarketReq) (resp *types.MarketResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
