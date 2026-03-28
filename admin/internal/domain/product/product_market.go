// Package product 商品领域层
package product

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ProductMarket 商品-市场关联实体
type ProductMarket struct {
	application.Model
	TenantID            int64            // 租户ID
	ProductID           int64            // 商品ID
	VariantID           *int64           // 变体ID，NULL表示基础商品
	MarketID            int64            // 市场ID
	IsEnabled           bool             // 是否在该市场可见
	StatusOverride      *Status          // 状态覆盖
	Price               decimal.Decimal  // 市场专属价格
	CompareAtPrice      *decimal.Decimal // 对比价格
	StockAlertThreshold int              // 库存预警阈值
	PublishedAt         *time.Time       // 发布时间
}

// TableName 表名
func (pm *ProductMarket) TableName() string {
	return "product_markets"
}

// NewProductMarket 创建商品市场关联
func NewProductMarket(productID, marketID int64) (*ProductMarket, error) {
	if productID <= 0 {
		return nil, code.ErrProductInvalidID
	}
	if marketID <= 0 {
		return nil, code.ErrMarketNotFound
	}

	return &ProductMarket{
		ProductID: productID,
		MarketID:  marketID,
		IsEnabled: false, // Default to disabled, requires price setup
	}, nil
}

// Enable 启用市场可见性
func (pm *ProductMarket) Enable() {
	pm.IsEnabled = true
	pm.UpdatedAt = time.Now().UTC()
}

// Disable 禁用市场可见性
func (pm *ProductMarket) Disable() {
	pm.IsEnabled = false
	pm.UpdatedAt = time.Now().UTC()
}

// SetPrice 设置价格
func (pm *ProductMarket) SetPrice(price decimal.Decimal) error {
	if price.LessThanOrEqual(decimal.Zero) {
		return code.ErrProductMarketPriceRequired
	}
	pm.Price = price
	pm.UpdatedAt = time.Now().UTC()
	return nil
}

// Publish 发布到市场
func (pm *ProductMarket) Publish() error {
	if pm.Price.IsZero() {
		return code.ErrProductMarketPriceRequired
	}
	pm.Enable()
	now := time.Now().UTC()
	pm.PublishedAt = &now
	return nil
}

// ProductMarketRepository 商品市场仓储接口
type ProductMarketRepository interface {
	Create(ctx context.Context, db *gorm.DB, pm *ProductMarket) error
	Update(ctx context.Context, db *gorm.DB, pm *ProductMarket) error
	Delete(ctx context.Context, db *gorm.DB, id int64) error
	FindByID(ctx context.Context, db *gorm.DB, id int64) (*ProductMarket, error)
	FindByProductAndMarket(ctx context.Context, db *gorm.DB, productID, marketID int64, variantID *int64) (*ProductMarket, error)
	FindByProductID(ctx context.Context, db *gorm.DB, productID int64) ([]*ProductMarket, error)
	FindByProductIDs(ctx context.Context, db *gorm.DB, productIDs []int64) ([]*ProductMarket, error)
	FindByMarketID(ctx context.Context, db *gorm.DB, marketID int64, query Query) ([]*ProductMarket, int64, error)
	BatchCreate(ctx context.Context, db *gorm.DB, pms []*ProductMarket) error
}
