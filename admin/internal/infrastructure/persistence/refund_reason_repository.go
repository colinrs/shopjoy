package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/code"
	"gorm.io/gorm"
)

type refundReasonRepo struct{}

func NewRefundReasonRepository() fulfillment.RefundReasonRepository {
	return &refundReasonRepo{}
}

// refundReasonModel represents the database model for RefundReason
type refundReasonModel struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement:false"`
	Code      string    `gorm:"column:code;size:50;not null;uniqueIndex"`
	Name      string    `gorm:"column:name;size:100;not null"`
	Sort      int       `gorm:"column:sort;not null;default:0"`
	IsActive  int       `gorm:"column:is_active;not null;default:1"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
}

func (refundReasonModel) TableName() string {
	return "refund_reasons"
}

func (m *refundReasonModel) toEntity() *fulfillment.RefundReason {
	return &fulfillment.RefundReason{
		Model:    application.Model{ID: m.ID, CreatedAt: m.CreatedAt.UTC()},
		Code:     m.Code,
		Name:     m.Name,
		Sort:     m.Sort,
		IsActive: m.IsActive == 1,
	}
}

func fromRefundReasonEntity(r *fulfillment.RefundReason) *refundReasonModel {
	isActive := 0
	if r.IsActive {
		isActive = 1
	}
	return &refundReasonModel{
		ID:        r.Model.ID,
		Code:      r.Code,
		Name:      r.Name,
		Sort:      r.Sort,
		IsActive:  isActive,
		CreatedAt: r.Model.CreatedAt,
	}
}

// Create inserts a new refund reason
func (r *refundReasonRepo) Create(ctx context.Context, db *gorm.DB, reason *fulfillment.RefundReason) error {
	model := fromRefundReasonEntity(reason)
	return db.WithContext(ctx).Create(model).Error
}

// FindByID finds a refund reason by ID
func (r *refundReasonRepo) FindByID(ctx context.Context, db *gorm.DB, id int64) (*fulfillment.RefundReason, error) {
	var model refundReasonModel
	err := db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrRefundReasonNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

// FindByCode finds a refund reason by code
func (r *refundReasonRepo) FindByCode(ctx context.Context, db *gorm.DB, codeStr string) (*fulfillment.RefundReason, error) {
	var model refundReasonModel
	err := db.WithContext(ctx).Where("code = ?", codeStr).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrRefundReasonNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

// FindAll finds all refund reasons
func (r *refundReasonRepo) FindAll(ctx context.Context, db *gorm.DB) ([]*fulfillment.RefundReason, error) {
	var models []refundReasonModel
	err := db.WithContext(ctx).Order("sort ASC, id ASC").Find(&models).Error
	if err != nil {
		return nil, err
	}

	reasons := make([]*fulfillment.RefundReason, len(models))
	for i, m := range models {
		reasons[i] = m.toEntity()
	}
	return reasons, nil
}

// FindActive finds all active refund reasons
func (r *refundReasonRepo) FindActive(ctx context.Context, db *gorm.DB) ([]*fulfillment.RefundReason, error) {
	var models []refundReasonModel
	err := db.WithContext(ctx).
		Where("is_active = 1").
		Order("sort ASC, id ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	reasons := make([]*fulfillment.RefundReason, len(models))
	for i, m := range models {
		reasons[i] = m.toEntity()
	}
	return reasons, nil
}

// Update updates a refund reason
func (r *refundReasonRepo) Update(ctx context.Context, db *gorm.DB, reason *fulfillment.RefundReason) error {
	model := fromRefundReasonEntity(reason)
	return db.WithContext(ctx).
		Model(&refundReasonModel{}).
		Where("id = ?", reason.Model.ID).
		Updates(map[string]any{
			"name":      model.Name,
			"sort":      model.Sort,
			"is_active": model.IsActive,
		}).Error
}
