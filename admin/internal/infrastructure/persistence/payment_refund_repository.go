package persistence

import (
	"context"
	"errors"

	"github.com/colinrs/shopjoy/admin/internal/domain/payment"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type paymentRefundRepository struct{}

func NewPaymentRefundRepository() payment.PaymentRefundRepository {
	return &paymentRefundRepository{}
}

func (r *paymentRefundRepository) Create(ctx context.Context, db *gorm.DB, refund *payment.PaymentRefund) error {
	return db.WithContext(ctx).Create(refund).Error
}

func (r *paymentRefundRepository) Update(ctx context.Context, db *gorm.DB, refund *payment.PaymentRefund) error {
	return db.WithContext(ctx).Save(refund).Error
}

func (r *paymentRefundRepository) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*payment.PaymentRefund, error) {
	var refund payment.PaymentRefund
	query := db.WithContext(ctx).Where("id = ?", id)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	err := query.First(&refund).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrPaymentRefundNotFound
		}
		return nil, err
	}
	return &refund, nil
}

func (r *paymentRefundRepository) FindByRefundNo(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, refundNo string) (*payment.PaymentRefund, error) {
	var refund payment.PaymentRefund
	query := db.WithContext(ctx).Where("refund_no = ?", refundNo)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	err := query.First(&refund).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrPaymentRefundNotFound
		}
		return nil, err
	}
	return &refund, nil
}

func (r *paymentRefundRepository) FindByIdempotencyKey(ctx context.Context, db *gorm.DB, idempotencyKey string) (*payment.PaymentRefund, error) {
	var refund payment.PaymentRefund
	err := db.WithContext(ctx).
		Where("idempotency_key = ?", idempotencyKey).
		First(&refund).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is not an error for idempotency check
		}
		return nil, err
	}
	return &refund, nil
}

func (r *paymentRefundRepository) FindByPaymentID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, paymentID int64) ([]*payment.PaymentRefund, error) {
	query := db.WithContext(ctx).Model(&payment.PaymentRefund{}).Where("payment_id = ?", paymentID)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var refunds []*payment.PaymentRefund
	err := query.Order("id DESC").Find(&refunds).Error
	return refunds, err
}

func (r *paymentRefundRepository) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query payment.RefundQuery) ([]*payment.PaymentRefund, int64, error) {
	query.Validate()

	dbQuery := db.WithContext(ctx).Model(&payment.PaymentRefund{})

	if tenantID != 0 {
		dbQuery = dbQuery.Where("tenant_id = ?", tenantID.Int64())
	}

	if query.OrderID > 0 {
		dbQuery = dbQuery.Where("order_id = ?", query.OrderID)
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

	var list []*payment.PaymentRefund
	err := dbQuery.Order("id DESC").
		Offset(query.Offset()).
		Limit(query.Limit()).
		Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *paymentRefundRepository) GetTotalRefundedAmount(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, paymentID int64) (decimal.Decimal, error) {
	query := db.WithContext(ctx).
		Model(&payment.PaymentRefund{}).
		Where("payment_id = ? AND status = ?", paymentID, payment.PaymentRefundStatusSucceeded)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	var total decimal.Decimal
	err := query.Select("COALESCE(SUM(amount), 0)").Scan(&total).Error
	return total, err
}
