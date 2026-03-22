package sku

import "time"

// SKUInfo contains parsed SKU information
type SKUInfo struct {
	TenantPrefix   string    // Tenant prefix
	ProductPrefix  string    // Product prefix
	TenantID       int64     // Decoded tenant ID
	CreatedAt      time.Time // Generation time (hour precision)
	RandomSequence string    // Random sequence part
	CompactCode    string    // Full compact code
}

// Generator generates SKU codes
type Generator interface {
	Generate(tenantID int64, tenantPrefix, productPrefix string) (string, error)
	GenerateWithRetry(tenantID int64, tenantPrefix, productPrefix string, maxRetry int) (string, error)
}

// Parser parses SKU codes
type Parser interface {
	Parse(code string) (*SKUInfo, error)
	ExtractTenantID(code string) (int64, error)
}

// Base62 character set: 0-9, A-Z, a-z
const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"