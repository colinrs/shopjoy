package promotion

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

// ==================== Promotion Status ====================
// Aligned with existing admin/internal/domain/promotion/entity.go

type Status int

const (
	StatusPending Status = iota // 0 - Draft/Pending state
	StatusActive                // 1 - Active and running
	StatusPaused                // 2 - Manually paused/inactive
	StatusEnded                 // 3 - Ended (time passed)
)

func (s Status) IsValid() bool {
	return s >= StatusPending && s <= StatusEnded
}

func (s Status) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusActive:
		return "active"
	case StatusPaused:
		return "paused"
	case StatusEnded:
		return "ended"
	default:
		return "unknown"
	}
}

// ==================== Promotion Type ====================
// Extended from existing with new MVP types

type Type int

const (
	TypeDiscount Type = iota // 0 - Generic discount (FIXED_DISCOUNT)
	TypeFlashSale            // 1 - Flash sale (Phase 2)
	TypeBundle               // 2 - Bundle promotion (Phase 2)
	TypeBuyXGetY             // 3 - Buy X Get Y (Phase 2)
	TypeFullReduce           // 4 - Tiered full reduction (NEW for MVP)
)

func (t Type) IsValid() bool {
	return t >= TypeDiscount && t <= TypeFullReduce
}

func (t Type) String() string {
	switch t {
	case TypeDiscount:
		return "discount"
	case TypeFlashSale:
		return "flash_sale"
	case TypeBundle:
		return "bundle"
	case TypeBuyXGetY:
		return "buy_x_get_y"
	case TypeFullReduce:
		return "full_reduce"
	default:
		return "unknown"
	}
}

// ==================== Scope Type ====================

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

// ==================== Rule Types ====================
// Aligned with existing admin/internal/domain/promotion/entity.go

type ConditionType int

const (
	ConditionMinAmount ConditionType = iota
	ConditionMinQuantity
)

func (t ConditionType) IsValid() bool {
	return t >= ConditionMinAmount && t <= ConditionMinQuantity
}

type ActionType int

const (
	ActionFixedAmount ActionType = iota
	ActionPercentage
)

func (t ActionType) IsValid() bool {
	return t >= ActionFixedAmount && t <= ActionPercentage
}

func (t ActionType) String() string {
	switch t {
	case ActionFixedAmount:
		return "FIXED_AMOUNT"
	case ActionPercentage:
		return "PERCENTAGE"
	default:
		return "UNKNOWN"
	}
}

// ==================== PromotionScope (Value Object) ====================

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

// ==================== PromotionRule (Entity) ====================

type PromotionRule struct {
	ID             int64         `json:"id"`
	PromotionID    int64         `json:"promotion_id"`
	ConditionType  ConditionType `json:"condition_type"`
	ConditionValue int64         `json:"condition_value"` // Threshold: cents for MIN_AMOUNT, count for MIN_QUANTITY
	ActionType     ActionType    `json:"action_type"`
	ActionValue    int64         `json:"action_value"`    // Discount: cents for FIXED_AMOUNT, basis points for PERCENTAGE (100 = 1%)
	MaxDiscount    int64         `json:"max_discount"`    // Maximum discount cap for percentage (cents)
	Currency       string        `json:"currency"`
	SortOrder      int           `json:"sort_order"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

// CalculateDiscount calculates the discount for a given amount
// ActionValue interpretation:
// - ActionFixedAmount: cents (100 = $1.00)
// - ActionPercentage: basis points (100 = 1%, 500 = 5%, 1000 = 10%)
func (r *PromotionRule) CalculateDiscount(matchedAmount int64) int64 {
	var discount int64

	switch r.ActionType {
	case ActionFixedAmount:
		discount = r.ActionValue
	case ActionPercentage:
		// Basis points: 100 = 1%, so divide by 10000
		discount = matchedAmount * r.ActionValue / 10000
	}

	if r.MaxDiscount > 0 && discount > r.MaxDiscount {
		discount = r.MaxDiscount
	}

	if discount > matchedAmount {
		discount = matchedAmount
	}

	return discount
}

// MeetsCondition checks if the amount/quantity meets the condition
func (r *PromotionRule) MeetsCondition(amount int64, quantity int) bool {
	switch r.ConditionType {
	case ConditionMinAmount:
		return amount >= r.ConditionValue
	case ConditionMinQuantity:
		return int64(quantity) >= r.ConditionValue
	default:
		return false
	}
}

// ==================== Promotion (Aggregate Root) ====================

type Promotion struct {
	ID          int64            `json:"id"`
	TenantID    shared.TenantID  `json:"tenant_id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Type        Type             `json:"type"`
	Status      Status           `json:"status"`
	Priority    int              `json:"priority"`
	StartAt     time.Time        `json:"start_at"`
	EndAt       time.Time        `json:"end_at"`
	Scope       PromotionScope   `json:"scope"`
	Currency    string           `json:"currency"`
	Rules       []PromotionRule  `json:"rules,omitempty"`
	Audit       shared.AuditInfo `json:"audit"`
	DeletedAt   *time.Time       `json:"deleted_at,omitempty"`
}

func (p *Promotion) TableName() string {
	return "promotions"
}

// IsActive checks if promotion is currently active
func (p *Promotion) IsActive() bool {
	if p.Status != StatusActive {
		return false
	}
	if p.DeletedAt != nil {
		return false
	}
	now := time.Now().UTC()
	return !now.Before(p.StartAt) && !now.After(p.EndAt)
}

// MatchesProduct checks if product is in promotion scope
func (p *Promotion) MatchesProduct(productID, categoryID, brandID int64) bool {
	return p.Scope.MatchesProduct(productID, categoryID, brandID)
}

// FindBestRule finds the best applicable rule for given amount
func (p *Promotion) FindBestRule(matchedAmount int64, quantity int) *PromotionRule {
	var bestRule *PromotionRule

	for i := range p.Rules {
		rule := &p.Rules[i]
		if !rule.MeetsCondition(matchedAmount, quantity) {
			continue
		}
		if bestRule == nil || rule.ConditionValue > bestRule.ConditionValue {
			bestRule = rule
		}
	}

	return bestRule
}

// ==================== Query Types ====================

type Query struct {
	shared.PageQuery
	TenantID shared.TenantID
	Name     string
	Status   *Status
	Type     *Type
}

// ==================== Repository Interface ====================
// Extended from existing Repository with new methods

type Repository interface {
	Create(ctx context.Context, db *gorm.DB, promotion *Promotion) error
	Update(ctx context.Context, db *gorm.DB, promotion *Promotion) error
	Delete(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Promotion, error)
	FindActive(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) ([]*Promotion, error)
	FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query Query) ([]*Promotion, int64, error)

	// Extended methods for MVP
	FindActiveByCurrency(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, currency string) ([]*Promotion, error)
	CreateRules(ctx context.Context, db *gorm.DB, rules []PromotionRule) error
	FindRulesByPromotionID(ctx context.Context, db *gorm.DB, promotionID int64) ([]PromotionRule, error)
	FindRulesByPromotionIDs(ctx context.Context, db *gorm.DB, promotionIDs []int64) (map[int64][]PromotionRule, error)
	UpdateRule(ctx context.Context, db *gorm.DB, rule *PromotionRule) error
	DeleteRule(ctx context.Context, db *gorm.DB, id int64) error
}