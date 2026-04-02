package persistence

import (
	"context"
	"errors"

	"github.com/colinrs/shopjoy/admin/internal/domain/payment"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type paymentRepository struct{}

func NewPaymentRepository() payment.PaymentRepository {
	return &paymentRepository{}
}

func (r *paymentRepository) Create(ctx context.Context, db *gorm.DB, p *payment.Payment) error {
	return db.WithContext(ctx).Create(p).Error
}

func (r *paymentRepository) Update(ctx context.Context, db *gorm.DB, p *payment.Payment) error {
	return db.WithContext(ctx).Save(p).Error
}

func (r *paymentRepository) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*payment.Payment, error) {
	var p payment.Payment
	query := db.WithContext(ctx).Where("id = ?", id)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	err := query.First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrPaymentNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepository) FindByOrderID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderID int64) (*payment.Payment, error) {
	var p payment.Payment
	query := db.WithContext(ctx).Where("order_id = ?", orderID)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	err := query.First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrPaymentNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepository) FindByPaymentNo(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, paymentNo string) (*payment.Payment, error) {
	var p payment.Payment
	query := db.WithContext(ctx).Where("payment_no = ?", paymentNo)
	if tenantID != 0 {
		query = query.Where("tenant_id = ?", tenantID.Int64())
	}
	err := query.First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrPaymentNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepository) FindByChannelPaymentID(ctx context.Context, db *gorm.DB, channelPaymentID string) (*payment.Payment, error) {
	var p payment.Payment
	err := db.WithContext(ctx).
		Where("channel_payment_id = ?", channelPaymentID).
		First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrPaymentNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepository) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query payment.PaymentQuery) ([]*payment.Payment, int64, error) {
	query.Validate()

	dbQuery := db.WithContext(ctx).Model(&payment.Payment{})

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

	var list []*payment.Payment
	err := dbQuery.Order("id DESC").
		Offset(query.Offset()).
		Limit(query.Limit()).
		Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
