package shipping

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
)

// TestStandardCarrier_Quote 验证基础费 + 体积重 + 税。
// zone: by_volume, first_unit=1000g first_fee=10, additional 1000g/5,
// taxable rate=0.10 (税外加)。单件 20*20*50cm=20000cm³/5000=4000g 体积重，
// 实重 1000g → 计费重 4000g × 1 = 4000g。
// 单位数 4000 > first 1000 → additional=(4000-1000+999)/1000=3 → base=10+5*3=25。
// 无附加费；税=25*0.10=2.5；total=27.5。计费重=4000。
func TestStandardCarrier_Quote(t *testing.T) {
	zone := &ShippingZone{
		Currency:          "USD",
		FeeType:           FeeTypeByVolume,
		FirstUnit:         1000,
		FirstFee:          decimal.RequireFromString("10"),
		AdditionalUnit:    1000,
		AdditionalFee:     decimal.RequireFromString("5"),
		VolumetricDivisor: 5000,
		Taxable:           true,
		TaxRate:           decimal.RequireFromString("0.10"),
		TaxIncluded:       false,
	}
	req := QuoteRequest{
		Zone: zone,
		Items: []CalculateItem{
			{Quantity: 1, Weight: 1000, Length: 200, Width: 200, Height: 500, Price: decimal.NewFromInt(100)},
		},
		OrderAmount: decimal.NewFromInt(100),
		ItemCount:   1,
		Address:     MatchInput{CountryCode: "US"},
	}

	got, err := StandardCarrier{}.Quote(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.BaseFee.Equal(decimal.RequireFromString("25")) {
		t.Errorf("expected BaseFee=25, got %s", got.BaseFee)
	}
	if !got.Tax.Equal(decimal.RequireFromString("2.5")) {
		t.Errorf("expected Tax=2.5, got %s", got.Tax)
	}
	if !got.Total.Equal(decimal.RequireFromString("27.5")) {
		t.Errorf("expected Total=27.5, got %s", got.Total)
	}
	if got.Weight != 4000 {
		t.Errorf("expected Weight=4000, got %d", got.Weight)
	}
	if got.Currency != "USD" {
		t.Errorf("expected Currency=USD, got %s", got.Currency)
	}
	if got.EstimatedDays != 7 {
		t.Errorf("expected EstimatedDays=7 (US), got %d", got.EstimatedDays)
	}
}

// TestStandardCarrier_Quote_NilZone 验证 zone 为空返回错误。
func TestStandardCarrier_Quote_NilZone(t *testing.T) {
	_, err := StandardCarrier{}.Quote(context.Background(), QuoteRequest{})
	if err == nil {
		t.Fatal("expected error for nil zone, got nil")
	}
}

// TestStandardCarrier_EstimateDays 验证 4 个国家档位 + 默认档。
func TestStandardCarrier_EstimateDays(t *testing.T) {
	c := StandardCarrier{}
	cases := []struct {
		country string
		want    int
	}{
		{"CN", 3},
		{"US", 7},
		{"JP", 5},
		{"BR", 10}, // default
		{"cn", 3},  // case-insensitive
	}
	for _, tc := range cases {
		if got := c.EstimateDays(tc.country); got != tc.want {
			t.Errorf("EstimateDays(%q)=%d, want %d", tc.country, got, tc.want)
		}
	}
}
