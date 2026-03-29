package promotion

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Status int

const (
	StatusPending Status = iota
	StatusActive
	StatusPaused
	StatusEnded
)

type Type int

const (
	TypeDiscount Type = iota
	TypeFlashSale
	TypeBundle
	TypeBuyXGetY
)

type Promotion struct {
	application.Model
	TenantID    shared.TenantID
	Name        string
	Description string
	Type        Type
	Status      Status
	Rules       []PromotionRule
	StartAt     time.Time
	EndAt       time.Time
	Audit       shared.AuditInfo `gorm:"embedded"`
}

func (p *Promotion) TableName() string {
	return "promotions"
}

func (p *Promotion) IsActive() bool {
	if p.Status != StatusActive {
		return false
	}
	now := time.Now().UTC()
	return !now.Before(p.StartAt) && !now.After(p.EndAt)
}

func (p *Promotion) CanApply(cartAmount decimal.Decimal) bool {
	if !p.IsActive() {
		return false
	}
	for _, rule := range p.Rules {
		if rule.MeetsCondition(cartAmount) {
			return true
		}
	}
	return false
}

func (p *Promotion) Apply(cartAmount decimal.Decimal) (decimal.Decimal, error) {
	if !p.IsActive() {
		return cartAmount, code.ErrPromotionExpired
	}

	bestDiscount := decimal.Zero
	for _, rule := range p.Rules {
		if discount, ok := rule.CalculateDiscount(cartAmount); ok {
			if discount.GreaterThan(bestDiscount) {
				bestDiscount = discount
			}
		}
	}

	return cartAmount.Sub(bestDiscount), nil
}

type PromotionRule struct {
	application.Model
	PromotionID    int64
	ConditionType  ConditionType
	ConditionValue decimal.Decimal
	ActionType     ActionType
	ActionValue    decimal.Decimal
	MaxDiscount    decimal.Decimal
}

type ConditionType int

const (
	ConditionMinAmount ConditionType = iota
	ConditionMinQuantity
)

type ActionType int

const (
	ActionFixedAmount ActionType = iota
	ActionPercentage
)

func (r *PromotionRule) MeetsCondition(cartAmount decimal.Decimal) bool {
	switch r.ConditionType {
	case ConditionMinAmount:
		return cartAmount.GreaterThanOrEqual(r.ConditionValue)
	default:
		return false
	}
}

func (r *PromotionRule) CalculateDiscount(cartAmount decimal.Decimal) (decimal.Decimal, bool) {
	if !r.MeetsCondition(cartAmount) {
		return decimal.Zero, false
	}

	var discount decimal.Decimal
	switch r.ActionType {
	case ActionFixedAmount:
		discount = r.ActionValue
	case ActionPercentage:
		// ActionValue is basis points (100 = 1%), so divide by 10000
		discount = cartAmount.Mul(r.ActionValue).Div(decimal.NewFromInt(10000))
	}

	if r.MaxDiscount.IsPositive() && discount.GreaterThan(r.MaxDiscount) {
		discount = r.MaxDiscount
	}

	if discount.GreaterThan(cartAmount) {
		discount = cartAmount
	}

	return discount, true
}

type Repository interface {
	Create(ctx context.Context, db *gorm.DB, promotion *Promotion) error
	Update(ctx context.Context, db *gorm.DB, promotion *Promotion) error
	FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*Promotion, error)
	FindActive(ctx context.Context, db *gorm.DB, tenantID shared.TenantID) ([]*Promotion, error)
	FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query Query) ([]*Promotion, int64, error)
}

type Query struct {
	shared.PageQuery
	Name   string
	Status Status
	Type   Type
}
