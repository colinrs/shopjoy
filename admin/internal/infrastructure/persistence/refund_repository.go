package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type refundRepo struct{}

func NewRefundRepository() fulfillment.RefundRepository {
	return &refundRepo{}
}

// refundModel represents the database model for Refund
type refundModel struct {
	ID           int64           `gorm:"column:id;primaryKey;autoIncrement:false"`
	TenantID     int64           `gorm:"column:tenant_id;not null;index"`
	OrderID      int64           `gorm:"column:order_id;not null;index"`
	RefundNo     string          `gorm:"column:refund_no;size:32;not null;uniqueIndex:uk_refund_no"`
	UserID       int64           `gorm:"column:user_id;not null;index"`
	Type         int              `gorm:"column:type;not null;default:1"`
	Status       int              `gorm:"column:status;not null;default:0;index"`
	ReasonType   string          `gorm:"column:reason_type;size:50;not null;default:''"`
	Reason       string          `gorm:"column:reason;size:500;not null;default:''"`
	Description string          `gorm:"column:description;type:text"`
	Images       string          `gorm:"column:images;type:json"`
	Amount       decimal.Decimal `gorm:"column:amount;type:decimal(19,4);not null;default:0"`
	Currency     string          `gorm:"column:currency;size:10;not null;default:'CNY'"`
	RejectReason string          `gorm:"column:reject_reason;size:500;not null;default:''"`
	ApprovedAt   *time.Time     `gorm:"column:approved_at"`
	ApprovedBy   int64           `gorm:"column:approved_by;not null;default:0"`
	CompletedAt  *time.Time     `gorm:"column:completed_at"`
	DeletedAt    *time.Time     `gorm:"column:deleted_at;index"`
	CreatedAt    time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;not null"`
}

func (refundModel) TableName() string {
	return "refunds"
}

func (m *refundModel) toEntity() *fulfillment.Refund {
	var deletedAt gorm.DeletedAt
	if m.DeletedAt != nil {
		deletedAt = gorm.DeletedAt{Time: *m.DeletedAt, Valid: true}
	}
	return &fulfillment.Refund{
		Model: application.Model{
			ID:        m.ID,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
			DeletedAt: deletedAt,
		},
		TenantID:     shared.TenantID(m.TenantID),
		OrderID:      m.OrderID,
		RefundNo:     m.RefundNo,
		UserID:       m.UserID,
		Type:         fulfillment.RefundType(m.Type),
		Status:       fulfillment.RefundStatus(m.Status),
		ReasonType:   m.ReasonType,
		Reason:       m.Reason,
		Description:  m.Description,
		Images:       m.Images,
		Amount:       m.Amount,
		Currency:     m.Currency,
		RejectReason: m.RejectReason,
		ApprovedAt:   m.ApprovedAt,
		ApprovedBy:   m.ApprovedBy,
		CompletedAt:  m.CompletedAt,
	}
}

func fromRefundEntity(r *fulfillment.Refund) *refundModel {
	var deletedAt *time.Time
	if r.Model.DeletedAt.Valid {
		deletedAt = &r.Model.DeletedAt.Time
	}
	return &refundModel{
		ID:           r.Model.ID,
		TenantID:     r.TenantID.Int64(),
		OrderID:      r.OrderID,
		RefundNo:     r.RefundNo,
		UserID:       r.UserID,
		Type:         int(r.Type),
		Status:       int(r.Status),
		ReasonType:   r.ReasonType,
		Reason:       r.Reason,
		Description:  r.Description,
		Images:       r.Images,
		Amount:       r.Amount,
		Currency:     r.Currency,
		RejectReason: r.RejectReason,
		ApprovedAt:   r.ApprovedAt,
		ApprovedBy:   r.ApprovedBy,
		CompletedAt:  r.CompletedAt,
		DeletedAt:    deletedAt,
		CreatedAt:    r.Model.CreatedAt,
		UpdatedAt:    r.Model.UpdatedAt,
	}
}

// Create inserts a new refund
func (r *refundRepo) Create(ctx context.Context, db *gorm.DB, refund *fulfillment.Refund) error {
	model := fromRefundEntity(refund)
	return db.WithContext(ctx).Create(model).Error
}

// Update updates an existing refund
func (r *refundRepo) Update(ctx context.Context, db *gorm.DB, refund *fulfillment.Refund) error {
	model := fromRefundEntity(refund)
	return db.WithContext(ctx).
		Model(&refundModel{}).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", refund.ID, refund.TenantID.Int64()).
		Updates(map[string]any{
			"status":        model.Status,
			"reject_reason": model.RejectReason,
			"approved_at":   model.ApprovedAt,
			"approved_by":   model.ApprovedBy,
			"completed_at":  model.CompletedAt,
			"updated_at":    model.UpdatedAt,
		}).Error
}

