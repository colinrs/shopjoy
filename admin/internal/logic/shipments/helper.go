package shipments

import (
	"fmt"
	"time"

	appfulfillment "github.com/colinrs/shopjoy/admin/internal/application/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

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
		ID:            s.ID,
		ShipmentNo:    s.ShipmentNo,
		OrderID:       s.OrderID,
		Status:        int8(s.Status),
		StatusText:    fulfillment.ShipmentStatus(s.Status).String(),
		Carrier:       s.Carrier,
		CarrierCode:   s.CarrierCode,
		TrackingNo:    s.TrackingNo,
		TrackingURL:   s.TrackingURL,
		ShippingCost:  appfulfillment.FormatInt64ToMoney(s.ShippingCost),
		Currency:      s.Currency,
		Weight:        formatWeight(s.Weight),
		ShippedAt:     formatTimeToRFC3339(s.ShippedAt),
		DeliveredAt:   formatTimeToRFC3339(s.DeliveredAt),
		Remark:        s.Remark,
		Items:         items,
		CreatedAt:     s.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     s.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:     s.CreatedBy,
	}
}

// formatWeight formats weight to string
func formatWeight(w float64) string {
	if w == 0 {
		return "0"
	}
	return fmt.Sprintf("%.2f", w)
}

// formatTimeToRFC3339 formats time to RFC3339 string
func formatTimeToRFC3339(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// parseFloat parses a string to float64
func parseFloat(s string) float64 {
	var f float64
	_, _ = fmt.Sscanf(s, "%f", &f)
	return f
}