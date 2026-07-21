package shipping

import "github.com/shopspring/decimal"

// CalculateTax derives (tax, total) from a monetary base.
//
//   - taxIncluded=false: the fee is net of tax, so tax = fee*rate and
//     total = fee + tax.
//   - taxIncluded=true: the fee already contains tax, so total = fee and
//     tax = fee - fee/(1+rate) (the embedded tax portion).
//
// A zero (or non-positive) rate always yields tax=0 and total=fee, regardless
// of taxIncluded.
//
// Callers must only invoke this when the zone is taxable; the logic layer
// guards this with the ShippingZone.Taxable flag.
func CalculateTax(fee, rate decimal.Decimal, taxIncluded bool) (tax, total decimal.Decimal) {
	if !rate.IsPositive() {
		return decimal.Zero, fee
	}
	if taxIncluded {
		net := fee.Div(decimal.NewFromInt(1).Add(rate))
		return fee.Sub(net), fee
	}
	tax = fee.Mul(rate)
	return tax, fee.Add(tax)
}