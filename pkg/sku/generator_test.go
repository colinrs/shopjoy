package sku

import (
	"strings"
	"testing"
)

func TestGenerator_Generate(t *testing.T) {
	cfg := DefaultConfig()
	gen := NewGenerator(cfg)

	code, err := gen.Generate(1, "NIKE", "SHOE")
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Check format: NIKE-SHOE-<compact>
	parts := strings.Split(code, "-")
	if len(parts) != 3 {
		t.Errorf("Expected 3 parts, got %d: %s", len(parts), code)
	}
	if parts[0] != "NIKE" {
		t.Errorf("Expected tenant prefix NIKE, got %s", parts[0])
	}
	if parts[1] != "SHOE" {
		t.Errorf("Expected product prefix SHOE, got %s", parts[1])
	}
	if len(parts[2]) != CompactCodeLength {
		t.Errorf("Expected compact code length %d, got %d", CompactCodeLength, len(parts[2]))
	}
}

func TestGenerator_GenerateNoPrefix(t *testing.T) {
	cfg := DefaultConfig()
	gen := NewGenerator(cfg)

	code, err := gen.Generate(12345, "", "")
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Should be just compact code
	if len(code) != CompactCodeLength {
		t.Errorf("Expected length %d, got %d: %s", CompactCodeLength, len(code), code)
	}
}

func TestGenerator_GenerateTenantPrefixOnly(t *testing.T) {
	cfg := DefaultConfig()
	gen := NewGenerator(cfg)

	code, err := gen.Generate(1, "NIKE", "")
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	parts := strings.Split(code, "-")
	if len(parts) != 2 {
		t.Errorf("Expected 2 parts, got %d: %s", len(parts), code)
	}
	if parts[0] != "NIKE" {
		t.Errorf("Expected prefix NIKE, got %s", parts[0])
	}
}

func TestGenerator_GenerateTenantIDEncoding(t *testing.T) {
	cfg := DefaultConfig()
	gen := NewGenerator(cfg)

	// Generate SKU for tenant ID 12345
	code, err := gen.Generate(12345, "", "")
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Extract tenant ID from compact code
	compactCode := code
	if idx := strings.LastIndex(code, "-"); idx >= 0 {
		compactCode = code[idx+1:]
	}

	tenantIDPart := compactCode[:TenantIDLength]
	decodedTenantID, err := DecodeBase62(tenantIDPart)
	if err != nil {
		t.Fatalf("DecodeBase62 failed: %v", err)
	}

	if decodedTenantID != 12345 {
		t.Errorf("Expected tenant ID 12345, got %d", decodedTenantID)
	}
}

func TestGenerator_GenerateUniqueness(t *testing.T) {
	cfg := DefaultConfig()
	gen := NewGenerator(cfg)

	codes := make(map[string]bool)
	// Test 200 generations - reasonable number to verify uniqueness
	// With 238,328 random combinations, birthday problem makes 200 safe
	for i := 0; i < 200; i++ {
		code, err := gen.Generate(1, "", "")
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}
		if codes[code] {
			t.Errorf("Duplicate code generated: %s", code)
		}
		codes[code] = true
	}
}

func TestGenerator_GenerateMultiTenantUniqueness(t *testing.T) {
	cfg := DefaultConfig()
	gen := NewGenerator(cfg)

	codes := make(map[string]bool)
	// Generate across multiple tenants - should be unique due to tenant ID encoding
	for tenantID := int64(1); tenantID <= 100; tenantID++ {
		for i := 0; i < 10; i++ {
			code, err := gen.Generate(tenantID, "", "")
			if err != nil {
				t.Fatalf("Generate failed: %v", err)
			}
			if codes[code] {
				t.Errorf("Duplicate code generated: %s", code)
			}
			codes[code] = true
		}
	}
}

func TestGenerator_GenerateTotalLengthLimit(t *testing.T) {
	cfg := DefaultConfig()
	gen := NewGenerator(cfg)

	// Max prefixes: 8 + 8 chars
	code, err := gen.Generate(1, "ABCDEFGH", "IJKLMNOP")
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Total: 8 + 1 + 8 + 1 + 10 = 28
	if len(code) > MaxTotalLength {
		t.Errorf("Code length %d exceeds max %d: %s", len(code), MaxTotalLength, code)
	}
}

func TestGenerator_GenerateMaxTenantID(t *testing.T) {
	cfg := DefaultConfig()
	gen := NewGenerator(cfg)

	// Test max valid tenant ID (14776335 = zzzz in Base62)
	code, err := gen.Generate(14776335, "", "")
	if err != nil {
		t.Fatalf("Generate failed for max tenant ID: %v", err)
	}

	// Verify tenant ID can be decoded
	compactCode := code
	if idx := strings.LastIndex(code, "-"); idx >= 0 {
		compactCode = code[idx+1:]
	}
	decodedTenantID, err := DecodeBase62(compactCode[:TenantIDLength])
	if err != nil {
		t.Fatalf("DecodeBase62 failed: %v", err)
	}
	if decodedTenantID != 14776335 {
		t.Errorf("Expected tenant ID 14776335, got %d", decodedTenantID)
	}
}

func TestGenerator_GenerateOverflowTenantID(t *testing.T) {
	cfg := DefaultConfig()
	gen := NewGenerator(cfg)

	// Test tenant ID that exceeds 4-char capacity
	_, err := gen.Generate(14776336, "", "") // Max is 14776335
	if err == nil {
		t.Error("Expected error for tenant ID exceeding capacity")
	}
}