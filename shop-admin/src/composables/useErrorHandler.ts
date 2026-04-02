/**
 * useErrorHandler Composable
 * Provides centralized error handling with ElMessage notifications
 */
import { ElMessage } from 'element-plus'

export function useErrorHandler() {
  /**
   * Handle an error by logging it and showing a user-friendly message
   * @param error - The error object or any value
   * @param customMessage - Optional custom message to display to user
   */
  const handleError = (error: unknown, customMessage?: string) => {
    console.error(error) // Keep for debugging in development

    const message = customMessage || (error instanceof Error ? error.message : 'An error occurred')
    ElMessage.error(message)

    // Could also send to error tracking service here
  }

  /**
   * Show a success message
   * @param message - The success message to display
   */
  const handleSuccess = (message: string) => {
    ElMessage.success(message)
  }

  return {
    handleError,
    handleSuccess
  }
}