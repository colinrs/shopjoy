package transactions

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/points"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExportPointsTransactionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	w      http.ResponseWriter
	r      *http.Request
}

func NewExportPointsTransactionsLogic(ctx context.Context, svcCtx *svc.ServiceContext, w http.ResponseWriter, r *http.Request) ExportPointsTransactionsLogic {
	return ExportPointsTransactionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		w:      w,
		r:      r,
	}
}

func (l *ExportPointsTransactionsLogic) ExportPointsTransactions(req *types.ExportPointsTransactionsReq) error {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Build query
	query := points.PointsTransactionQuery{
		UserID: req.UserID,
	}

	// Parse type filter - convert "earn"/"redeem" to TransactionType
	if req.Type != "" {
		switch req.Type {
		case "earn":
			query.Type = points.TransactionTypeEarn
		case "redeem":
			query.Type = points.TransactionTypeRedeem
		}
	}

	// Parse start time
	if req.StartTime != "" {
		t, err := time.Parse(time.RFC3339, req.StartTime)
		if err == nil {
			query.StartTime = &t
		}
	}

	// Parse end time
	if req.EndTime != "" {
		t, err := time.Parse(time.RFC3339, req.EndTime)
		if err == nil {
			query.EndTime = &t
		}
	}

	// Get transactions for export
	transactions, total, err := l.svcCtx.PointsService.ExportTransactions(l.ctx, shared.TenantID(tenantID), query)
	if err != nil {
		return err
	}

	if total > 10000 {
		return code.ErrPointsTransactionExportLimitExceed
	}

	// Set response headers for CSV download
	l.w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	l.w.Header().Set("Content-Disposition", "attachment; filename=points_transactions_export_"+time.Now().Format("20060102")+".csv")

	// Write UTF-8 BOM for Excel compatibility
	if _, err := l.w.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
		return err
	}

	// Create CSV writer
	writer := csv.NewWriter(l.w)
	defer writer.Flush()

	// Write header
	header := []string{
		"Transaction ID", "User ID", "Type", "Points", "Balance After", "Description", "Created At",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data rows
	for _, t := range transactions {
		record := []string{
			formatInt64(t.ID),
			formatInt64(t.UserID),
			string(t.Type),
			formatInt64(t.Points),
			formatInt64(t.BalanceAfter),
			t.Description,
			t.CreatedAt.Format(time.RFC3339),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func formatInt64(n int64) string {
	return fmt.Sprintf("%d", n)
}
