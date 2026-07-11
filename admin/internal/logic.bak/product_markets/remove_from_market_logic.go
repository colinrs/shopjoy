// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package product_markets

import (
	"context"
	"time"

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

func (l *RemoveFromMarketLogic) RemoveFromMarket(req *types.RemoveFromMarketReq) (resp *types.ProductMarketResp, err error) {
	db := l.svcCtx.DB
	repo := persistence.NewProductMarketRepository()

	pm, err := repo.FindByProductAndMarket(l.ctx, db, req.ProductID, req.MarketID, nil)
	if err != nil {
		return nil, err
	}

	marketRepo := persistence.NewMarketRepository()
	market, err := marketRepo.FindByID(l.ctx, db, pm.MarketID)
	if err != nil {
		return nil, err
	}

	if err := repo.Delete(l.ctx, db, pm.ID); err != nil {
		return nil, err
	}

	var compareAtPrice string
	if pm.CompareAtPrice != nil {
		compareAtPrice = pm.CompareAtPrice.String()
	}

	var publishedAt string
	if pm.PublishedAt != nil {
		publishedAt = pm.PublishedAt.Format(time.RFC3339)
	}

	return &types.ProductMarketResp{
		ID:                  pm.ID,
		ProductID:           pm.ProductID,
		MarketID:            pm.MarketID,
		MarketCode:          market.Code,
		MarketName:          market.Name,
		IsEnabled:           pm.IsEnabled,
		Price:               pm.Price.String(),
		CompareAtPrice:      compareAtPrice,
		Currency:            market.Currency,
		StockAlertThreshold: pm.StockAlertThreshold,
		PublishedAt:         publishedAt,
	}, nil
}
