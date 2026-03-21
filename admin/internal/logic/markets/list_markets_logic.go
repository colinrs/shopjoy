package markets

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/infrastructure/persistence"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListMarketsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMarketsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMarketsLogic {
	return &ListMarketsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMarketsLogic) ListMarkets() (resp *types.ListMarketsResp, err error) {
	repo := persistence.NewMarketRepository()
	markets, err := repo.FindAll(l.ctx, l.svcCtx.DB)
	if err != nil {
		return nil, err
	}

	list := make([]*types.MarketResponse, len(markets))
	for i, m := range markets {
		list[i] = toMarketResponse(m)
	}

	return &types.ListMarketsResp{
		List:  list,
		Total: int64(len(list)),
	}, nil
}
