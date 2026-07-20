package shipping_zones

import (
	"github.com/colinrs/shopjoy/admin/internal/domain/shipping"
	"github.com/colinrs/shopjoy/admin/internal/types"
)

// toStringI18n converts a wire []NameI18nEntry (array) into an entity
// StringI18n (map[locale]name). Empty/invalid entries are dropped.
//
// Exported as ToStringI18n so other logic packages (e.g. shipping_templates)
// can reuse it when building zone responses without duplicating logic.
func ToStringI18n(entries []types.NameI18nEntry) shipping.StringI18n {
	return toStringI18n(entries)
}

// FromStringI18n converts an entity StringI18n (map) back into a wire
// []NameI18nEntry (array). Used by response builders.
func FromStringI18n(s shipping.StringI18n) []types.NameI18nEntry {
	return fromStringI18n(s)
}
