import request from '@/utils/request'

// ==================== Types ====================

// Tier Configuration
export interface TierConfig {
  threshold: number | null  // null for last tier (unlimited)
  ratio: string             // Points per currency unit
}

// Earn Rule
export interface EarnRule {
  id: number
  name: string
  description: string
  scenario: 'ORDER_PAYMENT' | 'SIGN_IN' | 'PRODUCT_REVIEW' | 'FIRST_ORDER'
  calculation_type: 'FIXED' | 'RATIO' | 'TIERED'
  fixed_points: number
  ratio: string
  tiers: TierConfig[] | null
  condition_type: 'NONE' | 'NEW_USER' | 'FIRST_ORDER' | 'SPECIFIC_PRODUCTS' | 'MIN_AMOUNT'
  condition_value: Record<string, any> | null
  expiration_months: number
  status: 'draft' | 'active' | 'inactive'
  priority: number
  start_at: string | null
  end_at: string | null
  created_at: string
  updated_at: string
}

// Redeem Rule
export interface RedeemRule {
  id: number
  name: string
  description: string
  coupon_id: number
  coupon_name: string
  points_required: number
  total_stock: number
  used_stock: number
  per_user_limit: number
  status: 'inactive' | 'active'
  start_at: string | null
  end_at: string | null
  created_at: string
  updated_at: string
}

// Points Account
export interface PointsAccount {
  id: number
  user_id: number
  user_email?: string
  balance: number
  frozen_balance: number
  total_earned: number
  total_redeemed: number
  total_expired: number
  created_at: string
  updated_at: string
}

// Points Transaction
export interface PointsTransaction {
  id: number
  user_id: number
  account_id: number
  points: number              // Positive = earn, negative = deduct
  balance_after: number
  type: 'EARN' | 'REDEEM' | 'ADJUST' | 'EXPIRE' | 'FREEZE' | 'UNFREEZE'
  reference_type: string
  reference_id: string
  description: string
  expires_at: string | null
  created_at: string
}

// Points Redemption
export interface PointsRedemption {
  id: number
  user_id: number
  redeem_rule_id: number
  coupon_id: number
  coupon_name: string
  user_coupon_id: number | null
  points_used: number
  status: 'pending' | 'completed' | 'cancelled'
  created_at: string
  completed_at: string | null
}

// Points Statistics
export interface PointsStats {
  total_issued: number
  total_redeemed: number
  total_expired: number
  outstanding_balance: number
  redemption_rate: string
  active_users: number
  period_start: string
  period_end: string
}

// Trend Data Point
export interface TrendDataPoint {
  date: string
  earned: number
  redeemed: number
  expired: number
}

// Top User
export interface TopUser {
  user_id: number
  user_email?: string
  points_earned: number
  created_at: string
}

// Expiring Points
export interface ExpiringPoints {
  date: string
  points: number
  user_count: number
}

// ==================== Request/Response Types ====================

export interface ListEarnRulesParams {
  page?: number
  page_size?: number
  name?: string
  status?: string
  scenario?: string
  calculation_type?: string
}

export interface CreateEarnRuleParams {
  name: string
  description?: string
  scenario: string
  calculation_type: string
  fixed_points?: number
  ratio?: string
  tiers?: TierConfig[]
  condition_type?: string
  condition_value?: Record<string, any>
  expiration_months?: number
  status?: string  // 'draft' | 'active' | 'inactive'
  priority?: number
  start_at?: string
  end_at?: string
}

export interface UpdateEarnRuleParams extends CreateEarnRuleParams {
  id: number
}

export interface ListRedeemRulesParams {
  page?: number
  page_size?: number
  name?: string
  status?: string
}

export interface CreateRedeemRuleParams {
  name: string
  description?: string
  coupon_id: number
  points_required: number
  total_stock?: number
  per_user_limit?: number
  status?: string  // 'inactive' | 'active'
  start_at?: string
  end_at?: string
}

export interface UpdateRedeemRuleParams extends CreateRedeemRuleParams {
  id: number
}

export interface ListAccountsParams {
  page?: number
  page_size?: number
  user_id?: number
  email?: string
}

