import request from '@/utils/request'

export interface Brand {
  id: number
  name: string
  logo: string
  description: string
  website: string
  trademark_number: string
  trademark_country: string
  enable_page: number // 0=disabled, 1=enabled
  sort: number
  status: number
  product_count: number
  created_at: string
  updated_at: string
}

export interface ListBrandRequest {
  page?: number
  page_size?: number
  name?: string
  status?: number
}

export interface ListBrandResponse {
  list: Brand[]
  total: number
}

export interface CreateBrandRequest {
  name: string
  logo?: string
  description?: string
  website?: string
  trademark_number?: string
  trademark_country?: string
  enable_page?: boolean
  sort?: number
}

export interface UpdateBrandRequest {
  id: number
  name: string
  logo?: string
  description?: string
  website?: string
  trademark_number?: string
  trademark_country?: string
  enable_page?: boolean
  sort?: number
}

// Get brand list
export function getBrands(params?: ListBrandRequest) {
  return request<ListBrandResponse>({
    url: '/api/v1/brands',
    method: 'get',
    params
  })
}

// Get brand detail
export function getBrand(id: number) {
  return request<Brand>({
    url: `/api/v1/brands/${id}`,
    method: 'get'
  })
}

// Create brand
export function createBrand(data: CreateBrandRequest) {
  return request<{ id: number }>({
    url: '/api/v1/brands',
    method: 'post',
    data
  })
}

// Update brand
export function updateBrand(data: UpdateBrandRequest) {
  return request<Brand>({
    url: `/api/v1/brands/${data.id}`,
    method: 'put',
    data
  })
}

// Delete brand
export function deleteBrand(id: number) {
  return request<{ id: number }>({
    url: `/api/v1/brands/${id}`,
    method: 'delete'
  })
}

// Update brand status
export function updateBrandStatus(id: number, status: number) {
  return request<Brand>({
    url: `/api/v1/brands/${id}/status`,
    method: 'put',
    data: { id, status }
  })
}

// Toggle brand page
export function toggleBrandPage(id: number, enabled: boolean) {
  return request<Brand>({
    url: `/api/v1/brands/${id}/toggle-page`,
    method: 'put',
    data: { id, enabled }
  })
}

// Get brand product count
export function getBrandProductCount(id: number) {
  return request<{ count: number }>({
    url: `/api/v1/brands/${id}/product-count`,
    method: 'get'
  })
}

// Brand market visibility
export interface BrandMarketVisibility {
  brand_id: number
  markets: { market_id: number; is_visible: boolean }[]
}

export function getBrandMarketVisibility(id: number) {
  return request<BrandMarketVisibility>({
    url: `/api/v1/brands/${id}/market-visibility`,
    method: 'get'
  })
}

export function setBrandMarketVisibility(id: number, market_ids: number[], visible: boolean) {
  return request<{ id: number }>({
    url: `/api/v1/brands/${id}/market-visibility`,
    method: 'put',
    data: { id, market_ids, visible }
  })
}