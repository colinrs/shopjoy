package promotions

import (
	"strconv"
	"strings"
	"time"

	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
	"github.com/colinrs/shopjoy/admin/internal/types"
	pkgpromotion "github.com/colinrs/shopjoy/pkg/domain/promotion"
	"github.com/shopspring/decimal"
)

// =============================================================================
// Type / status / kind mappings (wire <-> domain)
// =============================================================================

func mapPromotionType(typeStr string) pkgpromotion.Type {
	switch typeStr {
	case "discount":
		return pkgpromotion.TypeDiscount
	case "flash_sale":
		return pkgpromotion.TypeFlashSale
	case "bundle":
		return pkgpromotion.TypeBundle
	case "buy_x_get_y":
		return pkgpromotion.TypeBuyXGetY
	default:
		return pkgpromotion.TypeDiscount
	}
}

// parsePromotionType is the new wire-side entry point. The wire sends
// lowercase ("discount", "flash_sale", …); normalize to the domain
// constant and validate.
func parsePromotionType(typeStr string) pkgpromotion.Type {
	return mapPromotionType(strings.ToLower(typeStr))
}

func mapPromotionKindToString(k pkgpromotion.Kind) string {
	switch k {
	case pkgpromotion.KindPromotion:
		return "promotion"
	case pkgpromotion.KindCoupon:
		return "coupon"
	default:
		return "promotion"
	}
}

// mapPromotionStatus converts a stored promotion status (int) to its
// wire-format string. Values match the .api file comment and the
// frontend's PromotionStatus type verbatim.
func mapPromotionStatus(status pkgpromotion.Status) string {
	switch status {
	case pkgpromotion.StatusPending:
		return "pending"
	case pkgpromotion.StatusActive:
		return "active"
	case pkgpromotion.StatusPaused:
		return "paused"
	case pkgpromotion.StatusEnded:
		return "ended"
	default:
		return "pending"
	}
}

// mapPromotionStatusToInt is the inverse of mapPromotionStatus; used to
// translate the status query parameter on the list endpoint.
func mapPromotionStatusToInt(statusStr string) pkgpromotion.Status {
	switch statusStr {
	case "pending":
		return pkgpromotion.StatusPending
	case "active":
		return pkgpromotion.StatusActive
	case "paused":
		return pkgpromotion.StatusPaused
	case "ended":
		return pkgpromotion.StatusEnded
	default:
		return pkgpromotion.StatusPending
	}
}

func mapPromotionTypeToString(t pkgpromotion.Type) string {
	switch t {
	case pkgpromotion.TypeDiscount:
		return "discount"
	case pkgpromotion.TypeFlashSale:
		return "flash_sale"
	case pkgpromotion.TypeBundle:
		return "bundle"
	case pkgpromotion.TypeBuyXGetY:
		return "buy_x_get_y"
	default:
		return "discount"
	}
}

func mapDiscountActionType(discountType string) pkgpromotion.ActionType {
	if strings.EqualFold(discountType, "percentage") {
		return pkgpromotion.ActionPercentage
	}
	return pkgpromotion.ActionFixedAmount
}

// mapCouponActionType maps the coupon's wire Type ("fixed_amount",
// "percentage", "free_shipping") onto the domain ActionType. The wire
// is largely unchanged from the pre-merge coupon form so the
// translation stays simple.
func mapCouponActionType(couponType string) pkgpromotion.ActionType {
	switch strings.ToLower(couponType) {
	case "percentage":
		return pkgpromotion.ActionPercentage
	case "free_shipping":
		return pkgpromotion.ActionFreeShipping
	default:
		return pkgpromotion.ActionFixedAmount
	}
}

func mapConditionType(ruleType string) pkgpromotion.ConditionType {
	switch strings.ToLower(ruleType) {
	case "amount":
		return pkgpromotion.ConditionMinAmount
	case "quantity":
		return pkgpromotion.ConditionMinQuantity
	default:
		return pkgpromotion.ConditionMinAmount
	}
}

func mapConditionTypeToString(conditionType pkgpromotion.ConditionType) string {
	switch conditionType {
	case pkgpromotion.ConditionMinAmount:
		return "amount"
	case pkgpromotion.ConditionMinQuantity:
		return "quantity"
	default:
		return "amount"
	}
}

func mapActionTypeIntToString(actionType pkgpromotion.ActionType) string {
	switch actionType {
	case pkgpromotion.ActionFixedAmount:
		return "fixed_amount"
	case pkgpromotion.ActionPercentage:
		return "percentage"
	case pkgpromotion.ActionFreeShipping:
		return "free_shipping"
	default:
		return "fixed_amount"
	}
}

// =============================================================================
// Money / decimal helpers
// =============================================================================

