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
 * Sign response from backend for direct (browser-to-Cloudinary) upload.
 * Returned by POST /api/v1/uploads/sign.
 */
export interface UploadSignResponse {
  cloud_name: string
  api_key: string
  timestamp: string
  signature: string
  folder: string
  public_id: string
  upload_preset?: string
  asset_id: string
  upload_url: string
}

/**
 * Body for POST /api/v1/uploads/confirm.
 * Backend records an asset after the browser has uploaded directly to Cloudinary.
 */
export interface UploadConfirmRequest {
  asset_id?: string
  public_id: string
  url: string
  filename: string
  size: number
  mime_type: string
  width: number
  height: number
  format: string
  category: string
}

/**
 * Pre-sign: backend generates a Cloudinary signature for direct upload.
 * @param p category + filename + mime_type
 */
export function signImage(p: {
  category: string
  filename: string
  mime_type: string
}): Promise<UploadSignResponse> {
  return request<UploadSignResponse>({
    url: '/api/v1/uploads/sign',
    method: 'post',
    params: p
  })
}

/**
 * Notify backend that the direct upload finished, so the asset gets recorded.
 */
export function confirmImage(b: UploadConfirmRequest): Promise<UploadResponse> {
  return request<UploadResponse>({
    url: '/api/v1/uploads/confirm',
    method: 'post',
    data: b
  })
}

/**
 * Get asset metadata by id.
 */
export function getImage(id: string): Promise<UploadResponse> {
  return request<UploadResponse>({
    url: `/api/v1/uploads/${id}`,
    method: 'get'
  })
}

/**
 * Delete a previously uploaded asset. Backend enforces tenant scoping.
 */
export function deleteImage(id: string): Promise<void> {
  return request<void>({
    url: `/api/v1/uploads/${id}`,
    method: 'delete'
  })
}

/**
 * Legacy proxy upload (kept for fallback / admin-only flows).
 * New code should use signImage + direct-to-Cloudinary + confirmImage.
 */
export function uploadImage(file: File, category: UploadCategory = 'product'): Promise<UploadResponse> {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('category', category)

  return request<UploadResponse>({
    url: '/api/v1/uploads',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}
