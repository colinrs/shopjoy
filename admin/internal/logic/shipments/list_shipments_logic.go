package shipments

import (
	"context"
	"time"

	appfulfillment "github.com/colinrs/shopjoy/admin/internal/application/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListShipmentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListShipmentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListShipmentsLogic {
	return ListShipmentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListShipmentsLogic) ListShipments(req *types.ListShipmentsReq) (resp *types.ListShipmentsResp, err error) {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Build query request
	queryReq := appfulfillment.QueryShipmentRequest{
		Page:         req.Page,
		PageSize:     req.PageSize,
		ShipmentNo:   req.ShipmentNo,
		OrderID:      req.OrderID,
		TrackingNo:   req.TrackingNo,
		CarrierCode:  req.CarrierCode,
		FulfillmentStatus: req.FulfillmentStatus,
	}

	// Parse status
	if req.Status > 0 {
		queryReq.Status = fulfillment.ShipmentStatus(req.Status)
	}

	// Parse time range
	if req.StartTime != "" {
		startTime, err := time.Parse(time.RFC3339, req.StartTime)
		if err == nil {
			queryReq.StartTime = startTime
		}
	}
	if req.EndTime != "" {
		endTime, err := time.Parse(time.RFC3339, req.EndTime)
		if err == nil {
			queryReq.EndTime = endTime
		}
	}

	// List shipments
	listResp, err := l.svcCtx.ShipmentApp.ListShipments(l.ctx, shared.TenantID(tenantID), queryReq)
	if err != nil {
		return nil, err
	}

	// Build response
	resp = &types.ListShipmentsResp{
		List:     make([]*types.ShipmentDetailResp, len(listResp.List)),
		Total:    listResp.Total,
		Page:     listResp.Page,
		PageSize: listResp.PageSize,
	}

	for i, s := range listResp.List {
		resp.List[i] = toShipmentDetailResp(s)
	}

	return resp, nil
}