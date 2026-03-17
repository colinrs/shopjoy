package payment

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

var (
	ErrPaymentNotFound = errors.New("payment not found")
	ErrInvalidAmount   = errors.New("invalid payment amount")
	ErrPaymentFailed   = errors.New("payment failed")
	ErrAlreadyPaid     = errors.New("payment already completed")
	ErrPaymentExpired  = errors.New("payment expired")
)

type Status int

const (
	StatusPending Status = iota
	StatusProcessing
	StatusSuccess
	StatusFailed
	StatusCancelled
	StatusRefunded
)

type Method int

const (
	MethodAlipay Method = iota
	MethodWechat
	MethodCreditCard
	MethodBankTransfer
	MethodCash
)

type Payment struct {
	ID            string
	TenantID      shared.TenantID
	OrderID       string
	UserID        int64
	Amount        shared.Money
	Status        Status
	Method        Method
	TransactionID string
	PaidAt        *time.Time
	ExpireAt      time.Time
	NotifyURL     string
	ReturnURL     string
	Audit         shared.AuditInfo
}

func (p *Payment) TableName() string {
	return "payments"
}

func (p *Payment) IsExpired() bool {
	return time.Now().After(p.ExpireAt)
}

func (p *Payment) Pay(transactionID string) error {
	if p.Status == StatusSuccess {
		return ErrAlreadyPaid
	}
	if p.Status == StatusFailed {
		return ErrPaymentFailed
	}
	if p.IsExpired() {
		return ErrPaymentExpired
	}

	now := time.Now().UTC()
	p.Status = StatusSuccess
	p.TransactionID = transactionID
	p.PaidAt = &now
	return nil
}

func (p *Payment) Fail(reason string) {
	p.Status = StatusFailed
}

func (p *Payment) Cancel() {
	if p.Status == StatusPending || p.Status == StatusProcessing {
		p.Status = StatusCancelled
	}
}

type Refund struct {
	ID            string
	TenantID      shared.TenantID
	PaymentID     string
	OrderID       string
	UserID        int64
	Amount        shared.Money
	Reason        string
	Status        RefundStatus
	TransactionID string
	RefundedAt    *time.Time
	Audit         shared.AuditInfo
}

type RefundStatus int

const (
	RefundStatusPending RefundStatus = iota
	RefundStatusProcessing
	RefundStatusCompleted
	RefundStatusFailed
)

func (r *Refund) TableName() string {
	return "payment_refunds"
}

func (r *Refund) Complete(transactionID string) {
	now := time.Now().UTC()
	r.Status = RefundStatusCompleted
	r.TransactionID = transactionID
	r.RefundedAt = &now
}

type Repository interface {
	Create(ctx context.Context, db *gorm.DB, payment *Payment) error
	Update(ctx context.Context, db *gorm.DB, payment *Payment) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id string) (*Payment, error)
	FindByOrderID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderID string) (*Payment, error)
	FindByTransactionID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, transactionID string) (*Payment, error)
}

type RefundRepository interface {
	Create(ctx context.Context, db *gorm.DB, refund *Refund) error
	Update(ctx context.Context, db *gorm.DB, refund *Refund) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id string) (*Refund, error)
	FindByPaymentID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, paymentID string) ([]*Refund, error)
}

func GeneratePaymentID(tenantID shared.TenantID) string {
	return fmt.Sprintf("PAY%s%d%d", time.Now().Format("20060102"), tenantID.Int64(), time.Now().UnixNano()%1000000)
}
