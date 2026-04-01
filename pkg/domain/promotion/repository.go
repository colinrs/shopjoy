package promotion

import (
	"context"

	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

// CouponRepository defines the interface for coupon persistence
type CouponRepository interface {
	Create(ctx context.Context, db *gorm.DB, coupon *Coupon) error
	Update(ctx context.Context, db *gorm.DB, coupon *Coupon) error
	Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Coupon, error)
	FindByCode(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, code string) (*Coupon, error)
	FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query CouponQuery) ([]*Coupon, int64, error)
	IncrementUsage(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
}

// UserCouponRepository defines the interface for user coupon persistence
type UserCouponRepository interface {
	Create(ctx context.Context, db *gorm.DB, userCoupon *UserCoupon) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*UserCoupon, error)
	FindByUserID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64, status *UserCouponStatus) ([]*UserCoupon, error)
	FindByUserAndCoupon(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64, couponID int64) ([]*UserCoupon, error)
	MarkUsed(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64, orderID int64) error
	CountUsageByUser(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, userID int64, couponID int64) (int, error)
}

// PromotionUsageRepository defines the interface for promotion usage persistence
type PromotionUsageRepository interface {
	Create(ctx context.Context, db *gorm.DB, usage *PromotionUsage) error
	FindByOrderID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderID int64) (*PromotionUsage, error)
	FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query PromotionUsageQuery) ([]*PromotionUsage, int64, error)
}

// PromotionUsageQuery for querying promotion usage
type PromotionUsageQuery struct {
	shared.PageQuery
	TenantID    shared.TenantID
	PromotionID *int64
	CouponID    *int64
	UserID      *int64
	OrderID     int64
}