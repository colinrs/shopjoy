import request from '@/utils/request'

// ==================== Types ====================

// Promotion Types
export type PromotionType = 'discount' | 'full_reduce' | 'flash_sale' | 'bundle'
export type PromotionStatus = 'draft' | 'active' | 'inactive' | 'expired'
export type DiscountType = 'percentage' | 'fixed_amount'
export type ConditionType = 'min_amount' | 'min_quantity'
export type ScopeType = 'storewide' | 'products' | 'categories' | 'brands'

// Coupon Types
export type CouponType = 'fixed_amount' | 'percentage'
export type CouponStatus = 'inactive' | 'active'

// User Coupon Types
export type UserCouponStatus = 'unused' | 'used' | 'expired'

// ==================== Promotion Interfaces ====================

export interface PromotionRule {
  id: number
  promotion_id: number
  rule_type: string
  operator: string
  value: string
  discount_type: string
  discount_value: string
  priority: number
  created_at: string
  updated_at: string
}

export interface Promotion {
  id: number
  name: string
  description: string
  type: PromotionType
  status: PromotionStatus
  start_time: string
  end_time: string
  discount_type: DiscountType
  discount_value: string
  min_order_amount: string
  max_discount: string
  usage_limit: number
  used_count: number
  per_user_limit: number
  product_ids: number[]
  category_ids: number[]
  market_ids: number[]
  tags: string[]
  created_at: string
  updated_at: string
}

export interface ListPromotionsParams {
  page: number
  page_size: number
  name?: string
  type?: PromotionType
  status?: PromotionStatus
  market_id?: number
}

export interface ListPromotionsResponse {
  list: Promotion[]
  total: number
  page: number
  page_size: number
}

export interface CreatePromotionRequest {
  name: string
  description?: string
  type: PromotionType
  start_time: string
  end_time: string
  discount_type?: DiscountType
  discount_value?: string
  min_order_amount?: string
  max_discount?: string
  usage_limit?: number
  per_user_limit?: number
  product_ids?: number[]
  category_ids?: number[]
  market_ids?: number[]
  tags?: string[]
}

export interface CreatePromotionResponse {
  id: number
}

export interface UpdatePromotionRequest {
  id: number
  name: string
  description?: string
  type: PromotionType
  start_time: string
  end_time: string
  discount_type?: DiscountType
  discount_value?: string
  min_order_amount?: string
  max_discount?: string
  usage_limit?: number
  per_user_limit?: number
  product_ids?: number[]
  category_ids?: number[]
  market_ids?: number[]
  tags?: string[]
}

export interface CreatePromotionRulesRequest {
  promotion_id: number
  rules: PromotionRuleRequest[]
}

export interface PromotionRuleRequest {
  rule_type: string
  operator: string
  value: string
  discount_type?: string
  discount_value?: string
  priority?: number
}

export interface CreatePromotionRulesResponse {
  ids: number[]
}

export interface ListPromotionRulesResponse {
  list: PromotionRule[]
  total: number
}

export interface UpdatePromotionRuleRequest {
  id: number
  rule_type: string
  operator: string
  value: string
  discount_type?: string
  discount_value?: string
  priority?: number
}

// ==================== Coupon Interfaces ====================

export interface Coupon {
  id: number
  code: string
  name: string
  description: string
  type: CouponType
  discount_value: string
  min_order_amount: string
  max_discount: string
  start_time: string
  end_time: string
  usage_limit: number
  used_count: number
  per_user_limit: number
  product_ids: string
  category_ids: string
  market_ids: string
  status: CouponStatus
  created_at: string
  updated_at: string
}

export interface ListCouponsParams {
  page: number
  page_size: number
  code?: string
  name?: string
  type?: CouponType
  status?: CouponStatus
}

export interface ListCouponsResponse {
  list: Coupon[]
  total: number
  page: number
  page_size: number
}

export interface CreateCouponRequest {
  code: string
  name: string
  description?: string
  type: CouponType
  discount_value: string
  min_order_amount?: string
  max_discount?: string
  start_time: string
  end_time: string
  usage_limit?: number
  per_user_limit?: number
  product_ids?: string
  category_ids?: string
  market_ids?: string
}

export interface CreateCouponResponse {
  id: number
}

export interface UpdateCouponRequest {
  id: number
  code: string
  name: string
  description?: string
  type: CouponType
  discount_value: string
  min_order_amount?: string
  max_discount?: string
  start_time: string
  end_time: string
  usage_limit?: number
  per_user_limit?: number
  product_ids?: string
  category_ids?: string
  market_ids?: string
}

export interface GenerateCouponCodesRequest {
  prefix?: string
  quantity: number
  length?: number
  coupon_config: string
}

export interface GenerateCouponCodesResponse {
  codes: string[]
  count: number
}

export interface CouponUsage {
  id: number
  coupon_id: number
  user_id: number
  order_id: number
  discount_amount: string
  used_at: string
}

