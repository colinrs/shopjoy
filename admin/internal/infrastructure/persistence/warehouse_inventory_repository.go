package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type warehouseInventoryRepo struct{}

func NewWarehouseInventoryRepository() product.WarehouseInventoryRepository {
	return &warehouseInventoryRepo{}
}

type warehouseInventoryModel struct {
	ID             int64  `gorm:"column:id;primaryKey"`
	TenantID       int64  `gorm:"column:tenant_id;not null;index"`
	SKUCode        string `gorm:"column:sku_code;type:varchar(50);not null;index"`
	WarehouseID    int64  `gorm:"column:warehouse_id;not null;index"`
	AvailableStock int    `gorm:"column:available_stock;not null;default:0"`
	LockedStock    int    `gorm:"column:locked_stock;not null;default:0"`
	CreatedAt      int64  `gorm:"column:created_at;not null"`
	UpdatedAt      int64  `gorm:"column:updated_at;not null"`
}

func (warehouseInventoryModel) TableName() string {
	return "warehouse_inventories"
}

func (m *warehouseInventoryModel) toEntity() *product.WarehouseInventory {
	return &product.WarehouseInventory{
		ID:             m.ID,
		TenantID:       shared.TenantID(m.TenantID),
		SKUCode:        m.SKUCode,
		WarehouseID:    m.WarehouseID,
		AvailableStock: m.AvailableStock,
		LockedStock:    m.LockedStock,
		Audit: shared.AuditInfo{
			CreatedAt: time.Unix(m.CreatedAt, 0).UTC(),
			UpdatedAt: time.Unix(m.UpdatedAt, 0).UTC(),
		},
	}
}

func fromWarehouseInventoryEntity(wi *product.WarehouseInventory) *warehouseInventoryModel {
	return &warehouseInventoryModel{
		ID:             wi.ID,
		TenantID:       wi.TenantID.Int64(),
		SKUCode:        wi.SKUCode,
		WarehouseID:    wi.WarehouseID,
		AvailableStock: wi.AvailableStock,
		LockedStock:    wi.LockedStock,
		CreatedAt:      wi.Audit.CreatedAt.Unix(),
		UpdatedAt:      wi.Audit.UpdatedAt.Unix(),
	}
}

func (r *warehouseInventoryRepo) Create(ctx context.Context, db *gorm.DB, wi *product.WarehouseInventory) error {
	model := fromWarehouseInventoryEntity(wi)
	return db.WithContext(ctx).Create(model).Error
}

func (r *warehouseInventoryRepo) Update(ctx context.Context, db *gorm.DB, wi *product.WarehouseInventory) error {
	model := fromWarehouseInventoryEntity(wi)
	return db.WithContext(ctx).Model(&warehouseInventoryModel{}).
		Where("id = ? AND tenant_id = ?", wi.ID, wi.TenantID.Int64()).
		Updates(map[string]any{
			"available_stock": model.AvailableStock,
			"locked_stock":    model.LockedStock,
			"updated_at":      model.UpdatedAt,
		}).Error
}

func (r *warehouseInventoryRepo) FindBySKUAndWarehouse(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, skuCode string, warehouseID int64) (*product.WarehouseInventory, error) {
	var model warehouseInventoryModel
	err := db.WithContext(ctx).
		Where("sku_code = ? AND warehouse_id = ? AND tenant_id = ?", skuCode, warehouseID, tenantID.Int64()).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *warehouseInventoryRepo) FindBySKU(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, skuCode string) ([]*product.WarehouseInventory, error) {
	var models []warehouseInventoryModel
	err := db.WithContext(ctx).
		Where("sku_code = ? AND tenant_id = ?", skuCode, tenantID.Int64()).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	result := make([]*product.WarehouseInventory, len(models))
	for i, m := range models {
		result[i] = m.toEntity()
	}
	return result, nil
}
