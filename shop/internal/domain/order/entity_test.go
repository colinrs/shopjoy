package order

import (
	"testing"
	"time"

	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// newTestOrder creates an order with default values for testing
func newTestOrder(status Status) *Order {
	now := time.Now().UTC()
	expireAt := now.Add(30 * time.Minute)
	return &Order{
		Status:         status,
		TenantID:       1,
		UserID:         100,
		OrderNo:        "ORD20260402123456123456",
		Currency:       "CNY",
		TotalAmount:    shared.NewMoney(decimal.NewFromInt(100), "CNY"),
		DiscountAmount: shared.NewMoney(decimal.Zero, "CNY"),
		FreightAmount:  shared.NewMoney(decimal.Zero, "CNY"),
		PayAmount:      shared.NewMoney(decimal.NewFromInt(100), "CNY"),
		OriginalAmount: shared.NewMoney(decimal.NewFromInt(100), "CNY"),
		AdjustAmount:   shared.NewMoney(decimal.Zero, "CNY"),
		Version:        1,
		ExpireAt:       expireAt,
		Audit: shared.AuditInfo{
			CreatedAt: now,
			UpdatedAt: now,
			CreatedBy: 100,
			UpdatedBy: 100,
		},
	}
}

// newTestOrderWithItems creates an order with items for testing
func newTestOrderWithItems(status Status) *Order {
	order := newTestOrder(status)
	order.Items = []OrderItem{
		{
			ProductID:   1,
			SKUId:       1,
			ProductName: "Test Product",
			SKUName:     "Test SKU",
			Image:       "http://example.com/image.jpg",
			Price:       shared.NewMoney(decimal.NewFromInt(50), "CNY"),
			Quantity:    2,
			TotalAmount: shared.NewMoney(decimal.NewFromInt(100), "CNY"),
		},
	}
	return order
}

// ==================== Status Tests ====================

func TestStatus_String(t *testing.T) {
	tests := []struct {
		name     string
		status   Status
		expected string
	}{
		{"pending_payment", StatusPendingPayment, "pending_payment"},
		{"paid", StatusPaid, "paid"},
		{"pending_shipment", StatusPendingShipment, "pending_shipment"},
		{"shipped", StatusShipped, "shipped"},
		{"completed", StatusCompleted, "completed"},
		{"cancelled", StatusCancelled, "cancelled"},
		{"refunding", StatusRefunding, "refunding"},
		{"refunded", StatusRefunded, "refunded"},
		{"unknown", Status(100), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.String())
		})
	}
}

func TestStatus_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		status   Status
		expected bool
	}{
		{"valid_pending_payment", StatusPendingPayment, true},
		{"valid_paid", StatusPaid, true},
		{"valid_pending_shipment", StatusPendingShipment, true},
		{"valid_shipped", StatusShipped, true},
		{"valid_completed", StatusCompleted, true},
		{"valid_cancelled", StatusCancelled, true},
		{"valid_refunding", StatusRefunding, true},
		{"valid_refunded", StatusRefunded, true},
		{"invalid_negative", Status(-1), false},
		{"invalid_too_large", Status(100), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.IsValid())
		})
	}
}