func parseMoneyToDecimal(s string) decimal.Decimal {
	if s == "" {
		return decimal.Zero
	}
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

// currencyWithDefault returns the wire-supplied currency or "CNY" if
// the frontend omitted it.
func currencyWithDefault(c string) string {
	if c == "" {
		return "CNY"
	}
	return c
}

// optionalInt64 returns a pointer to v, or nil if v is zero. Used to
// translate the wire "0 = unset" convention onto the domain's nullable
// MarketID pointer.
func optionalInt64(v int64) *int64 {
	if v == 0 {
		return nil
	}
	return &v
}

// parseTime parses an RFC3339 string into a UTC time.Time. Returns the
// zero value (and the parse error) when the input is empty or
// malformed; callers should propagate the error.
func parseTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, s)
}

// =============================================================================
// Scope helpers (shared by PROMOTION and COUPON create/update)
// =============================================================================

// buildPromotionScope maps the wire-level (scope_type, product_ids,
// category_ids, brand_ids) tuple onto the domain's PromotionScope.
// Empty IDs for the chosen scope are normalized to nil so the storage
// layer doesn't carry stale empty arrays.
//
// Wire values are lowercase ("products", "categories", "brands", …)
// but the domain constants are uppercase ("PRODUCTS", …). We uppercase
// before comparing so the form's emitted value matches the enum.
func buildPromotionScope(scopeType string, productIDs, categoryIDs, brandIDs []string) pkgpromotion.PromotionScope {
	normalize := func(s string) pkgpromotion.ScopeType {
		switch strings.ToUpper(s) {
		case string(pkgpromotion.ScopeTypeProducts):
			return pkgpromotion.ScopeTypeProducts
		case string(pkgpromotion.ScopeTypeCategories):
			return pkgpromotion.ScopeTypeCategories
		case string(pkgpromotion.ScopeTypeBrands):
			return pkgpromotion.ScopeTypeBrands
		case string(pkgpromotion.ScopeTypeStorewide):
			return pkgpromotion.ScopeTypeStorewide
		}
		return ""
	}

	st := normalize(scopeType)
	if st == "" {
		switch {
		case len(productIDs) > 0:
			st = pkgpromotion.ScopeTypeProducts
		case len(categoryIDs) > 0:
			st = pkgpromotion.ScopeTypeCategories
		case len(brandIDs) > 0:
			st = pkgpromotion.ScopeTypeBrands
		default:
			st = pkgpromotion.ScopeTypeStorewide
		}
	}

	var ids []int64
	switch st {
	case pkgpromotion.ScopeTypeProducts:
		ids = parseInt64Slice(productIDs)
	case pkgpromotion.ScopeTypeCategories:
		ids = parseInt64Slice(categoryIDs)
	case pkgpromotion.ScopeTypeBrands:
		ids = parseInt64Slice(brandIDs)
	default:
		ids = nil
	}
	return pkgpromotion.PromotionScope{Type: st, IDs: ids}
}

// buildCouponScope is a thin wrapper used by the coupon create/update
// logic. The wire shape currently does not bind scope IDs for coupons
// (product_ids / category_ids are JSON strings the form never
// re-populates), so the resulting PromotionScope carries the chosen
// scope_type with a nil ID slice.
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

