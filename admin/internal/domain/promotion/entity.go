package promotion

import (
	"context"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

var (
	ErrPromotionNotFound   = errors.New("promotion not found")
	ErrInvalidPromotion    = errors.New("invalid promotion")
	ErrPromotionExpired    = errors.New("promotion expired")
	ErrPromotionNotStarted = errors.New("promotion not started")
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
	ID          int64
	TenantID    shared.TenantID
	Name        string
	Description string
	Type        Type
	Status      Status
	Rules       []PromotionRule
	StartAt     time.Time
	EndAt       time.Time
	Audit       shared.AuditInfo
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

func (p *Promotion) CanApply(cartAmount shared.Money) bool {
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

func (p *Promotion) Apply(cartAmount shared.Money) (shared.Money, error) {
	if !p.IsActive() {
		return cartAmount, ErrPromotionExpired
	}

	bestDiscount := shared.NewMoney(0, cartAmount.Currency)
	for _, rule := range p.Rules {
		if discount, ok := rule.CalculateDiscount(cartAmount); ok {
			if discount.Amount > bestDiscount.Amount {
				bestDiscount = discount
			}
		}
	}

	return cartAmount.Subtract(bestDiscount)
}

type PromotionRule struct {
	ID             int64
	PromotionID    int64
	ConditionType  ConditionType
	ConditionValue int64
	ActionType     ActionType
	ActionValue    int64
	MaxDiscount    shared.Money
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

func (r *PromotionRule) MeetsCondition(cartAmount shared.Money) bool {
	switch r.ConditionType {
	case ConditionMinAmount:
		return cartAmount.Amount >= r.ConditionValue
	default:
		return false
	}
}

func (r *PromotionRule) CalculateDiscount(cartAmount shared.Money) (shared.Money, bool) {
	if !r.MeetsCondition(cartAmount) {
		return shared.Money{}, false
	}

	var discount shared.Money
	switch r.ActionType {
	case ActionFixedAmount:
		discount = shared.NewMoney(r.ActionValue, cartAmount.Currency)
	case ActionPercentage:
		percentage := float64(r.ActionValue) / 100
		discount = cartAmount.MultiplyFloat(percentage)
	}

	if r.MaxDiscount.Amount > 0 && discount.Amount > r.MaxDiscount.Amount {
		discount = r.MaxDiscount
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
