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

// mapWireConditionTypeToString is the inverse of mapWireConditionType.
// The wire shape uses "min_amount" / "min_quantity" so the form
// re-hydrates with the same enum strings the create endpoint accepts.
func mapWireConditionTypeToString(conditionType pkgpromotion.ConditionType) string {
	switch conditionType {
	case pkgpromotion.ConditionMinQuantity:
		return "min_quantity"
	case pkgpromotion.ConditionMinAmount:
		return "min_amount"
	default:
		return "min_amount"
	}
}

// mapWireActionTypeToString is the inverse of mapWireActionType. The
// wire shape uses "fixed_amount" / "percentage" / "free_shipping" so
// the form re-hydrates with the same enum strings the create endpoint
// accepts.
func mapWireActionTypeToString(actionType pkgpromotion.ActionType) string {
	return mapActionTypeIntToString(actionType)
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

// convertRuleReqsToDomain maps the wire-shaped rule request
// (condition_type / condition_value / action_type / action_value)
// onto the domain PromotionRule. OwnerKind / OwnerID are set by the
// caller (they aren't known until the promotion row exists).
func convertRuleReqsToDomain(reqs []types.PromotionRuleReq) []pkgpromotion.PromotionRule {
	if len(reqs) == 0 {
		return nil
	}
	out := make([]pkgpromotion.PromotionRule, 0, len(reqs))
	for _, r := range reqs {
		out = append(out, pkgpromotion.PromotionRule{
			ConditionType:  mapWireConditionType(r.ConditionType),
			ConditionValue: parseMoneyToDecimal(r.ConditionValue),
			ActionType:     mapWireActionType(r.ActionType),
			ActionValue:    parseMoneyToDecimal(r.ActionValue),
			MaxDiscount:    parseMoneyToDecimal(r.MaxDiscount),
			SortOrder:      r.SortOrder,
		})
	}
	return out
}

// convertRulesToResp maps app → wire types for rules. The wire
// PromotionRuleResp uses the unified shape
// (condition_type / condition_value / action_type / action_value /
// max_discount / sort_order).
func convertRulesToResp(rules []apppromotion.PromotionRuleResponse) []*types.PromotionRuleResp {
	if len(rules) == 0 {
		return nil
	}
	out := make([]*types.PromotionRuleResp, 0, len(rules))
	for _, r := range rules {
		out = append(out, &types.PromotionRuleResp{
			ID:             r.ID,
			ConditionType:  mapWireConditionTypeToString(r.ConditionType),
			ConditionValue: formatDecimalToString(r.ConditionValue),
			ActionType:     mapWireActionTypeToString(r.ActionType),
			ActionValue:    formatDecimalToString(r.ActionValue),
			MaxDiscount:    formatDecimalToString(r.MaxDiscount),
			SortOrder:      r.SortOrder,
		})
	}
	return out
}

// convertRulesPtrToResp mirrors convertRulesToResp but accepts the
// pointer slice returned by PromotionApp.GetRules / CreateRules.
func convertRulesPtrToResp(rules []*apppromotion.PromotionRuleResponse) []*types.PromotionRuleResp {
	if len(rules) == 0 {
		return nil
	}
	out := make([]*types.PromotionRuleResp, 0, len(rules))
	for _, r := range rules {
		if r == nil {
			continue
		}
		out = append(out, &types.PromotionRuleResp{
			ID:             r.ID,
			ConditionType:  mapWireConditionTypeToString(r.ConditionType),
			ConditionValue: formatDecimalToString(r.ConditionValue),
			ActionType:     mapWireActionTypeToString(r.ActionType),
			ActionValue:    formatDecimalToString(r.ActionValue),
			MaxDiscount:    formatDecimalToString(r.MaxDiscount),
			SortOrder:      r.SortOrder,
		})
	}
	return out
}

// =============================================================================
// PromotionResponse → PromotionDetailResp
// =============================================================================

// convertPromotionToDetailResp maps a unified PromotionResponse to the
// wire PromotionDetailResp. The wire type uses the unified shape
// (Kind / MarketID / Code / TotalCount / Rules) so the form can
// render either promotion or coupon detail from a single endpoint.
func convertPromotionToDetailResp(p *apppromotion.PromotionResponse) *types.PromotionDetailResp {
	status := mapPromotionStatus(p.Status)
	if isPromotionExpired(p.EndAt) {
		status = "expired"
	}

	code := ""
	if p.Code != nil {
		code = *p.Code
	}

	var totalCount int
	if p.TotalCount != nil {
		totalCount = *p.TotalCount
	}

	var marketID int64
	if p.MarketID != nil {
		marketID = *p.MarketID
	}

	var usedCount int
	if p.UsedCount != nil {
		usedCount = *p.UsedCount
	}

	resp := &types.PromotionDetailResp{
		ID:           p.ID,
		Kind:         mapPromotionKindToString(p.Kind),
		Name:         p.Name,
		Description:  p.Description,
		Code:         code,
		Type:         mapPromotionTypeToString(p.Type),
		Status:       status,
		MarketID:     marketID,
		Currency:     p.Currency,
		UsageLimit:   p.UsageLimit,
		UsedCount:    usedCount,
		PerUserLimit: p.PerUserLimit,
		TotalCount:   totalCount,
		ScopeType:    strings.ToLower(string(p.Scope.Type)),
		ScopeIDs:     idsToStrings(p.Scope.IDs),
		ExcludeIDs:   idsToStrings(p.Scope.ExcludeIDs),
		Tags:         p.Tags,
		Rules:        convertRulesToResp(p.Rules),
		StartTime:    p.StartAt.Format(time.RFC3339),
		EndTime:      p.EndAt.Format(time.RFC3339),
		CreatedAt:    p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    p.UpdatedAt.Format(time.RFC3339),
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