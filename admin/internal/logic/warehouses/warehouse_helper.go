package warehouses

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
)

func toWarehouseDetailResp(w *product.Warehouse) *types.WarehouseDetailResp {
	return &types.WarehouseDetailResp{
		ID:        w.ID,
		Code:      w.Code,
		Name:      w.Name,
		Country:   w.Country,
		Address:   w.Address,
		IsDefault: w.IsDefault,
		Status:    int8(w.Status), // #nosec G115 // status values are small (tinyint range)
		CreatedAt: w.Model.CreatedAt.Format(time.RFC3339),
		UpdatedAt: w.Model.UpdatedAt.Format(time.RFC3339),
	}
}

// requireWarehouseOwnership resolves the caller's TenantID from ctx, loads the
// warehouse by id, and verifies the row belongs to that tenant. It is the
// single enforcement point for cross-tenant isolation on the warehouse
// read/mutate-by-id paths (get, delete, set-default, update-status), mirroring
// the inline guard already used by UpdateWarehouse (FU-2).
//
// Contract:
//   - No TenantID in ctx (AuthMiddleware bypassed) or TenantID == 0 (platform
//     admin without an explicit X-Tenant-ID header) → code.ErrTenantNotFound.
//     Operating without a concrete tenant would leak rows across the isolation
//     boundary, so both cases are rejected up front.
//   - Warehouse missing, soft-deleted, or owned by a different tenant →
//     code.ErrInventoryWarehouseNotFound. Cross-tenant access is deliberately
//     surfaced as "not found" (not "forbidden") so a tenant cannot probe the
//     existence of another tenant's warehouses.
//
// svcCtx is passed (rather than the repo/db directly) so the TenantID guard
// runs before any svcCtx field is dereferenced; this keeps the guard unit
// testable with a nil svcCtx, matching the Create/Update logic guard tests.
func requireWarehouseOwnership(ctx context.Context, svcCtx *svc.ServiceContext, id int64) (*product.Warehouse, error) {
	tenantID, ok := contextx.GetTenantID(ctx)
	if !ok || tenantID == 0 {
		return nil, code.ErrTenantNotFound
	}

	warehouse, err := svcCtx.WarehouseRepo.FindByID(ctx, svcCtx.DB, id)
	if err != nil {
		return nil, err
	}
	if warehouse == nil {
		return nil, code.ErrInventoryWarehouseNotFound
	}

	// Defense in depth: refuse cross-tenant access.
	if int64(warehouse.TenantID) != tenantID {
		return nil, code.ErrInventoryWarehouseNotFound
	}

	return warehouse, nil
}
