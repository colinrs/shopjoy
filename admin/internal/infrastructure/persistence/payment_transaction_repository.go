package persistence

import (
	"context"
	"errors"

	"github.com/colinrs/shopjoy/admin/internal/domain/payment"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type paymentTransactionRepository struct{}

func NewPaymentTransactionRepository() payment.PaymentTransactionRepository {
	return &paymentTransactionRepository{}
}

func (r *paymentTransactionRepository) Create(ctx context.Context, db *gorm.DB, txn *payment.PaymentTransaction) error {
	return db.WithContext(ctx).Create(txn).Error
}

func (r *paymentTransactionRepository) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*payment.PaymentTransaction, error) {
	var txn payment.PaymentTransaction
	query := db.WithContext(ctx).Where("id = ?", id)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	err := query.First(&txn).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrTransactionNotFound
		}
		return nil, err
	}
	return &txn, nil
}

func (r *paymentTransactionRepository) FindByTransactionID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, txnID string) (*payment.PaymentTransaction, error) {
	var txn payment.PaymentTransaction
	query := db.WithContext(ctx).Where("transaction_id = ?", txnID)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	err := query.First(&txn).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrTransactionNotFound
		}
		return nil, err
	}
	return &txn, nil
}

func (r *paymentTransactionRepository) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query payment.TransactionQuery) ([]*payment.PaymentTransaction, int64, error) {
	query.Validate()

	dbQuery := db.WithContext(ctx).Model(&payment.PaymentTransaction{})

	if tenantID != 0 {
		dbQuery = dbQuery.Where("tenant_id = ?", tenantID.Int64())
	}

	if query.OrderID > 0 {
		dbQuery = dbQuery.Where("order_id = ?", query.OrderID)
	}
	if query.TransactionID != "" {
		dbQuery = dbQuery.Where("transaction_id LIKE ?", escapeLikePattern(query.TransactionID))
	}
	if query.PaymentMethod != "" {
		dbQuery = dbQuery.Where("payment_method = ?", query.PaymentMethod)
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

	var list []*payment.PaymentTransaction
	err := dbQuery.Order("id DESC").
		Offset(query.Offset()).
		Limit(query.Limit()).
		Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *paymentTransactionRepository) GetStats(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) (success, pending, failed int64, err error) {
	type statsResult struct {
		Status int64
		Count  int64
	}
	var results []statsResult

	query := db.WithContext(ctx).
		Model(&payment.PaymentTransaction{}).
		Select("status, COUNT(*) as count").
		Group("status")

	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}

	err = query.Scan(&results).Error
	if err != nil {
		return 0, 0, 0, err
	}

	for _, res := range results {
		switch res.Status {
		case int64(payment.TransactionStatusSucceeded):
			success = res.Count
		case int64(payment.TransactionStatusPending):
			pending = res.Count
		case int64(payment.TransactionStatusFailed):
			failed = res.Count
		}
	}

	return success, pending, failed, nil
}