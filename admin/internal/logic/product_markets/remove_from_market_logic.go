// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package product_markets

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/infrastructure/persistence"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveFromMarketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 从市场移除商品
func NewRemoveFromMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveFromMarketLogic {
	return &RemoveFromMarketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveFromMarketLogic) RemoveFromMarket(req *types.RemoveFromMarketReq) error {
	db := l.svcCtx.DB
	repo := persistence.NewProductMarketRepository()

	pm, err := repo.FindByProductAndMarket(l.ctx, db, req.ProductID, req.MarketID, nil)
	if err != nil {
		return err
	}

	return repo.Delete(l.ctx, db, pm.ID)
}
