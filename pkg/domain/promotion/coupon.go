package promotion

import (
	"time"

	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"
)

// UserCoupon is a per-user claim record. coupon_id points to promotions.id
// where kind = 'COUPON'.
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

func (uc *UserCoupon) TableName() string { return "user_coupons" }

func (uc *UserCoupon) IsExpired() bool {
	return time.Now().UTC().After(uc.ExpireAt)
}

func (uc *UserCoupon) CanUse() bool {
	return uc.Status == UserCouponStatusUnused && !uc.IsExpired()
}

// UserCouponStatus
type UserCouponStatus int

const (
	UserCouponStatusUnused  UserCouponStatus = iota // 0
	UserCouponStatusUsed                            // 1
	UserCouponStatusExpired                         // 2
)

// PromotionUsage records one (promotion|rule, coupon, order, user) hit.
// Either promotion_id or coupon_id (or both) may be set.
type PromotionUsage struct {
	ID             int64           `json:"id"`
	TenantID       shared.TenantID `json:"tenant_id"`
	PromotionID    int64           `json:"promotion_id"`
	RuleID         *int64          `json:"rule_id,omitempty"`
	OrderID        int64           `json:"order_id"`
	UserID         int64           `json:"user_id"`
	DiscountAmount decimal.Decimal `json:"discount_amount"`
	Currency       string          `json:"currency"`
	OriginalAmount decimal.Decimal `json:"original_amount"`
	FinalAmount    decimal.Decimal `json:"final_amount"`
	CouponID       *int64          `json:"coupon_id,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
}

func (pu *PromotionUsage) TableName() string { return "promotion_usage" }

// UserCouponQuery is the filter set for Repository.FindUserCoupons.
// All optional fields are pointers so callers can express "no filter"
// without colliding with iota-zero values.
type UserCouponQuery struct {
	Page     int
	Size     int
	UserID   *int64
	CouponID *int64
	Status   *UserCouponStatus
}

// UsageQuery is the filter set for Repository.FindPromotionUsage.
type UsageQuery struct {
	Page     int
	Size     int
	CouponID *int64
	UserID   *int64
}
