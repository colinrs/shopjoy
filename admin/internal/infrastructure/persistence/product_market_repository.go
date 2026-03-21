package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type productMarketRepo struct{}

func NewProductMarketRepository() product.ProductMarketRepository {
	return &productMarketRepo{}
}

type productMarketModel struct {
	ID                  int64          `gorm:"column:id;primaryKey;autoIncrement"`
	TenantID            int64          `gorm:"column:tenant_id;not null"`
	ProductID           int64          `gorm:"column:product_id;not null;uniqueIndex:uk_product_variant_market"`
	VariantID           *int64         `gorm:"column:variant_id;uniqueIndex:uk_product_variant_market"`
	MarketID            int64          `gorm:"column:market_id;not null;uniqueIndex:uk_product_variant_market"`
	IsEnabled           bool           `gorm:"column:is_enabled;not null;default:false"`
	StatusOverride      *int           `gorm:"column:status_override"`
	Price               string         `gorm:"column:price;type:decimal(19,4);not null"`
	CompareAtPrice      *string        `gorm:"column:compare_at_price;type:decimal(19,4)"`
	StockAlertThreshold int            `gorm:"column:stock_alert_threshold;not null;default:0"`
	PublishedAt         *time.Time     `gorm:"column:published_at"`
	CreatedAt           time.Time      `gorm:"column:created_at"`
	UpdatedAt           time.Time      `gorm:"column:updated_at"`
}

func (productMarketModel) TableName() string {
	return "product_markets"
}

func (m *productMarketModel) toEntity() *product.ProductMarket {
	pm := &product.ProductMarket{
		ID:                  m.ID,
		TenantID:            m.TenantID,
		ProductID:           m.ProductID,
		VariantID:           m.VariantID,
		MarketID:            m.MarketID,
		IsEnabled:           m.IsEnabled,
		StockAlertThreshold: m.StockAlertThreshold,
		PublishedAt:         m.PublishedAt,
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
	}

	pm.Price, _ = decimal.NewFromString(m.Price)
	if m.CompareAtPrice != nil {
		cap, _ := decimal.NewFromString(*m.CompareAtPrice)
		pm.CompareAtPrice = &cap
	}
	if m.StatusOverride != nil {
		status := product.Status(*m.StatusOverride)
		pm.StatusOverride = &status
	}

	return pm
}

func fromProductMarketEntity(pm *product.ProductMarket) *productMarketModel {
	m := &productMarketModel{
		ID:                  pm.ID,
		TenantID:            pm.TenantID,
		ProductID:           pm.ProductID,
		VariantID:           pm.VariantID,
		MarketID:            pm.MarketID,
		IsEnabled:           pm.IsEnabled,
		StockAlertThreshold: pm.StockAlertThreshold,
		PublishedAt:         pm.PublishedAt,
		CreatedAt:           pm.CreatedAt,
		UpdatedAt:           pm.UpdatedAt,
	}

	m.Price = pm.Price.String()
	if pm.CompareAtPrice != nil {
		cap := pm.CompareAtPrice.String()
		m.CompareAtPrice = &cap
	}
	if pm.StatusOverride != nil {
		status := int(*pm.StatusOverride)
		m.StatusOverride = &status
	}

	return m
}

func (r *productMarketRepo) Create(ctx context.Context, db *gorm.DB, pm *product.ProductMarket) error {
	model := fromProductMarketEntity(pm)
	return db.WithContext(ctx).Create(model).Error
}

func (r *productMarketRepo) Update(ctx context.Context, db *gorm.DB, pm *product.ProductMarket) error {
	model := fromProductMarketEntity(pm)
	return db.WithContext(ctx).Save(model).Error
}

func (r *productMarketRepo) Delete(ctx context.Context, db *gorm.DB, id int64) error {
	result := db.WithContext(ctx).Delete(&productMarketModel{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return code.ErrProductMarketNotFound
	}
	return nil
}

func (r *productMarketRepo) FindByID(ctx context.Context, db *gorm.DB, id int64) (*product.ProductMarket, error) {
	var model productMarketModel
	err := db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrProductMarketNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *productMarketRepo) FindByProductAndMarket(ctx context.Context, db *gorm.DB, productID, marketID int64, variantID *int64) (*product.ProductMarket, error) {
	var model productMarketModel
	query := db.WithContext(ctx).Where("product_id = ? AND market_id = ?", productID, marketID)
	if variantID != nil {
		query = query.Where("variant_id = ?", *variantID)
	} else {
		query = query.Where("variant_id IS NULL")
	}

	err := query.First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrProductMarketNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *productMarketRepo) FindByProductID(ctx context.Context, db *gorm.DB, productID int64) ([]*product.ProductMarket, error) {
	var models []productMarketModel
	err := db.WithContext(ctx).Where("product_id = ?", productID).Find(&models).Error
	if err != nil {
		return nil, err
	}

	pms := make([]*product.ProductMarket, len(models))
	for i, m := range models {
		pms[i] = m.toEntity()
	}
	return pms, nil
}

func (r *productMarketRepo) FindByProductIDs(ctx context.Context, db *gorm.DB, productIDs []int64) ([]*product.ProductMarket, error) {
	if len(productIDs) == 0 {
		return []*product.ProductMarket{}, nil
	}

	var models []productMarketModel
	err := db.WithContext(ctx).Where("product_id IN ?", productIDs).Find(&models).Error
	if err != nil {
		return nil, err
	}

	pms := make([]*product.ProductMarket, len(models))
	for i, m := range models {
		pms[i] = m.toEntity()
	}
	return pms, nil
}

func (r *productMarketRepo) FindByMarketID(ctx context.Context, db *gorm.DB, marketID int64, query product.Query) ([]*product.ProductMarket, int64, error) {
	dbQuery := db.WithContext(ctx).Model(&productMarketModel{}).Where("market_id = ?", marketID)

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []productMarketModel
	err := dbQuery.Order("created_at DESC").
		Offset(query.Offset()).
		Limit(query.Limit()).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	pms := make([]*product.ProductMarket, len(models))
	for i, m := range models {
		pms[i] = m.toEntity()
	}
	return pms, total, nil
}

func (r *productMarketRepo) BatchCreate(ctx context.Context, db *gorm.DB, pms []*product.ProductMarket) error {
	models := make([]*productMarketModel, len(pms))
	for i, pm := range pms {
		models[i] = fromProductMarketEntity(pm)
	}
	return db.WithContext(ctx).Create(&models).Error
}