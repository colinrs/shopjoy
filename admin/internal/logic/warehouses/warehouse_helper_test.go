package warehouses

import (
	"context"
	"errors"
	"testing"

	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
)

// TestWarehouseOwnershipGuard_RequiresTenantID pins the FU-2.1 contract for the
// read/mutate-by-id and list paths: every handler that operates on a tenant's
// warehouse MUST resolve a concrete TenantID from ctx before touching the repo.
//
// When AuthMiddleware is bypassed (no TenantID) or the caller is a platform
// admin without an explicit X-Tenant-ID header (TenantID == 0), the logic MUST
// short-circuit with code.ErrTenantNotFound instead of loading or listing rows
// that belong to another tenant.
//
// We exercise only the early-return guard here (with a nil *svc.ServiceContext)
// because reaching the repo requires a fully wired svcCtx; the cross-tenant
// mismatch path (loaded warehouse owned by a different tenant) is covered by
// requireWarehouseOwnership's shared logic and integration tests. If the guard
// ever regresses and the logic dereferences svcCtx, the recover() below turns
// the nil-pointer panic into a loud test failure.
func TestWarehouseOwnershipGuard_RequiresTenantID(t *testing.T) {
	ctxCases := []struct {
		name     string
		setupCtx func() context.Context
		reason   string
	}{
		{
			name:     "no TenantID in ctx",
			setupCtx: func() context.Context { return context.Background() },
			reason:   "AuthMiddleware bypassed (no SetTenantID call)",
		},
		{
			name:     "TenantID = 0",
			setupCtx: func() context.Context { return contextx.SetTenantID(context.Background(), 0) },
			reason:   "platform admin without X-Tenant-ID header",
		},
	}

	// Each op invokes one warehouse logic with a nil svcCtx. All must return
	// code.ErrTenantNotFound before dereferencing svcCtx.
	ops := []struct {
		name string
		call func(ctx context.Context) error
	}{
		{
			name: "GetWarehouse",
			call: func(ctx context.Context) error {
				l := NewGetWarehouseLogic(ctx, nil)
				_, err := l.GetWarehouse(&types.GetWarehouseReq{ID: 1})
				return err
			},
		},
		{
			name: "DeleteWarehouse",
			call: func(ctx context.Context) error {
				l := NewDeleteWarehouseLogic(ctx, nil)
				_, err := l.DeleteWarehouse(&types.GetWarehouseReq{ID: 1})
				return err
			},
		},
		{
			name: "SetDefaultWarehouse",
			call: func(ctx context.Context) error {
				l := NewSetDefaultWarehouseLogic(ctx, nil)
				_, err := l.SetDefaultWarehouse(&types.SetDefaultWarehouseReq{ID: 1})
				return err
			},
		},
		{
			name: "UpdateWarehouseStatus",
			call: func(ctx context.Context) error {
				l := NewUpdateWarehouseStatusLogic(ctx, nil)
				_, err := l.UpdateWarehouseStatus(&types.UpdateWarehouseStatusReq{ID: 1, Status: 1})
				return err
			},
		},
		{
			name: "ListWarehouses",
			call: func(ctx context.Context) error {
				l := NewListWarehousesLogic(ctx, nil)
				_, err := l.ListWarehouses(&types.ListWarehouseReq{})
				return err
			},
		},
	}

	for _, op := range ops {
		for _, cc := range ctxCases {
			t.Run(op.name+"/"+cc.name, func(t *testing.T) {
				defer func() {
					if r := recover(); r != nil {
						t.Fatalf("guard regressed — %s reached svcCtx before "+
							"TenantID check (reason: %s). recover: %v",
							op.name, cc.reason, r)
					}
				}()

				err := op.call(cc.setupCtx())
				if err == nil {
					t.Fatalf("%s() = nil err, want errors.Is == %v (reason: %s)",
						op.name, code.ErrTenantNotFound, cc.reason)
				}
				if !errors.Is(err, code.ErrTenantNotFound) {
					t.Fatalf("%s() err = %v, want errors.Is == %v (reason: %s)",
						op.name, err, code.ErrTenantNotFound, cc.reason)
				}
			})
		}
	}
}
