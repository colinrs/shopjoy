import request from '@/utils/request'

export interface Category {
  id: string
  parent_id: string
  name: string
  code: string
  level: number
  sort: number
  icon: string
  image: string
  seo_title: string
  seo_description: string
  status: number
  product_count: number
  created_at: string
  updated_at: string
}

export interface CategoryTree extends Category {
  children: CategoryTree[]
}

export interface CreateCategoryRequest {
  name: string
  parent_id?: string
  code?: string
  icon?: string
  image?: string
  seo_title?: string
  seo_description?: string
  sort?: number
}

export interface UpdateCategoryRequest {
  id: string
  name: string
  code?: string
  icon?: string
  image?: string
  seo_title?: string
  seo_description?: string
  sort?: number
}

export interface ListCategoryRequest {
  parent_id?: string
}

export interface ListCategoryResponse {
  list: Category[]
}

export interface UpdateCategoryStatusRequest {
  id: string
  status: number // 0=disabled, 1=enabled
}

export interface UpdateCategorySortRequest {
  sorts: { id: string; sort: number }[]
}

export interface MoveCategoryRequest {
  id: string
  new_parent_id: string
}

// Get category tree
export function getCategoryTree() {
  return request<CategoryTree[]>({
    url: '/api/v1/categories/tree',
    method: 'get'
  })
}

// Get category list
export function getCategories(params?: ListCategoryRequest) {
  return request<ListCategoryResponse>({
    url: '/api/v1/categories',
    method: 'get',
    params
  })
}

// Get category detail
export function getCategory(id: string) {
  return request<Category>({
    url: `/api/v1/categories/${id}`,
    method: 'get'
  })
}

// Create category
export function createCategory(data: CreateCategoryRequest) {
  return request<{ id: string }>({
    url: '/api/v1/categories',
    method: 'post',
    data
  })
}

// Update category
export function updateCategory(data: UpdateCategoryRequest) {
  return request<Category>({
    url: `/api/v1/categories/${data.id}`,
    method: 'put',
    data
  })
}

// Delete category
export function deleteCategory(id: string) {
  return request<{ id: string }>({
    url: `/api/v1/categories/${id}`,
    method: 'delete'
  })
}

// Update category status
export function updateCategoryStatus(id: string, status: number) {
  return request<Category>({
    url: `/api/v1/categories/${id}/status`,
    method: 'put',
    data: { id, status }
  })
}

// Update category sort
export function updateCategorySort(sorts: { id: string; sort: number }[]) {
  return request<{ id: string }>({
    url: '/api/v1/categories/sort',
    method: 'put',
    data: { sorts }
  })
}

// Move category
export function moveCategory(id: string, new_parent_id: string) {
  return request<Category>({
    url: `/api/v1/categories/${id}/move`,
    method: 'put',
    data: { id, new_parent_id }
  })
}

// Get category product count
export function getCategoryProductCount(id: string) {
  return request<{ count: number }>({
    url: `/api/v1/categories/${id}/product-count`,
    method: 'get'
  })
}

// Category market visibility
export interface CategoryMarketVisibility {
  category_id: string
  markets: { market_id: string; is_visible: boolean }[]
}

export function getCategoryMarketVisibility(id: string) {
  return request<CategoryMarketVisibility>({
    url: `/api/v1/categories/${id}/market-visibility`,
    method: 'get'
  })
}

export function setCategoryMarketVisibility(id: string, market_ids: string[], visible: boolean) {
  return request<{ id: string }>({
    url: `/api/v1/categories/${id}/market-visibility`,
    method: 'put',
    data: { id, market_ids, visible }
  })
}