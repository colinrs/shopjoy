package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type warehouseRepo struct{}

func NewWarehouseRepository() product.WarehouseRepository {
	return &warehouseRepo{}
}

type warehouseModel struct {
	ID        int64  `gorm:"column:id;primaryKey"`
	TenantID  int64  `gorm:"column:tenant_id;not null;index"`
	Code      string `gorm:"column:code;type:varchar(50);not null"`
	Name      string `gorm:"column:name;type:varchar(100);not null"`
	Country   string `gorm:"column:country;type:varchar(10)"`
	Address   string `gorm:"column:address;type:varchar(500)"`
	IsDefault bool   `gorm:"column:is_default;default:false"`
	Status    int8   `gorm:"column:status;not null;default:1"`
	CreatedAt int64  `gorm:"column:created_at;not null"`
	UpdatedAt int64  `gorm:"column:updated_at;not null"`
	DeletedAt *int64 `gorm:"column:deleted_at"`
}

func (warehouseModel) TableName() string {
	return "warehouses"
}

func (m *warehouseModel) toEntity() *product.Warehouse {
	return &product.Warehouse{
		ID:        m.ID,
		TenantID:  shared.TenantID(m.TenantID),
		Code:      m.Code,
		Name:      m.Name,
		Country:   m.Country,
		Address:   m.Address,
		IsDefault: m.IsDefault,
		Status:    shared.Status(m.Status),
		Audit: shared.AuditInfo{
			CreatedAt: time.Unix(m.CreatedAt, 0),
			UpdatedAt: time.Unix(m.UpdatedAt, 0),
		},
	}
}

func fromWarehouseEntity(w *product.Warehouse) *warehouseModel {
	now := time.Now().Unix()
	createdAt := now
	updatedAt := now
	if !w.Audit.CreatedAt.IsZero() {
		createdAt = w.Audit.CreatedAt.Unix()
	}
	if !w.Audit.UpdatedAt.IsZero() {
		updatedAt = w.Audit.UpdatedAt.Unix()
	}
	return &warehouseModel{
		ID:        w.ID,
		TenantID:  w.TenantID.Int64(),
		Code:      w.Code,
		Name:      w.Name,
		Country:   w.Country,
		Address:   w.Address,
		IsDefault: w.IsDefault,
		Status:    int8(w.Status),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func (r *warehouseRepo) Create(ctx context.Context, db *gorm.DB, w *product.Warehouse) error {
	model := fromWarehouseEntity(w)
	return db.WithContext(ctx).Create(model).Error
}

func (r *warehouseRepo) Update(ctx context.Context, db *gorm.DB, w *product.Warehouse) error {
	model := fromWarehouseEntity(w)
	return db.WithContext(ctx).Model(&warehouseModel{}).
		Where("id = ? AND tenant_id = ?", w.ID, w.TenantID.Int64()).
		Updates(map[string]any{
			"name":       model.Name,
			"country":    model.Country,
			"address":    model.Address,
			"is_default": model.IsDefault,
			"status":     model.Status,
			"updated_at": model.UpdatedAt,
		}).Error
}

func (r *warehouseRepo) Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error {
	now := time.Now().Unix()
	return db.WithContext(ctx).Model(&warehouseModel{}).
		Where("id = ? AND tenant_id = ?", id, tenantID.Int64()).
		Update("deleted_at", now).Error
}

func (r *warehouseRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*product.Warehouse, error) {
	var model warehouseModel
	err := db.WithContext(ctx).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", id, tenantID.Int64()).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *warehouseRepo) FindByCode(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, code string) (*product.Warehouse, error) {
	var model warehouseModel
	err := db.WithContext(ctx).
		Where("code = ? AND tenant_id = ? AND deleted_at IS NULL", code, tenantID.Int64()).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *warehouseRepo) FindAll(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) ([]*product.Warehouse, error) {
	var models []warehouseModel
	err := db.WithContext(ctx).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID.Int64()).
		Order("is_default DESC, code ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	warehouses := make([]*product.Warehouse, len(models))
	for i, m := range models {
		warehouses[i] = m.toEntity()
	}
	return warehouses, nil
}

func (r *warehouseRepo) FindDefault(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (*product.Warehouse, error) {
	var model warehouseModel
	err := db.WithContext(ctx).
		Where("tenant_id = ? AND is_default = ? AND deleted_at IS NULL", tenantID.Int64(), true).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.toEntity(), nil
}