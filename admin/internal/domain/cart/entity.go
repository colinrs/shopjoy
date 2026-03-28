package cart

import (
	"context"

	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

// ==================== Cart 购物车 ====================

// Cart 购物车
type Cart struct {
	application.Model
	TenantID  shared.TenantID `gorm:"column:tenant_id;not null;index"`
	UserID    *int64          `gorm:"column:user_id;index"`
	SessionID string          `gorm:"column:session_id;size:255"`
}

func (c *Cart) TableName() string {
	return "carts"
}

// ==================== CartItem 购物车商品 ====================

// CartItem 购物车商品
type CartItem struct {
	application.Model
	TenantID   shared.TenantID `gorm:"column:tenant_id;not null;index"`
	UserID     int64           `gorm:"column:user_id;not null;index"`
	CartID     int64           `gorm:"column:cart_id;not null;index"`
	ProductID  int64           `gorm:"column:product_id;not null;index"`
	SKUId       int64           `gorm:"column:sku_id;not null;index"`
	ProductName string          `gorm:"column:product_name;size:255;not null"`
	SKUName     string          `gorm:"column:sku_name;size:255;not null;default:''"`
	Image       string          `gorm:"column:image;size:500;not null;default:''"`
	Price       shared.Money    `gorm:"column:price;type:decimal(19,4);not null;embedded"`
	Quantity    int             `gorm:"column:quantity;not null;default:1"`
	TotalAmount shared.Money     `gorm:"column:total_amount;type:decimal(19,4);not null;embedded"`
	Selected    bool            `gorm:"column:selected;not null;default:true"`
	Audit       shared.AuditInfo `gorm:"embedded"`
}

func (ci *CartItem) TableName() string {
	return "cart_items"
}

// CalculateTotal 计算小计金额
func (ci *CartItem) CalculateTotal() {
	ci.TotalAmount = ci.Price.Multiply(ci.Quantity)
}

// ==================== Repository 接口 ====================

// Repository 购物车仓储接口
type Repository interface {
	// Cart operations
	FindOrCreateCart(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID *int64, sessionID string) (*Cart, error)
	FindByUserID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64) (*Cart, error)
	FindBySessionID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, sessionID string) (*Cart, error)
	Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error

	// CartItem operations
	AddItem(ctx context.Context, db *gorm.DB, item *CartItem) error
	UpdateItem(ctx context.Context, db *gorm.DB, item *CartItem) error
	RemoveItem(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
	FindItemsByCartID(ctx context.Context, db *gorm.DB, cartID int64) ([]CartItem, error)
	FindItemsByUserID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64) ([]CartItem, error)
	ClearItems(ctx context.Context, db *gorm.DB, cartID int64) error
}

// ==================== Query 查询条件 ====================

// Query 查询条件
type Query struct {
	shared.PageQuery
	UserID int64
}
