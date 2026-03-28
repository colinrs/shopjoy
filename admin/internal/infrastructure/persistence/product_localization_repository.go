package persistence

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type productLocalizationRepo struct{}

func NewProductLocalizationRepository() product.ProductLocalizationRepository {
	return &productLocalizationRepo{}
}

type productLocalizationModel struct {
	ID           int64  `gorm:"column:id;primaryKey"`
	TenantID     int64  `gorm:"column:tenant_id;not null;index"`
	ProductID    int64  `gorm:"column:product_id;not null;index"`
	LanguageCode string `gorm:"column:language_code;not null;size:10"`
	Name         string `gorm:"column:name;not null;size:255"`
	Description  string `gorm:"column:description;type:text"`
	CreatedAt    int64  `gorm:"column:created_at;not null"`
	UpdatedAt    int64  `gorm:"column:updated_at;not null"`
}

func (productLocalizationModel) TableName() string {
	return "product_localizations"
}

func (m *productLocalizationModel) toEntity() *product.ProductLocalization {
	return &product.ProductLocalization{
		Model:        application.Model{ID: m.ID, CreatedAt: time.Unix(m.CreatedAt, 0).UTC(), UpdatedAt: time.Unix(m.UpdatedAt, 0).UTC()},
		TenantID:     shared.TenantID(m.TenantID),
		ProductID:    m.ProductID,
		LanguageCode: m.LanguageCode,
		Name:         m.Name,
		Description:  m.Description,
	}
}

func fromProductLocalizationEntity(pl *product.ProductLocalization) *productLocalizationModel {
	return &productLocalizationModel{
		ID:           pl.Model.ID,
		TenantID:     pl.TenantID.Int64(),
		ProductID:    pl.ProductID,
		LanguageCode: pl.LanguageCode,
		Name:         pl.Name,
		Description:  pl.Description,
		CreatedAt:    pl.Model.CreatedAt.Unix(),
		UpdatedAt:    pl.Model.UpdatedAt.Unix(),
	}
}

func (r *productLocalizationRepo) Create(ctx context.Context, db *gorm.DB, localization *product.ProductLocalization) error {
	model := fromProductLocalizationEntity(localization)
	return db.WithContext(ctx).Create(model).Error
}

func (r *productLocalizationRepo) Update(ctx context.Context, db *gorm.DB, localization *product.ProductLocalization) error {
	return db.WithContext(ctx).Model(&productLocalizationModel{}).
		Where("id = ? AND tenant_id = ?", localization.Model.ID, localization.TenantID.Int64()).
		Updates(map[string]any{
			"language_code": localization.LanguageCode,
			"name":          localization.Name,
			"description":   localization.Description,
			"updated_at":    time.Now().UTC(),
		}).Error
}

func (r *productLocalizationRepo) Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error {
	return db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", id, tenantID.Int64()).
		Delete(&productLocalizationModel{}).Error
}

func (r *productLocalizationRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*product.ProductLocalization, error) {
	var model productLocalizationModel
	err := db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", id, tenantID.Int64()).
		First(&model).Error
	if err != nil {
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *productLocalizationRepo) FindByProductID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, productID int64) ([]*product.ProductLocalization, error) {
	var models []productLocalizationModel
	err := db.WithContext(ctx).
		Where("product_id = ? AND tenant_id = ?", productID, tenantID.Int64()).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	result := make([]*product.ProductLocalization, len(models))
	for i, m := range models {
		result[i] = m.toEntity()
	}
	return result, nil
}

func (r *productLocalizationRepo) FindByProductAndLanguage(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, productID int64, languageCode string) (*product.ProductLocalization, error) {
	var model productLocalizationModel
	err := db.WithContext(ctx).
		Where("product_id = ? AND tenant_id = ? AND language_code = ?", productID, tenantID.Int64(), languageCode).
		First(&model).Error
	if err != nil {
		return nil, err
	}
	return model.toEntity(), nil
}
