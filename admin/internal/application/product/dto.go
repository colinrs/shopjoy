package product

import (
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/shopspring/decimal"
)

// CreateProductRequest 创建商品请求
type CreateProductRequest struct {
	Name        string
	Description string
	Price       decimal.Decimal
	Currency    string
	CostPrice   decimal.Decimal
	CategoryID  int64
}

// UpdateProductRequest 更新商品请求
type UpdateProductRequest struct {
	ID              int64
	Name            string
	Description     string
	Price           decimal.Decimal
	Currency        string
	CostPrice       decimal.Decimal
	CategoryID      int64
	SKU             string
	Brand           string
	Tags            []string
	Images          []string
	IsMatrixProduct bool
	HSCode          string
	COO             string
	Weight          decimal.Decimal
	WeightUnit      string
	Length          decimal.Decimal
	Width           decimal.Decimal
	Height          decimal.Decimal
	DangerousGoods  []string
}

// UpdateStockRequest 更新库存请求
type UpdateStockRequest struct {
	ID       int64
	Quantity int
}

// DeductStockRequest 扣减库存请求
type DeductStockRequest struct {
	ID       int64
	Quantity int
}

// ProductResponse 商品响应
type ProductResponse struct {
	ID              int64           `json:"id"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	Price           decimal.Decimal `json:"price"`
	Currency        string          `json:"currency"`
	CostPrice       decimal.Decimal `json:"cost_price"`
	Stock           int             `json:"stock"`
	Status          string          `json:"status"`
	CategoryID      int64           `json:"category_id"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at"`
	SKU             string          `json:"sku"`
	Brand           string          `json:"brand"`
	Tags            []string        `json:"tags"`
	Images          []string        `json:"images"`
	IsMatrixProduct bool            `json:"is_matrix_product"`
	HSCode          string          `json:"hs_code"`
	COO             string          `json:"coo"`
	Weight          decimal.Decimal `json:"weight"`
	WeightUnit      string          `json:"weight_unit"`
	Length          decimal.Decimal `json:"length"`
	Width           decimal.Decimal `json:"width"`
	Height          decimal.Decimal `json:"height"`
	DangerousGoods  []string        `json:"dangerous_goods"`
	Markets         []MarketInfo    `json:"markets"`
}

// MarketInfo 市场信息
type MarketInfo struct {
	MarketID   int64  `json:"market_id"`
	MarketCode string `json:"market_code"`
	MarketName string `json:"market_name"`
	IsEnabled  bool   `json:"is_enabled"`
	Price      string `json:"price"`
	Currency   string `json:"currency"`
}

// ProductListResponse 商品列表响应
type ProductListResponse struct {
	List     []*ProductResponse `json:"list"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

// QueryProductRequest 查询商品请求
type QueryProductRequest struct {
	Name       string
	CategoryID int64
	Status     string
	MinPrice   *decimal.Decimal
	MaxPrice   *decimal.Decimal
	MarketID   int64
	Page       int
	PageSize   int
}

// BatchUpdateProductRequest 批量更新商品请求
type BatchUpdateProductRequest struct {
	ProductIDs []int64
	Fields     BatchProductFields
}

// BatchProductFields 批量更新字段
type BatchProductFields struct {
	Price      *decimal.Decimal
	Stock      *int
	Status     *product.Status
	CategoryID *int64
}

// ToDomainMoney 转换为领域层的 Money
func ToDomainMoney(amount decimal.Decimal, currency string) product.Money {
	return product.NewMoney(amount, currency)
}

// ToDomainMoneyFromString 从字符串（单位：元）转换为领域层的 Money
// 例如 "1.99" 表示 1.99 元
func ToDomainMoneyFromString(amountStr, currency string) (product.Money, error) {
	amount, err := decimal.NewFromString(amountStr)
	if err != nil {
		return product.Money{}, err
	}
	return product.NewMoney(amount, currency), nil
}

// ToDomainMoneyFromInt64 从 int64（单位为分）转换为领域层的 Money
func ToDomainMoneyFromInt64(amount int64, currency string) product.Money {
	return product.NewMoneyFromInt64(amount, currency)
}

// FromDomainProduct 从领域实体转换为响应DTO
func FromDomainProduct(p *product.Product) *ProductResponse {
	return &ProductResponse{
		ID:              int64(p.ID),
		Name:            p.Name,
		Description:     p.Description,
		Price:           p.Price.Amount,
		Currency:        p.Price.Currency,
		CostPrice:       p.CostPrice.Amount,
		Stock:           p.Stock,
		Status:          p.Status.String(),
		CategoryID:      p.CategoryID,
		CreatedAt:       p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       p.UpdatedAt.Format(time.RFC3339),
		SKU:             p.SKU,
		Brand:           p.Brand,
		Tags:            p.Tags,
		Images:          p.Images,
		IsMatrixProduct: p.IsMatrixProduct,
		HSCode:          p.HSCode,
		COO:             p.COO,
		Weight:          p.Weight,
		WeightUnit:      p.WeightUnit,
		Length:          p.Dimensions.Length,
		Width:           p.Dimensions.Width,
		Height:          p.Dimensions.Height,
		DangerousGoods:  p.DangerousGoods,
		Markets:         []MarketInfo{}, // Markets are loaded separately
	}
}

// ParseStatus 解析状态字符串
func ParseStatus(s string) *product.Status {
	switch s {
	case "draft":
		status := product.StatusDraft
		return &status
	case "on_sale":
		status := product.StatusOnSale
		return &status
	case "off_sale":
		status := product.StatusOffSale
		return &status
	case "deleted":
		status := product.StatusDeleted
		return &status
	default:
		return nil // 空字符串返回 nil，表示不过滤状态
	}
}

// BatchProductFail 批量更新失败项
type BatchProductFail struct {
	ProductID int64
	Code      int
	Message   string
}
