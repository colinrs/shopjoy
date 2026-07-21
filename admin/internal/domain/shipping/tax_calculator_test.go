package shipping

import (
	"testing"

	"github.com/shopspring/decimal"
)

// TestCalculateTax_NotIncluded verifies the tax-exclusive path:
// fee is net of tax, so tax = fee*rate and total = fee+tax.
func TestCalculateTax_NotIncluded(t *testing.T) {
	fee := decimal.RequireFromString("100.00")
	rate := decimal.RequireFromString("0.19")

	tax, total := CalculateTax(fee, rate, false /* taxIncluded */)

	if !tax.Equal(decimal.RequireFromString("19.00")) {
		t.Errorf("expected tax=19.00, got %s", tax)
	}
	if !total.Equal(decimal.RequireFromString("119.00")) {
		t.Errorf("expected total=119.00, got %s", total)
	}
}

// TestCalculateTax_Included verifies the tax-inclusive path:
// fee already contains tax → total == fee, tax = fee - fee/(1+rate).
func TestCalculateTax_Included(t *testing.T) {
	fee := decimal.RequireFromString("119.00")
	rate := decimal.RequireFromString("0.19")

	tax, total := CalculateTax(fee, rate, true /* taxIncluded */)

	// total must equal the inclusive fee unchanged.
	if !total.Equal(fee) {
		t.Errorf("expected total=119.00 (unchanged), got %s", total)
	}
	// tax = 119 - 119/1.19 = 119 - 100 = 19
	if !tax.Round(2).Equal(decimal.RequireFromString("19.00")) {
		t.Errorf("expected tax≈19.00, got %s", tax)
	}
}

// TestCalculateTax_ZeroRate verifies the zero-rate path:
// a zero tax rate must produce tax=0 and total=fee, regardless of taxIncluded.
func TestCalculateTax_ZeroRate(t *testing.T) {
	fee := decimal.RequireFromString("50.00")

	// taxIncluded=false
	tax, total := CalculateTax(fee, decimal.Zero, false)
	if !tax.IsZero() {
		t.Errorf("expected tax=0 (not-included), got %s", tax)
	}
	if !total.Equal(fee) {
		t.Errorf("expected total=fee=50.00 (not-included), got %s", total)
	}

	// taxIncluded=true (must still yield tax=0, total=fee)
	tax, total = CalculateTax(fee, decimal.Zero, true)
	if !tax.IsZero() {
		t.Errorf("expected tax=0 (included), got %s", tax)
	}
	if !total.Equal(fee) {
		t.Errorf("expected total=fee=50.00 (included), got %s", total)
	}
}