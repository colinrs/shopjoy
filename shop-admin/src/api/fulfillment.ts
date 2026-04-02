import request from '@/utils/request'
import { downloadFile } from '@/utils/download'

// Types
export type ShipmentStatus = '0' | '1' | '2' | '3' | '4' | '5'

// Status constants for use in templates
export const ShipmentStatusMap = {
  PENDING: '0',
  SHIPPED: '1',
  IN_TRANSIT: '2',
  DELIVERED: '3',
  FAILED: '4',
  CANCELLED: '5'
} as const

// Re-export FulfillmentStatus from order.ts to avoid duplication
// FulfillmentStatus: 0=pending, 1=partial_shipped, 2=shipped, 3=delivered
export type { FulfillmentStatus } from './order'

export type RefundType = 'full_refund' | 'partial_refund'

export type RefundStatus = '0' | '1' | '2' | '3' | '4'

export type OrderRefundStatus = '0' | '1' | '2' | '3' | '4'

export interface Shipment {
  id: number
  shipment_no: string
  order_id: number
  order_no: string
  status: ShipmentStatus
  carrier: string
  carrier_code: string
  tracking_no: string
  tracking_url: string
  shipping_cost: string
  currency: string
  weight: string
  shipped_at: string | null
  delivered_at: string | null
  remark: string
  created_at: string
  updated_at: string
  created_by: number
  created_by_name: string | null
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
  order_id: number
  carrier_code: string
  carrier?: string
  tracking_no: string
  shipping_cost?: string
  currency?: string
  weight?: string
  remark?: string
  items: {
    order_item_id: number
    quantity: number
  }[]
}

export interface UpdateShipmentRequest {
  carrier_code?: string
  carrier_name?: string
  tracking_no?: string
  shipping_cost?: string
  currency?: string
  weight?: string
  remark?: string
}

export interface ShipmentListParams {
  page?: number
  page_size?: number
  status?: ShipmentStatus
  carrier_code?: string
  tracking_no?: string
  start_time?: string
  end_time?: string
}

export interface Refund {
  id: number
  refund_no: string
  order_id: number
  order_no: string
  user_id: number
  user_name: string
  user_phone: string
  type: RefundType
  type_text: string
  status: RefundStatus
  status_text: string
  reason_type: string
  reason: string
  description: string
  images: string[]
  amount: string
  currency: string
  order_amount: string
  reject_reason: string
  approved_at: string | null
  approved_by: string | null
  completed_at: string | null
  created_at: string
  updated_at: string
  order_items: RefundOrderItem[]
}

export interface RefundOrderItem {
  id: number
  product_id: number
  product_name: string
  sku_name: string
  image: string
  quantity: number
  price: string
}

export interface RefundListParams {
  page?: number
  page_size?: number
  status?: RefundStatus
  reason_type?: string
  start_time?: string
  end_time?: string
}

export interface RejectRefundRequest {
  reject_reason: string
}

export interface Carrier {
  code: string
  name: string
  tracking_url: string
  is_active: boolean
}

export interface RefundReason {
  id: number
  code: string
  name: string
  sort: number
  is_active: boolean
}

