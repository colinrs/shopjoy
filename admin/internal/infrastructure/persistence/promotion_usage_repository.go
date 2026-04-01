package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/promotion"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type promotionUsageRepo struct{}

func NewPromotionUsageRepository() promotion.PromotionUsageRepository {
	return &promotionUsageRepo{}
}

// promotionUsageModel represents the database model for PromotionUsage
type promotionUsageModel struct {
	ID             int64           `gorm:"column:id;primaryKey;autoIncrement:false"`
	TenantID       int64           `gorm:"column:tenant_id;not null;index"`
	PromotionID    int64           `gorm:"column:promotion_id;not null;index"`
	RuleID         *int64          `gorm:"column:rule_id;index"`
	OrderID        int64           `gorm:"column:order_id;not null;index"`
	UserID         int64           `gorm:"column:user_id;not null;index"`
	DiscountAmount decimal.Decimal `gorm:"column:discount_amount;type:decimal(19,4);not null"`
	Currency       string          `gorm:"column:currency;size:10;not null"`
	OriginalAmount decimal.Decimal `gorm:"column:original_amount;type:decimal(19,4);not null"`
	FinalAmount    decimal.Decimal `gorm:"column:final_amount;type:decimal(19,4);not null"`
	CouponID       *int64           `gorm:"column:coupon_id;index"`
	CreatedAt      int64            `gorm:"column:created_at"`
}

func (promotionUsageModel) TableName() string {
	return "promotion_usage"
}

func (m *promotionUsageModel) toEntity() *promotion.PromotionUsage {
	return &promotion.PromotionUsage{
		ID:             m.ID,
		TenantID:       shared.TenantID(m.TenantID),
		PromotionID:    m.PromotionID,
		RuleID:         m.RuleID,
		OrderID:        m.OrderID,
		UserID:         m.UserID,
		DiscountAmount: m.DiscountAmount,
		Currency:       m.Currency,
		OriginalAmount: m.OriginalAmount,
		FinalAmount:    m.FinalAmount,
		CouponID:       m.CouponID,
		CreatedAt:      time.Unix(m.CreatedAt, 0),
	}
}

func fromPromotionUsageEntity(pu *promotion.PromotionUsage) *promotionUsageModel {
	return &promotionUsageModel{
		ID:             pu.ID,
		TenantID:       pu.TenantID.Int64(),
		PromotionID:    pu.PromotionID,
		RuleID:         pu.RuleID,
		OrderID:        pu.OrderID,
		UserID:         pu.UserID,
		DiscountAmount: pu.DiscountAmount,
		Currency:       pu.Currency,
		OriginalAmount: pu.OriginalAmount,
		FinalAmount:    pu.FinalAmount,
		CouponID:       pu.CouponID,
		CreatedAt:      pu.CreatedAt.Unix(),
	}
}

// Create inserts a new promotion usage record
func (r *promotionUsageRepo) Create(ctx context.Context, db *gorm.DB, usage *promotion.PromotionUsage) error {
	model := fromPromotionUsageEntity(usage)
	return db.WithContext(ctx).Create(model).Error
}

// FindByOrderID finds a promotion usage record by order ID
func (r *promotionUsageRepo) FindByOrderID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, orderID int64) (*promotion.PromotionUsage, error) {
	var model promotionUsageModel
	err := db.WithContext(ctx).
		Where("tenant_id = ?", tenantID.Int64()).
		Where("order_id = ?", orderID).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

// FindList finds promotion usage records with pagination and filters
func (r *promotionUsageRepo) FindList(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, query promotion.PromotionUsageQuery) ([]*promotion.PromotionUsage, int64, error) {
	query.Validate()

	dbQuery := db.WithContext(ctx).Model(&promotionUsageModel{}).
		Where("tenant_id = ?", tenantID.Int64())

	if query.PromotionID != nil {
		dbQuery = dbQuery.Where("promotion_id = ?", *query.PromotionID)
	}
	if query.CouponID != nil {
		dbQuery = dbQuery.Where("coupon_id = ?", *query.CouponID)
	}
	if query.UserID != nil {
		dbQuery = dbQuery.Where("user_id = ?", *query.UserID)
	}
	if query.OrderID != 0 {
		dbQuery = dbQuery.Where("order_id = ?", query.OrderID)
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []promotionUsageModel
	err := dbQuery.Order("created_at DESC").
		Offset(query.Offset()).
		Limit(query.Limit()).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	usages := make([]*promotion.PromotionUsage, len(models))
	for i, m := range models {
		usages[i] = m.toEntity()
	}
	return usages, total, nil
}