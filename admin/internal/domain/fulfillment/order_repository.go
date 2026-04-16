package fulfillment

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ==================== Enums ====================

// OrderStatus 订单状态
type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "pending_payment" // 待付款
	OrderStatusPaid           OrderStatus = "paid"            // 已付款
	OrderStatusShipped        OrderStatus = "shipped"         // 已发货
	OrderStatusDelivered      OrderStatus = "delivered"       // 已送达
	OrderStatusCancelled      OrderStatus = "cancelled"       // 已取消
	OrderStatusRefunded       OrderStatus = "refunded"        // 已退款
)

func (s OrderStatus) String() string {
	return string(s)
}

func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusPendingPayment, OrderStatusPaid, OrderStatusShipped,
		OrderStatusDelivered, OrderStatusCancelled, OrderStatusRefunded:
		return true
	default:
		return false
	}
}

func (s OrderStatus) Text() string {
	switch s {
	case OrderStatusPendingPayment:
		return "待付款"
	case OrderStatusPaid:
		return "已付款"
	case OrderStatusShipped:
		return "已发货"
	case OrderStatusDelivered:
		return "已送达"
	case OrderStatusCancelled:
		return "已取消"
	case OrderStatusRefunded:
		return "已退款"
	default:
		return "未知"
	}
}

// OrderFulfillmentStatus 订单履约状态
type OrderFulfillmentStatus int8

const (
	OrderFulfillmentStatusPending        OrderFulfillmentStatus = 0 // 待发货
	OrderFulfillmentStatusPartialShipped OrderFulfillmentStatus = 1 // 部分发货
	OrderFulfillmentStatusShipped        OrderFulfillmentStatus = 2 // 已发货
	OrderFulfillmentStatusDelivered      OrderFulfillmentStatus = 3 // 已送达
)

func (s OrderFulfillmentStatus) Text() string {
	switch s {
	case OrderFulfillmentStatusPending:
		return "待发货"
	case OrderFulfillmentStatusPartialShipped:
		return "部分发货"
	case OrderFulfillmentStatusShipped:
		return "已发货"
	case OrderFulfillmentStatusDelivered:
		return "已送达"
	default:
		return "未知"
	}
}

// OrderRefundStatus 订单退款状态
type OrderRefundStatus int8

const (
	OrderRefundStatusNone      OrderRefundStatus = 0 // 无退款
	OrderRefundStatusPending   OrderRefundStatus = 1 // 待处理
	OrderRefundStatusApproved  OrderRefundStatus = 2 // 已批准
	OrderRefundStatusRejected  OrderRefundStatus = 3 // 已拒绝
	OrderRefundStatusCompleted OrderRefundStatus = 4 // 已完成
)

func (s OrderRefundStatus) Text() string {
	switch s {
	case OrderRefundStatusNone:
		return "无"
	case OrderRefundStatusPending:
		return "待处理"
	case OrderRefundStatusApproved:
		return "已批准"
	case OrderRefundStatusRejected:
		return "已拒绝"
	case OrderRefundStatusCompleted:
		return "已完成"
	default:
		return "未知"
	}
}

// ==================== Order Entity ====================

