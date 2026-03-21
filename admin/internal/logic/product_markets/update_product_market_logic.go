// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package product_markets

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/infrastructure/persistence"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProductMarketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新商品市场配置
func NewUpdateProductMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductMarketLogic {
	return &UpdateProductMarketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProductMarketLogic) UpdateProductMarket(req *types.UpdateProductMarketReq) (resp *types.ProductMarketResp, err error) {
	db := l.svcCtx.DB
	repo := persistence.NewProductMarketRepository()

	pm, err := repo.FindByProductAndMarket(l.ctx, db, req.ProductID, req.MarketID, nil)
	if err != nil {
		return nil, err
	}

	if req.IsEnabled != nil {
		pm.IsEnabled = *req.IsEnabled
	}

	if req.Price != "" {
		pm.Price, _ = decimal.NewFromString(req.Price)
	}

	if req.CompareAtPrice != "" {
		cap, _ := decimal.NewFromString(req.CompareAtPrice)
		pm.CompareAtPrice = &cap
	}

	pm.StockAlertThreshold = req.StockAlertThreshold
	pm.UpdatedAt = time.Now()

	if err := repo.Update(l.ctx, db, pm); err != nil {
		return nil, err
	}

	// Get market info for response
	marketRepo := persistence.NewMarketRepository()
	market, err := marketRepo.FindByID(l.ctx, db, pm.MarketID)
	if err != nil {
		return nil, err
	}

	var compareAtPrice string
	if pm.CompareAtPrice != nil {
		compareAtPrice = pm.CompareAtPrice.String()
	}

	var publishedAt string
	if pm.PublishedAt != nil {
		publishedAt = pm.PublishedAt.Format("2006-01-02 15:04:05")
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