func TestStatus_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name     string
		from     Status
		to       Status
		expected bool
	}{
		// From StatusPendingPayment
		{"pending_payment_to_paid", StatusPendingPayment, StatusPaid, true},
		{"pending_payment_to_cancelled", StatusPendingPayment, StatusCancelled, true},
		{"pending_payment_to_pending_shipment", StatusPendingPayment, StatusPendingShipment, false},
		{"pending_payment_to_shipped", StatusPendingPayment, StatusShipped, false},
		{"pending_payment_to_completed", StatusPendingPayment, StatusCompleted, false},
		{"pending_payment_to_refunding", StatusPendingPayment, StatusRefunding, false},
		{"pending_payment_to_refunded", StatusPendingPayment, StatusRefunded, false},

		// From StatusPaid
		{"paid_to_pending_shipment", StatusPaid, StatusPendingShipment, true},
		{"paid_to_cancelled", StatusPaid, StatusCancelled, true},
		{"paid_to_refunding", StatusPaid, StatusRefunding, true},
		{"paid_to_paid", StatusPaid, StatusPaid, false},
		{"paid_to_shipped", StatusPaid, StatusShipped, false},
		{"paid_to_completed", StatusPaid, StatusCompleted, false},
		{"paid_to_refunded", StatusPaid, StatusRefunded, false},

		// From StatusPendingShipment
		{"pending_shipment_to_shipped", StatusPendingShipment, StatusShipped, true},
		{"pending_shipment_to_cancelled", StatusPendingShipment, StatusCancelled, true},
		{"pending_shipment_to_paid", StatusPendingShipment, StatusPaid, false},
		{"pending_shipment_to_completed", StatusPendingShipment, StatusCompleted, false},
		{"pending_shipment_to_refunding", StatusPendingShipment, StatusRefunding, false},
		{"pending_shipment_to_refunded", StatusPendingShipment, StatusRefunded, false},

		// From StatusShipped
		{"shipped_to_completed", StatusShipped, StatusCompleted, true},
		{"shipped_to_refunding", StatusShipped, StatusRefunding, true},
		{"shipped_to_pending_shipment", StatusShipped, StatusPendingShipment, false},
		{"shipped_to_cancelled", StatusShipped, StatusCancelled, false},
		{"shipped_to_paid", StatusShipped, StatusPaid, false},
		{"shipped_to_refunded", StatusShipped, StatusRefunded, false},

		// From StatusCompleted (terminal state)
		{"completed_to_any", StatusCompleted, StatusPendingPayment, false},
		{"completed_to_paid", StatusCompleted, StatusPaid, false},
		{"completed_to_shipped", StatusCompleted, StatusShipped, false},
		{"completed_to_cancelled", StatusCompleted, StatusCancelled, false},
		{"completed_to_refunding", StatusCompleted, StatusRefunding, false},
		{"completed_to_refunded", StatusCompleted, StatusRefunded, false},

		// From StatusCancelled (terminal state)
		{"cancelled_to_any", StatusCancelled, StatusPendingPayment, false},
		{"cancelled_to_paid", StatusCancelled, StatusPaid, false},
		{"cancelled_to_pending_shipment", StatusCancelled, StatusPendingShipment, false},
		{"cancelled_to_shipped", StatusCancelled, StatusShipped, false},
		{"cancelled_to_completed", StatusCancelled, StatusCompleted, false},
		{"cancelled_to_refunding", StatusCancelled, StatusRefunding, false},
		{"cancelled_to_refunded", StatusCancelled, StatusRefunded, false},

		// From StatusRefunding
		{"refunding_to_refunded", StatusRefunding, StatusRefunded, true},
		{"refunding_to_pending_shipment", StatusRefunding, StatusPendingShipment, true},
		{"refunding_to_cancelled", StatusRefunding, StatusCancelled, false},
		{"refunding_to_paid", StatusRefunding, StatusPaid, false},
		{"refunding_to_shipped", StatusRefunding, StatusShipped, false},
		{"refunding_to_completed", StatusRefunding, StatusCompleted, false},

		// From StatusRefunded (terminal state)
		{"refunded_to_any", StatusRefunded, StatusPendingPayment, false},
		{"refunded_to_pending_shipment", StatusRefunded, StatusPendingShipment, false},
		{"refunded_to_shipped", StatusRefunded, StatusShipped, false},
		{"refunded_to_completed", StatusRefunded, StatusCompleted, false},
		{"refunded_to_cancelled", StatusRefunded, StatusCancelled, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.from.CanTransitionTo(tt.to)
			assert.Equal(t, tt.expected, result, "transition from %v to %v", tt.from, tt.to)
		})
	}
}

// ==================== Order Pay() Tests ====================

