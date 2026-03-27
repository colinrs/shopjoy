package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/storefront"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type seoConfigRepo struct{}

func NewSEOConfigRepository() storefront.SEOConfigRepository {
	return &seoConfigRepo{}
}

type seoConfigModel struct {
	ID          int64     `gorm:"column:id;primaryKey"`
	TenantID    int64     `gorm:"column:tenant_id;not null;uniqueIndex:idx_tenant_page_type"`
	PageType    string    `gorm:"column:page_type;type:varchar(30);not null;uniqueIndex:idx_tenant_page_type;index"`
	PageID      *int64    `gorm:"column:page_id;uniqueIndex:idx_tenant_page_type"`
	Title       string    `gorm:"column:title;type:varchar(200);not null;default:''"`
	Description string    `gorm:"column:description;type:text;not null"`
	Keywords    string    `gorm:"column:keywords;type:varchar(500);not null;default:''"`
	CreatedAt   time.Time `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null"`
}

func (seoConfigModel) TableName() string {
	return "seo_configs"
}

func (m *seoConfigModel) toEntity() *storefront.SEOConfigEntity {
	return &storefront.SEOConfigEntity{
		Model:       application.Model{ID: m.ID},
		TenantID:    shared.TenantID(m.TenantID),
		PageType:    m.PageType,
		PageID:      m.PageID,
		Title:       m.Title,
		Description: m.Description,
		Keywords:    m.Keywords,
	}
}

func fromSEOConfigEntity(s *storefront.SEOConfigEntity) *seoConfigModel {
	now := time.Now().UTC()
	return &seoConfigModel{
		ID:          s.ID,
		TenantID:    s.TenantID.Int64(),
		PageType:    s.PageType,
		PageID:      s.PageID,
		Title:       s.Title,
		Description: s.Description,
		Keywords:    s.Keywords,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (r *seoConfigRepo) Save(ctx context.Context, db *gorm.DB, config *storefront.SEOConfigEntity) error {
	model := fromSEOConfigEntity(config)

	// Check if exists
	var existing seoConfigModel
	err := db.WithContext(ctx).
		Where("tenant_id = ? AND page_type = ? AND (page_id = ? OR (page_id IS NULL AND ? IS NULL))",
			model.TenantID, model.PageType, model.PageID, model.PageID).
		First(&existing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new
			return db.WithContext(ctx).Create(model).Error
		}
		return err
	}

	// Update existing
	config.ID = existing.ID
	model.ID = existing.ID
	return db.WithContext(ctx).Model(&seoConfigModel{}).
		Where("id = ?", existing.ID).
		Updates(map[string]interface{}{
			"title":       model.Title,
			"description": model.Description,
			"keywords":    model.Keywords,
			"updated_at":  model.UpdatedAt,
		}).Error
}

func (r *seoConfigRepo) FindByPageType(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, pageType string, pageID *int64) (*storefront.SEOConfigEntity, error) {
	var model seoConfigModel
	query := db.WithContext(ctx).
		Where("tenant_id = ? AND page_type = ?", tenantID.Int64(), pageType)

	if pageID != nil {
		query = query.Where("page_id = ?", *pageID)
	} else {
		query = query.Where("page_id IS NULL")
	}

	err := query.First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *seoConfigRepo) FindAll(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, page, pageSize int) ([]*storefront.SEOConfigEntity, int64, error) {
	var total int64
	if err := db.WithContext(ctx).Model(&seoConfigModel{}).
		Where("tenant_id = ?", tenantID.Int64()).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []seoConfigModel
	offset := (page - 1) * pageSize
	err := db.WithContext(ctx).
		Where("tenant_id = ?", tenantID.Int64()).
		Order("page_type ASC, page_id ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	configs := make([]*storefront.SEOConfigEntity, len(models))
	for i, m := range models {
		configs[i] = m.toEntity()
	}
	return configs, total, nil
}

func (r *seoConfigRepo) CountAll(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (int64, error) {
	var total int64
	err := db.WithContext(ctx).Model(&seoConfigModel{}).
		Where("tenant_id = ?", tenantID.Int64()).
		Count(&total).Error
	return total, err
}

func (r *seoConfigRepo) Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, pageType string, pageID *int64) error {
	query := db.WithContext(ctx).
		Where("tenant_id = ? AND page_type = ?", tenantID.Int64(), pageType)

	if pageID != nil {
		query = query.Where("page_id = ?", *pageID)
	} else {
		query = query.Where("page_id IS NULL")
	}

	return query.Delete(&seoConfigModel{}).Error
}