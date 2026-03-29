import request from '@/utils/request'

// User interface matching backend GetUserResponse
export interface User {
  id: number
  tenant_id: number
  email: string
  phone: string
  name: string
  avatar: string
  gender: number
  gender_text: string
  birthday: string | null
  status: number
  status_text: string
  review_count: number
  created_at: string
  updated_at: string
  last_login: string
}

// Extended user with points and order stats
export interface ExtendedUser extends User {
  points_balance: number
  order_count: number
  total_spent: string
}

// User detail with full information
export interface UserDetail extends ExtendedUser {
  frozen_points: number
  total_earned_points: number
  total_redeemed_points: number
  last_order_at: string
  default_address: UserAddress | null
}

// User address
export interface UserAddress {
  id: number
  user_id: number
  name: string
  phone: string
  country: string
  province: string
  city: string
  district: string
  detail: string
  postal_code: string
  is_default: boolean
  created_at: string
  updated_at: string
}

export interface UserStats {
  total: number
  active: number
  suspended: number
  new_today: number
}

export interface ListUsersParams {
  page: number
  page_size: number
  name?: string
  email?: string
  status?: number
  keyword?: string
}

export interface ListUsersResponse {
  list: User[]
  total: number
  page: number
  page_size: number
}

export interface ListExtendedUsersResponse {
  list: ExtendedUser[]
  total: number
  page: number
  page_size: number
}

export interface UpdateUserParams {
  name: string
  avatar?: string
}

export interface ResetPasswordResponse {
  temporary_password: string
}

export interface ListUserAddressesResponse {
  list: UserAddress[]
  total: number
}

// Get user info
export function getUserInfo(id: number) {
  return request<User>({
    url: `/api/v1/users/${id}`,
    method: 'get'
  })
}

// Get user list with pagination and filters
export function getUserList(params: ListUsersParams) {
  return request<ListUsersResponse>({
    url: '/api/v1/users',
    method: 'get',
    params
  })
}

// Update user info
export function updateUser(id: number, data: UpdateUserParams) {
  return request<User>({
    url: `/api/v1/users/${id}`,
    method: 'put',
    data
  })
}

// Suspend user
export function suspendUser(id: number) {
  return request<User>({
    url: `/api/v1/users/${id}/suspend`,
    method: 'post'
  })
}

// Activate user
export function activateUser(id: number) {
  return request<User>({
    url: `/api/v1/users/${id}/activate`,
    method: 'post'
  })
}

// Delete user (soft delete)
export function deleteUser(id: number) {
  return request({
    url: `/api/v1/users/${id}`,
    method: 'delete'
  })
}

// Reset user password
export function resetUserPassword(id: number) {
  return request<ResetPasswordResponse>({
    url: `/api/v1/users/${id}/reset-password`,
    method: 'post'
  })
}

// Get user statistics
export function getUserStats() {
  return request<UserStats>({
    url: '/api/v1/users/stats',
    method: 'get'
  })
}

// Enhanced user stats
export interface UserStatsEnhanced {
  total_users: number
  active_users: number
  suspended_users: number
  new_users_today: number
}

export function getUserStatsEnhanced() {
  return request<UserStatsEnhanced>({
    url: '/api/v1/users/stats/enhanced',
    method: 'get'
  })
}

// Get user detail with full information
export function getUserDetail(id: number) {
  return request<UserDetail>({
    url: `/api/v1/users/${id}/detail`,
    method: 'get'
  })
}

// Get user addresses
export function getUserAddresses(id: number, params?: { page?: number; page_size?: number }) {
  return request<ListUserAddressesResponse>({
    url: `/api/v1/users/${id}/addresses`,
    method: 'get',
    params
  })
}

// Suspend user with reason
export function suspendUserWithReason(id: number, reason: string) {
  return request<User>({
    url: `/api/v1/users/${id}/suspend-reason`,
    method: 'post',
    data: { reason }
  })
}

// Export users to file
export function exportUsers(params: ListUsersParams) {
  return request<Blob>({
    url: '/api/v1/users/export',
    method: 'get',
    params,
    responseType: 'blob'
  })
}

// Get enhanced user list with stats
export function getUserListEnhanced(params: ListUsersParams) {
  return request<ListExtendedUsersResponse>({
    url: '/api/v1/users/enhanced',
    method: 'get',
    params
  })
}