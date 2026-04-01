import request from '@/utils/request'

// Role interface matching backend RoleInfo
export interface Role {
  id: number
  name: string
  code: string
  description: string
  status: number // 0=disabled, 1=enabled
  status_text: string
  is_system: boolean
  created_at: string
  updated_at: string
}

// Permission interface matching backend PermissionInfo
export interface Permission {
  id: number
  name: string
  code: string
  type: number // 0=menu, 1=button, 2=api
  type_text: string
  parent_id: number
  path: string
  icon: string
  sort: number
}

// Role with permissions (from GetRole API)
export interface RoleWithPermissions extends Role {
  permissions: Permission[]
}

// List roles params
export interface ListRolesParams {
  page: number
  page_size: number
  name?: string
  code?: string
  status?: number
}

// List roles response
export interface ListRolesResponse {
  list: Role[]
  total: number
  page: number
  page_size: number
}

// Create role request
export interface CreateRoleParams {
  name: string
  code: string
  description?: string
  permission_ids?: number[]
}

// Update role request
export interface UpdateRoleParams {
  name?: string
  description?: string
}

// Update role permissions request
export interface UpdateRolePermissionsParams {
  permission_ids: number[]
}

// Update role status request
export interface UpdateRoleStatusParams {
  status: number // 0=disabled, 1=enabled
}

// Get role list
export function getRoleList(params: ListRolesParams) {
  return request<ListRolesResponse>({
    url: '/api/v1/roles',
    method: 'get',
    params
  })
}

// Get role detail with permissions
export function getRoleDetail(id: number) {
  return request<RoleWithPermissions>({
    url: `/api/v1/roles/${id}`,
    method: 'get'
  })
}

// Create role
export function createRole(data: CreateRoleParams) {
  return request<{ id: number }>({
    url: '/api/v1/roles',
    method: 'post',
    data
  })
}

// Update role
export function updateRole(id: number, data: UpdateRoleParams) {
  return request<Role>({
    url: `/api/v1/roles/${id}`,
    method: 'put',
    data
  })
}

// Delete role
export function deleteRole(id: number) {
  return request({
    url: `/api/v1/roles/${id}`,
    method: 'delete'
  })
}

// Update role status
export function updateRoleStatus(id: number, data: UpdateRoleStatusParams) {
  return request<Role>({
    url: `/api/v1/roles/${id}/status`,
    method: 'put',
    data
  })
}

// Update role permissions
export function updateRolePermissions(id: number, data: UpdateRolePermissionsParams) {
  return request({
    url: `/api/v1/roles/${id}/permissions`,
    method: 'put',
    data
  })
}

// Get all permissions list
export function getPermissionList() {
  return request<{ list: Permission[] }>({
    url: '/api/v1/permissions',
    method: 'get'
  })
}
