import request from '@/utils/request'

// Types
export interface ShippingTemplate {
  id: string
  name: string
  is_default: boolean
  is_active: boolean
  zone_count: number
  product_count: number
  category_count: number
  created_at: string
}

export interface ShippingZone {
  id: string
  template_id: string
  name: string
  regions: string[]
  fee_type: 'fixed' | 'by_count' | 'by_weight' | 'free'
  first_unit: number
  first_fee: string
  additional_unit: number
  additional_fee: string
  free_threshold_amount: string
  free_threshold_count: number
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
  fee_type: 'fixed' | 'by_count' | 'by_weight' | 'free'
  first_unit?: number
  first_fee?: string
  additional_unit?: number
  additional_fee?: string
  free_threshold_amount?: string
  free_threshold_count?: number
  sort?: number
}

export interface UpdateZoneRequest extends Partial<CreateZoneRequest> {}

export interface CalculateRequest {
  address: {
    province_code: string
    city_code: string
    district_code: string
  }
  items: Array<{
    product_id: string
    sku_id?: string
    quantity: number
    weight: number
    price: string
  }>
}

export interface FeeDetail {
  fee_type: string
  first_unit: number
  first_fee: string
  additional_unit: number
  additional_fee: string
  calculated_weight: number
  calculated_units: number
}

export interface CalculateResult {
  shipping_fee: string
  currency: string
  template_id: string
  template_name: string
  zone_name: string
  fee_detail: FeeDetail
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
export const calculateShippingFee = (data: CalculateRequest) => {
  return request.post<CalculateResult>('/api/v1/shipping/calculate', data)
}

// Regions
export const getRegions = async (parentCode?: string) => {
  const res = await request.get<{ list: Region[] }>('/api/v1/regions', { params: { parent_code: parentCode } })
  return res.list || []
}