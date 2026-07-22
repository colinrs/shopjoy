package shipping

import (
	"testing"

	"github.com/shopspring/decimal"
)

// TestSurchargeCalculator_RemoteAreaMatch verifies the remote-area path:
// a 5-digit Alaska postal code (99501) matches the pattern "^9[0-9]{4}$",
// so RemoteSurcharge (50) is applied. The fuel surcharge is 15% of a 20
// base fee (=3), giving Remote=50, Fuel=3, Total=53.
func TestSurchargeCalculator_RemoteAreaMatch(t *testing.T) {
	zone := &ShippingZone{
		RemoteSurcharge:   decimal.RequireFromString("50.00"),
		RemoteZipPatterns: StringArray{`^9[0-9]{4}$`},
		FuelSurchargePct:  decimal.RequireFromString("0.15"),
	}

	in := SurchargeInput{
		BaseFee:    decimal.RequireFromString("20.00"),
		PostalCode: "99501",
	}

	got := CalculateSurcharges(zone, in)

	if !got.Remote.Equal(decimal.RequireFromString("50")) {
		t.Errorf("expected Remote=50, got %s", got.Remote)
	}
	if !got.Fuel.Equal(decimal.RequireFromString("3")) {
		t.Errorf("expected Fuel=3, got %s", got.Fuel)
	}
	if !got.Total.Equal(decimal.RequireFromString("53")) {
		t.Errorf("expected Total=53, got %s", got.Total)
	}
}

// TestSurchargeCalculator_NoMatch verifies that a non-matching postal code
// (NY 10001) against "^9[0-9]{4}$" yields no remote surcharge. Fuel still
// applies (15% of 20 = 3).
func TestSurchargeCalculator_NoMatch(t *testing.T) {
	zone := &ShippingZone{
		RemoteSurcharge:   decimal.RequireFromString("50.00"),
		RemoteZipPatterns: StringArray{`^9[0-9]{4}$`},
		FuelSurchargePct:  decimal.RequireFromString("0.15"),
	}

	in := SurchargeInput{
		BaseFee:    decimal.RequireFromString("20.00"),
		PostalCode: "10001",
	}

	got := CalculateSurcharges(zone, in)

	if !got.Remote.IsZero() {
		t.Errorf("expected Remote=0, got %s", got.Remote)
	}
	if !got.Fuel.Equal(decimal.RequireFromString("3")) {
		t.Errorf("expected Fuel=3, got %s", got.Fuel)
	}
	if !got.Total.Equal(decimal.RequireFromString("3")) {
		t.Errorf("expected Total=3, got %s", got.Total)
	}
}

// TestSurchargeCalculator_EmptyPatterns verifies that when the zone has no
// remote zip patterns, no remote surcharge is added (even if RemoteSurcharge
// is non-zero and a postal code is supplied).
func TestSurchargeCalculator_EmptyPatterns(t *testing.T) {
	zone := &ShippingZone{
		RemoteSurcharge:   decimal.RequireFromString("50.00"),
		RemoteZipPatterns: nil,
		FuelSurchargePct:  decimal.RequireFromString("0.15"),
	}

	in := SurchargeInput{
		BaseFee:    decimal.RequireFromString("20.00"),
		PostalCode: "99501",
	}

	got := CalculateSurcharges(zone, in)

	if !got.Remote.IsZero() {
		t.Errorf("expected Remote=0 (no patterns), got %s", got.Remote)
	}
	if !got.Fuel.Equal(decimal.RequireFromString("3")) {
		t.Errorf("expected Fuel=3, got %s", got.Fuel)
	}
	if !got.Total.Equal(decimal.RequireFromString("3")) {
		t.Errorf("expected Total=3, got %s", got.Total)
	}
}

// TestSurchargeCalculator_ZeroBaseFee verifies that a zero base fee yields
// a zero fuel surcharge, even when the percentage is non-zero.
func TestSurchargeCalculator_ZeroBaseFee(t *testing.T) {
	zone := &ShippingZone{
		RemoteSurcharge:   decimal.RequireFromString("50.00"),
		RemoteZipPatterns: StringArray{`^9[0-9]{4}$`},
		FuelSurchargePct:  decimal.RequireFromString("0.15"),
	}

	in := SurchargeInput{
		BaseFee:    decimal.Zero,
		PostalCode: "99501",
	}

	got := CalculateSurcharges(zone, in)

	if !got.Remote.Equal(decimal.RequireFromString("50")) {
		t.Errorf("expected Remote=50, got %s", got.Remote)
	}
	if !got.Fuel.IsZero() {
		t.Errorf("expected Fuel=0 (zero base), got %s", got.Fuel)
	}
	if !got.Total.Equal(decimal.RequireFromString("50")) {
		t.Errorf("expected Total=50, got %s", got.Total)
	}
}

