package shipments

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

type UpdateShipmentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateShipmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateShipmentLogic {
	return UpdateShipmentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateShipmentLogic) UpdateShipment(req *types.UpdateShipmentReq) (resp *types.ShipmentDetailResp, err error) {
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

	// Build update request
	shippingCost, _ := decimal.NewFromString(req.ShippingCost)
	weight, _ := decimal.NewFromString(req.Weight)
	updateReq := appfulfillment.UpdateShipmentRequest{
		ID:           req.ID,
		CarrierCode:  req.CarrierCode,
		CarrierName:  req.CarrierName,
		TrackingNo:   req.TrackingNo,
		ShippingCost: shippingCost,
		Currency:     req.Currency,
		Weight:       weight,
		Remark:       req.Remark,
	}

	// Update shipment
	shipmentResp, err := l.svcCtx.ShipmentApp.UpdateShipment(l.ctx, shared.TenantID(tenantID), userID, updateReq)
	if err != nil {
		return nil, err
	}

	return toShipmentDetailResp(shipmentResp), nil
}
