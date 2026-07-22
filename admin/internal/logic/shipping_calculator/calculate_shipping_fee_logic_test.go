package shipping_calculator

import (
	"context"
	"testing"

	shippingapp "github.com/colinrs/shopjoy/admin/internal/application/shipping"
	"github.com/colinrs/shopjoy/admin/internal/domain/shipping"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/shopspring/decimal"
)

// TestCalculateShippingFee_FillsTaxAndTotal verifies the tax-exclusive path:
// zone taxable, tax_included=false, rate=0.19 → tax = fee*0.19, total = fee+tax.
func TestCalculateShippingFee_FillsTaxAndTotal(t *testing.T) {
	fee := decimal.RequireFromString("100.00")
	rate := decimal.RequireFromString("0.19")

	tax, total := shipping.CalculateTax(fee, rate, false /* taxIncluded */)

	if !tax.Equal(decimal.RequireFromString("19.00")) {
		t.Errorf("expected tax=19.00, got %s", tax)
	}
	if !total.Equal(decimal.RequireFromString("119.00")) {
		t.Errorf("expected total=119.00, got %s", total)
	}
}

// TestCalculateTax_TaxIncluded verifies the tax-inclusive path:
// fee already contains tax → total == fee, tax = fee - fee/(1+rate).
func TestCalculateTax_TaxIncluded(t *testing.T) {
	fee := decimal.RequireFromString("119.00")
	rate := decimal.RequireFromString("0.19")

	tax, total := shipping.CalculateTax(fee, rate, true /* taxIncluded */)

	// total must equal the inclusive fee unchanged.
	if !total.Equal(fee) {
		t.Errorf("expected total=119.00 (unchanged), got %s", total)
	}
	// tax = 119 - 119/1.19 = 119 - 100 = 19
	if !tax.Round(2).Equal(decimal.RequireFromString("19.00")) {
		t.Errorf("expected tax≈19.00, got %s", tax)
	}
}

// TestCalculateTax_NotTaxable is enforced by the caller passing rate=0.
// A zero rate must produce zero tax and total == fee.
func TestCalculateTax_ZeroRate(t *testing.T) {
	fee := decimal.RequireFromString("50.00")

	tax, total := shipping.CalculateTax(fee, decimal.Zero, false)

	if !tax.IsZero() {
		t.Errorf("expected tax=0, got %s", tax)
	}
	if !total.Equal(fee) {
		t.Errorf("expected total=fee=50.00, got %s", total)
	}
}

// TestCalculateShippingFee_CarrierCode verifies template carrier code passthrough
// and the "standard" fallback when empty.
func TestCalculateShippingFee_CarrierCode(t *testing.T) {
	if got := resolveCarrierCode("dhl"); got != "dhl" {
		t.Errorf("expected carrier=dhl, got %s", got)
	}
	if got := resolveCarrierCode(""); got != "standard" {
		t.Errorf("expected fallback carrier=standard, got %s", got)
	}
}

// TestResolveEstimatedDays_StandardRegistered checks the happy path: the
// default registry (registered with StandardCarrier) returns the documented
// country-tier days from StandardCarrier.EstimateDays.
func TestResolveEstimatedDays_StandardRegistered(t *testing.T) {
	registry := shipping.NewCarrierRegistry()

	cases := []struct {
		country string
		want    int
	}{
		{"CN", 3},
		{"US", 7},
		{"DE", 7},
		{"JP", 5},
		{"SG", 5},
		{"ZZ", 10}, // unknown country → default tier
		{"cn", 3},  // case-insensitive
	}
	for _, tc := range cases {
		t.Run(tc.country, func(t *testing.T) {
			got := resolveEstimatedDays(registry, "standard", tc.country)
			if got != tc.want {
				t.Errorf("estimateDays(standard, %q)=%d, want %d", tc.country, got, tc.want)
			}
		})
	}
}