export interface ListCouponUsageParams {
  id: number
  page: number
  page_size: number
}

export interface ListCouponUsageResponse {
  list: CouponUsage[]
  total: number
  page: number
  page_size: number
}

// ==================== User Coupon Interfaces ====================

export interface UserCoupon {
  id: number
  user_id: number
  coupon_id: number
  coupon_code: string
  coupon_name: string
  discount_type: string
  discount_value: string
  min_order_amount: string
  max_discount: string
  start_time: string
  end_time: string
  status: UserCouponStatus
  used_at?: string
  order_id?: number
  created_at: string
}

export interface ListUserCouponsParams {
  page: number
  page_size: number
  user_id?: number
  coupon_id?: number
  status?: UserCouponStatus
}

export interface ListUserCouponsResponse {
  list: UserCoupon[]
  total: number
  page: number
  page_size: number
}

export interface IssueUserCouponRequest {
  user_id: number
  coupon_id: number
}

export interface IssueUserCouponResponse {
  id: number
}

// ==================== Promotion API Functions ====================

export function getPromotionList(params: ListPromotionsParams) {
  return request<ListPromotionsResponse>({
    url: '/api/v1/promotions',
    method: 'get',
    params
  })
}

export function createPromotion(data: CreatePromotionRequest) {
  return request<CreatePromotionResponse>({
    url: '/api/v1/promotions',
    method: 'post',
    data
  })
}

export function getPromotion(id: number) {
  return request<Promotion>({
    url: `/api/v1/promotions/${id}`,
    method: 'get'
  })
}

export function updatePromotion(data: UpdatePromotionRequest) {
  return request<Promotion>({
    url: `/api/v1/promotions/${data.id}`,
    method: 'put',
    data
  })
}

export function deletePromotion(id: number) {
  return request<{ id: number }>({
    url: `/api/v1/promotions/${id}`,
    method: 'delete'
  })
}

export function activatePromotion(id: number) {
  return request<Promotion>({
    url: `/api/v1/promotions/${id}/activate`,
    method: 'post',
    data: {}
  })
}

export function deactivatePromotion(id: number) {
  return request<Promotion>({
    url: `/api/v1/promotions/${id}/deactivate`,
    method: 'post',
    data: {}
  })
}

export function getPromotionRules(promotionId: number) {
  return request<ListPromotionRulesResponse>({
    url: `/api/v1/promotions/${promotionId}/rules`,
    method: 'get'
  })
}

export function createPromotionRules(promotionId: number, data: PromotionRuleRequest[]) {
  return request<CreatePromotionRulesResponse>({
    url: `/api/v1/promotions/${promotionId}/rules`,
    method: 'post',
    data: { rules: data }
  })
}

export function updatePromotionRule(data: UpdatePromotionRuleRequest) {
  return request<PromotionRule>({
    url: `/api/v1/promotion-rules/${data.id}`,
    method: 'put',
    data
  })
}

export function deletePromotionRule(id: number) {
  return request<{ id: number }>({
    url: `/api/v1/promotion-rules/${id}`,
    method: 'delete'
  })
}

// ==================== Coupon API Functions ====================

export function getCouponList(params: ListCouponsParams) {
  return request<ListCouponsResponse>({
    url: '/api/v1/coupons',
    method: 'get',
    params
  })
}

export function createCoupon(data: CreateCouponRequest) {
  return request<CreateCouponResponse>({
    url: '/api/v1/coupons',
    method: 'post',
    data
  })
}

export function getCoupon(id: number) {
  return request<Coupon>({
    url: `/api/v1/coupons/${id}`,
    method: 'get'
  })
}

export function updateCoupon(data: UpdateCouponRequest) {
  return request<Coupon>({
    url: `/api/v1/coupons/${data.id}`,
    method: 'put',
    data
  })
}

export function deleteCoupon(id: number) {
  return request<{ id: number }>({
    url: `/api/v1/coupons/${id}`,
    method: 'delete'
  })
}

export function generateCouponCodes(data: GenerateCouponCodesRequest) {
  return request<GenerateCouponCodesResponse>({
    url: '/api/v1/coupons/generate',
    method: 'post',
    data
  })
}

export function getCouponUsage(id: number, params: { page: number; page_size: number }) {
  return request<ListCouponUsageResponse>({
    url: `/api/v1/coupons/${id}/usage`,
    method: 'get',
    params
  })
}

// ==================== User Coupon API Functions ====================

export function getUserCouponList(params: ListUserCouponsParams) {
  return request<ListUserCouponsResponse>({
    url: '/api/v1/user-coupons',
    method: 'get',
    params
  })
}

export function issueUserCoupon(data: IssueUserCouponRequest) {
  return request<IssueUserCouponResponse>({
    url: '/api/v1/user-coupons',
    method: 'post',
    data
  })
}