/**
 * useApi Composable
 * Provides centralized API error handling with ElMessage notifications
 */
import { ElMessage } from 'element-plus'
import { ref, type Ref } from 'vue'

/**
 * Options for API calls
 */
export interface ApiOptions<T = unknown> {
  /** Custom error message to show on failure */
  errorMessage?: string
  /** Custom success message to show on success */
  successMessage?: string
  /** Whether to show error message on failure (default: true) */
  showError?: boolean
  /** Whether to show success message on success (default: false) */
  showSuccess?: boolean
  /** Callback when the API call completes (both success and failure) */
  onComplete?: () => void
  /** Default value to return on error */
  defaultValue?: T
}

/**
 * Result type for API calls
 */
export interface ApiResult<T> {
  /** Whether the API call is currently in progress */
  loading: Ref<boolean>
  /** The data returned from the API call (or default value on error) */
  data: Ref<T | undefined>
  /** The error that occurred (if any) */
  error: Ref<Error | null>
  /** Execute the API call */
  execute: (...args: unknown[]) => Promise<T | undefined>
}

/**
 * Wrapper function for API calls with consistent error handling
 * Handles loading state, error messages, and provides typed responses
 */
export function useApi<T>(
  apiFn: (...args: unknown[]) => Promise<T>,
  options: ApiOptions<T> = {}
): ApiResult<T> {
  const loading = ref(false)
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const data = ref<any>(options.defaultValue)
  const error = ref<Error | null>(null)

  const execute = async (...args: unknown[]): Promise<T | undefined> => {
    loading.value = true
    error.value = null

    try {
      const result = await apiFn(...args) as T
      data.value = result

      if (options.showSuccess && options.successMessage) {
        ElMessage.success(options.successMessage)
      }

      return result
    } catch (err) {
      error.value = err as Error

      if (options.showError !== false && options.errorMessage) {
        ElMessage.error(options.errorMessage)
      }

      if (options.defaultValue !== undefined) {
        data.value = options.defaultValue
      }

      return options.defaultValue
    } finally {
      loading.value = false
      options.onComplete?.()
    }
  }

  return {
    loading,
    data,
    error,
    execute
  }
}

/**
 * Simple API wrapper that shows error messages automatically
 * Suitable for most API calls that don't need custom loading management
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export async function apiWrapper<T = any>(
  apiFn: (...args: unknown[]) => Promise<T>,
  args: unknown[],
  options: {
    errorMessage?: string
    defaultValue?: T
  } = {}
): Promise<T | undefined> {
  try {
    return await apiFn(...args)
  } catch (err) {
    const error = err as Error
    if (options.errorMessage) {
      ElMessage.error(options.errorMessage)
    } else {
      // Show the error message from the interceptor or a generic message
      ElMessage.error(error.message || 'An error occurred')
    }
    return options.defaultValue
  }
}

/**
 * Void API wrapper for fire-and-forget API calls
 * Shows error message if the call fails
 */
export async function voidApiCall(
  apiFn: (...args: unknown[]) => Promise<void>,
  args: unknown[],
  errorMessage: string = 'Operation failed'
): Promise<boolean> {
  try {
    await apiFn(...args)
    return true
  } catch {
    ElMessage.error(errorMessage)
    return false
  }
}

/**
 * Create a typed API handler with common error handling
 */
export function createApiHandler<TArgs extends unknown[], TResult>(
  apiFn: (...args: TArgs) => Promise<TResult>,
  options: {
    errorMessage?: string | ((args: TArgs) => string)
    successMessage?: string | ((args: TArgs) => string)
  } = {}
) {
  return {
    /**
     * Execute the API call with automatic error handling
     */
    async exec(...args: TArgs): Promise<TResult | undefined> {
      try {
        const result = await apiFn(...args)

        if (options.successMessage) {
          const message = typeof options.successMessage === 'function'
            ? options.successMessage(args)
            : options.successMessage
          ElMessage.success(message)
        }

        return result
      } catch (err) {
        const message = options.errorMessage
          ? typeof options.errorMessage === 'function'
            ? options.errorMessage(args)
            : options.errorMessage
          : (err as Error).message || 'An error occurred'

        ElMessage.error(message)
        return undefined
      }
    },

    /**
     * Execute and return boolean success status
     */
    async call(...args: TArgs): Promise<boolean> {
      try {
        await apiFn(...args)

        if (options.successMessage) {
          const message = typeof options.successMessage === 'function'
            ? options.successMessage(args)
            : options.successMessage
          ElMessage.success(message)
        }

        return true
      } catch (err) {
        const message = options.errorMessage
          ? typeof options.errorMessage === 'function'
            ? options.errorMessage(args)
            : options.errorMessage
          : (err as Error).message || 'An error occurred'

        ElMessage.error(message)
        return false
      }
    }
  }
}
