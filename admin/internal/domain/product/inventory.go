package product

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

// Warehouse represents a physical storage location
type Warehouse struct {
	application.Model
	TenantID  shared.TenantID
	Code      string // 仓库代码
	Name      string // 仓库名称
	Country   string // 所在国家 (ISO country code)
	Address   string // 详细地址
	IsDefault bool   // 是否默认仓库
	Status    shared.Status
	Audit     shared.AuditInfo `gorm:"embedded"`
}

func (w *Warehouse) TableName() string {
	return "warehouses"
}

func (w *Warehouse) Enable() {
	w.Status = shared.StatusEnabled
}

func (w *Warehouse) Disable() {
	w.Status = shared.StatusDisabled
}

func (w *Warehouse) IsEnabled() bool {
	return w.Status == shared.StatusEnabled
}

// WarehouseInventory tracks stock per SKU per warehouse
type WarehouseInventory struct {
	application.Model
	TenantID       shared.TenantID
	SKUCode        string // SKU代码
	WarehouseID    int64
	AvailableStock int              // 可用库存
	LockedStock    int              // 锁定库存
	Audit          shared.AuditInfo `gorm:"embedded"`
}

func (wi *WarehouseInventory) TableName() string {
	return "warehouse_inventories"
}

// TotalStock returns total stock (available + locked)
func (wi *WarehouseInventory) TotalStock() int {
	return wi.AvailableStock + wi.LockedStock
}

// CanDeduct checks if enough stock is available
func (wi *WarehouseInventory) CanDeduct(quantity int) bool {
	return wi.AvailableStock >= quantity
}

// InventoryLog records all stock changes
type InventoryLog struct {
	application.Model
	TenantID       shared.TenantID
	SKUCode        string
	ProductID      int64
	WarehouseID    int64  // 0 = total/summary
	ChangeType     string // manual, order, return, adjustment
	ChangeQuantity int    // positive = increase, negative = decrease
	BeforeStock    int
	AfterStock     int
	OrderNo        string
	Remark         string
	OperatorID     int64
}

func (il *InventoryLog) TableName() string {
	return "inventory_logs"
}

// Inventory change types
const (
	InventoryChangeManual     = "manual"
	InventoryChangeOrder      = "order"
	InventoryChangeReturn     = "return"
	InventoryChangeAdjustment = "adjustment"
)

// WarehouseRepository interface
type WarehouseRepository interface {
	Create(ctx context.Context, db *gorm.DB, warehouse *Warehouse) error
	Update(ctx context.Context, db *gorm.DB, warehouse *Warehouse) error
	Delete(ctx context.Context, db *gorm.DB, id int64) error
	FindByID(ctx context.Context, db *gorm.DB, id int64) (*Warehouse, error)
	FindByCode(ctx context.Context, db *gorm.DB, tenantID int64, code string) (*Warehouse, error)
	FindByTenant(ctx context.Context, db *gorm.DB, tenantID int64) ([]*Warehouse, error)
}

// WarehouseInventoryRepository interface
type WarehouseInventoryRepository interface {
	Create(ctx context.Context, db *gorm.DB, wi *WarehouseInventory) error
	Update(ctx context.Context, db *gorm.DB, wi *WarehouseInventory) error
	FindBySKUAndWarehouse(ctx context.Context, db *gorm.DB, skuCode string, warehouseID int64) (*WarehouseInventory, error)
	FindBySKU(ctx context.Context, db *gorm.DB, skuCode string) ([]*WarehouseInventory, error)
}

// InventoryLogRepository interface
type InventoryLogRepository interface {
	Create(ctx context.Context, db *gorm.DB, log *InventoryLog) error
	Find(ctx context.Context, db *gorm.DB, query InventoryLogQuery) ([]*InventoryLog, int64, error)
}

// InventoryLogQuery for querying inventory logs.
// Optional fields (ProductID, SKUCode, ChangeType) are filter conditions;
// when zero-valued, no filter is applied for that field.
type InventoryLogQuery struct {
	shared.PageQuery
	ProductID  int64
	SKUCode    string
	ChangeType string
	StartTime  time.Time
	EndTime    time.Time
}

// InventoryRepository interface for SKU-level inventory operations
type InventoryRepository interface {
	// SKU inventory
	GetSKUInventory(ctx context.Context, db *gorm.DB, skuCode string) (*SKU, error)
	UpdateSKUStock(ctx context.Context, db *gorm.DB, skuCode string, availableStock, lockedStock int) error
	LockStock(ctx context.Context, db *gorm.DB, skuCode string, quantity int) error
	DeductStock(ctx context.Context, db *gorm.DB, skuCode string, quantity int) error
	RestoreStock(ctx context.Context, db *gorm.DB, skuCode string, quantity int) error

	// Warehouse inventory
	GetWarehouseInventory(ctx context.Context, db *gorm.DB, skuCode string, warehouseID int64) (*WarehouseInventory, error)
	SetWarehouseStock(ctx context.Context, db *gorm.DB, skuCode string, warehouseID int64, stock int) error

	// Logs
	CreateLog(ctx context.Context, db *gorm.DB, log *InventoryLog) error
	GetLogs(ctx context.Context, db *gorm.DB, query InventoryLogQuery) ([]*InventoryLog, int64, error)

	// Low stock
	GetLowStockSKUs(ctx context.Context, db *gorm.DB) ([]*SKU, error)
}
