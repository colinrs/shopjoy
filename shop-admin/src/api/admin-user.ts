import request from '@/utils/request'

// Auth APIs
export function adminLogin(data: {
  account: string
  password: string
  ip?: string
}) {
  return request({
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
export function getAdminUsers(params: {
  page: number
  page_size: number
  keyword?: string
  type?: number
  status?: number
  tenant_id?: number
}) {
  return request({
    url: '/api/v1/admin-users',
    method: 'get',
    params
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
export function getCustomerList(params: {
  page: number
  page_size: number
  name?: string
  email?: string
}) {
  return request({
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