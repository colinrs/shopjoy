package shipping

import (
	"testing"

	"github.com/colinrs/shopjoy/admin/internal/domain/shipping"
)

// TestResolveZoneName_ExactMatch 精确匹配：acceptLanguage="en-US"，NameI18n["en-US"]
// 非空时直接返回该值。
func TestResolveZoneName_ExactMatch(t *testing.T) {
	zone := &shipping.ShippingZone{
		Name:     "华东",
		NameI18n: shipping.StringI18n{"en-US": "East China", "ja-JP": "華東"},
	}
	got := ResolveZoneName(zone, "en-US")
	if got != "East China" {
		t.Errorf("expected East China, got %s", got)
	}
}

// TestResolveZoneName_LanguageBaseMatch 语言部分匹配：acceptLanguage="en-GB" 找不到
// 精确 locale，但基础语言 "en" 命中 NameI18n["en-US"]。
func TestResolveZoneName_LanguageBaseMatch(t *testing.T) {
	zone := &shipping.ShippingZone{
		Name:     "华东",
		NameI18n: shipping.StringI18n{"en-US": "East China", "ja-JP": "華東"},
	}
	got := ResolveZoneName(zone, "en-GB")
	if got != "East China" {
		t.Errorf("expected East China (base match), got %s", got)
	}
}

// TestResolveZoneName_FirstNonEmpty 没有匹配项时返回 NameI18n 中第一个非空值。
//
// Note: Go map iteration 是随机顺序。为让断言稳定，本测试只放一条非空条目，
// 仍覆盖"无精确 + 无 base → 走首个非空"分支。
func TestResolveZoneName_FirstNonEmpty(t *testing.T) {
	zone := &shipping.ShippingZone{
		Name:     "华东",
		NameI18n: shipping.StringI18n{"en-US": "East China"},
	}
	got := ResolveZoneName(zone, "fr-FR")
	if got != "East China" {
		t.Errorf("expected first non-empty (East China), got %s", got)
	}
}

// TestResolveZoneName_EmptyAcceptLanguage acceptLanguage 为空直接走 fallback。
func TestResolveZoneName_EmptyAcceptLanguage(t *testing.T) {
	zone := &shipping.ShippingZone{
		Name:     "华东",
		NameI18n: shipping.StringI18n{"en-US": "East China"},
	}
	got := ResolveZoneName(zone, "")
	if got != "华东" {
		t.Errorf("expected fallback zone.Name (华东), got %s", got)
	}
}

// TestResolveZoneName_AllEmpty NameI18n 中所有值都为空，走 fallback。
func TestResolveZoneName_AllEmpty(t *testing.T) {
	zone := &shipping.ShippingZone{
		Name:     "华东",
		NameI18n: shipping.StringI18n{"en-US": "", "ja-JP": ""},
	}
	got := ResolveZoneName(zone, "en-US")
	if got != "华东" {
		t.Errorf("expected fallback zone.Name, got %s", got)
	}
}

// TestResolveZoneName_NilMap NameI18n 为 nil map，走 fallback。
func TestResolveZoneName_NilMap(t *testing.T) {
	zone := &shipping.ShippingZone{
		Name: "华东",
	}
	got := ResolveZoneName(zone, "en-US")
	if got != "华东" {
		t.Errorf("expected fallback zone.Name, got %s", got)
	}
}

// TestResolveZoneName_NilZone nil zone 直接返回空（防 panic）。
func TestResolveZoneName_NilZone(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("ResolveZoneName(nil, ...) must not panic, got %v", r)
		}
	}()
	got := ResolveZoneName(nil, "en-US")
	if got != "" {
		t.Errorf("expected empty string for nil zone, got %q", got)
	}
}

// TestResolveZoneName_ExactPrecedesBase 精确匹配优先于语言基础匹配。
func TestResolveZoneName_ExactPrecedesBase(t *testing.T) {
	zone := &shipping.ShippingZone{
		Name: "华东",
		NameI18n: shipping.StringI18n{
			"en":    "English Generic",
			"en-US": "American Specific",
		},
	}
	got := ResolveZoneName(zone, "en-US")
	if got != "American Specific" {
		t.Errorf("expected American Specific (exact), got %s", got)
	}
}

// TestResolveZoneName_EmptyStringInMapSkipped 在语言基础匹配中跳过空值条目。
func TestResolveZoneName_EmptyStringInMapSkipped(t *testing.T) {
	zone := &shipping.ShippingZone{
		Name: "华东",
		NameI18n: shipping.StringI18n{
			"en-US": "",
			"en-GB": "Britain",
		},
	}
	got := ResolveZoneName(zone, "en-AU")
	if got != "Britain" {
		t.Errorf("expected Britain (skip empty), got %s", got)
	}
}
