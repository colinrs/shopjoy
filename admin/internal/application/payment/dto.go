package payment

import (
	"strconv"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/payment"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
)

// PaymentStatsDTO represents payment statistics
type PaymentStatsDTO struct {
	TodayReceived       string                  `json:"today_received"`
	TodayGrowth         string                  `json:"today_growth"`
	PeriodReceived      string                  `json:"period_received"`
	RefundAmount        string                  `json:"refund_amount"`
	RefundRate          string                  `json:"refund_rate"`
	Currency            string                  `json:"currency"`
	ChannelDistribution []ChannelDistributionDTO `json:"channel_distribution"`
	Stats               TransactionStatsDTO     `json:"stats"`
}

// ChannelDistributionDTO represents channel distribution
type ChannelDistributionDTO struct {
	Name    string `json:"name"`
	Percent int    `json:"percent"`
	Amount  string `json:"amount"`
	Count   int64  `json:"count"`
	Color   string `json:"color"`
}

// TransactionStatsDTO represents transaction statistics
type TransactionStatsDTO struct {
	Success int64 `json:"success"`
	Pending int64 `json:"pending"`
	Failed  int64 `json:"failed"`
}

// ListTransactionsRequest represents the request to list transactions
type ListTransactionsRequest struct {
	Page          int
	PageSize      int
	OrderID       int64
	TransactionID string
	PaymentMethod payment.PaymentMethod
	Status        payment.TransactionStatus
	StartTime     time.Time
	EndTime       time.Time
}

// ListTransactionsResponse represents the response for listing transactions
type ListTransactionsResponse struct {
	List     []*TransactionDTO `json:"list"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	Stats    TransactionStatsDTO `json:"stats"`
}

// TransactionDTO represents a payment transaction
type TransactionDTO struct {
	ID                   int64  `json:"id"`
	TransactionID        string `json:"transaction_id"`
	OrderID              string `json:"order_id"`
	OrderNo              string `json:"order_no"`
	PaymentMethod        string `json:"payment_method"`
	PaymentMethodText    string `json:"payment_method_text"`
	ChannelTransactionID string `json:"channel_transaction_id"`
	Amount               string `json:"amount"`
	Currency             string `json:"currency"`
	TransactionFee       string `json:"transaction_fee"`
	Status               int8   `json:"status"`
	StatusText           string `json:"status_text"`
	CreatedAt            string `json:"created_at"`
	PaidAt               string `json:"paid_at,optional"`
}

// OrderPaymentDTO represents payment details for an order
type OrderPaymentDTO struct {
	PaymentID         int64                `json:"payment_id"`
	PaymentNo         string               `json:"payment_no"`
	PaymentMethod     string               `json:"payment_method"`
	PaymentMethodText string               `json:"payment_method_text"`
	ChannelIntentID   string               `json:"channel_intent_id"`
	ChannelPaymentID  string               `json:"channel_payment_id"`
	Amount            string               `json:"amount"`
	Currency          string               `json:"currency"`
	TransactionFee    string               `json:"transaction_fee"`
	FeeCurrency       string               `json:"fee_currency"`
	Status            int8                 `json:"status"`
	StatusText        string               `json:"status_text"`
	PaidAt            string               `json:"paid_at,optional"`
	RefundedAmount    string               `json:"refunded_amount"`
	Refunds           []*PaymentRefundDTO  `json:"refunds"`
}

// PaymentRefundDTO represents a payment refund
type PaymentRefundDTO struct {
	ID              int64  `json:"id"`
	RefundNo        string `json:"refund_no"`
	ChannelRefundID string `json:"channel_refund_id"`
	Amount          string `json:"amount"`
	Currency        string `json:"currency"`
	Status          int8   `json:"status"`
	StatusText      string `json:"status_text"`
	ReasonType      string `json:"reason_type"`
	Reason          string `json:"reason"`
	RefundedAt      string `json:"refunded_at,optional"`
	CreatedAt       string `json:"created_at"`
}

// InitiateRefundRequest represents the request to initiate a refund
type InitiateRefundRequest struct {
	OrderID        int64
	IdempotencyKey string
	Amount         string
	Currency       string
	ReasonType     string
	Reason         string
}

// InitiateRefundResponse represents the response for initiating a refund
type InitiateRefundResponse struct {
	RefundID        int64  `json:"refund_id"`
	RefundNo        string `json:"refund_no"`
	Amount          string `json:"amount"`
	Currency        string `json:"currency"`
	Status          int8   `json:"status"`
	StatusText      string `json:"status_text"`
	ChannelRefundID string `json:"channel_refund_id,optional"`
}

// WebhookEvent represents a webhook event
type WebhookEvent struct {
	EventID    string
	EventType  string
	ResourceID string
	RawPayload string
}

// toTransactionDTO converts a PaymentTransaction to TransactionDTO
func toTransactionDTO(txn *payment.PaymentTransaction) *TransactionDTO {
	return &TransactionDTO{
		ID:                   txn.ID,
		TransactionID:        txn.TransactionID,
		OrderID:              formatOrderID(txn.OrderID),
		OrderNo:              "", // Would need to fetch from order service
		PaymentMethod:        string(txn.PaymentMethod),
		PaymentMethodText:    getPaymentMethodText(txn.PaymentMethod),
		ChannelTransactionID: txn.ChannelTransactionID,
		Amount:               shared.NewMoney(txn.Amount, txn.Currency).String(),
		Currency:             txn.Currency,
		TransactionFee:       shared.NewMoney(txn.TransactionFee, txn.Currency).String(),
		Status:               int8(txn.Status),
		StatusText:           txn.Status.String(),
		CreatedAt:            txn.CreatedAt.Format(time.RFC3339),
		PaidAt:               formatTimeToString(txn.PaidAt),
	}
}

// toPaymentRefundDTO converts a PaymentRefund to PaymentRefundDTO
func toPaymentRefundDTO(refund *payment.PaymentRefund) *PaymentRefundDTO {
	return &PaymentRefundDTO{
		ID:              refund.ID,
		RefundNo:        refund.RefundNo,
		ChannelRefundID: refund.ChannelRefundID,
		Amount:          shared.NewMoney(refund.Amount, refund.Currency).String(),
		Currency:        refund.Currency,
		Status:          int8(refund.Status),
		StatusText:      refund.Status.String(),
		ReasonType:      refund.ReasonType,
		Reason:          refund.Reason,
		RefundedAt:      formatTimeToString(refund.RefundedAt),
		CreatedAt:       refund.CreatedAt.Format(time.RFC3339),
	}
}

// formatOrderID formats an order ID to string
func formatOrderID(orderID int64) string {
	return strconv.FormatInt(orderID, 10)
}