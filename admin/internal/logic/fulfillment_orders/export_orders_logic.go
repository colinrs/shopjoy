package fulfillment_orders

import (
	"context"
	"encoding/csv"
	"net/http"
	"strconv"
	"time"

	appfulfillment "github.com/colinrs/shopjoy/admin/internal/application/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExportOrdersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
	w      http.ResponseWriter
}

func NewExportOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext, w http.ResponseWriter, r *http.Request) ExportOrdersLogic {
	return ExportOrdersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		w:      w,
		r:      r,
	}
}

func (l *ExportOrdersLogic) ExportOrders(req *types.ExportOrdersReq) error {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Build query request
	queryReq := appfulfillment.QueryOrderRequest{
		Page:              1,
		PageSize:          10001, // Check if exceeds limit
		OrderNo:           req.OrderNo,
		UserID:            req.UserID,
		Status:            req.Status,
	}

	// Parse fulfillment status - convert from string to int8
	if req.FulfillmentStatus != "" {
		v, _ := strconv.ParseInt(req.FulfillmentStatus, 10, 8)
		queryReq.FulfillmentStatus = int8(v)
	}

	// Parse refund status - convert from string to int8
	if req.RefundStatus != "" {
		v, _ := strconv.ParseInt(req.RefundStatus, 10, 8)
		queryReq.RefundStatus = int8(v)
	}

	// Parse start time
	if req.StartTime != "" {
		t, err := time.Parse(time.RFC3339, req.StartTime)
		if err == nil {
			queryReq.StartTime = t
		}
	}

	// Parse end time
	if req.EndTime != "" {
		t, err := time.Parse(time.RFC3339, req.EndTime)
		if err == nil {
			queryReq.EndTime = t
		}
	}

	// Get orders for export
	rows, total, err := l.svcCtx.OrderFulfillmentApp.ExportOrders(l.ctx, shared.TenantID(tenantID), queryReq)
	if err != nil {
		return err
	}

	if total > 10000 {
		return code.ErrOrderExportLimitExceed
	}

	// Set response headers for CSV download
	l.w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	l.w.Header().Set("Content-Disposition", "attachment; filename=orders_"+time.Now().Format("20060102150405")+".csv")

	// Write UTF-8 BOM for Excel compatibility
	l.w.Write([]byte{0xEF, 0xBB, 0xBF})

	// Create CSV writer
	writer := csv.NewWriter(l.w)
	defer writer.Flush()

	// Write header
	header := []string{
		"订单号", "订单状态", "履约状态", "退款状态",
		"商品总额", "优惠金额", "运费", "实付金额",
		"收货人", "收货电话", "收货地址",
		"支付方式", "创建时间", "支付时间",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data rows
	for _, row := range rows {
		record := []string{
			row.OrderNo,
			row.Status,
			row.FulfillmentStatus,
			row.RefundStatus,
			formatDecimal(row.TotalAmount),
			formatDecimal(row.DiscountAmount),
			formatDecimal(row.ShippingFee),
			formatDecimal(row.PayAmount),
			row.ReceiverName,
			row.ReceiverPhone,
			row.ReceiverAddress,
			row.PaymentMethod,
			row.CreatedAt.Format(time.RFC3339),
			formatTimeForExport(row.PaidAt),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// formatTimeForExport formats time pointer to string for CSV export
func formatTimeForExport(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}