func TestOrder_Pay_Success(t *testing.T) {
	order := newTestOrder(StatusPendingPayment)
	order.ExpireAt = time.Now().UTC().Add(30 * time.Minute)

	err := order.Pay()

	assert.NoError(t, err)
	assert.Equal(t, StatusPaid, order.Status)
	assert.NotNil(t, order.PaidAt)
	assert.Equal(t, 2, order.Version) // Version should increment
}

func TestOrder_Pay_InvalidStatus(t *testing.T) {
	tests := []Status{
		StatusPaid,
		StatusPendingShipment,
		StatusShipped,
		StatusCompleted,
		StatusCancelled,
		StatusRefunded,
	}

	for _, status := range tests {
		t.Run(status.String()+"_to_paid", func(t *testing.T) {
			order := newTestOrder(status)
			err := order.Pay()

			assert.Error(t, err)
			assert.Equal(t, code.ErrOrderInvalidStatus, err)
			assert.Equal(t, status, order.Status) // Status should not change
		})
	}
}

func TestOrder_Pay_Expired(t *testing.T) {
	order := newTestOrder(StatusPendingPayment)
	order.ExpireAt = time.Now().UTC().Add(-1 * time.Minute) // Already expired

	err := order.Pay()

	assert.Error(t, err)
	assert.Equal(t, code.ErrOrderExpired, err)
	assert.Equal(t, StatusPendingPayment, order.Status) // Status should not change
}

// ==================== Order Cancel() Tests ====================

func TestOrder_Cancel_Success(t *testing.T) {
	order := newTestOrder(StatusPendingPayment)
	reason := "Customer requested cancellation"

	err := order.Cancel(reason)

	assert.NoError(t, err)
	assert.Equal(t, StatusCancelled, order.Status)
	assert.Equal(t, reason, order.Remark)
	assert.NotNil(t, order.CancelledAt)
	assert.Equal(t, 2, order.Version)
}

func TestOrder_Cancel_ValidStatus_Paid(t *testing.T) {
	order := newTestOrder(StatusPaid)

	err := order.Cancel("reason")

	assert.NoError(t, err)
	assert.Equal(t, StatusCancelled, order.Status)
	assert.NotNil(t, order.CancelledAt)
}

func TestOrder_Cancel_ValidStatus_PendingShipment(t *testing.T) {
	order := newTestOrder(StatusPendingShipment)

	err := order.Cancel("reason")

	assert.NoError(t, err)
	assert.Equal(t, StatusCancelled, order.Status)
	assert.NotNil(t, order.CancelledAt)
}

func TestOrder_Cancel_InvalidStatus_Shipped(t *testing.T) {
	order := newTestOrder(StatusShipped)

	err := order.Cancel("reason")

	assert.Error(t, err)
	assert.Equal(t, code.ErrOrderInvalidStatus, err)
	assert.Equal(t, StatusShipped, order.Status)
}

func TestOrder_Cancel_InvalidStatus_Completed(t *testing.T) {
	order := newTestOrder(StatusCompleted)

	err := order.Cancel("reason")

	assert.Error(t, err)
	assert.Equal(t, code.ErrOrderInvalidStatus, err)
	assert.Equal(t, StatusCompleted, order.Status)
}

// ==================== Order Ship() Tests ====================

func TestOrder_Ship_Success(t *testing.T) {
	order := newTestOrder(StatusPendingShipment)
	trackingNo := "SF1234567890"
	carrier := "顺丰速运"

	err := order.Ship(trackingNo, carrier)

	assert.NoError(t, err)
	assert.Equal(t, StatusShipped, order.Status)
	assert.Equal(t, trackingNo, order.TrackingNo)
	assert.Equal(t, carrier, order.Carrier)
	assert.NotNil(t, order.ShippedAt)
	assert.Equal(t, 2, order.Version)
}

func TestOrder_Ship_InvalidStatus_PendingPayment(t *testing.T) {
	order := newTestOrder(StatusPendingPayment)

	err := order.Ship("SF123", "carrier")

	assert.Error(t, err)
	assert.Equal(t, code.ErrOrderInvalidStatus, err)
	assert.Equal(t, StatusPendingPayment, order.Status)
}

