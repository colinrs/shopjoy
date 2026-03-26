package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/storefront"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type pageRepo struct{}

func NewPageRepository() storefront.PageRepository {
	return &pageRepo{}
}

type pageModel struct {
	ID          int64  `gorm:"column:id;primaryKey"`
	TenantID    int64  `gorm:"column:tenant_id;not null;index"`
	Name        string `gorm:"column:name;type:varchar(255);not null"`
	Slug        string `gorm:"column:slug;type:varchar(255);not null"`
	Type        int    `gorm:"column:type;not null;default:0;index"`
	Content     string `gorm:"column:content;type:longtext"`
	SEOTitle    string `gorm:"column:seo_title;type:varchar(255)"`
	SEODesc     string `gorm:"column:seo_description;type:text"`
	SEOKeywords string `gorm:"column:seo_keywords;type:varchar(500)"`
	Status      int8   `gorm:"column:status;not null;default:1;index"`
	Sort        int    `gorm:"column:sort;not null;default:0"`
	IsPublished int    `gorm:"column:is_published;not null;default:0"`
	PublishedAt *int64 `gorm:"column:published_at"`
	Version     int    `gorm:"column:version;not null;default:1"`
	CreatedAt   int64  `gorm:"column:created_at;not null"`
	UpdatedAt   int64  `gorm:"column:updated_at;not null"`
	CreatedBy   int64  `gorm:"column:created_by;not null;default:0"`
	UpdatedBy   int64  `gorm:"column:updated_by;not null;default:0"`
	DeletedAt   *int64 `gorm:"column:deleted_at;index"`
}

func (pageModel) TableName() string {
	return "pages"
}

func (m *pageModel) toEntity() *storefront.Page {
	return &storefront.Page{
		ID:      m.ID,
		TenantID: shared.TenantID(m.TenantID),
		Name:    m.Name,
		Slug:    m.Slug,
		Type:    storefront.PageType(m.Type),
		Content: m.Content,
		SEO: storefront.SEOConfig{
			Title:       m.SEOTitle,
			Description: m.SEODesc,
		},
		Status:      shared.Status(m.Status),
		Sort:        m.Sort,
		IsPublished: m.IsPublished == 1,
		PublishedAt: m.PublishedAt,
		Version:     m.Version,
		Audit: shared.AuditInfo{
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
			CreatedBy: m.CreatedBy,
			UpdatedBy: m.UpdatedBy,
		},
	}
}

func fromPageEntity(p *storefront.Page) *pageModel {
	isPublished := 0
	if p.IsPublished {
		isPublished = 1
	}

	return &pageModel{
		ID:          p.ID,
		TenantID:    p.TenantID.Int64(),
		Name:        p.Name,
		Slug:        p.Slug,
		Type:        int(p.Type),
		Content:     p.Content,
		SEOTitle:    p.SEO.Title,
		SEODesc:     p.SEO.Description,
		SEOKeywords: keywordsToJSONString(p.SEO.Keywords),
		Status:      int8(p.Status),
		Sort:        p.Sort,
		IsPublished: isPublished,
		PublishedAt: p.PublishedAt,
		Version:     p.Version,
		CreatedAt:   p.Audit.CreatedAt,
		UpdatedAt:   p.Audit.UpdatedAt,
		CreatedBy:   p.Audit.CreatedBy,
		UpdatedBy:   p.Audit.UpdatedBy,
	}
}

func keywordsToJSONString(keywords []string) string {
	if len(keywords) == 0 {
		return ""
	}
	data, _ := json.Marshal(keywords)
	return string(data)
}

func (r *pageRepo) Create(ctx context.Context, db *gorm.DB, page *storefront.Page) error {
	model := fromPageEntity(page)
	return db.WithContext(ctx).Create(model).Error
}

func (r *pageRepo) Update(ctx context.Context, db *gorm.DB, page *storefront.Page) error {
	model := fromPageEntity(page)
	return db.WithContext(ctx).Model(&pageModel{}).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", page.ID, page.TenantID.Int64()).
		Updates(map[string]interface{}{
			"name":          model.Name,
			"slug":          model.Slug,
			"type":          model.Type,
			"content":       model.Content,
			"seo_title":     model.SEOTitle,
			"seo_description": model.SEODesc,
			"seo_keywords":  model.SEOKeywords,
			"status":        model.Status,
			"sort":          model.Sort,
			"is_published":  model.IsPublished,
			"published_at":  model.PublishedAt,
			"version":       model.Version,
			"updated_at":    model.UpdatedAt,
			"updated_by":    model.UpdatedBy,
		}).Error
}

func (r *pageRepo) Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error {
	now := time.Now().Unix()
	return db.WithContext(ctx).Model(&pageModel{}).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", id, tenantID.Int64()).
		Update("deleted_at", now).Error
}

func (r *pageRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*storefront.Page, error) {
	var model pageModel
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

func (r *pageRepo) FindBySlug(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, slug string) (*storefront.Page, error) {
	var model pageModel
	err := db.WithContext(ctx).
		Where("slug = ? AND tenant_id = ? AND deleted_at IS NULL", slug, tenantID.Int64()).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *pageRepo) FindByType(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, pageType storefront.PageType) (*storefront.Page, error) {
	var model pageModel
	err := db.WithContext(ctx).
		Where("type = ? AND tenant_id = ? AND deleted_at IS NULL", int(pageType), tenantID.Int64()).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *pageRepo) FindAll(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, page, pageSize int) ([]*storefront.Page, int64, error) {
	var total int64
	if err := db.WithContext(ctx).Model(&pageModel{}).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID.Int64()).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []pageModel
	offset := (page - 1) * pageSize
	err := db.WithContext(ctx).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID.Int64()).
		Order("sort ASC, id DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	pages := make([]*storefront.Page, len(models))
	for i, m := range models {
		pages[i] = m.toEntity()
	}
	return pages, total, nil
}

func (r *pageRepo) CountAll(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (int64, error) {
	var total int64
	err := db.WithContext(ctx).Model(&pageModel{}).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID.Int64()).
		Count(&total).Error
	return total, err
}