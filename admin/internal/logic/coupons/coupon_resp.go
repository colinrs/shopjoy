package coupons

import (
	"context"
	"strings"
	"time"

	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
	"github.com/colinrs/shopjoy/admin/internal/types"
	pkgpromotion "github.com/colinrs/shopjoy/pkg/domain/promotion"
	"github.com/shopspring/decimal"
)

// mapCouponActionType maps the wire coupon Type ("fixed_amount",
// "percentage", "free_shipping") onto the domain ActionType.
func mapCouponActionType(t string) pkgpromotion.ActionType {
	switch strings.ToLower(t) {
	case "percentage":
		return pkgpromotion.ActionPercentage
	case "free_shipping":
		return pkgpromotion.ActionFreeShipping
	default:
		return pkgpromotion.ActionFixedAmount
	}
}

// buildCouponScope maps the wire-level scope hint onto the domain's
// PromotionScope. The form does not currently bind IDs for coupons
// (product_ids / category_ids are JSON strings the form never
// re-populates), so the resulting PromotionScope carries the chosen
// scope_type with a nil ID slice.
//
// Wire values are lowercase ("products", "categories", …); the
// domain constants are uppercase ("PRODUCTS", …). Normalize via
// strings.ToUpper before comparing to the enum.
func buildCouponScope(scopeType string) pkgpromotion.PromotionScope {
	switch strings.ToUpper(scopeType) {
	case string(pkgpromotion.ScopeTypeProducts):
		return pkgpromotion.PromotionScope{Type: pkgpromotion.ScopeTypeProducts}
	case string(pkgpromotion.ScopeTypeCategories):
		return pkgpromotion.PromotionScope{Type: pkgpromotion.ScopeTypeCategories}
	case string(pkgpromotion.ScopeTypeBrands):
		return pkgpromotion.PromotionScope{Type: pkgpromotion.ScopeTypeBrands}
	}
	return pkgpromotion.PromotionScope{Type: pkgpromotion.ScopeTypeStorewide}
}

// =============================================================================
// Money / decimal helpers
// =============================================================================

// parseMoney parses a decimal-formatted money string into a
// decimal.Decimal, returning zero for empty / invalid input.
func parseMoney(s string) decimal.Decimal {
	if s == "" {
		return decimal.Zero
	}
	v, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero
	}
	return v
}

// formatMoney formats a decimal with two fraction digits, returning
// "0" for the zero value (matches the pre-merge wire convention).
func formatMoney(v decimal.Decimal) string {
	if v.IsZero() {
		return "0"
	}
	return v.StringFixed(2)
}

// defaultCurrency returns the wire-supplied currency or "CNY" if
// the frontend omitted it.
func defaultCurrency(c string) string {
	if c == "" {
		return "CNY"
	}
	return c
}

// =============================================================================
// Context helpers
// =============================================================================

// tenantID extracts the tenant ID from context with a 0 fallback.
// (The admin API currently runs single-tenant; the 0 default lets
// the repo filter gracefully when the value is absent.)
func tenantID(ctx context.Context) int64 {
	if v := ctx.Value("tenant_id"); v != nil {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	return 0
}

// actorID is the audit-user helper. Missing values fall back to 0
// which the repo treats as "system".
func actorID(ctx context.Context) int64 {
	if v := ctx.Value("user_id"); v != nil {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	return 0
}

// =============================================================================
// Status / type mapping (coupon-specific enum values)
// =============================================================================

func mapCouponStatusToWire(s pkgpromotion.Status) string {
	switch s {
	case pkgpromotion.StatusPending:
		return "inactive"
	case pkgpromotion.StatusActive:
		return "active"
	case pkgpromotion.StatusPaused:
		return "inactive"
	case pkgpromotion.StatusEnded:
		return "depleted"
	default:
		return "inactive"
	}
}

// mapCouponTypeToWire converts the domain ActionType (which carries
// the discount play) back to the wire coupon Type string
// ("fixed_amount" / "percentage" / "free_shipping").
func mapCouponTypeToWire(t pkgpromotion.ActionType) string {
	switch t {
	case pkgpromotion.ActionPercentage:
		return "percentage"
	case pkgpromotion.ActionFreeShipping:
		return "free_shipping"
	default:
		return "fixed_amount"
	}
}

// isExpiredRFC3339 returns true if the given RFC3339 timestamp is
// already in the past.
func isExpiredRFC3339(s string) bool {
	if s == "" {
		return false
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return false
	}
	return time.Now().UTC().After(t)
}

// =============================================================================
// PromotionResponse → CouponDetailResp
// =============================================================================

// convertPromotionToCouponResp maps the unified PromotionResponse
// back onto the wire-shaped CouponDetailResp. The wire type still
// has the OLD shape (no Kind / MarketID / Code / TotalCount / Rules);
// Task 8 will regenerate them. Until then, fields not present on the
// wire type stay zero.
//
// DiscountType / DiscountValue / MinOrderAmount / MaxDiscount live on
// the first rule (one-rule-per-coupon form model); we surface them
// from there so the form re-hydrates after a refresh.
func convertPromotionToCouponResp(p *apppromotion.PromotionResponse) *types.CouponDetailResp {
	status := mapCouponStatusToWire(p.Status)
	if isExpiredRFC3339(p.EndAt.Format(time.RFC3339)) {
		status = "expired"
	}

	var (
		couponType      = "fixed_amount"
		discountValue   = "0"
		minOrderAmount  = "0"
		maxDiscount     = "0"
	)
	if len(p.Rules) > 0 {
		first := p.Rules[0]
		couponType = mapCouponTypeToWire(first.ActionType)
		discountValue = formatMoney(first.ActionValue)
		if first.ConditionType == pkgpromotion.ConditionMinAmount {
			minOrderAmount = formatMoney(first.ConditionValue)
		}
		maxDiscount = formatMoney(first.MaxDiscount)
	}

	var usageLimit, usedCount, perUserLimit int
	if p.TotalCount != nil {
		usageLimit = *p.TotalCount
	}
	if p.UsedCount != nil {
		usedCount = *p.UsedCount
	}
	perUserLimit = p.PerUserLimit

	code := ""
	if p.Code != nil {
		code = *p.Code
	}

	return &types.CouponDetailResp{
		ID:             p.ID,
		Code:           code,
		Name:           p.Name,
		Description:    p.Description,
		Type:           couponType,
		DiscountValue:  discountValue,
		MinOrderAmount: minOrderAmount,
		MaxDiscount:    maxDiscount,
		StartTime:      p.StartAt.Format(time.RFC3339),
		EndTime:        p.EndAt.Format(time.RFC3339),
		UsageLimit:     usageLimit,
		UsedCount:      usedCount,
		PerUserLimit:   perUserLimit,
		Status:         status,
		CreatedAt:      p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      p.UpdatedAt.Format(time.RFC3339),
	}
}