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

type BatchCreateShipmentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchCreateShipmentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) BatchCreateShipmentsLogic {
	return BatchCreateShipmentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchCreateShipmentsLogic) BatchCreateShipments(req *types.BatchCreateShipmentsReq) (resp *types.BatchCreateShipmentsResp, err error) {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Get user ID from context
	userID := contextx.GetCurrentUserID(l.ctx)

	// Build batch shipment items
	items := make([]appfulfillment.BatchShipmentItem, len(req.Shipments))
	for i, s := range req.Shipments {
		items[i] = appfulfillment.BatchShipmentItem{
			OrderID:    s.OrderID,
			TrackingNo: s.TrackingNo,
		}
	}

	// Batch create shipments
	result, err := l.svcCtx.ShipmentApp.BatchCreateShipments(l.ctx, shared.TenantID(tenantID), userID, req.CarrierCode, req.CarrierName, items)
	if err != nil {
		return nil, err
	}

	// Build response
	resp = &types.BatchCreateShipmentsResp{
		Total:   result.Total,
		Success: result.Success,
		Failed:  result.Failed,
		Results: make([]types.BatchShipmentResultResp, len(result.Results)),
	}

	for i, r := range result.Results {
		resp.Results[i] = types.BatchShipmentResultResp{
			OrderID:    r.OrderID,
			ShipmentID: r.ShipmentID,
			ShipmentNo: r.ShipmentNo,
			Success:    r.Success,
			Error:      r.Error,
		}
	}

	return resp, nil
}
