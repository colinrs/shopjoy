package shipments

import (
	"context"
	"encoding/csv"
	"net/http"
	"time"

	appfulfillment "github.com/colinrs/shopjoy/admin/internal/application/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExportShipmentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	w      http.ResponseWriter
	r      *http.Request
}

func NewExportShipmentsLogic(ctx context.Context, svcCtx *svc.ServiceContext, w http.ResponseWriter, r *http.Request) ExportShipmentsLogic {
	return ExportShipmentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		w:      w,
		r:      r,
	}
}

func (l *ExportShipmentsLogic) ExportShipments(req *types.ExportShipmentsReq) error {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Build query request with large page size for export
	queryReq := appfulfillment.QueryShipmentRequest{
		Page:      1,
		PageSize:  10001, // Check if exceeds limit
		TrackingNo: req.TrackingNo,
		CarrierCode: req.CarrierCode,
	}

	// Parse status - convert from string to domain enum
	if req.Status != "" {
		queryReq.Status = fulfillment.ParseShipmentStatus(req.Status)
	}

	// Parse start time
	if req.StartTime != "" {
		startTime, err := time.Parse(time.RFC3339, req.StartTime)
		if err == nil {
			queryReq.StartTime = startTime
		}
	}

	// Parse end time
	if req.EndTime != "" {
		endTime, err := time.Parse(time.RFC3339, req.EndTime)
		if err == nil {
			queryReq.EndTime = endTime
		}
	}

	// Get shipments for export
	listResp, err := l.svcCtx.ShipmentApp.ListShipments(l.ctx, shared.TenantID(tenantID), queryReq)
	if err != nil {
		return err
	}

	if listResp.Total > 10000 {
		return code.ErrShipmentExportLimitExceed
	}

	// Set response headers for CSV download
	l.w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	l.w.Header().Set("Content-Disposition", "attachment; filename=shipments_export_"+time.Now().Format("20060102")+".csv")

	// Write UTF-8 BOM for Excel compatibility
	l.w.Write([]byte{0xEF, 0xBB, 0xBF})

	// Create CSV writer
	writer := csv.NewWriter(l.w)
	defer writer.Flush()

	// Write header
	header := []string{
		"发货单号", "订单号", "物流公司", "物流单号", "状态",
		"发货时间", "送达时间",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data rows
	for _, row := range listResp.List {
		record := []string{
			row.ShipmentNo,
			row.OrderID,
			row.Carrier,
			row.TrackingNo,
			fulfillment.ShipmentStatus(row.Status).String(),
			formatTimeToRFC3339(row.ShippedAt),
			formatTimeToRFC3339(row.DeliveredAt),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}