package persistence

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type brandMarketRepo struct{}

func NewBrandMarketRepository() product.BrandMarketRepository {
	return &brandMarketRepo{}
}

type brandMarketModel struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	TenantID  int64     `gorm:"column:tenant_id;not null;index"`
	BrandID   int64     `gorm:"column:brand_id;not null;index"`
	MarketID  int64     `gorm:"column:market_id;not null;index"`
	IsVisible bool      `gorm:"column:is_visible;not null;default:true"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
}

func (brandMarketModel) TableName() string {
	return "brand_markets"
}

func (m *brandMarketModel) toEntity() *product.BrandMarket {
	return &product.BrandMarket{
		Model:     application.Model{ID: m.ID, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt},
		TenantID:  shared.TenantID(m.TenantID),
		BrandID:   m.BrandID,
		MarketID:  m.MarketID,
		IsVisible: m.IsVisible,
	}
}

func fromBrandMarketEntity(bm *product.BrandMarket) *brandMarketModel {
	return &brandMarketModel{
		ID:        bm.Model.ID,
		TenantID:  bm.TenantID.Int64(),
		BrandID:   bm.BrandID,
		MarketID:  bm.MarketID,
		IsVisible: bm.IsVisible,
		CreatedAt: bm.Model.CreatedAt,
		UpdatedAt: bm.Model.UpdatedAt,
	}
}

func (r *brandMarketRepo) Create(ctx context.Context, db *gorm.DB, bm *product.BrandMarket) error {
	model := fromBrandMarketEntity(bm)
	return db.WithContext(ctx).Create(model).Error
}

func (r *brandMarketRepo) Update(ctx context.Context, db *gorm.DB, bm *product.BrandMarket) error {
	return db.WithContext(ctx).Model(&brandMarketModel{}).
		Where("id = ?", bm.Model.ID).
		Updates(map[string]any{
			"is_visible": bm.IsVisible,
			"updated_at": time.Now().UTC(),
		}).Error
}

func (r *brandMarketRepo) FindByBrand(ctx context.Context, db *gorm.DB, brandID int64) ([]*product.BrandMarket, error) {
	var models []brandMarketModel
	err := db.WithContext(ctx).
		Where("brand_id = ?", brandID).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	result := make([]*product.BrandMarket, len(models))
	for i, m := range models {
		result[i] = m.toEntity()
	}
	return result, nil
}

func (r *brandMarketRepo) DeleteByBrand(ctx context.Context, db *gorm.DB, brandID int64) error {
	return db.WithContext(ctx).
		Where("brand_id = ?", brandID).
		Delete(&brandMarketModel{}).Error
}

func (r *brandMarketRepo) BatchCreate(ctx context.Context, db *gorm.DB, items []*product.BrandMarket) error {
	if len(items) == 0 {
		return nil
	}
	models := make([]*brandMarketModel, len(items))
	for i, item := range items {
		models[i] = fromBrandMarketEntity(item)
	}
	return db.WithContext(ctx).CreateInBatches(models, 100).Error
}
