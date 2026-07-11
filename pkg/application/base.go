package application

import (
	"context"

	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
)

type CommandHandler[C any, R any] interface {
	Handle(ctx context.Context, cmd C) (R, error)
}

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, query Q) (R, error)
}

type UnitOfWork interface {
	Begin() error
	Commit() error
	Rollback() error
}

func GetTenantID(ctx context.Context) (shared.TenantID, error) {
	tenantID, ok := contextx.GetTenantIDValueObject(ctx)
	if !ok {
		return 0, shared.ErrInvalidTenantID
	}
	return tenantID, nil
}

func MustGetTenantID(ctx context.Context) shared.TenantID {
	tenantID, ok := contextx.GetTenantIDValueObject(ctx)
	if !ok {
		panic("tenant ID not found in context")
	}
	return tenantID
}
