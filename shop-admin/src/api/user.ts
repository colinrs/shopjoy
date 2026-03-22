import request from '@/utils/request'

// User interface matching backend GetUserResponse
export interface User {
  id: number
  email: string
  phone: string
  name: string
  avatar: string
  status: number
  created_at: string
  last_login: string
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
}

export interface ListUsersResponse {
  list: User[]
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