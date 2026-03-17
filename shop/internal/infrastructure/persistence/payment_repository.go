package persistence

import (
	"context"
	"errors"

	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/shop/internal/domain/payment"
	"gorm.io/gorm"
)

type PaymentRepository struct{}

func NewPaymentRepository() payment.Repository {
	return &PaymentRepository{}
}

func (r *PaymentRepository) Create(ctx context.Context, db *gorm.DB, p *payment.Payment) error {
	return db.WithContext(ctx).Create(p).Error
}

func (r *PaymentRepository) Update(ctx context.Context, db *gorm.DB, p *payment.Payment) error {
	return db.WithContext(ctx).Save(p).Error
}

func (r *PaymentRepository) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id string) (*payment.Payment, error) {
	var p payment.Payment
	err := db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID.Int64()).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, payment.ErrPaymentNotFound
	}
	return &p, err
}

func (r *PaymentRepository) FindByOrderID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderID string) (*payment.Payment, error) {
	var p payment.Payment
	err := db.WithContext(ctx).Where("order_id = ? AND tenant_id = ?", orderID, tenantID.Int64()).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, payment.ErrPaymentNotFound
	}
	return &p, err
}

func (r *PaymentRepository) FindByTransactionID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, transactionID string) (*payment.Payment, error) {
	var p payment.Payment
	err := db.WithContext(ctx).Where("transaction_id = ? AND tenant_id = ?", transactionID, tenantID.Int64()).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, payment.ErrPaymentNotFound
	}
	return &p, err
}