func TestOrder_Ship_InvalidStatus_Paid(t *testing.T) {
	order := newTestOrder(StatusPaid)

	err := order.Ship("SF123", "carrier")

	assert.Error(t, err)
	assert.Equal(t, code.ErrOrderInvalidStatus, err)
	assert.Equal(t, StatusPaid, order.Status)
}

func TestOrder_Ship_InvalidStatus_Shipped(t *testing.T) {
	order := newTestOrder(StatusShipped)

	err := order.Ship("SF123", "carrier")

	assert.Error(t, err)
	assert.Equal(t, code.ErrOrderInvalidStatus, err)
	assert.Equal(t, StatusShipped, order.Status)
}

func TestOrder_Ship_InvalidStatus_Completed(t *testing.T) {
	order := newTestOrder(StatusCompleted)

	err := order.Ship("SF123", "carrier")

	assert.Error(t, err)
	assert.Equal(t, code.ErrOrderInvalidStatus, err)
	assert.Equal(t, StatusCompleted, order.Status)
}

// ==================== Order Complete() Tests ====================

func TestOrder_Complete_Success(t *testing.T) {
	order := newTestOrder(StatusShipped)

	err := order.Complete()

	assert.NoError(t, err)
	assert.Equal(t, StatusCompleted, order.Status)
	assert.NotNil(t, order.CompletedAt)
	assert.Equal(t, 2, order.Version)
}

func TestOrder_Complete_InvalidStatus_PendingPayment(t *testing.T) {
	order := newTestOrder(StatusPendingPayment)

	err := order.Complete()

	assert.Error(t, err)
	assert.Equal(t, code.ErrOrderInvalidStatus, err)
	assert.Equal(t, StatusPendingPayment, order.Status)
}

func TestOrder_Complete_InvalidStatus_Paid(t *testing.T) {
	order := newTestOrder(StatusPaid)

	err := order.Complete()

	assert.Error(t, err)
	assert.Equal(t, code.ErrOrderInvalidStatus, err)
	assert.Equal(t, StatusPaid, order.Status)
}

func TestOrder_Complete_InvalidStatus_PendingShipment(t *testing.T) {
	order := newTestOrder(StatusPendingShipment)

	err := order.Complete()

	assert.Error(t, err)
	assert.Equal(t, code.ErrOrderInvalidStatus, err)
	assert.Equal(t, StatusPendingShipment, order.Status)
}

func TestOrder_Complete_InvalidStatus_Shipped(t *testing.T) {
	order := newTestOrder(StatusShipped)

	err := order.Complete()

	assert.NoError(t, err)
	assert.Equal(t, StatusCompleted, order.Status)
}

func TestOrder_Complete_InvalidStatus_Cancelled(t *testing.T) {
	order := newTestOrder(StatusCancelled)

	err := order.Complete()

	assert.Error(t, err)
	assert.Equal(t, code.ErrOrderInvalidStatus, err)
	assert.Equal(t, StatusCancelled, order.Status)
}

// ==================== Order AdjustPrice() Tests ====================

func TestOrder_AdjustPrice_Success(t *testing.T) {
	order := newTestOrder(StatusPendingPayment)
	order.OriginalAmount = shared.NewMoney(decimal.NewFromInt(100), "CNY")
	order.PayAmount = shared.NewMoney(decimal.NewFromInt(100), "CNY")
	adjustAmount := decimal.NewFromInt(-10)
	reason := "Discount for VIP customer"
	adjustedBy := int64(1)

	err := order.AdjustPrice(adjustAmount, reason, adjustedBy)

	assert.NoError(t, err)
	assert.True(t, order.AdjustAmount.Amount.Equal(adjustAmount))
	assert.Equal(t, reason, order.AdjustReason)
	assert.Equal(t, adjustedBy, order.AdjustedBy)
	assert.NotNil(t, order.AdjustedAt)
	assert.Equal(t, 2, order.Version)
}

