import request from '@/utils/request'

// ===================== Types =====================

// Order status (backend: admin/internal/domain/fulfillment/order_repository.go OrderStatus)
// Values: pending_payment, paid, shipped, delivered, cancelled, refunded
export type OrderStatus = 'pending_payment' | 'paid' | 'shipped' | 'delivered' | 'cancelled' | 'refunded'

// Fulfillment status
export type FulfillmentStatus = 'pending' | 'partial_shipped' | 'shipped' | 'delivered'

// Order list query parameters
export interface OrderListParams {
  page?: number
  page_size?: number
  status?: OrderStatus
  fulfillment_status?: FulfillmentStatus
  order_no?: string
  user_phone?: string
  start_time?: string
  end_time?: string
}

// Order item
export interface OrderItem {
  order_item_id: number
  product_id: number
  product_name: string
  sku_id: number
  sku_name: string
  image: string
  quantity: number
  shipped_qty: number
  pending_qty: number
  unit_price: string
  total_price: string
}

// Shipping address
export interface ShippingAddress {
  receiver_name: string
  receiver_phone: string
  province: string
  city: string
  district: string
  address: string
  full_address: string
}

// Payment info
export interface PaymentInfo {
  payment_method: string
  payment_method_name: string
  payment_no: string
  paid_at: string
  pay_amount: string
  discount_amount: string
  currency: string
}

// Shipment info
export interface ShipmentInfo {
  shipment_id: number
  shipment_no: string
  carrier_code: string
  carrier_name: string
  tracking_no: string
  status: string
  shipped_at: string
  delivered_at: string
}

// Order list item
export interface Order {
  order_id: string
  order_no: string
  status: OrderStatus
  status_text: string
  fulfillment_status: FulfillmentStatus
  fulfillment_text: string
  refund_status: string
  refund_text: string
  user_id: number
  user_name: string
  user_phone: string
  total_amount: string
  pay_amount: string
  discount_amount: string
  currency: string
  item_count: number
  payment_method: string
  payment_method_name: string
  order_source: string
  buyer_remark: string
  seller_remark: string
  created_at: string
  paid_at: string
  items: OrderItem[]
}

// Order detail
export interface OrderDetail extends Order {
  shipping_address: ShippingAddress
  payment_info: PaymentInfo
  shipments: ShipmentInfo[]
}

// Ship order request item
export interface ShipOrderItem {
  order_item_id: number
  quantity: number
}

// Ship order request
export interface ShipOrderRequest {
  carrier_code: string
  carrier_name?: string
  tracking_no: string
  shipping_cost?: string
  currency?: string
  weight?: string
  remark?: string
  items?: ShipOrderItem[] // Empty means all items
}

// Ship order response
export interface ShipOrderResponse {
  shipment_id: number
  shipment_no: string
}

// Adjust price request
export interface AdjustPriceRequest {
  adjust_amount: string // Positive = increase, Negative = decrease
  reason: string // Required, max 200 characters
}

// Adjust price response
export interface AdjustPriceResponse {
  order_id: number
  original_amount: string
  adjust_amount: string
  new_pay_amount: string
  adjust_reason: string
  adjusted_at: string
}

// Update order remark request
export interface UpdateOrderRemarkRequest {
  remark: string // Max 500 characters
}

// Fulfillment summary
export interface FulfillmentSummary {
  pending_shipment: number
  partial_shipped: number
  shipped: number
  delivered: number
  pending_refund: number
  refunding: number
  total_orders: number
  today_orders: number
  today_gmv: string
}

// Export orders parameters
export interface ExportOrdersParams {
  order_no?: string
  user_id?: number
  status?: OrderStatus
  fulfillment_status?: FulfillmentStatus
  refund_status?: string
  start_time?: string
  end_time?: string
}

// ===================== API Functions =====================

/**
 * Get order list with fulfillment status filtering
 */
export const getOrderList = (params: OrderListParams) => {
  return request.get<{ list: Order[]; total: number; page: number; page_size: number }>(
    '/api/v1/orders',
    { params }
  )
}

/**
 * Get order detail with fulfillment information
 */
export const getOrderDetail = (orderId: string) => {
  return request.get<OrderDetail>(`/api/v1/orders/${orderId}`)
}

/**
 * Ship order (create shipment from order)
 */
export const shipOrder = (orderId: string, data: ShipOrderRequest) => {
  return request.put<ShipOrderResponse>(`/api/v1/orders/${orderId}/ship`, data)
}

/**
 * Update order remark
 */
export const updateOrderRemark = (orderId: string, data: UpdateOrderRemarkRequest) => {
  return request.put<{ order_id: number }>(`/api/v1/orders/${orderId}/remark`, data)
}

/**
 * Adjust order price
 */
export const adjustOrderPrice = (orderId: string, data: AdjustPriceRequest) => {
  return request.put<AdjustPriceResponse>(`/api/v1/orders/${orderId}/adjust-price`, data)
}

/**
 * Export orders
 */
export const exportOrders = (params: ExportOrdersParams) => {
  return request.get<Blob>('/api/v1/orders/export', {
    params,
    responseType: 'blob'
  })
}

/**
 * Get fulfillment summary statistics
 */
export const getFulfillmentSummary = () => {
  return request.get<FulfillmentSummary>('/api/v1/orders/fulfillment-summary')
}

/**
 * Cancel order
 */
export const cancelOrder = (orderId: string, reason: string) => {
  return request.put<void>(`/api/v1/orders/${orderId}/cancel`, { reason })
}

/**
 * Remind payment
 */
export const remindPayment = (orderId: string) => {
  return request.post<void>(`/api/v1/orders/${orderId}/remind-payment`)
}

/**
 * Batch ship orders - delegates to fulfillment module
 * Note: Use batchCreateShipments from fulfillment.ts for the actual API call
 */
export const batchShipOrders = (data: {
  order_ids: string[]
  carrier_code: string
  tracking_no_start: string
}) => {
  return request.post<{ shipments: { order_id: string; shipment_id: number }[] }>(
    '/api/v1/shipments/batch',
    {
      order_ids: data.order_ids,
      carrier_code: data.carrier_code,
      tracking_no_start: data.tracking_no_start
    }
  )
}

// Re-export Carrier type and getCarrierList from fulfillment module to avoid duplication
export type { Carrier } from './fulfillment'
export { getCarrierList } from './fulfillment'