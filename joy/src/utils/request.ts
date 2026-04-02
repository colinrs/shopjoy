import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios'
import { useUserStore } from '@/stores/user'
import { handleCustomerError } from './error-codes'
import type { ApiResponse, ErrorResponse } from './types'

// Simple toast notification for customer side
const showToast = (message: string, type: 'error' | 'success' | 'info' = 'error') => {
  // Create a simple toast element
  const existingToast = document.getElementById('app-toast')
  if (existingToast) {
    existingToast.remove()
  }

  const toast = document.createElement('div')
  toast.id = 'app-toast'
  toast.textContent = message
  toast.style.cssText = `
    position: fixed;
    top: 20px;
    left: 50%;
    transform: translateX(-50%);
    padding: 12px 24px;
    border-radius: 8px;
    font-size: 14px;
    z-index: 10000;
    max-width: 80%;
    text-align: center;
    animation: fadeInOut 3s ease;
    ${type === 'error' ? 'background: #fef2f2; color: #dc2626; border: 1px solid #fecaca;' : ''}
    ${type === 'success' ? 'background: #f0fdf4; color: #16a34a; border: 1px solid #bbf7d0;' : ''}
    ${type === 'info' ? 'background: #eff6ff; color: #2563eb; border: 1px solid #bfdbfe;' : ''}
  `

  // Add animation styles
  const style = document.createElement('style')
  style.textContent = `
    @keyframes fadeInOut {
      0% { opacity: 0; transform: translateX(-50%) translateY(-10px); }
      10% { opacity: 1; transform: translateX(-50%) translateY(0); }
      90% { opacity: 1; transform: translateX(-50%) translateY(0); }
      100% { opacity: 0; transform: translateX(-50%) translateY(-10px); }
    }
  `
  document.head.appendChild(style)
  document.body.appendChild(toast)

  // Remove after 3 seconds
  setTimeout(() => {
    toast.remove()
  }, 3000)
}

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '',
  timeout: 10000
})

// Request interceptor
request.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    // Customer side - use tenant_id from userInfo, localStorage, or default to '1'
    // TODO: Derive tenant from domain/subdomain when multi-domain support is implemented
    const tenantId = userStore.userInfo?.tenant_id || localStorage.getItem('tenant_id') || '1'
    config.headers['X-Tenant-ID'] = String(tenantId)
    return config
  },
  (error: AxiosError) => {
    return Promise.reject(error)
  }
)

// Response interceptor
request.interceptors.response.use(
  (response) => {
    const res = response.data as ApiResponse

    // Business error code check (code !== 0)
    if (res.code !== 0) {
      const errorResponse: ErrorResponse = {
        code: res.code,
        msg: res.msg,
        httpStatus: response.status,
        data: res.data
      }

      const errorConfig = handleCustomerError(errorResponse)

      // Handle based on action type
      switch (errorConfig.action) {
        case 'logout':
          showToast(errorConfig.message, 'error')
          const userStore = useUserStore()
          userStore.clearToken()
          window.location.href = errorConfig.redirectPath || '/login'
          break
        case 'redirect':
          showToast(errorConfig.message, 'error')
          window.location.href = errorConfig.redirectPath || '/'
          break
        case 'toast':
          showToast(errorConfig.message, 'error')
          break
        case 'none':
          // Silent handling
          break
        default:
          showToast(errorConfig.message, 'error')
      }

      return Promise.reject(new Error(errorConfig.message))
    }

    return res.data
  },
  (error: AxiosError<ApiResponse>) => {
    const { response } = error

    // Handle HTTP errors (non-200 responses)
    if (response) {
      const errorResponse: ErrorResponse = {
        code: response.data?.code,
        msg: response.data?.msg,
        message: error.message,
        httpStatus: response.status
      }

      const errorConfig = handleCustomerError(errorResponse)

      switch (errorConfig.action) {
        case 'logout':
          showToast(errorConfig.message, 'error')
          const userStore = useUserStore()
          userStore.clearToken()
          window.location.href = errorConfig.redirectPath || '/login'
          break
        case 'redirect':
          showToast(errorConfig.message, 'error')
          window.location.href = errorConfig.redirectPath || '/'
          break
        case 'toast':
          showToast(errorConfig.message, 'error')
          break
        default:
          showToast(errorConfig.message, 'error')
      }

      return Promise.reject(error)
    }

    // Network error (no response)
    showToast('网络连接失败，请检查网络', 'error')
    return Promise.reject(error)
  }
)

export default request
export { showToast }