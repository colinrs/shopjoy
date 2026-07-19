package promotions

import (
	"time"

	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
	"github.com/colinrs/shopjoy/admin/internal/types"
	pkgpromotion "github.com/colinrs/shopjoy/pkg/domain/promotion"
	"github.com/shopspring/decimal"
)

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
	case "full_reduce":
		return pkgpromotion.TypeFullReduce
	default:
		return pkgpromotion.TypeDiscount
	}
}

// mapPromotionStatus converts a stored promotion status (int) to its
// wire-format string. Values match the .api file comment and the
// frontend's PromotionStatus type verbatim.
func mapPromotionStatus(status int) string {
	switch pkgpromotion.Status(status) {
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
	case pkgpromotion.TypeFullReduce:
		return "full_reduce"
	default:
		return "discount"
	}
}

func mapDiscountActionType(discountType string) pkgpromotion.ActionType {
	if discountType == "percentage" {
		return pkgpromotion.ActionPercentage
	}
	return pkgpromotion.ActionFixedAmount
}

func mapConditionType(ruleType string) pkgpromotion.ConditionType {
	switch ruleType {
	case "amount":
		return pkgpromotion.ConditionMinAmount
	case "quantity":
		return pkgpromotion.ConditionMinQuantity
	default:
		return pkgpromotion.ConditionMinAmount
	}
}

func mapConditionTypeToString(conditionType int) string {
	switch pkgpromotion.ConditionType(conditionType) {
	case pkgpromotion.ConditionMinAmount:
		return "amount"
	case pkgpromotion.ConditionMinQuantity:
		return "quantity"
	default:
		return "amount"
	}
}

func mapActionTypeIntToString(actionType int) string {
	switch pkgpromotion.ActionType(actionType) {
	case pkgpromotion.ActionFixedAmount:
		return "fixed_amount"
	case pkgpromotion.ActionPercentage:
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

// convertPromotionToDetailResp converts the application-layer response to
// the wire-format response. The status is computed: if the promotion is
// past its EndAt, "expired" is returned regardless of the stored status
// (the frontend renders this as a non-toggleable tag).
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
		if pkgpromotion.ConditionType(first.ConditionType) == pkgpromotion.ConditionMinAmount && !first.ConditionValue.IsZero() {
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
	scopeIDs := pkgpromotion.ScopeType(p.ScopeType)
	if len(p.ScopeIDs) > 0 {
		switch scopeIDs {
		case pkgpromotion.ScopeTypeProducts:
			productIDs = p.ScopeIDs
		case pkgpromotion.ScopeTypeCategories:
			categoryIDs = p.ScopeIDs
		case pkgpromotion.ScopeTypeBrands:
			brandIDs = p.ScopeIDs
		}
	}

	return &types.PromotionDetailResp{
		ID:             p.ID,
		Name:           p.Name,
		Description:    p.Description,
		Type:           mapPromotionTypeToString(pkgpromotion.Type(p.Type)),
		Status:         status,
		StartTime:      p.StartAt,
		EndTime:        p.EndAt,
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
		ScopeType:      p.ScopeType,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}

// isPromotionExpired returns true if the promotion's EndAt is in the
// past. EndAt is an RFC3339 string (as emitted by the application
// layer). A malformed or zero value is treated as not expired so that
// the stored status is preserved.
func isPromotionExpired(endAt string) bool {
	if endAt == "" {
		return false
	}
	t, err := time.Parse(time.RFC3339, endAt)
	if err != nil {
		return false
	}
	return time.Now().UTC().After(t)
}
