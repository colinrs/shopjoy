package payment

import (
	"context"
	"errors"
	"testing"

	"github.com/colinrs/shopjoy/admin/internal/domain/payment"
	"github.com/colinrs/shopjoy/admin/internal/domain/shop"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// =============================================================================
// Mock Implementations
// =============================================================================

// mockPaymentRepo is a mock implementation of payment.PaymentRepository
type mockPaymentRepo struct {
	findByChannelPaymentIDFunc func(ctx context.Context, db *gorm.DB, channelPaymentID string) (*payment.Payment, error)
	updateFunc                 func(ctx context.Context, db *gorm.DB, payment *payment.Payment) error
}

func (m *mockPaymentRepo) Create(ctx context.Context, db *gorm.DB, p *payment.Payment) error {
	return nil
}

func (m *mockPaymentRepo) Update(ctx context.Context, db *gorm.DB, p *payment.Payment) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, db, p)
	}
	return nil
}

func (m *mockPaymentRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*payment.Payment, error) {
	return nil, nil
}

func (m *mockPaymentRepo) FindByOrderID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderID int64) (*payment.Payment, error) {
	return nil, nil
}

func (m *mockPaymentRepo) FindByPaymentNo(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, paymentNo string) (*payment.Payment, error) {
	return nil, nil
}

func (m *mockPaymentRepo) FindByChannelPaymentID(ctx context.Context, db *gorm.DB, channelPaymentID string) (*payment.Payment, error) {
	if m.findByChannelPaymentIDFunc != nil {
		return m.findByChannelPaymentIDFunc(ctx, db, channelPaymentID)
	}
	return nil, nil
}

func (m *mockPaymentRepo) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query payment.PaymentQuery) ([]*payment.Payment, int64, error) {
	return nil, 0, nil
}

// mockRefundRepo is a mock implementation of payment.PaymentRefundRepository
type mockRefundRepo struct{}

func (m *mockRefundRepo) Create(ctx context.Context, db *gorm.DB, r *payment.PaymentRefund) error {
	return nil
}

func (m *mockRefundRepo) Update(ctx context.Context, db *gorm.DB, r *payment.PaymentRefund) error {
	return nil
}

func (m *mockRefundRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*payment.PaymentRefund, error) {
	return nil, nil
}

func (m *mockRefundRepo) FindByRefundNo(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, refundNo string) (*payment.PaymentRefund, error) {
	return nil, nil
}

func (m *mockRefundRepo) FindByIdempotencyKey(ctx context.Context, db *gorm.DB, idempotencyKey string) (*payment.PaymentRefund, error) {
	return nil, nil
}

func (m *mockRefundRepo) FindByPaymentID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, paymentID int64) ([]*payment.PaymentRefund, error) {
	return nil, nil
}

func (m *mockRefundRepo) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query payment.RefundQuery) ([]*payment.PaymentRefund, int64, error) {
	return nil, 0, nil
}

func (m *mockRefundRepo) GetTotalRefundedAmount(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, paymentID int64) (decimal.Decimal, error) {
	return decimal.Zero, nil
}

// mockTransactionRepo is a mock implementation of payment.PaymentTransactionRepository
type mockTransactionRepo struct{}

func (m *mockTransactionRepo) Create(ctx context.Context, db *gorm.DB, t *payment.PaymentTransaction) error {
	return nil
}

func (m *mockTransactionRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*payment.PaymentTransaction, error) {
	return nil, nil
}

func (m *mockTransactionRepo) FindByTransactionID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, txnID string) (*payment.PaymentTransaction, error) {
	return nil, nil
}

func (m *mockTransactionRepo) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query payment.TransactionQuery) ([]*payment.PaymentTransaction, int64, error) {
	return nil, 0, nil
}

func (m *mockTransactionRepo) GetStats(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (int64, int64, int64, error) {
	return 0, 0, 0, nil
}

// mockWebhookEventRepo is a mock implementation of payment.WebhookEventRepository
type mockWebhookEventRepo struct {
	findByEventIDFunc func(ctx context.Context, db *gorm.DB, eventID string) (*payment.WebhookEvent, error)
	createFunc        func(ctx context.Context, db *gorm.DB, event *payment.WebhookEvent) error
	markProcessedFunc func(ctx context.Context, db *gorm.DB, id int64, processed int8, errorMsg string) error
}

func (m *mockWebhookEventRepo) Create(ctx context.Context, db *gorm.DB, event *payment.WebhookEvent) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, db, event)
	}
	return nil
}

