package shipping_calculator

import (
	"testing"

	"github.com/colinrs/shopjoy/admin/internal/domain/shipping"
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

// TestCalculateShippingFee_WeightBreakdown verifies CalculatedWeight uses the
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