// Order 订单实体
type Order struct {
	ID                int64                  `gorm:"column:id;primaryKey"`
	TenantID          shared.TenantID        `gorm:"column:tenant_id;not null;index"`
	OrderNo           string                 `gorm:"column:order_no;not null;uniqueIndex:uk_order_no"`
	UserID            int64                  `gorm:"column:user_id;not null;index"`
	Status            OrderStatus            `gorm:"column:status;not null;default:'pending_payment';index"`
	FulfillmentStatus OrderFulfillmentStatus `gorm:"column:fulfillment_status;not null;default:0;index"`
	RefundStatus      OrderRefundStatus      `gorm:"column:refund_status;not null;default:0;index"`
	TotalAmount       decimal.Decimal        `gorm:"column:total_amount;type:decimal(19,4);not null"`    // 商品总金额
	DiscountAmount    decimal.Decimal        `gorm:"column:discount_amount;type:decimal(19,4);not null"` // 优惠金额
	ShippingFee       decimal.Decimal        `gorm:"column:shipping_fee;type:decimal(19,4);not null"`    // 运费
	PayAmount         decimal.Decimal        `gorm:"column:pay_amount;type:decimal(19,4);not null"`      // 实付金额
	Currency          string                 `gorm:"column:currency;not null;default:'CNY'"`
	MerchantRemark    string                 `gorm:"column:merchant_remark;not null;default:''"`                   // 商家内部备注
	Remark            string                 `gorm:"column:remark;not null;default:''"`                            // 用户备注
	OriginalAmount    decimal.Decimal        `gorm:"column:original_amount;type:decimal(19,4);not null;default:0"` // 改价前原金额
	AdjustAmount      decimal.Decimal        `gorm:"column:adjust_amount;type:decimal(19,4);not null;default:0"`   // 改价金额
	AdjustReason      string                 `gorm:"column:adjust_reason;not null;default:''"`                     // 改价原因
	AdjustedBy        int64                  `gorm:"column:adjusted_by;not null;default:0"`                        // 改价操作人ID
	AdjustedAt        *time.Time             `gorm:"column:adjusted_at"`                                           // 改价时间
	Version           int                    `gorm:"column:version;not null;default:1"`                            // 乐观锁版本号
	PaymentMethod     string                 `gorm:"column:payment_method;not null;default:''"`                    // 支付方式
	Source            string                 `gorm:"column:source;not null;default:''"`                            // 订单来源
	// Receiver info
	ReceiverName    string `gorm:"column:receiver_name;not null"`
	ReceiverPhone   string `gorm:"column:receiver_phone;not null"`
	ReceiverAddress string `gorm:"column:receiver_address;not null"`
	// Timestamps
	PaidAt      *time.Time `gorm:"column:paid_at"`
	ShippedAt   *time.Time `gorm:"column:shipped_at"`
	DeliveredAt *time.Time `gorm:"column:delivered_at"`
	CancelledAt *time.Time `gorm:"column:cancelled_at"`
	CancelledBy int64      `gorm:"column:cancelled_by;not null;default:0"` // 取消操作人ID
	// Audit info
	Audit     shared.AuditInfo `gorm:"embedded"`
	DeletedAt *int64           `gorm:"column:deleted_at;index"`
	// Relations
	Items []OrderItem `gorm:"foreignKey:OrderID"`
}

func (o *Order) TableName() string {
	return "orders"
}

// CanAdjustPrice 检查订单是否可以改价
// 只有待付款状态的订单可以改价
func (o *Order) CanAdjustPrice() bool {
	return o.Status == OrderStatusPendingPayment
}

// AdjustPrice 改价
// adjustAmount: 改价金额，正数为涨价，负数为降价
// 规则：|adjustAmount| <= originalAmount * 20%
func (o *Order) AdjustPrice(adjustAmount decimal.Decimal, reason string, adjustedBy int64) error {
	if !o.CanAdjustPrice() {
		return code.ErrOrderCannotAdjustPrice
	}

	if reason == "" {
		return code.ErrOrderAdjustReasonRequired
	}

	// Check if adjustment exceeds 20% limit
	maxAdjust := o.OriginalAmount.Mul(decimal.RequireFromString("0.2"))
	if adjustAmount.IsPositive() && adjustAmount.GreaterThan(maxAdjust) {
		return code.ErrOrderAdjustAmountExceed
	}
	if adjustAmount.IsNegative() && adjustAmount.Abs().GreaterThan(maxAdjust) {
		return code.ErrOrderAdjustAmountExceed
	}

	// Ensure pay amount doesn't go negative
	newPayAmount := o.OriginalAmount.Add(adjustAmount)
	if newPayAmount.IsNegative() {
		return code.ErrOrderAdjustAmountExceed
	}

	o.AdjustAmount = adjustAmount
	o.AdjustReason = reason
	o.AdjustedBy = adjustedBy
	now := time.Now().UTC()
	o.AdjustedAt = &now
	o.PayAmount = newPayAmount

	return nil
}

// UpdateRemark 更新商家备注
func (o *Order) UpdateRemark(remark string) {
	if len(remark) > 500 {
		remark = remark[:500]
	}
	o.MerchantRemark = remark
}

// IsPaid 检查订单是否已支付
func (o *Order) IsPaid() bool {
	return o.Status == OrderStatusPaid ||
		o.Status == OrderStatusShipped ||
		o.Status == OrderStatusDelivered
}

// ==================== OrderItem Entity ====================

// OrderItem 订单明细实体
type OrderItem struct {
	ID          int64           `gorm:"column:id;primaryKey"`
	TenantID    shared.TenantID `gorm:"column:tenant_id;not null;index"`
	OrderID     int64           `gorm:"column:order_id;not null;index"`
	ProductID   int64           `gorm:"column:product_id;not null;index"`
	SKUID       int64           `gorm:"column:sku_id;not null;index"`
	ProductName string          `gorm:"column:product_name;not null"`
	SKUName     string          `gorm:"column:sku_name;not null"`
	Image       string          `gorm:"column:image"`
	Quantity    int             `gorm:"column:quantity;not null"`
	UnitPrice   decimal.Decimal `gorm:"column:unit_price;type:decimal(19,4);not null"`  // 单价
	TotalPrice  decimal.Decimal `gorm:"column:total_price;type:decimal(19,4);not null"` // 总价
	Currency    string          `gorm:"column:currency;not null;default:'CNY'"`
	CreatedAt   time.Time       `gorm:"column:created_at;not null"`
}

