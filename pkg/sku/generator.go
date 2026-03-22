package sku

import (
	"crypto/rand"
	"math/big"
	"strings"
	"time"

	"github.com/colinrs/shopjoy/pkg/code"
)

type generator struct {
	config Config
}

// NewGenerator creates a new SKU generator
func NewGenerator(config Config) Generator {
	return &generator{config: config}
}

// Generate generates a new SKU code
func (g *generator) Generate(tenantID int64, tenantPrefix, productPrefix string) (string, error) {
	// Validate and normalize prefixes
	if err := ValidatePrefix(tenantPrefix); err != nil {
		return "", err
	}
	if err := ValidatePrefix(productPrefix); err != nil {
		return "", err
	}

	tenantPrefix = NormalizePrefix(tenantPrefix)
	productPrefix = NormalizePrefix(productPrefix)

	// Generate compact code
	compactCode, err := g.generateCompactCode(tenantID)
	if err != nil {
		return "", err
	}

	// Build full code
	var parts []string
	if tenantPrefix != "" {
		parts = append(parts, tenantPrefix)
	}
	if productPrefix != "" {
		parts = append(parts, productPrefix)
	}
	parts = append(parts, compactCode)

	result := strings.Join(parts, "-")

	// Validate total length
	if len(result) > g.config.MaxTotalLength {
		return "", code.ErrSKUCodeTooLong
	}

	return result, nil
}

// GenerateWithRetry generates a new SKU code with retry on collision
func (g *generator) GenerateWithRetry(tenantID int64, tenantPrefix, productPrefix string, maxRetry int) (string, error) {
	var lastErr error
	for i := 0; i < maxRetry; i++ {
		code, err := g.Generate(tenantID, tenantPrefix, productPrefix)
		if err == nil {
			return code, nil
		}
		lastErr = err
	}
	return "", lastErr
}

// generateCompactCode generates the 10-character compact code
func (g *generator) generateCompactCode(tenantID int64) (string, error) {
	// Encode tenant ID (4 chars)
	tenantPart := EncodeBase62(tenantID, TenantIDLength)

	// Calculate hour offset from epoch
	now := time.Now().UTC()
	hourOffset := int64(now.Sub(g.config.Epoch).Hours())

	// Encode timestamp (3 chars)
	timestampPart := EncodeBase62(hourOffset, TimestampLength)

	// Generate random sequence (3 chars)
	randomPart, err := g.generateRandomSequence()
	if err != nil {
		return "", code.ErrSKUGenerateFailed
	}

	return tenantPart + timestampPart + randomPart, nil
}

// generateRandomSequence generates a 3-character random Base62 string
func (g *generator) generateRandomSequence() (string, error) {
	max := big.NewInt(62 * 62 * 62) // 238328 combinations
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return EncodeBase62(n.Int64(), RandomLength), nil
}