package promotions

import (
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

func mapPromotionStatus(status int) string {
	switch pkgpromotion.Status(status) {
	case pkgpromotion.StatusPending:
		return "draft"
	case pkgpromotion.StatusActive:
		return "active"
	case pkgpromotion.StatusPaused:
		return "inactive"
	case pkgpromotion.StatusEnded:
		return "expired"
	default:
		return "draft"
	}
}

func mapPromotionStatusToInt(statusStr string) pkgpromotion.Status {
	switch statusStr {
	case "draft":
		return pkgpromotion.StatusPending
	case "active":
		return pkgpromotion.StatusActive
	case "inactive":
		return pkgpromotion.StatusPaused
	case "expired":
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

func mapActionTypeToString(actionType pkgpromotion.ActionType) string {
	switch actionType {
	case pkgpromotion.ActionPercentage:
		return "percentage"
	case pkgpromotion.ActionFixedAmount:
		return "fixed_amount"
	default:
		return "fixed_amount"
	}
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

func convertPromotionToDetailResp(p *apppromotion.PromotionResponse) *types.PromotionDetailResp {
	return &types.PromotionDetailResp{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Type:        mapPromotionTypeToString(pkgpromotion.Type(p.Type)),
		Status:      mapPromotionStatus(p.Status),
		StartTime:   p.StartAt,
		EndTime:     p.EndAt,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}