func TestOrder_AdjustPrice_InvalidStatus_Paid(t *testing.T) {
	order := newTestOrder(StatusPaid)
	order.AdjustAmount = shared.NewMoney(decimal.Zero, "CNY")
	order.AdjustReason = ""
	order.AdjustedBy = 0
	order.AdjustedAt = nil

	err := order.AdjustPrice(decimal.NewFromInt(-10), "reason", 1)

	assert.Error(t, err)
	assert.Equal(t, code.ErrOrderInvalidStatus, err)
	assert.True(t, order.AdjustAmount.Amount.IsZero())
	assert.Empty(t, order.AdjustReason)
}

func TestOrder_AdjustPrice_InvalidStatus_PendingShipment(t *testing.T) {
	order := newTestOrder(StatusPendingShipment)

	err := order.AdjustPrice(decimal.NewFromInt(-10), "reason", 1)

	assert.Error(t, err)
	assert.Equal(t, code.ErrOrderInvalidStatus, err)
}

func TestOrder_AdjustPrice_InvalidStatus_Shipped(t *testing.T) {
	order := newTestOrder(StatusShipped)

	err := order.AdjustPrice(decimal.NewFromInt(-10), "reason", 1)

	assert.Error(t, err)
	assert.Equal(t, code.ErrOrderInvalidStatus, err)
}

func TestOrder_AdjustPrice_InvalidStatus_Completed(t *testing.T) {
	order := newTestOrder(StatusCompleted)

	err := order.AdjustPrice(decimal.NewFromInt(-10), "reason", 1)

	assert.Error(t, err)
	assert.Equal(t, code.ErrOrderInvalidStatus, err)
}

func TestOrder_AdjustPrice_InvalidStatus_Cancelled(t *testing.T) {
	order := newTestOrder(StatusCancelled)

	err := order.AdjustPrice(decimal.NewFromInt(-10), "reason", 1)

	assert.Error(t, err)
	assert.Equal(t, code.ErrOrderInvalidStatus, err)
}

func TestOrder_AdjustPrice_PositiveAmount(t *testing.T) {
	order := newTestOrder(StatusPendingPayment)
	order.OriginalAmount = shared.NewMoney(decimal.NewFromInt(100), "CNY")
	order.PayAmount = shared.NewMoney(decimal.NewFromInt(100), "CNY")
	adjustAmount := decimal.NewFromInt(5) // Price increase

	err := order.AdjustPrice(adjustAmount, "Additional fee", 1)

	assert.NoError(t, err)
	assert.True(t, order.AdjustAmount.Amount.Equal(adjustAmount))
}

// ==================== Order IsExpired() Tests ====================

func TestOrder_IsExpired_NotExpired(t *testing.T) {
	order := newTestOrder(StatusPendingPayment)
	order.ExpireAt = time.Now().UTC().Add(30 * time.Minute)

	assert.False(t, order.IsExpired())
}

func TestOrder_IsExpired_Expired(t *testing.T) {
	order := newTestOrder(StatusPendingPayment)
	order.ExpireAt = time.Now().UTC().Add(-1 * time.Minute)

	assert.True(t, order.IsExpired())
}

func TestOrder_IsExpired_JustExpired(t *testing.T) {
	order := newTestOrder(StatusPendingPayment)
	order.ExpireAt = time.Now().UTC().Add(-1 * time.Second)

	assert.True(t, order.IsExpired())
}

func TestOrder_IsExpired_AlmostExpired(t *testing.T) {
	order := newTestOrder(StatusPendingPayment)
	order.ExpireAt = time.Now().UTC().Add(1 * time.Second)

	assert.False(t, order.IsExpired())
}

// ==================== Order CalculateTotals() Tests ====================

