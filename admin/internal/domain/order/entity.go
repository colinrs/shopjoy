package order

import (
	"context"
	"fmt"
	"time"

	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Status int

const (
	StatusPendingPayment Status = iota
	StatusPaid
	StatusPendingShipment
	StatusShipped
	StatusCompleted
	StatusCancelled
	StatusRefunding
	StatusRefunded
)

func (s Status) String() string {
	switch s {
	case StatusPendingPayment:
		return "pending_payment"
	case StatusPaid:
		return "paid"
	case StatusPendingShipment:
		return "pending_shipment"
	case StatusShipped:
		return "shipped"
	case StatusCompleted:
		return "completed"
	case StatusCancelled:
		return "cancelled"
	case StatusRefunding:
		return "refunding"
	case StatusRefunded:
		return "refunded"
	default:
		return "unknown"
	}
}

type Order struct {
	ID             int64
	TenantID       shared.TenantID
	UserID         int64
	OrderNo        string
	Status         Status
	TotalAmount    shared.Money `gorm:"embedded"`
	DiscountAmount shared.Money `gorm:"embedded"`
	FreightAmount  shared.Money `gorm:"embedded"`
	PayAmount      shared.Money `gorm:"embedded"`
	Currency       string
	Items          []OrderItem
	Address        ShippingAddress
	Remark         string
	ExpireAt       time.Time
	PaidAt         *time.Time
	ShippedAt      *time.Time
	CompletedAt    *time.Time
	CancelledAt    *time.Time
	DeletedAt      *int64
	Audit          shared.AuditInfo `gorm:"embedded"`
}

func (o *Order) TableName() string {
	return "orders"
}

func (o *Order) CalculateTotals() {
	total := shared.NewMoney(decimal.Zero, o.Currency)
	for _, item := range o.Items {
		itemTotal := item.Price.Multiply(item.Quantity)
		total, _ = total.Add(itemTotal)
	}
	o.TotalAmount = total

	discount := o.DiscountAmount
	freight := o.FreightAmount

	payAmount := total
	payAmount, _ = payAmount.Subtract(discount)
	payAmount, _ = payAmount.Add(freight)
	o.PayAmount = payAmount
}

func (o *Order) Pay(paymentID string) error {
	if o.Status != StatusPendingPayment {
		return code.ErrOrderInvalidStatus
	}
	if time.Now().After(o.ExpireAt) {
		return code.ErrOrderExpired
	}
	now := time.Now().UTC()
	o.Status = StatusPaid
	o.PaidAt = &now
	return nil
}

func (o *Order) Cancel(reason string) error {
	if o.Status != StatusPendingPayment {
		return code.ErrOrderInvalidStatus
	}
	now := time.Now().UTC()
	o.Status = StatusCancelled
	o.Remark = reason
	o.CancelledAt = &now
	return nil
}

func (o *Order) Ship(trackingNo, carrier string) error {
	if o.Status != StatusPaid && o.Status != StatusPendingShipment {
		return code.ErrOrderInvalidStatus
	}
	now := time.Now().UTC()
	o.Status = StatusShipped
	o.ShippedAt = &now
	o.Address.TrackingNo = trackingNo
	o.Address.Carrier = carrier
	return nil
}

func (o *Order) Complete() error {
	if o.Status != StatusShipped {
		return code.ErrOrderInvalidStatus
	}
	now := time.Now().UTC()
	o.Status = StatusCompleted
	o.CompletedAt = &now
	return nil
}

type OrderItem struct {
	ID          int64
	OrderID     int64 `gorm:"column:order_id;not null;index:idx_order_id"`
	ProductID   int64
	SKUId       int64
	ProductName string
	SKUName     string
	Image       string
	Price       shared.Money `gorm:"embedded"`
	Quantity    int
	TotalAmount shared.Money `gorm:"embedded"`
	DeletedAt   *int64
}

func (oi *OrderItem) TableName() string {
	return "order_items"
}

type ShippingAddress struct {
	Name       string
	Phone      string
	Province   string
	City       string
	District   string
	Address    string
	ZipCode    string
	TrackingNo string
	Carrier    string
}

type Repository interface {
	Create(ctx context.Context, db *gorm.DB, order *Order) error
	Update(ctx context.Context, db *gorm.DB, order *Order) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Order, error)
	FindByOrderNo(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderNo string) (*Order, error)
	FindByUserID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64, query Query) ([]*Order, int64, error)
	FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query Query) ([]*Order, int64, error)
	UpdateStatus(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64, status Status) error
}

type Query struct {
	shared.PageQuery
	UserID    int64
	Status    Status
	StartTime *time.Time
	EndTime   *time.Time
}

func GenerateOrderNo(tenantID shared.TenantID) string {
	return fmt.Sprintf("%s%d%d", time.Now().Format("20060102"), tenantID.Int64(), time.Now().UnixNano()%1000000)
}
