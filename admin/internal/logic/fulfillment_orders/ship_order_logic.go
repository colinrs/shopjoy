package fulfillment_orders

import (
	"context"

	appfulfillment "github.com/colinrs/shopjoy/admin/internal/application/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"

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
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok && !contextx.IsPlatformAdmin(l.ctx) {
		return nil, code.ErrUnauthorized
	}
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Get user ID from context
	userID := contextx.GetCurrentUserID(l.ctx)

	// Parse shipping cost
	shippingCost, _ := decimal.NewFromString(req.ShippingCost)
	weight, _ := decimal.NewFromString(req.Weight)

	// Build ship order request
	shipReq := appfulfillment.ShipOrderRequest{
		CarrierCode:  req.CarrierCode,
		CarrierName:  req.CarrierName,
		TrackingNo:   req.TrackingNo,
		ShippingCost: shippingCost,
		Currency:     req.Currency,
		Weight:       weight,
		Remark:       req.Remark,
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

	// Ship order
	shipmentResp, err := l.svcCtx.OrderFulfillmentApp.ShipOrder(l.ctx, shared.TenantID(tenantID), userID, req.ID, shipReq)
	if err != nil {
		return nil, err
	}

	return &types.ShipOrderResp{
		ShipmentID: shipmentResp.ID,
		ShipmentNo: shipmentResp.ShipmentNo,
	}, nil
}