// TestSurchargeCalculator_EmptyPostalCode verifies that an empty postal code
// produces no remote surcharge (no patterns match).
func TestSurchargeCalculator_EmptyPostalCode(t *testing.T) {
	zone := &ShippingZone{
		RemoteSurcharge:   decimal.RequireFromString("50.00"),
		RemoteZipPatterns: StringArray{`^9[0-9]{4}$`},
		FuelSurchargePct:  decimal.RequireFromString("0.15"),
	}

	in := SurchargeInput{
		BaseFee:    decimal.RequireFromString("20.00"),
		PostalCode: "",
	}

	got := CalculateSurcharges(zone, in)

	if !got.Remote.IsZero() {
		t.Errorf("expected Remote=0 (empty postcode), got %s", got.Remote)
	}
	if !got.Fuel.Equal(decimal.RequireFromString("3")) {
		t.Errorf("expected Fuel=3, got %s", got.Fuel)
	}
	if !got.Total.Equal(decimal.RequireFromString("3")) {
		t.Errorf("expected Total=3, got %s", got.Total)
	}
}

// TestSurchargeCalculator_FirstPatternMatches verifies that when multiple
// patterns are configured, the FIRST matching pattern applies (subsequent
// patterns are skipped).
func TestSurchargeCalculator_FirstPatternMatches(t *testing.T) {
	zone := &ShippingZone{
		RemoteSurcharge:   decimal.RequireFromString("25.00"),
		RemoteZipPatterns: StringArray{`^99[0-9]{3}$`, `^9[0-9]{4}$`},
		FuelSurchargePct:  decimal.Zero,
	}

	in := SurchargeInput{
		BaseFee:    decimal.RequireFromString("20.00"),
		PostalCode: "99501",
	}

	got := CalculateSurcharges(zone, in)

	if !got.Remote.Equal(decimal.RequireFromString("25")) {
		t.Errorf("expected Remote=25 (first match), got %s", got.Remote)
	}
	if !got.Total.Equal(decimal.RequireFromString("25")) {
		t.Errorf("expected Total=25, got %s", got.Total)
	}
}

// TestSurchargeCalculator_InvalidPatternSkipped verifies that a pattern that
// fails to compile is silently skipped (continue), allowing the next valid
// pattern to match.
func TestSurchargeCalculator_InvalidPatternSkipped(t *testing.T) {
	zone := &ShippingZone{
		RemoteSurcharge:   decimal.RequireFromString("40.00"),
		RemoteZipPatterns: StringArray{`[invalid`, `^9[0-9]{4}$`},
		FuelSurchargePct:  decimal.Zero,
	}

	in := SurchargeInput{
		BaseFee:    decimal.RequireFromString("20.00"),
		PostalCode: "99501",
	}

	got := CalculateSurcharges(zone, in)

	if !got.Remote.Equal(decimal.RequireFromString("40")) {
		t.Errorf("expected Remote=40 (skip bad pattern), got %s", got.Remote)
	}
	if !got.Total.Equal(decimal.RequireFromString("40")) {
		t.Errorf("expected Total=40, got %s", got.Total)
	}
}

// TestSurchargeCalculator_ZeroRemoteFee verifies that a zero RemoteSurcharge
// on the zone produces zero remote even when a postal code matches.
func TestSurchargeCalculator_ZeroRemoteFee(t *testing.T) {
	zone := &ShippingZone{
		RemoteSurcharge:   decimal.Zero,
		RemoteZipPatterns: StringArray{`^9[0-9]{4}$`},
		FuelSurchargePct:  decimal.RequireFromString("0.10"),
	}

	in := SurchargeInput{
		BaseFee:    decimal.RequireFromString("100.00"),
		PostalCode: "99501",
	}

	got := CalculateSurcharges(zone, in)

	if !got.Remote.IsZero() {
		t.Errorf("expected Remote=0 (zero zone fee), got %s", got.Remote)
	}
	if !got.Fuel.Equal(decimal.RequireFromString("10")) {
		t.Errorf("expected Fuel=10, got %s", got.Fuel)
	}
	if !got.Total.Equal(decimal.RequireFromString("10")) {
		t.Errorf("expected Total=10, got %s", got.Total)
	}
}

// TestSurchargeCalculator_NilZone verifies that a nil zone pointer is handled
// gracefully, yielding all-zero surcharges.
func TestSurchargeCalculator_NilZone(t *testing.T) {
	in := SurchargeInput{
		BaseFee:    decimal.RequireFromString("20.00"),
		PostalCode: "99501",
	}

	got := CalculateSurcharges(nil, in)

	if !got.Remote.IsZero() {
		t.Errorf("expected Remote=0 (nil zone), got %s", got.Remote)
	}
	if !got.Fuel.IsZero() {
		t.Errorf("expected Fuel=0 (nil zone), got %s", got.Fuel)
	}
	if !got.Total.IsZero() {
		t.Errorf("expected Total=0 (nil zone), got %s", got.Total)
	}
}
