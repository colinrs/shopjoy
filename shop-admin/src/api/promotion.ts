import request from '@/utils/request'

// ==================== Types ====================

// Unified Promotion/Coupon discriminator
export type PromotionKind = 'promotion' | 'coupon'

// Generic field enums (matching backend wire values)
export type PromotionType = 'discount' | 'flash_sale' | 'bundle' | 'buy_x_get_y'
export type PromotionStatus = 'pending' | 'active' | 'paused' | 'ended' | 'expired'
export type PromotionActionType = 'fixed_amount' | 'percentage' | 'free_shipping'
export type PromotionConditionType = 'min_amount' | 'min_quantity'
export type ScopeType = 'storewide' | 'products' | 'categories' | 'brands'

// User Coupon Types (still served by /user-coupons endpoints)
export type UserCouponStatus = 'unused' | 'used' | 'expired'

// ==================== Promotion Interfaces ====================

/**
 * PromotionRule — wire shape of promotion_rules rows.
 * Discount fields moved off of `promotions` onto the first matching rule.
 */
export interface PromotionRule {
  id: string
  condition_type: PromotionConditionType
  condition_value: string
  action_type: PromotionActionType
  action_value: string
  max_discount?: string
  sort_order?: number
}

/**
 * Unified Promotion entity. Discriminated by `kind`:
 * - `kind === 'promotion'` — activity/promotion (no code, no total_count)
 * - `kind === 'coupon'`    — coupon (code, total_count populated)
 *
 * Matches the backend `PromotionDetailResp` after the
 * promotion×coupon merge refactor (Tasks 1-8).
 */
export interface Promotion {
  id: string
  kind: PromotionKind
  name: string
  description: string
  code?: string                          // coupon only
  type: string
  status: PromotionStatus
  market_id?: string                     // optional
  currency: string
  usage_limit: number
  used_count: number
  per_user_limit: number
  total_count?: number                   // coupon only
  scope_type: ScopeType
  scope_ids: string[]
  exclude_ids: string[]
  tags: string[]
  rules?: PromotionRule[]                // nested rules
  start_time: string
  end_time: string
  created_at: string
  updated_at: string
}

export interface ListPromotionsParams {
  page: number
  page_size: number
  name?: string
  kind?: PromotionKind
  type?: string
  status?: PromotionStatus | string
  market_id?: string
}

export interface ListPromotionsResponse {
  list: Promotion[]
  total: number
  page: number
  page_size: number
}

export interface CreatePromotionRequest {
  kind: PromotionKind
  name: string
  description?: string
  code?: string
  type?: string
  market_id?: string
  currency?: string
  usage_limit?: number
  per_user_limit?: number
  total_count?: number
  scope_type?: ScopeType
  scope_ids?: string[]
  exclude_ids?: string[]
  tags?: string[]
  rules?: PromotionRuleRequest[]
  start_time: string
  end_time: string
}

export interface CreatePromotionResponse {
  id: string
}

export interface UpdatePromotionRequest {
  id: string
  kind: PromotionKind
  name: string
  description?: string
  code?: string
  type?: string
  market_id?: string
  currency?: string
  usage_limit?: number
  per_user_limit?: number
  total_count?: number
  scope_type?: ScopeType
  scope_ids?: string[]
  exclude_ids?: string[]
  tags?: string[]
  rules?: PromotionRuleRequest[]
  start_time: string
  end_time: string
}

export interface PromotionRuleRequest {
  condition_type: PromotionConditionType
  condition_value: string
  action_type: PromotionActionType
  action_value: string
  max_discount?: string
  sort_order?: number
}

// Rule CRUD — owner_kind tells the backend whether the parent is a
// promotion or a coupon (same table layout).
export interface ListPromotionRulesResponse {
  list: PromotionRule[]
  total: number
}

export interface CreatePromotionRulesRequest {
  owner_kind: PromotionKind
  rules: PromotionRuleRequest[]
}

export interface CreatePromotionRulesResponse {
  ids: string[]
}

export interface UpdatePromotionRuleRequest {
  id: string
  condition_type: PromotionConditionType
  condition_value: string
  action_type: PromotionActionType
  action_value: string
  max_discount?: string
  sort_order?: number
}

