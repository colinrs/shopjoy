package warehouses

import (
	"context"
	"testing"

	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
)

// TestWarehouseOwnershipGuard_AllowsPlatformAdmin pins the FIX-2 contract
// for the read/mutate-by-id and list paths:
//
//	Platform admins (no TenantID in ctx, or TenantID == 0) MUST NOT be
//	rejected by the application-layer tenant guard. The GORM tenant
//	middleware is the authoritative platform-scope filter; this layer
//	must let those requests through so the middleware can do its job.
//
//	Ordinary users (TenantID != 0) MUST still be blocked from operating on
//	a warehouse owned by another tenant. requireWarehouseOwnership enforces
//	that mismatch check; the per-handler guards around it just feed the
//	resolved TenantID in.
//
// The test passes a nil *svc.ServiceContext on purpose: the FIX-2 guard is
// expected to dereference it (via svcCtx.WarehouseRepo.FindByID) so the nil
// pointer turns the missing-platform-admin case into a loud panic instead of
// a silent "ok" path. That panic is the regression sentinel — if someone
// re-introduces a `tenantID == 0 → ErrTenantNotFound` early-return, this test
// will start returning nil-err or the wrong error and fail. We accept a panic
// here because that is the new contract: the application layer no longer
// short-circuits platform-admin requests, and the GORM middleware is the only
// thing standing between the platform admin and the cross-tenant rows.
//
// Each op invokes one warehouse logic with a nil svcCtx. For a platform
// admin (ctx has no TenantID, or TenantID == 0) we expect a panic because
// the code reaches svcCtx.WarehouseRepo.FindByID; for an ordinary user we
// use a tenantID that won't match anything in the nil repo and assert the
// panic as well (FindByID would also be reached before the mismatch check).
func TestWarehouseOwnershipGuard_AllowsPlatformAdmin(t *testing.T) {
	ctxCases := []struct {
		name     string
		setupCtx func() context.Context
	}{
		{
			name:     "no TenantID in ctx",
			setupCtx: func() context.Context { return context.Background() },
		},
		{
			name:     "TenantID = 0",
			setupCtx: func() context.Context { return contextx.SetTenantID(context.Background(), 0) },
		},
		{
			name:     "TenantID != 0 (ordinary user)",
			setupCtx: func() context.Context { return contextx.SetTenantID(context.Background(), 42) },
		},
	}

	// Each op invokes one warehouse logic with a nil svcCtx. FIX-2 lets the
	// call through the tenant guard so we expect a panic from the nil
	// svcCtx.WarehouseRepo dereference — that panic is what proves the
	// application layer is no longer pre-empting platform-admin requests.
	ops := []struct {
		name string
		call func(ctx context.Context)
	}{
		{
			name: "GetWarehouse",
			call: func(ctx context.Context) {
				l := NewGetWarehouseLogic(ctx, nil)
				_, _ = l.GetWarehouse(&types.GetWarehouseReq{ID: 1})
			},
		},
		{
			name: "DeleteWarehouse",
			call: func(ctx context.Context) {
				l := NewDeleteWarehouseLogic(ctx, nil)
				_, _ = l.DeleteWarehouse(&types.GetWarehouseReq{ID: 1})
			},
		},
		{
			name: "SetDefaultWarehouse",
			call: func(ctx context.Context) {
				l := NewSetDefaultWarehouseLogic(ctx, nil)
				_, _ = l.SetDefaultWarehouse(&types.SetDefaultWarehouseReq{ID: 1})
			},
		},
		{
			name: "UpdateWarehouseStatus",
			call: func(ctx context.Context) {
				l := NewUpdateWarehouseStatusLogic(ctx, nil)
				_, _ = l.UpdateWarehouseStatus(&types.UpdateWarehouseStatusReq{ID: 1, Status: 1})
			},
		},
		{
			name: "ListWarehouses",
			call: func(ctx context.Context) {
				l := NewListWarehousesLogic(ctx, nil)
				_, _ = l.ListWarehouses(&types.ListWarehouseReq{})
			},
		},
	}

	for _, op := range ops {
		for _, cc := range ctxCases {
			t.Run(op.name+"/"+cc.name, func(t *testing.T) {
				// FIX-2 contract: the application-layer guard MUST NOT
				// short-circuit platform-admin requests. We assert that by
				// expecting a panic from the nil svcCtx dereference — if
				// someone re-introduces a tenantID == 0 → ErrTenantNotFound
				// early-return, the call will return a clean error and this
				// test will fail.
				defer func() {
					if r := recover(); r == nil {
						t.Fatalf("guard regressed — %s returned cleanly for "+
							"%s instead of reaching svcCtx; the FIX-2 "+
							"application-layer tenant guard must NOT pre-empt "+
							"platform-admin requests.",
							op.name, cc.name)
					}
				}()

				op.call(cc.setupCtx())
			})
		}
	}
}