func (m *mockWebhookEventRepo) FindByEventID(ctx context.Context, db *gorm.DB, eventID string) (*payment.WebhookEvent, error) {
	if m.findByEventIDFunc != nil {
		return m.findByEventIDFunc(ctx, db, eventID)
	}
	return nil, nil
}

func (m *mockWebhookEventRepo) MarkProcessed(ctx context.Context, db *gorm.DB, id int64, processed int8, errorMsg string) error {
	if m.markProcessedFunc != nil {
		return m.markProcessedFunc(ctx, db, id, processed, errorMsg)
	}
	return nil
}

// mockPaymentSettingsRepo is a mock implementation of shop.PaymentSettingsRepository
type mockPaymentSettingsRepo struct{}

func (m *mockPaymentSettingsRepo) Create(ctx context.Context, db *gorm.DB, settings *shop.PaymentSettings) error {
	return nil
}

func (m *mockPaymentSettingsRepo) Update(ctx context.Context, db *gorm.DB, settings *shop.PaymentSettings) error {
	return nil
}

func (m *mockPaymentSettingsRepo) Save(ctx context.Context, db *gorm.DB, settings *shop.PaymentSettings) error {
	return nil
}

func (m *mockPaymentSettingsRepo) FindByShopID(ctx context.Context, db *gorm.DB, shopID int64) (*shop.PaymentSettings, error) {
	return nil, nil
}

// mockIDGen is a mock implementation of snowflake.Snowflake
type mockIDGen struct {
	nextIDFunc func(ctx context.Context) (int64, error)
}

func (m *mockIDGen) GetNodeID() int64 {
	return 1
}

func (m *mockIDGen) NextID(ctx context.Context) (int64, error) {
	if m.nextIDFunc != nil {
		return m.nextIDFunc(ctx)
	}
	return 1, nil
}

func (m *mockIDGen) NextIDs(ctx context.Context, n int64) ([]int64, error) {
	return nil, nil
}

// =============================================================================
// Test Cases
// =============================================================================

