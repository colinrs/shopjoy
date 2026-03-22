package sku

// EncodeBase62 encodes an integer to a Base62 string with fixed length
func EncodeBase62(n int64, length int) string {
	if length <= 0 {
		length = 1
	}

	result := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		result[i] = base62Chars[n%62]
		n /= 62
	}

	return string(result)
}

// DecodeBase62 decodes a Base62 string to an integer
func DecodeBase62(s string) int64 {
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
		}
	}
	return result
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