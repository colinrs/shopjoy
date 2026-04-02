package refunds

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"time"

	appfulfillment "github.com/colinrs/shopjoy/admin/internal/application/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExportRefundsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	w      http.ResponseWriter
	r      *http.Request
}

func NewExportRefundsLogic(ctx context.Context, svcCtx *svc.ServiceContext, w http.ResponseWriter, r *http.Request) ExportRefundsLogic {
	return ExportRefundsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		w:      w,
		r:      r,
	}
}

func (l *ExportRefundsLogic) ExportRefunds(req *types.ExportRefundsReq) error {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return err
	}

	// Parse time filters
	var startTime, endTime time.Time
	if req.StartTime != "" {
		startTime, _ = parseTime(req.StartTime)
	}
	if req.EndTime != "" {
		endTime, _ = parseTime(req.EndTime)
	}

	// Build query request - use large page size for export
	// Note: OrderNo is the external order number (string) and cannot be used for OrderID (int64) filtering
	queryReq := appfulfillment.QueryRefundRequest{
		Page:       1,
		PageSize:   10001, // Check if exceeds limit
		RefundNo:   req.RefundNo,
		Status:     fulfillment.ParseRefundStatus(req.Status),
		ReasonType: req.ReasonType,
		StartTime:  startTime,
		EndTime:    endTime,
	}

	// Get refunds for export
	listResp, err := l.svcCtx.RefundApp.ListRefunds(l.ctx, shared.TenantID(tenantID), queryReq)
	if err != nil {
		return err
	}

	if listResp.Total > 10000 {
		return code.ErrOrderExportLimitExceed
	}

	// Set response headers for CSV download
	l.w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	l.w.Header().Set("Content-Disposition", "attachment; filename=refunds_export_"+time.Now().Format("20060102")+".csv")

	// Write UTF-8 BOM for Excel compatibility
	if _, err := l.w.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
		return err
	}

	// Create CSV writer
	writer := csv.NewWriter(l.w)
	defer writer.Flush()

	// Write header
	header := []string{
		"Refund No", "Order No", "Amount", "Currency", "Reason", "Reason Type",
		"Status", "Applied At", "Processed At", "Approved At", "Rejected Reason",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data rows
	for _, r := range listResp.List {
		record := []string{
			r.RefundNo,
			fmt.Sprintf("%d", r.OrderID),
			r.Amount,
			r.Currency,
			r.Reason,
			r.ReasonType,
			r.StatusText,
			r.CreatedAt,
			utils.FormatTimeStringForExport(r.CompletedAt),
			utils.FormatTimeStringForExport(r.ApprovedAt),
			r.RejectReason,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}
