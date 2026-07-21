package shipping

import "github.com/shopspring/decimal"

// IOSSThresholdEUR is the EU IOSS low-value consignment cap (€150).
// Per EU regulation 2021/847, IOSS may only be used when the intrinsic value
// of the consignment is strictly below this threshold (i.e. value < 150.00 EUR).
const IOSSThresholdEUR = 150

// IOSSReason values describe why an order did NOT qualify for IOSS.
// An empty string ("") means IOSS does apply.
const (
	// IOSSReasonNotApplicable: the zone has IossApplicable=false (merchant opt-out).
	IOSSReasonNotApplicable = "not_applicable"
	// IOSSReasonCurrencyNotEUR: the zone is denominated in a non-EUR currency.
	IOSSReasonCurrencyNotEUR = "currency_not_eur"
	// IOSSReasonExceedsThreshold: the order value is >= €150.
	IOSSReasonExceedsThreshold = "exceeds_threshold"
)

// EvaluateIOSS decides whether a given order qualifies for EU IOSS (Import
// One-Stop Shop) declaration. IOSS is a B2C cross-border VAT simplification
// mechanism available only for consignments whose intrinsic value is strictly
// under €150 (IOSSThresholdEUR), denominated in EUR, AND the destination zone
// has explicitly opted in (IossApplicable=true).
//
// The returned reason is empty when IOSS applies; otherwise it identifies the
// first failing condition (checked in the order: opt-out → currency → value).
// This lets the caller both gate the declaration AND surface a human-readable
// reason to operators / end users.
func EvaluateIOSS(orderValue decimal.Decimal, iossApplicable bool, currency string) (requiresIOSS bool, reason string) {
	if !iossApplicable {
		return false, IOSSReasonNotApplicable
	}
	if currency != "EUR" {
		return false, IOSSReasonCurrencyNotEUR
	}
	if orderValue.GreaterThanOrEqual(decimal.NewFromInt(IOSSThresholdEUR)) {
		return false, IOSSReasonExceedsThreshold
	}
	return true, ""
}
