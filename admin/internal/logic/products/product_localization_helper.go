package products

import (
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

// toProductLocalizationResp 序列化 ProductLocalization 为 API 响应。
//
// 时间字段必须读取 loc.Model.{CreatedAt,UpdatedAt}——这是持久化层
// productLocalizationModel.toEntity 映射器实际填充的字段。ProductLocalization
// 同时嵌入 application.Model 与 shared.AuditInfo，但映射器只填充 Model，
// AuditInfo 保持零值。禁止读取 loc.AuditInfo.CreatedAt（会得到 "0001-01-01"）。
func toProductLocalizationResp(loc *product.ProductLocalization) *types.ProductLocalizationResp {
	return &types.ProductLocalizationResp{
		ID:           loc.ID,
		ProductID:    loc.ProductID,
		LanguageCode: loc.LanguageCode,
		Name:         loc.Name,
		Description:  loc.Description,
		CreatedAt:    loc.Model.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    loc.Model.UpdatedAt.Format(time.RFC3339),
	}
}
