package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"gorm.io/gorm"
)

type productRepo struct{}

func NewProductRepository() product.Repository {
	return &productRepo{}
}

type productModel struct {
	ID          int64  `gorm:"column:id;primaryKey;autoIncrement:false"`
	Name        string `gorm:"column:name;size:200;not null;index"`
	Description string `gorm:"column:description;type:text"`
	Price       int64  `gorm:"column:price;not null"`
	CostPrice   int64  `gorm:"column:cost_price"`
	Currency    string `gorm:"column:currency;size:10;default:'CNY'"`
	Stock       int    `gorm:"column:stock;default:0"`
	Status      int    `gorm:"column:status;default:0;index"`
	CategoryID  int64  `gorm:"column:category_id;index"`
	CreatedAt   int64  `gorm:"column:created_at"`
	UpdatedAt   int64  `gorm:"column:updated_at"`
}

func (productModel) TableName() string {
	return "products"
}

func (m *productModel) toEntity() *product.Product {
	return &product.Product{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Price:       product.NewMoney(m.Price, m.Currency),
		CostPrice:   product.NewMoney(m.CostPrice, m.Currency),
		Stock:       m.Stock,
		Status:      product.Status(m.Status),
		CategoryID:  m.CategoryID,
	}
}

func fromEntity(p *product.Product) *productModel {
	return &productModel{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price.Amount,
		CostPrice:   p.CostPrice.Amount,
		Currency:    p.Price.Currency,
		Stock:       p.Stock,
		Status:      int(p.Status),
		CategoryID:  p.CategoryID,
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
			"name":        model.Name,
			"description": model.Description,
			"price":       model.Price,
			"cost_price":  model.CostPrice,
			"currency":    model.Currency,
			"stock":       model.Stock,
			"status":      model.Status,
			"category_id": model.CategoryID,
			"updated_at":  model.UpdatedAt,
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
		return product.ErrProductNotFound
	}
	return nil
}

func (r *productRepo) FindByID(ctx context.Context, db *gorm.DB, id int64) (*product.Product, error) {
	var model productModel
	err := db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, product.ErrProductNotFound
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
	if query.Status.IsValid() {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}
	if query.MinPrice != nil {
		dbQuery = dbQuery.Where("price >= ?", *query.MinPrice)
	}
	if query.MaxPrice != nil {
		dbQuery = dbQuery.Where("price <= ?", *query.MaxPrice)
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []productModel
	err := dbQuery.Order("created_at DESC").
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
		return product.ErrInsufficientStock
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
