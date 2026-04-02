import axios, { type AxiosError } from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { t } from '@/plugins/i18n'

/**
 * Download a file from the given URL with query parameters
 * @param url - The URL to download from
 * @param params - Query parameters
 * @param filename - Optional filename for the download
 */
export async function downloadFile(
  url: string,
  params: Record<string, unknown> = {},
  filename?: string
): Promise<void> {
  const userStore = useUserStore()

  try {
    const response = await axios.get(url, {
      params,
      headers: {
        Authorization: userStore.token ? `Bearer ${userStore.token}` : undefined,
        'X-Tenant-ID': '1'
      },
      responseType: 'blob'
    })

    // Get filename from response header or use provided filename
    const contentDisposition = response.headers['content-disposition']
    let downloadFilename = filename

    if (!downloadFilename && contentDisposition) {
      const match = contentDisposition.match(/filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/)
      if (match) {
        downloadFilename = decodeURIComponent(match[1].replace(/['"]/g, ''))
      }
    }

    // Create blob link and trigger download
    const blob = new Blob([response.data])
    const downloadUrl = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = downloadUrl

    if (downloadFilename) {
      link.download = downloadFilename
    } else {
      // If no filename, try to extract from URL
      const urlParts = url.split('/')
      link.download = urlParts[urlParts.length - 1].split('?')[0] || 'download'
    }

    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(downloadUrl)

    ElMessage.success(t('common.exportSuccess'))
  } catch (error: unknown) {
    // Try to parse error response as JSON to get error message
    const axiosError = error as AxiosError
    if (axiosError.response?.data) {
      try {
        const text = await (axiosError.response.data as Blob).text()
        const json = JSON.parse(text)
        ElMessage.error(json.msg || t('common.exportFailed'))
      } catch {
        ElMessage.error(t('common.exportFailed'))
      }
    } else {
      ElMessage.error(t('common.exportFailed'))
    }
    throw error
  }
}

/**
 * Open export URL in a new window (for server-side downloads)
 * @param url - The URL to open
 * @param params - Query parameters
 */
export function openExportUrl(url: string, params: Record<string, unknown> = {}): void {
  // Build URL with params
  const stringParams: Record<string, string> = {}
  for (const [key, value] of Object.entries(params)) {
    stringParams[key] = String(value ?? '')
  }
  const queryString = new URLSearchParams(stringParams).toString()
  const fullUrl = queryString ? `${url}?${queryString}` : url

  // Open in new window
  window.open(fullUrl, '_blank')

  ElMessage.success(t('common.exporting'))
}
