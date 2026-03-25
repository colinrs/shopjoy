package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/shipping"
	"github.com/colinrs/shopjoy/pkg/code"
	"gorm.io/gorm"
)

// ShippingTemplateRepository 运费模板仓储接口
type ShippingTemplateRepository interface {
	// Template operations
	Create(ctx context.Context, db *gorm.DB, template *shipping.ShippingTemplate) error
	Update(ctx context.Context, db *gorm.DB, template *shipping.ShippingTemplate) error
	Delete(ctx context.Context, db *gorm.DB, tenantID, id int64) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID, id int64) (*shipping.ShippingTemplate, error)
	FindByIDWithDetails(ctx context.Context, db *gorm.DB, tenantID, id int64) (*shipping.ShippingTemplate, []*shipping.ShippingZone, []*shipping.ShippingTemplateMapping, error)
	FindList(ctx context.Context, db *gorm.DB, tenantID int64, name string, isActive *bool, page, pageSize int) ([]*shipping.ShippingTemplate, int64, error)
	FindDefault(ctx context.Context, db *gorm.DB, tenantID int64) (*shipping.ShippingTemplate, error)
	SetDefault(ctx context.Context, db *gorm.DB, tenantID, id int64) error
	UnsetAllDefault(ctx context.Context, db *gorm.DB, tenantID int64) error

	// Zone operations
	CreateZone(ctx context.Context, db *gorm.DB, zone *shipping.ShippingZone) error
	UpdateZone(ctx context.Context, db *gorm.DB, zone *shipping.ShippingZone) error
	DeleteZone(ctx context.Context, db *gorm.DB, id int64) error
	FindZoneByID(ctx context.Context, db *gorm.DB, id int64) (*shipping.ShippingZone, error)
	FindZonesByTemplateID(ctx context.Context, db *gorm.DB, templateID int64) ([]*shipping.ShippingZone, error)
	ReorderZones(ctx context.Context, db *gorm.DB, templateID int64, zoneIDs []int64) error
	FindZoneByCityCode(ctx context.Context, db *gorm.DB, tenantID int64, cityCode string) ([]*shipping.ShippingZone, error)

	// Mapping operations
	CreateMapping(ctx context.Context, db *gorm.DB, mapping *shipping.ShippingTemplateMapping) error
	DeleteMapping(ctx context.Context, db *gorm.DB, id int64) error
	FindMappingByID(ctx context.Context, db *gorm.DB, id int64) (*shipping.ShippingTemplateMapping, error)
	FindMappingsByTemplateID(ctx context.Context, db *gorm.DB, templateID int64) ([]*shipping.ShippingTemplateMapping, error)
	FindMappingByTarget(ctx context.Context, db *gorm.DB, targetType shipping.TargetType, targetID int64) (*shipping.ShippingTemplateMapping, error)

	// Statistics
	CountZonesByTemplateID(ctx context.Context, db *gorm.DB, templateID int64) (int64, error)
	CountProductsByTemplateID(ctx context.Context, db *gorm.DB, templateID int64) (int64, error)
	CountCategoriesByTemplateID(ctx context.Context, db *gorm.DB, templateID int64) (int64, error)
}

// shippingTemplateRepo 运费模板仓储实现
type shippingTemplateRepo struct{}

// NewShippingTemplateRepository 创建运费模板仓储
func NewShippingTemplateRepository() ShippingTemplateRepository {
	return &shippingTemplateRepo{}
}

// Create 创建运费模板
func (r *shippingTemplateRepo) Create(ctx context.Context, db *gorm.DB, template *shipping.ShippingTemplate) error {
	now := time.Now().UTC()
	template.CreatedAt = now
	template.UpdatedAt = now
	return db.WithContext(ctx).Create(template).Error
}

// Update 更新运费模板
func (r *shippingTemplateRepo) Update(ctx context.Context, db *gorm.DB, template *shipping.ShippingTemplate) error {
	template.UpdatedAt = time.Now().UTC()
	return db.WithContext(ctx).Save(template).Error
}

// Delete 删除运费模板
func (r *shippingTemplateRepo) Delete(ctx context.Context, db *gorm.DB, tenantID, id int64) error {
	result := db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Delete(&shipping.ShippingTemplate{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return code.ErrShippingTemplateNotFound
	}
	return nil
}

// FindByID 根据ID查找运费模板
func (r *shippingTemplateRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID, id int64) (*shipping.ShippingTemplate, error) {
	var template shipping.ShippingTemplate
	err := db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&template).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrShippingTemplateNotFound
		}
		return nil, err
	}
	return &template, nil
}

