import axios, { AxiosError, InternalAxiosRequestConfig, AxiosRequestConfig } from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { handleAdminError } from './error-codes'
import type { ApiResponse, ErrorResponse } from './types'

// Custom request interface that returns T directly (not AxiosResponse<T>)
interface RequestInstance {
  <T = unknown>(config: AxiosRequestConfig): Promise<T>
  get<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T>
  post<T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T>
  put<T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T>
  delete<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T>
}

const axiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '',
  timeout: 10000
})

// Request interceptor
axiosInstance.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    config.headers['X-Tenant-ID'] = '1'
    return config
  },
  (error: AxiosError) => {
    return Promise.reject(error)
  }
)

// Response interceptor
axiosInstance.interceptors.response.use(
  (response): any => {
    const res = response.data as ApiResponse

    // Business error code check (code !== 0)
    if (res.code !== 0) {
      const errorResponse: ErrorResponse = {
        code: res.code,
        msg: res.msg,
        httpStatus: response.status,
        data: res.data
      }

      const errorConfig = handleAdminError(errorResponse)

      // Handle based on action type
      switch (errorConfig.action) {
        case 'logout':
          ElMessage.error(errorConfig.message)
          const userStore = useUserStore()
          userStore.clearToken()
          window.location.href = errorConfig.redirectPath || '/login'
          break
        case 'redirect':
          ElMessage.error(errorConfig.message)
          window.location.href = errorConfig.redirectPath || '/'
          break
        case 'toast':
          ElMessage.error(errorConfig.message)
          break
        case 'none':
          // Silent handling
          break
        default:
          ElMessage.error(errorConfig.message)
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

      const errorConfig = handleAdminError(errorResponse)

      switch (errorConfig.action) {
        case 'logout':
          ElMessage.error(errorConfig.message)
          const userStore = useUserStore()
          userStore.clearToken()
          window.location.href = errorConfig.redirectPath || '/login'
          break
        case 'redirect':
          ElMessage.error(errorConfig.message)
          window.location.href = errorConfig.redirectPath || '/'
          break
        case 'toast':
          ElMessage.error(errorConfig.message)
          break
        default:
          ElMessage.error(errorConfig.message)
      }

      return Promise.reject(error)
    }

    // Network error (no response)
    ElMessage.error('网络连接失败，请检查网络')
    return Promise.reject(error)
  }
)

// Export typed request wrapper
const request = Object.assign(
  (config: AxiosRequestConfig) => axiosInstance(config) as Promise<any>,
  {
    get: (url: string, config?: AxiosRequestConfig) => axiosInstance.get(url, config) as Promise<any>,
    post: (url: string, data?: unknown, config?: AxiosRequestConfig) => axiosInstance.post(url, data, config) as Promise<any>,
    put: (url: string, data?: unknown, config?: AxiosRequestConfig) => axiosInstance.put(url, data, config) as Promise<any>,
    delete: (url: string, config?: AxiosRequestConfig) => axiosInstance.delete(url, config) as Promise<any>
  }
) as RequestInstance

export default request