package promotion

import (
	"time"

	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"
)

// Kind discriminates between system-driven promotions and claim-based coupons.
type Kind string

const (
	KindPromotion Kind = "PROMOTION"
	KindCoupon    Kind = "COUPON"
)

func (k Kind) IsValid() bool {
	return k == KindPromotion || k == KindCoupon
}

// Status (unified for both kinds)
type Status int

const (
	StatusPending Status = iota // 0
	StatusActive                // 1
	StatusPaused                // 2
	StatusEnded                 // 3 - depleted coupons also surface as ended (see IsActive)
)

func (s Status) IsValid() bool {
	return s >= StatusPending && s <= StatusEnded
}

// Type (marketing play). COUPONs always use TypeDiscount (=0).
type Type int

const (
	TypeDiscount  Type = iota // 0
	TypeFlashSale             // 1
	TypeBundle                // 2
	TypeBuyXGetY              // 3
)

func (t Type) IsValid() bool {
	return t >= TypeDiscount && t <= TypeBuyXGetY
}

// ScopeType enumerates the kinds of product scope a promotion can target.
// (MARKET scope was removed — promotions now use the top-level market_id column.)
type ScopeType string

const (
	ScopeTypeStorewide  ScopeType = "STOREWIDE"
	ScopeTypeProducts   ScopeType = "PRODUCTS"
	ScopeTypeCategories ScopeType = "CATEGORIES"
	ScopeTypeBrands     ScopeType = "BRANDS"
)

func (t ScopeType) IsValid() bool {
	switch t {
	case ScopeTypeStorewide, ScopeTypeProducts, ScopeTypeCategories, ScopeTypeBrands:
		return true
	default:
		return false
	}
}

// PromotionScope (Value Object)
type PromotionScope struct {
	Type       ScopeType `json:"type"`
	IDs        []int64   `json:"ids,omitempty"`
	ExcludeIDs []int64   `json:"exclude_ids,omitempty"`
}

// MatchesProduct checks if a product matches the scope
func (s *PromotionScope) MatchesProduct(productID, categoryID, brandID int64) bool {
	for _, excludeID := range s.ExcludeIDs {
		if excludeID == productID {
			return false
		}
	}

	switch s.Type {
	case ScopeTypeStorewide:
		return true
	case ScopeTypeProducts:
		for _, id := range s.IDs {
			if id == productID {
				return true
			}
		}
		return false
	case ScopeTypeCategories:
		for _, id := range s.IDs {
			if id == categoryID {
				return true
			}
		}
		return false
	case ScopeTypeBrands:
		for _, id := range s.IDs {
			if id == brandID {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// Promotion is the aggregate root for both system promotions and user-claimable
// coupons. Coupon-specific fields are nullable; semantics activate when Kind == KindCoupon.
type Promotion struct {
	ID           int64            `json:"id"`
	TenantID     shared.TenantID  `json:"tenant_id"`
	Kind         Kind             `json:"kind"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	Code         *string          `json:"code,omitempty"`
	Type         Type             `json:"type"`
	Status       Status           `json:"status"`
	Priority     int              `json:"priority"`
	MarketID     *int64           `json:"market_id,omitempty"`
	Currency     string           `json:"currency"`
	TotalCount   *int             `json:"total_count,omitempty"`
	UsedCount    *int             `json:"used_count,omitempty"`
	UsageLimit   int              `json:"usage_limit"`
	PerUserLimit int              `json:"per_user_limit"`
	Tags         []string         `json:"tags,omitempty" gorm:"type:json"`
	Scope        PromotionScope   `json:"scope"`
	StartAt      time.Time        `json:"start_at"`
	EndAt        time.Time        `json:"end_at"`
	Rules        []PromotionRule  `json:"rules,omitempty"`
	Audit        shared.AuditInfo `json:"audit"`
	DeletedAt    *time.Time       `json:"deleted_at,omitempty"`
}

func (p *Promotion) TableName() string { return "promotions" }

// IsActive returns true if the promotion is currently usable. For COUPONs,
// this also checks inventory (used_count < total_count).
func (p *Promotion) IsActive() bool {
	if p.Status != StatusActive || p.DeletedAt != nil {
		return false
	}
	if p.Kind == KindCoupon && p.TotalCount != nil && p.UsedCount != nil {
		if *p.UsedCount >= *p.TotalCount {
			return false // depleted
		}
	}
	now := time.Now().UTC()
	return !now.Before(p.StartAt) && !now.After(p.EndAt)
}

// MatchesMarket returns true if the promotion applies to the given market.
// A NULL market_id means "applies to all markets".
func (p *Promotion) MatchesMarket(marketID int64) bool {
	return p.MarketID == nil || *p.MarketID == marketID
}

// MatchesScope delegates to Scope (kind-agnostic).
func (p *Promotion) MatchesScope(productID, categoryID, brandID int64) bool {
	return p.Scope.MatchesProduct(productID, categoryID, brandID)
}

// FindBestRule returns the rule with the highest ConditionValue that still
// meets the condition. Multi-tier support for COUPONs lives here.
func (p *Promotion) FindBestRule(matchedAmount decimal.Decimal, quantity int) *PromotionRule {
	var best *PromotionRule
	for i := range p.Rules {
		rule := &p.Rules[i]
		if !rule.MeetsCondition(matchedAmount, quantity) {
			continue
		}
		if best == nil || rule.ConditionValue.GreaterThan(best.ConditionValue) {
			best = rule
		}
	}
	return best
}

// CalculateDiscount is a convenience wrapper around FindBestRule.
func (p *Promotion) CalculateDiscount(matchedAmount decimal.Decimal, quantity int) decimal.Decimal {
	rule := p.FindBestRule(matchedAmount, quantity)
	if rule == nil {
		return decimal.Zero
	}
	return rule.CalculateDiscount(matchedAmount)
}

// Issue creates a UserCoupon from this promotion. Returns ErrPromotionInvalidKind
// if invoked on a non-COUPON.
func (p *Promotion) Issue(userID int64, now time.Time) (*UserCoupon, error) {
	if p.Kind != KindCoupon {
		return nil, code.ErrPromotionInvalidKind
	}
	if !p.IsActive() {
		return nil, code.ErrCouponExpired
	}
	return &UserCoupon{
		TenantID:   p.TenantID,
		UserID:     userID,
		CouponID:   p.ID,
		Status:     UserCouponStatusUnused,
		ReceivedAt: now,
		ExpireAt:   p.EndAt,
	}, nil
}

// ConsumeInventory increments used_count in memory (pre-check). Persistence
// happens via repo.IncrementUsedCount which uses an atomic SQL check to
// prevent overselling:
//
//	UPDATE promotions
//	SET used_count = used_count + 1
//	WHERE id = ? AND kind = 'COUPON' AND (total_count IS NULL OR used_count < total_count)
//
// Caller must roll back the in-memory value if the SQL fails.
func (p *Promotion) ConsumeInventory() error {
	if p.Kind != KindCoupon || p.UsedCount == nil || p.TotalCount == nil {
		return code.ErrPromotionInvalidKind
	}
	if *p.UsedCount >= *p.TotalCount {
		return code.ErrCouponUsedUp
	}
	*p.UsedCount++
	return nil
}
