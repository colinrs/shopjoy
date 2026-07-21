package shipping

import (
	"testing"

	"github.com/colinrs/shopjoy/pkg/application"
)

// newZone is a test helper that builds a ShippingZone with embedded Model ID.
// Go's shorthand promoted-field init ({ID: 1, ...}) doesn't compile against
// application.Model here, so we use explicit Model initialization.
func newZone(id int64, name string, regions Regions, sort int) *ShippingZone {
	return &ShippingZone{
		Model:   application.Model{ID: id},
		Name:    name,
		Regions: regions,
		Sort:    sort,
	}
}

func TestZoneMatcher_MatchByCountry(t *testing.T) {
	zones := []*ShippingZone{
		newZone(1, "US", Regions{"US"}, 10),
		newZone(2, "EU-DE", Regions{"DE"}, 5),
		newZone(3, "EU-FR", Regions{"FR"}, 5),
	}
	matcher := NewZoneMatcher(zones)

	got := matcher.Match(MatchInput{CountryCode: "US"})
	if got == nil || got.ID != 1 {
		t.Errorf("expected US zone, got %v", got)
	}
}

func TestZoneMatcher_MatchByProvince(t *testing.T) {
	zones := []*ShippingZone{
		newZone(1, "US-CA", Regions{"US-CA"}, 5),
		newZone(2, "US-TX", Regions{"US-TX"}, 5),
		newZone(3, "US-Default", Regions{"US"}, 100),
	}
	matcher := NewZoneMatcher(zones)

	got := matcher.Match(MatchInput{CountryCode: "US", ProvinceCode: "US-CA"})
	if got == nil || got.ID != 1 {
		t.Errorf("expected US-CA, got %v", got)
	}
}

func TestZoneMatcher_FallbackToCountry(t *testing.T) {
	zones := []*ShippingZone{
		newZone(1, "US-Default", Regions{"US"}, 100),
		newZone(2, "US-FL", Regions{"US-FL"}, 5),
	}
	matcher := NewZoneMatcher(zones)

	got := matcher.Match(MatchInput{CountryCode: "US", ProvinceCode: "US-WY"})
	if got == nil || got.ID != 1 {
		t.Errorf("expected US-Default fallback, got %v", got)
	}
}

func TestZoneMatcher_NoMatch(t *testing.T) {
	zones := []*ShippingZone{
		newZone(1, "US", Regions{"US"}, 100),
	}
	matcher := NewZoneMatcher(zones)
	got := matcher.Match(MatchInput{CountryCode: "JP"})
	if got != nil {
		t.Errorf("expected no match for JP, got %v", got)
	}
}