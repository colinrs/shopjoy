/**
 * Shared status type utilities for Element Plus tag types.
 * Each domain has its own status values and type mappings.
 */

/**
 * Review status types: pending -> warning, approved -> success, hidden -> info
 */
export function getReviewStatusType(status: string): string {
  const types: Record<string, string> = {
    pending: 'warning',
    approved: 'success',
    hidden: 'info'
  }
  return types[status] || 'info'
}

/**
 * Product status types: on_sale -> success, off_sale -> warning, draft -> info
 */
export function getProductStatusType(status: string): string {
  const types: Record<string, string> = {
    on_sale: 'success',
    off_sale: 'warning',
    draft: 'info'
  }
  return types[status] || 'info'
}

/**
 * Payment transaction status types: 0 -> warning, 1 -> success, 2 -> danger
 */
export function getPaymentStatusType(status: number): string {
  const types: Record<number, string> = {
    0: 'warning',
    1: 'success',
    2: 'danger'
  }
  return types[status] || 'info'
}

/**
 * Promotion status types: active -> success, paused -> warning, pending/ended -> info
 */
export function getPromotionStatusType(status: string): string {
  const types: Record<string, string> = {
    active: 'success',
    paused: 'warning',
    pending: 'info',
    ended: 'info'
  }
  return types[status] || 'info'
}

/**
 * Order status types for dashboard: pending_payment -> warning, paid -> primary,
 * shipped -> info, delivered -> success, cancelled -> danger, refunded -> danger
 */
export function getOrderStatusType(status: string): string {
  const types: Record<string, string> = {
    pending_payment: 'warning',
    paid: 'primary',
    shipped: 'info',
    delivered: 'success',
    cancelled: 'danger',
    refunded: 'danger'
  }
  return types[status] || 'info'
}
