import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref<any>(null)
  const cartCount = ref(0)

  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const clearToken = () => {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
  }

  const setCartCount = (count: number) => {
    cartCount.value = count
  }

  return {
    token,
    userInfo,
    cartCount,
    setToken,
    clearToken,
    setCartCount
  }
})
