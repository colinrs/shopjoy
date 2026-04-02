package promotion

import (
	"context"
	"sort"

	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ==================== DTOs ====================

type CartItem struct {
	ProductID  int64           `json:"product_id"`
	CategoryID int64           `json:"category_id"`
	BrandID    int64           `json:"brand_id"`
	SKU        string          `json:"sku"`
	Quantity   int             `json:"quantity"`
	UnitPrice  decimal.Decimal `json:"unit_price"`
	LineTotal  decimal.Decimal `json:"line_total"`
}

type CalculateRequest struct {
	TenantID   shared.TenantID `json:"tenant_id"`
	UserID     int64           `json:"user_id"`
	CartItems  []CartItem      `json:"cart_items"`
	Currency   string          `json:"currency"`
	CouponCode string          `json:"coupon_code,omitempty"`
}

type CalculateResult struct {
	OriginalTotal     decimal.Decimal    `json:"original_total"`
	PromotionDiscount decimal.Decimal    `json:"promotion_discount"`
	CouponDiscount    decimal.Decimal    `json:"coupon_discount"`
	FinalTotal        decimal.Decimal    `json:"final_total"`
	AppliedPromotions []AppliedPromotion `json:"applied_promotions"`
	AppliedCoupon     *AppliedCoupon     `json:"applied_coupon,omitempty"`
}

type AppliedPromotion struct {
	PromotionID    int64           `json:"promotion_id"`
	PromotionName  string          `json:"promotion_name"`
	RuleID         int64           `json:"rule_id"`
	DiscountType   string          `json:"discount_type"`
	DiscountAmount decimal.Decimal `json:"discount_amount"`
}

type AppliedCoupon struct {
	CouponID       int64           `json:"coupon_id"`
	CouponName     string          `json:"coupon_name"`
	Code           string          `json:"code"`
	DiscountType   string          `json:"discount_type"`
	DiscountAmount decimal.Decimal `json:"discount_amount"`
}

// ==================== CalculationService ====================

type CalculationService struct {
	promotionRepo  Repository
	couponRepo     CouponRepository
	userCouponRepo UserCouponRepository
}

func NewCalculationService(
	promotionRepo Repository,
	couponRepo CouponRepository,
	userCouponRepo UserCouponRepository,
) *CalculationService {
	return &CalculationService{
		promotionRepo:  promotionRepo,
		couponRepo:     couponRepo,
		userCouponRepo: userCouponRepo,
	}
}

func (s *CalculationService) CalculateDiscount(
	ctx context.Context,
	db *gorm.DB,
	req *CalculateRequest,
) (*CalculateResult, error) {
	var originalTotal decimal.Decimal
	for _, item := range req.CartItems {
		originalTotal = originalTotal.Add(item.LineTotal)
	}

	result := &CalculateResult{
		OriginalTotal:     originalTotal,
		PromotionDiscount: decimal.Zero,
		CouponDiscount:    decimal.Zero,
		FinalTotal:        originalTotal,
		AppliedPromotions: []AppliedPromotion{},
	}

	if originalTotal.IsZero() {
		return result, nil
	}

	// Fetch active promotions filtered by currency
	promotions, err := s.promotionRepo.FindActiveByCurrency(ctx, db, req.TenantID, req.Currency)
	if err != nil {
		return nil, err
	}

	// Sort by priority (descending), then by ID (ascending)
	sort.Slice(promotions, func(i, j int) bool {
		if promotions[i].Priority != promotions[j].Priority {
			return promotions[i].Priority > promotions[j].Priority
		}
		return promotions[i].ID < promotions[j].ID
	})

	// Batch load rules for all promotions (avoid N+1 query)
	promotionIDs := make([]int64, len(promotions))
	for i, p := range promotions {
		promotionIDs[i] = p.ID
	}
	rulesByPromotion, err := s.promotionRepo.FindRulesByPromotionIDs(ctx, db, promotionIDs)
	if err != nil {
		return nil, err
	}
	for i := range promotions {
		promotions[i].Rules = rulesByPromotion[promotions[i].ID]
	}

	// Apply promotions
	for _, promo := range promotions {
		applied := s.applyPromotion(promo, req.CartItems)
		if applied != nil {
			result.AppliedPromotions = append(result.AppliedPromotions, *applied)
			result.PromotionDiscount = result.PromotionDiscount.Add(applied.DiscountAmount)
		}
	}

	// Apply coupon if provided
	if req.CouponCode != "" {
		appliedCoupon, err := s.applyCoupon(ctx, db, req)
		if err != nil {
			return nil, err
		}
		if appliedCoupon != nil {
			result.AppliedCoupon = appliedCoupon
			result.CouponDiscount = appliedCoupon.DiscountAmount
		}
	}

	// Calculate final total
	totalDiscount := result.PromotionDiscount.Add(result.CouponDiscount)
	if totalDiscount.GreaterThan(originalTotal) {
		totalDiscount = originalTotal
	}
	result.FinalTotal = originalTotal.Sub(totalDiscount)

	return result, nil
}

func (s *CalculationService) applyPromotion(
	promo *Promotion,
	items []CartItem,
) *AppliedPromotion {
	var matchedAmount decimal.Decimal
	var matchedQuantity int
	for _, item := range items {
		if promo.MatchesProduct(item.ProductID, item.CategoryID, item.BrandID) {
			matchedAmount = matchedAmount.Add(item.LineTotal)
			matchedQuantity += item.Quantity
		}
	}

	if matchedAmount.IsZero() {
		return nil
	}

	rule := promo.FindBestRule(matchedAmount, matchedQuantity)
	if rule == nil {
		return nil
	}

	discount := rule.CalculateDiscount(matchedAmount)
	if discount.IsZero() {
		return nil
	}

	return &AppliedPromotion{
		PromotionID:    promo.ID,
		PromotionName:  promo.Name,
		RuleID:         rule.ID,
		DiscountType:   rule.ActionType.String(),
		DiscountAmount: discount,
	}
}

func (s *CalculationService) applyCoupon(
	ctx context.Context,
	db *gorm.DB,
	req *CalculateRequest,
) (*AppliedCoupon, error) {
	coupon, err := s.couponRepo.FindByCode(ctx, db, req.TenantID, req.CouponCode)
	if err != nil {
		return nil, err
	}

	if !coupon.IsActive() {
		return nil, nil
	}

	if coupon.Currency != "" && coupon.Currency != req.Currency {
		return nil, nil
	}

	var cartTotal decimal.Decimal
	var scopeMatchedTotal decimal.Decimal
	for _, item := range req.CartItems {
		cartTotal = cartTotal.Add(item.LineTotal)
		if coupon.MatchesProduct(item.ProductID, item.CategoryID, item.BrandID) {
			scopeMatchedTotal = scopeMatchedTotal.Add(item.LineTotal)
		}
	}

	if cartTotal.LessThan(coupon.MinAmount) {
		return nil, nil
	}

	usageCount, err := s.userCouponRepo.CountUsageByUser(ctx, db, req.TenantID, req.UserID, coupon.ID)
	if err != nil {
		return nil, err
	}
	if coupon.PerUserLimit > 0 && usageCount >= coupon.PerUserLimit {
		return nil, nil
	}

	discountBase := cartTotal
	if scopeMatchedTotal.IsPositive() {
		discountBase = scopeMatchedTotal
	}

	discount := coupon.CalculateDiscount(discountBase)
	if discount.IsZero() {
		return nil, nil
	}

	return &AppliedCoupon{
		CouponID:       coupon.ID,
		CouponName:     coupon.Name,
		Code:           coupon.Code,
		DiscountType:   coupon.Type.String(),
		DiscountAmount: discount,
	}, nil
}
