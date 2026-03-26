package cart

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type Cart struct {
	ID        int64
	TenantID  shared.TenantID
	UserID    int64
	SessionID string
	Items     []CartItem
	UpdatedAt time.Time
	DeletedAt *int64
}

func (c *Cart) TableName() string {
	return "carts"
}

func (c *Cart) AddItem(productID, skuID int64, quantity int, price shared.Money, productName, skuName, image string) error {
	if quantity <= 0 {
		return code.ErrCartInvalidQuantity
	}

	for i := range c.Items {
		if c.Items[i].SKUId == skuID {
			c.Items[i].Quantity += quantity
			c.Items[i].UpdateTotal()
			c.UpdatedAt = time.Now().UTC()
			return nil
		}
	}

	item := CartItem{
		TenantID:    c.TenantID,
		UserID:      c.UserID,
		ProductID:   productID,
		SKUId:       skuID,
		ProductName: productName,
		SKUName:     skuName,
		Image:       image,
		Price:       price,
		Quantity:    quantity,
	}
	item.UpdateTotal()
	c.Items = append(c.Items, item)
	c.UpdatedAt = time.Now().UTC()
	return nil
}

func (c *Cart) UpdateItem(skuID int64, quantity int) error {
	if quantity <= 0 {
		return c.RemoveItem(skuID)
	}

	for i := range c.Items {
		if c.Items[i].SKUId == skuID {
			c.Items[i].Quantity = quantity
			c.Items[i].UpdateTotal()
			c.UpdatedAt = time.Now().UTC()
			return nil
		}
	}
	return code.ErrCartItemNotFound
}

func (c *Cart) RemoveItem(skuID int64) error {
	for i := range c.Items {
		if c.Items[i].SKUId == skuID {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			c.UpdatedAt = time.Now().UTC()
			return nil
		}
	}
	return code.ErrCartItemNotFound
}

func (c *Cart) Clear() {
	c.Items = nil
	c.UpdatedAt = time.Now().UTC()
}

func (c *Cart) GetTotal() shared.Money {
	total := shared.NewMoney(0, "CNY")
	for _, item := range c.Items {
		total, _ = total.Add(item.TotalAmount)
	}
	return total
}

func (c *Cart) IsEmpty() bool {
	return len(c.Items) == 0
}

type CartItem struct {
	ID          int64
	TenantID    shared.TenantID
	UserID      int64
	CartID      int64
	ProductID   int64
	SKUId       int64
	ProductName string
	SKUName     string
	Image       string
	Price       shared.Money `gorm:"embedded"`
	Quantity    int
	TotalAmount shared.Money `gorm:"embedded"`
	Selected    bool
	DeletedAt   *int64
	Audit       shared.AuditInfo `gorm:"embedded"`
}

func (ci *CartItem) TableName() string {
	return "cart_items"
}

func (ci *CartItem) UpdateTotal() {
	ci.TotalAmount = ci.Price.Multiply(ci.Quantity)
}

type Repository interface {
	FindByUserID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64) (*Cart, error)
	FindBySessionID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, sessionID string) (*Cart, error)
	Save(ctx context.Context, db *gorm.DB, cart *Cart) error
	SaveItem(ctx context.Context, db *gorm.DB, item *CartItem) error
	DeleteItem(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, itemID int64) error
	Clear(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64) error
	Merge(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, fromSessionID string, toUserID int64) error
}
