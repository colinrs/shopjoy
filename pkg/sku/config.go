package sku

import "time"

const (
	// Compact code structure
	CompactCodeLength = 10
	TenantIDLength    = 4
	TimestampLength   = 3
	RandomLength      = 3

	// Prefix limits
	MaxTenantPrefixLength  = 8
	MaxProductPrefixLength = 8
	MaxTotalLength         = 28

	// Default epoch: 2024-01-01 00:00:00 UTC
	DefaultEpoch = 1704067200
)

// Config holds SKU generator configuration
type Config struct {
	Epoch                  time.Time
	MaxTenantPrefixLength  int
	MaxProductPrefixLength int
	MaxTotalLength         int
	MaxRetry               int
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		Epoch:                  time.Unix(DefaultEpoch, 0).UTC(),
		MaxTenantPrefixLength:  MaxTenantPrefixLength,
		MaxProductPrefixLength: MaxProductPrefixLength,
		MaxTotalLength:         MaxTotalLength,
		MaxRetry:               3,
	}
}