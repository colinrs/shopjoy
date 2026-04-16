package shipments

import (
	"context"

	appfulfillment "github.com/colinrs/shopjoy/admin/internal/application/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateShipmentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateShipmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateShipmentLogic {
	return CreateShipmentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateShipmentLogic) CreateShipment(req *types.CreateShipmentReq) (resp *types.CreateShipmentResp, err error) {
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

	// Build request
	shippingCost, _ := decimal.NewFromString(req.ShippingCost)
	createReq := appfulfillment.CreateShipmentRequest{
		OrderID:      req.OrderID,
		CarrierCode:  req.CarrierCode,
		CarrierName:  req.CarrierName,
		TrackingNo:   req.TrackingNo,
		ShippingCost: shippingCost,
		Currency:     req.Currency,
		Remark:       req.Remark,
	}

	// Parse weight
	weight, _ := decimal.NewFromString(req.Weight)
	createReq.Weight = weight

	// Build items
	if len(req.Items) > 0 {
		createReq.Items = make([]appfulfillment.CreateShipmentItemRequest, len(req.Items))
		for i, item := range req.Items {
			createReq.Items[i] = appfulfillment.CreateShipmentItemRequest{
				OrderItemID: item.OrderItemID,
				Quantity:    item.Quantity,
			}
		}
	}

	// Create shipment
	shipmentResp, err := l.svcCtx.ShipmentApp.CreateShipment(l.ctx, shared.TenantID(tenantID), userID, createReq)
	if err != nil {
		return nil, err
	}

	return &types.CreateShipmentResp{
		ID:         shipmentResp.ID,
		ShipmentNo: shipmentResp.ShipmentNo,
	}, nil
}

// mapShipmentStatusToInt maps shipment status to int8
func mapShipmentStatusToInt(status fulfillment.ShipmentStatus) int8 {
	return int8(status) // #nosec G115 // status values are small (tinyint range)
}

// mapShipmentStatusToString maps shipment status to string
func mapShipmentStatusToString(status int8) string {
	return fulfillment.ShipmentStatus(status).String()
}
