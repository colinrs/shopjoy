import request from '@/utils/request'

// ===================== Shop Settings Types =====================

export interface ShopSettings {
  id: number
  name: string
  code: string
  logo: string
  description: string
  contact_name: string
  contact_phone: string
  contact_email: string
  address: string
  domain: string
  custom_domain: string
  primary_color: string
  secondary_color: string
  favicon: string
  default_currency: string
  default_language: string
  timezone: string
  status: number
  status_text: string
  plan: number
  plan_text: string
  expire_at: string
  created_at: string
  updated_at: string
}

export interface UpdateShopSettingsParams {
  name?: string
  logo?: string
  description?: string
  contact_name?: string
  contact_phone?: string
  contact_email?: string
  address?: string
  custom_domain?: string
  primary_color?: string
  secondary_color?: string
  favicon?: string
  default_language?: string
  timezone?: string
}

// ===================== Business Hours Types =====================

export interface BusinessHours {
  day_of_week: number // 0=Sunday, 1=Monday...
  open_time: string // "09:00"
  close_time: string // "18:00"
  is_closed: boolean
}

export interface UpdateBusinessHoursParams {
  hours: BusinessHours[]
}

// ===================== Notification Settings Types =====================

export interface NotificationSettings {
  order_created: boolean
  order_paid: boolean
  order_shipped: boolean
  order_cancelled: boolean
  low_stock_alert: boolean
  low_stock_threshold: number
  refund_requested: boolean
  new_review: boolean
}

export interface UpdateNotificationSettingsParams extends NotificationSettings {}

// ===================== Payment Settings Types =====================

export interface PaymentSettings {
  stripe_enabled: boolean
  stripe_public_key: string
}

export interface UpdatePaymentSettingsParams {
  stripe_enabled: boolean
  stripe_secret_key?: string
}

// ===================== Shipping Settings Types =====================

export interface ShippingSettings {
  free_shipping_threshold: string
  default_shipping_fee: string
  currency: string
}

export interface UpdateShippingSettingsParams {
  free_shipping_threshold?: string
  default_shipping_fee?: string
}

// ===================== API Functions =====================

/**
 * Get shop settings
 */
export function getShopSettings(): Promise<ShopSettings> {
  return request({
    url: '/api/v1/shop',
    method: 'get'
  })
}

/**
 * Update shop settings
 */
export function updateShopSettings(data: UpdateShopSettingsParams): Promise<ShopSettings> {
  return request({
    url: '/api/v1/shop',
    method: 'put',
    data
  })
}

/**
 * Get business hours
 */
export function getBusinessHours(): Promise<BusinessHours[]> {
  return request({
    url: '/api/v1/shop/business-hours',
    method: 'get'
  })
}

/**
 * Update business hours
 */
export function updateBusinessHours(data: UpdateBusinessHoursParams): Promise<void> {
  return request({
    url: '/api/v1/shop/business-hours',
    method: 'put',
    data
  })
}

/**
 * Get notification settings
 */
export function getNotificationSettings(): Promise<NotificationSettings> {
  return request({
    url: '/api/v1/shop/notifications',
    method: 'get'
  })
}

/**
 * Update notification settings
 */
export function updateNotificationSettings(data: UpdateNotificationSettingsParams): Promise<NotificationSettings> {
  return request({
    url: '/api/v1/shop/notifications',
    method: 'put',
    data
  })
}

/**
 * Get payment settings
 */
export function getPaymentSettings(): Promise<PaymentSettings> {
  return request({
    url: '/api/v1/shop/payment',
    method: 'get'
  })
}

/**
 * Update payment settings
 */
export function updatePaymentSettings(data: UpdatePaymentSettingsParams): Promise<PaymentSettings> {
  return request({
    url: '/api/v1/shop/payment',
    method: 'put',
    data
  })
}

/**
 * Get shipping settings
 */
export function getShippingSettings(): Promise<ShippingSettings> {
  return request({
    url: '/api/v1/shop/shipping',
    method: 'get'
  })
}

/**
 * Update shipping settings
 */
export function updateShippingSettings(data: UpdateShippingSettingsParams): Promise<ShippingSettings> {
  return request({
    url: '/api/v1/shop/shipping',
    method: 'put',
    data
  })
}