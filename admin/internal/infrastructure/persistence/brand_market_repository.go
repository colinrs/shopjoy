package persistence

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type brandMarketRepo struct{}

func NewBrandMarketRepository() product.BrandMarketRepository {
	return &brandMarketRepo{}
}

type brandMarketModel struct {
	ID        int64 `gorm:"column:id;primaryKey"`
	TenantID  int64 `gorm:"column:tenant_id;not null;index"`
	BrandID   int64 `gorm:"column:brand_id;not null;index"`
	MarketID  int64 `gorm:"column:market_id;not null;index"`
	IsVisible bool  `gorm:"column:is_visible;not null;default:true"`
	CreatedAt int64 `gorm:"column:created_at;not null"`
	UpdatedAt int64 `gorm:"column:updated_at;not null"`
}

func (brandMarketModel) TableName() string {
	return "brand_markets"
}

func (m *brandMarketModel) toEntity() *product.BrandMarket {
	return &product.BrandMarket{
		ID:        m.ID,
		TenantID:  shared.TenantID(m.TenantID),
		BrandID:   m.BrandID,
		MarketID:  m.MarketID,
		IsVisible: m.IsVisible,
		Audit: shared.AuditInfo{
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		},
	}
}

func fromBrandMarketEntity(bm *product.BrandMarket) *brandMarketModel {
	return &brandMarketModel{
		ID:        bm.ID,
		TenantID:  bm.TenantID.Int64(),
		BrandID:   bm.BrandID,
		MarketID:  bm.MarketID,
		IsVisible: bm.IsVisible,
		CreatedAt: bm.Audit.CreatedAt,
		UpdatedAt: bm.Audit.UpdatedAt,
	}
}

func (r *brandMarketRepo) Create(ctx context.Context, db *gorm.DB, bm *product.BrandMarket) error {
	model := fromBrandMarketEntity(bm)
	return db.WithContext(ctx).Create(model).Error
}

func (r *brandMarketRepo) Update(ctx context.Context, db *gorm.DB, bm *product.BrandMarket) error {
	return db.WithContext(ctx).Model(&brandMarketModel{}).
		Where("id = ?", bm.ID).
		Updates(map[string]any{
			"is_visible": bm.IsVisible,
			"updated_at": time.Now().Unix(),
		}).Error
}

func (r *brandMarketRepo) FindByBrand(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, brandID int64) ([]*product.BrandMarket, error) {
	var models []brandMarketModel
	err := db.WithContext(ctx).
		Where("brand_id = ? AND tenant_id = ?", brandID, tenantID.Int64()).
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

func (r *brandMarketRepo) DeleteByBrand(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, brandID int64) error {
	return db.WithContext(ctx).
		Where("brand_id = ? AND tenant_id = ?", brandID, tenantID.Int64()).
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
