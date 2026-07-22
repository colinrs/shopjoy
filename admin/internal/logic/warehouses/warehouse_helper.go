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
//   - No TenantID in ctx (AuthMiddleware bypassed) → defaults to 0 so the
//     tenant-scope check below becomes a no-op for unauthenticated callers.
//     We deliberately do NOT return code.ErrTenantNotFound here; the GORM
//     tenant middleware is the authoritative platform-scope filter, and this
//     helper only needs to block ordinary users from touching another tenant's
//     warehouse.
//   - TenantID != 0 (ordinary user) AND warehouse owned by a different tenant
//     → code.ErrInventoryWarehouseNotFound. Cross-tenant access is deliberately
//     surfaced as "not found" (not "forbidden") so a tenant cannot probe the
//     existence of another tenant's warehouses.
//   - Warehouse missing or soft-deleted → code.ErrInventoryWarehouseNotFound.
//
// svcCtx is passed (rather than the repo/db directly) so the TenantID guard
// runs before any svcCtx field is dereferenced; this keeps the guard unit
// testable with a nil svcCtx, matching the Create/Update logic guard tests.
func requireWarehouseOwnership(ctx context.Context, svcCtx *svc.ServiceContext, id int64) (*product.Warehouse, error) {
	tenantID, _ := contextx.GetTenantID(ctx)

	warehouse, err := svcCtx.WarehouseRepo.FindByID(ctx, svcCtx.DB, id)
	if err != nil {
		return nil, err
	}
	if warehouse == nil {
		return nil, code.ErrInventoryWarehouseNotFound
	}

	// Defense in depth: refuse cross-tenant access for ordinary users.
	// FIX-2: a platform admin (tenantID == 0) is allowed through this check
	// — the GORM tenant middleware enforces the broader platform-scope
	// invariant, so this layer only blocks a tenant from touching another
	// tenant's row.
	if tenantID != 0 && int64(warehouse.TenantID) != tenantID {
		return nil, code.ErrInventoryWarehouseNotFound
	}

	return warehouse, nil
}