func TestHandleWebhook_PaymentIntentSucceeded(t *testing.T) {
	tests := []struct {
		name                 string
		event                *WebhookEvent
		setupMockPaymentRepo func() *mockPaymentRepo
		setupMockWebhookRepo func() *mockWebhookEventRepo
		wantErr              bool
		errMsg               string
		paymentStatusCheck   func(t *testing.T, payment *payment.Payment)
	}{
		{
			name: "successful payment intent update",
			event: &WebhookEvent{
				EventID:    "evt_123",
				EventType:  "payment_intent.succeeded",
				ResourceID: "pi_abc123",
				RawPayload: `{"id": "pi_abc123", "amount": 9999}`,
			},
			setupMockPaymentRepo: func() *mockPaymentRepo {
				return &mockPaymentRepo{
					findByChannelPaymentIDFunc: func(ctx context.Context, db *gorm.DB, channelPaymentID string) (*payment.Payment, error) {
						if channelPaymentID != "pi_abc123" {
							t.Errorf("expected channelPaymentID=pi_abc123, got %s", channelPaymentID)
						}
						return &payment.Payment{
							Model:           application.Model{ID: 1},
							TenantID:        1,
							OrderID:         100,
							PaymentNo:       "PAY123",
							Amount:          decimal.NewFromInt(100),
							Currency:        "USD",
							Status:          payment.PaymentStatusPending,
							ChannelIntentID: "pi_abc123",
						}, nil
					},
					updateFunc: func(ctx context.Context, db *gorm.DB, p *payment.Payment) error {
						return nil
					},
				}
			},
			setupMockWebhookRepo: func() *mockWebhookEventRepo {
				return &mockWebhookEventRepo{
					findByEventIDFunc: func(ctx context.Context, db *gorm.DB, eventID string) (*payment.WebhookEvent, error) {
						return nil, nil // Event not found, will be created
					},
					createFunc: func(ctx context.Context, db *gorm.DB, event *payment.WebhookEvent) error {
						if event.EventID != "evt_123" {
							t.Errorf("expected EventID=evt_123, got %s", event.EventID)
						}
						if event.EventType != "payment_intent.succeeded" {
							t.Errorf("expected EventType=payment_intent.succeeded, got %s", event.EventType)
						}
						return nil
					},
					markProcessedFunc: func(ctx context.Context, db *gorm.DB, id int64, processed int8, errorMsg string) error {
						if processed != 1 {
							t.Errorf("expected processed=1, got %d", processed)
						}
						if errorMsg != "" {
							t.Errorf("expected empty errorMsg, got %s", errorMsg)
						}
						return nil
					},
				}
			},
			wantErr:            false,
			paymentStatusCheck: func(t *testing.T, p *payment.Payment) {},
		},
		{
			name: "payment not found returns error",
			event: &WebhookEvent{
				EventID:    "evt_124",
				EventType:  "payment_intent.succeeded",
				ResourceID: "pi_notfound",
				RawPayload: `{"id": "pi_notfound"}`,
			},
			setupMockPaymentRepo: func() *mockPaymentRepo {
				return &mockPaymentRepo{
					findByChannelPaymentIDFunc: func(ctx context.Context, db *gorm.DB, channelPaymentID string) (*payment.Payment, error) {
						return nil, errors.New("payment not found")
					},
				}
			},
			setupMockWebhookRepo: func() *mockWebhookEventRepo {
				return &mockWebhookEventRepo{
					findByEventIDFunc: func(ctx context.Context, db *gorm.DB, eventID string) (*payment.WebhookEvent, error) {
						return nil, nil
					},
					createFunc: func(ctx context.Context, db *gorm.DB, event *payment.WebhookEvent) error {
						return nil
					},
					markProcessedFunc: func(ctx context.Context, db *gorm.DB, id int64, processed int8, errorMsg string) error {
						// Should be marked as failed
						if processed != 2 {
							t.Errorf("expected processed=2 for failure, got %d", processed)
						}
						return nil
					},
				}
			},
			wantErr: true,
			errMsg:  "payment not found",
		},
		{
			name: "webhook event already processed - idempotent behavior",
			event: &WebhookEvent{
				EventID:    "evt_125",
				EventType:  "payment_intent.succeeded",
				ResourceID: "pi_alreadyprocessed",
				RawPayload: `{"id": "pi_alreadyprocessed"}`,
			},
			setupMockPaymentRepo: func() *mockPaymentRepo {
				return &mockPaymentRepo{}
			},
			setupMockWebhookRepo: func() *mockWebhookEventRepo {
				return &mockWebhookEventRepo{
					findByEventIDFunc: func(ctx context.Context, db *gorm.DB, eventID string) (*payment.WebhookEvent, error) {
						// Return existing processed event
						return &payment.WebhookEvent{
							Model:      application.Model{ID: 99},
							EventID:    "evt_125",
							EventType:  "payment_intent.succeeded",
							Processed:  1, // Already processed
							RawPayload: `{"id": "pi_alreadyprocessed"}`,
						}, nil
					},
				}
			},
			wantErr: false, // Should return nil (idempotent - skip processing)
		},
		{
			name: "webhook event update fails",
			event: &WebhookEvent{
				EventID:    "evt_126",
				EventType:  "payment_intent.succeeded",
				ResourceID: "pi_updatefail",
				RawPayload: `{"id": "pi_updatefail"}`,
			},
			setupMockPaymentRepo: func() *mockPaymentRepo {
				return &mockPaymentRepo{
					findByChannelPaymentIDFunc: func(ctx context.Context, db *gorm.DB, channelPaymentID string) (*payment.Payment, error) {
						return &payment.Payment{
							Model:           application.Model{ID: 1},
							Status:          payment.PaymentStatusPending,
							ChannelIntentID: "pi_updatefail",
						}, nil
					},
					updateFunc: func(ctx context.Context, db *gorm.DB, p *payment.Payment) error {
						return errors.New("database update failed")
					},
				}
			},
			setupMockWebhookRepo: func() *mockWebhookEventRepo {
				return &mockWebhookEventRepo{
					findByEventIDFunc: func(ctx context.Context, db *gorm.DB, eventID string) (*payment.WebhookEvent, error) {
						return nil, nil
					},
					createFunc: func(ctx context.Context, db *gorm.DB, event *payment.WebhookEvent) error {
						return nil
					},
					markProcessedFunc: func(ctx context.Context, db *gorm.DB, id int64, processed int8, errorMsg string) error {
						if processed != 2 {
							t.Errorf("expected processed=2 for failure, got %d", processed)
						}
						return nil
					},
				}
			},
			wantErr: true,
			errMsg:  "database update failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &service{
				db:                  nil,
				paymentRepo:         tt.setupMockPaymentRepo(),
				refundRepo:          &mockRefundRepo{},
				transactionRepo:     &mockTransactionRepo{},
				webhookEventRepo:    tt.setupMockWebhookRepo(),
				paymentSettingsRepo: &mockPaymentSettingsRepo{},
				idGen:               &mockIDGen{},
			}

			err := svc.HandleWebhook(context.Background(), tt.event)

			if tt.wantErr {
				if err == nil {
					t.Errorf("HandleWebhook() error = nil, want error containing %q", tt.errMsg)
					return
				}
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("HandleWebhook() error = %q, want error containing %q", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("HandleWebhook() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestHandleWebhook_PaymentIntentFailed(t *testing.T) {
	tests := []struct {
		name                 string
		event                *WebhookEvent
		setupMockPaymentRepo func() *mockPaymentRepo
		setupMockWebhookRepo func() *mockWebhookEventRepo
		wantErr              bool
		errMsg               string
	}{
		{
			name: "failed payment intent updates payment status",
			event: &WebhookEvent{
				EventID:    "evt_fail_001",
				EventType:  "payment_intent.payment_failed",
				ResourceID: "pi_failed123",
				RawPayload: `{"id": "pi_failed123", "last_payment_error": {"message": "Card declined"}}`,
			},
			setupMockPaymentRepo: func() *mockPaymentRepo {
				return &mockPaymentRepo{
					findByChannelPaymentIDFunc: func(ctx context.Context, db *gorm.DB, channelPaymentID string) (*payment.Payment, error) {
						return &payment.Payment{
							Model:           application.Model{ID: 2},
							TenantID:        1,
							OrderID:         101,
							PaymentNo:       "PAY124",
							Amount:          decimal.NewFromInt(50),
							Currency:        "USD",
							Status:          payment.PaymentStatusPending,
							ChannelIntentID: "pi_failed123",
						}, nil
					},
					updateFunc: func(ctx context.Context, db *gorm.DB, p *payment.Payment) error {
						// Verify payment was marked as failed
						if p.Status != payment.PaymentStatusFailed {
							t.Errorf("expected Status=PaymentStatusFailed, got %d", p.Status)
						}
						if p.FailedReason != "Payment failed via webhook" {
							t.Errorf("expected FailedReason='Payment failed via webhook', got %q", p.FailedReason)
						}
						return nil
					},
				}
			},
			setupMockWebhookRepo: func() *mockWebhookEventRepo {
				return &mockWebhookEventRepo{
					findByEventIDFunc: func(ctx context.Context, db *gorm.DB, eventID string) (*payment.WebhookEvent, error) {
						return nil, nil
					},
					createFunc: func(ctx context.Context, db *gorm.DB, event *payment.WebhookEvent) error {
						return nil
					},
					markProcessedFunc: func(ctx context.Context, db *gorm.DB, id int64, processed int8, errorMsg string) error {
						if processed != 1 {
							t.Errorf("expected processed=1, got %d", processed)
						}
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "failed payment intent - payment not found",
			event: &WebhookEvent{
				EventID:    "evt_fail_002",
				EventType:  "payment_intent.payment_failed",
				ResourceID: "pi_notexist",
				RawPayload: `{"id": "pi_notexist"}`,
			},
			setupMockPaymentRepo: func() *mockPaymentRepo {
				return &mockPaymentRepo{
					findByChannelPaymentIDFunc: func(ctx context.Context, db *gorm.DB, channelPaymentID string) (*payment.Payment, error) {
						return nil, errors.New("record not found")
					},
				}
			},
			setupMockWebhookRepo: func() *mockWebhookEventRepo {
				return &mockWebhookEventRepo{
					findByEventIDFunc: func(ctx context.Context, db *gorm.DB, eventID string) (*payment.WebhookEvent, error) {
						return nil, nil
					},
					createFunc: func(ctx context.Context, db *gorm.DB, event *payment.WebhookEvent) error {
						return nil
					},
					markProcessedFunc: func(ctx context.Context, db *gorm.DB, id int64, processed int8, errorMsg string) error {
						if processed != 2 {
							t.Errorf("expected processed=2 for failure, got %d", processed)
						}
						return nil
					},
				}
			},
			wantErr: true,
			errMsg:  "record not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &service{
				db:                  nil,
				paymentRepo:         tt.setupMockPaymentRepo(),
				refundRepo:          &mockRefundRepo{},
				transactionRepo:     &mockTransactionRepo{},
				webhookEventRepo:    tt.setupMockWebhookRepo(),
				paymentSettingsRepo: &mockPaymentSettingsRepo{},
				idGen:               &mockIDGen{},
			}

			err := svc.HandleWebhook(context.Background(), tt.event)

			if tt.wantErr {
				if err == nil {
					t.Errorf("HandleWebhook() error = nil, want error containing %q", tt.errMsg)
					return
				}
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("HandleWebhook() error = %q, want error containing %q", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("HandleWebhook() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestHandleWebhook_ChargeRefunded(t *testing.T) {
	tests := []struct {
		name                 string
		event                *WebhookEvent
		setupMockPaymentRepo func() *mockPaymentRepo
		setupMockWebhookRepo func() *mockWebhookEventRepo
		wantErr              bool
		errMsg               string
	}{
		{
			name: "refund webhook marks payment as refunded",
			event: &WebhookEvent{
				EventID:    "evt_refund_001",
				EventType:  "charge.refunded",
				ResourceID: "ch_refunded123",
				RawPayload: `{"id": "ch_refunded123", "refunded": true}`,
			},
			setupMockPaymentRepo: func() *mockPaymentRepo {
				return &mockPaymentRepo{
					findByChannelPaymentIDFunc: func(ctx context.Context, db *gorm.DB, channelPaymentID string) (*payment.Payment, error) {
						return &payment.Payment{
							Model:            application.Model{ID: 3},
							TenantID:         1,
							OrderID:          102,
							PaymentNo:        "PAY125",
							Amount:           decimal.NewFromInt(200),
							Currency:         "USD",
							Status:           payment.PaymentStatusSuccess,
							ChannelPaymentID: "ch_refunded123",
						}, nil
					},
					updateFunc: func(ctx context.Context, db *gorm.DB, p *payment.Payment) error {
						// Verify payment was marked as refunded (partial=true in this case)
						if p.Status != payment.PaymentStatusPartiallyRefunded {
							t.Errorf("expected Status=PaymentStatusPartiallyRefunded, got %d", p.Status)
						}
						return nil
					},
				}
			},
			setupMockWebhookRepo: func() *mockWebhookEventRepo {
				return &mockWebhookEventRepo{
					findByEventIDFunc: func(ctx context.Context, db *gorm.DB, eventID string) (*payment.WebhookEvent, error) {
						return nil, nil
					},
					createFunc: func(ctx context.Context, db *gorm.DB, event *payment.WebhookEvent) error {
						return nil
					},
					markProcessedFunc: func(ctx context.Context, db *gorm.DB, id int64, processed int8, errorMsg string) error {
						if processed != 1 {
							t.Errorf("expected processed=1, got %d", processed)
						}
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "refund webhook - payment not found",
			event: &WebhookEvent{
				EventID:    "evt_refund_002",
				EventType:  "charge.refunded",
				ResourceID: "ch_notfound",
				RawPayload: `{"id": "ch_notfound"}`,
			},
			setupMockPaymentRepo: func() *mockPaymentRepo {
				return &mockPaymentRepo{
					findByChannelPaymentIDFunc: func(ctx context.Context, db *gorm.DB, channelPaymentID string) (*payment.Payment, error) {
						return nil, errors.New("payment not found")
					},
				}
			},
			setupMockWebhookRepo: func() *mockWebhookEventRepo {
				return &mockWebhookEventRepo{
					findByEventIDFunc: func(ctx context.Context, db *gorm.DB, eventID string) (*payment.WebhookEvent, error) {
						return nil, nil
					},
					createFunc: func(ctx context.Context, db *gorm.DB, event *payment.WebhookEvent) error {
						return nil
					},
					markProcessedFunc: func(ctx context.Context, db *gorm.DB, id int64, processed int8, errorMsg string) error {
						if processed != 2 {
							t.Errorf("expected processed=2 for failure, got %d", processed)
						}
						return nil
					},
				}
			},
			wantErr: true,
			errMsg:  "payment not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &service{
				db:                  nil,
				paymentRepo:         tt.setupMockPaymentRepo(),
				refundRepo:          &mockRefundRepo{},
				transactionRepo:     &mockTransactionRepo{},
				webhookEventRepo:    tt.setupMockWebhookRepo(),
				paymentSettingsRepo: &mockPaymentSettingsRepo{},
				idGen:               &mockIDGen{},
			}

			err := svc.HandleWebhook(context.Background(), tt.event)

			if tt.wantErr {
				if err == nil {
					t.Errorf("HandleWebhook() error = nil, want error containing %q", tt.errMsg)
					return
				}
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("HandleWebhook() error = %q, want error containing %q", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("HandleWebhook() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestHandleWebhook_Idempotent(t *testing.T) {
	// Test that duplicate webhooks are handled idempotently
	t.Run("duplicate event returns nil without reprocessing", func(t *testing.T) {
		callCount := 0

		mockPaymentRepo := &mockPaymentRepo{
			findByChannelPaymentIDFunc: func(ctx context.Context, db *gorm.DB, channelPaymentID string) (*payment.Payment, error) {
				callCount++
				return &payment.Payment{
					Model:            application.Model{ID: 1},
					Status:           payment.PaymentStatusSuccess,
					ChannelPaymentID: "pi_success123",
				}, nil
			},
			updateFunc: func(ctx context.Context, db *gorm.DB, p *payment.Payment) error {
				return nil
			},
		}

		mockWebhookRepo := &mockWebhookEventRepo{
			findByEventIDFunc: func(ctx context.Context, db *gorm.DB, eventID string) (*payment.WebhookEvent, error) {
				// Simulate existing processed event
				if eventID == "evt_duplicate" {
					return &payment.WebhookEvent{
						Model:      application.Model{ID: 1},
						EventID:    "evt_duplicate",
						EventType:  "payment_intent.succeeded",
						Processed:  1, // Already processed
						RawPayload: `{"id": "pi_success123"}`,
					}, nil
				}
				return nil, nil
			},
		}

		svc := &service{
			db:                  nil,
			paymentRepo:         mockPaymentRepo,
			refundRepo:          &mockRefundRepo{},
			transactionRepo:     &mockTransactionRepo{},
			webhookEventRepo:    mockWebhookRepo,
			paymentSettingsRepo: &mockPaymentSettingsRepo{},
			idGen:               &mockIDGen{},
		}

		event := &WebhookEvent{
			EventID:    "evt_duplicate",
			EventType:  "payment_intent.succeeded",
			ResourceID: "pi_success123",
			RawPayload: `{"id": "pi_success123"}`,
		}

		// First call - should skip due to already processed
		err := svc.HandleWebhook(context.Background(), event)
		if err != nil {
			t.Errorf("HandleWebhook() error = %v, want nil for already processed event", err)
		}

		// Verify payment repo was NOT called because event was already processed
		if callCount != 0 {
			t.Errorf("paymentRepo.FindByChannelPaymentID called %d times, expected 0 for already processed event", callCount)
		}
	})

	t.Run("unprocessed event is processed", func(t *testing.T) {
		callCount := 0

		mockPaymentRepo := &mockPaymentRepo{
			findByChannelPaymentIDFunc: func(ctx context.Context, db *gorm.DB, channelPaymentID string) (*payment.Payment, error) {
				callCount++
				return &payment.Payment{
					Model:           application.Model{ID: 1},
					Status:          payment.PaymentStatusPending,
					ChannelIntentID: "pi_new123",
				}, nil
			},
			updateFunc: func(ctx context.Context, db *gorm.DB, p *payment.Payment) error {
				return nil
			},
		}

		mockWebhookRepo := &mockWebhookEventRepo{
			findByEventIDFunc: func(ctx context.Context, db *gorm.DB, eventID string) (*payment.WebhookEvent, error) {
				return nil, nil // Event not found, will be created
			},
			createFunc: func(ctx context.Context, db *gorm.DB, event *payment.WebhookEvent) error {
				return nil
			},
			markProcessedFunc: func(ctx context.Context, db *gorm.DB, id int64, processed int8, errorMsg string) error {
				return nil
			},
		}

		svc := &service{
			db:                  nil,
			paymentRepo:         mockPaymentRepo,
			refundRepo:          &mockRefundRepo{},
			transactionRepo:     &mockTransactionRepo{},
			webhookEventRepo:    mockWebhookRepo,
			paymentSettingsRepo: &mockPaymentSettingsRepo{},
			idGen:               &mockIDGen{},
		}

		event := &WebhookEvent{
			EventID:    "evt_new",
			EventType:  "payment_intent.succeeded",
			ResourceID: "pi_new123",
			RawPayload: `{"id": "pi_new123"}`,
		}

		err := svc.HandleWebhook(context.Background(), event)
		if err != nil {
			t.Errorf("HandleWebhook() error = %v, want nil", err)
		}

		// Verify payment repo WAS called
		if callCount != 1 {
			t.Errorf("paymentRepo.FindByChannelPaymentID called %d times, expected 1", callCount)
		}
	})
}

func TestHandleWebhook_UnknownEventType(t *testing.T) {
	// Test that unknown event types are marked as processed without error
	t.Run("unknown event type is marked as processed", func(t *testing.T) {
		webhookRepo := &mockWebhookEventRepo{
			findByEventIDFunc: func(ctx context.Context, db *gorm.DB, eventID string) (*payment.WebhookEvent, error) {
				return nil, nil
			},
			createFunc: func(ctx context.Context, db *gorm.DB, event *payment.WebhookEvent) error {
				return nil
			},
			markProcessedFunc: func(ctx context.Context, db *gorm.DB, id int64, processed int8, errorMsg string) error {
				// Unknown events should still be marked as processed (not failed)
				if processed != 1 {
					t.Errorf("expected processed=1 for unknown event type, got %d", processed)
				}
				return nil
			},
		}

		svc := &service{
			db:                  nil,
			paymentRepo:         &mockPaymentRepo{},
			refundRepo:          &mockRefundRepo{},
			transactionRepo:     &mockTransactionRepo{},
			webhookEventRepo:    webhookRepo,
			paymentSettingsRepo: &mockPaymentSettingsRepo{},
			idGen:               &mockIDGen{},
		}

		event := &WebhookEvent{
			EventID:    "evt_unknown",
			EventType:  "customer.created", // Unknown event type
			ResourceID: "cus_123",
			RawPayload: `{"id": "cus_123"}`,
		}

		err := svc.HandleWebhook(context.Background(), event)
		if err != nil {
			t.Errorf("HandleWebhook() error = %v, want nil for unknown event type", err)
		}
	})
}

func TestHandleWebhook_ExistingUnprocessedEvent(t *testing.T) {
	// Test that when an event exists but is not processed, it gets processed
	t.Run("existing unprocessed event is processed", func(t *testing.T) {
		paymentUpdateCalled := false

		mockPaymentRepo := &mockPaymentRepo{
			findByChannelPaymentIDFunc: func(ctx context.Context, db *gorm.DB, channelPaymentID string) (*payment.Payment, error) {
				return &payment.Payment{
					Model:            application.Model{ID: 1},
					Status:           payment.PaymentStatusSuccess,
					ChannelPaymentID: "pi_existing",
				}, nil
			},
			updateFunc: func(ctx context.Context, db *gorm.DB, p *payment.Payment) error {
				paymentUpdateCalled = true
				return nil
			},
		}

		mockWebhookRepo := &mockWebhookEventRepo{
			findByEventIDFunc: func(ctx context.Context, db *gorm.DB, eventID string) (*payment.WebhookEvent, error) {
				// Return existing event that is NOT processed
				return &payment.WebhookEvent{
					Model:      application.Model{ID: 50},
					EventID:    "evt_existing",
					EventType:  "payment_intent.succeeded",
					Processed:  0, // Not yet processed
					RawPayload: `{"id": "pi_existing"}`,
				}, nil
			},
			markProcessedFunc: func(ctx context.Context, db *gorm.DB, id int64, processed int8, errorMsg string) error {
				if id != 50 {
					t.Errorf("expected id=50 (existing event), got %d", id)
				}
				return nil
			},
		}

		svc := &service{
			db:                  nil,
			paymentRepo:         mockPaymentRepo,
			refundRepo:          &mockRefundRepo{},
			transactionRepo:     &mockTransactionRepo{},
			webhookEventRepo:    mockWebhookRepo,
			paymentSettingsRepo: &mockPaymentSettingsRepo{},
			idGen:               &mockIDGen{},
		}

		event := &WebhookEvent{
			EventID:    "evt_existing",
			EventType:  "payment_intent.succeeded",
			ResourceID: "pi_existing",
			RawPayload: `{"id": "pi_existing"}`,
		}

		err := svc.HandleWebhook(context.Background(), event)
		if err != nil {
			t.Errorf("HandleWebhook() error = %v, want nil", err)
		}

		if !paymentUpdateCalled {
			t.Error("paymentRepo.Update was not called for existing unprocessed event")
		}
	})
}

// =============================================================================
// Webhook Event Creation Tests
// =============================================================================

func TestHandleWebhook_EventCreation(t *testing.T) {
	t.Run("webhook event is created with correct fields", func(t *testing.T) {
		var createdEvent *payment.WebhookEvent

		mockWebhookRepo := &mockWebhookEventRepo{
			findByEventIDFunc: func(ctx context.Context, db *gorm.DB, eventID string) (*payment.WebhookEvent, error) {
				return nil, nil
			},
			createFunc: func(ctx context.Context, db *gorm.DB, event *payment.WebhookEvent) error {
				createdEvent = event
				return nil
			},
			markProcessedFunc: func(ctx context.Context, db *gorm.DB, id int64, processed int8, errorMsg string) error {
				return nil
			},
		}

		mockPaymentRepo := &mockPaymentRepo{
			findByChannelPaymentIDFunc: func(ctx context.Context, db *gorm.DB, channelPaymentID string) (*payment.Payment, error) {
				return &payment.Payment{
					Model:            application.Model{ID: 1},
					Status:           payment.PaymentStatusSuccess,
					ChannelPaymentID: "pi_create_test",
				}, nil
			},
			updateFunc: func(ctx context.Context, db *gorm.DB, p *payment.Payment) error {
				return nil
			},
		}

		svc := &service{
			db:                  nil,
			paymentRepo:         mockPaymentRepo,
			refundRepo:          &mockRefundRepo{},
			transactionRepo:     &mockTransactionRepo{},
			webhookEventRepo:    mockWebhookRepo,
			paymentSettingsRepo: &mockPaymentSettingsRepo{},
			idGen:               &mockIDGen{},
		}

		event := &WebhookEvent{
			EventID:    "evt_create_test",
			EventType:  "payment_intent.succeeded",
			ResourceID: "pi_create_test",
			RawPayload: `{"id": "pi_create_test", "amount": 5000}`,
		}

		err := svc.HandleWebhook(context.Background(), event)
		if err != nil {
			t.Errorf("HandleWebhook() error = %v, want nil", err)
		}

		if createdEvent == nil {
			t.Fatal("createFunc was not called")
		}

		if createdEvent.EventID != "evt_create_test" {
			t.Errorf("EventID = %q, want %q", createdEvent.EventID, "evt_create_test")
		}
		if createdEvent.EventType != "payment_intent.succeeded" {
			t.Errorf("EventType = %q, want %q", createdEvent.EventType, "payment_intent.succeeded")
		}
		if createdEvent.ResourceID != "pi_create_test" {
			t.Errorf("ResourceID = %q, want %q", createdEvent.ResourceID, "pi_create_test")
		}
		if createdEvent.RawPayload != `{"id": "pi_create_test", "amount": 5000}` {
			t.Errorf("RawPayload = %q, want %q", createdEvent.RawPayload, `{"id": "pi_create_test", "amount": 5000}`)
		}
	})
}

// =============================================================================
// Helper Functions
// =============================================================================

// contains checks if s contains substr
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
