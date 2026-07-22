package warehouses

import (
	"context"
	"testing"

	"github.com/colinrs/shopjoy/pkg/contextx"
)

// TestCreateWarehouseLogic_AllowsPlatformAdmin pins the FIX-2 contract for
// the create path: a platform admin (no TenantID in ctx, or TenantID == 0)
// MUST NOT be rejected by the application-layer tenant guard. The GORM
// tenant middleware is responsible for platform-scope filtering, so this
// layer must let those requests through.
//
// We assert that by expecting a panic from the nil svcCtx dereference — if
// someone re-introduces a tenantID == 0 → ErrTenantNotFound early-return,
// the call will return a clean error and this test will fail. That panic is
// the regression sentinel: the application layer no longer pre-empts
// platform-admin requests, and the GORM middleware is the only thing
// standing between the platform admin and the cross-tenant rows.
func TestCreateWarehouseLogic_AllowsPlatformAdmin(t *testing.T) {
	cases := []struct {
		name     string
		setupCtx func() context.Context
	}{
		{
			name: "no TenantID in ctx (AuthMiddleware bypassed)",
			setupCtx: func() context.Context {
				return context.Background()
			},
		},
		{
			name: "TenantID = 0 (platform admin without X-Tenant-ID)",
			setupCtx: func() context.Context {
				return contextx.SetTenantID(context.Background(), 0)
			},
		},
		{
			name: "TenantID != 0 (ordinary user)",
			setupCtx: func() context.Context {
				return contextx.SetTenantID(context.Background(), 42)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Fatalf("guard regressed — CreateWarehouse returned cleanly "+
						"for %s instead of reaching svcCtx; the FIX-2 "+
						"application-layer tenant guard must NOT pre-empt "+
						"platform-admin requests.", tc.name)
				}
			}()

			l := NewCreateWarehouseLogic(tc.setupCtx(), nil)
			_, _ = l.CreateWarehouse(nil) // req can be nil: we expect to panic before it is read
		})
	}
}

// TestUpdateWarehouseLogic_AllowsPlatformAdmin pins the FIX-2 contract for
// the update path: a platform admin (no TenantID in ctx, or TenantID == 0)
// MUST NOT be rejected by the application-layer tenant guard. The
// cross-tenant check that follows is intentionally skipped for tenantID == 0
// so platform admins can manage warehouses they own; ordinary users
// (TenantID != 0) still cannot touch another tenant's warehouse because the
// mismatch check below is unconditional for them.
func TestUpdateWarehouseLogic_AllowsPlatformAdmin(t *testing.T) {
	cases := []struct {
		name     string
		setupCtx func() context.Context
	}{
		{
			name: "no TenantID in ctx (AuthMiddleware bypassed)",
			setupCtx: func() context.Context {
				return context.Background()
			},
		},
		{
			name: "TenantID = 0 (platform admin without X-Tenant-ID)",
			setupCtx: func() context.Context {
				return contextx.SetTenantID(context.Background(), 0)
			},
		},
		{
			name: "TenantID != 0 (ordinary user)",
			setupCtx: func() context.Context {
				return contextx.SetTenantID(context.Background(), 42)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Fatalf("guard regressed — UpdateWarehouse returned cleanly "+
						"for %s instead of reaching svcCtx; the FIX-2 "+
						"application-layer tenant guard must NOT pre-empt "+
						"platform-admin requests.", tc.name)
				}
			}()

			l := NewUpdateWarehouseLogic(tc.setupCtx(), nil)
			_, _ = l.UpdateWarehouse(nil) // req can be nil: we expect to panic before it is read
		})
	}
}
