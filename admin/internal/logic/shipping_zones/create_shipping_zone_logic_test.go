package shipping_zones

import (
	"reflect"
	"testing"

	"github.com/colinrs/shopjoy/admin/internal/domain/shipping"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/shopspring/decimal"
)

// TestToStringI18n verifies the wire []NameI18nEntry → entity StringI18n
// conversion. Empty entries are dropped; nil/empty input returns nil.
func TestToStringI18n(t *testing.T) {
	cases := []struct {
		name string
		in   []types.NameI18nEntry
		want shipping.StringI18n
	}{
		{
			name: "nil input → nil",
			in:   nil,
			want: nil,
		},
		{
			name: "empty input → nil",
			in:   []types.NameI18nEntry{},
			want: nil,
		},
		{
			name: "single entry",
			in:   []types.NameI18nEntry{{Locale: "en-US", Name: "Standard"}},
			want: shipping.StringI18n{"en-US": "Standard"},
		},
		{
			name: "multiple entries",
			in: []types.NameI18nEntry{
				{Locale: "en-US", Name: "East"},
				{Locale: "zh-CN", Name: "东区"},
			},
			want: shipping.StringI18n{"en-US": "East", "zh-CN": "东区"},
		},
		{
			name: "empty locale dropped",
			in:   []types.NameI18nEntry{{Locale: "", Name: "x"}, {Locale: "en", Name: "y"}},
			want: shipping.StringI18n{"en": "y"},
		},
		{
			name: "empty name dropped",
			in:   []types.NameI18nEntry{{Locale: "en", Name: ""}, {Locale: "de", Name: "D"}},
			want: shipping.StringI18n{"de": "D"},
		},
		{
			name: "all invalid → nil",
			in: []types.NameI18nEntry{
				{Locale: "", Name: "x"},
				{Locale: "en", Name: ""},
			},
			want: nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := toStringI18n(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("toStringI18n(%v) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

// TestFromStringI18n verifies the entity StringI18n → wire []NameI18nEntry
// round-trip preserves entries.
func TestFromStringI18n(t *testing.T) {
	t.Run("nil → nil", func(t *testing.T) {
		got := fromStringI18n(nil)
		if got != nil {
			t.Errorf("fromStringI18n(nil) = %v, want nil", got)
		}
	})
	t.Run("empty → nil", func(t *testing.T) {
		got := fromStringI18n(shipping.StringI18n{})
		if got != nil {
			t.Errorf("fromStringI18n(empty) = %v, want nil", got)
		}
	})
	t.Run("populated map", func(t *testing.T) {
		in := shipping.StringI18n{"en-US": "Standard", "zh-CN": "标准"}
		got := fromStringI18n(in)
		if len(got) != 2 {
			t.Fatalf("expected 2 entries, got %d", len(got))
		}
		// map iteration order is non-deterministic; build a set to compare
		seen := make(map[types.NameI18nEntry]bool, len(got))
		for _, e := range got {
			seen[e] = true
		}
		want := []types.NameI18nEntry{
			{Locale: "en-US", Name: "Standard"},
			{Locale: "zh-CN", Name: "标准"},
		}
		for _, w := range want {
			if !seen[w] {
				t.Errorf("missing entry %v in result %v", w, got)
			}
		}
	})
}

// TestFromStringI18n_ToStringI18n_Roundtrip verifies bidirectional conversion
// preserves locale/name pairs (order independent).
func TestFromStringI18n_ToStringI18n_Roundtrip(t *testing.T) {
	original := []types.NameI18nEntry{
		{Locale: "en-US", Name: "Standard"},
		{Locale: "zh-CN", Name: "标准"},
		{Locale: "ja-JP", Name: "通常"},
	}
	entity := toStringI18n(original)
	wire := fromStringI18n(entity)
	if len(wire) != len(original) {
		t.Fatalf("roundtrip lost entries: original=%d wire=%d", len(original), len(wire))
	}
	seen := make(map[types.NameI18nEntry]bool, len(wire))
	for _, e := range wire {
		seen[e] = true
	}
	for _, e := range original {
		if !seen[e] {
			t.Errorf("roundtrip dropped %v", e)
		}
	}
}

// TestParseAmount verifies empty/invalid strings → zero decimal.
func TestParseAmount(t *testing.T) {
	cases := []struct {
		name, in string
		wantStr  string // decimal string for comparison
	}{
		{"empty → 0", "", "0"},
		{"invalid → 0", "not-a-number", "0"},
		{"valid integer", "100", "100"},
		{"valid decimal", "5.50", "5.5"},
		{"zero", "0", "0"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseAmount(tc.in)
			if got.String() != tc.wantStr {
				t.Errorf("parseAmount(%q) = %s, want %s", tc.in, got.String(), tc.wantStr)
			}
		})
	}
}

// TestDefaultCurrency verifies empty input defaults to "CNY".
func TestDefaultCurrency(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"", "CNY"},
		{"USD", "USD"},
		{"EUR", "EUR"},
	}
	for _, c := range cases {
		if got := defaultCurrency(c.in); got != c.want {
			t.Errorf("defaultCurrency(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// TestDefaultInt verifies zero input defaults to the supplied fallback.
func TestDefaultInt(t *testing.T) {
	cases := []struct {
		in, def, want int
	}{
		{0, 5000, 5000},
		{6000, 5000, 6000},
		{1, 5000, 1},
	}
	for _, c := range cases {
		if got := defaultInt(c.in, c.def); got != c.want {
			t.Errorf("defaultInt(%d, %d) = %d, want %d", c.in, c.def, got, c.want)
		}
	}
}

// TestShippingZone_ValidateByVolume verifies the entity Validate() rejects
// by_volume zones without VolumetricDivisor > 0.
func TestShippingZone_ValidateByVolume(t *testing.T) {
	// Missing VolumetricDivisor → must fail.
	z := &shipping.ShippingZone{
		Name:            "EU Zone",
		Regions:         shipping.Regions{"DE"},
		FeeType:         shipping.FeeTypeByVolume,
		FirstUnit:       1,
		FirstFee:        decimal.RequireFromString("5.00"),
		AdditionalUnit:  1,
		AdditionalFee:   decimal.RequireFromString("2.00"),
	}
	if err := z.Validate(); err == nil {
		t.Error("expected Validate() to fail when by_volume lacks VolumetricDivisor")
	}

	// With VolumetricDivisor > 0 → passes.
	z.VolumetricDivisor = 5000
	if err := z.Validate(); err != nil {
		t.Errorf("expected Validate() to pass, got %v", err)
	}
}

// TestFeeType_IsValidV2 verifies the new validator recognises by_volume.
func TestFeeType_IsValidV2(t *testing.T) {
	valid := []shipping.FeeType{
		shipping.FeeTypeFixed,
		shipping.FeeTypeByCount,
		shipping.FeeTypeByWeight,
		shipping.FeeTypeByVolume,
		shipping.FeeTypeFree,
	}
	for _, f := range valid {
		if !f.IsValidV2() {
			t.Errorf("IsValidV2(%q) = false, want true", f)
		}
	}

	// Unknown values are invalid.
	invalid := []shipping.FeeType{"", "unknown", "magic"}
	for _, f := range invalid {
		if f.IsValidV2() {
			t.Errorf("IsValidV2(%q) = true, want false", f)
		}
	}
}

// TestFeeType_IsValid_Legacy verifies IsValid (without by_volume) is preserved
// for any callers still using the legacy validator.
func TestFeeType_IsValid_Legacy(t *testing.T) {
	// by_volume is intentionally NOT accepted by the legacy IsValid.
	if shipping.FeeTypeByVolume.IsValid() {
		t.Error("IsValid(by_volume) should be false in legacy validator")
	}
	// Original 4 values still valid.
	for _, f := range []shipping.FeeType{
		shipping.FeeTypeFixed,
		shipping.FeeTypeByCount,
		shipping.FeeTypeByWeight,
		shipping.FeeTypeFree,
	} {
		if !f.IsValid() {
			t.Errorf("IsValid(%q) = false, want true", f)
		}
	}
}

// TestDefaultVolumetricDivisor verifies the entity exports the expected
// industry-standard default constant (5000 cm³/kg).
func TestDefaultVolumetricDivisor(t *testing.T) {
	if shipping.DefaultVolumetricDivisor != 5000 {
		t.Errorf("DefaultVolumetricDivisor = %d, want 5000", shipping.DefaultVolumetricDivisor)
	}
}
