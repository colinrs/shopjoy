package promotion

import (
	"time"

	"github.com/shopspring/decimal"
)

// ConditionType enumerates the types of rule conditions.
type ConditionType int

const (
	ConditionMinAmount ConditionType = iota // 0 - minimum cart amount
	ConditionMinQuantity                     // 1 - minimum item quantity
)

func (t ConditionType) IsValid() bool {
	return t >= ConditionMinAmount && t <= ConditionMinQuantity
}

// ActionType enumerates the discount actions a rule can perform.
type ActionType int

const (
	ActionFixedAmount ActionType = iota // 0 - subtract ActionValue as money
	ActionPercentage                    // 1 - multiply by ActionValue/10000 (basis points)
	ActionFreeShipping                   // 2 - waive shipping (ActionValue unused)
)

func (t ActionType) IsValid() bool {
	return t >= ActionFixedAmount && t <= ActionFreeShipping
}

// PromotionRule belongs to a Promotion (of any Kind) via owner_kind+owner_id.
// Rules are ordered by sort_order; the best-matching rule wins via FindBestRule.
type PromotionRule struct {
	ID             int64           `json:"id"`
	OwnerKind      Kind            `json:"owner_kind"`
	OwnerID        int64           `json:"owner_id"`
	ConditionType  ConditionType   `json:"condition_type"`
	ConditionValue decimal.Decimal `json:"condition_value"`
	ActionType     ActionType      `json:"action_type"`
	ActionValue    decimal.Decimal `json:"action_value"`
	MaxDiscount    decimal.Decimal `json:"max_discount"`
	Currency       string          `json:"currency"`
	SortOrder      int             `json:"sort_order"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// CalculateDiscount applies the rule's action to matchedAmount and caps by MaxDiscount.
// ActionPercentage: ActionValue is basis points (100 = 1%, so divide by 10000).
// ActionFreeShipping: returns zero (the discount is the shipping fee, handled elsewhere).
func (r *PromotionRule) CalculateDiscount(matchedAmount decimal.Decimal) decimal.Decimal {
	var discount decimal.Decimal
	switch r.ActionType {
	case ActionFixedAmount:
		discount = r.ActionValue
	case ActionPercentage:
		discount = matchedAmount.Mul(r.ActionValue).Div(decimal.NewFromInt(10000))
	case ActionFreeShipping:
		return decimal.Zero
	}

	if r.MaxDiscount.IsPositive() && discount.GreaterThan(r.MaxDiscount) {
		discount = r.MaxDiscount
	}
	if discount.GreaterThan(matchedAmount) {
		discount = matchedAmount
	}
	return discount
}

// MeetsCondition returns true if amount/quantity clears the threshold.
func (r *PromotionRule) MeetsCondition(amount decimal.Decimal, quantity int) bool {
	switch r.ConditionType {
	case ConditionMinAmount:
		return amount.GreaterThanOrEqual(r.ConditionValue)
	case ConditionMinQuantity:
		return decimal.NewFromInt(int64(quantity)).GreaterThanOrEqual(r.ConditionValue)
	default:
		return false
	}
}
