package product

import (
	"github.com/colinrs/shopjoy/admin/internal/domain/product"
)

// CreateProductRequest 创建商品请求
type CreateProductRequest struct {
	Name        string
	Description string
	Price       int64
	Currency    string
	CostPrice   int64
	CategoryID  int64
}

// UpdateProductRequest 更新商品请求
type UpdateProductRequest struct {
	ID          int64
	Name        string
	Description string
	Price       int64
	Currency    string
	CategoryID  int64
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
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Currency    string `json:"currency"`
	CostPrice   int64  `json:"cost_price"`
	Stock       int    `json:"stock"`
	Status      string `json:"status"`
	CategoryID  int64  `json:"category_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
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
	MinPrice   *int64
	MaxPrice   *int64
	MarketID   int64
	Page       int
	PageSize   int
}

// ToDomainMoney 转换为领域层的 Money
func ToDomainMoney(amount int64, currency string) product.Money {
	return product.NewMoney(amount, currency)
}

// FromDomainProduct 从领域实体转换为响应DTO
func FromDomainProduct(p *product.Product) *ProductResponse {
	return &ProductResponse{
		ID:          int64(p.ID),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price.Amount,
		Currency:    p.Price.Currency,
		CostPrice:   p.CostPrice.Amount,
		Stock:       p.Stock,
		Status:      p.Status.String(),
		CategoryID:  p.CategoryID,
		CreatedAt:   p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02 15:04:05"),
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
