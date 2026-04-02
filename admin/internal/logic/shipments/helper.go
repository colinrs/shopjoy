package shipments

import (
	"fmt"
	"strconv"
	"time"

	appfulfillment "github.com/colinrs/shopjoy/admin/internal/application/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/utils"
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
		Status:        fulfillment.ShipmentStatus(s.Status).String(),
		Carrier:       s.Carrier,
		CarrierCode:   s.CarrierCode,
		TrackingNo:    s.TrackingNo,
		TrackingURL:   s.TrackingURL,
		ShippingCost:  s.ShippingCost,
		Currency:      s.Currency,
		Weight:        s.Weight,
		ShippedAt:     utils.FormatTimeToRFC3339(s.ShippedAt),
		DeliveredAt:   utils.FormatTimeToRFC3339(s.DeliveredAt),
		Remark:        s.Remark,
		Items:         items,
		CreatedAt:     s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     s.UpdatedAt.Format(time.RFC3339),
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

// parseFloat parses a string to float64
func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}