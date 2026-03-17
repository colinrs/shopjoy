import request from '@/utils/request'

export function login(data: any) {
  return request({
    url: '/users/login',
    method: 'post',
    data
  })
}

export function register(data: any) {
  return request({
    url: '/users/register',
    method: 'post',
    data
  })
}

export function getUserInfo(id: number) {
  return request({
    url: `/users/${id}`,
    method: 'get'
  })
}

export function getUserList(params: any) {
  return request({
    url: '/users',
    method: 'get',
    params
  })
}

export function updateUser(id: number, data: any) {
  return request({
    url: `/users/${id}`,
    method: 'put',
    data
  })
}

export function changePassword(data: any) {
  return request({
    url: '/users/password',
    method: 'put',
    data
  })
}
