import request from '@/utils/request'

export function getProductList(params: {
  page: number
  page_size: number
  name?: string
  category_id?: number
  status?: string
  min_price?: number
  max_price?: number
}) {
  return request({
    url: '/api/v1/products',
    method: 'get',
    params
  })
}

export function createProduct(data: {
  name: string
  description?: string
  price: number
  currency?: string
  cost_price?: number
  category_id: number
}) {
  return request({
    url: '/api/v1/products',
    method: 'post',
    data
  })
}

export function getProduct(id: number) {
  return request({
    url: `/api/v1/products/${id}`,
    method: 'get'
  })
}

export function updateProduct(id: number, data: {
  name: string
  description?: string
  price: number
  currency?: string
  category_id: number
}) {
  return request({
    url: `/api/v1/products/${id}`,
    method: 'put',
    data
  })
}

export function putOnSale(id: number) {
  return request({
    url: `/api/v1/products/${id}/on-sale`,
    method: 'post',
    data: {}
  })
}

export function takeOffSale(id: number) {
  return request({
    url: `/api/v1/products/${id}/off-sale`,
    method: 'post',
    data: {}
  })
}

export function updateStock(id: number, quantity: number) {
  return request({
    url: `/api/v1/products/${id}/stock`,
    method: 'put',
    data: { quantity }
  })
}