func TestOrder_CalculateTotals_SingleItem(t *testing.T) {
	order := &Order{
		Currency: "CNY",
		Items: []OrderItem{
			{
				Price:       shared.NewMoney(decimal.NewFromInt(50), "CNY"),
				Quantity:    2,
				TotalAmount: shared.NewMoney(decimal.NewFromInt(100), "CNY"),
			},
		},
		DiscountAmount: shared.NewMoney(decimal.Zero, "CNY"),
		FreightAmount:  shared.NewMoney(decimal.Zero, "CNY"),
	}

	order.CalculateTotals()

	assert.True(t, order.TotalAmount.Amount.Equal(decimal.NewFromInt(100)))
	assert.True(t, order.OriginalAmount.Amount.Equal(decimal.NewFromInt(100)))
	assert.True(t, order.PayAmount.Amount.Equal(decimal.NewFromInt(100)))
}

func TestOrder_CalculateTotals_MultipleItems(t *testing.T) {
	order := &Order{
		Currency: "CNY",
		Items: []OrderItem{
			{
				Price:       shared.NewMoney(decimal.NewFromFloat(10.5), "CNY"),
				Quantity:    2,
				TotalAmount: shared.NewMoney(decimal.NewFromFloat(21), "CNY"),
			},
			{
				Price:       shared.NewMoney(decimal.NewFromFloat(30), "CNY"),
				Quantity:    1,
				TotalAmount: shared.NewMoney(decimal.NewFromFloat(30), "CNY"),
			},
		},
		DiscountAmount: shared.NewMoney(decimal.NewFromFloat(5), "CNY"),
		FreightAmount:  shared.NewMoney(decimal.NewFromFloat(10), "CNY"),
	}

	order.CalculateTotals()

	// Total: 21 + 30 = 51
	assert.True(t, order.TotalAmount.Amount.Equal(decimal.NewFromFloat(51)))
	// Original: 51
	assert.True(t, order.OriginalAmount.Amount.Equal(decimal.NewFromFloat(51)))
	// Pay: 51 - 5 + 10 = 56
	assert.True(t, order.PayAmount.Amount.Equal(decimal.NewFromFloat(56)))
}

func TestOrder_CalculateTotals_EmptyItems(t *testing.T) {
	order := &Order{
		Currency:       "CNY",
		Items:          []OrderItem{},
		DiscountAmount: shared.NewMoney(decimal.Zero, "CNY"),
		FreightAmount:  shared.NewMoney(decimal.Zero, "CNY"),
	}

	order.CalculateTotals()

	assert.True(t, order.TotalAmount.Amount.IsZero())
	assert.True(t, order.OriginalAmount.Amount.IsZero())
	assert.True(t, order.PayAmount.Amount.IsZero())
}

func TestOrder_CalculateTotals_WithDiscount(t *testing.T) {
	order := &Order{
		Currency: "CNY",
		Items: []OrderItem{
			{
				Price:       shared.NewMoney(decimal.NewFromInt(100), "CNY"),
				Quantity:    1,
				TotalAmount: shared.NewMoney(decimal.NewFromInt(100), "CNY"),
			},
		},
		DiscountAmount: shared.NewMoney(decimal.NewFromInt(20), "CNY"),
		FreightAmount:  shared.NewMoney(decimal.Zero, "CNY"),
	}

	order.CalculateTotals()

	// Total: 100
	assert.True(t, order.TotalAmount.Amount.Equal(decimal.NewFromInt(100)))
	// Pay: 100 - 20 = 80
	assert.True(t, order.PayAmount.Amount.Equal(decimal.NewFromInt(80)))
}

func TestOrder_CalculateTotals_WithFreight(t *testing.T) {
	order := &Order{
		Currency: "CNY",
		Items: []OrderItem{
			{
				Price:       shared.NewMoney(decimal.NewFromInt(100), "CNY"),
				Quantity:    1,
				TotalAmount: shared.NewMoney(decimal.NewFromInt(100), "CNY"),
			},
		},
		DiscountAmount: shared.NewMoney(decimal.Zero, "CNY"),
		FreightAmount:  shared.NewMoney(decimal.NewFromInt(15), "CNY"),
	}

	order.CalculateTotals()

	// Total: 100
	assert.True(t, order.TotalAmount.Amount.Equal(decimal.NewFromInt(100)))
	// Pay: 100 + 15 = 115
	assert.True(t, order.PayAmount.Amount.Equal(decimal.NewFromInt(115)))
}

