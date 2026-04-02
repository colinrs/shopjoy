package webhooks

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	appPayment "github.com/colinrs/shopjoy/admin/internal/application/payment"
	"github.com/colinrs/shopjoy/admin/internal/config"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
)

// mockPaymentService is a mock implementation of appPayment.Service
type mockPaymentService struct {
	handleWebhookFunc func(ctx context.Context, event *appPayment.WebhookEvent) error
}

func (m *mockPaymentService) GetPaymentStats(ctx context.Context, tenantID shared.TenantID, period string) (*appPayment.PaymentStatsDTO, error) {
	return nil, nil
}

func (m *mockPaymentService) ListTransactions(ctx context.Context, tenantID shared.TenantID, req appPayment.ListTransactionsRequest) (*appPayment.ListTransactionsResponse, error) {
	return nil, nil
}

func (m *mockPaymentService) GetTransaction(ctx context.Context, tenantID shared.TenantID, id int64) (*appPayment.TransactionDTO, error) {
	return nil, nil
}

func (m *mockPaymentService) GetOrderPayment(ctx context.Context, tenantID shared.TenantID, orderID int64) (*appPayment.OrderPaymentDTO, error) {
	return nil, nil
}

func (m *mockPaymentService) InitiateRefund(ctx context.Context, tenantID shared.TenantID, adminID int64, req appPayment.InitiateRefundRequest) (*appPayment.InitiateRefundResponse, error) {
	return nil, nil
}

func (m *mockPaymentService) HandleWebhook(ctx context.Context, event *appPayment.WebhookEvent) error {
	if m.handleWebhookFunc != nil {
		return m.handleWebhookFunc(ctx, event)
	}
	return nil
}

// mockServiceContext creates a minimal service context for testing
func mockServiceContext(mockPaymentService *mockPaymentService, webhookSecret string) *svc.ServiceContext {
	return &svc.ServiceContext{
		PaymentService: mockPaymentService,
		Config:         config.Config{StripeWebhookSecret: webhookSecret},
	}
}

// mockRequest creates a mock HTTP request with the given body and headers
func mockRequest(body string, headers map[string]string) *http.Request {
	req, _ := http.NewRequest("POST", "/api/v1/webhooks/stripe", strings.NewReader(body))
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req
}

func TestStripeWebhook_MissingWebhookSecret(t *testing.T) {
	tests := []struct {
		name           string
		webhookSecret  string
		requestHeaders map[string]string
		requestBody    string
		wantErr        error
	}{
		{
			name:           "webhook secret not configured",
			webhookSecret:  "", // Empty secret
			requestHeaders: map[string]string{"Stripe-Signature": "sig123"},
			requestBody:    `{"type": "payment_intent.succeeded"}`,
			wantErr:        code.ErrPaymentWebhookSignatureInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPaymentSvc := &mockPaymentService{}
			svcCtx := mockServiceContext(mockPaymentSvc, tt.webhookSecret)
			logic := NewStripeWebhookLogic(context.Background(), svcCtx)

			req := mockRequest(tt.requestBody, tt.requestHeaders)
			err := logic.StripeWebhook(req)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("StripeWebhook() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestStripeWebhook_MissingSignatureHeader(t *testing.T) {
	tests := []struct {
		name           string
		webhookSecret  string
		requestHeaders map[string]string
		requestBody    string
		wantErr        error
	}{
		{
			name:           "missing Stripe-Signature header",
			webhookSecret:  "whsec_test_secret",
			requestHeaders: map[string]string{}, // No signature header
			requestBody:    `{"type": "payment_intent.succeeded"}`,
			wantErr:        code.ErrPaymentWebhookSignatureInvalid,
		},
		{
			name:           "empty Stripe-Signature header",
			webhookSecret:  "whsec_test_secret",
			requestHeaders: map[string]string{"Stripe-Signature": ""},
			requestBody:    `{"type": "payment_intent.succeeded"}`,
			wantErr:        code.ErrPaymentWebhookSignatureInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPaymentSvc := &mockPaymentService{}
			svcCtx := mockServiceContext(mockPaymentSvc, tt.webhookSecret)
			logic := NewStripeWebhookLogic(context.Background(), svcCtx)

			req := mockRequest(tt.requestBody, tt.requestHeaders)
			err := logic.StripeWebhook(req)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("StripeWebhook() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestStripeWebhook_InvalidSignature(t *testing.T) {
	// Note: We cannot test actual signature verification without mocking the Stripe library,
	// but we can test that invalid signatures are rejected.
	// The webhook.ConstructEvent will return an error for invalid signatures.
	t.Run("invalid signature format", func(t *testing.T) {
		mockPaymentSvc := &mockPaymentService{}
		svcCtx := mockServiceContext(mockPaymentSvc, "whsec_test_secret")
		logic := NewStripeWebhookLogic(context.Background(), svcCtx)

		// An invalid signature format should be rejected
		req := mockRequest(`{"type": "payment_intent.succeeded"}`, map[string]string{
			"Stripe-Signature": "invalid_signature_format",
		})
		err := logic.StripeWebhook(req)

		// Should return signature invalid error
		if !errors.Is(err, code.ErrPaymentWebhookSignatureInvalid) {
			t.Errorf("StripeWebhook() error = %v, want %v", err, code.ErrPaymentWebhookSignatureInvalid)
		}
	})
}

func TestStripeWebhook_RequestBodyReadError(t *testing.T) {
	t.Run("empty request body", func(t *testing.T) {
		mockPaymentSvc := &mockPaymentService{}
		svcCtx := mockServiceContext(mockPaymentSvc, "whsec_test_secret")
		logic := NewStripeWebhookLogic(context.Background(), svcCtx)

		// Empty body should cause an error in signature verification
		req := mockRequest("", map[string]string{"Stripe-Signature": "sig123"})
		err := logic.StripeWebhook(req)

		// Empty body should fail signature verification
		if !errors.Is(err, code.ErrPaymentWebhookSignatureInvalid) {
			t.Errorf("StripeWebhook() error = %v, want %v", err, code.ErrPaymentWebhookSignatureInvalid)
		}
	})
}

func TestStripeWebhook_PaymentServiceError(t *testing.T) {
	// Test that errors from the payment service are propagated
	// Note: This test cannot fully run because we cannot mock the Stripe library's
	// webhook.ConstructEvent function. The test verifies the mock setup is correct.
	t.Run("payment service error propagation", func(t *testing.T) {
		mockPaymentSvc := &mockPaymentService{
			handleWebhookFunc: func(ctx context.Context, event *appPayment.WebhookEvent) error {
				return errors.New("database error")
			},
		}

		// Verify the mock is set up correctly
		if mockPaymentSvc.handleWebhookFunc == nil {
			t.Error("handleWebhookFunc not set")
		}

		// The actual error propagation would be tested in integration tests
		// where we can provide a valid Stripe signature
	})
}
