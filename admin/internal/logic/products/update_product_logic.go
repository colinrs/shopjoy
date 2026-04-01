package products

import (
	"context"

	appProduct "github.com/colinrs/shopjoy/admin/internal/application/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
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
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// 平台管理员设置 tenantID = 0 以访问所有数据
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

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

	productResp, err := l.svcCtx.ProductService.UpdateProduct(l.ctx, shared.TenantID(tenantID), updateReq)
	if err != nil {
		return nil, err
	}

	return convertToProductDetailResp(productResp), nil
}

func parseDecimal(s string) decimal.Decimal {
	if s == "" {
		return decimal.Zero
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero
	}
	return d
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
