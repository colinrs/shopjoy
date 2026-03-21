package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type categoryRepo struct{}

func NewCategoryRepository() product.CategoryRepository {
	return &categoryRepo{}
}

type categoryModel struct {
	ID             int64  `gorm:"column:id;primaryKey"`
	TenantID       int64  `gorm:"column:tenant_id;not null;index"`
	ParentID       int64  `gorm:"column:parent_id;default:0;index"`
	Name           string `gorm:"column:name;type:varchar(100);not null"`
	Code           string `gorm:"column:code;type:varchar(50)"`
	Level          int    `gorm:"column:level;not null;default:1"`
	Sort           int    `gorm:"column:sort;default:0"`
	Icon           string `gorm:"column:icon;type:varchar(500)"`
	Image          string `gorm:"column:image;type:varchar(500)"`
	SeoTitle       string `gorm:"column:seo_title;type:varchar(200)"`
	SeoDescription string `gorm:"column:seo_description;type:varchar(500)"`
	Status         int8   `gorm:"column:status;not null;default:1"`
	CreatedAt      int64  `gorm:"column:created_at;not null"`
	UpdatedAt      int64  `gorm:"column:updated_at;not null"`
	CreatedBy      int64  `gorm:"column:created_by"`
	UpdatedBy      int64  `gorm:"column:updated_by"`
	DeletedAt      *int64 `gorm:"column:deleted_at;index"`
}

func (categoryModel) TableName() string {
	return "categories"
}

func (m *categoryModel) toEntity() *product.Category {
	return &product.Category{
		ID:             m.ID,
		TenantID:       shared.TenantID(m.TenantID),
		ParentID:       m.ParentID,
		Name:           m.Name,
		Code:           m.Code,
		Level:          m.Level,
		Sort:           m.Sort,
		Icon:           m.Icon,
		Image:          m.Image,
		SeoTitle:       m.SeoTitle,
		SeoDescription: m.SeoDescription,
		Status:         product.CategoryStatus(m.Status),
		Audit: shared.AuditInfo{
			CreatedAt: time.Unix(m.CreatedAt, 0),
			UpdatedAt: time.Unix(m.UpdatedAt, 0),
			CreatedBy: m.CreatedBy,
			UpdatedBy: m.UpdatedBy,
		},
	}
}

func fromCategoryEntity(c *product.Category) *categoryModel {
	now := time.Now().Unix()
	createdAt := now
	updatedAt := now
	if !c.Audit.CreatedAt.IsZero() {
		createdAt = c.Audit.CreatedAt.Unix()
	}
	if !c.Audit.UpdatedAt.IsZero() {
		updatedAt = c.Audit.UpdatedAt.Unix()
	}
	return &categoryModel{
		ID:             c.ID,
		TenantID:       c.TenantID.Int64(),
		ParentID:       c.ParentID,
		Name:           c.Name,
		Code:           c.Code,
		Level:          c.Level,
		Sort:           c.Sort,
		Icon:           c.Icon,
		Image:          c.Image,
		SeoTitle:       c.SeoTitle,
		SeoDescription: c.SeoDescription,
		Status:         int8(c.Status),
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		CreatedBy:      c.Audit.CreatedBy,
		UpdatedBy:      c.Audit.UpdatedBy,
	}
}

func (r *categoryRepo) Create(ctx context.Context, db *gorm.DB, c *product.Category) error {
	model := fromCategoryEntity(c)
	return db.WithContext(ctx).Create(model).Error
}

func (r *categoryRepo) Update(ctx context.Context, db *gorm.DB, c *product.Category) error {
	model := fromCategoryEntity(c)
	return db.WithContext(ctx).Model(&categoryModel{}).
		Where("id = ? AND tenant_id = ?", c.ID, c.TenantID.Int64()).
		Updates(map[string]interface{}{
			"name":            model.Name,
			"code":            model.Code,
			"parent_id":       model.ParentID,
			"level":           model.Level,
			"sort":            model.Sort,
			"icon":            model.Icon,
			"image":           model.Image,
			"seo_title":       model.SeoTitle,
			"seo_description": model.SeoDescription,
			"status":          model.Status,
			"updated_at":      model.UpdatedAt,
			"updated_by":      model.UpdatedBy,
		}).Error
}

func (r *categoryRepo) Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error {
	now := time.Now().Unix()
	return db.WithContext(ctx).Model(&categoryModel{}).
		Where("id = ? AND tenant_id = ?", id, tenantID.Int64()).
		Update("deleted_at", now).Error
}

func (r *categoryRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*product.Category, error) {
	var model categoryModel
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

func (r *categoryRepo) FindByParentID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, parentID int64) ([]*product.Category, error) {
	var models []categoryModel
	err := db.WithContext(ctx).
		Where("parent_id = ? AND tenant_id = ? AND deleted_at IS NULL", parentID, tenantID.Int64()).
		Order("sort ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	categories := make([]*product.Category, len(models))
	for i, m := range models {
		categories[i] = m.toEntity()
	}
	return categories, nil
}

func (r *categoryRepo) FindAll(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) ([]*product.Category, error) {
	var models []categoryModel
	err := db.WithContext(ctx).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID.Int64()).
		Order("sort ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	categories := make([]*product.Category, len(models))
	for i, m := range models {
		categories[i] = m.toEntity()
	}
	return categories, nil
}

func (r *categoryRepo) FindTree(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) ([]*product.Category, error) {
	return r.FindAll(ctx, db, tenantID)
}

func (r *categoryRepo) FindByCode(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, code string) (*product.Category, error) {
	var model categoryModel
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

func (r *categoryRepo) GetProductCount(ctx context.Context, db *gorm.DB, categoryID int64) (int64, error) {
	var count int64
	err := db.WithContext(ctx).Table("products").
		Where("category_id = ?", categoryID).
		Count(&count).Error
	return count, err
}

func (r *categoryRepo) UpdateSort(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, sorts []product.CategorySort) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, s := range sorts {
			if err := tx.Model(&categoryModel{}).
				Where("id = ? AND tenant_id = ?", s.ID, tenantID.Int64()).
				Update("sort", s.Sort).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *categoryRepo) Move(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64, newParentID int64) error {
	// Calculate new level
	var parentLevel int = 0
	if newParentID > 0 {
		var parent categoryModel
		if err := db.WithContext(ctx).
			Where("id = ? AND tenant_id = ?", newParentID, tenantID.Int64()).
			First(&parent).Error; err != nil {
			return err
		}
		parentLevel = parent.Level
	}

	return db.WithContext(ctx).Model(&categoryModel{}).
		Where("id = ? AND tenant_id = ?", id, tenantID.Int64()).
		Updates(map[string]interface{}{
			"parent_id": newParentID,
			"level":     parentLevel + 1,
		}).Error
}
