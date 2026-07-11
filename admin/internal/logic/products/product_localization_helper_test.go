package products

import (
	"testing"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
)

// TestToProductLocalizationResp_ReadsModelTime 验证序列化器读取的是
// application.Model.{CreatedAt,UpdatedAt}，而不是 shared.AuditInfo。
//
// ProductLocalization 同时嵌入 application.Model 和 shared.AuditInfo，但
// 持久化层 productLocalizationModel.toEntity 映射器只填充 Model，AuditInfo
// 保持零值。历史 bug：logic 读取 loc.AuditInfo.CreatedAt，导致 API 返回
// "0001-01-01T00:00:00Z"。修复后必须从 loc.Model 读取。
//
// 该测试模拟 GORM 从 DB 读取后的状态（Model 已填充、AuditInfo 为零值），
// 若有人再次把序列化改回读取 AuditInfo，此测试会失败。
func TestToProductLocalizationResp_ReadsModelTime(t *testing.T) {
	created := time.Date(2026, 7, 5, 19, 24, 38, 0, time.UTC)
	updated := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)

	cases := []struct {
		name string
		loc  *product.ProductLocalization
	}{
		{
			name: "model populated, audit zero (simulates DB load)",
			loc: &product.ProductLocalization{
				Model: application.Model{
					ID:        999,
					CreatedAt: created,
					UpdatedAt: updated,
				},
				TenantID:     1,
				ProductID:    42,
				LanguageCode: "en",
				Name:         "Widget",
				Description:  "A widget",
				// AuditInfo 故意留空——映射器不会填充它。
				AuditInfo: shared.AuditInfo{},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resp := toProductLocalizationResp(tc.loc)

			wantCreated := created.Format(time.RFC3339)
			wantUpdated := updated.Format(time.RFC3339)

			if resp.CreatedAt != wantCreated {
				t.Errorf("CreatedAt = %q, want %q (must read application.Model.CreatedAt, not AuditInfo)", resp.CreatedAt, wantCreated)
			}
			if resp.UpdatedAt != wantUpdated {
				t.Errorf("UpdatedAt = %q, want %q (must read application.Model.UpdatedAt, not AuditInfo)", resp.UpdatedAt, wantUpdated)
			}
		})
	}
}
