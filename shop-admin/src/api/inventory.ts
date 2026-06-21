import request from '@/utils/request'

// Warehouse types
export interface Warehouse {
  id: string
  code: string
  name: string
  country: string
  address: string
  is_default: boolean
  status: number
  created_at: string
  updated_at: string
}

export interface CreateWarehouseRequest {
  code: string
  name: string
  country?: string
  address?: string
  is_default?: boolean
}

export interface UpdateWarehouseRequest {
  id: string
  name: string
  country?: string
  address?: string
  is_default?: boolean
}

// Inventory types
export interface WarehouseInventoryItem {
  warehouse_id: string
  warehouse_name: string
  available_stock: number
  locked_stock: number
}

export interface SKUInventory {
  sku_code: string
  product_id: string
  total_stock: number
  available_stock: number
  locked_stock: number
  safety_stock: number
  is_low_stock: boolean
  warehouses: WarehouseInventoryItem[]
}

export interface InventoryLog {
  id: string
  sku_code: string
  product_id: string
  warehouse_id: string
  change_type: string
  change_quantity: number
  before_stock: number
  after_stock: number
  order_no: string
  remark: string
  operator_id: string
  created_at: string
}

export interface LowStockSKU {
  sku_code: string
  product_id: string
  product_name: string
  available_stock: number
  safety_stock: number
}

// Warehouse API functions
export function getWarehouses() {
  return request<{ list: Warehouse[] }>({
    url: '/api/v1/warehouses',
    method: 'get'
  })
}

export function getWarehouse(id: string) {
  return request<Warehouse>({
    url: `/api/v1/warehouses/${id}`,
    method: 'get'
  })
}

export function createWarehouse(data: CreateWarehouseRequest) {
  return request<{ id: string }>({
    url: '/api/v1/warehouses',
    method: 'post',
    data
  })
}

export function updateWarehouse(data: UpdateWarehouseRequest) {
  return request<Warehouse>({
    url: `/api/v1/warehouses/${data.id}`,
    method: 'put',
    data
  })
}

export function updateWarehouseStatus(id: string, status: number) {
  return request<Warehouse>({
    url: `/api/v1/warehouses/${id}/status`,
    method: 'put',
    data: { id, status }
  })
}

export function setDefaultWarehouse(id: string) {
  return request<Warehouse>({
    url: `/api/v1/warehouses/${id}/default`,
    method: 'put'
  })
}

export function deleteWarehouse(id: string) {
  return request<{ id: string }>({
    url: `/api/v1/warehouses/${id}`,
    method: 'delete'
  })
}

// Inventory API functions
export function getSKUInventory(skuCode: string) {
  return request<SKUInventory>({
    url: `/api/v1/inventory/sku/${skuCode}`,
    method: 'get'
  })
}

export function updateSKUStock(data: {
  sku_code: string
  warehouse_id?: string
  available_stock: number
  remark?: string
}) {
  return request<{ id: string }>({
    url: '/api/v1/inventory/stock',
    method: 'put',
    data
  })
}

export function adjustStock(data: {
  sku_code: string
  warehouse_id: string
  quantity: number
  remark?: string
}) {
  return request<{ id: string }>({
    url: '/api/v1/inventory/adjust',
    method: 'post',
    data
  })
}

export function getInventoryLogs(params?: {
  page?: number
  page_size?: number
  sku_code?: string
  product_id?: string
  type?: string
}) {
  return request<{ list: InventoryLog[]; total: number }>({
    url: '/api/v1/inventory/logs',
    method: 'get',
    params
  })
}

export interface ExportInventoryLogsParams {
  sku_code?: string
  product_id?: string
  type?: string
  [key: string]: unknown
}

/**
 * Export inventory logs - returns URL and params for download utility
 */
export function exportInventoryLogsUrl(params: ExportInventoryLogsParams): { url: string; params: ExportInventoryLogsParams } {
  return {
    url: '/api/v1/inventory/logs/export',
    params
  }
}

export function getLowStockSKUs(params?: { page?: number; page_size?: number }) {
  return request<{ list: LowStockSKU[]; total: number }>({
    url: '/api/v1/inventory/low-stock',
    method: 'get',
    params
  })
}

export function batchUpdateSafetyStock(items: { sku_code: string; safety_stock: number }[]) {
  return request<{ id: string }>({
    url: '/api/v1/inventory/safety-stock',
    method: 'put',
    data: { items }
  })
}