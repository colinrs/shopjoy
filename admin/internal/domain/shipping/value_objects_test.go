package shipping

import (
	"encoding/json"
	"testing"
)

func TestStringI18n_Marshal(t *testing.T) {
	i18n := StringI18n{"zh-CN": "华东", "en-US": "East China", "ja-JP": "華東"}
	data, err := json.Marshal(i18n)
	if err != nil {
		t.Fatal(err)
	}
	var got StringI18n
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got["en-US"] != "East China" {
		t.Errorf("expected East China, got %s", got["en-US"])
	}
}

func TestStringI18n_ValueScan(t *testing.T) {
	src := StringI18n{"en-US": "Standard", "zh-CN": "标准"}
	v, err := src.Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	bytes, ok := v.([]byte)
	if !ok {
		t.Fatalf("Value returned %T, want []byte", v)
	}
	var got StringI18n
	if err := got.Scan(bytes); err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if got["en-US"] != "Standard" || got["zh-CN"] != "标准" {
		t.Errorf("roundtrip mismatch: got %v", got)
	}

	// Scan from string (alternative driver return)
	var got2 StringI18n
	if err := got2.Scan(string(bytes)); err != nil {
		t.Fatalf("Scan(string): %v", err)
	}
	if got2["zh-CN"] != "标准" {
		t.Errorf("Scan(string) mismatch: got %v", got2)
	}

	// nil roundtrip
	var nilI18n StringI18n
	if err := nilI18n.Scan(nil); err != nil {
		t.Fatalf("Scan(nil): %v", err)
	}
	if nilI18n != nil {
		t.Errorf("Scan(nil) should leave nil, got %v", nilI18n)
	}

	// unsupported type
	var bad StringI18n
	if err := bad.Scan(123); err == nil {
		t.Error("Scan(int) should error")
	}
}

func TestStringI18n_Get(t *testing.T) {
	i18n := StringI18n{
		"en-US":   "East Region",
		"zh-CN":   "华东区域",
		"ja-JP":   "華東",
		"empty-v": "",
	}

	cases := []struct {
		locale, fallback, want string
	}{
		{"zh-CN", "en-US", "华东区域"},      // exact
		{"en-US", "ja-JP", "East Region"},  // exact
		{"ja-JP", "en-US", "華東"},          // exact
		{"fr-FR", "en-US", "East Region"},  // fallback exact hit
		{"de-DE", "en-US", "East Region"},  // fallback exact hit
		{"fr-FR", "ja-JP", "華東"},          // both miss → first non-empty
		{"", "en-US", "East Region"},       // empty locale → fallback hit (skip empty-key collision)
		{"en-US", "", "East Region"},       // exact hit regardless of fallback
		{"empty-v", "en-US", "East Region"}, // exact hit but empty value → fallback
	}

	for _, c := range cases {
		if got := i18n.Get(c.locale, c.fallback); got != c.want {
			t.Errorf("Get(%q,%q)=%q, want %q", c.locale, c.fallback, got, c.want)
		}
	}

	// nil map → ""
	var nilI18n StringI18n
	if got := nilI18n.Get("en-US", "zh-CN"); got != "" {
		t.Errorf("nil.Get should return empty, got %q", got)
	}

	// all empty values → ""
	allEmpty := StringI18n{"a": "", "b": ""}
	if got := allEmpty.Get("a", "b"); got != "" {
		t.Errorf("all-empty Get should return empty, got %q", got)
	}
}

func TestStringArray_Contains(t *testing.T) {
	arr := StringArray{"110000", "310000", "US-CA"}
	if !arr.Contains("110000") {
		t.Error("expected contains 110000")
	}
	if arr.Contains("999999") {
		t.Error("expected not contains 999999")
	}
}

func TestStringArray_ValueScan(t *testing.T) {
	src := StringArray{"99*", "88*", "70*"}
	v, err := src.Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	bytes, ok := v.([]byte)
	if !ok {
		t.Fatalf("Value returned %T, want []byte", v)
	}
	var got StringArray
	if err := got.Scan(bytes); err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if !got.Contains("99*") || !got.Contains("88*") || !got.Contains("70*") {
		t.Errorf("roundtrip mismatch: got %v", got)
	}

	// nil
	var nilArr StringArray
	if err := nilArr.Scan(nil); err != nil {
		t.Fatalf("Scan(nil): %v", err)
	}
	if nilArr != nil {
		t.Errorf("Scan(nil) should leave nil, got %v", nilArr)
	}

	// unsupported type
	var bad StringArray
	if err := bad.Scan(42); err == nil {
		t.Error("Scan(int) should error")
	}
}