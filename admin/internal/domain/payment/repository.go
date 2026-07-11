package payment

import (
	"context"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// PaymentRepository 支付仓储接口
type PaymentRepository interface {
	Create(ctx context.Context, db *gorm.DB, payment *Payment) error
	Update(ctx context.Context, db *gorm.DB, payment *Payment) error
	FindByID(ctx context.Context, db *gorm.DB,  id int64) (*Payment, error)
	FindByOrderID(ctx context.Context, db *gorm.DB,  orderID int64) (*Payment, error)
	FindByPaymentNo(ctx context.Context, db *gorm.DB,  paymentNo string) (*Payment, error)
	FindByChannelPaymentID(ctx context.Context, db *gorm.DB, channelPaymentID string) (*Payment, error)
	FindList(ctx context.Context, db *gorm.DB,  query PaymentQuery) ([]*Payment, int64, error)
}

// PaymentRefundRepository 退款仓储接口
type PaymentRefundRepository interface {
	Create(ctx context.Context, db *gorm.DB, refund *PaymentRefund) error
	Update(ctx context.Context, db *gorm.DB, refund *PaymentRefund) error
	FindByID(ctx context.Context, db *gorm.DB,  id int64) (*PaymentRefund, error)
	FindByRefundNo(ctx context.Context, db *gorm.DB,  refundNo string) (*PaymentRefund, error)
	FindByIdempotencyKey(ctx context.Context, db *gorm.DB, idempotencyKey string) (*PaymentRefund, error)
	FindByPaymentID(ctx context.Context, db *gorm.DB,  paymentID int64) ([]*PaymentRefund, error)
	FindList(ctx context.Context, db *gorm.DB,  query RefundQuery) ([]*PaymentRefund, int64, error)
	GetTotalRefundedAmount(ctx context.Context, db *gorm.DB,  paymentID int64) (decimal.Decimal, error)
}

// PaymentTransactionRepository 交易仓储接口
type PaymentTransactionRepository interface {
	Create(ctx context.Context, db *gorm.DB, txn *PaymentTransaction) error
	FindByID(ctx context.Context, db *gorm.DB,  id int64) (*PaymentTransaction, error)
	FindByTransactionID(ctx context.Context, db *gorm.DB,  txnID string) (*PaymentTransaction, error)
	FindList(ctx context.Context, db *gorm.DB,  query TransactionQuery) ([]*PaymentTransaction, int64, error)
	GetStats(ctx context.Context, db *gorm.DB) (success, pending, failed int64, err error)
}

// WebhookEventRepository Webhook事件仓储接口
type WebhookEventRepository interface {
	Create(ctx context.Context, db *gorm.DB, event *WebhookEvent) error
	FindByEventID(ctx context.Context, db *gorm.DB, eventID string) (*WebhookEvent, error)
	MarkProcessed(ctx context.Context, db *gorm.DB, id int64, processed int8, errorMsg string) error
}
