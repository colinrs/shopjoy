package utils

import (
	"time"

	"github.com/shopspring/decimal"
)

// FormatAmount formats a decimal amount to a string with 2 decimal places.
func FormatAmount(amount decimal.Decimal) string {
	if amount.IsZero() {
		return "0.00"
	}
	return amount.StringFixed(2)
}

// FormatAmountWithCurrency formats a decimal amount with currency prefix.
// The currency parameter is ignored in this implementation as the format
// is primarily for display purposes.
func FormatAmountWithCurrency(amount decimal.Decimal, currency string) string {
	return FormatAmount(amount)
}

// FormatDecimal formats a decimal.Decimal to string representation.
func FormatDecimal(d decimal.Decimal) string {
	return d.String()
}

// FormatDecimalToString formats a decimal.Decimal to string with 2 decimal places.
func FormatDecimalToString(v decimal.Decimal) string {
	if v.IsZero() {
		return "0"
	}
	return v.StringFixed(2)
}

// FormatTimeToRFC3339 formats a time.Time pointer to RFC3339 string.
// Returns empty string if t is nil.
func FormatTimeToRFC3339(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

// FormatTimeForExport formats a time.Time pointer to export format (YYYY-MM-DD HH:mm:ss).
// Returns empty string if t is nil.
func FormatTimeForExport(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// FormatTimeStringForExport is a no-op function that returns the input string.
// It exists for compatibility with code that passes pre-formatted time strings.
func FormatTimeStringForExport(t string) string {
	return t
}
