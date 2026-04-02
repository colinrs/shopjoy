package inventory

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type ExportInventoryLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	w      http.ResponseWriter
	r      *http.Request
}

func NewExportInventoryLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext, w http.ResponseWriter, r *http.Request) ExportInventoryLogsLogic {
	return ExportInventoryLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		w:      w,
		r:      r,
	}
}

func (l *ExportInventoryLogsLogic) ExportInventoryLogs(req *types.GetInventoryLogsReq) error {
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// For export, use a large page size to get all records
	query := product.InventoryLogQuery{
		PageQuery: shared.PageQuery{
			Page:     1,
			PageSize: 10000, // Limit for export
		},
		SKUCode:    req.SKUCode,
		ChangeType: req.Type,
	}

	var logs []*product.InventoryLog
	var total int64
	var err error

	if req.ProductID > 0 {
		logs, total, err = l.svcCtx.InventoryLogRepo.FindByProduct(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ProductID, query)
	} else if req.SKUCode != "" {
		logs, total, err = l.svcCtx.InventoryLogRepo.FindBySKU(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.SKUCode, query)
	} else {
		// No filter - get all (limited)
		logs, total, err = l.svcCtx.InventoryLogRepo.FindAll(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), query)
	}

	if err != nil {
		return err
	}

	if total > 10000 {
		total = 10000 // Cap at export limit
	}

	// Set response headers for CSV download
	l.w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	l.w.Header().Set("Content-Disposition", "attachment; filename=inventory_logs_export_"+time.Now().Format("20060102")+".csv")

	// Write UTF-8 BOM for Excel compatibility
	if _, err := l.w.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
		return err
	}

	// Create CSV writer
	writer := csv.NewWriter(l.w)
	defer writer.Flush()

	// Write header
	header := []string{
		"ID", "SKU Code", "Product ID", "Warehouse ID", "Change Type",
		"Change Quantity", "Before Stock", "After Stock", "Order No",
		"Remark", "Operator ID", "Created At",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data rows
	for _, log := range logs {
		record := []string{
			fmt.Sprintf("%d", log.ID),
			log.SKUCode,
			fmt.Sprintf("%d", log.ProductID),
			fmt.Sprintf("%d", log.WarehouseID),
			log.ChangeType,
			fmt.Sprintf("%d", log.ChangeQuantity),
			fmt.Sprintf("%d", log.BeforeStock),
			fmt.Sprintf("%d", log.AfterStock),
			log.OrderNo,
			log.Remark,
			fmt.Sprintf("%d", log.OperatorID),
			log.CreatedAt.Format(time.RFC3339),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}
