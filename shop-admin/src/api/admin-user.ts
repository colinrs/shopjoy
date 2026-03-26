import request from '@/utils/request'

// Admin User interface
export interface AdminUser {
  id: number
  email: string
  mobile: string
  real_name: string
  avatar: string
  type: number
  type_text: string
  status: number
  created_at: string
  updated_at: string
}

export interface AdminUserDetail extends AdminUser {
  roles: AdminRole[]
}

export interface AdminRole {
  id: number
  name: string
  code: string
}

// Auth APIs
export function adminLogin(data: {
  account: string
  password: string
  ip?: string
}) {
  return request<{ access_token: string; user: AdminUser }>({
    url: '/api/v1/auth/login',
    method: 'post',
    data
  })
}

export function registerTenantAdmin(data: {
  email: string
  mobile?: string
  real_name: string
  password: string
  tenant_id?: number
}) {
  return request({
    url: '/api/v1/auth/register',
    method: 'post',
    data
  })
}

// Admin User APIs
export interface GetAdminUsersParams {
  page: number
  page_size: number
  keyword?: string
  type?: number
  status?: number
  tenant_id?: number
}

export function getAdminUsers(params: GetAdminUsersParams) {
  return request<{ list: AdminUser[]; total: number; page: number; page_size: number }>({
    url: '/api/v1/admin-users',
    method: 'get',
    params
  })
}

export function getAdminUserDetail(id: number) {
  return request<AdminUserDetail>({
    url: `/api/v1/admin-users/${id}`,
    method: 'get'
  })
}

export interface CreateAdminUserParams {
  email: string
  mobile?: string
  real_name: string
  password: string
  type: number
  tenant_id?: number
}

export function createAdminUser(data: CreateAdminUserParams) {
  return request<AdminUser>({
    url: '/api/v1/admin-users',
    method: 'post',
    data
  })
}

export interface UpdateAdminUserParams {
  email?: string
  mobile?: string
  real_name?: string
  avatar?: string
}

export function updateAdminUser(id: number, data: UpdateAdminUserParams) {
  return request<AdminUser>({
    url: `/api/v1/admin-users/${id}`,
    method: 'put',
    data
  })
}

export function deleteAdminUser(id: number) {
  return request({
    url: `/api/v1/admin-users/${id}`,
    method: 'delete'
  })
}

export function resetAdminPassword(id: number) {
  return request<{ temporary_password: string }>({
    url: `/api/v1/admin-users/${id}/reset-password`,
    method: 'post'
  })
}

export function assignRoles(id: number, roleIds: number[]) {
  return request({
    url: `/api/v1/admin-users/${id}/roles`,
    method: 'put',
    data: { role_ids: roleIds }
  })
}

export interface RoleListResponse {
  list: AdminRole[]
  total: number
  page: number
  page_size: number
}

export function getAvailableRoles() {
  return request<RoleListResponse>({
    url: '/api/v1/roles',
    method: 'get'
  })
}

export function disableAdminUser(id: number) {
  return request({
    url: `/api/v1/admin-users/${id}/disable`,
    method: 'post',
    data: {}
  })
}

export function enableAdminUser(id: number) {
  return request({
    url: `/api/v1/admin-users/${id}/enable`,
    method: 'post',
    data: {}
  })
}

export function changeAdminPassword(data: {
  old_password: string
  new_password: string
  confirm_password: string
}) {
  return request({
    url: '/api/v1/admin-users/password',
    method: 'put',
    data
  })
}

export function updateAdminProfile(data: {
  user_id: number
  real_name?: string
  avatar?: string
  mobile?: string
  email?: string
}) {
  return request({
    url: '/api/v1/admin-users/profile',
    method: 'put',
    data
  })
}

// Customer/User APIs
export interface Customer {
  id: number
  name: string
  email: string
  mobile?: string
  avatar?: string
  status: number
  created_at: string
}

export function getCustomerList(params: {
  page: number
  page_size: number
  name?: string
  email?: string
}) {
  return request<{ list: Customer[]; total: number }>({
    url: '/api/v1/users',
    method: 'get',
    params
  })
}

export function getCustomerDetail(id: number) {
  return request({
    url: `/api/v1/users/${id}`,
    method: 'get'
  })
}

export function updateCustomer(id: number, data: {
  name: string
  avatar?: string
}) {
  return request({
    url: `/api/v1/users/${id}`,
    method: 'put',
    data
  })
}