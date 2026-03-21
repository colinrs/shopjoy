package markets

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/infrastructure/persistence"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetMarketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMarketLogic {
	return &GetMarketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMarketLogic) GetMarket(req *types.GetMarketReq) (resp *types.MarketResponse, err error) {
	repo := persistence.NewMarketRepository()
	m, err := repo.FindByID(l.ctx, l.svcCtx.DB, req.ID)
	if err != nil {
		return nil, err
	}
	return toMarketResponse(m), nil
}