// ==================== Coupon Legacy Interfaces ====================
// NOTE: `CouponDetailResp` was merged into `Promotion` after the
// promotion×coupon merge refactor. The legacy `/api/v1/coupons`
// endpoints still exist for backwards compatibility and continue
// to return the unified `PromotionDetailResp` shape (with kind="coupon").
// Callers should use the unified `Promotion` type instead.
export type Coupon = Promotion

// ==================== Coupon Usage Interfaces ====================

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
  id: string
  coupon_id: string
  user_id: string
  order_id: string
  discount_amount: string
  used_at: string
}

export interface ListCouponUsageParams {
  id: string
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
  id: string
  user_id: string
  coupon_id: string
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
  order_id?: string
  created_at: string
}

export interface ListUserCouponsParams {
  page: number
  page_size: number
  user_id?: string
  coupon_id?: string
  status?: UserCouponStatus
}

export interface ListUserCouponsResponse {
  list: UserCoupon[]
  total: number
  page: number
  page_size: number
}

export interface IssueUserCouponRequest {
  user_id: string
  coupon_id: string
}

export interface IssueUserCouponResponse {
  id: string
}

export interface BatchIssueUserCouponRequest {
  coupon_id: string
  user_ids: string[]
}

export interface BatchIssueUserCouponResponse {
  issued: number
  user_coupon_ids: string[]
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

export function getPromotion(id: string) {
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

export function deletePromotion(id: string) {
  return request<{ id: string }>({
    url: `/api/v1/promotions/${id}`,
    method: 'delete'
  })
}

export function activatePromotion(id: string) {
  return request<Promotion>({
    url: `/api/v1/promotions/${id}/activate`,
    method: 'post',
    data: {}
  })
}

export function deactivatePromotion(id: string) {
  return request<Promotion>({
    url: `/api/v1/promotions/${id}/deactivate`,
    method: 'post',
    data: {}
  })
}

export function getPromotionRules(promotionId: string) {
  return request<ListPromotionRulesResponse>({
    url: `/api/v1/promotions/${promotionId}/rules`,
    method: 'get'
  })
}

export function createPromotionRules(
  promotionId: string,
  ownerKind: PromotionKind,
  rules: PromotionRuleRequest[]
) {
  return request<CreatePromotionRulesResponse>({
    url: `/api/v1/promotions/${promotionId}/rules`,
    method: 'post',
    data: { owner_kind: ownerKind, rules }
  })
}

export function updatePromotionRule(data: UpdatePromotionRuleRequest) {
  return request<PromotionRule>({
    url: `/api/v1/promotion-rules/${data.id}`,
    method: 'put',
    data
  })
}

export function deletePromotionRule(id: string) {
  return request<{ id: string }>({
    url: `/api/v1/promotion-rules/${id}`,
    method: 'delete'
  })
}

// ==================== Coupon API Functions (legacy endpoints) ====================

export function getCouponList(params: { page: number; page_size: number; name?: string; type?: string; status?: string }) {
  return request<ListPromotionsResponse>({
    url: '/api/v1/coupons',
    method: 'get',
    params
  })
}

export function createCoupon(data: CreatePromotionRequest) {
  return request<CreatePromotionResponse>({
    url: '/api/v1/coupons',
    method: 'post',
    data
  })
}

export function getCoupon(id: string) {
  return request<Promotion>({
    url: `/api/v1/coupons/${id}`,
    method: 'get'
  })
}

export function updateCoupon(data: UpdatePromotionRequest) {
  return request<Promotion>({
    url: `/api/v1/coupons/${data.id}`,
    method: 'put',
    data
  })
}

export function deleteCoupon(id: string) {
  return request<{ id: string }>({
    url: `/api/v1/coupons/${id}`,
    method: 'delete'
  })
}

export function activateCoupon(id: string) {
  return request<Promotion>({
    url: `/api/v1/coupons/${id}/activate`,
    method: 'post'
  })
}

export function deactivateCoupon(id: string) {
  return request<Promotion>({
    url: `/api/v1/coupons/${id}/deactivate`,
    method: 'post'
  })
}

export function generateCouponCodes(data: GenerateCouponCodesRequest) {
  return request<GenerateCouponCodesResponse>({
    url: '/api/v1/coupons/generate',
    method: 'post',
    data
  })
}

export function getCouponUsage(id: string, params: { page: number; page_size: number }) {
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

export function batchIssueUserCoupon(data: BatchIssueUserCouponRequest) {
  return request<BatchIssueUserCouponResponse>({
    url: '/api/v1/user-coupons/batch',
    method: 'post',
    data
  })
}
