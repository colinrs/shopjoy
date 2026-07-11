package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type skuModel struct {
	ID             int64           `gorm:"column:id;primaryKey;autoIncrement:false"`
	TenantID       int64           `gorm:"column:tenant_id;not null;index:idx_sku_tenant"`
	ProductID      int64           `gorm:"column:product_id;not null;index:idx_sku_product"`
	Code           string          `gorm:"column:code;size:64;not null;uniqueIndex:idx_sku_code"`
	PriceAmount    decimal.Decimal `gorm:"column:price_amount;type:decimal(19,4);not null"`
	PriceCurrency  string          `gorm:"column:price_currency;size:3;not null"`
	Stock          int             `gorm:"column:stock;not null;default:0"`
	AvailableStock int             `gorm:"column:available_stock;not null;default:0"`
	LockedStock    int             `gorm:"column:locked_stock;not null;default:0"`
	SafetyStock    int             `gorm:"column:safety_stock;not null;default:0"`
	PreSaleEnabled bool            `gorm:"column:pre_sale_enabled;not null;default:false"`
	Attributes     string          `gorm:"column:attributes;type:text"`
	Status         int             `gorm:"column:status;not null;default:1"`
	CreatedAt      time.Time       `gorm:"column:created_at;not null"`
	UpdatedAt      time.Time       `gorm:"column:updated_at;not null"`
}

func (m *skuModel) TableName() string {
	return "skus"
}

func (m *skuModel) toEntity() *product.SKU {
	var attributes map[string]string
	if m.Attributes != "" {
		_ = json.Unmarshal([]byte(m.Attributes), &attributes)
	}

	return &product.SKU{
		Model: application.Model{
			ID:        m.ID,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		},
		TenantID:       shared.TenantID(m.TenantID),
		ProductID:      m.ProductID,
		Code:           m.Code,
		Price:          shared.NewMoney(m.PriceAmount, m.PriceCurrency),
		Stock:          m.Stock,
		AvailableStock: m.AvailableStock,
		LockedStock:    m.LockedStock,
		SafetyStock:    m.SafetyStock,
		PreSaleEnabled: m.PreSaleEnabled,
		Attributes:     attributes,
		Status:         shared.Status(m.Status),
		Audit: shared.AuditInfo{
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		},
	}
}

func fromSKUEntity(sku *product.SKU) *skuModel {
	attributesJSON, _ := json.Marshal(sku.Attributes)

	return &skuModel{
		ID:             sku.Model.ID,
		TenantID:       sku.TenantID.Int64(),
		ProductID:      sku.ProductID,
		Code:           sku.Code,
		PriceAmount:    sku.Price.Amount,
		PriceCurrency:  sku.Price.Currency,
		Stock:          sku.Stock,
		AvailableStock: sku.AvailableStock,
		LockedStock:    sku.LockedStock,
		SafetyStock:    sku.SafetyStock,
		PreSaleEnabled: sku.PreSaleEnabled,
		Attributes:     string(attributesJSON),
		Status:         int(sku.Status),
		CreatedAt:      sku.Model.CreatedAt,
		UpdatedAt:      sku.Model.UpdatedAt,
	}
}

type skuRepo struct{}

func NewSKURepository() product.SKURepository {
	return &skuRepo{}
}

func (r *skuRepo) Create(ctx context.Context, db *gorm.DB, sku *product.SKU) error {
	model := fromSKUEntity(sku)
	return db.WithContext(ctx).Create(model).Error
}

func (r *skuRepo) Update(ctx context.Context, db *gorm.DB, sku *product.SKU) error {
	model := fromSKUEntity(sku)
	return db.WithContext(ctx).
		Model(&skuModel{}).
		Where("id = ?", sku.Model.ID).
		Updates(map[string]interface{}{
			"code":             model.Code,
			"price_amount":     model.PriceAmount,
			"price_currency":   model.PriceCurrency,
			"stock":            model.Stock,
			"available_stock":  model.AvailableStock,
			"locked_stock":     model.LockedStock,
			"safety_stock":     model.SafetyStock,
			"pre_sale_enabled": model.PreSaleEnabled,
			"attributes":       model.Attributes,
			"status":           model.Status,
			"updated_at":       model.UpdatedAt,
		}).Error
}

func (r *skuRepo) Delete(ctx context.Context, db *gorm.DB,  id int64) error {
	result := db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&skuModel{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return code.ErrProductNotFound
	}
	return nil
}

