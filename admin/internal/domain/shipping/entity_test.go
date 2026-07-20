package shipping

import "testing"

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