// parseInt64Slice converts []string of decimal numeric IDs into
// []int64, ignoring entries that fail to parse.
func parseInt64Slice(ss []string) []int64 {
	if len(ss) == 0 {
		return nil
	}
	out := make([]int64, 0, len(ss))
	for _, s := range ss {
		if n, err := strconv.ParseInt(s, 10, 64); err == nil {
			out = append(out, n)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// idsToStrings converts []int64 IDs to []string (decimal form). Used
// when surfacing the domain Scope back into the wire response. Returns
// nil for an empty slice so the JSON tag can emit "omitted".
func idsToStrings(ids []int64) []string {
	if len(ids) == 0 {
		return nil
	}
	out := make([]string, len(ids))
	for i, n := range ids {
		out[i] = strconv.FormatInt(n, 10)
	}
	return out
}

// =============================================================================
// Rule conversion (domain <-> wire)
// =============================================================================

// convertRuleReqsToDomain maps the wire-shaped rule request (rule_type,
// value, discount_type, discount_value) onto the domain
// PromotionRule. OwnerKind / OwnerID are set by the caller (they
// aren't known until the promotion row exists).
func convertRuleReqsToDomain(reqs []types.PromotionRuleReq) []pkgpromotion.PromotionRule {
	if len(reqs) == 0 {
		return nil
	}
	out := make([]pkgpromotion.PromotionRule, 0, len(reqs))
	for _, r := range reqs {
		out = append(out, pkgpromotion.PromotionRule{
			ConditionType:  mapConditionType(r.RuleType),
			ConditionValue: parseMoneyToDecimal(r.Value),
			ActionType:     mapDiscountActionType(r.DiscountType),
			ActionValue:    parseMoneyToDecimal(r.DiscountValue),
		})
	}
	return out
}

// convertRulesToResp maps app → wire types for rules. The wire
// PromotionRuleResp keeps the old flat shape (rule_type /
// discount_type / discount_value / value / priority) so the form can
// re-hydrate unchanged.
func convertRulesToResp(rules []*apppromotion.PromotionRuleResponse) []*types.PromotionRuleResp {
	if len(rules) == 0 {
		return nil
	}
	out := make([]*types.PromotionRuleResp, 0, len(rules))
	for _, r := range rules {
		out = append(out, &types.PromotionRuleResp{
			ID:            r.ID,
			RuleType:      mapConditionTypeToString(r.ConditionType),
			Operator:      "gte",
			Value:         formatDecimalToString(r.ConditionValue),
			DiscountType:  mapActionTypeIntToString(r.ActionType),
			DiscountValue: formatDecimalToString(r.ActionValue),
			Priority:      0,
			CreatedAt:     "",
			UpdatedAt:     "",
		})
	}
	return out
}

// =============================================================================
// PromotionResponse → PromotionDetailResp
// =============================================================================

// convertPromotionToDetailResp maps a unified PromotionResponse to the
// wire PromotionDetailResp. The wire type still has the OLD shape
// (no Kind / MarketID / Code / TotalCount / Rules fields); Task 8
// will regenerate them. Until then, those fields stay zero/empty so
// existing forms continue to render.
func convertPromotionToDetailResp(p *apppromotion.PromotionResponse) *types.PromotionDetailResp {
	status := mapPromotionStatus(p.Status)
	if isPromotionExpired(p.EndAt) {
		status = "expired"
	}

	// DiscountType/DiscountValue/MinOrderAmount/MaxDiscount aren't
	// stored on the Promotion row; they're stored as PromotionRules.
	// Surface the first rule's values when present so the form
	// re-renders correctly after the user refreshes.
	var (
		discountType, discountValue, minOrderAmount, maxDiscount string
	)
	if len(p.Rules) > 0 {
		first := p.Rules[0]
		discountType = mapActionTypeIntToString(first.ActionType)
		if !first.ActionValue.IsZero() {
			discountValue = first.ActionValue.StringFixed(2)
		}
		if first.ConditionType == pkgpromotion.ConditionMinAmount && !first.ConditionValue.IsZero() {
			minOrderAmount = first.ConditionValue.StringFixed(2)
		}
		if !first.MaxDiscount.IsZero() {
			maxDiscount = first.MaxDiscount.StringFixed(2)
		}
	}

	// Split the stored Scope.IDs back into the per-type wire field
	// so the form re-hydrates correctly after a refresh. Only the
	// array matching ScopeType is populated — the others stay nil
	// to mirror how the form sends a single non-empty ID array.
	var productIDs, categoryIDs, brandIDs []string
	if len(p.Scope.IDs) > 0 {
		switch p.Scope.Type {
		case pkgpromotion.ScopeTypeProducts:
			productIDs = idsToStrings(p.Scope.IDs)
		case pkgpromotion.ScopeTypeCategories:
			categoryIDs = idsToStrings(p.Scope.IDs)
		case pkgpromotion.ScopeTypeBrands:
			brandIDs = idsToStrings(p.Scope.IDs)
		}
	}

	resp := &types.PromotionDetailResp{
		ID:             p.ID,
		Name:           p.Name,
		Description:    p.Description,
		Type:           mapPromotionTypeToString(p.Type),
		Status:         status,
		StartTime:      p.StartAt.Format(time.RFC3339),
		EndTime:        p.EndAt.Format(time.RFC3339),
		DiscountType:   discountType,
		DiscountValue:  discountValue,
		MinOrderAmount: minOrderAmount,
		MaxDiscount:    maxDiscount,
		Currency:       p.Currency,
		UsageLimit:     p.UsageLimit,
		PerUserLimit:   p.PerUserLimit,
		ProductIDs:     productIDs,
		CategoryIDs:    categoryIDs,
		BrandIDs:       brandIDs,
		Tags:           p.Tags,
		ScopeType:      string(p.Scope.Type),
		CreatedAt:      p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      p.UpdatedAt.Format(time.RFC3339),
	}
	return resp
}

// isPromotionExpired returns true if the promotion's EndAt is in the
// past.
func isPromotionExpired(endAt time.Time) bool {
	if endAt.IsZero() {
		return false
	}
	return time.Now().UTC().After(endAt)
}