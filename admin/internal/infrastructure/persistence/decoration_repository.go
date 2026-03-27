package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/storefront"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type decorationRepo struct{}

func NewDecorationRepository() storefront.DecorationRepository {
	return &decorationRepo{}
}

type decorationModel struct {
	ID          int64     `gorm:"column:id;primaryKey"`
	TenantID    int64     `gorm:"column:tenant_id;not null;index"`
	PageID      int64     `gorm:"column:page_id;not null;index"`
	BlockType   string    `gorm:"column:block_type;type:varchar(50);not null;index"`
	BlockConfig string    `gorm:"column:block_config;type:text;not null"`
	SortOrder   int       `gorm:"column:sort_order;not null"`
	IsActive    int       `gorm:"column:is_active;not null;default:1"`
	CreatedAt   time.Time `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null"`
}

func (decorationModel) TableName() string {
	return "decorations"
}

func (m *decorationModel) toEntity() *storefront.Decoration {
	var config storefront.BlockConfig
	if m.BlockConfig != "" {
		json.Unmarshal([]byte(m.BlockConfig), &config)
	}
	return &storefront.Decoration{
		Model:       application.Model{ID: m.ID},
		TenantID:    shared.TenantID(m.TenantID),
		PageID:      m.PageID,
		BlockType:   m.BlockType,
		BlockConfig: config,
		SortOrder:   m.SortOrder,
		IsActive:    m.IsActive == 1,
	}
}

func fromDecorationEntity(d *storefront.Decoration) *decorationModel {
	now := time.Now().UTC()
	config, _ := json.Marshal(d.BlockConfig)

	isActive := 0
	if d.IsActive {
		isActive = 1
	}

	return &decorationModel{
		ID:          d.ID,
		TenantID:    d.TenantID.Int64(),
		PageID:      d.PageID,
		BlockType:   d.BlockType,
		BlockConfig: string(config),
		SortOrder:   d.SortOrder,
		IsActive:    isActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (r *decorationRepo) Create(ctx context.Context, db *gorm.DB, d *storefront.Decoration) error {
	model := fromDecorationEntity(d)
	return db.WithContext(ctx).Create(model).Error
}

func (r *decorationRepo) Update(ctx context.Context, db *gorm.DB, d *storefront.Decoration) error {
	model := fromDecorationEntity(d)
	config, _ := json.Marshal(d.BlockConfig)

	return db.WithContext(ctx).Model(&decorationModel{}).
		Where("id = ? AND tenant_id = ?", d.ID, d.TenantID.Int64()).
		Updates(map[string]interface{}{
			"block_type":   model.BlockType,
			"block_config": string(config),
			"sort_order":   model.SortOrder,
			"is_active":    model.IsActive,
			"updated_at":   model.UpdatedAt,
		}).Error
}

func (r *decorationRepo) Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error {
	return db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", id, tenantID.Int64()).
		Delete(&decorationModel{}).Error
}

func (r *decorationRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*storefront.Decoration, error) {
	var model decorationModel
	err := db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", id, tenantID.Int64()).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *decorationRepo) FindByPageID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, pageID int64) ([]*storefront.Decoration, error) {
	var models []decorationModel
	err := db.WithContext(ctx).
		Where("page_id = ? AND tenant_id = ? AND is_active = 1", pageID, tenantID.Int64()).
		Order("sort_order ASC, id ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	decorations := make([]*storefront.Decoration, len(models))
	for i, m := range models {
		decorations[i] = m.toEntity()
	}
	return decorations, nil
}

func (r *decorationRepo) Reorder(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orders []storefront.BlockOrder) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, order := range orders {
			if err := tx.Model(&decorationModel{}).
				Where("id = ? AND tenant_id = ?", order.ID, tenantID.Int64()).
				Update("sort_order", order.SortOrder).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *decorationRepo) DeleteByPageID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, pageID int64) error {
	return db.WithContext(ctx).
		Where("page_id = ? AND tenant_id = ?", pageID, tenantID.Int64()).
		Delete(&decorationModel{}).Error
}