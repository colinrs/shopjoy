package persistence

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type categoryMarketRepo struct{}

func NewCategoryMarketRepository() product.CategoryMarketRepository {
	return &categoryMarketRepo{}
}

type categoryMarketModel struct {
	ID         int64 `gorm:"column:id;primaryKey"`
	TenantID   int64 `gorm:"column:tenant_id;not null;index"`
	CategoryID int64 `gorm:"column:category_id;not null;index"`
	MarketID   int64 `gorm:"column:market_id;not null;index"`
	IsVisible  bool  `gorm:"column:is_visible;not null;default:true"`
	CreatedAt  int64 `gorm:"column:created_at;not null"`
	UpdatedAt  int64 `gorm:"column:updated_at;not null"`
}

func (categoryMarketModel) TableName() string {
	return "category_markets"
}

func (m *categoryMarketModel) toEntity() *product.CategoryMarket {
	return &product.CategoryMarket{
		ID:         m.ID,
		TenantID:   shared.TenantID(m.TenantID),
		CategoryID: m.CategoryID,
		MarketID:   m.MarketID,
		IsVisible:  m.IsVisible,
		Audit: shared.AuditInfo{
			CreatedAt: time.Unix(m.CreatedAt, 0),
			UpdatedAt: time.Unix(m.UpdatedAt, 0),
		},
	}
}

func fromCategoryMarketEntity(cm *product.CategoryMarket) *categoryMarketModel {
	now := time.Now().Unix()
	createdAt := now
	updatedAt := now
	if !cm.Audit.CreatedAt.IsZero() {
		createdAt = cm.Audit.CreatedAt.Unix()
	}
	if !cm.Audit.UpdatedAt.IsZero() {
		updatedAt = cm.Audit.UpdatedAt.Unix()
	}
	return &categoryMarketModel{
		ID:         cm.ID,
		TenantID:   cm.TenantID.Int64(),
		CategoryID: cm.CategoryID,
		MarketID:   cm.MarketID,
		IsVisible:  cm.IsVisible,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}

func (r *categoryMarketRepo) Create(ctx context.Context, db *gorm.DB, cm *product.CategoryMarket) error {
	model := fromCategoryMarketEntity(cm)
	return db.WithContext(ctx).Create(model).Error
}

func (r *categoryMarketRepo) Update(ctx context.Context, db *gorm.DB, cm *product.CategoryMarket) error {
	return db.WithContext(ctx).Model(&categoryMarketModel{}).
		Where("id = ?", cm.ID).
		Updates(map[string]interface{}{
			"is_visible": cm.IsVisible,
			"updated_at": time.Now().Unix(),
		}).Error
}

func (r *categoryMarketRepo) FindByCategory(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, categoryID int64) ([]*product.CategoryMarket, error) {
	var models []categoryMarketModel
	err := db.WithContext(ctx).
		Where("category_id = ? AND tenant_id = ?", categoryID, tenantID.Int64()).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	result := make([]*product.CategoryMarket, len(models))
	for i, m := range models {
		result[i] = m.toEntity()
	}
	return result, nil
}

func (r *categoryMarketRepo) DeleteByCategory(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, categoryID int64) error {
	return db.WithContext(ctx).
		Where("category_id = ? AND tenant_id = ?", categoryID, tenantID.Int64()).
		Delete(&categoryMarketModel{}).Error
}

func (r *categoryMarketRepo) BatchCreate(ctx context.Context, db *gorm.DB, items []*product.CategoryMarket) error {
	if len(items) == 0 {
		return nil
	}
	models := make([]*categoryMarketModel, len(items))
	for i, item := range items {
		models[i] = fromCategoryMarketEntity(item)
	}
	return db.WithContext(ctx).CreateInBatches(models, 100).Error
}