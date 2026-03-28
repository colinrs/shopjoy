package order

import (
	"context"
	"fmt"
	"time"

	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

// ==================== Status 订单状态 ====================

// Status 订单状态
type Status int

const (
	StatusPendingPayment Status = iota // 待支付
	StatusPaid                       // 已支付
	StatusPendingShipment            // 待发货
	StatusShipped                    // 已发货
	StatusCompleted                  // 已完成
	StatusCancelled                  // 已取消
	StatusRefunding                 // 退款中
	StatusRefunded                  // 已退款
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

func (s Status) IsValid() bool {
	return s >= StatusPendingPayment && s <= StatusRefunded
}

// CanTransitionTo 检查状态是否可以转换
func (s Status) CanTransitionTo(target Status) bool {
	transitions := map[Status][]Status{
		StatusPendingPayment: {StatusPaid, StatusCancelled},
		StatusPaid:         {StatusPendingShipment, StatusCancelled, StatusRefunding},
		StatusPendingShipment: {StatusShipped, StatusCancelled},
		StatusShipped:       {StatusCompleted, StatusRefunding},
		StatusCompleted:     {},
		StatusCancelled:     {},
		StatusRefunding:    {StatusRefunded, StatusPendingShipment},
		StatusRefunded:     {},
	}

	allowed, ok := transitions[s]
	if !ok {
		return false
	}

	for _, status := range allowed {
		if status == target {
			return true
		}
	}
	return false
}

// ==================== FulfillmentStatus 履约状态 ====================

// FulfillmentStatus 订单履约状态
type FulfillmentStatus int

const (
	FulfillmentStatusPending FulfillmentStatus = iota // 待发货
	FulfillmentStatusPartialShipped                   // 部分发货
	FulfillmentStatusShipped                          // 已发货
	FulfillmentStatusDelivered                        // 已送达
)

func (s FulfillmentStatus) String() string {
	switch s {
	case FulfillmentStatusPending:
		return "pending"
	case FulfillmentStatusPartialShipped:
		return "partial_shipped"
	case FulfillmentStatusShipped:
		return "shipped"
	case FulfillmentStatusDelivered:
		return "delivered"
	default:
		return "unknown"
	}
}

// ==================== RefundStatus 退款状态 ====================

// RefundStatus 退款状态
type RefundStatus int

const (
	RefundStatusNone RefundStatus = iota // 无退款
	RefundStatusPending                  // 待处理
	RefundStatusApproved                 // 已批准
	RefundStatusRejected                 // 已拒绝
	RefundStatusCompleted                // 已完成
)

func (s RefundStatus) String() string {
	switch s {
	case RefundStatusNone:
		return "none"
	case RefundStatusPending:
		return "pending"
	case RefundStatusApproved:
		return "approved"
	case RefundStatusRejected:
		return "rejected"
	case RefundStatusCompleted:
		return "completed"
	default:
		return "unknown"
	}
}

// ==================== Order 订单实体 ====================

// Order 订单实体
type Order struct {
	application.Model
	ID                  string          `gorm:"column:id;primaryKey;size:64"`
	TenantID            shared.TenantID `gorm:"column:tenant_id;not null;index:idx_tenant_id"`
	UserID              int64           `gorm:"column:user_id;not null;index:idx_user_id"`
	OrderNo             string          `gorm:"column:order_no;not null;uniqueIndex:uk_order_no;size:64"`
	Status              Status          `gorm:"column:status;not null;default:0;index:idx_status"`
	FulfillmentStatus   FulfillmentStatus `gorm:"column:fulfillment_status;not null;default:0;index:idx_fulfillment_status"`
	RefundStatus        RefundStatus    `gorm:"column:refund_status;not null;default:0;index:idx_refund_status"`
	TotalAmount         shared.Money    `gorm:"column:total_amount;type:decimal(19,4);not null;embedded"`
	DiscountAmount      shared.Money    `gorm:"column:discount_amount;type:decimal(19,4);not null;embedded"`
	FreightAmount       shared.Money    `gorm:"column:freight_amount;type:decimal(19,4);not null;embedded"`
	PayAmount           shared.Money    `gorm:"column:pay_amount;type:decimal(19,4);not null;embedded"`
	OriginalAmount      shared.Money    `gorm:"column:original_amount;type:decimal(19,4);not null;embedded"`
	AdjustAmount        shared.Money    `gorm:"column:adjust_amount;type:decimal(19,4);not null;embedded"`
	AdjustReason        string          `gorm:"column:adjust_reason;size:200;not null;default:''"`
	AdjustedBy           int64           `gorm:"column:adjusted_by;not null;default:0"`
	AdjustedAt          *time.Time      `gorm:"column:adjusted_at"`
	Version             int             `gorm:"column:version;not null;default:1"`
	PaymentMethod       string          `gorm:"column:payment_method;size:32;not null;default:''"`
	Source              string          `gorm:"column:source;size:32;not null;default:''"`
	Currency            string          `gorm:"column:currency;size:10;not null;default:'CNY'"`
	Address             ShippingAddress `gorm:"embedded"`
	TrackingNo          string          `gorm:"column:tracking_no;size:100;not null;default:''"`
	Carrier             string          `gorm:"column:carrier;size:50;not null;default:''"`
	Remark              string          `gorm:"column:remark;type:text"`
	MerchantRemark     string          `gorm:"column:merchant_remark;size:500;not null;default:''"`
	ExpireAt           time.Time       `gorm:"column:expire_at;not null"`
	PaidAt             *time.Time      `gorm:"column:paid_at"`
	ShippedAt          *time.Time      `gorm:"column:shipped_at"`
	CompletedAt        *time.Time      `gorm:"column:completed_at"`
	CancelledAt        *time.Time      `gorm:"column:cancelled_at"`
	Audit              shared.AuditInfo `gorm:"embedded"`
	Items              []OrderItem     `gorm:"foreignKey:OrderID"`
}

func (o *Order) TableName() string {
	return "orders"
}

// CalculateTotals 计算订单金额
func (o *Order) CalculateTotals() {
	total := shared.NewMoney(0, o.Currency)
	for _, item := range o.Items {
		itemTotal := item.Price.Multiply(item.Quantity)
		total, _ = total.Add(itemTotal)
	}
	o.TotalAmount = total
	o.OriginalAmount = total

	discount := o.DiscountAmount
	freight := o.FreightAmount

	payAmount := total
	payAmount, _ = payAmount.Subtract(discount)
	payAmount, _ = payAmount.Add(freight)
	o.PayAmount = payAmount
}

// Pay 支付成功
func (o *Order) Pay() error {
	if !o.Status.CanTransitionTo(StatusPaid) {
		return code.ErrOrderInvalidStatus
	}
	if time.Now().UTC().After(o.ExpireAt) {
		return code.ErrOrderExpired
	}
	now := time.Now().UTC()
	o.Status = StatusPaid
	o.PaidAt = &now
	o.Version++
	return nil
}

// Cancel 取消订单
func (o *Order) Cancel(reason string) error {
	if !o.Status.CanTransitionTo(StatusCancelled) {
		return code.ErrOrderInvalidStatus
	}
	now := time.Now().UTC()
	o.Status = StatusCancelled
	o.Remark = reason
	o.CancelledAt = &now
	o.Version++
	return nil
}

// Ship 发货
func (o *Order) Ship(trackingNo, carrier string) error {
	if !o.Status.CanTransitionTo(StatusShipped) {
		return code.ErrOrderInvalidStatus
	}
	now := time.Now().UTC()
	o.Status = StatusShipped
	o.ShippedAt = &now
	o.TrackingNo = trackingNo
	o.Carrier = carrier
	o.Version++
	return nil
}

// Complete 完成订单
func (o *Order) Complete() error {
	if !o.Status.CanTransitionTo(StatusCompleted) {
		return code.ErrOrderInvalidStatus
	}
	now := time.Now().UTC()
	o.Status = StatusCompleted
	o.CompletedAt = &now
	o.Version++
	return nil
}

// AdjustPrice 调整价格
func (o *Order) AdjustPrice(amount int64, reason string, adjustedBy int64) error {
	if o.Status != StatusPendingPayment {
		return code.ErrOrderInvalidStatus
	}
	now := time.Now().UTC()
	o.AdjustAmount = shared.NewMoney(amount, o.Currency)
	o.AdjustReason = reason
	o.AdjustedBy = adjustedBy
	o.AdjustedAt = &now
	o.Version++
	return nil
}

// IsExpired 检查是否过期
func (o *Order) IsExpired() bool {
	return time.Now().UTC().After(o.ExpireAt)
}

// ==================== ShippingAddress 收货地址 ====================

// ShippingAddress 收货地址
type ShippingAddress struct {
	Name        string `gorm:"column:address_name;size:100;not null;default:''"`
	Phone       string `gorm:"column:address_phone;size:20;not null;default:''"`
	Province    string `gorm:"column:address_province;size:50;not null;default:''"`
	City        string `gorm:"column:address_city;size:50;not null;default:''"`
	District    string `gorm:"column:address_district;size:50;not null;default:''"`
	Detail      string `gorm:"column:address_detail;type:text"`
	ZipCode     string `gorm:"column:address_zipcode;size:20;not null;default:''"`
	TrackingNo  string `gorm:"column:tracking_no;size:100;not null;default:''"`
	Carrier     string `gorm:"column:carrier;size:50;not null;default:''"`
}

// ==================== OrderItem 订单商品 ====================

// OrderItem 订单商品
type OrderItem struct {
	application.Model
	OrderID     string       `gorm:"column:order_id;not null;size:64;index:idx_order_id"`
	ProductID   int64        `gorm:"column:product_id;not null;index:idx_product_id"`
	SKUId       int64        `gorm:"column:sku_id;not null;index:idx_sku_id"`
	ProductName string       `gorm:"column:product_name;size:255;not null"`
	SKUName     string       `gorm:"column:sku_name;size:255;not null;default:''"`
	Image       string       `gorm:"column:image;size:500;not null;default:''"`
	Price       shared.Money `gorm:"column:price;type:decimal(19,4);not null;embedded"`
	Quantity    int          `gorm:"column:quantity;not null;default:1"`
	TotalAmount shared.Money `gorm:"column:total_amount;type:decimal(19,4);not null;embedded"`
}

func (oi *OrderItem) TableName() string {
	return "order_items"
}

// ==================== Repository 接口 ====================

// Repository 订单仓储接口
type Repository interface {
	Create(ctx context.Context, db *gorm.DB, order *Order) error
	Update(ctx context.Context, db *gorm.DB, order *Order) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id string) (*Order, error)
	FindByOrderNo(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderNo string) (*Order, error)
	FindByUserID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64, query Query) ([]*Order, int64, error)
	FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query Query) ([]*Order, int64, error)
	UpdateStatus(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id string, status Status) error
}

// Query 查询条件
type Query struct {
	shared.PageQuery
	UserID    int64
	Status    *Status
	StartTime *time.Time
	EndTime   *time.Time
}

// ==================== 辅助函数 ====================

// GenerateOrderNo 生成订单号
func GenerateOrderNo(tenantID shared.TenantID) string {
	return fmt.Sprintf("ORD%s%d%d", time.Now().Format("20060102"), tenantID.Int64(), time.Now().UnixNano()%1000000)
}
