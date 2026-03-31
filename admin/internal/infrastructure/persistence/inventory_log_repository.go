package persistence

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type inventoryLogRepo struct{}

func NewInventoryLogRepository() product.InventoryLogRepository {
	return &inventoryLogRepo{}
}

type inventoryLogModel struct {
	ID             int64     `gorm:"column:id;primaryKey"`
	TenantID       int64     `gorm:"column:tenant_id;not null;index"`
	SKUCode        string    `gorm:"column:sku_code;type:varchar(50);not null;index"`
	ProductID      int64     `gorm:"column:product_id;not null;index"`
	WarehouseID    int64     `gorm:"column:warehouse_id;not null;default:0"`
	ChangeType     string    `gorm:"column:change_type;type:varchar(20);not null"`
	ChangeQuantity int       `gorm:"column:change_quantity;not null"`
	BeforeStock    int       `gorm:"column:before_stock;not null"`
	AfterStock     int       `gorm:"column:after_stock;not null"`
	OrderNo        string    `gorm:"column:order_no;type:varchar(50)"`
	Remark         string    `gorm:"column:remark;type:varchar(500)"`
	OperatorID     int64     `gorm:"column:operator_id;not null"`
	CreatedAt      time.Time `gorm:"column:created_at;not null;index"`
}

func (inventoryLogModel) TableName() string {
	return "inventory_logs"
}

func (m *inventoryLogModel) toEntity() *product.InventoryLog {
	return &product.InventoryLog{
		Model:         application.Model{ID: m.ID},
		TenantID:      shared.TenantID(m.TenantID),
		SKUCode:       m.SKUCode,
		ProductID:     m.ProductID,
		WarehouseID:   m.WarehouseID,
		ChangeType:    m.ChangeType,
		ChangeQuantity: m.ChangeQuantity,
		BeforeStock:   m.BeforeStock,
		AfterStock:    m.AfterStock,
		OrderNo:       m.OrderNo,
		Remark:        m.Remark,
		OperatorID:    m.OperatorID,
	}
}

func fromInventoryLogEntity(il *product.InventoryLog) *inventoryLogModel {
	return &inventoryLogModel{
		ID:              il.ID,
		TenantID:        il.TenantID.Int64(),
		SKUCode:         il.SKUCode,
		ProductID:       il.ProductID,
		WarehouseID:     il.WarehouseID,
		ChangeType:      il.ChangeType,
		ChangeQuantity:  il.ChangeQuantity,
		BeforeStock:     il.BeforeStock,
		AfterStock:      il.AfterStock,
		OrderNo:         il.OrderNo,
		Remark:          il.Remark,
		OperatorID:      il.OperatorID,
		CreatedAt:       time.Now().UTC(),
	}
}

func (r *inventoryLogRepo) Create(ctx context.Context, db *gorm.DB, log *product.InventoryLog) error {
	model := fromInventoryLogEntity(log)
	return db.WithContext(ctx).Create(model).Error
}

func (r *inventoryLogRepo) FindBySKU(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, skuCode string, query product.InventoryLogQuery) ([]*product.InventoryLog, int64, error) {
	var models []inventoryLogModel
	var total int64

	dbQuery := db.WithContext(ctx).Model(&inventoryLogModel{}).
		Where("sku_code = ? AND tenant_id = ?", skuCode, tenantID.Int64())

	if query.ChangeType != "" {
		dbQuery = dbQuery.Where("change_type = ?", query.ChangeType)
	}
	if !query.StartTime.IsZero() {
		dbQuery = dbQuery.Where("created_at >= ?", query.StartTime)
	}
	if !query.EndTime.IsZero() {
		dbQuery = dbQuery.Where("created_at <= ?", query.EndTime)
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := dbQuery.Order("created_at DESC").Offset(offset).Limit(query.PageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*product.InventoryLog, len(models))
	for i, m := range models {
		result[i] = m.toEntity()
	}
	return result, total, nil
}

func (r *inventoryLogRepo) FindByProduct(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, productID int64, query product.InventoryLogQuery) ([]*product.InventoryLog, int64, error) {
	var models []inventoryLogModel
	var total int64

	dbQuery := db.WithContext(ctx).Model(&inventoryLogModel{}).
		Where("product_id = ? AND tenant_id = ?", productID, tenantID.Int64())

	if query.ChangeType != "" {
		dbQuery = dbQuery.Where("change_type = ?", query.ChangeType)
	}
	if !query.StartTime.IsZero() {
		dbQuery = dbQuery.Where("created_at >= ?", query.StartTime)
	}
	if !query.EndTime.IsZero() {
		dbQuery = dbQuery.Where("created_at <= ?", query.EndTime)
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := dbQuery.Order("created_at DESC").Offset(offset).Limit(query.PageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*product.InventoryLog, len(models))
	for i, m := range models {
		result[i] = m.toEntity()
	}
	return result, total, nil
}

func (r *inventoryLogRepo) FindAll(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query product.InventoryLogQuery) ([]*product.InventoryLog, int64, error) {
	var models []inventoryLogModel
	var total int64

	dbQuery := db.WithContext(ctx).Model(&inventoryLogModel{}).
		Where("tenant_id = ?", tenantID.Int64())

	if query.ChangeType != "" {
		dbQuery = dbQuery.Where("change_type = ?", query.ChangeType)
	}
	if !query.StartTime.IsZero() {
		dbQuery = dbQuery.Where("created_at >= ?", query.StartTime)
	}
	if !query.EndTime.IsZero() {
		dbQuery = dbQuery.Where("created_at <= ?", query.EndTime)
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := dbQuery.Order("created_at DESC").Offset(offset).Limit(query.PageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*product.InventoryLog, len(models))
	for i, m := range models {
		result[i] = m.toEntity()
	}
	return result, total, nil
}
