package fulfillment_orders

import (
	"fmt"
	"time"
	"unicode/utf8"

	appfulfillment "github.com/colinrs/shopjoy/admin/internal/application/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/utils"
	"github.com/shopspring/decimal"
)

// parseMoneyToDecimal parses a money string to decimal.Decimal
func parseMoneyToDecimal(s string) (decimal.Decimal, error) {
	if s == "" {
		return decimal.Zero, nil
	}
	v, err := decimal.NewFromString(s)
	return v, err
}

// truncateString truncates a string to maxChars characters (UTF-8 safe)
func truncateString(s string, maxChars int) string {
	if utf8.RuneCountInString(s) <= maxChars {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxChars])
}

// formatFloatToString formats a float64 to string
func formatFloatToString(w float64) string {
	return fmt.Sprintf("%.2f", w)
}

// toShipmentDetailResp converts ShipmentResponse to types.ShipmentDetailResp
func toShipmentDetailResp(s *appfulfillment.ShipmentResponse) *types.ShipmentDetailResp {
	items := make([]*types.ShipmentItemResp, len(s.Items))
	for i, item := range s.Items {
		items[i] = &types.ShipmentItemResp{
			ID:          item.ID,
			ShipmentID:  item.ShipmentID,
			OrderItemID: item.OrderItemID,
			ProductID:   item.ProductID,
			SKUID:       item.SKUID,
			ProductName: item.ProductName,
			SKUName:     item.SKUName,
			Image:       item.Image,
			Quantity:    item.Quantity,
		}
	}

	return &types.ShipmentDetailResp{
		ID:           s.ID,
		ShipmentNo:   s.ShipmentNo,
		OrderID:      s.OrderID,
		Status:       fulfillment.ShipmentStatus(s.Status).String(),
		Carrier:      s.Carrier,
		CarrierCode:  s.CarrierCode,
		TrackingNo:   s.TrackingNo,
		TrackingURL:  s.TrackingURL,
		ShippingCost: s.ShippingCost,
		Currency:     s.Currency,
		Weight:       s.Weight,
		ShippedAt:    utils.FormatTimeToRFC3339(s.ShippedAt),
		DeliveredAt:  utils.FormatTimeToRFC3339(s.DeliveredAt),
		Remark:       s.Remark,
		Items:        items,
		CreatedAt:    s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    s.UpdatedAt.Format(time.RFC3339),
		CreatedBy:    s.CreatedBy,
	}
}

// formatWeight formats weight to string
func formatWeight(w float64) string {
	if w == 0 {
		return "0"
	}
	return fmt.Sprintf("%.2f", w)
}

// toRefundDetailResp converts RefundResponse to types.RefundDetailResp
func toRefundDetailResp(r *appfulfillment.RefundResponse) *types.RefundDetailResp {
	return &types.RefundDetailResp{
		ID:             r.ID,
		RefundNo:       r.RefundNo,
		OrderID:        r.OrderID,
		UserID:         r.UserID,
		UserName:       r.UserName,
		Type:           fulfillment.RefundType(r.Type).String(),
		TypeText:       r.TypeText,
		Status:         fulfillment.RefundStatus(r.Status).String(),
		StatusText:     r.StatusText,
		ReasonType:     r.ReasonType,
		Reason:         r.Reason,
		Description:    r.Description,
		Images:         r.Images,
		Amount:         utils.FormatAmount(r.Amount),
		Currency:       r.Currency,
		RejectReason:   r.RejectReason,
		ApprovedAt:     utils.FormatTimeToRFC3339(r.ApprovedAt),
		ApprovedBy:     r.ApprovedBy,
		ApprovedByName: r.ApprovedByName,
		CompletedAt:    utils.FormatTimeToRFC3339(r.CompletedAt),
		CreatedAt:      r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      r.UpdatedAt.Format(time.RFC3339),
	}
}
