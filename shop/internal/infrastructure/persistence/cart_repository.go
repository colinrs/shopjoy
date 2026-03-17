package persistence

import (
	"context"
	"errors"

	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/shop/internal/domain/cart"
	"gorm.io/gorm"
)

type CartRepository struct{}

func NewCartRepository() cart.Repository {
	return &CartRepository{}
}

func (r *CartRepository) FindByUserID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64) (*cart.Cart, error) {
	var c cart.Cart
	err := db.WithContext(ctx).Where("tenant_id = ? AND user_id = ?", tenantID.Int64(), userID).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c = cart.Cart{
			TenantID: tenantID,
			UserID:   userID,
		}
		return &c, nil
	}

	if err == nil {
		var items []cart.CartItem
		db.WithContext(ctx).Where("cart_id = ?", c.ID).Find(&items)
		c.Items = items
	}

	return &c, err
}

func (r *CartRepository) FindBySessionID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, sessionID string) (*cart.Cart, error) {
	var c cart.Cart
	err := db.WithContext(ctx).Where("tenant_id = ? AND session_id = ?", tenantID.Int64(), sessionID).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c = cart.Cart{
			TenantID:  tenantID,
			SessionID: sessionID,
		}
		return &c, nil
	}

	if err == nil {
		var items []cart.CartItem
		db.WithContext(ctx).Where("cart_id = ?", c.ID).Find(&items)
		c.Items = items
	}

	return &c, err
}

func (r *CartRepository) Save(ctx context.Context, db *gorm.DB, c *cart.Cart) error {
	return db.WithContext(ctx).Save(c).Error
}

func (r *CartRepository) SaveItem(ctx context.Context, db *gorm.DB, item *cart.CartItem) error {
	return db.WithContext(ctx).Save(item).Error
}

func (r *CartRepository) DeleteItem(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, itemID int64) error {
	return db.WithContext(ctx).Where("id = ? AND tenant_id = ?", itemID, tenantID.Int64()).Delete(&cart.CartItem{}).Error
}

func (r *CartRepository) Clear(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var c cart.Cart
		if err := tx.Where("tenant_id = ? AND user_id = ?", tenantID.Int64(), userID).First(&c).Error; err != nil {
			return err
		}
		if err := tx.Where("cart_id = ?", c.ID).Delete(&cart.CartItem{}).Error; err != nil {
			return err
		}
		return tx.Delete(&c).Error
	})
}

func (r *CartRepository) Merge(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, fromSessionID string, toUserID int64) error {
	sessionCart, err := r.FindBySessionID(ctx, db, tenantID, fromSessionID)
	if err != nil {
		return err
	}

	if len(sessionCart.Items) == 0 {
		return nil
	}

	userCart, err := r.FindByUserID(ctx, db, tenantID, toUserID)
	if err != nil {
		return err
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if userCart.ID == 0 {
			userCart.TenantID = tenantID
			userCart.UserID = toUserID
			if err := tx.Create(userCart).Error; err != nil {
				return err
			}
		}

		for _, item := range sessionCart.Items {
			item.ID = 0
			item.CartID = userCart.ID
			item.UserID = toUserID
			if err := tx.Create(&item).Error; err != nil {
				return err
			}
		}

		return tx.Where("session_id = ?", fromSessionID).Delete(&cart.Cart{}).Error
	})
}
