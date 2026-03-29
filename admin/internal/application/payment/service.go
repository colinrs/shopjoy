package payment

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/payment"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/snowflake"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Service is the payment application service interface
type Service interface {
	// GetPaymentStats returns payment statistics
	GetPaymentStats(ctx context.Context, tenantID shared.TenantID, period string) (*PaymentStatsDTO, error)
	// ListTransactions returns a list of payment transactions
	ListTransactions(ctx context.Context, tenantID shared.TenantID, req ListTransactionsRequest) (*ListTransactionsResponse, error)
	// GetTransaction returns a single transaction by ID
	GetTransaction(ctx context.Context, tenantID shared.TenantID, id int64) (*TransactionDTO, error)
	// GetOrderPayment returns payment details for an order
	GetOrderPayment(ctx context.Context, tenantID shared.TenantID, orderID int64) (*OrderPaymentDTO, error)
	// InitiateRefund initiates a refund for an order
	InitiateRefund(ctx context.Context, tenantID shared.TenantID, adminID int64, req InitiateRefundRequest) (*InitiateRefundResponse, error)
	// HandleWebhook handles Stripe webhook events
	HandleWebhook(ctx context.Context, event *WebhookEvent) error
}

type service struct {
	db                *gorm.DB
	paymentRepo       payment.PaymentRepository
	refundRepo        payment.PaymentRefundRepository
	transactionRepo   payment.PaymentTransactionRepository
	webhookEventRepo  payment.WebhookEventRepository
	idGen             snowflake.Snowflake
}

// NewService creates a new payment application service
func NewService(
	db *gorm.DB,
	paymentRepo payment.PaymentRepository,
	refundRepo payment.PaymentRefundRepository,
	transactionRepo payment.PaymentTransactionRepository,
	webhookEventRepo payment.WebhookEventRepository,
	idGen snowflake.Snowflake,
) Service {
	return &service{
		db:               db,
		paymentRepo:      paymentRepo,
		refundRepo:       refundRepo,
		transactionRepo:  transactionRepo,
		webhookEventRepo: webhookEventRepo,
		idGen:            idGen,
	}
}

// GetPaymentStats returns payment statistics
func (s *service) GetPaymentStats(ctx context.Context, tenantID shared.TenantID, period string) (*PaymentStatsDTO, error) {
	// Parse period to get time range
	now := time.Now().UTC()
	var startTime time.Time
	switch period {
	case "7d":
		startTime = now.AddDate(0, 0, -7)
	case "30d":
		startTime = now.AddDate(0, 0, -30)
	case "90d":
		startTime = now.AddDate(0, 0, -90)
	default:
		startTime = now.AddDate(0, 0, -7)
	}

	// Get transaction stats
	success, pending, failed, err := s.transactionRepo.GetStats(ctx, s.db, tenantID)
	if err != nil {
		return nil, err
	}

	// Get payment stats for the period
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	yesterdayStart := todayStart.AddDate(0, 0, -1)

	// Today's transactions
	todayQuery := payment.TransactionQuery{
		StartTime: todayStart,
		EndTime:   now,
	}
	todayTxns, _, err := s.transactionRepo.FindList(ctx, s.db, tenantID, todayQuery)
	if err != nil {
		return nil, err
	}

	var todayReceived decimal.Decimal
	for _, txn := range todayTxns {
		if txn.Status == payment.TransactionStatusSucceeded {
			todayReceived = todayReceived.Add(txn.Amount)
		}
	}

	// Yesterday's transactions for growth calculation
	yesterdayQuery := payment.TransactionQuery{
		StartTime: yesterdayStart,
		EndTime:   todayStart,
	}
	yesterdayTxns, _, err := s.transactionRepo.FindList(ctx, s.db, tenantID, yesterdayQuery)
	if err != nil {
		return nil, err
	}

	var yesterdayReceived decimal.Decimal
	for _, txn := range yesterdayTxns {
		if txn.Status == payment.TransactionStatusSucceeded {
			yesterdayReceived = yesterdayReceived.Add(txn.Amount)
		}
	}

	// Calculate growth rate
	var todayGrowth float64
	if yesterdayReceived.IsPositive() {
		todayGrowth, _ = todayReceived.Sub(yesterdayReceived).Div(yesterdayReceived).Mul(decimal.NewFromInt(100)).Float64()
	} else if todayReceived.IsPositive() {
		todayGrowth = 100
	}

	// Period transactions
	periodQuery := payment.TransactionQuery{
		StartTime: startTime,
		EndTime:   now,
	}
	periodTxns, _, err := s.transactionRepo.FindList(ctx, s.db, tenantID, periodQuery)
	if err != nil {
		return nil, err
	}

	var periodReceived decimal.Decimal
	channelMap := make(map[string]struct {
		Count  int64
		Amount decimal.Decimal
	})
	for _, txn := range periodTxns {
		if txn.Status == payment.TransactionStatusSucceeded {
			periodReceived = periodReceived.Add(txn.Amount)
			data := channelMap[string(txn.PaymentMethod)]
			data.Count++
			data.Amount = data.Amount.Add(txn.Amount)
			channelMap[string(txn.PaymentMethod)] = data
		}
	}

	// Calculate channel distribution
	channelDistribution := make([]ChannelDistributionDTO, 0, len(channelMap))
	for method, data := range channelMap {
		percent := 0
		if periodReceived.IsPositive() {
			p, _ := data.Amount.Div(periodReceived).Mul(decimal.NewFromInt(100)).Float64()
			percent = int(p)
		}
		channelDistribution = append(channelDistribution, ChannelDistributionDTO{
			Name:    getPaymentMethodText(payment.PaymentMethod(method)),
			Percent: percent,
			Amount:  shared.NewMoney(data.Amount, "USD").String(),
			Count:   data.Count,
			Color:   getChannelColor(payment.PaymentMethod(method)),
		})
	}

	// Calculate refund amount (simplified - would need actual refund tracking)
	refundAmount := int64(0)
	refundRate := "0.00%"

	return &PaymentStatsDTO{
		TodayReceived:       shared.NewMoney(todayReceived, "USD").String(),
		TodayGrowth:         formatPercentage(todayGrowth),
		PeriodReceived:      shared.NewMoney(periodReceived, "USD").String(),
		RefundAmount:        shared.NewMoney(decimal.NewFromInt(refundAmount), "USD").String(),
		RefundRate:          refundRate,
		Currency:            "USD",
		ChannelDistribution: channelDistribution,
		Stats: TransactionStatsDTO{
			Success: success,
			Pending: pending,
			Failed:  failed,
		},
	}, nil
}

