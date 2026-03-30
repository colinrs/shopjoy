package payments

import (
	"context"
	"encoding/csv"
	"net/http"
	"strconv"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/payment"
	appPayment "github.com/colinrs/shopjoy/admin/internal/application/payment"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExportPaymentTransactionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
	w      http.ResponseWriter
}

func NewExportPaymentTransactionsLogic(ctx context.Context, svcCtx *svc.ServiceContext, w http.ResponseWriter, r *http.Request) ExportPaymentTransactionsLogic {
	return ExportPaymentTransactionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		w:      w,
		r:      r,
	}
}

func (l *ExportPaymentTransactionsLogic) ExportPaymentTransactions(req *types.ExportPaymentTransactionsReq) error {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Build query request
	queryReq := appPayment.ListTransactionsRequest{
		Page:          1,
		PageSize:      10001, // Check if exceeds limit
		TransactionID: req.TransactionID,
		PaymentMethod: payment.PaymentMethod(req.PaymentMethod),
		Status:        payment.TransactionStatus(req.Status),
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

	// Get transactions for export
	result, err := l.svcCtx.PaymentService.ListTransactions(l.ctx, shared.TenantID(tenantID), queryReq)
	if err != nil {
		return err
	}

	if result.Total > 10000 {
		return code.ErrPaymentExportLimitExceed
	}

	// Set response headers for CSV download
	l.w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	l.w.Header().Set("Content-Disposition", "attachment; filename=payment_transactions_export_"+time.Now().Format("20060102")+".csv")

	// Write UTF-8 BOM for Excel compatibility
	l.w.Write([]byte{0xEF, 0xBB, 0xBF})

	// Create CSV writer
	writer := csv.NewWriter(l.w)
	defer writer.Flush()

	// Write header
	header := []string{
		"Transaction ID",
		"Order No",
		"Amount",
		"Currency",
		"Payment Method",
		"Status",
		"Paid At",
		"Created At",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data rows
	for _, txn := range result.List {
		record := []string{
			txn.TransactionID,
			txn.OrderNo,
			txn.Amount,
			txn.Currency,
			txn.PaymentMethodText,
			txn.StatusText,
			txn.PaidAt,
			txn.CreatedAt,
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

// formatDecimal formats decimal.Decimal to string for export
func formatDecimal(d decimal.Decimal) string {
	return d.String()
}

// parseTime parses time string to *time.Time
func parseTime(timeStr string) (*time.Time, error) {
	if timeStr == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// formatStatusText returns human-readable status text
func formatStatusText(status int8) string {
	switch status {
	case 0:
		return "pending"
	case 1:
		return "succeeded"
	case 2:
		return "failed"
	default:
		return "unknown"
	}
}

// parsePaymentMethod returns human-readable payment method text
func parsePaymentMethodText(method string) string {
	switch method {
	case "stripe":
		return "Stripe"
	case "alipay":
		return "Alipay"
	case "wechat":
		return "WeChat Pay"
	default:
		return method
	}
}

// formatOrderID formats order ID to string
func formatOrderID(orderID int64) string {
	return strconv.FormatInt(orderID, 10)
}
