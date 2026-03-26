package payment

import (
	"fmt"
	"time"

	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
)

// ==================== Enums ====================

// PaymentStatus 支付状态
type PaymentStatus int

const (
	PaymentStatusPending         PaymentStatus = iota // 0 待支付
	PaymentStatusProcessing                           // 1 支付处理中
	PaymentStatusSuccess                              // 2 支付成功
	PaymentStatusFailed                               // 3 支付失败
	PaymentStatusCancelled                            // 4 已取消
	PaymentStatusRefunded                             // 5 已退款
	PaymentStatusPartiallyRefunded                    // 6 部分退款
	PaymentStatusRequiresAction                       // 7 待验证 (3D Secure)
)

func (s PaymentStatus) String() string {
	switch s {
	case PaymentStatusPending:
		return "pending"
	case PaymentStatusProcessing:
		return "processing"
	case PaymentStatusSuccess:
		return "success"
	case PaymentStatusFailed:
		return "failed"
	case PaymentStatusCancelled:
		return "cancelled"
	case PaymentStatusRefunded:
		return "refunded"
	case PaymentStatusPartiallyRefunded:
		return "partially_refunded"
	case PaymentStatusRequiresAction:
		return "requires_action"
	default:
		return "unknown"
	}
}

func (s PaymentStatus) IsValid() bool {
	return s >= PaymentStatusPending && s <= PaymentStatusRequiresAction
}

// TransactionStatus 交易状态
type TransactionStatus int

const (
	TransactionStatusPending TransactionStatus = iota // 0 待处理
	TransactionStatusSucceeded                        // 1 成功
	TransactionStatusFailed                           // 2 失败
)

