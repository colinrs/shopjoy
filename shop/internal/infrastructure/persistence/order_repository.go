package persistence

import (
	"context"
	"errors"

	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/shop/internal/domain/order"
	"gorm.io/gorm"
)

type OrderRepository struct{}

func NewOrderRepository() order.Repository {
	return &OrderRepository{}
}

func (r *OrderRepository) Create(ctx context.Context, db *gorm.DB, o *order.Order) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(o).Error; err != nil {
			return err
		}
		for i := range o.Items {
			o.Items[i].OrderID = o.ID
			if err := tx.Create(&o.Items[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *OrderRepository) Update(ctx context.Context, db *gorm.DB, o *order.Order) error {
	return db.WithContext(ctx).Save(o).Error
}

func (r *OrderRepository) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id string) (*order.Order, error) {
	var o order.Order
	err := db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID.Int64()).First(&o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, order.ErrOrderNotFound
	}
	if err == nil {
		var items []order.OrderItem
		db.WithContext(ctx).Where("order_id = ?", o.ID).Find(&items)
		o.Items = items
	}
	return &o, err
}

func (r *OrderRepository) FindByOrderNo(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderNo string) (*order.Order, error) {
	var o order.Order
	err := db.WithContext(ctx).Where("order_no = ? AND tenant_id = ?", orderNo, tenantID.Int64()).First(&o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, order.ErrOrderNotFound
	}
	if err == nil {
		var items []order.OrderItem
		db.WithContext(ctx).Where("order_id = ?", o.ID).Find(&items)
		o.Items = items
	}
	return &o, err
}

func (r *OrderRepository) FindByUserID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64, query order.Query) ([]*order.Order, int64, error) {
	var orders []*order.Order
	var total int64

	dbQuery := db.WithContext(ctx).Model(&order.Order{}).
		Where("tenant_id = ? AND user_id = ?", tenantID.Int64(), userID)

	if query.Status != 0 {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := dbQuery.Order("created_at DESC").Offset(query.Offset()).Limit(query.Limit()).Find(&orders).Error
	return orders, total, err
}

func (r *OrderRepository) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query order.Query) ([]*order.Order, int64, error) {
	var orders []*order.Order
	var total int64

	dbQuery := db.WithContext(ctx).Model(&order.Order{}).Where("tenant_id = ?", tenantID.Int64())

	if query.Status != 0 {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := dbQuery.Order("created_at DESC").Offset(query.Offset()).Limit(query.Limit()).Find(&orders).Error
	return orders, total, err
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id string, status order.Status) error {
	return db.WithContext(ctx).Model(&order.Order{}).
		Where("id = ? AND tenant_id = ?", id, tenantID.Int64()).
		Update("status", status).Error
}
