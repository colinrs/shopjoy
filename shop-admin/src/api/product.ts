import request from '@/utils/request'

// Product market info embedded in product detail
export interface ProductMarketInfo {
  market_id: number
  market_code: string
  market_name: string
  is_enabled: boolean
  price: string
  currency: string
}

// Product interface matching backend ProductDetailResp
export interface Product {
  id: number
  name: string
  description: string
  price: number
  currency: string
  cost_price: number
  stock: number
  status: string
  category_id: number
  created_at: string
  updated_at: string
  // Compliance fields
  sku: string
  brand: string
  tags: string[]
  images: string[]
  is_matrix_product: boolean
  hs_code: string
  coo: string
  weight: string
  weight_unit: string
  length: string
  width: string
  height: string
  dangerous_goods: string[]
  // Markets
  markets: ProductMarketInfo[]
}

// Product market response from product market API
export interface ProductMarket {
  id: number
  product_id: number
  market_id: number
  market_code: string
  market_name: string
  is_enabled: boolean
  price: string
  compare_at_price: string
  currency: string
  stock_alert_threshold: number
  published_at: string
}

export interface ListProductsParams {
  page: number
  page_size: number
  name?: string
  category_id?: number
  status?: string
  min_price?: number
  max_price?: number
  market_id?: number
}

export interface ListProductsResponse {
  list: Product[]
  total: number
  page: number
  page_size: number
}

export interface CreateProductRequest {
  name: string
  description?: string
  price: number
  currency?: string
  cost_price?: number
  category_id: number
  // New fields
  sku?: string
  brand?: string
  tags?: string[]
  images?: string[]
  is_matrix_product?: boolean
  // Compliance fields
  hs_code?: string
  coo?: string
  weight?: string
  weight_unit?: string
  length?: string
  width?: string
  height?: string
  dangerous_goods?: string[]
}

export interface CreateProductResponse {
  id: number
}

export interface UpdateProductRequest {
  id: number
  name: string
  description?: string
  price: number
  currency?: string
  category_id: number
  // New fields
  sku?: string
  brand?: string
  tags?: string[]
  images?: string[]
  is_matrix_product?: boolean
  // Compliance fields
  hs_code?: string
  coo?: string
  weight?: string
  weight_unit?: string
  length?: string
  width?: string
  height?: string
  dangerous_goods?: string[]
}

export interface UpdateProductMarketRequest {
  is_enabled?: boolean
  price?: string
  compare_at_price?: string
  stock_alert_threshold?: number
}

export interface PushToMarketRequest {
  market_ids: number[]
  prices: string[]
}

export interface PushToMarketResponse {
  success: number[]
  failed: number[]
}

export function getProductList(params: ListProductsParams) {
  return request<ListProductsResponse>({
    url: '/api/v1/products',
    method: 'get',
    params
  })
}

export function createProduct(data: CreateProductRequest) {
  return request<CreateProductResponse>({
    url: '/api/v1/products',
    method: 'post',
    data
  })
}

export function getProduct(id: number) {
  return request<Product>({
    url: `/api/v1/products/${id}`,
    method: 'get'
  })
}

export function updateProduct(data: UpdateProductRequest) {
  return request<Product>({
    url: `/api/v1/products/${data.id}`,
    method: 'put',
    data
  })
}

export function putOnSale(id: number) {
  return request<Product>({
    url: `/api/v1/products/${id}/on-sale`,
    method: 'post',
    data: {}
  })
}

export function takeOffSale(id: number) {
  return request<Product>({
    url: `/api/v1/products/${id}/off-sale`,
    method: 'post',
    data: {}
  })
}

export function updateStock(id: number, quantity: number) {
  return request<Product>({
    url: `/api/v1/products/${id}/stock`,
    method: 'put',
    data: { quantity }
  })
}

// Product Market API functions

export function getProductMarkets(productId: number) {
  return request<{ list: ProductMarket[] }>({
    url: `/api/v1/products/${productId}/markets`,
    method: 'get'
  })
}

export function updateProductMarket(productId: number, marketId: number, data: UpdateProductMarketRequest) {
  return request<ProductMarket>({
    url: `/api/v1/products/${productId}/markets/${marketId}`,
    method: 'put',
    data
  })
}

export function pushToMarket(productId: number, data: PushToMarketRequest) {
  return request<PushToMarketResponse>({
    url: `/api/v1/products/${productId}/push-to-market`,
    method: 'post',
    data
  })
}

export function removeFromMarket(productId: number, marketId: number) {
  return request({
    url: `/api/v1/products/${productId}/markets/${marketId}`,
    method: 'delete'
  })
}