export interface ListTransactionsParams {
  page?: number
  page_size?: number
  user_id?: number
  account_id?: number
  type?: string
  start_time?: string
  end_time?: string
}

export interface ListRedemptionsParams {
  page?: number
  page_size?: number
  user_id?: number
  status?: string
  start_time?: string
  end_time?: string
}

export interface AdjustPointsParams {
  adjustment_type: 'ADD' | 'DEDUCT'
  points: number
  reason: string
}

export interface AdjustPointsResponse {
  transaction_id: number
  points: number
  balance_after: number
  created_at: string
}

export interface PaginatedResponse<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

export interface EarnRulesListResponse extends PaginatedResponse<EarnRule> {
  stats: {
    total: number
    active: number
  }
}

export interface RedeemRulesListResponse extends PaginatedResponse<RedeemRule> {
  stats: {
    total: number
    active: number
    total_redeemed: number
  }
}

export interface AccountsListResponse extends PaginatedResponse<PointsAccount> {
  stats: {
    total: number
    total_balance: number
    active: number
  }
}

export interface TransactionsListResponse extends PaginatedResponse<PointsTransaction> {
  stats: {
    total_earned: number
    total_redeemed: number
  }
}

export interface TrendResponse {
  data: TrendDataPoint[]
  period: string
}

export interface TopUsersResponse {
  list: TopUser[]
  period: string
}

export interface ExpiringPointsResponse {
  list: ExpiringPoints[]
  total_points: number
}

// ==================== API Functions ====================

// ----- Statistics -----

/**
 * Get points statistics overview
 */
export const getPointsStats = (period: string = '7d') => {
  return request.get<PointsStats>('/api/v1/points/stats', { params: { period } })
}

/**
 * Get points trend data
 */
export const getPointsTrend = (period: string = '7d', granularity: 'daily' | 'weekly' = 'daily') => {
  return request.get<TrendResponse>('/api/v1/points/stats/trend', {
    params: { period, granularity }
  })
}

/**
 * Get top points earners
 */
export const getTopUsers = (period: string = '7d', limit: number = 10) => {
  return request.get<TopUsersResponse>('/api/v1/points/stats/top-users', {
    params: { period, limit }
  })
}

/**
 * Get expiring points preview
 */
export const getExpiringPoints = (days: number = 30) => {
  return request.get<ExpiringPointsResponse>('/api/v1/points/stats/expiring', {
    params: { days }
  })
}

// ----- Earn Rules -----

/**
 * Get earn rules list
 */
export const getEarnRules = (params: ListEarnRulesParams) => {
  return request.get<EarnRulesListResponse>('/api/v1/points/earn-rules', { params })
}

/**
 * Get earn rule detail
 */
export const getEarnRule = (id: number) => {
  return request.get<EarnRule>(`/api/v1/points/earn-rules/${id}`)
}

/**
 * Create earn rule
 */
export const createEarnRule = (data: CreateEarnRuleParams) => {
  return request.post<EarnRule>('/api/v1/points/earn-rules', data)
}

/**
 * Update earn rule
 */
export const updateEarnRule = (data: UpdateEarnRuleParams) => {
  const { id, ...params } = data
  return request.put<EarnRule>(`/api/v1/points/earn-rules/${id}`, params)
}

/**
 * Delete earn rule
 */
export const deleteEarnRule = (id: number) => {
  return request.delete(`/api/v1/points/earn-rules/${id}`)
}

/**
 * Activate earn rule
 */
export const activateEarnRule = (id: number) => {
  return request.post<EarnRule>(`/api/v1/points/earn-rules/${id}/activate`)
}

/**
 * Deactivate earn rule
 */
export const deactivateEarnRule = (id: number) => {
  return request.post<EarnRule>(`/api/v1/points/earn-rules/${id}/deactivate`)
}

// ----- Redeem Rules -----

/**
 * Get redeem rules list
 */
export const getRedeemRules = (params: ListRedeemRulesParams) => {
  return request.get<RedeemRulesListResponse>('/api/v1/points/redeem-rules', { params })
}

