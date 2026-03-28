package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type brandRepo struct{}

func NewBrandRepository() product.BrandRepository {
	return &brandRepo{}
}

type brandModel struct {
	ID               int64  `gorm:"column:id;primaryKey"`
	TenantID         int64  `gorm:"column:tenant_id;not null;index"`
	Name             string `gorm:"column:name;type:varchar(100);not null"`
	Logo             string `gorm:"column:logo;type:varchar(500)"`
	Description      string `gorm:"column:description;type:text"`
	Website          string `gorm:"column:website;type:varchar(500)"`
	Sort             int    `gorm:"column:sort;default:0"`
	EnablePage       bool   `gorm:"column:enable_page;default:false"`
	TrademarkNumber  string `gorm:"column:trademark_number;type:varchar(100)"`
	TrademarkCountry string `gorm:"column:trademark_country;type:varchar(10)"`
	Status           int8   `gorm:"column:status;not null;default:1"`
	CreatedAt        int64  `gorm:"column:created_at;not null"`
	UpdatedAt        int64  `gorm:"column:updated_at;not null"`
	CreatedBy        int64  `gorm:"column:created_by"`
	UpdatedBy        int64  `gorm:"column:updated_by"`
	DeletedAt        *int64 `gorm:"column:deleted_at;index"`
}

func (brandModel) TableName() string {
	return "brands"
}

func (m *brandModel) toEntity() *product.Brand {
	return &product.Brand{
		TenantID:         shared.TenantID(m.TenantID),
		Name:             m.Name,
		Logo:             m.Logo,
		Description:      m.Description,
		Website:          m.Website,
		Sort:             m.Sort,
		EnablePage:       m.EnablePage,
		TrademarkNumber:  m.TrademarkNumber,
		TrademarkCountry: m.TrademarkCountry,
		Status:           shared.Status(m.Status),
		Audit: shared.AuditInfo{
			CreatedAt: time.Unix(m.CreatedAt, 0).UTC(),
			UpdatedAt: time.Unix(m.UpdatedAt, 0).UTC(),
			CreatedBy: m.CreatedBy,
			UpdatedBy: m.UpdatedBy,
		},
	}
}

func fromBrandEntity(b *product.Brand) *brandModel {
	return &brandModel{
		ID:               b.Model.ID,
		TenantID:         b.TenantID.Int64(),
		Name:             b.Name,
		Logo:             b.Logo,
		Description:      b.Description,
		Website:          b.Website,
		Sort:             b.Sort,
		EnablePage:       b.EnablePage,
		TrademarkNumber:  b.TrademarkNumber,
		TrademarkCountry: b.TrademarkCountry,
		Status:           int8(b.Status),
		CreatedAt:        b.Audit.CreatedAt.Unix(),
		UpdatedAt:        b.Audit.UpdatedAt.Unix(),
		CreatedBy:        b.Audit.CreatedBy,
		UpdatedBy:        b.Audit.UpdatedBy,
	}
}

func (r *brandRepo) Create(ctx context.Context, db *gorm.DB, b *product.Brand) error {
	model := fromBrandEntity(b)
	return db.WithContext(ctx).Create(model).Error
}

func (r *brandRepo) Update(ctx context.Context, db *gorm.DB, b *product.Brand) error {
	model := fromBrandEntity(b)
	return db.WithContext(ctx).Model(&brandModel{}).
		Where("id = ? AND tenant_id = ?", b.Model.ID, b.TenantID.Int64()).
		Updates(map[string]interface{}{
			"name":              model.Name,
			"logo":              model.Logo,
			"description":       model.Description,
			"website":           model.Website,
			"sort":              model.Sort,
			"enable_page":       model.EnablePage,
			"trademark_number":  model.TrademarkNumber,
			"trademark_country": model.TrademarkCountry,
			"status":            model.Status,
			"updated_at":        model.UpdatedAt,
			"updated_by":        model.UpdatedBy,
		}).Error
}

func (r *brandRepo) Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error {
	now := time.Now().UTC()
	return db.WithContext(ctx).Model(&brandModel{}).
		Where("id = ? AND tenant_id = ?", id, tenantID.Int64()).
		Update("deleted_at", now).Error
}

func (r *brandRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*product.Brand, error) {
	var model brandModel
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

func (r *brandRepo) FindByName(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, name string) (*product.Brand, error) {
	var model brandModel
	err := db.WithContext(ctx).
		Where("name = ? AND tenant_id = ? AND deleted_at IS NULL", name, tenantID.Int64()).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *brandRepo) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query product.BrandQuery) ([]*product.Brand, int64, error) {
	var models []brandModel
	var total int64

	tx := db.WithContext(ctx).Model(&brandModel{}).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID.Int64())

	if query.Name != "" {
		tx = tx.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Status != 0 {
		tx = tx.Where("status = ?", query.Status)
	}

	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := tx.Order("sort ASC, id DESC").
		Offset(int(offset)).
		Limit(int(query.PageSize)).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	brands := make([]*product.Brand, len(models))
	for i, m := range models {
		brands[i] = m.toEntity()
	}
	return brands, total, nil
}

func (r *brandRepo) GetProductCount(ctx context.Context, db *gorm.DB, brandID int64) (int64, error) {
	var count int64
	err := db.WithContext(ctx).Table("products").
		Where("brand_id = ?", brandID).
		Count(&count).Error
	return count, err
}
