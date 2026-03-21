package product

import (
	"context"

	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type CategoryStatus int

const (
	CategoryStatusDisabled CategoryStatus = iota
	CategoryStatusEnabled
)

type Category struct {
	ID             int64
	TenantID       shared.TenantID
	ParentID       int64
	Name           string
	Code           string
	Level          int
	Sort           int
	Icon           string
	Image          string
	SeoTitle       string // SEO标题
	SeoDescription string // SEO描述
	Status         CategoryStatus
	Audit          shared.AuditInfo `gorm:"embedded"`
}

func (c *Category) TableName() string {
	return "categories"
}

func (c *Category) Enable() {
	c.Status = CategoryStatusEnabled
}

func (c *Category) Disable() {
	c.Status = CategoryStatusDisabled
}

func (c *Category) IsActive() bool {
	return c.Status == CategoryStatusEnabled
}

func (c *Category) IsRoot() bool {
	return c.ParentID == 0
}

type Brand struct {
	ID               int64
	TenantID         shared.TenantID
	Name             string
	Logo             string
	Description      string
	Website          string
	Sort             int
	EnablePage       bool   // 是否启用品牌专区
	TrademarkNumber  string // 商标号
	TrademarkCountry string // 商标注册国家 (ISO country code)
	Status           shared.Status
	Audit            shared.AuditInfo `gorm:"embedded"`
}

func (b *Brand) TableName() string {
	return "brands"
}

func (b *Brand) Enable() {
	b.Status = shared.StatusEnabled
}

func (b *Brand) Disable() {
	b.Status = shared.StatusDisabled
}

func (b *Brand) IsEnabled() bool {
	return b.Status == shared.StatusEnabled
}

func (b *Brand) TogglePage(enabled bool) {
	b.EnablePage = enabled
}

type Attribute struct {
	ID         int64
	TenantID   shared.TenantID
	Name       string
	Code       string
	InputType  AttributeInputType
	Options    []string
	IsRequired bool
	Status     shared.Status
	Audit      shared.AuditInfo `gorm:"embedded"`
}

type AttributeInputType int

const (
	InputTypeText AttributeInputType = iota
	InputTypeSelect
	InputTypeMultiSelect
	InputTypeNumber
	InputTypeBoolean
)

type SKU struct {
	ID             int64
	TenantID       shared.TenantID
	ProductID      int64
	Code           string
	Price          shared.Money `gorm:"embedded"`
	Stock          int          // 总库存
	AvailableStock int          // 可用库存
	LockedStock    int          // 锁定库存
	SafetyStock    int          // 安全库存阈值
	PreSaleEnabled bool         // 是否开启预售
	Attributes     map[string]string
	Status         shared.Status
	Audit          shared.AuditInfo `gorm:"embedded"`
}

func (s *SKU) TableName() string {
	return "skus"
}

func (s *SKU) IsAvailable() bool {
	return s.Status == shared.StatusEnabled && s.AvailableStock > 0
}

func (s *SKU) CanDeduct(quantity int) bool {
	return s.AvailableStock >= quantity
}

// DeductStock deducts from available stock (for payment success)
func (s *SKU) DeductStock(quantity int) error {
	if s.AvailableStock < quantity {
		return code.ErrInventoryInsufficientStock
	}
	s.AvailableStock -= quantity
	s.Stock = s.AvailableStock + s.LockedStock
	return nil
}

// LockStock moves stock from available to locked (for order creation)
func (s *SKU) LockStock(quantity int) error {
	if s.AvailableStock < quantity {
		return code.ErrInventoryInsufficientStock
	}
	s.AvailableStock -= quantity
	s.LockedStock += quantity
	s.Stock = s.AvailableStock + s.LockedStock
	return nil
}

// UnlockStock moves stock from locked to available (for order cancellation)
func (s *SKU) UnlockStock(quantity int) error {
	if s.LockedStock < quantity {
		return code.ErrInventoryInsufficientLockedStock
	}
	s.LockedStock -= quantity
	s.AvailableStock += quantity
	s.Stock = s.AvailableStock + s.LockedStock
	return nil
}

// DeductLockedStock deducts from locked stock (for payment success when locked)
func (s *SKU) DeductLockedStock(quantity int) error {
	if s.LockedStock < quantity {
		return code.ErrInventoryInsufficientLockedStock
	}
	s.LockedStock -= quantity
	s.Stock = s.AvailableStock + s.LockedStock
	return nil
}

// RestoreStock adds to available stock
func (s *SKU) RestoreStock(quantity int) {
	s.AvailableStock += quantity
	s.Stock = s.AvailableStock + s.LockedStock
}

// IsLowStock checks if available stock is below safety threshold
func (s *SKU) IsLowStock() bool {
	return s.SafetyStock > 0 && s.AvailableStock < s.SafetyStock
}

type CategoryRepository interface {
	Create(ctx context.Context, db *gorm.DB, category *Category) error
	Update(ctx context.Context, db *gorm.DB, category *Category) error
	Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Category, error)
	FindByParentID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, parentID int64) ([]*Category, error)
	FindAll(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) ([]*Category, error)
	FindTree(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) ([]*Category, error)
	FindByCode(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, code string) (*Category, error)
	GetProductCount(ctx context.Context, db *gorm.DB, categoryID int64) (int64, error)
	UpdateSort(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, sorts []CategorySort) error
	Move(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64, newParentID int64) error
}

type BrandRepository interface {
	Create(ctx context.Context, db *gorm.DB, brand *Brand) error
	Update(ctx context.Context, db *gorm.DB, brand *Brand) error
	Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Brand, error)
	FindByName(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, name string) (*Brand, error)
	FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query BrandQuery) ([]*Brand, int64, error)
	GetProductCount(ctx context.Context, db *gorm.DB, brandID int64) (int64, error)
}

type BrandQuery struct {
	shared.PageQuery
	Name   string
	Status shared.Status
}

// CategoryMarket represents category visibility in specific markets
type CategoryMarket struct {
	ID         int64
	TenantID   shared.TenantID
	CategoryID int64
	MarketID   int64
	IsVisible  bool
	Audit      shared.AuditInfo `gorm:"embedded"`
}

func (cm *CategoryMarket) TableName() string {
	return "category_markets"
}

// BrandMarket represents brand visibility in specific markets
type BrandMarket struct {
	ID        int64
	TenantID  shared.TenantID
	BrandID   int64
	MarketID  int64
	IsVisible bool
	Audit     shared.AuditInfo `gorm:"embedded"`
}

func (bm *BrandMarket) TableName() string {
	return "brand_markets"
}

// CategorySort for batch sort updates
type CategorySort struct {
	ID   int64
	Sort int
}

// CategoryMarketRepository interface
type CategoryMarketRepository interface {
	Create(ctx context.Context, db *gorm.DB, cm *CategoryMarket) error
	Update(ctx context.Context, db *gorm.DB, cm *CategoryMarket) error
	FindByCategory(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, categoryID int64) ([]*CategoryMarket, error)
	DeleteByCategory(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, categoryID int64) error
	BatchCreate(ctx context.Context, db *gorm.DB, items []*CategoryMarket) error
}

// BrandMarketRepository interface
type BrandMarketRepository interface {
	Create(ctx context.Context, db *gorm.DB, bm *BrandMarket) error
	Update(ctx context.Context, db *gorm.DB, bm *BrandMarket) error
	FindByBrand(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, brandID int64) ([]*BrandMarket, error)
	DeleteByBrand(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, brandID int64) error
	BatchCreate(ctx context.Context, db *gorm.DB, items []*BrandMarket) error
}