/**
 * Get redeem rule detail
 */
export const getRedeemRule = (id: number) => {
  return request.get<RedeemRule>(`/api/v1/points/redeem-rules/${id}`)
}

/**
 * Create redeem rule
 */
export const createRedeemRule = (data: CreateRedeemRuleParams) => {
  return request.post<RedeemRule>('/api/v1/points/redeem-rules', data)
}

/**
 * Update redeem rule
 */
export const updateRedeemRule = (data: UpdateRedeemRuleParams) => {
  const { id, ...params } = data
  return request.put<RedeemRule>(`/api/v1/points/redeem-rules/${id}`, params)
}

/**
 * Delete redeem rule
 */
export const deleteRedeemRule = (id: number) => {
  return request.delete(`/api/v1/points/redeem-rules/${id}`)
}

/**
 * Activate redeem rule
 */
export const activateRedeemRule = (id: number) => {
  return request.post<RedeemRule>(`/api/v1/points/redeem-rules/${id}/activate`)
}

/**
 * Deactivate redeem rule
 */
export const deactivateRedeemRule = (id: number) => {
  return request.post<RedeemRule>(`/api/v1/points/redeem-rules/${id}/deactivate`)
}

// ----- Points Accounts -----

/**
 * Get points accounts list
 */
export const getPointsAccounts = (params: ListAccountsParams) => {
  return request.get<AccountsListResponse>('/api/v1/points/accounts', { params })
}

/**
 * Get account detail
 */
export const getPointsAccount = (id: number) => {
  return request.get<PointsAccount>(`/api/v1/points/accounts/${id}`)
}

/**
 * Get account by user ID
 */
export const getAccountByUser = (userId: number) => {
  return request.get<PointsAccount>(`/api/v1/points/accounts/by-user/${userId}`)
}

/**
 * Get account transactions
 */
export const getAccountTransactions = (accountId: number, params: Omit<ListTransactionsParams, 'account_id'>) => {
  return request.get<TransactionsListResponse>(`/api/v1/points/accounts/${accountId}/transactions`, { params })
}

/**
 * Manual points adjustment
 */
export const adjustPoints = (accountId: number, data: AdjustPointsParams) => {
  return request.post<AdjustPointsResponse>(`/api/v1/points/accounts/${accountId}/adjust`, data)
}

// ----- Transactions -----

/**
 * Get all transactions
 */
export const getTransactions = (params: ListTransactionsParams) => {
  return request.get<TransactionsListResponse>('/api/v1/points/transactions', { params })
}

/**
 * Get transaction detail
 */
export const getTransaction = (id: number) => {
  return request.get<PointsTransaction>(`/api/v1/points/transactions/${id}`)
}

// ----- Redemptions -----

/**
 * Get redemptions list
 */
export const getRedemptions = (params: ListRedemptionsParams) => {
  return request.get<PaginatedResponse<PointsRedemption>>('/api/v1/points/redemptions', { params })
}

/**
 * Get redemption detail
 */
export const getRedemption = (id: number) => {
  return request.get<PointsRedemption>(`/api/v1/points/redemptions/${id}`)
}

// ----- Coupon Selection Helper -----

/**
 * Get available coupons for redeem rules
 * This reuses the promotion API to get active coupons
 */
export interface AvailableCoupon {
  id: number
  code: string
  name: string
  discount_value: string
  type?: 'fixed_amount' | 'percentage'
  start_time: string
  end_time: string
  status: string
}

export const getAvailableCoupons = (params: { page?: number; page_size?: number; name?: string }) => {
  return request.get<{ list: AvailableCoupon[]; total: number }>('/api/v1/coupons', { params })
}

// Export types
export interface ExportPointsTransactionsParams {
  user_id?: number
  type?: string
  start_time?: string
  end_time?: string
  [key: string]: unknown
}

/**
 * Export points transactions - returns URL and params for download utility
 */
export function exportPointsTransactionsUrl(params: ExportPointsTransactionsParams): { url: string; params: ExportPointsTransactionsParams } {
  return {
    url: '/api/v1/points/transactions/export',
    params
  }
}