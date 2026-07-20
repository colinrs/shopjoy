import request from '@/utils/request'

// Types

// 模板：所有 ID 为 string（后端 int64 序列化为 string），金额为 string（元）
export interface ShippingTemplate {
  id: string
  tenant_id: string
  market_id: string        // 0 = 全市场
  currency: string         // ISO 4217
  carrier_code: string
  warehouse_id: string
  name: string
  is_default: boolean
  is_active: boolean
  zone_count: number
  created_at: string
}

// 区域名称多语言条目
export interface NameI18nEntry {
  locale: string           // 'en-US', 'ja-JP', 'zh-CN' ...
  name: string
}

// 运费区域：与后端 ShippingZone 1:1
export interface ShippingZone {
  id: string
  tenant_id: string
  template_id: string
  market_id: string
  currency: string
  name: string
  name_i18n?: NameI18nEntry[]
  regions: string[]
  fee_type: 'fixed' | 'by_count' | 'by_weight' | 'by_volume' | 'free'
  first_unit: number
  first_fee: string
  additional_unit: number
  additional_fee: string
  free_threshold_amount: string
  free_threshold_count: number
  taxable: boolean
  tax_rate: string
  tax_included: boolean
  ioss_applicable: boolean
  remote_surcharge: string
  remote_zip_patterns: string[]
  fuel_surcharge_pct: string
  volumetric_divisor: number
  sort: number
}

export interface TemplateMapping {
  id: string
  template_id: string
  target_type: 'product' | 'category'
  target_id: string
  target_name?: string
}

export interface TemplateDetail extends ShippingTemplate {
  zones: ShippingZone[]
  mappings: TemplateMapping[]
}

export interface TemplateListParams {
  page?: number
  page_size?: number
  name?: string
  is_active?: boolean
}

export interface CreateTemplateRequest {
  name: string
  is_default?: boolean
}

export interface UpdateTemplateRequest {
  name?: string
  is_active?: boolean
}

export interface CreateZoneRequest {
  name: string
  regions: string[]
  fee_type: 'fixed' | 'by_count' | 'by_weight' | 'by_volume' | 'free'
  first_unit?: number
  first_fee?: string
  additional_unit?: number
  additional_fee?: string
  free_threshold_amount?: string
  free_threshold_count?: number
  sort?: number
}

export interface UpdateZoneRequest extends Partial<CreateZoneRequest> {}

// ---- 运费试算（Calculator）：与后端 CalculateShippingFee 契约 1:1 ----

export interface CalculatorAddress {
  country_code: string         // ISO 3166-1 alpha-2 (REQUIRED)
  province_code?: string
  city_code: string
  district_code?: string
  postal_code?: string
}

export interface CalculatorItem {
  product_id: string
  sku_id?: string
  quantity: number
  weight: number              // 克
  length?: number             // mm
  width?: number              // mm
  height?: number             // mm
  price: string
}

export interface CalculateShippingFeeReq {
  market_id: string           // REQUIRED
  address: CalculatorAddress
  items: CalculatorItem[]
}

export interface FeeCalculationDetail {
  fee_type: string
  first_unit: number
  first_fee: string
  additional_unit: number
  additional_fee: string
  calculated_weight: number
  volumetric_weight?: number
  calculated_units: number
  applied_surcharge?: string
  applied_tax?: string
}

export interface CalculateShippingFeeResp {
  shipping_fee: string
  tax?: string
  total: string
  currency: string
  price_includes_tax: boolean
  template_id: string
  template_name: string
  zone_name: string
  carrier_code: string
  estimated_days?: number
  fee_detail: FeeCalculationDetail
}

export interface Region {
  code: string
  name: string
  level: number
  parent_code: string
  children?: Region[]
}

// API Functions

// Templates
export const getShippingTemplates = (params?: TemplateListParams) => {
  return request.get<{ list: ShippingTemplate[]; total: number }>('/api/v1/shipping-templates', { params })
}

export const getShippingTemplate = (id: string) => {
  return request.get<TemplateDetail>(`/api/v1/shipping-templates/${id}`)
}

export const createShippingTemplate = (data: CreateTemplateRequest) => {
  return request.post<{ id: string; name: string }>('/api/v1/shipping-templates', data)
}

export const updateShippingTemplate = (id: string, data: UpdateTemplateRequest) => {
  return request.put<ShippingTemplate>(`/api/v1/shipping-templates/${id}`, data)
}

export const deleteShippingTemplate = (id: string) => {
  return request.delete(`/api/v1/shipping-templates/${id}`)
}

export const setDefaultTemplate = (id: string) => {
  return request.put(`/api/v1/shipping-templates/${id}/set-default`)
}

// Zones
export const createShippingZone = (templateId: string, data: CreateZoneRequest) => {
  return request.post<ShippingZone>(`/api/v1/shipping-templates/${templateId}/zones`, data)
}

export const updateShippingZone = (id: string, data: UpdateZoneRequest) => {
  return request.put<ShippingZone>(`/api/v1/shipping-zones/${id}`, data)
}

export const deleteShippingZone = (id: string) => {
  return request.delete(`/api/v1/shipping-zones/${id}`)
}

export const reorderZones = (templateId: string, zoneIds: string[]) => {
  return request.put(`/api/v1/shipping-templates/${templateId}/zones/reorder`, { zone_ids: zoneIds })
}

// Mappings
export const getTemplateMappings = async (templateId: string) => {
  const res = await request.get<{ list: TemplateMapping[] }>(`/api/v1/shipping-templates/${templateId}/mappings`)
  return res.list || []
}

export const createTemplateMapping = (data: { template_id: string; target_type: 'product' | 'category'; target_id: string }) => {
  return request.post<TemplateMapping>('/api/v1/shipping-template-mappings', data)
}

export const deleteTemplateMapping = (id: string) => {
  return request.delete(`/api/v1/shipping-template-mappings/${id}`)
}

// Calculator
export const calculateShippingFee = (data: CalculateShippingFeeReq) => {
  return request.post<CalculateShippingFeeResp>('/api/v1/shipping/calculate', data)
}

// Regions
export const getRegions = async (parentCode?: string) => {
  const res = await request.get<{ list: Region[] }>('/api/v1/regions', { params: { parent_code: parentCode } })
  return res.list || []
}
