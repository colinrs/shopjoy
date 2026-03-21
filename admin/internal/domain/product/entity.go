// Package product 商品领域层
// 包含实体、值对象和仓储接口定义
package product

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Status 商品状态（值对象）
type Status int

const (
	StatusDraft   Status = iota // 草稿
	StatusOnSale                // 上架
	StatusOffSale               // 下架
	StatusDeleted               // 已删除
)

func (s Status) String() string {
	switch s {
	case StatusDraft:
		return "draft"
	case StatusOnSale:
		return "on_sale"
	case StatusOffSale:
		return "off_sale"
	case StatusDeleted:
		return "deleted"
	default:
		return "unknown"
	}
}

// IsValid 验证状态是否有效
func (s Status) IsValid() bool {
	return s >= StatusDraft && s <= StatusDeleted
}

// CanTransitionTo 检查状态是否可以转换到目标状态
func (s Status) CanTransitionTo(target Status) bool {
	transitions := map[Status][]Status{
		StatusDraft:   {StatusOnSale, StatusDeleted},
		StatusOnSale:  {StatusOffSale, StatusDeleted},
		StatusOffSale: {StatusOnSale, StatusDeleted},
	}

	allowed, ok := transitions[s]
	if !ok {
		return false
	}

	for _, status := range allowed {
		if status == target {
			return true
		}
	}
	return false
}

// Money 金额（值对象）
type Money struct {
	Amount   int64  // 单位为分，避免浮点数精度问题
	Currency string // 币种，如 CNY
}

// NewMoney 创建金额
func NewMoney(amount int64, currency string) Money {
	if currency == "" {
		currency = "CNY"
	}
	return Money{Amount: amount, Currency: currency}
}

// Add 金额相加
func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, code.ErrProductCurrencyMismatch
	}
	return NewMoney(m.Amount+other.Amount, m.Currency), nil
}

// Subtract 金额相减
func (m Money) Subtract(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, code.ErrProductCurrencyMismatch
	}
	if m.Amount < other.Amount {
		return Money{}, code.ErrProductInsufficientAmount
	}
	return NewMoney(m.Amount-other.Amount, m.Currency), nil
}

// Equals 金额相等
func (m Money) Equals(other Money) bool {
	return m.Amount == other.Amount && m.Currency == other.Currency
}

// Dimensions 尺寸（值对象）
type Dimensions struct {
	Length decimal.Decimal
	Width  decimal.Decimal
	Height decimal.Decimal
	Unit   string // cm
}

// Product 商品实体
type Product struct {
	ID              int64          // 商品ID
	SKU             string         // SKU代码
	Name            string         // 商品名称
	Description     string         // 商品描述
	Price           Money          `gorm:"embedded"` // 售价
	CostPrice       Money          `gorm:"embedded"` // 成本价
	Stock           int            // 库存
	Status          Status         // 状态
	CategoryID      int64          // 分类ID
	Brand           string         // 品牌
	Tags            []string       `gorm:"type:json"` // 标签
	Images          []string       `gorm:"type:json"` // 图片列表
	IsMatrixProduct bool           // 是否有变体

	// Compliance fields (cross-border)
	HSCode         string          // HS编码
	COO            string          // 原产国
	Weight         decimal.Decimal // 重量
	WeightUnit     string          // 重量单位
	Dimensions     Dimensions      `gorm:"embedded"` // 尺寸
	DangerousGoods []string        `gorm:"type:json"` // 危险品标识

	CreatedAt       time.Time      // 创建时间
	UpdatedAt       time.Time      // 更新时间
}

// TableName 表名
func (p *Product) TableName() string {
	return "products"
}

