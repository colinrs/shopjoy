package coupons

import (
	"strings"
	"time"

	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
	"github.com/colinrs/shopjoy/admin/internal/types"
	pkgcoupon "github.com/colinrs/shopjoy/pkg/domain/promotion"
	"github.com/shopspring/decimal"
)

func mapCouponType(typeStr string) pkgcoupon.CouponType {
	switch typeStr {
	case "fixed_amount":
		return pkgcoupon.CouponTypeFixedAmount
	case "percentage":
		return pkgcoupon.CouponTypePercentage
	default:
		return pkgcoupon.CouponTypeFixedAmount
	}
}

func mapCouponStatus(status int) string {
	switch pkgcoupon.CouponStatus(status) {
	case pkgcoupon.CouponStatusInactive:
		return "inactive"
	case pkgcoupon.CouponStatusActive:
		return "active"
	case pkgcoupon.CouponStatusExpired:
		return "expired"
	case pkgcoupon.CouponStatusDepleted:
		return "depleted"
	default:
		return "inactive"
	}
}

func mapCouponStatusToInt(statusStr string) pkgcoupon.CouponStatus {
	switch statusStr {
	case "active":
		return pkgcoupon.CouponStatusActive
	case "inactive":
		return pkgcoupon.CouponStatusInactive
	case "expired":
		return pkgcoupon.CouponStatusExpired
	case "depleted":
		return pkgcoupon.CouponStatusDepleted
	default:
		return pkgcoupon.CouponStatusInactive
	}
}

func mapCouponTypeToString(t pkgcoupon.CouponType) string {
	switch t {
	case pkgcoupon.CouponTypeFixedAmount:
		return "fixed_amount"
	case pkgcoupon.CouponTypePercentage:
		return "percentage"
	default:
		return "fixed_amount"
	}
}

func parseMoneyToDecimal(s string) decimal.Decimal {
	if s == "" {
		return decimal.Zero
	}
	// Try parsing as decimal
	v, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero
	}
	return v
}

func formatDecimalToString(v decimal.Decimal) string {
	if v.IsZero() {
		return "0"
	}
	return v.StringFixed(2)
}

// buildCouponScope maps the wire-level scope hint onto the domain's
// PromotionScope. product_ids / category_ids / market_ids on the
// wire are JSON-stringified arrays; we don't currently decode them
// because no coupon edit form binds them — a future form that does
// can replace the Scope assignment here with the full mapping.
//
// Wire values are lowercase ("products", "categories", …); the
// domain constants are uppercase ("PRODUCTS", …). Normalize via
// strings.ToUpper before comparing to the enum.
func buildCouponScope(scopeType string) pkgcoupon.PromotionScope {
	switch strings.ToUpper(scopeType) {
	case string(pkgcoupon.ScopeTypeProducts):
		return pkgcoupon.PromotionScope{Type: pkgcoupon.ScopeTypeProducts}
	case string(pkgcoupon.ScopeTypeCategories):
		return pkgcoupon.PromotionScope{Type: pkgcoupon.ScopeTypeCategories}
	case string(pkgcoupon.ScopeTypeBrands):
		return pkgcoupon.PromotionScope{Type: pkgcoupon.ScopeTypeBrands}
	}
	return pkgcoupon.PromotionScope{Type: pkgcoupon.ScopeTypeStorewide}
}

// convertCouponToDetailResp converts the application-layer response to
// the wire-format response. The status is computed: if the coupon is
// past its EndAt, "expired" is returned regardless of the stored
// status (the frontend renders this as a non-toggleable tag).
func convertCouponToDetailResp(c *apppromotion.CouponResponse) *types.CouponDetailResp {
	status := mapCouponStatus(c.Status)
	if isCouponExpired(c.EndAt) {
		status = "expired"
	}
	return &types.CouponDetailResp{
		ID:             c.ID,
		Code:           c.Code,
		Name:           c.Name,
		Description:    c.Description,
		Type:           mapCouponTypeToString(pkgcoupon.CouponType(c.Type)),
		DiscountValue:  formatDecimalToString(c.Value),
		MinOrderAmount: formatDecimalToString(c.MinAmount),
		MaxDiscount:    formatDecimalToString(c.MaxDiscount),
		StartTime:      c.StartAt,
		EndTime:        c.EndAt,
		UsageLimit:     c.TotalCount,
		UsedCount:      c.UsedCount,
		PerUserLimit:   c.PerUserLimit,
		Status:         status,
		ScopeType:      c.ScopeType,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
}

// isCouponExpired returns true if the coupon's EndAt is in the past.
// EndAt is an RFC3339 string (as emitted by the application layer). A
// malformed or zero value is treated as not expired so that the stored
// status is preserved.
func isCouponExpired(endAt string) bool {
	if endAt == "" {
		return false
	}
	t, err := time.Parse(time.RFC3339, endAt)
	if err != nil {
		return false
	}
	return time.Now().UTC().After(t)
}
