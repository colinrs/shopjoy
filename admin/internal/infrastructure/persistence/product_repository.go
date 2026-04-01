package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type productRepo struct{}

func NewProductRepository() product.Repository {
	return &productRepo{}
}

type productModel struct {
	application.Model
	TenantID        int64           `gorm:"column:tenant_id;not null;index"`
	SKU             string          `gorm:"column:sku;size:64;uniqueIndex"`
	Name            string          `gorm:"column:name;size:200;not null;index"`
	Description     string          `gorm:"column:description;type:text"`
	Price           decimal.Decimal `gorm:"column:price;type:decimal(19,4);not null;default:0"`
	CostPrice       decimal.Decimal `gorm:"column:cost_price;type:decimal(19,4);not null;default:0"`
	Currency        string          `gorm:"column:currency;size:10;default:'CNY'"`
	Stock           int             `gorm:"column:stock;default:0"`
	Status          int             `gorm:"column:status;default:0;index"`
	CategoryID      int64           `gorm:"column:category_id;index"`
	Brand           string          `gorm:"column:brand;size:64"`
	SKUPrefix       string          `gorm:"column:sku_prefix;size:8;default:''"`
	Tags            string          `gorm:"column:tags;type:json"`
	Images          string          `gorm:"column:images;type:json"`
	IsMatrixProduct bool            `gorm:"column:is_matrix_product;default:false"`
	HSCode          string          `gorm:"column:hs_code;size:20"`
	COO             string          `gorm:"column:coo;size:10"`
	Weight          string          `gorm:"column:weight;type:decimal(10,2)"`
	WeightUnit      string          `gorm:"column:weight_unit;size:10;default:'g'"`
	Length          string          `gorm:"column:length;type:decimal(10,2)"`
	Width           string          `gorm:"column:width;type:decimal(10,2)"`
	Height          string          `gorm:"column:height;type:decimal(10,2)"`
	DangerousGoods  string          `gorm:"column:dangerous_goods;type:json"`
}

func (productModel) TableName() string {
	return "products"
}

func (m *productModel) toEntity() *product.Product {
	// Parse JSON arrays
	var tags, images, dangerousGoods []string
	if m.Tags != "" {
		json.Unmarshal([]byte(m.Tags), &tags)
	}
	if m.Images != "" {
		json.Unmarshal([]byte(m.Images), &images)
	}
	if m.DangerousGoods != "" {
		json.Unmarshal([]byte(m.DangerousGoods), &dangerousGoods)
	}

	// Parse decimals
	weight, _ := decimal.NewFromString(m.Weight)
	length, _ := decimal.NewFromString(m.Length)
	width, _ := decimal.NewFromString(m.Width)
	height, _ := decimal.NewFromString(m.Height)

	entity := &product.Product{
		TenantID:        shared.TenantID(m.TenantID),
		SKU:             m.SKU,
		Name:            m.Name,
		Description:     m.Description,
		Price:           product.NewMoney(m.Price, m.Currency),
		CostPrice:       product.NewMoney(m.CostPrice, m.Currency),
		Stock:           m.Stock,
		Status:          product.Status(m.Status),
		CategoryID:      m.CategoryID,
		Brand:           m.Brand,
		SKUPrefix:       m.SKUPrefix,
		Tags:            tags,
		Images:          images,
		IsMatrixProduct: m.IsMatrixProduct,
		HSCode:          m.HSCode,
		COO:             m.COO,
		Weight:          weight,
		WeightUnit:      m.WeightUnit,
		Dimensions: product.Dimensions{
			Length: length,
			Width:  width,
			Height: height,
			Unit:   "cm",
		},
		DangerousGoods: dangerousGoods,
	}
	entity.ID = m.ID
	entity.CreatedAt = m.CreatedAt
	entity.UpdatedAt = m.UpdatedAt
	entity.DeletedAt = m.DeletedAt
	return entity
}

func fromEntity(p *product.Product) *productModel {
	// Serialize JSON arrays
	tagsJSON, _ := json.Marshal(p.Tags)
	imagesJSON, _ := json.Marshal(p.Images)
	dangerousGoodsJSON, _ := json.Marshal(p.DangerousGoods)

	model := &productModel{
		TenantID:        p.TenantID.Int64(),
		SKU:             p.SKU,
		Name:            p.Name,
		Description:     p.Description,
		Price:           p.Price.Amount,
		CostPrice:       p.CostPrice.Amount,
		Currency:        p.Price.Currency,
		Stock:           p.Stock,
		Status:          int(p.Status),
		CategoryID:      p.CategoryID,
		Brand:           p.Brand,
		SKUPrefix:       p.SKUPrefix,
		Tags:            string(tagsJSON),
		Images:          string(imagesJSON),
		IsMatrixProduct: p.IsMatrixProduct,
		HSCode:          p.HSCode,
		COO:             p.COO,
		Weight:          p.Weight.String(),
		WeightUnit:      p.WeightUnit,
		Length:          p.Dimensions.Length.String(),
		Width:           p.Dimensions.Width.String(),
		Height:          p.Dimensions.Height.String(),
		DangerousGoods:  string(dangerousGoodsJSON),
	}
	model.ID = p.ID
	model.CreatedAt = p.CreatedAt
	model.UpdatedAt = p.UpdatedAt
	return model
}