// NewProduct 创建新商品
func NewProduct(id int64, name, description string, price Money, categoryID int64) (*Product, error) {
	if name == "" {
		return nil, code.ErrProductEmptyName
	}
	if price.Amount <= 0 {
		return nil, code.ErrProductInvalidPrice
	}
	if id <= 0 {
		return nil, code.ErrProductInvalidID
	}

	now := time.Now()
	return &Product{
		ID:          id,
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       0,
		Status:      StatusDraft,
		CategoryID:  categoryID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// NewProductWithCompliance 创建带合规信息的商品
func NewProductWithCompliance(id int64, name, description, sku string, price Money, categoryID int64) (*Product, error) {
	if name == "" {
		return nil, code.ErrProductEmptyName
	}
	if price.Amount <= 0 {
		return nil, code.ErrProductInvalidPrice
	}
	if id <= 0 {
		return nil, code.ErrProductInvalidID
	}

	now := time.Now()
	return &Product{
		ID:          id,
		SKU:         sku,
		Name:        name,
		Description: description,
		Price:       price,
		Status:      StatusDraft,
		CategoryID:  categoryID,
		WeightUnit:  "g",
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// SetCompliance 设置合规信息
func (p *Product) SetCompliance(hsCode, coo string, weight decimal.Decimal, weightUnit string, dims Dimensions) {
	p.HSCode = hsCode
	p.COO = coo
	p.Weight = weight
	p.WeightUnit = weightUnit
	p.Dimensions = dims
	p.UpdatedAt = time.Now()
}

// HasComplianceInfo 检查是否有合规信息
func (p *Product) HasComplianceInfo() bool {
	return p.HSCode != "" && p.COO != "" && !p.Weight.IsZero()
}

// IsDangerousGoods 检查是否为危险品
func (p *Product) IsDangerousGoods() bool {
	return len(p.DangerousGoods) > 0
}

// PutOnSale 上架商品
func (p *Product) PutOnSale() error {
	if p.Status == StatusDeleted {
		return code.ErrProductDeleted
	}
	if !p.Status.CanTransitionTo(StatusOnSale) {
		return code.ErrProductInvalidStatusTransition
	}
	if p.Stock <= 0 {
		return code.ErrProductNoStock
	}
	p.Status = StatusOnSale
	p.UpdatedAt = time.Now()
	return nil
}

// TakeOffSale 下架商品
func (p *Product) TakeOffSale() error {
	if p.Status == StatusDeleted {
		return code.ErrProductDeleted
	}
	if !p.Status.CanTransitionTo(StatusOffSale) {
		return code.ErrProductInvalidStatusTransition
	}
	p.Status = StatusOffSale
	p.UpdatedAt = time.Now()
	return nil
}

// UpdateStock 更新库存
func (p *Product) UpdateStock(quantity int) error {
	if p.Status == StatusDeleted {
		return code.ErrProductDeleted
	}
	if quantity < 0 {
		return code.ErrProductNegativeStock
	}
	p.Stock = quantity
	p.UpdatedAt = time.Now()
	return nil
}

// DeductStock 扣减库存
func (p *Product) DeductStock(quantity int) error {
	if p.Status != StatusOnSale {
		return code.ErrProductNotOnSale
	}
	if quantity <= 0 {
		return code.ErrProductInvalidQuantity
	}
	if p.Stock < quantity {
		return code.ErrProductInsufficientStock
	}
	p.Stock -= quantity
	p.UpdatedAt = time.Now()
	return nil
}

// UpdatePrice 更新价格
func (p *Product) UpdatePrice(newPrice Money) error {
	if p.Status == StatusDeleted {
		return code.ErrProductDeleted
	}
	if newPrice.Amount <= 0 {
		return code.ErrProductInvalidPrice
	}
	p.Price = newPrice
	p.UpdatedAt = time.Now()
	return nil
}

// SoftDelete 软删除
func (p *Product) SoftDelete() error {
	if p.Status == StatusDeleted {
		return code.ErrProductDeleted
	}
	p.Status = StatusDeleted
	p.UpdatedAt = time.Now()
	return nil
}

// IsOnSale 是否在售
func (p *Product) IsOnSale() bool {
	return p.Status == StatusOnSale && p.Stock > 0
}

type DBProvider interface {
	DB() *gorm.DB
}

type Repository interface {
	Create(ctx context.Context, db *gorm.DB, product *Product) error
	Update(ctx context.Context, db *gorm.DB, product *Product) error
	Delete(ctx context.Context, db *gorm.DB, id int64) error
	FindByID(ctx context.Context, db *gorm.DB, id int64) (*Product, error)
	FindByIDs(ctx context.Context, db *gorm.DB, ids []int64) ([]*Product, error)
	FindList(ctx context.Context, db *gorm.DB, query Query) ([]*Product, int64, error)
	UpdateStock(ctx context.Context, db *gorm.DB, id int64, delta int) error
	Exists(ctx context.Context, db *gorm.DB, id int64) (bool, error)
}

// Query 查询条件（值对象）
type Query struct {
	Name       string
	CategoryID int64
	Status     Status
	MinPrice   *int64
	MaxPrice   *int64
	MarketID   int64
	Page       int
	PageSize   int
}

// Validate 验证查询条件
func (q Query) Validate() error {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 100 {
		q.PageSize = 20
	}
	return nil
}

// Offset 计算偏移量
func (q Query) Offset() int {
	return (q.Page - 1) * q.PageSize
}

// Limit 返回限制数
func (q Query) Limit() int {
	return q.PageSize
}