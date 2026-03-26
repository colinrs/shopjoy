import request from '@/utils/request'

// Upload category types
export type UploadCategory = 'product' | 'banner' | 'avatar'

// Upload response from backend
export interface UploadResponse {
  id: string
  url: string
  filename: string
  category: string
  size: number
  mime_type: string
  width: number
  height: number
  created_at: string
}

/**
 * Upload image file
 * @param file - File object from input
 * @param category - Upload category (product/banner/avatar)
 * @returns Upload response with file ID and URL
 */
export function uploadImage(file: File, category: UploadCategory = 'product'): Promise<UploadResponse> {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('category', category)

  return request({
    url: '/api/v1/uploads',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

/**
 * Delete uploaded image
 * @param id - File ID to delete
 */
export function deleteImage(id: string): Promise<void> {
  return request({
    url: `/api/v1/uploads/${id}`,
    method: 'delete'
  })
}