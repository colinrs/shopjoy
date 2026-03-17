package tenant

import (
	"context"
	"fmt"

	"github.com/colinrs/shopjoy/pkg/domain/shared"
)

type contextKey struct{}

var tenantContextKey = &contextKey{}

func WithContext(ctx context.Context, tenantID shared.TenantID) context.Context {
	return context.WithValue(ctx, tenantContextKey, tenantID)
}

func FromContext(ctx context.Context) (shared.TenantID, bool) {
	tenantID, ok := ctx.Value(tenantContextKey).(shared.TenantID)
	return tenantID, ok
}

func MustFromContext(ctx context.Context) shared.TenantID {
	tenantID, ok := FromContext(ctx)
	if !ok {
		panic("tenant ID not found in context")
	}
	return tenantID
}

type Provider interface {
	GetCurrentTenantID() shared.TenantID
}

type contextProvider struct {
	ctx context.Context
}

func NewContextProvider(ctx context.Context) Provider {
	return &contextProvider{ctx: ctx}
}

func (p *contextProvider) GetCurrentTenantID() shared.TenantID {
	return MustFromContext(p.ctx)
}

type Tenant struct {
	ID           shared.TenantID
	Name         string
	Code         string
	Status       shared.Status
	Domain       string
	ContactName  string
	ContactPhone string
	ContactEmail string
	Audit        shared.AuditInfo
}

func (t *Tenant) IsActive() bool {
	return t.Status == shared.StatusEnabled
}

func (t *Tenant) Validate() error {
	if !t.ID.IsValid() {
		return shared.ErrInvalidTenantID
	}
	if t.Name == "" {
		return fmt.Errorf("tenant name is required")
	}
	if t.Code == "" {
		return fmt.Errorf("tenant code is required")
	}
	return nil
}

type Repository interface {
	FindByID(ctx context.Context, id shared.TenantID) (*Tenant, error)
	FindByCode(ctx context.Context, code string) (*Tenant, error)
	FindByDomain(ctx context.Context, domain string) (*Tenant, error)
	Save(ctx context.Context, tenant *Tenant) error
	Update(ctx context.Context, tenant *Tenant) error
}
