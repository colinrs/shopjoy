package sku

import (
	"testing"
)

func TestEncodeBase62(t *testing.T) {
	tests := []struct {
		input    int64
		length   int
		expected string
	}{
		{0, 4, "0000"},
		{1, 4, "0001"},
		{10, 4, "000A"},
		{35, 4, "000Z"},
		{36, 4, "000a"},
		{61, 4, "000z"},
		{62, 4, "0010"},
		{12345, 4, "03D7"},
		{1000000, 4, "4C92"},
		{14776335, 4, "zzzz"}, // Max 4-digit
		{0, 3, "000"},
		{238327, 3, "zzz"}, // Max 3-digit
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := EncodeBase62(tt.input, tt.length)
			if result != tt.expected {
				t.Errorf("EncodeBase62(%d, %d) = %q, want %q", tt.input, tt.length, result, tt.expected)
			}
		})
	}
}

func TestDecodeBase62(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"0000", 0},
		{"0001", 1},
		{"000A", 10},
		{"000Z", 35},
		{"000a", 36},
		{"000z", 61},
		{"0010", 62},
		{"03D7", 12345},
		{"4C92", 1000000},
		{"zzzz", 14776335},
		{"000", 0},
		{"zzz", 238327},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := DecodeBase62(tt.input)
			if result != tt.expected {
				t.Errorf("DecodeBase62(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	for i := int64(0); i < 10000; i++ {
		encoded := EncodeBase62(i, 4)
		decoded := DecodeBase62(encoded)
		if decoded != i {
			t.Errorf("Round trip failed: %d -> %q -> %d", i, encoded, decoded)
		}
	}
}

func TestIsBase62(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"0000", true},
		{"ABC123abc", true},
		{"zzzz", true},
		{"ZZZZ", true},
		{"", true},
		{"ABC-123", false},
		{"ABC_123", false},
		{"ABC 123", false},
		{"中文", false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := IsBase62(tt.input)
			if result != tt.expected {
				t.Errorf("IsBase62(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}