func (r *skuRepo) FindByID(ctx context.Context, db *gorm.DB,  id int64) (*product.SKU, error) {
	var model skuModel
	err := db.WithContext(ctx).
		Where("id = ?", id).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrProductNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *skuRepo) FindByCode(ctx context.Context, db *gorm.DB,  skuCode string) (*product.SKU, error) {
	var model skuModel
	err := db.WithContext(ctx).
		Where("code = ?", skuCode).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrProductNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *skuRepo) FindByProductID(ctx context.Context, db *gorm.DB,  productID int64) ([]*product.SKU, error) {
	var models []skuModel
	query := db.WithContext(ctx).Where("product_id = ?", productID)

	// Platform admin (tenantID == 0) can access all tenant data

	err := query.Find(&models).Error
	if err != nil {
		return nil, err
	}

	skus := make([]*product.SKU, len(models))
	for i, m := range models {
		skus[i] = m.toEntity()
	}
	return skus, nil
}

func (r *skuRepo) FindList(ctx context.Context, db *gorm.DB, query product.SKUQuery) ([]*product.SKU, int64, error) {
	dbQuery := db.WithContext(ctx).Model(&skuModel{})

	if query.ProductID > 0 {
		dbQuery = dbQuery.Where("product_id = ?", query.ProductID)
	}
	if query.Code != "" {
		dbQuery = dbQuery.Where("code LIKE ?", "%"+query.Code+"%")
	}
	if query.Status != 0 {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []skuModel
	offset := (query.Page - 1) * query.PageSize
	if offset < 0 {
		offset = 0
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}

	err := dbQuery.Order("created_at DESC").
		Offset(offset).
		Limit(query.PageSize).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	skus := make([]*product.SKU, len(models))
	for i, m := range models {
		skus[i] = m.toEntity()
	}
	return skus, total, nil
}

// FindLowStock finds SKUs where available_stock < safety_stock AND safety_stock > 0
func (r *skuRepo) FindLowStock(ctx context.Context, db *gorm.DB,  page, pageSize int) ([]*product.SKU, int64, error) {
	dbQuery := db.WithContext(ctx).Model(&skuModel{}).
		Where("safety_stock > 0 AND available_stock < safety_stock").
		Where("status = ?", shared.StatusEnabled)

	// Tenant filter: platform admin (tenantID == 0) can access all tenant data

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var models []skuModel
	err := dbQuery.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	skus := make([]*product.SKU, len(models))
	for i, m := range models {
		skus[i] = m.toEntity()
	}
	return skus, total, nil
}

// Search returns a lightweight SKU item list joined with product name for the dropdown.
// Only enabled SKUs are returned, filtered by LIKE on sku code.
func (r *skuRepo) Search(ctx context.Context, db *gorm.DB, keyword string, page, pageSize int) ([]*product.SKUItem, int64, error) {
	dbQuery := db.WithContext(ctx).Model(&skuModel{}).
		Select("skus.code AS sku_code, skus.product_id, products.name AS product_name, skus.safety_stock").
		Joins("LEFT JOIN products ON products.id = skus.product_id AND products.deleted_at IS NULL").
		Where("skus.status = ?", shared.StatusEnabled)

	// Manual tenant filter for JOINed table (plugin only handles primary table)
	if tenantID, ok := contextx.GetTenantID(ctx); ok && tenantID != 0 {
		dbQuery = dbQuery.Where("products.tenant_id = ?", tenantID)
	}

	if keyword != "" {
		dbQuery = dbQuery.Where("skus.code LIKE ?", "%"+keyword+"%")
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 50 {
		pageSize = 50
	}
	offset := (page - 1) * pageSize

	type row struct {
		SKUCode     string `gorm:"column:sku_code"`
		ProductID   int64  `gorm:"column:product_id"`
		ProductName string `gorm:"column:product_name"`
		SafetyStock int    `gorm:"column:safety_stock"`
	}
	var rows []row
	if err := dbQuery.Order("skus.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	items := make([]*product.SKUItem, len(rows))
	for i, r := range rows {
		items[i] = &product.SKUItem{
			SKUCode:     r.SKUCode,
			ProductID:   r.ProductID,
			ProductName: r.ProductName,
			SafetyStock: r.SafetyStock,
		}
	}
	return items, total, nil
}