// FindByIDWithDetails 根据ID查找运费模板详情（包含区域和关联）
func (r *shippingTemplateRepo) FindByIDWithDetails(ctx context.Context, db *gorm.DB, tenantID, id int64) (*shipping.ShippingTemplate, []*shipping.ShippingZone, []*shipping.ShippingTemplateMapping, error) {
	template, err := r.FindByID(ctx, db, tenantID, id)
	if err != nil {
		return nil, nil, nil, err
	}

	zones, err := r.FindZonesByTemplateID(ctx, db, id)
	if err != nil {
		return nil, nil, nil, err
	}

	mappings, err := r.FindMappingsByTemplateID(ctx, db, id)
	if err != nil {
		return nil, nil, nil, err
	}

	return template, zones, mappings, nil
}

// FindList 查找运费模板列表
func (r *shippingTemplateRepo) FindList(ctx context.Context, db *gorm.DB, tenantID int64, name string, isActive *bool, page, pageSize int) ([]*shipping.ShippingTemplate, int64, error) {
	query := db.WithContext(ctx).Model(&shipping.ShippingTemplate{}).
		Where("tenant_id = ?", tenantID)

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var templates []*shipping.ShippingTemplate
	offset := (page - 1) * pageSize
	err := query.Order("is_default DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&templates).Error
	if err != nil {
		return nil, 0, err
	}

	return templates, total, nil
}

// FindDefault 查找默认运费模板
func (r *shippingTemplateRepo) FindDefault(ctx context.Context, db *gorm.DB, tenantID int64) (*shipping.ShippingTemplate, error) {
	var template shipping.ShippingTemplate
	err := db.WithContext(ctx).
		Where("tenant_id = ? AND is_default = ? AND is_active = ?", tenantID, true, true).
		First(&template).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 没有默认模板不是错误
		}
		return nil, err
	}
	return &template, nil
}

// SetDefault 设置默认模板
func (r *shippingTemplateRepo) SetDefault(ctx context.Context, db *gorm.DB, tenantID, id int64) error {
	return db.WithContext(ctx).Model(&shipping.ShippingTemplate{}).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Update("is_default", true).Error
}

// UnsetAllDefault 取消所有默认模板
func (r *shippingTemplateRepo) UnsetAllDefault(ctx context.Context, db *gorm.DB, tenantID int64) error {
	return db.WithContext(ctx).Model(&shipping.ShippingTemplate{}).
		Where("tenant_id = ?", tenantID).
		Update("is_default", false).Error
}

// CreateZone 创建配送区域
func (r *shippingTemplateRepo) CreateZone(ctx context.Context, db *gorm.DB, zone *shipping.ShippingZone) error {
	now := time.Now().UTC()
	zone.CreatedAt = now
	zone.UpdatedAt = now
	return db.WithContext(ctx).Create(zone).Error
}

// UpdateZone 更新配送区域
func (r *shippingTemplateRepo) UpdateZone(ctx context.Context, db *gorm.DB, zone *shipping.ShippingZone) error {
	zone.UpdatedAt = time.Now().UTC()
	return db.WithContext(ctx).Save(zone).Error
}

// DeleteZone 删除配送区域
func (r *shippingTemplateRepo) DeleteZone(ctx context.Context, db *gorm.DB, id int64) error {
	return db.WithContext(ctx).Delete(&shipping.ShippingZone{}, id).Error
}

// FindZoneByID 根据ID查找配送区域
func (r *shippingTemplateRepo) FindZoneByID(ctx context.Context, db *gorm.DB, id int64) (*shipping.ShippingZone, error) {
	var zone shipping.ShippingZone
	err := db.WithContext(ctx).First(&zone, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrShippingZoneNotFound
		}
		return nil, err
	}
	return &zone, nil
}

// FindZonesByTemplateID 根据模板ID查找配送区域列表
func (r *shippingTemplateRepo) FindZonesByTemplateID(ctx context.Context, db *gorm.DB, templateID int64) ([]*shipping.ShippingZone, error) {
	var zones []*shipping.ShippingZone
	err := db.WithContext(ctx).
		Where("template_id = ?", templateID).
		Order("sort ASC, id ASC").
		Find(&zones).Error
	return zones, err
}

