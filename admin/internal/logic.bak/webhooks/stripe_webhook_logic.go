package webhooks

import (
	"context"
	"io"
	"net/http"

	appPayment "github.com/colinrs/shopjoy/admin/internal/application/payment"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
	"github.com/zeromicro/go-zero/core/logx"
)

type StripeWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStripeWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) StripeWebhookLogic {
	return StripeWebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StripeWebhookLogic) StripeWebhook(r *http.Request) error {
	// Get webhook secret from config
	webhookSecret := l.svcCtx.Config.StripeWebhookSecret
	if webhookSecret == "" {
		l.Logger.Errorf("Stripe webhook secret not configured")
		return code.ErrPaymentWebhookSignatureInvalid
	}

	// Get signature header
	stripeSignature := r.Header.Get("Stripe-Signature")
	if stripeSignature == "" {
		l.Logger.Errorf("Missing Stripe-Signature header")
		return code.ErrPaymentWebhookSignatureInvalid
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		l.Logger.Errorf("Failed to read webhook body: %v", err)
		return err
	}
	defer r.Body.Close()

	// Verify Stripe webhook signature
	event, err := webhook.ConstructEvent(body, stripeSignature, webhookSecret)
	if err != nil {
		l.Logger.Errorf("Failed to verify webhook signature: %v", err)
		return code.ErrPaymentWebhookSignatureInvalid
	}

	// Extract event information
	eventID := event.ID
	eventType := string(event.Type)
	resourceID := extractResourceID(event)

	// Create webhook event for processing
	webhookEvent := &appPayment.WebhookEvent{
		EventID:    eventID,
		EventType:  eventType,
		ResourceID: resourceID,
		RawPayload: string(body),
	}

	// Process webhook through payment service
	if err := l.svcCtx.PaymentService.HandleWebhook(l.ctx, webhookEvent); err != nil {
		l.Logger.Errorf("Failed to handle webhook event %s: %v", eventID, err)
		return err
	}

	l.Logger.Infof("Successfully processed webhook event %s of type %s", eventID, eventType)
	return nil
}

// extractResourceID extracts the resource ID from a Stripe event
func extractResourceID(event stripe.Event) string {
	switch event.Type {
	case "payment_intent.succeeded", "payment_intent.payment_failed",
		"payment_intent.created", "payment_intent.canceled":
		if event.Data.Object["id"] != nil {
			if id, ok := event.Data.Object["id"].(string); ok {
				return id
			}
		}
	case "charge.succeeded", "charge.failed", "charge.refunded":
		if event.Data.Object["id"] != nil {
			if id, ok := event.Data.Object["id"].(string); ok {
				return id
			}
		}
	case "refund.created", "refund.updated":
		if event.Data.Object["id"] != nil {
			if id, ok := event.Data.Object["id"].(string); ok {
				return id
			}
		}
	}
	return ""
}
