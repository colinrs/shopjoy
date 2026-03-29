import request from '@/utils/request'

// ===================== Types =====================

// Dashboard overview response
export interface DashboardOverview {
  today_orders: number
  today_sales: string
  today_growth: string
  yesterday_sales: string
  total_products: number
  total_users: number
  new_users_today: number
  currency: string
}

// Sales trend data point
export interface SalesTrendData {
  date: string
  sales: string
  orders: number
}

// Sales trend response
export interface SalesTrendResponse {
  period: string
  data: SalesTrendData[]
  currency: string
}

// Order status item
export interface OrderStatusItem {
  status: string
  status_text: string
  count: number
  percentage: string
  color: string
}

// Order status distribution response
export interface OrderStatusDistribution {
  list: OrderStatusItem[]
  total: number
}

// Top product item
export interface TopProductItem {
  product_id: number
  product_name: string
  image: string
  sales: number
  revenue: string
}

// Top products response
export interface TopProductsResponse {
  list: TopProductItem[]
  currency: string
}

// Pending order item
export interface PendingOrderItem {
  order_id: number
  order_no: string
  pay_amount: string
  status: string
  status_text: string
  created_at: string
  user_name?: string
}

// Pending orders response
export interface PendingOrdersResponse {
  list: PendingOrderItem[]
  total: number
}

// Activity item
export interface ActivityItem {
  id: number
  type: string // 'order_created' | 'payment_received' | 'product_low_stock' etc.
  content: string
  time: string
  operator?: string
}

// Recent activities response
export interface RecentActivitiesResponse {
  list: ActivityItem[]
}

// Combined dashboard response
export interface DashboardResponse {
  overview: DashboardOverview
  status_distribution: OrderStatusDistribution
  pending_orders: PendingOrderItem[]
  top_products: TopProductItem[]
  recent_activities: ActivityItem[]
}

// Query parameters
export interface SalesTrendParams {
  period?: 'week' | 'month' | 'year'
}

export interface TopProductsParams {
  limit?: number
  period?: 'week' | 'month' | 'all'
}

export interface PendingOrdersParams {
  limit?: number
}

export interface RecentActivitiesParams {
  limit?: number
}

// ===================== API Functions =====================

/**
 * Get combined dashboard data (aggregated endpoint)
 */
export const getDashboard = () => {
  return request.get<DashboardResponse>('/api/v1/dashboard')
}

/**
 * Get dashboard overview statistics
 */
export const getDashboardOverview = () => {
  return request.get<DashboardOverview>('/api/v1/dashboard/overview')
}

/**
 * Get sales trend data
 */
export const getSalesTrend = (params?: SalesTrendParams) => {
  return request.get<SalesTrendResponse>('/api/v1/dashboard/sales-trend', { params })
}

/**
 * Get order status distribution
 */
export const getOrderStatusDistribution = () => {
  return request.get<OrderStatusDistribution>('/api/v1/dashboard/order-status')
}

/**
 * Get top selling products
 */
export const getTopProducts = (params?: TopProductsParams) => {
  return request.get<TopProductsResponse>('/api/v1/dashboard/top-products', { params })
}

/**
 * Get pending payment orders
 */
export const getPendingOrders = (params?: PendingOrdersParams) => {
  return request.get<PendingOrdersResponse>('/api/v1/dashboard/pending-orders', { params })
}

/**
 * Get recent activities
 */
export const getRecentActivities = (params?: RecentActivitiesParams) => {
  return request.get<RecentActivitiesResponse>('/api/v1/dashboard/activities', { params })
}