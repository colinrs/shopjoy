package coupon

import (
	"context"
	"fmt"
	"time"

	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Status int

const (
	StatusInactive Status = iota
	StatusActive
	StatusExpired
	StatusDepleted
)

type Type int

const (
	TypeFixedAmount Type = iota
	TypePercentage
	TypeFreeShipping
)

type Coupon struct {
	application.Model
	TenantID     shared.TenantID
	Name         string
	Code         string
	Description  string
	Type         Type
	Value        decimal.Decimal `gorm:"column:value;type:decimal(19,4);not null;default:0"`
	MinAmount    decimal.Decimal `gorm:"column:min_amount;type:decimal(19,4);not null;default:0"`
	MaxDiscount  decimal.Decimal `gorm:"column:max_discount;type:decimal(19,4);not null;default:0"`
	TotalCount   int
	UsedCount    int
	PerUserLimit int
	Status       Status
	StartAt      time.Time
	EndAt        time.Time
	Audit        shared.AuditInfo `gorm:"embedded"`
}

func (c *Coupon) TableName() string {
	return "coupons"
}

func (c *Coupon) GenerateCode() string {
	c.Code = fmt.Sprintf("CPN%s%d", time.Now().Format("20060102"), c.ID)
	return c.Code
}

func (c *Coupon) IsActive() bool {
	if c.Status != StatusActive {
		return false
	}
	if c.TotalCount > 0 && c.UsedCount >= c.TotalCount {
		return false
	}
	now := time.Now().UTC()
	return !now.Before(c.StartAt) && !now.After(c.EndAt)
}

func (c *Coupon) CanUse(userID int64, cartAmount shared.Money) error {
	if !c.IsActive() {
		return code.ErrCouponExpired
	}
	if cartAmount.Amount.LessThan(c.MinAmount) {
		return code.ErrCouponAmountBelowMin
	}
	return nil
}

func (c *Coupon) CalculateDiscount(cartAmount shared.Money) shared.Money {
	var discount shared.Money
	switch c.Type {
	case TypeFixedAmount:
		discount = shared.NewMoney(c.Value, cartAmount.Currency)
	case TypePercentage:
		// c.Value is percentage (e.g., 20.00 = 20%)
		percentage := c.Value.Div(decimal.NewFromInt(100))
		discountDecimal := cartAmount.Amount.Mul(percentage)
		discount = shared.NewMoney(discountDecimal, cartAmount.Currency)
	case TypeFreeShipping:
		discount = shared.NewMoney(decimal.Zero, cartAmount.Currency)
	}

	if c.MaxDiscount.IsPositive() && discount.Amount.GreaterThan(c.MaxDiscount) {
		discount = shared.NewMoney(c.MaxDiscount, cartAmount.Currency)
	}

	if discount.Amount.GreaterThan(cartAmount.Amount) {
		discount = cartAmount
	}

	return discount
}

func (c *Coupon) Use() error {
	if !c.IsActive() {
		return code.ErrCouponExpired
	}
	if c.TotalCount > 0 && c.UsedCount >= c.TotalCount {
		return code.ErrCouponUsedUp
	}
	c.UsedCount++
	return nil
}

type UserCoupon struct {
	application.Model
	TenantID   shared.TenantID
	UserID     int64
	CouponID   int64
	Status     UserCouponStatus
	UsedAt     *time.Time
	OrderID    int64
	ReceivedAt time.Time
	ExpireAt   time.Time
}

type UserCouponStatus int

const (
	UserCouponStatusUnused UserCouponStatus = iota
	UserCouponStatusUsed
	UserCouponStatusExpired
)

func (uc *UserCoupon) TableName() string {
	return "user_coupons"
}

func (uc *UserCoupon) IsExpired() bool {
	return time.Now().After(uc.ExpireAt)
}

func (uc *UserCoupon) CanUse() bool {
	return uc.Status == UserCouponStatusUnused && !uc.IsExpired()
}

type CouponRepository interface {
	Create(ctx context.Context, db *gorm.DB, coupon *Coupon) error
	Update(ctx context.Context, db *gorm.DB, coupon *Coupon) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Coupon, error)
	FindByCode(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, code string) (*Coupon, error)
	FindActive(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) ([]*Coupon, error)
	FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query Query) ([]*Coupon, int64, error)
}

type UserCouponRepository interface {
	Create(ctx context.Context, db *gorm.DB, userCoupon *UserCoupon) error
	FindByUserID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64, status UserCouponStatus) ([]*UserCoupon, error)
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*UserCoupon, error)
	Use(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64, orderID string) error
}

type Query struct {
	shared.PageQuery
	Name   string
	Status Status
	Type   Type
}
