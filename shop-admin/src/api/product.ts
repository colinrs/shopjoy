import request from '@/utils/request'

// Product market info embedded in product detail
export interface ProductMarketInfo {
  market_id: string
  market_code: string
  market_name: string
  is_enabled: boolean
  price: string
  currency: string
}

export type ProductStatus = 'draft' | 'on_sale' | 'off_sale' | 'deleted'

// Product interface matching backend ProductDetailResp
export interface Product {
  id: string
  name: string
  description: string
  price: string
  currency: string
  cost_price: string
  stock: number
  status: ProductStatus
  category_id: string
  category_path: string
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
  id: string
  product_id: string
  market_id: string
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
  category_id?: string
  status?: string
  min_price?: string
  max_price?: string
  market_id?: string
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
  price: string
  currency?: string
  cost_price?: string
  category_id: string
  stock: number
  // New fields
  sku: string
  brand?: string
  tags?: string[]
  images?: string[]
  is_matrix_product?: boolean
  // Compliance fields
  hs_code?: string
  coo?: string
  weight?: number
  weight_unit?: string
  length?: number
  width?: number
  height?: number
  dangerous_goods?: string[]
}

export interface CreateProductResponse {
  id: string
}

export interface UpdateProductRequest {
  id: string
  name: string
  description?: string
  price: string
  currency?: string
  category_id: string
  // New fields
  sku?: string
  brand?: string
  tags?: string[]
  images?: string[]
  is_matrix_product?: boolean
  // Compliance fields
  hs_code?: string
  coo?: string
  weight?: number
  weight_unit?: string
  length?: number
  width?: number
  height?: number
  dangerous_goods?: string[]
}

export interface UpdateProductMarketRequest {
  is_enabled?: boolean
  price?: string
  compare_at_price?: string
  stock_alert_threshold?: number
}

export interface PushToMarketRequest {
  market_ids: string[]
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

export function getProduct(id: string) {
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

export function putOnSale(id: string) {
  return request<Product>({
    url: `/api/v1/products/${id}/on-sale`,
    method: 'post',
    data: {}
  })
}

export function takeOffSale(id: string) {
  return request<Product>({
    url: `/api/v1/products/${id}/off-sale`,
    method: 'post',
    data: {}
  })
}

export function updateStock(id: string, quantity: number) {
  return request<Product>({
    url: `/api/v1/products/${id}/stock`,
    method: 'put',
    data: { quantity }
  })
}

// Product Market API functions

export function getProductMarkets(productId: string) {
  return request<{ list: ProductMarket[] }>({
    url: `/api/v1/products/${productId}/markets`,
    method: 'get'
  })
}

export function updateProductMarket(productId: string, marketId: string, data: UpdateProductMarketRequest) {
  return request<ProductMarket>({
    url: `/api/v1/products/${productId}/markets/${marketId}`,
    method: 'put',
    data
  })
}

export function pushToMarket(productId: string, data: PushToMarketRequest) {
  return request<PushToMarketResponse>({
    url: `/api/v1/products/${productId}/push-to-market`,
    method: 'post',
    data
  })
}

export function removeFromMarket(productId: string, marketId: string) {
  return request<ProductMarket>({
    url: `/api/v1/products/${productId}/markets/${marketId}`,
    method: 'delete'
  })
}

// SKU (Variant) API functions

export interface SKU {
  id: string
  product_id: string
  code: string
  price: string
  currency: string
  stock: number
  available_stock: number
  locked_stock: number
  safety_stock: number
  pre_sale_enabled: boolean
  attributes: Record<string, string>
  status: string
  is_low_stock: boolean
  created_at: string
  updated_at: string
}

export interface CreateSKURequest {
  product_id: string
  code: string
  price: string
  currency?: string
  stock?: number
  safety_stock?: number
  pre_sale_enabled?: boolean
  attributes?: Record<string, string>
}

export interface UpdateSKURequest {
  id: string
  code?: string
  price?: string
  currency?: string
  stock?: number
  safety_stock?: number
  pre_sale_enabled?: boolean
  status?: string
  attributes?: Record<string, string>
}

export function getSKUsByProduct(productId: string) {
  return request<{ list: SKU[]; total: number }>({
    url: `/api/v1/products/${productId}/skus`,
    method: 'get'
  })
}

export function createSKU(data: CreateSKURequest) {
  return request<{ id: string }>({
    url: '/api/v1/skus',
    method: 'post',
    data
  })
}

export function updateSKU(data: UpdateSKURequest) {
  return request<SKU>({
    url: `/api/v1/skus/${data.id}`,
    method: 'put',
    data
  })
}

export function deleteSKU(id: string) {
  return request<{ id: string }>({
    url: `/api/v1/skus/${id}`,
    method: 'delete'
  })
}

export function getSKU(id: string) {
  return request<SKU>({
    url: `/api/v1/skus/${id}`,
    method: 'get'
  })
}

// Product Localization API functions

export interface ProductLocalization {
  id: string
  product_id: string
  language_code: string
  name: string
  description: string
  created_at: string
  updated_at: string
}

export interface CreateProductLocalizationRequest {
  product_id: string
  language_code: string
  name?: string
  description?: string
}

export interface UpdateProductLocalizationRequest {
  id: string
  name?: string
  description?: string
}

export function getProductLocalizations(productId: string) {
  return request<{ list: ProductLocalization[]; total: number }>({
    url: `/api/v1/products/${productId}/localizations`,
    method: 'get'
  })
}

export function createProductLocalization(data: CreateProductLocalizationRequest) {
  return request<{ id: string }>({
    url: '/api/v1/product-localizations',
    method: 'post',
    data
  })
}

export function updateProductLocalization(data: UpdateProductLocalizationRequest) {
  return request<ProductLocalization>({
    url: `/api/v1/product-localizations/${data.id}`,
    method: 'put',
    data
  })
}

export function deleteProductLocalization(id: string) {
  return request<{ id: string }>({
    url: `/api/v1/product-localizations/${id}`,
    method: 'delete'
  })
}

export function deleteProduct(id: string) {
  return request<{ id: string }>({
    url: `/api/v1/products/${id}`,
    method: 'delete'
  })
}

// Export types
export interface ExportProductsParams {
  name?: string
  category_id?: string
  status?: string
  min_price?: string
  max_price?: string
  market_id?: string
  [key: string]: unknown
}

/**
 * Export products - returns URL and params for download utility
 */
export function exportProductsUrl(params: ExportProductsParams): { url: string; params: ExportProductsParams } {
  return {
    url: '/api/v1/products/export',
    params
  }
}

// ===================== Batch Operations =====================

// Batch update product fields
export interface BatchUpdateProductFields {
  price?: string
  stock?: number
  status?: 'on_sale' | 'off_sale'
  category_id?: string
}

// Batch update product request
export interface BatchUpdateProductRequest {
  product_ids: string[]
  update_fields: BatchUpdateProductFields
}

// Batch update product response
export interface BatchUpdateProductResponse {
  success: number[]
  failed: { product_id: string; code: number; message: string }[]
}

// Batch update products
export function batchUpdateProducts(data: BatchUpdateProductRequest) {
  return request<BatchUpdateProductResponse>({
    url: '/api/v1/products/batch-update',
    method: 'post',
    data
  })
}