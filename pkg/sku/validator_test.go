package sku

import (
	"testing"

	"github.com/colinrs/shopjoy/pkg/code"
)

func TestValidatePrefix(t *testing.T) {
	tests := []struct {
		input    string
		expected error
	}{
		{"NIKE", nil},
		{"SHOE", nil},
		{"A1", nil},
		{"PROD", nil},
		{"AbCdEf", nil},
		{"", nil}, // Empty is valid (no prefix)
		{"123", code.ErrSKUPrefixStartsWithNumber},
		{"1ABC", code.ErrSKUPrefixStartsWithNumber},
		{"NIKE-SHOE", code.ErrSKUPrefixInvalid}, // Contains hyphen
		{"NIKE_2024", code.ErrSKUPrefixInvalid}, // Contains underscore
		{"NIKE 2024", code.ErrSKUPrefixInvalid}, // Contains space
		{"NIKE中文", code.ErrSKUPrefixInvalid},    // Contains non-ASCII
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			err := ValidatePrefix(tt.input)
			if tt.expected == nil && err != nil {
				t.Errorf("ValidatePrefix(%q) = %v, want nil", tt.input, err)
			}
			if tt.expected != nil && err == nil {
				t.Errorf("ValidatePrefix(%q) = nil, want %v", tt.input, tt.expected)
			}
		})
	}
}

func TestValidatePrefixLength(t *testing.T) {
	// Test max length
	validPrefix := "ABCDEFGH" // 8 chars
	if err := ValidatePrefix(validPrefix); err != nil {
		t.Errorf("ValidatePrefix(%q) should be valid", validPrefix)
	}

	// Test too long
	longPrefix := "ABCDEFGHI" // 9 chars
	if err := ValidatePrefix(longPrefix); err != code.ErrSKUPrefixTooLong {
		t.Errorf("ValidatePrefix(%q) = %v, want ErrSKUPrefixTooLong", longPrefix, err)
	}
}

func TestNormalizePrefix(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"nike", "NIKE"},
		{"Nike", "NIKE"},
		{"NIKE", "NIKE"},
		{"AbCdEf", "ABCDEF"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := NormalizePrefix(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizePrefix(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
