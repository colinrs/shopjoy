package persistence

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/region"
	"github.com/colinrs/shopjoy/pkg/application"
	"gorm.io/gorm"
)

type regionRepo struct{}

// NewRegionRepository 创建区域仓储实现
func NewRegionRepository() region.RegionRepository {
	return &regionRepo{}
}

// regionModel 数据库模型
type regionModel struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement:false"`
	TenantID      int64     `gorm:"column:tenant_id;not null;default:0"`
	Code          string    `gorm:"column:code;size:20;not null"`
	Name          string    `gorm:"column:name;size:100;not null"`
	Level         int       `gorm:"column:level;not null;default:1"`
	ParentCode    string    `gorm:"column:parent_code;size:20;not null;default:''"`
	CountryCode   string    `gorm:"column:country_code;size:2;not null;default:'CN'"`
	PostalPattern string    `gorm:"column:postal_pattern;size:255;not null;default:''"`
	Sort          int       `gorm:"column:sort;not null;default:0"`
	IsActive      int       `gorm:"column:is_active;not null;default:1"`
	CreatedAt     time.Time `gorm:"column:created_at;not null"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null"`
}

func (regionModel) TableName() string {
	return "regions"
}

func (m *regionModel) toEntity() *region.Region {
	return &region.Region{
		Model:         application.Model{ID: m.ID, CreatedAt: m.CreatedAt.UTC(), UpdatedAt: m.UpdatedAt.UTC()},
		TenantID:      m.TenantID,
		Code:          m.Code,
		Name:          m.Name,
		Level:         m.Level,
		ParentCode:    m.ParentCode,
		CountryCode:   m.CountryCode,
		PostalPattern: m.PostalPattern,
		Sort:          m.Sort,
		IsActive:      m.IsActive == 1,
	}
}

// FindTree 按查询条件返回扁平区域列表（按 sort, code 排序）。
//
// 行为说明：
//   - TenantID=0：仅查询平台预置数据
//   - CountryCode 非空：按 ISO 国家码过滤
//   - ParentCode 非空：返回该父级下的所有子级（不限层级）
//   - ParentCode 空：返回顶级（Level=1，即国家列表）
//   - Level > 0：在 ParentCode 过滤基础上再按 level 过滤
func (r *regionRepo) FindTree(ctx context.Context, db *gorm.DB, query region.RegionQuery) ([]*region.Region, error) {
	tx := db.WithContext(ctx).Model(&regionModel{})

	// 默认仅查平台预置数据（tenant_id=0），除非显式传入其他值
	tx = tx.Where("tenant_id = ?", query.TenantID)

	if query.CountryCode != "" {
		tx = tx.Where("country_code = ?", query.CountryCode)
	}

	if query.ParentCode != "" {
		tx = tx.Where("parent_code = ?", query.ParentCode)
	} else if query.Level <= 0 {
		// 无父级且无 level 过滤时，默认返回顶级（国家）
		tx = tx.Where("level = ?", region.RegionLevelCountry)
	}

	if query.Level > 0 && query.ParentCode != "" {
		tx = tx.Where("level = ?", query.Level)
	}

	if query.IsActive != nil {
		isActive := 0
		if *query.IsActive {
			isActive = 1
		}
		tx = tx.Where("is_active = ?", isActive)
	}

	var models []regionModel
	if err := tx.Order("sort ASC, code ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	entities := make([]*region.Region, 0, len(models))
	for i := range models {
		entities = append(entities, models[i].toEntity())
	}
	return entities, nil
}
