package tenant

import (
	"fmt"
	"net/http"

	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	HeaderTenantID   = "X-Tenant-ID"
	HeaderTenantCode = "X-Tenant-Code"
)

type Middleware struct {
	tenantRepo Repository
}

func NewMiddleware(tenantRepo Repository) *Middleware {
	return &Middleware{tenantRepo: tenantRepo}
}

func (m *Middleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tenantIDStr := r.Header.Get(HeaderTenantID)
		tenantCode := r.Header.Get(HeaderTenantCode)

		if tenantIDStr == "" && tenantCode == "" {
			logx.WithContext(ctx).Errorf("tenant ID or code is required")
			http.Error(w, "tenant required", http.StatusBadRequest)
			return
		}

		var tenantID shared.TenantID
		if tenantIDStr != "" {
			var id int64
			if _, err := fmt.Sscanf(tenantIDStr, "%d", &id); err != nil {
				logx.WithContext(ctx).Errorf("invalid tenant ID: %s", tenantIDStr)
				http.Error(w, "invalid tenant ID", http.StatusBadRequest)
				return
			}
			tenantID = shared.TenantID(id)
		} else {
			tenant, err := m.tenantRepo.FindByCode(ctx, tenantCode)
			if err != nil {
				logx.WithContext(ctx).Errorf("tenant not found: %s", tenantCode)
				http.Error(w, "tenant not found", http.StatusNotFound)
				return
			}
			tenantID = tenant.ID
		}

		if !tenantID.IsValid() {
			logx.WithContext(ctx).Errorf("invalid tenant ID")
			http.Error(w, "invalid tenant", http.StatusBadRequest)
			return
		}

		ctx = WithContext(ctx, tenantID)
		next(w, r.WithContext(ctx))
	}
}
