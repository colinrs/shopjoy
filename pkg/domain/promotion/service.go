package promotion

import (
	"context"
	"sort"

	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

// ==================== DTOs ====================

type CartItem struct {
	ProductID  int64 `json:"product_id"`
	CategoryID int64 `json:"category_id"`
	BrandID    int64 `json:"brand_id"`
	SKU        string `json:"sku"`
	Quantity   int    `json:"quantity"`
	UnitPrice  int64  `json:"unit_price"`
	LineTotal  int64  `json:"line_total"`
}

type CalculateRequest struct {
	TenantID   shared.TenantID `json:"tenant_id"`
	UserID     int64           `json:"user_id"`
	CartItems  []CartItem      `json:"cart_items"`
	Currency   string          `json:"currency"`
	CouponCode string          `json:"coupon_code,omitempty"`
}

type CalculateResult struct {
	OriginalTotal     int64              `json:"original_total"`
	PromotionDiscount int64              `json:"promotion_discount"`
	CouponDiscount    int64              `json:"coupon_discount"`
	FinalTotal        int64              `json:"final_total"`
	AppliedPromotions []AppliedPromotion `json:"applied_promotions"`
	AppliedCoupon     *AppliedCoupon     `json:"applied_coupon,omitempty"`
}

type AppliedPromotion struct {
	PromotionID    int64  `json:"promotion_id"`
	PromotionName  string `json:"promotion_name"`
	RuleID         int64  `json:"rule_id"`
	DiscountType   string `json:"discount_type"`
	DiscountAmount int64  `json:"discount_amount"`
}

type AppliedCoupon struct {
	CouponID       int64  `json:"coupon_id"`
	CouponName     string `json:"coupon_name"`
	Code           string `json:"code"`
	DiscountType   string `json:"discount_type"`
	DiscountAmount int64  `json:"discount_amount"`
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
	var originalTotal int64
	for _, item := range req.CartItems {
		originalTotal += item.LineTotal
	}

	result := &CalculateResult{
		OriginalTotal:     originalTotal,
		PromotionDiscount: 0,
		CouponDiscount:    0,
		FinalTotal:        originalTotal,
		AppliedPromotions: []AppliedPromotion{},
	}

	if originalTotal == 0 {
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
			result.PromotionDiscount += applied.DiscountAmount
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
	totalDiscount := result.PromotionDiscount + result.CouponDiscount
	if totalDiscount > originalTotal {
		totalDiscount = originalTotal
	}
	result.FinalTotal = originalTotal - totalDiscount

	return result, nil
}

func (s *CalculationService) applyPromotion(
	promo *Promotion,
	items []CartItem,
) *AppliedPromotion {
	var matchedAmount int64
	var matchedQuantity int
	for _, item := range items {
		if promo.MatchesProduct(item.ProductID, item.CategoryID, item.BrandID) {
			matchedAmount += item.LineTotal
			matchedQuantity += item.Quantity
		}
	}

	if matchedAmount == 0 {
		return nil
	}

	rule := promo.FindBestRule(matchedAmount, matchedQuantity)
	if rule == nil {
		return nil
	}

	discount := rule.CalculateDiscount(matchedAmount)
	if discount == 0 {
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

	var cartTotal int64
	var scopeMatchedTotal int64
	for _, item := range req.CartItems {
		cartTotal += item.LineTotal
		if coupon.MatchesProduct(item.ProductID, item.CategoryID, item.BrandID) {
			scopeMatchedTotal += item.LineTotal
		}
	}

	if cartTotal < coupon.MinAmount {
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
	if scopeMatchedTotal > 0 {
		discountBase = scopeMatchedTotal
	}

	discount := coupon.CalculateDiscount(discountBase)
	if discount == 0 {
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