// ListTransactions returns a list of payment transactions
func (s *service) ListTransactions(ctx context.Context, tenantID shared.TenantID, req ListTransactionsRequest) (*ListTransactionsResponse, error) {
	query := payment.TransactionQuery{
		PageQuery: shared.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		TransactionID: req.TransactionID,
		PaymentMethod: req.PaymentMethod,
	}

	if req.Status.IsValid() {
		query.Status = req.Status
	}
	if !req.StartTime.IsZero() {
		query.StartTime = req.StartTime
	}
	if !req.EndTime.IsZero() {
		query.EndTime = req.EndTime
	}
	if req.OrderID > 0 {
		query.OrderID = req.OrderID
	}

	query.PageQuery.Validate()

	transactions, total, err := s.transactionRepo.FindList(ctx, s.db, tenantID, query)
	if err != nil {
		return nil, err
	}

	list := make([]*TransactionDTO, len(transactions))
	for i, txn := range transactions {
		list[i] = toTransactionDTO(txn)
	}

	// Get stats
	success, pending, failed, err := s.transactionRepo.GetStats(ctx, s.db, tenantID)
	if err != nil {
		return nil, err
	}

	return &ListTransactionsResponse{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Stats: TransactionStatsDTO{
			Success: success,
			Pending: pending,
			Failed:  failed,
		},
	}, nil
}

// GetTransaction returns a single transaction by ID
func (s *service) GetTransaction(ctx context.Context, tenantID shared.TenantID, id int64) (*TransactionDTO, error) {
	txn, err := s.transactionRepo.FindByID(ctx, s.db, tenantID, id)
	if err != nil {
		return nil, err
	}
	return toTransactionDTO(txn), nil
}

// GetOrderPayment returns payment details for an order
func (s *service) GetOrderPayment(ctx context.Context, tenantID shared.TenantID, orderID int64) (*OrderPaymentDTO, error) {
	paymentEntity, err := s.paymentRepo.FindByOrderID(ctx, s.db, tenantID, orderID)
	if err != nil {
		return nil, err
	}

	// Get refunds for this payment
	refunds, err := s.refundRepo.FindByPaymentID(ctx, s.db, tenantID, paymentEntity.ID)
	if err != nil {
		return nil, err
	}

	// Calculate total refunded amount
	totalRefunded, err := s.refundRepo.GetTotalRefundedAmount(ctx, s.db, tenantID, paymentEntity.ID)
	if err != nil {
		return nil, err
	}

	refundDTOs := make([]*PaymentRefundDTO, len(refunds))
	for i, refund := range refunds {
		refundDTOs[i] = toPaymentRefundDTO(refund)
	}

	return &OrderPaymentDTO{
		PaymentID:         paymentEntity.ID,
		PaymentNo:         paymentEntity.PaymentNo,
		PaymentMethod:     string(paymentEntity.PaymentMethod),
		PaymentMethodText: getPaymentMethodText(paymentEntity.PaymentMethod),
		ChannelIntentID:   paymentEntity.ChannelIntentID,
		ChannelPaymentID:  paymentEntity.ChannelPaymentID,
		Amount:            shared.NewMoney(paymentEntity.Amount, paymentEntity.Currency).String(),
		Currency:          paymentEntity.Currency,
		TransactionFee:    shared.NewMoney(paymentEntity.TransactionFee, paymentEntity.FeeCurrency).String(),
		FeeCurrency:       paymentEntity.FeeCurrency,
		Status:            int8(paymentEntity.Status),
		StatusText:        paymentEntity.Status.String(),
		PaidAt:            timestampToString(paymentEntity.PaidAt),
		RefundedAmount:    shared.NewMoney(totalRefunded, paymentEntity.Currency).String(),
		Refunds:           refundDTOs,
	}, nil
}

