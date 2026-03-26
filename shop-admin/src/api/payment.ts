import request from '@/utils/request'

// Types

export interface PaymentStats {
  today_received: string
  today_growth: string
  period_received: string
  refund_amount: string
  refund_rate: string
  currency: string
  channel_distribution: ChannelDistribution[]
}

export interface ChannelDistribution {
  name: string
  percent: string
  amount: string
  count: number
  color: string
}

export interface Transaction {
  id: number
  transaction_id: string
  order_id: string
  order_no: string
  payment_method: string
  payment_method_text: string
  channel_transaction_id: string
  amount: string
  currency: string
  transaction_fee: string
  status: number
  status_text: string
  created_at: string
  paid_at: string | null
}

export interface TransactionListParams {
  page: number
  page_size: number
  order_no?: string
  transaction_id?: string
  payment_method?: string
  status?: number
  start_time?: string
  end_time?: string
}

export interface TransactionListResponse {
  list: Transaction[]
  total: number
  page: number
  page_size: number
  stats: {
    success: number
    pending: number
    failed: number
  }
}

export interface OrderPayment {
  payment_id: number
  payment_no: string
  order_id: string
  order_no: string
  payment_method: string
  payment_method_text: string
  channel_intent_id: string
  channel_payment_id: string
  amount: string
  currency: string
  transaction_fee: string
  fee_currency: string
  status: number
  status_text: string
  paid_at: string | null
  refunded_amount: string
  refunds: PaymentRefund[]
}

export interface PaymentRefund {
  id: number
  refund_no: string
  channel_refund_id: string
  amount: string
  currency: string
  status: number
  status_text: string
  reason_type: string
  reason: string
  refunded_at: string | null
  created_at: string
}

export interface InitiateRefundParams {
  idempotency_key: string
  amount: string
  reason_type: string
  reason?: string
}

export interface InitiateRefundResponse {
  refund_id: number
  refund_no: string
  amount: string
  currency: string
  status: number
  status_text: string
  channel_refund_id: string | null
}

// API Functions

/**
 * Get payment statistics for dashboard
 * @param period - Time period: '7d', '30d', '90d'
 */
export const getPaymentStats = (period: string = '7d') => {
  return request.get<PaymentStats>('/api/v1/payments/stats', { params: { period } })
}

/**
 * Get transaction list with filtering and pagination
 */
export const getTransactionList = (params: TransactionListParams) => {
  return request.get<TransactionListResponse>('/api/v1/payments/transactions', { params })
}

/**
 * Get transaction detail by ID
 */
export const getTransactionDetail = (id: number) => {
  return request.get<Transaction>(`/api/v1/payments/transactions/${id}`)
}

/**
 * Get payment details for an order
 */
export const getOrderPayment = (orderId: number) => {
  return request.get<OrderPayment>(`/api/v1/orders/${orderId}/payment`)
}

/**
 * Initiate a refund for an order
 */
export const initiateRefund = (orderId: number, data: InitiateRefundParams) => {
  return request.post<InitiateRefundResponse>(`/api/v1/orders/${orderId}/refund`, data)
}