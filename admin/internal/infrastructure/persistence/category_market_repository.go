package persistence

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type categoryMarketRepo struct{}

func NewCategoryMarketRepository() product.CategoryMarketRepository {
	return &categoryMarketRepo{}
}

type categoryMarketModel struct {
	ID         int64     `gorm:"column:id;primaryKey"`
	TenantID   int64     `gorm:"column:tenant_id;not null;index"`
	CategoryID int64     `gorm:"column:category_id;not null;index"`
	MarketID   int64     `gorm:"column:market_id;not null;index"`
	IsVisible  bool      `gorm:"column:is_visible;not null;default:true"`
	CreatedAt  time.Time `gorm:"column:created_at;not null"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null"`
}

func (categoryMarketModel) TableName() string {
	return "category_markets"
}

func (m *categoryMarketModel) toEntity() *product.CategoryMarket {
	return &product.CategoryMarket{
		Model:      application.Model{ID: m.ID, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt},
		TenantID:   shared.TenantID(m.TenantID),
		CategoryID: m.CategoryID,
		MarketID:   m.MarketID,
		IsVisible:  m.IsVisible,
	}
}

func fromCategoryMarketEntity(cm *product.CategoryMarket) *categoryMarketModel {
	return &categoryMarketModel{
		ID:         cm.Model.ID,
		TenantID:   cm.TenantID.Int64(),
		CategoryID: cm.CategoryID,
		MarketID:   cm.MarketID,
		IsVisible:  cm.IsVisible,
		CreatedAt:  cm.Model.CreatedAt,
		UpdatedAt:  cm.Model.UpdatedAt,
	}
}

func (r *categoryMarketRepo) Create(ctx context.Context, db *gorm.DB, cm *product.CategoryMarket) error {
	model := fromCategoryMarketEntity(cm)
	return db.WithContext(ctx).Create(model).Error
}

func (r *categoryMarketRepo) Update(ctx context.Context, db *gorm.DB, cm *product.CategoryMarket) error {
	return db.WithContext(ctx).Model(&categoryMarketModel{}).
		Where("id = ?", cm.Model.ID).
		Updates(map[string]interface{}{
			"is_visible": cm.IsVisible,
			"updated_at": time.Now().UTC(),
		}).Error
}

func (r *categoryMarketRepo) FindByCategory(ctx context.Context, db *gorm.DB,  categoryID int64) ([]*product.CategoryMarket, error) {
	var models []categoryMarketModel
	err := db.WithContext(ctx).
		Where("category_id = ?", categoryID).
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

func (r *categoryMarketRepo) DeleteByCategory(ctx context.Context, db *gorm.DB,  categoryID int64) error {
	return db.WithContext(ctx).
		Where("category_id = ?", categoryID).
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
