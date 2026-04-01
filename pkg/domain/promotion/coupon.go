package promotion

import (
	"time"

	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"
)

// ==================== Coupon Type ====================

type CouponType int

const (
	CouponTypeFixedAmount CouponType = iota
	CouponTypePercentage
)

func (t CouponType) IsValid() bool {
	return t >= CouponTypeFixedAmount && t <= CouponTypePercentage
}

func (t CouponType) String() string {
	switch t {
	case CouponTypeFixedAmount:
		return "FIXED_AMOUNT"
	case CouponTypePercentage:
		return "PERCENTAGE"
	default:
		return "UNKNOWN"
	}
}

// ==================== Coupon Status ====================

type CouponStatus int

const (
	CouponStatusInactive CouponStatus = iota
	CouponStatusActive
	CouponStatusExpired
	CouponStatusDepleted
)

func (s CouponStatus) IsValid() bool {
	return s >= CouponStatusInactive && s <= CouponStatusDepleted
}

// ==================== UserCoupon Status ====================

type UserCouponStatus int

const (
	UserCouponStatusUnused UserCouponStatus = iota
	UserCouponStatusUsed
	UserCouponStatusExpired
)

// ==================== Coupon (Aggregate Root) ====================

type Coupon struct {
	ID           int64            `json:"id"`
	TenantID     shared.TenantID  `json:"tenant_id"`
	Name         string           `json:"name"`
	Code         string           `json:"code"`
	Description  string           `json:"description"`
	Type         CouponType       `json:"type"`
	Value        decimal.Decimal  `json:"value"`
	MinAmount    decimal.Decimal  `json:"min_amount"`
	MaxDiscount  decimal.Decimal  `json:"max_discount"`
	Currency     string           `json:"currency"`
	TotalCount   int              `json:"total_count"`
	UsedCount    int              `json:"used_count"`
	PerUserLimit int              `json:"per_user_limit"`
	Status       CouponStatus     `json:"status"`
	StartAt      time.Time        `json:"start_at"`
	EndAt        time.Time        `json:"end_at"`
	Scope        PromotionScope   `json:"scope"`
	Audit        shared.AuditInfo `json:"audit"`
	DeletedAt    *int64           `json:"deleted_at,omitempty"`
}

func (c *Coupon) TableName() string {
	return "coupons"
}

// IsActive checks if coupon is currently active
func (c *Coupon) IsActive() bool {
	if c.Status != CouponStatusActive {
		return false
	}
	if c.DeletedAt != nil {
		return false
	}
	if c.TotalCount > 0 && c.UsedCount >= c.TotalCount {
		return false
	}
	now := time.Now().UTC()
	return !now.Before(c.StartAt) && !now.After(c.EndAt)
}

// CalculateDiscount calculates the discount for a given amount
// Value interpretation:
// - CouponTypeFixedAmount: amount (e.g., 100 = 100.00)
// - CouponTypePercentage: basis points (100 = 1%, 500 = 5%, 1000 = 10%)
func (c *Coupon) CalculateDiscount(cartAmount decimal.Decimal) decimal.Decimal {
	var discount decimal.Decimal

	switch c.Type {
	case CouponTypeFixedAmount:
		discount = c.Value
	case CouponTypePercentage:
		// Basis points: 100 = 1%, so divide by 10000
		discount = cartAmount.Mul(c.Value).Div(decimal.NewFromInt(10000))
	}

	if c.MaxDiscount.IsPositive() && discount.GreaterThan(c.MaxDiscount) {
		discount = c.MaxDiscount
	}

	if discount.GreaterThan(cartAmount) {
		discount = cartAmount
	}

	return discount
}

// MatchesProduct checks if product is in coupon scope
func (c *Coupon) MatchesProduct(productID, categoryID, brandID int64) bool {
	return c.Scope.MatchesProduct(productID, categoryID, brandID)
}

// ==================== UserCoupon (Entity) ====================

type UserCoupon struct {
	ID         int64            `json:"id"`
	TenantID   shared.TenantID  `json:"tenant_id"`
	UserID     int64            `json:"user_id"`
	CouponID   int64            `json:"coupon_id"`
	Status     UserCouponStatus `json:"status"`
	UsedAt     *time.Time       `json:"used_at,omitempty"`
	OrderID    int64            `json:"order_id"`
	ReceivedAt time.Time        `json:"received_at"`
	ExpireAt   time.Time        `json:"expire_at"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

func (uc *UserCoupon) TableName() string {
	return "user_coupons"
}

// IsExpired checks if user coupon is expired
func (uc *UserCoupon) IsExpired() bool {
	return time.Now().UTC().After(uc.ExpireAt)
}

// CanUse checks if user coupon can be used
func (uc *UserCoupon) CanUse() bool {
	return uc.Status == UserCouponStatusUnused && !uc.IsExpired()
}

// ==================== PromotionUsage (Entity) ====================

type PromotionUsage struct {
	ID             int64            `json:"id"`
	TenantID       shared.TenantID  `json:"tenant_id"`
	PromotionID    int64            `json:"promotion_id"`
	RuleID         *int64           `json:"rule_id,omitempty"`
	OrderID        int64            `json:"order_id"`
	UserID         int64            `json:"user_id"`
	DiscountAmount decimal.Decimal  `json:"discount_amount"`
	Currency       string           `json:"currency"`
	OriginalAmount decimal.Decimal  `json:"original_amount"`
	FinalAmount    decimal.Decimal  `json:"final_amount"`
	CouponID       *int64           `json:"coupon_id,omitempty"`
	CreatedAt      time.Time        `json:"created_at"`
}

func (pu *PromotionUsage) TableName() string {
	return "promotion_usage"
}

// ==================== Query Types ====================

type CouponQuery struct {
	shared.PageQuery
	TenantID shared.TenantID
	Name     string
	Code     string
	Status   *CouponStatus
	Type     *CouponType
}
