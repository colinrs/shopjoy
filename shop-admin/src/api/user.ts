import request from '@/utils/request'

export function getUserInfo(id: number) {
  return request({
    url: `/api/v1/users/${id}`,
    method: 'get'
  })
}

export function getUserList(params: {
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

export function updateUser(id: number, data: {
  name: string
  avatar?: string
}) {
  return request({
    url: `/api/v1/users/${id}`,
    method: 'put',
    data
  })
}