// InitiateRefund initiates a refund for an order
func (s *service) InitiateRefund(ctx context.Context, tenantID shared.TenantID, adminID int64, req InitiateRefundRequest) (*InitiateRefundResponse, error) {
	// Parse amount from string to int64 (cents)
	// The amount string can be "99.99 USD" or just "99.99"
	amount, err := parseMoneyString(req.Amount)
	if err != nil {
		return nil, code.ErrPaymentInvalidAmount
	}

	// Validate refund amount
	if amount.IsZero() || amount.IsNegative() {
		return nil, code.ErrPaymentInvalidAmount
	}

	// Check idempotency key
	existingRefund, err := s.refundRepo.FindByIdempotencyKey(ctx, s.db, req.IdempotencyKey)
	if err != nil {
		return nil, err
	}
	if existingRefund != nil {
		// Return existing refund details for idempotency
		return &InitiateRefundResponse{
			RefundID:        existingRefund.ID,
			RefundNo:        existingRefund.RefundNo,
			Amount:          shared.NewMoney(existingRefund.Amount, existingRefund.Currency).String(),
			Currency:        existingRefund.Currency,
			Status:          int8(existingRefund.Status),
			StatusText:      existingRefund.Status.String(),
			ChannelRefundID: existingRefund.ChannelRefundID,
		}, nil
	}

	// Get payment for the order
	paymentEntity, err := s.paymentRepo.FindByOrderID(ctx, s.db, tenantID, req.OrderID)
	if err != nil {
		return nil, err
	}

	// Validate payment status
	if !paymentEntity.CanRefund() {
		return nil, code.ErrPaymentOrderNotPaid
	}

	// Check currency match
	if req.Currency != "" && req.Currency != paymentEntity.Currency {
		return nil, code.ErrRefundCurrencyMismatch
	}

	// Calculate refundable amount
	totalRefunded, err := s.refundRepo.GetTotalRefundedAmount(ctx, s.db, tenantID, paymentEntity.ID)
	if err != nil {
		return nil, err
	}

	refundableAmount := paymentEntity.Amount.Sub(totalRefunded)
	if amount.GreaterThan(refundableAmount) {
		return nil, code.ErrPaymentRefundAmountExceeded
	}

	// Validate refund reason
	if req.ReasonType == "" {
		return nil, code.ErrPaymentRefundReasonRequired
	}

	// Create refund record
	refundID, err := s.idGen.NextID(ctx)
	if err != nil {
		return nil, err
	}

	refund := payment.NewPaymentRefund(
		tenantID,
		req.OrderID,
		paymentEntity.ID,
		req.IdempotencyKey,
		amount,
		paymentEntity.Currency,
		req.ReasonType,
		req.Reason,
		adminID,
	)
	refund.ID = refundID

	// Start database transaction for refund creation, update, and payment status update
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Ensure transaction is rolled back on error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Save refund
	if err := s.refundRepo.Create(ctx, tx, refund); err != nil {
		tx.Rollback()
		return nil, err
	}

	// TODO: Call Stripe API to initiate refund (actual Stripe integration will be separate)
	// For now, mark as succeeded directly
	channelRefundID := "re_stub_" + refund.RefundNo
	refund.MarkSucceeded(channelRefundID, decimal.Zero)
	if err := s.refundRepo.Update(ctx, tx, refund); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Update payment status
	isPartial := amount.LessThan(refundableAmount)
	paymentEntity.MarkRefunded(isPartial)
	if err := s.paymentRepo.Update(ctx, tx, paymentEntity); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &InitiateRefundResponse{
		RefundID:        refund.ID,
		RefundNo:        refund.RefundNo,
		Amount:          shared.NewMoney(refund.Amount, refund.Currency).String(),
		Currency:        refund.Currency,
		Status:          int8(refund.Status),
		StatusText:      refund.Status.String(),
		ChannelRefundID: refund.ChannelRefundID,
	}, nil
}