func (r *productRepo) Create(ctx context.Context, db *gorm.DB, p *product.Product) error {
	model := fromEntity(p)
	return db.WithContext(ctx).Create(model).Error
}

func (r *productRepo) Update(ctx context.Context, db *gorm.DB, p *product.Product) error {
	model := fromEntity(p)
	return db.WithContext(ctx).
		Model(&productModel{}).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", p.ID, p.TenantID.Int64()).
		Updates(map[string]interface{}{
			"sku":               model.SKU,
			"name":              model.Name,
			"description":       model.Description,
			"price":             model.Price,
			"cost_price":        model.CostPrice,
			"currency":          model.Currency,
			"stock":             model.Stock,
			"status":            model.Status,
			"category_id":       model.CategoryID,
			"brand":             model.Brand,
			"tags":              model.Tags,
			"images":            model.Images,
			"is_matrix_product": model.IsMatrixProduct,
			"hs_code":           model.HSCode,
			"coo":               model.COO,
			"weight":            model.Weight,
			"weight_unit":       model.WeightUnit,
			"length":            model.Length,
			"width":             model.Width,
			"height":            model.Height,
			"dangerous_goods":   model.DangerousGoods,
			"updated_at":        model.UpdatedAt,
		}).Error
}

func (r *productRepo) Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error {
	query := db.WithContext(ctx).Model(&productModel{}).Where("id = ? AND deleted_at IS NULL", id)
	// 平台管理员 (tenantID == 0) 可删除所有租户数据
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	now := time.Now().UTC()
	result := query.Update("deleted_at", now)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return code.ErrProductNotFound
	}
	return nil
}

func (r *productRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*product.Product, error) {
	query := db.WithContext(ctx).Where("deleted_at IS NULL")
	// 平台管理员 (tenantID == 0) 可访问所有租户数据
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var model productModel
	err := query.First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrProductNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *productRepo) FindByIDs(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, ids []int64) ([]*product.Product, error) {
	if len(ids) == 0 {
		return []*product.Product{}, nil
	}

	query := db.WithContext(ctx).Where("deleted_at IS NULL")
	// 平台管理员 (tenantID == 0) 可访问所有租户数据
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}

	var models []productModel
	err := query.
		Where("id IN ?", ids).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	products := make([]*product.Product, len(models))
	for i, m := range models {
		products[i] = m.toEntity()
	}
	return products, nil
}

func (r *productRepo) FindList(ctx context.Context, db *gorm.DB, query product.Query) ([]*product.Product, int64, error) {
	query.Validate()

	dbQuery := db.WithContext(ctx).Model(&productModel{}).Where("deleted_at IS NULL")

	// 租户过滤：平台管理员 (TenantID == 0) 可访问所有租户数据
	if query.TenantID != 0 {
		dbQuery = dbQuery.Where("tenant_id = ?", query.TenantID.Int64())
	}

	if query.Name != "" {
		dbQuery = dbQuery.Where("name LIKE ?", fmt.Sprintf("%%%s%%", query.Name))
	}
	if query.CategoryID > 0 {
		dbQuery = dbQuery.Where("category_id = ?", query.CategoryID)
	}
	if query.Status != nil && query.Status.IsValid() {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}
	if query.MinPrice != nil {
		dbQuery = dbQuery.Where("price >= ?", *query.MinPrice)
	}
	if query.MaxPrice != nil {
		dbQuery = dbQuery.Where("price <= ?", *query.MaxPrice)
	}

	// Filter by market - use subquery to avoid duplicate rows from JOIN
	if query.MarketID > 0 {
		dbQuery = dbQuery.Where("id IN (SELECT DISTINCT product_id FROM product_markets WHERE market_id = ?)", query.MarketID)
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []productModel
	err := dbQuery.Order("products.created_at DESC").
		Offset(query.Offset()).
		Limit(query.Limit()).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	products := make([]*product.Product, len(models))
	for i, m := range models {
		products[i] = m.toEntity()
	}
	return products, total, nil
}

func (r *productRepo) UpdateStock(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64, delta int) error {
	query := db.WithContext(ctx).Model(&productModel{}).
		Where("id = ? AND status = ? AND deleted_at IS NULL", id, product.StatusOnSale)
	// 租户过滤：平台管理员 (tenantID == 0) 可操作所有租户数据
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	result := query.Where("stock + ? >= 0", delta).
		UpdateColumn("stock", gorm.Expr("stock + ?", delta))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return code.ErrProductInsufficientStock
	}
	return nil
}

func (r *productRepo) Exists(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (bool, error) {
	query := db.WithContext(ctx).Model(&productModel{}).
		Where("id = ? AND deleted_at IS NULL", id)
	// 租户过滤：平台管理员 (tenantID == 0) 可访问所有租户数据
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
