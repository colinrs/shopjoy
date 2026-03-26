package warehouses

import (
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

func toWarehouseDetailResp(w *product.Warehouse) *types.WarehouseDetailResp {
	return &types.WarehouseDetailResp{
		ID:        w.ID,
		Code:      w.Code,
		Name:      w.Name,
		Country:   w.Country,
		Address:   w.Address,
		IsDefault: w.IsDefault,
		Status:    int8(w.Status),
		CreatedAt: time.Unix(w.Audit.CreatedAt, 0).Format(time.RFC3339),
		UpdatedAt: time.Unix(w.Audit.UpdatedAt, 0).Format(time.RFC3339),
	}
}