export interface FulfillmentStatistics {
  overview: {
    total_shipments: number
    pending_shipments: number
    total_refunds: number
    pending_refunds: number
    refund_rate: string
    delivery_success_rate: string
    refund_amount: string
    currency: string
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

export interface BatchCreateShipmentItem {
  order_id: number
  tracking_no: string
}

export interface BatchCreateShipmentsRequest {
  carrier_code: string
  carrier_name?: string
  shipments: BatchCreateShipmentItem[]
}

export const batchCreateShipments = (data: BatchCreateShipmentsRequest) => {
  return request.post<Shipment[]>('/api/v1/shipments/batch', data)
}

export const updateShipmentStatus = (id: number, status: ShipmentStatus) => {
  return request.put(`/api/v1/shipments/${id}/status`, { status })
}

export const updateShipment = (id: number, data: UpdateShipmentRequest) => {
  return request.put<Shipment>(`/api/v1/shipments/${id}`, data)
}

export const getOrderShipments = (orderId: number) => {
  return request.get<Shipment[]>(`/api/v1/orders/${orderId}/shipments`)
}

// Refunds
export const getRefundList = (params: RefundListParams) => {
  return request.get<{ list: Refund[]; total: number }>('/api/v1/refunds', { params })
}

export const getRefundDetail = (id: number) => {
  return request.get<Refund>(`/api/v1/refunds/${id}`)
}

export const approveRefund = (refundId: number) => {
  return request.put(`/api/v1/refunds/${refundId}/approve`)
}

export const rejectRefund = (refundId: number, rejectReason: string) => {
  return request.put(`/api/v1/refunds/${refundId}/reject`, { reject_reason: rejectReason })
}

// Reference Data
export const getCarrierList = () => {
  return request.get<Carrier[]>('/api/v1/carriers')
}

export const getRefundReasonList = () => {
  return request.get<RefundReason[]>('/api/v1/refund-reasons')
}

// Statistics
export const getFulfillmentStatistics = (params: { period?: string; start_date?: string; end_date?: string }) => {
  // Convert frontend params to backend format
  const backendParams: Record<string, string> = {}
  if (params.period) {
    backendParams.period = params.period
  }
  if (params.start_date) {
    backendParams.start_time = params.start_date
  }
  if (params.end_date) {
    backendParams.end_time = params.end_date
  }
  return request.get<FulfillmentStatistics>('/api/v1/fulfillment/statistics', { params: backendParams })
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

// Export types
export interface ExportRefundsParams {
  order_no?: string
  refund_no?: string
  status?: string
  reason_type?: string
  start_time?: string
  end_time?: string
  [key: string]: unknown
}

export interface ExportShipmentsParams {
  tracking_no?: string
  status?: string
  carrier_code?: string
  start_time?: string
  end_time?: string
  [key: string]: unknown
}

export interface ExportFulfillmentStatisticsParams {
  period?: string
  start_date?: string
  end_date?: string
  [key: string]: unknown
}

/**
 * Export refunds - returns URL and params for download utility
 */
export function exportRefundsUrl(params: ExportRefundsParams): { url: string; params: ExportRefundsParams } {
  return {
    url: '/api/v1/refunds/export',
    params
  }
}

/**
 * Export shipments - returns URL and params for download utility
 */
export function exportShipmentsUrl(params: ExportShipmentsParams): { url: string; params: ExportShipmentsParams } {
  return {
    url: '/api/v1/shipments/export',
    params
  }
}

/**
 * Export fulfillment statistics - returns URL and params for download utility
 */
export function exportFulfillmentStatisticsUrl(params: ExportFulfillmentStatisticsParams): { url: string; params: ExportFulfillmentStatisticsParams } {
  return {
    url: '/api/v1/fulfillment/statistics/export',
    params
  }
}

/**
 * Export shipments - calls the API and downloads the file
 */
export async function exportShipments(params: ExportShipmentsParams): Promise<void> {
  const { url, params: queryParams } = exportShipmentsUrl(params)
  await downloadFile(url, queryParams)
}

/**
 * Export fulfillment statistics - calls the API and downloads the file
 */
export async function exportFulfillmentStatistics(params: ExportFulfillmentStatisticsParams): Promise<void> {
  const { url, params: queryParams } = exportFulfillmentStatisticsUrl(params)
  await downloadFile(url, queryParams)
}

// Batch Update Tracking
export interface BatchUpdateTrackingRequest {
  shipment_ids: number[]
  carrier_code: string
  tracking_no: string
  weight?: string
}

export interface BatchUpdateTrackingResponse {
  success: number[]
  failed: { shipment_id: number; code: number; message: string }[]
}

export const batchUpdateTracking = (data: BatchUpdateTrackingRequest) => {
  return request.post<BatchUpdateTrackingResponse>('/api/v1/shipments/batch-tracking', data)
}

// Cancel Shipment
export interface CancelShipmentResponse {
  id: number
  shipment_no: string
  status: string
  status_text: string
  cancelled_at: string
  reason: string
}

export const cancelShipment = (id: number, reason: string) => {
  return request.put<CancelShipmentResponse>(`/api/v1/shipments/${id}/cancel`, { reason })
}