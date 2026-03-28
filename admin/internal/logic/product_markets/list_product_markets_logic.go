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

type ListProductMarketsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取商品市场配置列表
func NewListProductMarketsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListProductMarketsLogic {
	return &ListProductMarketsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListProductMarketsLogic) ListProductMarkets(req *types.ListProductMarketsReq) (resp *types.ListProductMarketsResp, err error) {
	db := l.svcCtx.DB
	repo := persistence.NewProductMarketRepository()

	productMarkets, err := repo.FindByProductID(l.ctx, db, req.ProductID)
	if err != nil {
		return nil, err
	}

	// Get market info
	marketRepo := persistence.NewMarketRepository()
	markets, err := marketRepo.FindAll(l.ctx, db)
	if err != nil {
		return nil, err
	}

	marketMap := make(map[int64]*types.MarketResponse)
	for _, m := range markets {
		marketMap[int64(m.ID)] = &types.MarketResponse{
			ID:              int64(m.ID),
			Code:            m.Code,
			Name:            m.Name,
			Currency:        m.Currency,
			DefaultLanguage: m.DefaultLanguage,
			Flag:            m.Flag,
			IsActive:        m.IsActive,
			IsDefault:       m.IsDefault,
		}
	}

	list := make([]*types.ProductMarketResp, 0, len(productMarkets))
	for _, pm := range productMarkets {
		market, ok := marketMap[pm.MarketID]
		if !ok {
			continue
		}

		var compareAtPrice string
		if pm.CompareAtPrice != nil {
			compareAtPrice = pm.CompareAtPrice.String()
		}

		var publishedAt string
		if pm.PublishedAt != nil {
			publishedAt = pm.PublishedAt.Format(time.RFC3339)
		}

		list = append(list, &types.ProductMarketResp{
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
		})
	}

	return &types.ListProductMarketsResp{List: list}, nil
}
