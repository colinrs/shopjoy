import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import type { AdminUser } from '@/api/admin-user'

/**
 * User info stored in the user store
 * Matches the backend AdminUser response structure from admin login
 */
export type UserInfo = AdminUser

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')

  // Restore userInfo from localStorage if available
  const savedUserInfo = localStorage.getItem('user_info')
  const userInfo = ref<UserInfo | null>(savedUserInfo ? JSON.parse(savedUserInfo) : null)

  // Persist userInfo to localStorage whenever it changes
  watch(userInfo, (val) => {
    if (val) {
      localStorage.setItem('user_info', JSON.stringify(val))
    } else {
      localStorage.removeItem('user_info')
    }
  }, { deep: true })

  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const clearToken = () => {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user_info')
    localStorage.removeItem('target_tenant_id')
  }

  return {
    token,
    userInfo,
    setToken,
    clearToken
  }
})
