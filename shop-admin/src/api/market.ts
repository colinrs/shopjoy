import request from '@/utils/request'

export interface Market {
  id: number
  code: string
  name: string
  currency: string
  default_language: string
  flag: string
  is_active: boolean
  is_default: boolean
  tax_rules: {
    vat_rate: string
    gst_rate: string
    ioss_enabled: boolean
    include_tax: boolean
  }
  created_at: string
  updated_at: string
}

export interface ListMarketsResponse {
  list: Market[]
  total: number
}

export interface CreateMarketRequest {
  code: string
  name: string
  currency: string
  default_language?: string
  flag?: string
  tax_rules?: {
    vat_rate?: string
    gst_rate?: string
    ioss_enabled?: boolean
    include_tax?: boolean
  }
}

export interface UpdateMarketRequest {
  id: number
  name?: string
  is_active?: boolean
  tax_rules?: {
    vat_rate?: string
    gst_rate?: string
    ioss_enabled?: boolean
    include_tax?: boolean
  }
}

export function getMarkets() {
  return request<ListMarketsResponse>({
    url: '/api/v1/markets',
    method: 'get'
  })
}

export function getMarket(id: number) {
  return request<Market>({
    url: `/api/v1/markets/${id}`,
    method: 'get'
  })
}

export function createMarket(data: CreateMarketRequest) {
  return request<Market>({
    url: '/api/v1/markets',
    method: 'post',
    data
  })
}

export function updateMarket(data: UpdateMarketRequest) {
  return request<Market>({
    url: `/api/v1/markets/${data.id}`,
    method: 'put',
    data
  })
}

export function deleteMarket(id: number) {
  return request({
    url: `/api/v1/markets/${id}`,
    method: 'delete'
  })
}