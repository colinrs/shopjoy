package persistence

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/payment"
	"gorm.io/gorm"
)

type webhookEventRepository struct{}

func NewWebhookEventRepository() payment.WebhookEventRepository {
	return &webhookEventRepository{}
}

func (r *webhookEventRepository) Create(ctx context.Context, db *gorm.DB, event *payment.WebhookEvent) error {
	return db.WithContext(ctx).Create(event).Error
}

func (r *webhookEventRepository) FindByEventID(ctx context.Context, db *gorm.DB, eventID string) (*payment.WebhookEvent, error) {
	var event payment.WebhookEvent
	err := db.WithContext(ctx).
		Where("event_id = ?", eventID).
		First(&event).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

func (r *webhookEventRepository) MarkProcessed(ctx context.Context, db *gorm.DB, id int64, processed int8, errorMsg string) error {
	now := time.Now().UTC()
	return db.WithContext(ctx).
		Model(&payment.WebhookEvent{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"processed":     processed,
			"error_message": errorMsg,
			"processed_at":  now,
		}).Error
}