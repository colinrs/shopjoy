package promotion

import (
	"github.com/colinrs/shopjoy/pkg/domain/shared"
)

// Query is the filter set for FindList. All optional fields are pointers
// because Kind/Status/Type are iota-based enums whose zero values are
// legitimate members (PROMOTION, ACTIVE, DISCOUNT) — using != 0 as a
// "filter set" sentinel would silently drop those.
type Query struct {
	shared.PageQuery
	TenantID    shared.TenantID
	Name        string
	Kind        *Kind
	Status      *Status
	Type        *Type
	MarketID    *int64
	ExpiredOnly bool
}
