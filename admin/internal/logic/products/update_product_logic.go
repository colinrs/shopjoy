package products

import (
	"context"

	appProduct "github.com/colinrs/shopjoy/admin/internal/application/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/shopspring/decimal"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateProductLogic {
	return UpdateProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProductLogic) UpdateProduct(req *types.UpdateProductReq) (resp *types.ProductDetailResp, err error) {
	// 从 context 获取 tenantID

	// 解析价格字符串（单位：元）
	price, err := appProduct.ToDomainMoneyFromString(req.Price, req.Currency)
	if err != nil {
		return nil, err
	}

	updateReq := appProduct.UpdateProductRequest{
		ID:              req.ID,
		Name:            req.Name,
		Description:     req.Description,
		Price:           price.Amount,
		Currency:        req.Currency,
		CategoryID:      req.CategoryID,
		SKU:             req.SKU,
		Brand:           req.Brand,
		Tags:            req.Tags,
		Images:          req.Images,
		IsMatrixProduct: req.IsMatrixProduct,
		HSCode:          req.HSCode,
		COO:             req.COO,
		Weight:          parseDecimal(req.Weight),
		WeightUnit:      req.WeightUnit,
		Length:          parseDecimal(req.Length),
		Width:           parseDecimal(req.Width),
		Height:          parseDecimal(req.Height),
		DangerousGoods:  req.DangerousGoods,
	}

	productResp, err := l.svcCtx.ProductService.UpdateProduct(l.ctx, updateReq)
	if err != nil {
		return nil, err
	}

	resp = convertToProductDetailResp(productResp)
	resp.CategoryPath = buildCategoryPath(l.ctx, l.svcCtx.DB, productResp.CategoryID)
	return resp, nil
}

func parseDecimal(f float64) decimal.Decimal {
	return decimal.NewFromFloat(f)
}

func convertToProductDetailResp(p *appProduct.ProductResponse) *types.ProductDetailResp {
	// Convert market info
	markets := make([]types.ProductMarketInfo, len(p.Markets))
	for i, m := range p.Markets {
		markets[i] = types.ProductMarketInfo{
			MarketID:   m.MarketID,
			MarketCode: m.MarketCode,
			MarketName: m.MarketName,
			IsEnabled:  m.IsEnabled,
			Price:      m.Price,
			Currency:   m.Currency,
		}
	}

	return &types.ProductDetailResp{
		ID:              p.ID,
		Name:            p.Name,
		Description:     p.Description,
		Price:           p.Price.String(),
		Currency:        p.Currency,
		CostPrice:       p.CostPrice.String(),
		Stock:           p.Stock,
		Status:          p.Status,
		CategoryID:      p.CategoryID,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
		SKU:             p.SKU,
		Brand:           p.Brand,
		Tags:            p.Tags,
		Images:          p.Images,
		IsMatrixProduct: p.IsMatrixProduct,
		HSCode:          p.HSCode,
		COO:             p.COO,
		Weight:          p.Weight.String(),
		WeightUnit:      p.WeightUnit,
		Length:          p.Length.String(),
		Width:           p.Width.String(),
		Height:          p.Height.String(),
		DangerousGoods:  p.DangerousGoods,
		Markets:         markets,
	}
}
