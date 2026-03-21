package tenant

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"gorm.io/gorm"
)

type Status int

const (
	StatusPending Status = iota
	StatusActive
	StatusSuspended
	StatusExpired
)

type PlanType int

const (
	PlanFree PlanType = iota
	PlanBasic
	PlanPro
	PlanEnterprise
)

type Tenant struct {
	ID           int64
	Name         string
	Code         string
	Status       Status
	Plan         PlanType
	Domain       string
	CustomDomain string
	Logo         string
	ContactName  string
	ContactPhone string
	ContactEmail string
	Address      string
	ExpireAt     *time.Time
	Audit        shared.AuditInfo `gorm:"embedded"`
}

func (t *Tenant) TableName() string {
	return "tenants"
}

func (t *Tenant) IsActive() bool {
	if t.Status != StatusActive {
		return false
	}
	if t.ExpireAt != nil && time.Now().After(*t.ExpireAt) {
		return false
	}
	return true
}

func (t *Tenant) Activate() {
	t.Status = StatusActive
}

func (t *Tenant) Suspend() error {
	if t.Status == StatusExpired {
		return code.ErrTenantCannotSuspendExpired
	}
	t.Status = StatusSuspended
	return nil
}

func (t *Tenant) Expire() {
	t.Status = StatusExpired
}

func (t *Tenant) GetEffectiveDomain() string {
	if t.CustomDomain != "" {
		return t.CustomDomain
	}
	return t.Domain
}

type Repository interface {
	Create(ctx context.Context, db *gorm.DB, tenant *Tenant) error
	Update(ctx context.Context, db *gorm.DB, tenant *Tenant) error
	FindByID(ctx context.Context, db *gorm.DB, id int64) (*Tenant, error)
	FindByCode(ctx context.Context, db *gorm.DB, code string) (*Tenant, error)
	FindByDomain(ctx context.Context, db *gorm.DB, domain string) (*Tenant, error)
	FindList(ctx context.Context, db *gorm.DB, query Query) ([]*Tenant, int64, error)
	ExistsByCode(ctx context.Context, db *gorm.DB, code string) (bool, error)
	ExistsByDomain(ctx context.Context, db *gorm.DB, domain string) (bool, error)
}

type Query struct {
	shared.PageQuery
	Name   string
	Code   string
	Status Status
	Plan   PlanType
}
