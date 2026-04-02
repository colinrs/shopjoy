// Package shared 共享内核 - 领域事件
package shared

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"
)

// EventType 事件类型
type EventType string

const (
	// 租户事件
	EventTenantCreated EventType = "tenant.created"
	EventTenantUpdated EventType = "tenant.updated"
	EventTenantDeleted EventType = "tenant.deleted"

	// 用户事件
	EventUserRegistered EventType = "user.registered"
	EventUserLogin      EventType = "user.login"
	EventUserLogout     EventType = "user.logout"
	EventUserUpdated    EventType = "user.updated"

	// 商品事件
	EventProductCreated EventType = "product.created"
	EventProductUpdated EventType = "product.updated"
	EventProductDeleted EventType = "product.deleted"
	EventProductOnSale  EventType = "product.on_sale"
	EventProductOffSale EventType = "product.off_sale"
	EventStockChanged   EventType = "stock.changed"
	EventStockLow       EventType = "stock.low"

	// 订单事件
	EventOrderCreated   EventType = "order.created"
	EventOrderPaid      EventType = "order.paid"
	EventOrderShipped   EventType = "order.shipped"
	EventOrderCompleted EventType = "order.completed"
	EventOrderCancelled EventType = "order.cancelled"
	EventOrderRefunded  EventType = "order.refunded"

	// 促销事件
	EventPromotionStarted EventType = "promotion.started"
	EventPromotionEnded   EventType = "promotion.ended"
	EventCouponUsed       EventType = "coupon.used"
	EventCouponExpired    EventType = "coupon.expired"

	// 支付事件
	EventPaymentSuccess  EventType = "payment.success"
	EventPaymentFailed   EventType = "payment.failed"
	EventPaymentRefunded EventType = "payment.refunded"
)

// DomainEvent 领域事件接口
type DomainEvent interface {
	EventType() EventType
	EventID() string
	OccurredAt() time.Time
	TenantID() TenantID
	Payload() []byte
}

// BaseEvent 基础事件实现
type BaseEvent struct {
	Type       EventType `json:"type"`
	ID         string    `json:"id"`
	Tenant     TenantID  `json:"tenant_id"`
	OccurredOn time.Time `json:"occurred_at"`
	Data       []byte    `json:"data"`
}

// EventType 返回事件类型
func (e BaseEvent) EventType() EventType {
	return e.Type
}

// EventID 返回事件ID
func (e BaseEvent) EventID() string {
	return e.ID
}

// OccurredAt 返回发生时间
func (e BaseEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// TenantID 返回租户ID
func (e BaseEvent) TenantID() TenantID {
	return e.Tenant
}

// Payload 返回事件数据
func (e BaseEvent) Payload() []byte {
	return e.Data
}

// NewBaseEvent 创建基础事件
func NewBaseEvent(eventType EventType, tenantID TenantID, payload interface{}) (BaseEvent, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return BaseEvent{}, err
	}

	return BaseEvent{
		Type:       eventType,
		ID:         GenerateEventID(),
		Tenant:     tenantID,
		OccurredOn: time.Now().UTC(),
		Data:       data,
	}, nil
}

// EventPublisher 事件发布接口
type EventPublisher interface {
	Publish(event DomainEvent) error
	PublishBatch(events []DomainEvent) error
}

// EventSubscriber 事件订阅接口
type EventSubscriber interface {
	Subscribe(eventType EventType, handler EventHandler) error
}

// EventHandler 事件处理器
type EventHandler func(event DomainEvent) error

// EventBus 事件总线接口
type EventBus interface {
	EventPublisher
	EventSubscriber
}

// GenerateEventID 生成事件ID（使用crypto/rand生成安全随机数）
func GenerateEventID() string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return fmt.Sprintf("evt_%d_%d", time.Now().UnixNano(), binary.BigEndian.Uint64(b[:]))
}

// ==================== 具体事件定义 ====================

// ProductCreatedEvent 商品创建事件
type ProductCreatedEvent struct {
	ProductID    int64  `json:"product_id"`
	Name         string `json:"name"`
	CategoryID   int64  `json:"category_id"`
	Price        int64  `json:"price"`
	InitialStock int    `json:"initial_stock"`
}

// OrderCreatedEvent 订单创建事件
type OrderCreatedEvent struct {
	OrderID     string      `json:"order_id"`
	UserID      int64       `json:"user_id"`
	TotalAmount int64       `json:"total_amount"`
	Currency    string      `json:"currency"`
	Items       []OrderItem `json:"items"`
}

// OrderItem 订单项（事件用）
type OrderItem struct {
	ProductID int64 `json:"product_id"`
	SKUId     int64 `json:"sku_id"`
	Quantity  int   `json:"quantity"`
	Price     int64 `json:"price"`
}

// StockChangedEvent 库存变更事件
type StockChangedEvent struct {
	ProductID int64  `json:"product_id"`
	SKUId     int64  `json:"sku_id"`
	OldStock  int    `json:"old_stock"`
	NewStock  int    `json:"new_stock"`
	Delta     int    `json:"delta"`
	Reason    string `json:"reason"`
}

// CouponUsedEvent 优惠券使用事件
type CouponUsedEvent struct {
	CouponID int64  `json:"coupon_id"`
	Code     string `json:"code"`
	UserID   int64  `json:"user_id"`
	OrderID  string `json:"order_id"`
	Discount int64  `json:"discount"`
}

// PaymentSuccessEvent 支付成功事件
type PaymentSuccessEvent struct {
	OrderID       string    `json:"order_id"`
	PaymentID     string    `json:"payment_id"`
	Amount        int64     `json:"amount"`
	Currency      string    `json:"currency"`
	PaymentMethod string    `json:"payment_method"`
	PaidAt        time.Time `json:"paid_at"`
}