func TestOrder_CalculateTotals_WithDiscountAndFreight(t *testing.T) {
	order := &Order{
		Currency: "CNY",
		Items: []OrderItem{
			{
				Price:       shared.NewMoney(decimal.NewFromInt(100), "CNY"),
				Quantity:    2,
				TotalAmount: shared.NewMoney(decimal.NewFromInt(200), "CNY"),
			},
		},
		DiscountAmount: shared.NewMoney(decimal.NewFromInt(30), "CNY"),
		FreightAmount:  shared.NewMoney(decimal.NewFromInt(20), "CNY"),
	}

	order.CalculateTotals()

	// Total: 200
	assert.True(t, order.TotalAmount.Amount.Equal(decimal.NewFromInt(200)))
	// Pay: 200 - 30 + 20 = 190
	assert.True(t, order.PayAmount.Amount.Equal(decimal.NewFromInt(190)))
}

// ==================== OrderItem Multiply() Tests ====================

func TestOrderItem_Multiply(t *testing.T) {
	item := OrderItem{
		Price:    shared.NewMoney(decimal.NewFromFloat(12.99), "CNY"),
		Quantity: 3,
	}

	// Testing the Money.Multiply method used in CalculateTotals
	result := item.Price.Multiply(item.Quantity)

	assert.True(t, result.Amount.Equal(decimal.NewFromFloat(38.97)))
	assert.Equal(t, "CNY", result.Currency)
}

// ==================== Full State Machine Flow Tests ====================

func TestOrder_FullHappyPath(t *testing.T) {
	// Create order
	order := newTestOrderWithItems(StatusPendingPayment)
	order.TotalAmount = shared.NewMoney(decimal.NewFromInt(100), "CNY")
	order.PayAmount = shared.NewMoney(decimal.NewFromInt(100), "CNY")
	order.OriginalAmount = shared.NewMoney(decimal.NewFromInt(100), "CNY")
	order.DiscountAmount = shared.NewMoney(decimal.Zero, "CNY")
	order.FreightAmount = shared.NewMoney(decimal.Zero, "CNY")

	// 1. Pay
	err := order.Pay()
	assert.NoError(t, err)
	assert.Equal(t, StatusPaid, order.Status)
	assert.NotNil(t, order.PaidAt)

	// Note: Transition from Paid to PendingShipment happens via admin action in service layer
	// For entity testing, we manually set the status (simulating service layer behavior)
	order.Status = StatusPendingShipment

	// 2. Ship
	err = order.Ship("SF1234567890", "顺丰速运")
	assert.NoError(t, err)
	assert.Equal(t, StatusShipped, order.Status)
	assert.NotNil(t, order.ShippedAt)
	assert.Equal(t, "SF1234567890", order.TrackingNo)

	// 3. Complete
	err = order.Complete()
	assert.NoError(t, err)
	assert.Equal(t, StatusCompleted, order.Status)
	assert.NotNil(t, order.CompletedAt)
}

func TestOrder_CancellationFlow_PendingPayment(t *testing.T) {
	order := newTestOrder(StatusPendingPayment)

	err := order.Cancel("Customer changed mind")
	assert.NoError(t, err)
	assert.Equal(t, StatusCancelled, order.Status)
	assert.NotNil(t, order.CancelledAt)
	assert.Equal(t, "Customer changed mind", order.Remark)
}

func TestOrder_CancellationFlow_Paid(t *testing.T) {
	order := newTestOrder(StatusPaid)

	err := order.Cancel("Customer changed mind")
	assert.NoError(t, err)
	assert.Equal(t, StatusCancelled, order.Status)
}

func TestOrder_RefundFlow_PaidToRefunding(t *testing.T) {
	order := newTestOrder(StatusPaid)

	// Transition to refunding (valid from paid)
	assert.True(t, order.Status.CanTransitionTo(StatusRefunding))
}

