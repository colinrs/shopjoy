package promotion

import (
	"context"

	"gorm.io/gorm"
)

// Repository is the persistence boundary for promotions and their rules.
// All methods are kind-agnostic unless suffixed (e.g., Issue* methods are
// only valid when OwnerKind == KindCoupon).
type Repository interface {
	// Generic CRUD
	Create(ctx context.Context, db *gorm.DB, p *Promotion) error
	Update(ctx context.Context, db *gorm.DB, p *Promotion) error
	Delete(ctx context.Context, db *gorm.DB, id int64) error
	FindByID(ctx context.Context, db *gorm.DB, id int64) (*Promotion, error)
	FindByCode(ctx context.Context, db *gorm.DB, code string) (*Promotion, error)
	FindList(ctx context.Context, db *gorm.DB, query Query) ([]*Promotion, int64, error)

	// Rules
	CreateRules(ctx context.Context, db *gorm.DB, ownerKind Kind, ownerID int64, rules []PromotionRule) error
	FindRulesByOwner(ctx context.Context, db *gorm.DB, ownerKind Kind, ownerID int64) ([]PromotionRule, error)
	// FindRulesByOwners batches the per-owner lookup so list endpoints can
	// hydrate rules in a single query rather than N+1.
	FindRulesByOwners(ctx context.Context, db *gorm.DB, ownerIDs []int64) (map[int64][]PromotionRule, error)
	UpdateRule(ctx context.Context, db *gorm.DB, rule *PromotionRule) error
	DeleteRule(ctx context.Context, db *gorm.DB, id int64) error
	DeleteRulesByOwner(ctx context.Context, db *gorm.DB, ownerKind Kind, ownerID int64) error

	// Coupon-specific
	FindActiveCoupons(ctx context.Context, db *gorm.DB, marketID *int64) ([]*Promotion, error)
	IncrementUsedCount(ctx context.Context, db *gorm.DB, couponID int64) error
	IssueUserCoupon(ctx context.Context, db *gorm.DB, uc *UserCoupon) error
	FindUserCoupons(ctx context.Context, db *gorm.DB, query UserCouponQuery) ([]*UserCoupon, int64, error)

	// Usage (existing)
	FindPromotionUsage(ctx context.Context, db *gorm.DB, query UsageQuery) ([]*PromotionUsage, int64, error)
}