// TestResolveEstimatedDays_UnknownCarrierFallsBack ensures that an unregistered
// carrier code (e.g. a future "dhl" carrier before it is wired into the
// registry) falls back to the 5-day default rather than silently returning 0.
func TestResolveEstimatedDays_UnknownCarrierFallsBack(t *testing.T) {
	registry := shipping.NewCarrierRegistry()
	got := resolveEstimatedDays(registry, "dhl", "US")
	if got != 5 {
		t.Errorf("expected fallback=5 for unknown carrier, got %d", got)
	}
}

// TestResolveEstimatedDays_NilRegistryIsSafe verifies defensive handling: a nil
// registry (e.g. an older ServiceContext from before Phase 7-2 wiring) yields
// the fallback instead of panicking.
func TestResolveEstimatedDays_NilRegistryIsSafe(t *testing.T) {
	got := resolveEstimatedDays(nil, "standard", "US")
	if got != 5 {
		t.Errorf("expected fallback=5 for nil registry, got %d", got)
	}
}

// TestEvaluateIOSS_AllFourReasons exercises the four canonical IOSS outcomes
// returned by shipping.EvaluateIOSS so the wire string values are pinned.
// These reason strings are part of the API contract — the frontend renders
// them directly to explain why IOSS did or did not apply.
func TestEvaluateIOSS_AllFourReasons(t *testing.T) {
	cases := []struct {
		name           string
		orderValue     string
		iossApplicable bool
		currency       string
		wantReason     string
	}{
		{
			name:           "EUR + under threshold + ioss on → empty reason (applies)",
			orderValue:     "100.00",
			iossApplicable: true,
			currency:       "EUR",
			wantReason:     "",
		},
		{
			name:           "EUR + >= threshold → exceeds_threshold",
			orderValue:     "150.00",
			iossApplicable: true,
			currency:       "EUR",
			wantReason:     shipping.IOSSReasonExceedsThreshold,
		},
		{
			name:           "non-EUR currency → currency_not_eur",
			orderValue:     "50.00",
			iossApplicable: true,
			currency:       "USD",
			wantReason:     shipping.IOSSReasonCurrencyNotEUR,
		},
		{
			name:           "ioss_applicable=false → not_applicable (wins over currency/value)",
			orderValue:     "50.00",
			iossApplicable: false,
			currency:       "EUR",
			wantReason:     shipping.IOSSReasonNotApplicable,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			orderValue := decimal.RequireFromString(tc.orderValue)
			_, reason := shipping.EvaluateIOSS(orderValue, tc.iossApplicable, tc.currency)
			if reason != tc.wantReason {
				t.Errorf("EvaluateIOSS(%s, applicable=%v, %q) reason=%q, want %q",
					tc.orderValue, tc.iossApplicable, tc.currency, reason, tc.wantReason)
			}
		})
	}
}

// TestCalculateShippingFee_SurchargeAddedToFee verifies the logic-layer contract
// that surcharges are computed off the base fee and added on top: the remote
// surcharge (flat) plus the fuel surcharge (pct * base) sum into the final fee.
func TestCalculateShippingFee_SurchargeAddedToFee(t *testing.T) {
	baseFee := decimal.RequireFromString("100.00")
	zone := &shipping.ShippingZone{
		RemoteSurcharge:   decimal.RequireFromString("15.00"),
		RemoteZipPatterns: []string{"^99"},
		FuelSurchargePct:  decimal.RequireFromString("0.10"),
	}

	sc := shipping.CalculateSurcharges(zone, shipping.SurchargeInput{
		BaseFee:    baseFee,
		PostalCode: "99001",
	})

	// remote=15.00 (flat), fuel=100*0.10=10.00, total=25.00
	if !sc.Remote.Equal(decimal.RequireFromString("15.00")) {
		t.Errorf("expected remote=15.00, got %s", sc.Remote)
	}
	if !sc.Fuel.Equal(decimal.RequireFromString("10.00")) {
		t.Errorf("expected fuel=10.00, got %s", sc.Fuel)
	}
	if !sc.Total.Equal(decimal.RequireFromString("25.00")) {
		t.Errorf("expected surcharge total=25.00, got %s", sc.Total)
	}

	// Final fee = base + surcharge total, mirroring the logic layer.
	finalFee := baseFee.Add(sc.Total)
	if !finalFee.Equal(decimal.RequireFromString("125.00")) {
		t.Errorf("expected final fee=125.00, got %s", finalFee)
	}
	if got := formatAmount(sc.Total); got != "25.00" {
		t.Errorf("expected applied_surcharge=25.00, got %s", got)
	}
}