// ReorderZones 重新排序配送区域
func (r *shippingTemplateRepo) ReorderZones(ctx context.Context, db *gorm.DB, templateID int64, zoneIDs []int64) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, zoneID := range zoneIDs {
			if err := tx.Model(&shipping.ShippingZone{}).
				Where("id = ? AND template_id = ?", zoneID, templateID).
				Update("sort", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// FindZoneByCityCode 根据城市代码查找匹配的配送区域
func (r *shippingTemplateRepo) FindZoneByCityCode(ctx context.Context, db *gorm.DB, tenantID int64, cityCode string) ([]*shipping.ShippingZone, error) {
	var zones []*shipping.ShippingZone
	err := db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("JSON_CONTAINS(regions, ?)", `"`+cityCode+`"`).
		Order("sort ASC").
		Find(&zones).Error
	return zones, err
}

// CreateMapping 创建模板关联
func (r *shippingTemplateRepo) CreateMapping(ctx context.Context, db *gorm.DB, mapping *shipping.ShippingTemplateMapping) error {
	now := time.Now().UTC()
	mapping.CreatedAt = now
	mapping.UpdatedAt = now
	return db.WithContext(ctx).Create(mapping).Error
}

// DeleteMapping 删除模板关联
func (r *shippingTemplateRepo) DeleteMapping(ctx context.Context, db *gorm.DB, id int64) error {
	return db.WithContext(ctx).Delete(&shipping.ShippingTemplateMapping{}, id).Error
}

// FindMappingByID 根据ID查找模板关联
func (r *shippingTemplateRepo) FindMappingByID(ctx context.Context, db *gorm.DB, id int64) (*shipping.ShippingTemplateMapping, error) {
	var mapping shipping.ShippingTemplateMapping
	err := db.WithContext(ctx).First(&mapping, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrShippingMappingNotFound
		}
		return nil, err
	}
	return &mapping, nil
}

// FindMappingsByTemplateID 根据模板ID查找关联列表
func (r *shippingTemplateRepo) FindMappingsByTemplateID(ctx context.Context, db *gorm.DB, templateID int64) ([]*shipping.ShippingTemplateMapping, error) {
	var mappings []*shipping.ShippingTemplateMapping
	err := db.WithContext(ctx).
		Where("template_id = ?", templateID).
		Find(&mappings).Error
	return mappings, err
}

// FindMappingByTarget 根据目标查找模板关联
func (r *shippingTemplateRepo) FindMappingByTarget(ctx context.Context, db *gorm.DB, targetType shipping.TargetType, targetID int64) (*shipping.ShippingTemplateMapping, error) {
	var mapping shipping.ShippingTemplateMapping
	err := db.WithContext(ctx).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		First(&mapping).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 没有关联不是错误
		}
		return nil, err
	}
	return &mapping, nil
}

// CountZonesByTemplateID 统计模板的配送区域数量
func (r *shippingTemplateRepo) CountZonesByTemplateID(ctx context.Context, db *gorm.DB, templateID int64) (int64, error) {
	var count int64
	err := db.WithContext(ctx).
		Model(&shipping.ShippingZone{}).
		Where("template_id = ?", templateID).
		Count(&count).Error
	return count, err
}

// CountProductsByTemplateID 统计模板关联的商品数量
func (r *shippingTemplateRepo) CountProductsByTemplateID(ctx context.Context, db *gorm.DB, templateID int64) (int64, error) {
	var count int64
	err := db.WithContext(ctx).
		Model(&shipping.ShippingTemplateMapping{}).
		Where("template_id = ? AND target_type = ?", templateID, shipping.TargetTypeProduct).
		Count(&count).Error
	return count, err
}

// CountCategoriesByTemplateID 统计模板关联的分类数量
func (r *shippingTemplateRepo) CountCategoriesByTemplateID(ctx context.Context, db *gorm.DB, templateID int64) (int64, error) {
	var count int64
	err := db.WithContext(ctx).
		Model(&shipping.ShippingTemplateMapping{}).
		Where("template_id = ? AND target_type = ?", templateID, shipping.TargetTypeCategory).
		Count(&count).Error
	return count, err
}