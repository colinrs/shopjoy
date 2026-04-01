import request from '@/utils/request'

// ===================== Types =====================

// Order status (backend: admin/internal/domain/fulfillment/order_repository.go OrderStatus)
// Values: pending_payment, paid, shipped, delivered, cancelled, refunded
export type OrderStatus = 'pending_payment' | 'paid' | 'shipped' | 'delivered' | 'cancelled' | 'refunded'

// Fulfillment status (backend: "0"=pending, "1"=partial_shipped, "2"=shipped, "3"=delivered)
export type FulfillmentStatus = '0' | '1' | '2' | '3'

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

// Payment info (backend only has payment_method and payment_method_text)
export interface PaymentInfo {
  payment_method: string
  payment_method_text: string
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
  order_id: number
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
  currency: string
  item_count: number
  payment_method: string
  payment_method_text: string
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
export const getOrderDetail = (orderId: number) => {
  return request.get<OrderDetail>(`/api/v1/orders/${orderId}`)
}

/**
 * Ship order (create shipment from order)
 */
export const shipOrder = (orderId: number, data: ShipOrderRequest) => {
  return request.put<ShipOrderResponse>(`/api/v1/orders/${orderId}/ship`, data)
}

/**
 * Update order remark
 */
export const updateOrderRemark = (orderId: number, data: UpdateOrderRemarkRequest) => {
  return request.put<{ order_id: number }>(`/api/v1/orders/${orderId}/remark`, data)
}

/**
 * Adjust order price
 */
export const adjustOrderPrice = (orderId: number, data: AdjustPriceRequest) => {
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
 * Cancel order
 */
export const cancelOrder = (orderId: number, reason: string) => {
  return request.put<void>(`/api/v1/orders/${orderId}/cancel`, { reason })
}

/**
 * Remind payment
 */
export const remindPayment = (orderId: number) => {
  return request.post<void>(`/api/v1/orders/${orderId}/remind-payment`)
}

/**
 * Batch ship orders - delegates to fulfillment module
 * Note: Use batchCreateShipments from fulfillment.ts for the actual API call
 */
export const batchShipOrders = (data: {
  order_ids: number[]
  carrier_code: string
  tracking_no_start: string
}) => {
  return request.post<{ shipments: { order_id: number; shipment_id: number }[] }>(
    '/api/v1/shipments/batch',
    {
      order_ids: data.order_ids,
      carrier_code: data.carrier_code,
      tracking_no_start: data.tracking_no_start
    }
  )
}

// ===================== Batch Operations =====================

// Batch cancel order request
export interface BatchCancelOrderRequest {
  order_ids: number[]
  reason: string
}

// Batch cancel order response
export interface BatchCancelOrderResponse {
  success: number[]
  failed: { order_id: number; code: number; message: string }[]
}

// Batch cancel orders
export const batchCancelOrders = (data: BatchCancelOrderRequest) => {
  return request.post<BatchCancelOrderResponse>('/api/v1/orders/batch-cancel', data)
}

// Re-export Carrier type and getCarrierList from fulfillment module to avoid duplication
export type { Carrier } from './fulfillment'
export { getCarrierList } from './fulfillment'