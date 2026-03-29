import request from '@/utils/request'

// Types
export interface Shipment {
  id: number
  shipment_no: string
  order_id: string
  order_no: string
  status: number // 0=pending, 1=shipped, 2=in_transit, 3=delivered, 4=failed
  carrier: string
  carrier_code: string
  tracking_no: string
  shipping_cost: string
  shipping_currency: string
  weight: string
  shipped_at: string | null
  delivered_at: string | null
  remark: string
  created_at: string
  items: ShipmentItem[]
}

export interface ShipmentItem {
  id: number
  product_id: number
  sku_id: number
  product_name: string
  sku_name: string
  image: string
  quantity: number
}

export interface CreateShipmentRequest {
  order_id: string
  carrier_code: string
  carrier?: string
  tracking_no: string
  shipping_cost?: string
  shipping_currency?: string
  weight?: string
  remark?: string
  items: {
    order_item_id: number
    quantity: number
  }[]
}

export interface ShipmentListParams {
  page?: number
  page_size?: number
  status?: number
  carrier_code?: string
  tracking_no?: string
  start_time?: string
  end_time?: string
}

export interface Refund {
  id: number
  refund_no: string
  order_id: string
  order_no: string
  user_id: number
  user_name: string
  user_phone: string
  type: number // 1=full_refund, 2=partial_refund
  status: number // 0=pending, 1=approved, 2=rejected, 3=completed, 4=cancelled
  reason_type: string
  reason: string
  description: string
  images: string[]
  amount: number
  currency: string
  reject_reason: string
  approved_at: string | null
  approved_by: string | null
  completed_at: string | null
  created_at: string
  order_items: RefundOrderItem[]
}

export interface RefundOrderItem {
  id: number
  product_id: number
  product_name: string
  sku_name: string
  image: string
  quantity: number
  price: number
}

export interface RefundListParams {
  page?: number
  page_size?: number
  status?: number
  reason_type?: string
  start_time?: string
  end_time?: string
}

export interface ApproveRefundRequest {
  refund_id: number
}

export interface RejectRefundRequest {
  refund_id: number
  reject_reason: string
}

export interface Carrier {
  code: string
  name: string
  tracking_url: string
  is_active: boolean
}

export interface RefundReason {
  code: string
  name: string
}

export interface FulfillmentStatistics {
  overview: {
    total_shipments: number
    pending_shipments: number
    total_refunds: number
    pending_refunds: number
    refund_rate: number
    delivery_success_rate: number
  }
  refund_rate_trend: {
    date: string
    rate: number
  }[]
  refund_reasons: {
    reason_type: string
    reason_name: string
    count: number
    percentage: number
  }[]
  problem_products: {
    product_id: number
    product_name: string
    image: string
    total_sales: number
    refund_count: number
    refund_rate: number
  }[]
  carrier_performance: {
    carrier_code: string
    carrier_name: string
    total_shipments: number
    delivery_success_rate: number
    avg_delivery_time: number
  }[]
}

// API Functions

// Shipments
export const getShipmentList = (params: ShipmentListParams) => {
  return request.get<{ list: Shipment[]; total: number }>('/api/v1/shipments', { params })
}

export const getShipmentDetail = (id: number) => {
  return request.get<Shipment>(`/api/v1/shipments/${id}`)
}

export const createShipment = (data: CreateShipmentRequest) => {
  return request.post<Shipment>('/api/v1/shipments', data)
}

export const batchCreateShipments = (data: {
  order_ids: string[]
  carrier_code: string
  tracking_no_start: string
}) => {
  return request.post<Shipment[]>('/api/v1/shipments/batch', data)
}

export const updateShipmentStatus = (id: number, status: number) => {
  return request.put(`/api/v1/shipments/${id}/status`, { status })
}

export const getOrderShipments = (orderId: string) => {
  return request.get<Shipment[]>(`/api/v1/orders/${orderId}/shipments`)
}

// Refunds
export const getRefundList = (params: RefundListParams) => {
  return request.get<{ list: Refund[]; total: number }>('/api/v1/refunds', { params })
}

export const getRefundDetail = (id: number) => {
  return request.get<Refund>(`/api/v1/refunds/${id}`)
}

export const approveRefund = (data: ApproveRefundRequest) => {
  return request.put(`/api/v1/refunds/${data.refund_id}/approve`)
}

export const rejectRefund = (data: RejectRefundRequest) => {
  return request.put(`/api/v1/refunds/${data.refund_id}/reject`, { reject_reason: data.reject_reason })
}

// Reference Data
export const getCarrierList = () => {
  return request.get<Carrier[]>('/api/v1/carriers')
}

export const getRefundReasonList = () => {
  return request.get<RefundReason[]>('/api/v1/refund-reasons')
}

// Statistics
export const getFulfillmentStatistics = (params: { start_date?: string; end_date?: string }) => {
  return request.get<FulfillmentStatistics>('/api/v1/refunds/statistics', { params })
}

// Fulfillment Summary
export const getFulfillmentSummary = () => {
  return request.get<{
    pending_shipment: number
    partial_shipped: number
    shipped: number
    delivered: number
    pending_refund: number
  }>('/api/v1/orders/fulfillment-summary')
}