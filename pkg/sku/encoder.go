package sku

import (
	"github.com/colinrs/shopjoy/pkg/code"
)

// maxValues for each length
var maxValues = map[int]int64{
	3: 238327,    // 62^3 - 1
	4: 14776335,  // 62^4 - 1
	10: 8392993658683402, // 62^10 - 1
}

// EncodeBase62 encodes an integer to a Base62 string with fixed length.
// Returns error if the value exceeds the capacity for the specified length.
func EncodeBase62(n int64, length int) (string, error) {
	if length <= 0 {
		length = 1
	}

	// Validate value fits in the specified length
	maxVal, ok := maxValues[length]
	if !ok {
		// Calculate max for unknown lengths
		maxVal = 0
		for i := 0; i < length; i++ {
			maxVal = maxVal*62 + 61
		}
	}

	if n < 0 || n > maxVal {
		return "", code.ErrSKUGenerateFailed
	}

	result := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		result[i] = base62Chars[n%62]
		n /= 62
	}

	return string(result), nil
}

// DecodeBase62 decodes a Base62 string to an integer.
// Returns error if the input contains invalid characters.
func DecodeBase62(s string) (int64, error) {
	var result int64
	for _, c := range s {
		result *= 62
		switch {
		case c >= '0' && c <= '9':
			result += int64(c - '0')
		case c >= 'A' && c <= 'Z':
			result += int64(c - 'A' + 10)
		case c >= 'a' && c <= 'z':
			result += int64(c - 'a' + 36)
		default:
			return 0, code.ErrSKUParseFailed
		}
	}
	return result, nil
}

// IsBase62 checks if a string contains only Base62 characters
func IsBase62(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= 'a' && c <= 'z')) {
			return false
		}
	}
	return true
}