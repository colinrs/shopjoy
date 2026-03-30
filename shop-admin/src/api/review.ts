import request from '@/utils/request'

// ==================== Types ====================

// Review Types
export type ReviewStatus = 'pending' | 'approved' | 'hidden'

// ==================== Review Interfaces ====================

export interface ReviewListItem {
  id: number
  order_id: number
  product_id: number
  product_name: string
  sku_code: string
  user_name: string
  is_anonymous: boolean
  is_verified: boolean
  quality_rating: number
  value_rating: number
  overall_rating: string
  content: string
  images: string[]
  status: ReviewStatus
  is_featured: boolean
  helpful_count: number
  has_reply: boolean
  created_at: string
}

export interface ListReviewsParams {
  product_id?: number
  status?: string
  rating_min?: number
  rating_max?: number
  has_image?: boolean
  keyword?: string
  start_time?: string
  end_time?: string
  page: number
  page_size: number
}

export interface ListReviewsResponse {
  list: ReviewListItem[]
  total: number
  page: number
  page_size: number
}

export interface ReviewDetail {
  id: number
  tenant_id: number
  order_id: number
  product_id: number
  product_name: string
  sku_code: string
  user_id: number
  user_name: string
  is_anonymous: boolean
  is_verified: boolean
  quality_rating: number
  value_rating: number
  overall_rating: string
  content: string
  images: string[]
  status: ReviewStatus
  is_featured: boolean
  helpful_count: number
  created_at: string
  updated_at: string
  reply?: ReviewReply
}

export interface ReviewReply {
  id: number
  content: string
  admin_name: string
  created_at: string
  updated_at: string
}

export interface ReviewStats {
  total_reviews: number
  pending_reviews: number
  approved_reviews: number
  hidden_reviews: number
  average_rating: string
  quality_avg_rating: string
  value_avg_rating: string
  five_star_count: number
  four_star_count: number
  three_star_count: number
  two_star_count: number
  one_star_count: number
  with_image_count: number
  reply_rate: number
  featured_count: number
}

export interface ProductStats {
  product_id: number
  total_reviews: number
  average_rating: string
  quality_avg_rating: string
  value_avg_rating: string
  rating_distribution: {
    "1": number
    "2": number
    "3": number
    "4": number
    "5": number
  }
  with_image_count: number
  reply_count: number
  reply_rate: number
}

export interface CreateReplyRequest {
  content: string
}

export interface UpdateReplyRequest {
  content: string
}

export interface HideReviewRequest {
  reason?: string
}

export interface ToggleFeaturedRequest {
  is_featured: boolean
}

export interface BatchApproveRequest {
  ids: number[]
}

export interface BatchHideRequest {
  ids: number[]
  reason?: string
}

export interface BatchOperationResponse {
  success_count: number
  failed_count: number
  errors: string[]
}

// ==================== Review API Functions ====================

export function getReviewList(params: ListReviewsParams) {
  return request<ListReviewsResponse>({
    url: '/api/v1/reviews',
    method: 'get',
    params
  })
}

export function getReviewDetail(id: number) {
  return request<ReviewDetail>({
    url: `/api/v1/reviews/${id}`,
    method: 'get'
  })
}

export function approveReview(id: number) {
  return request<{ id: number; status: string; updated_at: string }>({
    url: `/api/v1/reviews/${id}/approve`,
    method: 'put',
    data: {}
  })
}

export function hideReview(id: number, data: HideReviewRequest) {
  return request<{ id: number; status: string; updated_at: string }>({
    url: `/api/v1/reviews/${id}/hide`,
    method: 'put',
    data
  })
}

export function showReview(id: number) {
  return request<{ id: number; status: string; updated_at: string }>({
    url: `/api/v1/reviews/${id}/show`,
    method: 'put',
    data: {}
  })
}

export function deleteReview(id: number) {
  return request<{ id: number; status: string; deleted_at: string }>({
    url: `/api/v1/reviews/${id}`,
    method: 'delete'
  })
}

export function toggleFeatured(id: number, data: ToggleFeaturedRequest) {
  return request<{ id: number; is_featured: boolean; updated_at: string }>({
    url: `/api/v1/reviews/${id}/featured`,
    method: 'put',
    data
  })
}

export function createReply(id: number, data: CreateReplyRequest) {
  return request<ReviewReply>({
    url: `/api/v1/reviews/${id}/reply`,
    method: 'post',
    data
  })
}

export function updateReply(id: number, data: UpdateReplyRequest) {
  return request<ReviewReply>({
    url: `/api/v1/reviews/${id}/reply`,
    method: 'put',
    data
  })
}

export function deleteReply(id: number) {
  return request<{ success: boolean }>({
    url: `/api/v1/reviews/${id}/reply`,
    method: 'delete'
  })
}

export function batchApprove(data: BatchApproveRequest) {
  return request<BatchOperationResponse>({
    url: '/api/v1/reviews/batch-approve',
    method: 'post',
    data
  })
}

export function batchHide(data: BatchHideRequest) {
  return request<BatchOperationResponse>({
    url: '/api/v1/reviews/batch-hide',
    method: 'post',
    data
  })
}

export function getReviewStats() {
  return request<ReviewStats>({
    url: '/api/v1/reviews/stats',
    method: 'get'
  })
}

export function getProductStats(productId: number) {
  return request<ProductStats>({
    url: `/api/v1/reviews/product/${productId}/stats`,
    method: 'get'
  })
}