package shipments

import (
	"context"

	appfulfillment "github.com/colinrs/shopjoy/admin/internal/application/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

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
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Get user ID from context
	userID := contextx.GetCurrentUserID(l.ctx)

	// Build update request
	updateReq := appfulfillment.UpdateShipmentRequest{
		ID:           req.ID,
		CarrierCode:  req.CarrierCode,
		CarrierName:  req.CarrierName,
		TrackingNo:   req.TrackingNo,
		ShippingCost: appfulfillment.FormatMoneyToInt64(req.ShippingCost),
		Currency:     req.Currency,
		Remark:       req.Remark,
	}

	// Parse weight
	if req.Weight != "" {
		updateReq.Weight = parseFloat(req.Weight)
	}

	// Update shipment
	shipmentResp, err := l.svcCtx.ShipmentApp.UpdateShipment(l.ctx, shared.TenantID(tenantID), userID, updateReq)
	if err != nil {
		return nil, err
	}

	return toShipmentDetailResp(shipmentResp), nil
}