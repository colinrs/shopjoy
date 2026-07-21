package shipping

import (
	"testing"

	"github.com/shopspring/decimal"
)

// TestEvaluateIOSS_Applicable covers the happy path:
// EU-zone, EUR currency, order value strictly under the €150 threshold → IOSS applies.
func TestEvaluateIOSS_Applicable(t *testing.T) {
	orderValue := decimal.RequireFromString("149.99")
	requiresIOSS, reason := EvaluateIOSS(orderValue, true /* iossApplicable */, "EUR")
	if !requiresIOSS {
		t.Errorf("expected IOSS applicable, got reason=%q", reason)
	}
	if reason != "" {
		t.Errorf("expected empty reason when IOSS applies, got %q", reason)
	}
}

// TestEvaluateIOSS_AtThreshold covers the boundary case:
// exactly €150 must NOT use IOSS (the EU regulation caps IOSS at < €150).
func TestEvaluateIOSS_AtThreshold(t *testing.T) {
	orderValue := decimal.RequireFromString("150.00")
	requiresIOSS, reason := EvaluateIOSS(orderValue, true, "EUR")
	if requiresIOSS {
		t.Errorf("expected IOSS not applicable at threshold, got reason=%q", reason)
	}
	if reason != IOSSReasonExceedsThreshold {
		t.Errorf("expected reason=%q at threshold, got %q", IOSSReasonExceedsThreshold, reason)
	}
}

// TestEvaluateIOSS_ExceedsThreshold covers the over-threshold case:
// order > €150 must not use IOSS; the parcel needs standard customs clearance.
func TestEvaluateIOSS_ExceedsThreshold(t *testing.T) {
	orderValue := decimal.RequireFromString("250.00")
	requiresIOSS, reason := EvaluateIOSS(orderValue, true, "EUR")
	if requiresIOSS {
		t.Errorf("expected IOSS not applicable when over threshold, got reason=%q", reason)
	}
	if reason != IOSSReasonExceedsThreshold {
		t.Errorf("expected reason=%q, got %q", IOSSReasonExceedsThreshold, reason)
	}
}

// TestEvaluateIOSS_NotApplicable covers the opt-out case:
// the merchant turned off IOSS for this zone, regardless of currency or value.
func TestEvaluateIOSS_NotApplicable(t *testing.T) {
	orderValue := decimal.RequireFromString("50.00")
	requiresIOSS, reason := EvaluateIOSS(orderValue, false /* iossApplicable */, "EUR")
	if requiresIOSS {
		t.Errorf("expected IOSS not applicable when zone flag is off, got reason=%q", reason)
	}
	if reason != IOSSReasonNotApplicable {
		t.Errorf("expected reason=%q, got %q", IOSSReasonNotApplicable, reason)
	}
}

// TestEvaluateIOSS_NonEURCurrency covers the currency guard:
// only EUR-denominated zones can use IOSS. Other currencies must skip it
// even when the zone flag is on and the order value is well under the cap.
func TestEvaluateIOSS_NonEURCurrency(t *testing.T) {
	orderValue := decimal.RequireFromString("50.00")
	requiresIOSS, reason := EvaluateIOSS(orderValue, true, "USD")
	if requiresIOSS {
		t.Errorf("expected IOSS not applicable for non-EUR currency, got reason=%q", reason)
	}
	if reason != IOSSReasonCurrencyNotEUR {
		t.Errorf("expected reason=%q, got %q", IOSSReasonCurrencyNotEUR, reason)
	}
}
