package shipping

import (
	"encoding/json"
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
		NameI18n:            json.RawMessage(`{"en":"US Zone","zh":"美国区域"}`),
		Regions:             Regions{"US-NY", "US-CA"},
		FeeType:             FeeTypeByWeight,
		Taxable:             true,
		TaxRate:             decimal.NewFromFloat(0.0825),
		TaxIncluded:         false,
		IossApplicable:      true,
		RemoteSurcharge:     decimal.NewFromInt(500),
		RemoteZipPatterns:   json.RawMessage(`["99*","88*"]`),
		FuelSurchargePct:    decimal.NewFromFloat(0.1500),
		VolumetricDivisor:   6000,
		Sort:                5,
	}

	if zone.Currency != "USD" {
		t.Errorf("expected Currency=USD, got %s", zone.Currency)
	}
	if !zone.Taxable {
		t.Error("expected Taxable=true")
	}
	if !zone.TaxRate.Equal(decimal.NewFromFloat(0.0825)) {
		t.Errorf("expected TaxRate=0.0825, got %s", zone.TaxRate)
	}
	if !zone.TaxIncluded == zone.TaxIncluded {
		t.Errorf("expected TaxIncluded=false, got %v", zone.TaxIncluded)
	}
	if !zone.IossApplicable {
		t.Error("expected IossApplicable=true")
	}
	if !zone.RemoteSurcharge.Equal(decimal.NewFromInt(500)) {
		t.Errorf("expected RemoteSurcharge=500, got %s", zone.RemoteSurcharge)
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