// FindByID finds a refund by ID
func (r *refundRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*fulfillment.Refund, error) {
	query := db.WithContext(ctx).Where("deleted_at IS NULL")
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var model refundModel
	err := query.First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrRefundNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

// FindByRefundNo finds a refund by refund number
func (r *refundRepo) FindByRefundNo(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, refundNo string) (*fulfillment.Refund, error) {
	query := db.WithContext(ctx).Model(&refundModel{}).Where("refund_no = ? AND deleted_at IS NULL", refundNo)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var model refundModel
	err := query.First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrRefundNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

// FindByOrderID finds all refunds for an order
func (r *refundRepo) FindByOrderID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderID int64) ([]*fulfillment.Refund, error) {
	query := db.WithContext(ctx).Model(&refundModel{}).Where("order_id = ? AND deleted_at IS NULL", orderID)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var models []refundModel
	err := query.Order("created_at DESC").Find(&models).Error
	if err != nil {
		return nil, err
	}

	refunds := make([]*fulfillment.Refund, len(models))
	for i, m := range models {
		refunds[i] = m.toEntity()
	}
	return refunds, nil
}

// FindPendingByOrderID finds the pending refund for an order
func (r *refundRepo) FindPendingByOrderID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderID int64) (*fulfillment.Refund, error) {
	query := db.WithContext(ctx).Model(&refundModel{}).
		Where("order_id = ? AND status = ? AND deleted_at IS NULL", orderID, fulfillment.RefundStatusPending)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var model refundModel
	err := query.First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No pending refund is not an error
		}
		return nil, err
	}
	return model.toEntity(), nil
}

// FindByUserID finds refunds for a user with pagination
func (r *refundRepo) FindByUserID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64, query fulfillment.RefundQuery) ([]*fulfillment.Refund, int64, error) {
	query.Validate()

	dbQuery := db.WithContext(ctx).Model(&refundModel{}).
		Where("user_id = ? AND deleted_at IS NULL", userID)

	if tenantID != 0 {
		dbQuery = dbQuery.Where("tenant_id = ?", tenantID.Int64())
	}

	if query.Status.IsValid() {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}
	if !query.StartTime.IsZero() {
		dbQuery = dbQuery.Where("created_at >= ?", query.StartTime)
	}
	if !query.EndTime.IsZero() {
		dbQuery = dbQuery.Where("created_at < ?", query.EndTime)
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []refundModel
	err := dbQuery.Order("created_at DESC").
		Offset(query.Offset()).
		Limit(query.Limit()).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	refunds := make([]*fulfillment.Refund, len(models))
	for i, m := range models {
		refunds[i] = m.toEntity()
	}
	return refunds, total, nil
}

// FindList finds refunds with pagination and filters
func (r *refundRepo) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query fulfillment.RefundQuery) ([]*fulfillment.Refund, int64, error) {
	query.Validate()

	dbQuery := db.WithContext(ctx).Model(&refundModel{}).Where("deleted_at IS NULL")

	if tenantID != 0 {
		dbQuery = dbQuery.Where("tenant_id = ?", tenantID.Int64())
	}

	if query.OrderID != 0 {
		dbQuery = dbQuery.Where("order_id = ?", query.OrderID)
	}
	if query.RefundNo != "" {
		dbQuery = dbQuery.Where("refund_no LIKE ?", escapeLikePattern(query.RefundNo))
	}
	if query.UserID > 0 {
		dbQuery = dbQuery.Where("user_id = ?", query.UserID)
	}
	if query.Status.IsValid() {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}
	if query.ReasonType != "" {
		dbQuery = dbQuery.Where("reason_type = ?", query.ReasonType)
	}
	if !query.StartTime.IsZero() {
		dbQuery = dbQuery.Where("created_at >= ?", query.StartTime)
	}
	if !query.EndTime.IsZero() {
		dbQuery = dbQuery.Where("created_at < ?", query.EndTime)
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []refundModel
	err := dbQuery.Order("created_at DESC").
		Offset(query.Offset()).
		Limit(query.Limit()).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	refunds := make([]*fulfillment.Refund, len(models))
	for i, m := range models {
		refunds[i] = m.toEntity()
	}
	return refunds, total, nil
}

// Delete soft deletes a refund
func (r *refundRepo) Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error {
	query := db.WithContext(ctx).Model(&refundModel{}).Where("id = ? AND deleted_at IS NULL", id)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	now := time.Now().UTC()
	result := query.Update("deleted_at", now)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return code.ErrRefundNotFound
	}
	return nil
}

// CountByStatus counts refunds by status
func (r *refundRepo) CountByStatus(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, status fulfillment.RefundStatus) (int64, error) {
	query := db.WithContext(ctx).Model(&refundModel{}).Where("status = ? AND deleted_at IS NULL", status)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var count int64
	err := query.Count(&count).Error
	return count, err
}