// HandleWebhook handles Stripe webhook events
func (s *service) HandleWebhook(ctx context.Context, event *WebhookEvent) error {
	// Check if event already processed
	existingEvent, err := s.webhookEventRepo.FindByEventID(ctx, s.db, event.EventID)
	if err != nil {
		return err
	}
	if existingEvent != nil && existingEvent.IsProcessed() {
		// Event already processed, skip
		return nil
	}

	// Create webhook event record
	webhookEvent := &payment.WebhookEvent{
		EventID:    event.EventID,
		EventType:  event.EventType,
		ResourceID: event.ResourceID,
		RawPayload: event.RawPayload,
	}

	if existingEvent == nil {
		if err := s.webhookEventRepo.Create(ctx, s.db, webhookEvent); err != nil {
			return err
		}
	} else {
		webhookEvent = existingEvent
	}

	// Process event based on type
	var processErr error
	switch event.EventType {
	case "payment_intent.succeeded":
		processErr = s.handlePaymentIntentSucceeded(ctx, event.ResourceID, event.RawPayload)
	case "payment_intent.payment_failed":
		processErr = s.handlePaymentIntentFailed(ctx, event.ResourceID, event.RawPayload)
	case "charge.refunded":
		processErr = s.handleChargeRefunded(ctx, event.ResourceID, event.RawPayload)
	default:
		// Unknown event type, mark as processed
	}

	if processErr != nil {
		// Mark as failed
		s.webhookEventRepo.MarkProcessed(ctx, s.db, webhookEvent.ID, 2, processErr.Error())
		return processErr
	}

	// Mark as processed
	return s.webhookEventRepo.MarkProcessed(ctx, s.db, webhookEvent.ID, 1, "")
}

func (s *service) handlePaymentIntentSucceeded(ctx context.Context, paymentIntentID string, rawPayload string) error {
	// Find payment by channel intent ID
	paymentEntity, err := s.paymentRepo.FindByChannelPaymentID(ctx, s.db, paymentIntentID)
	if err != nil {
		return err
	}

	// Mark payment as succeeded
	// Fee calculation would come from webhook payload in real implementation
	paymentEntity.MarkSuccess(paymentIntentID, decimal.Zero, paymentEntity.Currency)
	return s.paymentRepo.Update(ctx, s.db, paymentEntity)
}

func (s *service) handlePaymentIntentFailed(ctx context.Context, paymentIntentID string, rawPayload string) error {
	// Find payment by channel intent ID
	paymentEntity, err := s.paymentRepo.FindByChannelPaymentID(ctx, s.db, paymentIntentID)
	if err != nil {
		return err
	}

	// Mark payment as failed
	paymentEntity.MarkFailed("Payment failed via webhook")
	return s.paymentRepo.Update(ctx, s.db, paymentEntity)
}

func (s *service) handleChargeRefunded(ctx context.Context, chargeID string, rawPayload string) error {
	// Find payment by channel payment ID
	paymentEntity, err := s.paymentRepo.FindByChannelPaymentID(ctx, s.db, chargeID)
	if err != nil {
		return err
	}

	// Find refund by channel refund ID
	// In real implementation, would parse webhook payload to get refund details
	// For now, just update payment status to partially refunded
	paymentEntity.MarkRefunded(true)
	return s.paymentRepo.Update(ctx, s.db, paymentEntity)
}

// Helper functions

func getPaymentMethodText(method payment.PaymentMethod) string {
	switch method {
	case payment.PaymentMethodStripe:
		return "Stripe"
	case payment.PaymentMethodAlipay:
		return "Alipay"
	case payment.PaymentMethodWechat:
		return "WeChat Pay"
	default:
		return string(method)
	}
}

func getChannelColor(method payment.PaymentMethod) string {
	switch method {
	case payment.PaymentMethodStripe:
		return "#635BFF"
	case payment.PaymentMethodAlipay:
		return "#1677FF"
	case payment.PaymentMethodWechat:
		return "#07C160"
	default:
		return "#999999"
	}
}

func formatTimeToString(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

func timestampToString(timestamp *time.Time) string {
	if timestamp == nil {
		return ""
	}
	return timestamp.Format(time.RFC3339)
}

func formatPercentage(value float64) string {
	if value >= 0 {
		return fmt.Sprintf("+%.1f%%", value)
	}
	return fmt.Sprintf("%.1f%%", value)
}

// parseMoneyString parses a money string to decimal.Decimal
// Supports formats: "99.99 USD", "99.99", "100"
func parseMoneyString(s string) (decimal.Decimal, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return decimal.Zero, nil
	}

	// Remove currency suffix if present
	parts := strings.Split(s, " ")
	amountStr := parts[0]

	// Parse using decimal for precision
	d, err := decimal.NewFromString(amountStr)
	if err != nil {
		return decimal.Zero, err
	}
	return d, nil
}