package products

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExportProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
	w      http.ResponseWriter
}

func NewExportProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext, w http.ResponseWriter, r *http.Request) ExportProductsLogic {
	return ExportProductsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		w:      w,
		r:      r,
	}
}

func (l *ExportProductsLogic) ExportProducts(req *types.ExportProductsReq) error {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Build query request - use large page size for export
	query := product.Query{
		TenantID:   shared.TenantID(tenantID),
		Name:       req.Name,
		CategoryID: req.CategoryID,
		Status:     parseProductStatus(req.Status),
		MarketID:   req.MarketID,
		Page:       1,
		PageSize:   10001, // Fetch one more to check if limit exceeded
	}

	// Parse price range
	if req.MinPrice != "" {
		if minPrice, err := decimal.NewFromString(req.MinPrice); err == nil {
			query.MinPrice = &minPrice
		}
	}
	if req.MaxPrice != "" {
		if maxPrice, err := decimal.NewFromString(req.MaxPrice); err == nil {
			query.MaxPrice = &maxPrice
		}
	}

	// Get products for export
	products, total, err := l.svcCtx.ProductRepo.FindList(l.ctx, l.svcCtx.DB, query)
	if err != nil {
		return err
	}

	if total > 10000 {
		return code.ErrProductExportLimitExceed
	}

	// Set response headers for CSV download
	l.w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	l.w.Header().Set("Content-Disposition", "attachment; filename=products_"+time.Now().Format("20060102150405")+".csv")

	// Write UTF-8 BOM for Excel compatibility
	l.w.Write([]byte{0xEF, 0xBB, 0xBF})

	// Create CSV writer
	writer := csv.NewWriter(l.w)
	defer writer.Flush()

	// Write header
	header := []string{
		"ID", "SKU", "名称", "描述", "价格", "成本价", "币种",
		"库存", "状态", "分类ID", "品牌",
		"标签", "图片", "是否矩阵产品",
		"HS编码", "原产国", "重量", "重量单位",
		"长度", "宽度", "高度", "危险品",
		"创建时间", "更新时间",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data rows
	for _, p := range products {
		record := []string{
			fmt.Sprintf("%d", p.ID),
			sanitizeCSVField(p.SKU),
			sanitizeCSVField(p.Name),
			sanitizeCSVField(p.Description),
			p.Price.Amount.String(),
			p.CostPrice.Amount.String(),
			p.Price.Currency,
			fmt.Sprintf("%d", p.Stock),
			p.Status.String(),
			fmt.Sprintf("%d", p.CategoryID),
			sanitizeCSVField(p.Brand),
			sanitizeCSVField(joinTags(p.Tags)),
			sanitizeCSVField(joinImages(p.Images)),
			fmt.Sprintf("%v", p.IsMatrixProduct),
			sanitizeCSVField(p.HSCode),
			sanitizeCSVField(p.COO),
			p.Weight.String(),
			sanitizeCSVField(p.WeightUnit),
			p.Dimensions.Length.String(),
			p.Dimensions.Width.String(),
			p.Dimensions.Height.String(),
			sanitizeCSVField(joinDangerousGoods(p.DangerousGoods)),
			p.CreatedAt.Format(time.RFC3339),
			p.UpdatedAt.Format(time.RFC3339),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// parseProductStatus parses status string to product.Status
func parseProductStatus(s string) *product.Status {
	switch s {
	case "draft":
		status := product.StatusDraft
		return &status
	case "on_sale":
		status := product.StatusOnSale
		return &status
	case "off_sale":
		status := product.StatusOffSale
		return &status
	case "deleted":
		status := product.StatusDeleted
		return &status
	default:
		return nil
	}
}

// joinTags converts tags slice to JSON string for CSV export
func joinTags(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	b, _ := json.Marshal(tags)
	return string(b)
}

// joinImages converts images slice to JSON string for CSV export
func joinImages(images []string) string {
	if len(images) == 0 {
		return ""
	}
	b, _ := json.Marshal(images)
	return string(b)
}

// joinDangerousGoods converts dangerous goods slice to JSON string for CSV export
func joinDangerousGoods(dg []string) string {
	if len(dg) == 0 {
		return ""
	}
	b, _ := json.Marshal(dg)
	return string(b)
}

// sanitizeCSVField sanitizes a CSV field to prevent CSV injection attacks.
// If the field starts with formula characters (=, +, -, @, \t), it prefixes with a single quote.
func sanitizeCSVField(field string) string {
	if len(field) == 0 {
		return field
	}
	firstChar := field[0]
	// Characters that could trigger formula execution in Excel
	if firstChar == '=' || firstChar == '+' || firstChar == '-' || firstChar == '@' || firstChar == '\t' {
		return "'" + field
	}
	return field
}
