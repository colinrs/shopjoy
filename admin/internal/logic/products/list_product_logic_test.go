package products

import (
	"testing"

	"github.com/shopspring/decimal"

	appProduct "github.com/colinrs/shopjoy/admin/internal/application/product"
)

// TestConvertToProductDetailRespWithMarkets_PopulatesImages is the regression
// guard for the "product list thumbnails don't show" bug: the LIST endpoint's
// converter previously returned resp.Images == nil even when the underlying
// ProductResponse had images, because it only set basic fields + markets.
func TestConvertToProductDetailRespWithMarkets_PopulatesImages(t *testing.T) {
	p := &appProduct.ProductResponse{
		ID:         42,
		Name:       "Test Product",
		Price:      decimal.NewFromInt(99),
		Currency:   "USD",
		SKU:        "SKU-1",
		Brand:      "Acme",
		Tags:       []string{"hot"},
		Images:     []string{"https://cdn.example.com/a.jpg", "https://cdn.example.com/b.jpg"},
		HSCode:     "HS-1",
		COO:        "CN",
		Weight:     decimal.NewFromInt(100),
		WeightUnit: "g",
		Length:     decimal.NewFromInt(10),
		Width:      decimal.NewFromInt(20),
		Height:     decimal.NewFromInt(30),
	}

	resp := convertToProductDetailRespWithMarkets(p, nil, nil)

	if got := len(resp.Images); got != 2 {
		t.Fatalf("Images length = %d, want 2 (list endpoint was dropping images)", got)
	}
	if resp.Images[0] != "https://cdn.example.com/a.jpg" || resp.Images[1] != "https://cdn.example.com/b.jpg" {
		t.Errorf("Images = %v, want both URLs preserved", resp.Images)
	}
	if resp.SKU != "SKU-1" {
		t.Errorf("SKU = %q, want %q", resp.SKU, "SKU-1")
	}
	if resp.Brand != "Acme" {
		t.Errorf("Brand = %q, want %q", resp.Brand, "Acme")
	}
	if len(resp.Tags) != 1 || resp.Tags[0] != "hot" {
		t.Errorf("Tags = %v, want [hot]", resp.Tags)
	}
}
