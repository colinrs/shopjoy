package shipping

import (
	"regexp"

	"github.com/shopspring/decimal"
)

// SurchargeInput captures the inputs required to evaluate both surcharges:
// the post-tax (or pre-tax, per zone policy) base fee and the destination
// postal code used for remote-area pattern matching.
type SurchargeInput struct {
	BaseFee    decimal.Decimal
	PostalCode string
}

// SurchargeBreakdown is the result of applying both surcharges to a base fee.
// Total is the sum of Remote and Fuel and is intended to be added on top of
// the base fee by the caller.
type SurchargeBreakdown struct {
	Remote decimal.Decimal
	Fuel   decimal.Decimal
	Total  decimal.Decimal
}

// CalculateSurcharges evaluates the remote-area surcharge and the fuel
// surcharge for a given zone and returns both amounts plus their sum.
//
// Remote surcharge: applied when the zone configures a positive
// RemoteSurcharge AND at least one pattern in RemoteZipPatterns. The first
// pattern whose compiled regexp matches the supplied PostalCode wins; any
// pattern that fails to compile is skipped (continue) rather than aborting
// evaluation.
//
// Fuel surcharge: applied when zone.FuelSurchargePct is positive. Computed
// as BaseFee * FuelSurchargePct and rounded to 2 decimal places to keep
// money math stable at the cents boundary.
//
// A nil zone pointer is tolerated and produces an all-zero breakdown.
//
// Note: callers must only invoke this when the zone is one that actually
// needs surcharge evaluation; the logic layer is responsible for the gating.
func CalculateSurcharges(zone *ShippingZone, in SurchargeInput) SurchargeBreakdown {
	var breakdown SurchargeBreakdown
	if zone == nil {
		return breakdown
	}

	// Remote surcharge: requires a non-zero surcharge amount, at least one
	// pattern, and a non-empty postal code to evaluate.
	if zone.RemoteSurcharge.IsPositive() && len(zone.RemoteZipPatterns) > 0 && in.PostalCode != "" {
		for _, pattern := range zone.RemoteZipPatterns {
			re, err := regexp.Compile(pattern)
			if err != nil {
				continue // skip invalid patterns silently
			}
			if re.MatchString(in.PostalCode) {
				breakdown.Remote = zone.RemoteSurcharge
				break // first match wins
			}
		}
	}

	// Fuel surcharge: percentage of the base fee. A zero base or zero pct
	// yields zero fuel without invoking the multiplication.
	if zone.FuelSurchargePct.IsPositive() && in.BaseFee.IsPositive() {
		fuel := in.BaseFee.Mul(zone.FuelSurchargePct).Round(2)
		breakdown.Fuel = fuel
	}

	breakdown.Total = breakdown.Remote.Add(breakdown.Fuel)
	return breakdown
}