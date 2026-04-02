import request from '@/utils/request'

// Warehouse types
export interface Warehouse {
  id: number
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
  id: number
  name: string
  country?: string
  address?: string
  is_default?: boolean
}

// Inventory types
export interface WarehouseInventoryItem {
  warehouse_id: number
  warehouse_name: string
  available_stock: number
  locked_stock: number
}

export interface SKUInventory {
  sku_code: string
  product_id: number
  total_stock: number
  available_stock: number
  locked_stock: number
  safety_stock: number
  is_low_stock: boolean
  warehouses: WarehouseInventoryItem[]
}

export interface InventoryLog {
  id: number
  sku_code: string
  product_id: number
  warehouse_id: number
  change_type: string
  change_quantity: number
  before_stock: number
  after_stock: number
  order_no: string
  remark: string
  operator_id: number
  created_at: string
}

export interface LowStockSKU {
  sku_code: string
  product_id: number
  product_name: string
  available_stock: number
  safety_stock: number
}

// Warehouse API functions
export function getWarehouses() {
  return request<Warehouse[]>({
    url: '/api/v1/warehouses',
    method: 'get'
  })
}

export function getWarehouse(id: number) {
  return request<Warehouse>({
    url: `/api/v1/warehouses/${id}`,
    method: 'get'
  })
}

export function createWarehouse(data: CreateWarehouseRequest) {
  return request<{ id: number }>({
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

export function updateWarehouseStatus(id: number, status: number) {
  return request<Warehouse>({
    url: `/api/v1/warehouses/${id}/status`,
    method: 'put',
    data: { id, status }
  })
}

export function setDefaultWarehouse(id: number) {
  return request<Warehouse>({
    url: `/api/v1/warehouses/${id}/default`,
    method: 'put'
  })
}

export function deleteWarehouse(id: number) {
  return request<{ id: number }>({
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
  warehouse_id?: number
  available_stock: number
  remark?: string
}) {
  return request<{ id: number }>({
    url: '/api/v1/inventory/stock',
    method: 'put',
    data
  })
}

export function adjustStock(data: {
  sku_code: string
  warehouse_id: number
  quantity: number
  remark?: string
}) {
  return request<{ id: number }>({
    url: '/api/v1/inventory/adjust',
    method: 'post',
    data
  })
}

export function getInventoryLogs(params?: {
  page?: number
  page_size?: number
  sku_code?: string
  product_id?: number
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
  product_id?: number
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
  return request<{ id: number }>({
    url: '/api/v1/inventory/safety-stock',
    method: 'put',
    data: { items }
  })
}