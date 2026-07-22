package warehouses

import (
	"context"
	"errors"
	"testing"

	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
)

// TestCreateWarehouseLogic_RequiresTenantID pins the FU-2 contract:
//
//	When the AuthMiddleware is bypassed (or runs without injecting a
//	TenantID) the CreateWarehouse logic MUST short-circuit with
//	code.ErrTenantNotFound instead of writing a row with tenant_id = 0.
//
//	This is the #1 cross-tenant isolation bug class identified in
//	the Phase 8 review. The earlier implementation did not call
//	contextx.GetTenantID and silently persisted tenant_id = 0, which
//	made every tenant's warehouses visible to every other tenant.
//
// We exercise only the early-return guard here because the rest of
//	CreateWarehouse requires a fully wired *svc.ServiceContext
//	(WarehouseRepo, IDGen, DB) and is covered by handler-level tests.
func TestCreateWarehouseLogic_RequiresTenantID(t *testing.T) {
	cases := []struct {
		name       string
		setupCtx   func() context.Context
		wantErrIs  error
		wantReason string
	}{
		{
			name: "no TenantID in ctx → ErrTenantNotFound",
			setupCtx: func() context.Context {
				return context.Background()
			},
			wantErrIs:  code.ErrTenantNotFound,
			wantReason: "AuthMiddleware bypassed (no SetTenantID call)",
		},
		{
			name: "TenantID = 0 → ErrTenantNotFound",
			setupCtx: func() context.Context {
				return contextx.SetTenantID(context.Background(), 0)
			},
			wantErrIs:  code.ErrTenantNotFound,
			wantReason: "platform admin without X-Tenant-ID header",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// We don't need a real svcCtx because the guard short-circuits
			// before any field is touched. Passing nil here would still
			// dereference inside CreateWarehouse if we got past the guard;
			// the guard is what we're verifying.
			//
			// To keep the test honest, also confirm a panic-recovery path
			// does not silently absorb the nil svcCtx: if the guard ever
			// regresses and the logic reaches l.svcCtx.WarehouseRepo,
			// the test will fail loudly instead of returning the wrong
			// error. This is intentional — we WANT regression to be
			// obvious, not silently passing.
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("guard regressed — CreateWarehouse reached "+
						"svcCtx before TenantID check (reason: %s). "+
						"recover: %v", tc.wantReason, r)
				}
			}()

			l := NewCreateWarehouseLogic(tc.setupCtx(), nil)

			_, err := l.CreateWarehouse(nil) // req can be nil: guard runs first
			if err == nil {
				t.Fatalf("CreateWarehouse() = nil err, want errors.Is == %v "+
					"(reason: %s)", tc.wantErrIs, tc.wantReason)
			}
			if !errors.Is(err, tc.wantErrIs) {
				t.Fatalf("CreateWarehouse() err = %v, want errors.Is == %v "+
					"(reason: %s)", err, tc.wantErrIs, tc.wantReason)
			}
		})
	}
}

// TestUpdateWarehouseLogic_RequiresTenantID pins the FU-2 contract for the
// update path: the caller MUST have a TenantID in ctx, and the loaded
// warehouse MUST belong to that TenantID (cross-tenant updates are
// rejected with ErrInventoryWarehouseNotFound).
//
// We exercise only the early-return guard for symmetry with the create
// test. The cross-tenant mismatch path requires a working
//*svc.ServiceContext and is covered by integration tests.
func TestUpdateWarehouseLogic_RequiresTenantID(t *testing.T) {
	cases := []struct {
		name       string
		setupCtx   func() context.Context
		wantErrIs  error
		wantReason string
	}{
		{
			name: "no TenantID in ctx → ErrTenantNotFound",
			setupCtx: func() context.Context {
				return context.Background()
			},
			wantErrIs:  code.ErrTenantNotFound,
			wantReason: "AuthMiddleware bypassed",
		},
		{
			name: "TenantID = 0 → ErrTenantNotFound",
			setupCtx: func() context.Context {
				return contextx.SetTenantID(context.Background(), 0)
			},
			wantErrIs:  code.ErrTenantNotFound,
			wantReason: "platform admin without X-Tenant-ID header",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("guard regressed — UpdateWarehouse reached "+
						"svcCtx before TenantID check (reason: %s). "+
						"recover: %v", tc.wantReason, r)
				}
			}()

			l := NewUpdateWarehouseLogic(tc.setupCtx(), nil)

			_, err := l.UpdateWarehouse(nil) // req can be nil: guard runs first
			if err == nil {
				t.Fatalf("UpdateWarehouse() = nil err, want errors.Is == %v "+
					"(reason: %s)", tc.wantErrIs, tc.wantReason)
			}
			if !errors.Is(err, tc.wantErrIs) {
				t.Fatalf("UpdateWarehouse() err = %v, want errors.Is == %v "+
					"(reason: %s)", err, tc.wantErrIs, tc.wantReason)
			}
		})
	}
}