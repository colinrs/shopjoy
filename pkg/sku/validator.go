package sku

import (
	"github.com/colinrs/shopjoy/pkg/code"
)

// ValidatePrefix validates a SKU prefix
// Rules:
// 1. Only letters and numbers allowed
// 2. Length 0-8 characters (empty means no prefix)
// 3. Cannot start with a number
func ValidatePrefix(prefix string) error {
	if prefix == "" {
		return nil
	}

	// Check length
	if len(prefix) > MaxTenantPrefixLength {
		return code.ErrSKUPrefixTooLong
	}

	// Check first character is not a digit
	if prefix[0] >= '0' && prefix[0] <= '9' {
		return code.ErrSKUPrefixStartsWithNumber
	}

	// Check all characters are alphanumeric
	for _, c := range prefix {
		if !((c >= 'A' && c <= 'Z') ||
			(c >= 'a' && c <= 'z') ||
			(c >= '0' && c <= '9')) {
			return code.ErrSKUPrefixInvalid
		}
	}

	return nil
}

// NormalizePrefix converts prefix to uppercase
func NormalizePrefix(prefix string) string {
	result := make([]byte, len(prefix))
	for i, c := range prefix {
		if c >= 'a' && c <= 'z' {
			result[i] = byte(c - 32) // Convert to uppercase
		} else {
			result[i] = byte(c)
		}
	}
	return string(result)
}
