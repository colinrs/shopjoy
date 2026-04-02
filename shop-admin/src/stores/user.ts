import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { AdminUser } from '@/api/admin-user'

/**
 * User info stored in the user store
 * Matches the backend AdminUser response structure from admin login
 */
export type UserInfo = AdminUser

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref<UserInfo | null>(null)

  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const clearToken = () => {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
  }

  return {
    token,
    userInfo,
    setToken,
    clearToken
  }
})
