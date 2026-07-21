package shipping

import (
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
)

func TestShippingTemplate_MarketIsolation(t *testing.T) {
	tmpl := &ShippingTemplate{
		TenantID:  1,
		MarketID:  2, // US market
		Currency:  "USD",
		Name:      "US Standard",
		IsDefault: true,
		IsActive:  true,
	}
	if tmpl.MarketID != 2 {
		t.Errorf("expected MarketID=2, got %d", tmpl.MarketID)
	}
	if tmpl.Currency != "USD" {
		t.Errorf("expected Currency=USD, got %s", tmpl.Currency)
	}
}

func TestShippingZone_Taxable(t *testing.T) {
	zone := &ShippingZone{
		TenantID: 1,
		Currency: "EUR",
		FeeType:  FeeTypeByWeight,
		Taxable:  true,
	}
	if !zone.Taxable {
		t.Error("expected Taxable=true")
	}
	if zone.Currency != "EUR" {
		t.Errorf("expected Currency=EUR, got %s", zone.Currency)
	}
}

func TestShippingZone_I18nAndSurchargeFields(t *testing.T) {
	zone := &ShippingZone{
		TenantID:            1,
		TemplateID:          10,
		MarketID:            2,
		Currency:            "USD",
		Name:                "US Zone",
		NameI18n:            StringI18n{"en": "US Zone", "zh": "美国区域"},
		Regions:             Regions{"US-NY", "US-CA"},
		FeeType:             FeeTypeByWeight,
		Taxable:             true,
		TaxRate:             decimal.NewFromFloat(0.0825),
		TaxIncluded:         false,
		IossApplicable:      true,
		RemoteSurcharge:     decimal.NewFromInt(500),
		RemoteZipPatterns:   StringArray{"99*", "88*"},
		FuelSurchargePct:    decimal.NewFromFloat(0.1500),
		VolumetricDivisor:   6000,
		Sort:                5,
	}

	if zone.Currency != "USD" {
		t.Errorf("expected Currency=USD, got %s", zone.Currency)
	}
	wantNameI18n := StringI18n{"en": "US Zone", "zh": "美国区域"}
	if !reflect.DeepEqual(zone.NameI18n, wantNameI18n) {
		t.Errorf("expected NameI18n=%v, got %v", wantNameI18n, zone.NameI18n)
	}
	if !zone.Taxable {
		t.Error("expected Taxable=true")
	}
	if !zone.TaxRate.Equal(decimal.NewFromFloat(0.0825)) {
		t.Errorf("expected TaxRate=0.0825, got %s", zone.TaxRate)
	}
	if zone.TaxIncluded {
		t.Error("expected TaxIncluded=false")
	}
	if !zone.IossApplicable {
		t.Error("expected IossApplicable=true")
	}
	if !zone.RemoteSurcharge.Equal(decimal.NewFromInt(500)) {
		t.Errorf("expected RemoteSurcharge=500, got %s", zone.RemoteSurcharge)
	}
	wantZipPatterns := StringArray{"99*", "88*"}
	if !reflect.DeepEqual(zone.RemoteZipPatterns, wantZipPatterns) {
		t.Errorf("expected RemoteZipPatterns=%v, got %v", wantZipPatterns, zone.RemoteZipPatterns)
	}
	if !zone.FuelSurchargePct.Equal(decimal.NewFromFloat(0.1500)) {
		t.Errorf("expected FuelSurchargePct=0.1500, got %s", zone.FuelSurchargePct)
	}
	if zone.VolumetricDivisor != 6000 {
		t.Errorf("expected VolumetricDivisor=6000, got %d", zone.VolumetricDivisor)
	}
	if zone.Sort != 5 {
		t.Errorf("expected Sort=5, got %d", zone.Sort)
	}
	if zone.MarketID != 2 {
		t.Errorf("expected MarketID=2, got %d", zone.MarketID)
	}
}

func TestShippingZone_CalculateFee_ByVolume(t *testing.T) {
	zone := &ShippingZone{
		FeeType:           FeeTypeByVolume,
		FirstUnit:         1000, // 首件 1000g
		FirstFee:          decimal.RequireFromString("10.00"),
		AdditionalUnit:    500,
		AdditionalFee:     decimal.RequireFromString("3.00"),
		VolumetricDivisor: 5000,
	}
	// 商品: 长宽高 200mm × 200mm × 200mm → 体积 = 8,000,000 mm³ = 8000 cm³ → 体积重 1600g
	items := []CalculateItem{
		{ProductID: 1, Quantity: 1, Weight: 500, Length: 200, Width: 200, Height: 200, Price: decimal.NewFromInt(20)},
	}
	fee := zone.CalculateFee(items, decimal.NewFromInt(20), 1)
	// 实重 500g, 体积重 1600g, 取大 = 1600g
	// 超过首件 1000g 600g, 600/500=1.2, 向上取整 = 2 续件
	// fee = 10 + 3*2 = 16
	if !fee.Equal(decimal.RequireFromString("16.00")) {
		t.Errorf("expected 16.00, got %s", fee)
	}
}
