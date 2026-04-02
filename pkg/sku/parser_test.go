package sku

import (
	"strings"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	cfg := DefaultConfig()
	parser := NewParser(cfg)

	// Generate a code first
	gen := NewGenerator(cfg)
	code, _ := gen.Generate(12345, "NIKE", "SHOE")

	info, err := parser.Parse(code)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if info.TenantPrefix != "NIKE" {
		t.Errorf("Expected tenant prefix NIKE, got %s", info.TenantPrefix)
	}
	if info.ProductPrefix != "SHOE" {
		t.Errorf("Expected product prefix SHOE, got %s", info.ProductPrefix)
	}
	if info.TenantID != 12345 {
		t.Errorf("Expected tenant ID 12345, got %d", info.TenantID)
	}
	if len(info.CompactCode) != CompactCodeLength {
		t.Errorf("Expected compact code length %d, got %d", CompactCodeLength, len(info.CompactCode))
	}
}

func TestParser_ParseNoPrefix(t *testing.T) {
	cfg := DefaultConfig()
	parser := NewParser(cfg)

	gen := NewGenerator(cfg)
	code, _ := gen.Generate(12345, "", "")

	info, err := parser.Parse(code)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if info.TenantPrefix != "" {
		t.Errorf("Expected empty tenant prefix, got %s", info.TenantPrefix)
	}
	if info.ProductPrefix != "" {
		t.Errorf("Expected empty product prefix, got %s", info.ProductPrefix)
	}
	if info.TenantID != 12345 {
		t.Errorf("Expected tenant ID 12345, got %d", info.TenantID)
	}
}

func TestParser_ExtractTenantID(t *testing.T) {
	cfg := DefaultConfig()
	parser := NewParser(cfg)

	tests := []struct {
		code string
	}{
		{"NIKE-SHOE-0001001ABC"},
		{"NIKE-00D7002ABC"}, // 10-char compact code
		{"4C92003XYZ"},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			// Extract tenant ID directly from known compact code
			compactCode := tt.code
			if idx := strings.LastIndex(tt.code, "-"); idx >= 0 {
				compactCode = tt.code[idx+1:]
			}
			expected, err := DecodeBase62(compactCode[:TenantIDLength])
			if err != nil {
				t.Fatalf("DecodeBase62 failed: %v", err)
			}

			tenantID, err := parser.ExtractTenantID(tt.code)
			if err != nil {
				t.Fatalf("ExtractTenantID failed: %v", err)
			}
			if tenantID != expected {
				t.Errorf("Expected tenant ID %d, got %d", expected, tenantID)
			}
		})
	}
}

func TestParser_ParseInvalidCodes(t *testing.T) {
	cfg := DefaultConfig()
	parser := NewParser(cfg)

	tests := []string{
		"",                 // Empty
		"NIKE",             // No compact code
		"NIKE-SHOE",        // No compact code
		"ABC-123-XYZ",      // Invalid characters in compact
		"NIKE-123",         // Compact too short
		"NIKE-12345678901", // Compact too long
	}

	for _, code := range tests {
		t.Run(code, func(t *testing.T) {
			_, err := parser.Parse(code)
			if err == nil {
				t.Errorf("Expected error for code %q", code)
			}
		})
	}
}

func TestParser_ExtractTenantIDInvalidCodes(t *testing.T) {
	cfg := DefaultConfig()
	parser := NewParser(cfg)

	tests := []string{
		"",
		"NIKE",
		"ABC-123", // Too short
	}

	for _, code := range tests {
		t.Run(code, func(t *testing.T) {
			_, err := parser.ExtractTenantID(code)
			if err == nil {
				t.Errorf("Expected error for code %q", code)
			}
		})
	}
}