func (s TransactionStatus) String() string {
	switch s {
	case TransactionStatusPending:
		return "pending"
	case TransactionStatusSucceeded:
		return "succeeded"
	case TransactionStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

func (s TransactionStatus) IsValid() bool {
	return s >= TransactionStatusPending && s <= TransactionStatusFailed
}

// PaymentRefundStatus 退款状态
type PaymentRefundStatus int

const (
	PaymentRefundStatusPending PaymentRefundStatus = iota // 0 待处理
	PaymentRefundStatusSucceeded                           // 1 退款成功
	PaymentRefundStatusFailed                              // 2 退款失败
)

func (s PaymentRefundStatus) String() string {
	switch s {
	case PaymentRefundStatusPending:
		return "pending"
	case PaymentRefundStatusSucceeded:
		return "succeeded"
	case PaymentRefundStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

func (s PaymentRefundStatus) IsValid() bool {
	return s >= PaymentRefundStatusPending && s <= PaymentRefundStatusFailed
}

// PaymentMethod 支付方式
type PaymentMethod string

const (
	PaymentMethodStripe PaymentMethod = "stripe"
	PaymentMethodAlipay PaymentMethod = "alipay"
	PaymentMethodWechat PaymentMethod = "wechat"
)

// SupportedCurrencies for Stripe (Phase 1)
var SupportedCurrencies = map[string]bool{
	"USD": true,
	"EUR": true,
	"GBP": true,
	"JPY": true,
	"SGD": true,
	"HKD": true,
}

// ==================== Payment ====================

// Payment 支付记录
type Payment struct {
	ID               int64            `gorm:"column:id;primaryKey"`
	TenantID         shared.TenantID  `gorm:"column:tenant_id;not null;index:idx_tenant_order"`
	OrderID          int64            `gorm:"column:order_id;not null;index:idx_tenant_order"`
	PaymentNo        string           `gorm:"column:payment_no;not null;uniqueIndex:uk_payment_no"`
	PaymentMethod    PaymentMethod    `gorm:"column:payment_method;not null"`
	ChannelIntentID  string           `gorm:"column:channel_intent_id;not null;default:'';index"`
	ChannelPaymentID string           `gorm:"column:channel_payment_id;not null;default:'';index"`
	Amount           int64            `gorm:"column:amount;not null;default:0"` // 分
	Currency         string           `gorm:"column:currency;not null;default:'USD'"`
	Status           PaymentStatus    `gorm:"column:status;not null;default:0;index"`
	TransactionFee   int64            `gorm:"column:transaction_fee;not null;default:0"`
	FeeCurrency      string           `gorm:"column:fee_currency;not null;default:'USD'"`
	PaidAt           int64            `gorm:"column:paid_at"`
	FailedAt         int64            `gorm:"column:failed_at"`
	FailedReason     string           `gorm:"column:failed_reason;not null;default:''"`
	Audit            shared.AuditInfo `gorm:"embedded"`
	DeletedAt        *int64          `gorm:"column:deleted_at;index"`
}

func (p *Payment) TableName() string {
	return "order_payments"
}

// NewPayment 创建支付记录
func NewPayment(tenantID shared.TenantID, orderID int64, method PaymentMethod, amount int64, currency string) *Payment {
	now := time.Now().Unix()
	return &Payment{
		TenantID:      tenantID,
		OrderID:       orderID,
		PaymentNo:     GeneratePaymentNo(tenantID, 1),
		PaymentMethod: method,
		Amount:        amount,
		Currency:      currency,
		Status:        PaymentStatusPending,
		FeeCurrency:   currency,
		Audit: shared.AuditInfo{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

// MarkSuccess 标记支付成功
func (p *Payment) MarkSuccess(channelPaymentID string, fee int64, feeCurrency string) error {
	if p.Status != PaymentStatusPending && p.Status != PaymentStatusProcessing {
		return code.ErrPaymentAlreadyPaid
	}
	now := time.Now().Unix()
	p.Status = PaymentStatusSuccess
	p.ChannelPaymentID = channelPaymentID
	p.TransactionFee = fee
	p.FeeCurrency = feeCurrency
	p.PaidAt = now
	p.Audit.UpdatedAt = now
	return nil
}

// MarkFailed 标记支付失败
func (p *Payment) MarkFailed(reason string) error {
	if p.Status != PaymentStatusPending && p.Status != PaymentStatusProcessing {
		return code.ErrPaymentAlreadyPaid
	}
	now := time.Now().Unix()
	p.Status = PaymentStatusFailed
	p.FailedAt = now
	p.FailedReason = reason
	p.Audit.UpdatedAt = now
	return nil
}

// MarkRefunded 标记已退款
func (p *Payment) MarkRefunded(partial bool) error {
	if p.Status != PaymentStatusSuccess && p.Status != PaymentStatusPartiallyRefunded {
		return code.ErrPaymentOrderNotPaid
	}
	now := time.Now().Unix()
	if partial {
		p.Status = PaymentStatusPartiallyRefunded
	} else {
		p.Status = PaymentStatusRefunded
	}
	p.Audit.UpdatedAt = now
	return nil
}

// CanRefund 是否可以退款
func (p *Payment) CanRefund() bool {
	return p.Status == PaymentStatusSuccess || p.Status == PaymentStatusPartiallyRefunded
}

// GeneratePaymentNo 生成支付单号
func GeneratePaymentNo(tenantID shared.TenantID, sequence int) string {
	dateStr := time.Now().UTC().Format("20060102")
	return fmt.Sprintf("PAY%s%06d", dateStr, sequence)
}

// ==================== PaymentRefund ====================

// PaymentRefund 支付退款记录
type PaymentRefund struct {
	ID                  int64                `gorm:"column:id;primaryKey"`
	TenantID            shared.TenantID      `gorm:"column:tenant_id;not null;index:idx_tenant_order"`
	OrderID             int64                `gorm:"column:order_id;not null;index:idx_tenant_order"`
	PaymentID           int64                `gorm:"column:payment_id;not null;index"`
	FulfillmentRefundID int64                `gorm:"column:fulfillment_refund_id"`
	RefundNo            string               `gorm:"column:refund_no;not null;uniqueIndex:uk_refund_no"`
	IdempotencyKey      string               `gorm:"column:idempotency_key;not null;default:'';uniqueIndex:uk_idempotency_key"`
	ChannelRefundID     string               `gorm:"column:channel_refund_id;not null;default:'';index"`
	Amount              int64                `gorm:"column:amount;not null;default:0"`
	Currency            string               `gorm:"column:currency;not null;default:'USD'"`
	RefundFee           int64                `gorm:"column:refund_fee;not null;default:0"`
	Status              PaymentRefundStatus  `gorm:"column:status;not null;default:0"`
	ReasonType          string               `gorm:"column:reason_type;not null;default:''"`
	Reason              string               `gorm:"column:reason;not null;default:''"`
	RefundedAt          int64                `gorm:"column:refunded_at"`
	CreatedAt           int64                `gorm:"column:created_at;not null"`
	CreatedBy           int64                `gorm:"column:created_by;not null;default:0"`
	DeletedAt           *int64              `gorm:"column:deleted_at;index"`
}

func (r *PaymentRefund) TableName() string {
	return "payment_refunds"
}

// NewPaymentRefund 创建退款记录
func NewPaymentRefund(tenantID shared.TenantID, orderID, paymentID int64, idempotencyKey string, amount int64, currency, reasonType, reason string, createdBy int64) *PaymentRefund {
	return &PaymentRefund{
		TenantID:       tenantID,
		OrderID:        orderID,
		PaymentID:      paymentID,
		RefundNo:       GenerateRefundNo(tenantID, 1),
		IdempotencyKey: idempotencyKey,
		Amount:         amount,
		Currency:       currency,
		Status:         PaymentRefundStatusPending,
		ReasonType:     reasonType,
		Reason:         reason,
		CreatedAt:      time.Now().Unix(),
		CreatedBy:      createdBy,
	}
}

// MarkSucceeded 标记退款成功
func (r *PaymentRefund) MarkSucceeded(channelRefundID string, refundFee int64) error {
	now := time.Now().Unix()
	r.Status = PaymentRefundStatusSucceeded
	r.ChannelRefundID = channelRefundID
	r.RefundFee = refundFee
	r.RefundedAt = now
	return nil
}

// MarkFailed 标记退款失败
func (r *PaymentRefund) MarkFailed() error {
	r.Status = PaymentRefundStatusFailed
	return nil
}

// GenerateRefundNo 生成退款单号
func GenerateRefundNo(tenantID shared.TenantID, sequence int) string {
	dateStr := time.Now().UTC().Format("20060102")
	return fmt.Sprintf("PRF%s%06d", dateStr, sequence)
}

// ==================== PaymentTransaction ====================

// PaymentTransaction 支付交易记录
type PaymentTransaction struct {
	ID                   int64             `gorm:"column:id;primaryKey"`
	TenantID             shared.TenantID   `gorm:"column:tenant_id;not null;index:idx_tenant_order"`
	OrderID              int64             `gorm:"column:order_id;not null;index:idx_tenant_order"`
	PaymentID            int64             `gorm:"column:payment_id;not null;index"`
	TransactionID        string            `gorm:"column:transaction_id;not null;uniqueIndex:uk_transaction_id"`
	PaymentMethod        PaymentMethod     `gorm:"column:payment_method;not null"`
	ChannelTransactionID string            `gorm:"column:channel_transaction_id;not null;default:'';index"`
	Amount               int64             `gorm:"column:amount;not null;default:0"`
	Currency             string            `gorm:"column:currency;not null;default:'USD'"`
	Status               TransactionStatus `gorm:"column:status;not null;default:0"`
	TransactionFee       int64             `gorm:"column:transaction_fee;not null;default:0"`
	PaidAt               int64             `gorm:"column:paid_at"`
	FailedReason         string            `gorm:"column:failed_reason;not null;default:''"`
	CreatedAt            int64             `gorm:"column:created_at;not null"`
	DeletedAt            *int64           `gorm:"column:deleted_at;index"`
}

func (t *PaymentTransaction) TableName() string {
	return "payment_transactions"
}

// ==================== WebhookEvent ====================

// WebhookEvent Webhook事件记录
type WebhookEvent struct {
	ID           int64            `gorm:"column:id;primaryKey"`
	TenantID     shared.TenantID  `gorm:"column:tenant_id;not null;default:0"`
	EventID      string           `gorm:"column:event_id;not null;uniqueIndex:uk_event_id"`
	EventType    string           `gorm:"column:event_type;not null;index:idx_tenant_event"`
	ResourceID   string           `gorm:"column:resource_id;not null;default:'';index"`
	Processed    int8             `gorm:"column:processed;not null;default:0"` // 0=pending, 1=processed, 2=failed
	RawPayload   string           `gorm:"column:raw_payload;type:text"`
	ErrorMessage string           `gorm:"column:error_message;not null;default:''"`
	CreatedAt    int64            `gorm:"column:created_at;not null"`
	ProcessedAt  int64            `gorm:"column:processed_at"`
}

func (e *WebhookEvent) TableName() string {
	return "webhook_events"
}

// IsProcessed 是否已处理
func (e *WebhookEvent) IsProcessed() bool {
	return e.Processed == 1
}

// ==================== Query Types ====================

// PaymentQuery 支付查询参数
type PaymentQuery struct {
	shared.PageQuery
	OrderID   int64
	Status    PaymentStatus
	StartTime time.Time
	EndTime   time.Time
}

// TransactionQuery 交易查询参数
type TransactionQuery struct {
	shared.PageQuery
	OrderID       int64
	TransactionID string
	PaymentMethod PaymentMethod
	Status        TransactionStatus
	StartTime     time.Time
	EndTime       time.Time
}

// RefundQuery 退款查询参数
type RefundQuery struct {
	shared.PageQuery
	OrderID   int64
	Status    PaymentRefundStatus
	StartTime time.Time
	EndTime   time.Time
}

// PaymentStats 支付统计
type PaymentStats struct {
	TodayReceived  int64   `json:"today_received"`
	TodayGrowth    float64 `json:"today_growth"`
	PeriodReceived int64   `json:"period_received"`
	RefundAmount   int64   `json:"refund_amount"`
	RefundRate     float64 `json:"refund_rate"`
	Currency       string  `json:"currency"`
}

// ChannelDistribution 渠道分布
type ChannelDistribution struct {
	Name    string `json:"name"`
	Percent int    `json:"percent"`
	Amount  int64  `json:"amount"`
	Count   int64  `json:"count"`
	Color   string `json:"color"`
}