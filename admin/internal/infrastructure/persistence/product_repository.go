package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type productRepo struct{}

func NewProductRepository() product.Repository {
	return &productRepo{}
}

type productModel struct {
	ID              int64  `gorm:"column:id;primaryKey;autoIncrement:false"`
	SKU             string `gorm:"column:sku;size:64;uniqueIndex"`
	Name            string `gorm:"column:name;size:200;not null;index"`
	Description     string `gorm:"column:description;type:text"`
	Price           int64  `gorm:"column:price;not null"`
	CostPrice       int64  `gorm:"column:cost_price"`
	Currency        string `gorm:"column:currency;size:10;default:'CNY'"`
	Stock           int    `gorm:"column:stock;default:0"`
	Status          int    `gorm:"column:status;default:0;index"`
	CategoryID      int64  `gorm:"column:category_id;index"`
	Brand           string `gorm:"column:brand;size:64"`
	Tags            string `gorm:"column:tags;type:json"`
	Images          string `gorm:"column:images;type:json"`
	IsMatrixProduct bool   `gorm:"column:is_matrix_product;default:false"`
	HSCode          string `gorm:"column:hs_code;size:20"`
	COO             string `gorm:"column:coo;size:10"`
	Weight          string `gorm:"column:weight;type:decimal(10,2)"`
	WeightUnit      string `gorm:"column:weight_unit;size:10;default:'g'"`
	Length          string `gorm:"column:length;type:decimal(10,2)"`
	Width           string `gorm:"column:width;type:decimal(10,2)"`
	Height          string `gorm:"column:height;type:decimal(10,2)"`
	DangerousGoods  string `gorm:"column:dangerous_goods;type:json"`
	CreatedAt       int64  `gorm:"column:created_at"`
	UpdatedAt       int64  `gorm:"column:updated_at"`
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

	return &product.Product{
		ID:              m.ID,
		SKU:             m.SKU,
		Name:            m.Name,
		Description:     m.Description,
		Price:           product.NewMoney(m.Price, m.Currency),
		CostPrice:       product.NewMoney(m.CostPrice, m.Currency),
		Stock:           m.Stock,
		Status:          product.Status(m.Status),
		CategoryID:      m.CategoryID,
		Brand:           m.Brand,
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
		CreatedAt:      time.Unix(m.CreatedAt, 0),
		UpdatedAt:      time.Unix(m.UpdatedAt, 0),
	}
}

func fromEntity(p *product.Product) *productModel {
	// Serialize JSON arrays
	tagsJSON, _ := json.Marshal(p.Tags)
	imagesJSON, _ := json.Marshal(p.Images)
	dangerousGoodsJSON, _ := json.Marshal(p.DangerousGoods)

	return &productModel{
		ID:              p.ID,
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
		CreatedAt:       p.CreatedAt.Unix(),
		UpdatedAt:       p.UpdatedAt.Unix(),
	}
}

func (r *productRepo) Create(ctx context.Context, db *gorm.DB, p *product.Product) error {
	model := fromEntity(p)
	return db.WithContext(ctx).Create(model).Error
}

func (r *productRepo) Update(ctx context.Context, db *gorm.DB, p *product.Product) error {
	model := fromEntity(p)
	return db.WithContext(ctx).
		Where("id = ? AND status != ?", p.ID, product.StatusDeleted).
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

func (r *productRepo) Delete(ctx context.Context, db *gorm.DB, id int64) error {
	result := db.WithContext(ctx).Model(&productModel{}).
		Where("id = ?", id).
		Update("status", product.StatusDeleted)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return code.ErrProductNotFound
	}
	return nil
}

func (r *productRepo) FindByID(ctx context.Context, db *gorm.DB, id int64) (*product.Product, error) {
	var model productModel
	err := db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrProductNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *productRepo) FindByIDs(ctx context.Context, db *gorm.DB, ids []int64) ([]*product.Product, error) {
	if len(ids) == 0 {
		return []*product.Product{}, nil
	}

	var models []productModel
	err := db.WithContext(ctx).
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

	dbQuery := db.WithContext(ctx).Model(&productModel{})

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

func (r *productRepo) UpdateStock(ctx context.Context, db *gorm.DB, id int64, delta int) error {
	result := db.WithContext(ctx).Model(&productModel{}).
		Where("id = ? AND status = ? AND stock + ? >= 0", id, product.StatusOnSale, delta).
		UpdateColumn("stock", gorm.Expr("stock + ?", delta))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return code.ErrProductInsufficientStock
	}
	return nil
}

func (r *productRepo) Exists(ctx context.Context, db *gorm.DB, id int64) (bool, error) {
	var count int64
	err := db.WithContext(ctx).Model(&productModel{}).
		Where("id = ? AND status != ?", id, product.StatusDeleted).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
