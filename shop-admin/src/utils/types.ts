/**
 * Common error handling types
 */

export interface ErrorResponse {
  code?: number
  msg?: string
  message?: string
  httpStatus?: number
  data?: unknown
}

export interface ErrorCodeHandler {
  (error: ErrorResponse): {
    message: string
    action?: 'logout' | 'redirect' | 'toast' | 'none'
    redirectPath?: string
  }
}

export interface ApiResponse<T = unknown> {
  code: number
  msg: string
  data: T
}