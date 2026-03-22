package fulfillment_orders

import (
	"context"
	"fmt"

	appfulfillment "github.com/colinrs/shopjoy/admin/internal/application/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShipOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShipOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) ShipOrderLogic {
	return ShipOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShipOrderLogic) ShipOrder(req *types.ShipOrderReq) (resp *types.ShipOrderResp, err error) {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Get user ID from context
	userID := contextx.GetCurrentUserID(l.ctx)

	// Build ship order request
	shipReq := appfulfillment.ShipOrderRequest{
		CarrierCode:  req.CarrierCode,
		CarrierName:  req.CarrierName,
		TrackingNo:   req.TrackingNo,
		ShippingCost: parseMoneyToInt64(req.ShippingCost),
		Currency:     req.Currency,
		Remark:       req.Remark,
	}

	// Parse weight
	if req.Weight != "" {
		shipReq.Weight = parseFloat(req.Weight)
	}

	// Build items
	if len(req.Items) > 0 {
		shipReq.Items = make([]appfulfillment.CreateShipmentItemRequest, len(req.Items))
		for i, item := range req.Items {
			shipReq.Items[i] = appfulfillment.CreateShipmentItemRequest{
				OrderItemID: item.OrderItemID,
				Quantity:    item.Quantity,
			}
		}
	}

	// Convert order ID to string
	orderID := fmt.Sprintf("%d", req.ID)

	// Ship order
	shipmentResp, err := l.svcCtx.OrderFulfillmentApp.ShipOrder(l.ctx, shared.TenantID(tenantID), userID, orderID, shipReq)
	if err != nil {
		return nil, err
	}

	return &types.ShipOrderResp{
		ShipmentID: shipmentResp.ID,
		ShipmentNo: shipmentResp.ShipmentNo,
	}, nil
}

// parseMoneyToInt64 parses a money string to int64 (cents)
func parseMoneyToInt64(s string) int64 {
	if s == "" {
		return 0
	}
	var v int64
	_, _ = fmt.Sscanf(s, "%d", &v)
	return v
}

// parseFloat parses a string to float64
func parseFloat(s string) float64 {
	var f float64
	_, _ = fmt.Sscanf(s, "%f", &f)
	return f
}