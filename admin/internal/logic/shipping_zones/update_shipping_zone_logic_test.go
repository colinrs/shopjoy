package shipping_zones

import "testing"

// TestApplyBoolPtrOverride pins the Important/bool-pointer contract used
// by UpdateShippingZone.Taxable/TaxIncluded/IossApplicable:
//
//   - nil pointer → caller did not supply the field; the current value is preserved.
//   - non-nil pointer (true) → caller explicitly set true; current value is overwritten.
//   - non-nil pointer (false) → caller explicitly set false; current value is overwritten (this is the fix: the old non-pointer bool silently dropped false-on-true).
//
// The fix exists because prior to this Task, the wire type was `bool` with
// `if req.Taxable { zone.Taxable = true }`. A caller wanting to disable a
// previously-enabled flag would set `taxable: false` in the request — but the
// `if` branch only fired on true, leaving the entity's true value in place.
// Switching to *bool + applyBoolPtrOverride makes the explicit-false case work.
func TestApplyBoolPtrOverride(t *testing.T) {
	cases := []struct {
		name      string
		current   bool
		override  *bool
		wantValue bool
	}{
		{
			name:      "nil override, current true → keep true",
			current:   true,
			override:  nil,
			wantValue: true,
		},
		{
			name:      "nil override, current false → keep false",
			current:   false,
			override:  nil,
			wantValue: false,
		},
		{
			name:      "true override, current false → set true",
			current:   false,
			override:  boolPtr(true),
			wantValue: true,
		},
		{
			name:      "true override, current true → set true",
			current:   true,
			override:  boolPtr(true),
			wantValue: true,
		},
		{
			// The headline case: previously dropped, now works.
			name:      "false override, current true → set false (explicit disable works!)",
			current:   true,
			override:  boolPtr(false),
			wantValue: false,
		},
		{
			name:      "false override, current false → set false",
			current:   false,
			override:  boolPtr(false),
			wantValue: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := applyBoolPtrOverride(tc.current, tc.override)
			if got != tc.wantValue {
				t.Errorf("applyBoolPtrOverride(current=%v, override=%v) = %v, want %v",
					tc.current, deref(tc.override), got, tc.wantValue)
			}
		})
	}
}

// TestApplyBoolPtrOverride_ExplicitFalseOnThreeToggleFlags simulates the
// end-to-end story: a zone with all three toggles currently true receives
// an UpdateShippingZoneReq whose body sets all three to false. The naive
// `if req.X { zone.X = true }` would leave all three true; the new
// pointer-based handling must flip all three to false.
func TestApplyBoolPtrOverride_ExplicitFalseOnThreeToggleFlags(t *testing.T) {
	type toggle struct {
		name     string
		current  bool
		override *bool
	}

	falsePtr := boolPtr(false)
	falses := []*bool{falsePtr, falsePtr, falsePtr}

	toggles := []toggle{
		{"Taxable", true, falses[0]},
		{"TaxIncluded", true, falses[1]},
		{"IossApplicable", true, falses[2]},
	}

	after := make([]bool, len(toggles))
	for i, tg := range toggles {
		after[i] = applyBoolPtrOverride(tg.current, tg.override)
	}

	for i, tg := range toggles {
		if after[i] != false {
			t.Errorf("%s: applyBoolPtrOverride(current=true, override=&false) = %v, want false",
				tg.name, after[i])
		}
	}
}

func boolPtr(b bool) *bool { return &b }

// deref is a small helper for nicer error messages.
func deref(b *bool) string {
	if b == nil {
		return "<nil>"
	}
	if *b {
		return "&true"
	}
	return "&false"
}
