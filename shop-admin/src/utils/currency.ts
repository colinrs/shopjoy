/**
 * Currency formatting helpers.
 *
 * The backend stores monetary values as decimal-formatted strings
 * (`"10.00"`) tagged with an ISO 4217 currency code on the parent
 * promotion / coupon row. The frontend never parses these into numbers
 * — we just decorate them with the right symbol for display.
 */

const SYMBOLS: Record<string, string> = {
  CNY: '¥',
  USD: '$',
  EUR: '€',
  JPY: '¥',
  GBP: '£',
  SGD: 'S$',
  HKD: 'HK$'
}

/** Returns the display symbol for an ISO 4217 currency code, or the code itself if unknown. */
export function currencySymbol(code: string | undefined): string {
  if (!code) return ''
  return SYMBOLS[code.toUpperCase()] ?? code
}