func (oi *OrderItem) TableName() string {
	return "order_items"
}

// ==================== Query Types ====================

// OrderQuery 订单查询参数
type OrderQuery struct {
	shared.PageQuery
	OrderNo           string
	UserID            int64
	UserName          string
	Status            OrderStatus
	FulfillmentStatus OrderFulfillmentStatus
	RefundStatus      OrderRefundStatus
	StartTime         time.Time
	EndTime           time.Time
}

// ==================== Repository Interfaces ====================

// OrderRepository 订单仓储接口
type OrderRepository interface {
	// FindByID 根据ID查询订单
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Order, error)
	// FindByOrderNo 根据订单号查询订单
	FindByOrderNo(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderNo string) (*Order, error)
	// FindList 分页查询订单列表
	FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query OrderQuery) ([]*Order, int64, error)
	// UpdateWithVersion 带乐观锁的更新
	// 返回 ErrOrderVersionConflict 如果版本冲突
	UpdateWithVersion(ctx context.Context, db *gorm.DB, order *Order) error
	// UpdateRemark 更新商家备注
	UpdateRemark(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderID int64, remark string) error
	// CountTodayOrders 统计今日订单数
	CountTodayOrders(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (int64, error)
	// SumTodayGMV 统计今日GMV（已支付订单的总金额）
	SumTodayGMV(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (decimal.Decimal, error)
	// FindForExport 导出订单（最多10000条）
	FindForExport(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query OrderQuery) ([]*Order, error)
	// CountByStatus 按状态统计订单数量
	CountByStatus(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) ([]OrderStatusCount, error)
	// FindPendingOrders 查询待付款订单
	FindPendingOrders(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, limit int) ([]*Order, error)
	// CountPendingOrders 统计待付款订单数量
	CountPendingOrders(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (int64, error)
	// FindRecentOrders 查询最近创建的订单
	FindRecentOrders(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, limit int) ([]*Order, error)
	// FindRecentPaidOrders 查询最近已支付的订单
	FindRecentPaidOrders(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, limit int) ([]*Order, error)
	// SumGMVByDateRange 按日期范围统计GMV
	SumGMVByDateRange(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, start, end time.Time, statuses []OrderStatus) (decimal.Decimal, error)
	// FindTopProducts 查询热销商品
	FindTopProducts(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, startTime time.Time, limit int) ([]*TopProduct, error)
	// FindSalesTrend 查询销售趋势
	FindSalesTrend(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, startDate, endDate time.Time) ([]*DailySalesTrend, error)
}

// OrderItemRepository 订单明细仓储接口
type OrderItemRepository interface {
	// FindByOrderID 根据订单ID查询明细
	FindByOrderID(ctx context.Context, db *gorm.DB, orderID int64) ([]OrderItem, error)
	// FindByOrderIDs 批量查询订单明细
	FindByOrderIDs(ctx context.Context, db *gorm.DB, orderIDs []int64) (map[int64][]OrderItem, error)
}

// ==================== Export Types ====================

// OrderExportRow 订单导出行
type OrderExportRow struct {
	OrderNo           string
	Status            string
	FulfillmentStatus string
	RefundStatus      string
	TotalAmount       decimal.Decimal
	DiscountAmount    decimal.Decimal
	ShippingFee       decimal.Decimal
	PayAmount         decimal.Decimal
	ReceiverName      string
	ReceiverPhone     string
	ReceiverAddress   string
	PaymentMethod     string
	CreatedAt         time.Time
	PaidAt            *time.Time
}

// AdjustPriceResponse 改价响应
type AdjustPriceResponse struct {
	OrderID        int64           `json:"order_id"`
	OriginalAmount decimal.Decimal `json:"original_amount"`
	AdjustAmount   decimal.Decimal `json:"adjust_amount"`
	NewPayAmount   decimal.Decimal `json:"new_pay_amount"`
	AdjustReason   string          `json:"adjust_reason"`
	AdjustedAt     string          `json:"adjusted_at"`
}

// ==================== Dashboard Types ====================

// OrderStatusCount 订单状态统计
type OrderStatusCount struct {
	Status OrderStatus
	Count  int64
}

// TopProduct 热销商品
type TopProduct struct {
	ProductID   int64
	ProductName string
	Image       string
	Sales       int64
	Revenue     decimal.Decimal
}

// DailySalesTrend 日销售趋势
type DailySalesTrend struct {
	Date   string // 日期字符串，格式为 YYYY-MM-DD
	Sales  decimal.Decimal
	Orders int64
}
