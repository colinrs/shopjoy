package payment

import (
	"context"

	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/tenant"
	"github.com/colinrs/shopjoy/shop/internal/domain/payment"
	"gorm.io/gorm"
)

type Service interface {
	CreatePayment(ctx context.Context, orderID string, userID int64, amount shared.Money, method payment.Method) (*payment.Payment, error)
	ProcessPaymentCallback(ctx context.Context, paymentID, transactionID string, success bool) error
	GetPayment(ctx context.Context, paymentID string) (*payment.Payment, error)
}

type service struct {
	db          *gorm.DB
	paymentRepo payment.Repository
}

func NewService(db *gorm.DB, paymentRepo payment.Repository) Service {
	return &service{
		db:          db,
		paymentRepo: paymentRepo,
	}
}

func (s *service) CreatePayment(ctx context.Context, orderID string, userID int64, amount shared.Money, method payment.Method) (*payment.Payment, error) {
	tenantID, ok := tenant.FromContext(ctx)
	if !ok {
		return nil, shared.ErrInvalidTenantID
	}

	p := &payment.Payment{
		ID:       payment.GeneratePaymentID(tenantID),
		TenantID: tenantID,
		OrderID:  orderID,
		UserID:   userID,
		Amount:   amount,
		Status:   payment.StatusPending,
		Method:   method,
		Audit:    shared.NewAuditInfo(userID),
	}

	if err := s.paymentRepo.Create(ctx, s.db, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *service) ProcessPaymentCallback(ctx context.Context, paymentID, transactionID string, success bool) error {
	tenantID, ok := tenant.FromContext(ctx)
	if !ok {
		return shared.ErrInvalidTenantID
	}

	p, err := s.paymentRepo.FindByID(ctx, s.db, tenantID, paymentID)
	if err != nil {
		return err
	}

	if success {
		if err := p.Pay(transactionID); err != nil {
			return err
		}
	} else {
		p.Fail("payment failed")
	}

	return s.paymentRepo.Update(ctx, s.db, p)
}

func (s *service) GetPayment(ctx context.Context, paymentID string) (*payment.Payment, error) {
	tenantID, ok := tenant.FromContext(ctx)
	if !ok {
		return nil, shared.ErrInvalidTenantID
	}
	return s.paymentRepo.FindByID(ctx, s.db, tenantID, paymentID)
}
