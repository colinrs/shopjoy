package shipping

import (
	"math"
	"testing"
)

func TestWeightConverter_ToGrams(t *testing.T) {
	tests := []struct {
		name  string
		unit  WeightUnit
		value float64
		want  int
	}{
		{"gram passthrough", WeightUnitGram, 500, 500},
		{"gram rounds", WeightUnitGram, 500.6, 501},
		{"kilogram to grams", WeightUnitKilogram, 1.5, 1500},
		{"pound to grams", WeightUnitPound, 1, 454},           // 453.59237 -> 454
		{"pound multiple", WeightUnitPound, 2, 907},           // 907.18474 -> 907
		{"ounce to grams", WeightUnitOunce, 1, 28},            // 28.349523125 -> 28
		{"ounce multiple", WeightUnitOunce, 16, 454},          // 453.5924 -> 454
		{"zero", WeightUnitGram, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewWeightConverter(tt.unit)
			got := c.ToGrams(tt.value)
			if got != tt.want {
				t.Errorf("ToGrams(%v) with unit %v = %d, want %d", tt.value, tt.unit, got, tt.want)
			}
		})
	}
}

func TestWeightConverter_FromGrams(t *testing.T) {
	const eps = 1e-6
	tests := []struct {
		name  string
		unit  WeightUnit
		grams int
		want  float64
	}{
		{"gram passthrough", WeightUnitGram, 500, 500},
		{"grams to kilogram", WeightUnitKilogram, 1500, 1.5},
		{"grams to pound", WeightUnitPound, 454, 454.0 / 453.59237},
		{"grams to ounce", WeightUnitOunce, 28, 28.0 / 28.349523125},
		{"zero", WeightUnitKilogram, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewWeightConverter(tt.unit)
			got := c.FromGrams(tt.grams)
			if math.Abs(got-tt.want) > eps {
				t.Errorf("FromGrams(%d) with unit %v = %v, want %v", tt.grams, tt.unit, got, tt.want)
			}
		})
	}
}