// chargeable weight (max of real and volumetric) and VolumetricWeight reports
// the accumulated volumetric weight for debug display.
func TestCalculateShippingFee_WeightBreakdown(t *testing.T) {
	divisor := 5000
	items := []shipping.CalculateItem{
		// real 500g; volumetric: 200*200*200mm = 8000cm³ → 1600g. chargeable=1600.
		{Quantity: 2, Weight: 500, Length: 200, Width: 200, Height: 200},
		// real 300g; no dims → volumetric 0. chargeable=300.
		{Quantity: 1, Weight: 300},
	}

	chargeable, volumetric := sumWeights(items, divisor)

	// chargeable = 1600*2 + 300*1 = 3500
	if chargeable != 3500 {
		t.Errorf("expected chargeable weight=3500, got %d", chargeable)
	}
	// volumetric = 1600*2 + 0*1 = 3200
	if volumetric != 3200 {
		t.Errorf("expected volumetric weight=3200, got %d", volumetric)
	}
}

// TestGetAcceptLanguage_AbsentReturnsEmpty verifies that the helper returns ""
// when the handler has not injected Accept-Language into ctx. The resolver
// treats "" as "no locale signal" and falls back to zone.Name, so production
// behaviour is identical to the pre-i18n wiring.
func TestGetAcceptLanguage_AbsentReturnsEmpty(t *testing.T) {
	if got := contextx.GetAcceptLanguage(context.Background()); got != "" {
		t.Errorf("expected empty string when ctx has no accept-language value, got %q", got)
	}
}

// TestGetAcceptLanguage_Injected verifies the helper reads the value the
// handler would inject via contextx.SetAcceptLanguage(ctx, ...).
func TestGetAcceptLanguage_Injected(t *testing.T) {
	ctx := contextx.SetAcceptLanguage(context.Background(), "en-US,en;q=0.9")
	if got := contextx.GetAcceptLanguage(ctx); got != "en-US,en;q=0.9" {
		t.Errorf("expected injected accept-language to round-trip, got %q", got)
	}
}

// TestZoneNameResolution_RespectsAcceptLanguage models the response-building
// step's ZoneName computation without requiring a DB. It proves that the
// helper + shippingapp.ResolveZoneName composition behaves correctly end-to-end
// when ctx carries an Accept-Language header value.
func TestZoneNameResolution_RespectsAcceptLanguage(t *testing.T) {
	zone := &shipping.ShippingZone{
		Name:     "华东",
		NameI18n: shipping.StringI18n{"en-US": "East China"},
	}

	cases := []struct {
		name string
		ctx  context.Context
		want string
	}{
		{
			name: "no accept-language in ctx → fallback to zone.Name",
			ctx:  context.Background(),
			want: "华东",
		},
		{
			name: "exact match en-US → returns en-US name",
			ctx:  contextx.SetAcceptLanguage(context.Background(), "en-US"),
			want: "East China",
		},
		{
			name: "language-base match en-GB → returns en-US name (base=en)",
			ctx:  contextx.SetAcceptLanguage(context.Background(), "en-GB"),
			want: "East China",
		},
		{
			name: "no match for fr-FR → first non-empty",
			ctx:  contextx.SetAcceptLanguage(context.Background(), "fr-FR"),
			want: "East China",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := shippingapp.ResolveZoneName(zone, contextx.GetAcceptLanguage(tc.ctx))
			if got != tc.want {
				t.Errorf("ResolveZoneName(%q)=%q, want %q", contextx.GetAcceptLanguage(tc.ctx), got, tc.want)
			}
		})
	}
}
