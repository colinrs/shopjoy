package sku

import (
	"strings"
	"time"

	"github.com/colinrs/shopjoy/pkg/code"
)

type parser struct {
	config Config
}

// NewParser creates a new SKU parser
func NewParser(config Config) Parser {
	return &parser{config: config}
}

// Parse parses a SKU code and returns its components
func (p *parser) Parse(skuCode string) (*SKUInfo, error) {
	if skuCode == "" {
		return nil, code.ErrSKUParseFailed
	}

	// Split by hyphen
	parts := strings.Split(skuCode, "-")

	// Last part should be compact code
	compactCode := parts[len(parts)-1]

	// Validate compact code
	if len(compactCode) != CompactCodeLength {
		return nil, code.ErrSKUParseFailed
	}
	if !IsBase62(compactCode) {
		return nil, code.ErrSKUParseFailed
	}

	// Decode tenant ID (first 4 chars)
	tenantID, err := DecodeBase62(compactCode[:TenantIDLength])
	if err != nil {
		return nil, code.ErrSKUParseFailed
	}

	// Decode timestamp (next 3 chars)
	hourOffset, err := DecodeBase62(compactCode[TenantIDLength : TenantIDLength+TimestampLength])
	if err != nil {
		return nil, code.ErrSKUParseFailed
	}
	createdAt := p.config.Epoch.Add(time.Duration(hourOffset) * time.Hour)

	// Extract random sequence (last 3 chars)
	randomSequence := compactCode[TenantIDLength+TimestampLength:]

	// Extract prefixes
	var tenantPrefix, productPrefix string
	prefixes := parts[:len(parts)-1]
	if len(prefixes) > 0 {
		tenantPrefix = prefixes[0]
	}
	if len(prefixes) > 1 {
		productPrefix = prefixes[1]
	}

	return &SKUInfo{
		TenantPrefix:   tenantPrefix,
		ProductPrefix:  productPrefix,
		TenantID:       tenantID,
		CreatedAt:      createdAt,
		RandomSequence: randomSequence,
		CompactCode:    compactCode,
	}, nil
}

// ExtractTenantID extracts only the tenant ID from a SKU code
func (p *parser) ExtractTenantID(skuCode string) (int64, error) {
	if skuCode == "" {
		return 0, code.ErrSKUParseFailed
	}

	// Find the last hyphen
	idx := strings.LastIndex(skuCode, "-")
	compactCode := skuCode
	if idx >= 0 {
		compactCode = skuCode[idx+1:]
	}

	// Validate compact code
	if len(compactCode) != CompactCodeLength {
		return 0, code.ErrSKUParseFailed
	}
	if !IsBase62(compactCode) {
		return 0, code.ErrSKUParseFailed
	}

	// Decode tenant ID (first 4 chars)
	return DecodeBase62(compactCode[:TenantIDLength])
}