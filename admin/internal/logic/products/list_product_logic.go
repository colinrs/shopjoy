package products

import (
	"context"
	"time"

	appProduct "github.com/colinrs/shopjoy/admin/internal/application/product"
	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListProductLogic {
	return ListProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListProductLogic) ListProduct(req *types.ListProductReq) (resp *types.ListProductResp, err error) {
	// 从 context 获取 tenantID
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok && !contextx.IsPlatformAdmin(l.ctx) {
		return nil, code.ErrUnauthorized
	}
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	queryReq := appProduct.QueryProductRequest{
		Name:       req.Name,
		CategoryID: req.CategoryID,
		Status:     req.Status,
		Page:       req.Page,
		PageSize:   req.PageSize,
		MarketID:   req.MarketID,
	}

	// 解析价格范围（单位：元）
	if req.MinPrice != "" {
		if minPrice, err := decimal.NewFromString(req.MinPrice); err == nil {
			queryReq.MinPrice = &minPrice
		}
	}
	if req.MaxPrice != "" {
		if maxPrice, err := decimal.NewFromString(req.MaxPrice); err == nil {
			queryReq.MaxPrice = &maxPrice
		}
	}

	listResp, err := l.svcCtx.ProductService.GetProductList(l.ctx, shared.TenantID(tenantID), queryReq)
	if err != nil {
		return nil, err
	}

	// Collect product IDs for batch loading market info
	productIDs := make([]int64, len(listResp.List))
	for i, p := range listResp.List {
		productIDs[i] = p.ID
	}

	// Load market associations for all products
	productMarketsMap := make(map[int64][]*product.ProductMarket)
	if len(productIDs) > 0 {
		productMarkets, err := l.svcCtx.ProductMarketRepo.FindByProductIDs(l.ctx, l.svcCtx.DB, productIDs)
		if err != nil {
			return nil, err
		}

		for _, pm := range productMarkets {
			productMarketsMap[pm.ProductID] = append(productMarketsMap[pm.ProductID], pm)
		}
	}

	// Collect unique market IDs
	marketIDsMap := make(map[int64]bool)
	for _, pms := range productMarketsMap {
		for _, pm := range pms {
			marketIDsMap[pm.MarketID] = true
		}
	}

	marketIDs := make([]int64, 0, len(marketIDsMap))
	for id := range marketIDsMap {
		marketIDs = append(marketIDs, id)
	}

	// Load market details
	marketsMap := make(map[int64]*types.MarketResponse)
	if len(marketIDs) > 0 {
		markets, err := l.svcCtx.MarketRepo.FindByIDs(l.ctx, l.svcCtx.DB, marketIDs)
		if err != nil {
			return nil, err
		}

		for _, m := range markets {
			marketsMap[int64(m.ID)] = &types.MarketResponse{
				ID:              int64(m.ID),
				Code:            m.Code,
				Name:            m.Name,
				Currency:        m.Currency,
				DefaultLanguage: m.DefaultLanguage,
				Flag:            m.Flag,
				IsActive:        m.IsActive,
				IsDefault:       m.IsDefault,
				CreatedAt:       m.CreatedAt.Format(time.RFC3339),
				UpdatedAt:       m.UpdatedAt.Format(time.RFC3339),
			}
		}
	}

	list := make([]*types.ProductDetailResp, len(listResp.List))
	for i, p := range listResp.List {
		list[i] = convertToProductDetailRespWithMarkets(p, productMarketsMap[p.ID], marketsMap)
	}

	return &types.ListProductResp{
		List:     list,
		Total:    listResp.Total,
		Page:     listResp.Page,
		PageSize: listResp.PageSize,
	}, nil
}

func convertToProductDetailRespWithMarkets(p *appProduct.ProductResponse, pms []*product.ProductMarket, marketsMap map[int64]*types.MarketResponse) *types.ProductDetailResp {
	resp := &types.ProductDetailResp{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price.String(),
		Currency:    p.Currency,
		CostPrice:   p.CostPrice.String(),
		Stock:       p.Stock,
		Status:      p.Status,
		CategoryID:  p.CategoryID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}

	// Add market info
	if len(pms) > 0 {
		markets := make([]types.ProductMarketInfo, 0, len(pms))
		for _, pm := range pms {
			marketInfo, ok := marketsMap[pm.MarketID]
			if !ok {
				continue
			}
			markets = append(markets, types.ProductMarketInfo{
				MarketID:   pm.MarketID,
				MarketCode: marketInfo.Code,
				MarketName: marketInfo.Name,
				IsEnabled:  pm.IsEnabled,
				Price:      pm.Price.String(),
				Currency:   marketInfo.Currency,
			})
		}
		resp.Markets = markets
	}

	return resp
}