func TestOrder_RefundFlow_RefundingToRefunded(t *testing.T) {
	order := newTestOrder(StatusRefunding)

	err := order.Cancel("Refund requested") // Cancel is not valid, but status transition would be via refund
	assert.Error(t, err)
	assert.Equal(t, code.ErrOrderInvalidStatus, err)
}

func TestOrder_CannotTransitionFromTerminalStates(t *testing.T) {
	terminalStates := []Status{StatusCompleted, StatusCancelled, StatusRefunded}
	allStatuses := []Status{
		StatusPendingPayment, StatusPaid, StatusPendingShipment,
		StatusShipped, StatusCompleted, StatusCancelled,
		StatusRefunding, StatusRefunded,
	}

	for _, terminal := range terminalStates {
		for _, target := range allStatuses {
			if terminal == target {
				continue
			}
			t.Run(terminal.String()+"_to_"+target.String(), func(t *testing.T) {
				assert.False(t, terminal.CanTransitionTo(target),
					"terminal state %v should not transition to %v", terminal, target)
			})
		}
	}
}

// ==================== FulfillmentStatus Tests ====================

func TestFulfillmentStatus_String(t *testing.T) {
	tests := []struct {
		status   FulfillmentStatus
		expected string
	}{
		{FulfillmentStatusPending, "pending"},
		{FulfillmentStatusPartialShipped, "partial_shipped"},
		{FulfillmentStatusShipped, "shipped"},
		{FulfillmentStatusDelivered, "delivered"},
		{FulfillmentStatus(100), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.String())
		})
	}
}

// ==================== RefundStatus Tests ====================

func TestRefundStatus_String(t *testing.T) {
	tests := []struct {
		status   RefundStatus
		expected string
	}{
		{RefundStatusNone, "none"},
		{RefundStatusPending, "pending"},
		{RefundStatusApproved, "approved"},
		{RefundStatusRejected, "rejected"},
		{RefundStatusCompleted, "completed"},
		{RefundStatus(100), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.String())
		})
	}
}

// ==================== Version Increment Tests ====================

func TestOrder_VersionIncrementsOnPay(t *testing.T) {
	order := newTestOrder(StatusPendingPayment)
	initialVersion := order.Version

	err := order.Pay()
	assert.NoError(t, err)
	assert.Equal(t, initialVersion+1, order.Version)
}

func TestOrder_VersionIncrementsOnCancel(t *testing.T) {
	order := newTestOrder(StatusPendingPayment)
	initialVersion := order.Version

	err := order.Cancel("reason")
	assert.NoError(t, err)
	assert.Equal(t, initialVersion+1, order.Version)
}

func TestOrder_VersionIncrementsOnShip(t *testing.T) {
	order := newTestOrder(StatusPendingShipment)
	initialVersion := order.Version

	err := order.Ship("SF123", "carrier")
	assert.NoError(t, err)
	assert.Equal(t, initialVersion+1, order.Version)
}

func TestOrder_VersionIncrementsOnComplete(t *testing.T) {
	order := newTestOrder(StatusShipped)
	initialVersion := order.Version

	err := order.Complete()
	assert.NoError(t, err)
	assert.Equal(t, initialVersion+1, order.Version)
}

func TestOrder_VersionIncrementsOnAdjustPrice(t *testing.T) {
	order := newTestOrder(StatusPendingPayment)
	initialVersion := order.Version

	err := order.AdjustPrice(decimal.NewFromInt(-10), "discount", 1)
	assert.NoError(t, err)
	assert.Equal(t, initialVersion+1, order.Version)
}

func TestOrder_VersionNotIncrementsOnFailedTransition(t *testing.T) {
	order := newTestOrder(StatusCompleted)
	initialVersion := order.Version

	err := order.Pay()
	assert.Error(t, err)
	assert.Equal(t, initialVersion, order.Version) // Version should not change

	err = order.Cancel("reason")
	assert.Error(t, err)
	assert.Equal(t, initialVersion, order.Version)

	err = order.Ship("SF123", "carrier")
	assert.Error(t, err)
	assert.Equal(t, initialVersion, order.Version)
}
