import { ref, type Ref } from 'vue'
import { nanoid } from 'nanoid'
import imageCompression from 'browser-image-compression'
import { signImage, confirmImage, deleteImage } from '@/api/upload'
import type { UploadResponse } from '@/api/upload'

export type UploadStatus =
  | 'pending'
  | 'signing'
  | 'uploading'
  | 'confirming'
  | 'done'
  | 'error'
  | 'cancelled'

export interface UploadItem {
  uid: string
  status: UploadStatus
  progress: number // 0..100
  url?: string
  assetId?: string
  publicId?: string
  error?: string
  raw?: File
  preview?: string // object URL (revoked on remove)
  width?: number
  height?: number
  mimeType?: string
}

export interface UseImageUploadOptions {
  category: string
  maxSize?: number
  concurrency?: number
  compress?: boolean
}

/**
 * Read a File's dimensions from a data URL. Returns 0×0 on failure.
 */
function readImageDimensions(dataUrl: string): Promise<{ width: number; height: number }> {
  return new Promise((resolve) => {
    const img = new Image()
    img.onload = () => resolve({ width: img.naturalWidth, height: img.naturalHeight })
    img.onerror = () => resolve({ width: 0, height: 0 })
    img.src = dataUrl
  })
}

/**
 * Composable wrapping the sign → xhr → confirm direct-upload flow.
 *
 * State machine per item:
 *   pending → signing → uploading → confirming → done
 *                                                   └──► (or → error on failure, → cancelled on abort)
 *
 * Concurrency is capped via a tiny semaphore so the caller can dump N files in
 * without saturating the network.
 */
export function useImageUpload(opts: UseImageUploadOptions) {
  const items: Ref<UploadItem[]> = ref([])

  const category = opts.category
  // 5 MB upper bound (reserved for future size-guard).
  void opts.maxSize
  const concurrency = opts.concurrency ?? 3
  const compress = opts.compress ?? true

  // Semaphore state
  let inFlight = 0
  const queue: Array<() => void> = []

  function acquire(): Promise<void> {
    if (inFlight < concurrency) {
      inFlight++
      return Promise.resolve()
    }
    return new Promise<void>((resolve) => queue.push(() => { inFlight++; resolve() }))
  }
  function release(): void {
    inFlight--
    const next = queue.shift()
    if (next) next()
  }

  /**
   * Actually upload a single file:
   *   1. compress (optional, only for >1MB)
   *   2. signImage → get signed params + asset_id + upload_url
   *   3. XHR POST file to upload_url with progress
   *   4. confirmImage → persist asset metadata
   */
  async function runUpload(item: UploadItem, signal?: AbortSignal): Promise<void> {
    if (!item.raw) throw new Error('No file to upload')

    // --- 1. Compress ---
    if (compress && item.raw.size > 1024 * 1024) {
      try {
        item.raw = await imageCompression(item.raw, {
          maxSizeMB: 2,
          maxWidthOrHeight: 2048,
          useWebWorker: true
        })
      } catch (e) {
        // Compression failures are non-fatal — upload original instead.
        // Keep the original File as-is.
      }
    }

    item.status = 'signing'
    item.progress = 0

    // --- 2. Sign ---
    const sign = await signImage({
      category,
      filename: item.raw.name,
      mime_type: item.raw.type
    })
    if (signal?.aborted) throwObject()

    // --- 3. Direct XHR PUT/POST to Cloudinary ---
    item.status = 'uploading'
    const { url, width, height, publicId, format } = await uploadToCloudinary(
      sign.upload_url,
      item.raw,
      sign,
      signal,
      (p) => {
        item.progress = p
      }
    )
    if (width && height) {
      item.width = width
      item.height = height
    }

    // --- 4. Confirm ---
    item.status = 'confirming'
    const confirmed: UploadResponse = await confirmImage({
      asset_id: sign.asset_id,
      // Use the full public_id (folder + uuid) from Cloudinary's
      // response — the backend's RegisterAsset requires the folder
      // prefix. The bare sign.public_id would fail the prefix check.
      public_id: publicId,
      url,
      filename: item.raw.name,
      size: item.raw.size,
      mime_type: item.raw.type,
      width: item.width ?? 0,
      height: item.height ?? 0,
      // Prefer Cloudinary's detected format; fall back to the
      // extension parsed from the MIME type if Cloudinary didn't
      // report one (e.g. older responses).
      format: format || extFromMime(item.raw.type),
      category
    })

    item.status = 'done'
    item.progress = 100
    item.url = confirmed.url
    item.assetId = confirmed.id
    item.publicId = publicId
  }

  function throwObject(): never {
    const e: any = new Error('Upload cancelled')
    e.name = 'AbortError'
    throw e
  }

  function extFromMime(mime: string): string {
    if (!mime) return ''
    const idx = mime.indexOf('/')
    return idx >= 0 ? mime.slice(idx + 1) : ''
  }

  /**
   * Upload one file. Returns the UploadItem once it reaches 'done' state.
   * The item is registered into `items` as 'pending' immediately so callers can
   * render progress; concurrency is bounded by `concurrency`.
   */
  async function uploadOne(file: File): Promise<UploadItem> {
    const preview = URL.createObjectURL(file)
    const item: UploadItem = {
      uid: nanoid(),
      status: 'pending',
      progress: 0,
      raw: file,
      preview,
      mimeType: file.type
    }
    items.value.push(item)

    const controller = new AbortController()
    ;(item as any)._controller = controller

    await acquire()
    try {
      // Pre-compute dimensions for the confirm body.
      try {
        const dims = await readImageDimensions(preview)
        item.width = dims.width
        item.height = dims.height
      } catch {
        // ignore
      }

      await runUpload(item, controller.signal)
    } catch (err: any) {
      if (err?.name === 'AbortError' || controller.signal.aborted) {
        item.status = 'cancelled'
        item.error = 'cancelled'
      } else {
        item.status = 'error'
        item.error = err?.message || String(err)
      }
      // Rethrow to surface to caller; the item itself carries .status/error.
      throw err
    } finally {
      release()
    }
    return item
  }

  /**
   * Cancel an in-flight upload (xhr.abort()). Safe to call multiple times.
   */
  function cancel(item: UploadItem): void {
    const ctrl = (item as any)._controller as AbortController | undefined
    if (ctrl) ctrl.abort()
    item.status = 'cancelled'
  }

  /**
   * Retry a failed upload by re-running runUpload with the stored raw file.
   */
  async function retry(item: UploadItem): Promise<UploadItem> {
    if (!item.raw) throw new Error('No raw file to retry')
    if (item.status === 'uploading' || item.status === 'signing' || item.status === 'confirming') {
      return item // already in flight
    }
    item.status = 'pending'
    item.error = undefined
    item.progress = 0

    const controller = new AbortController()
    ;(item as any)._controller = controller

    await acquire()
    try {
      await runUpload(item, controller.signal)
    } catch (err: any) {
      if (err?.name === 'AbortError' || controller.signal.aborted) {
        item.status = 'cancelled'
      } else {
        item.status = 'error'
        item.error = err?.message || String(err)
      }
      throw err
    } finally {
      release()
    }
    return item
  }

  /**
   * Remove an item: deletes the asset on the backend (if uploaded) and revokes
   * any object URL.
   */
  async function remove(item: UploadItem): Promise<void> {
    if (item.preview) {
      try { URL.revokeObjectURL(item.preview) } catch { /* noop */ }
    }
    if (item.assetId) {
      try {
        await deleteImage(item.assetId)
      } catch {
        // Best-effort — backend may already be in a clean state.
      }
    }
    items.value = items.value.filter((it) => it.uid !== item.uid)
  }

  /**
   * Cancel everything, revoke object URLs, delete anything that reached 'done'.
   */
  async function clearAll(): Promise<void> {
    const snapshot = items.value.slice()
    for (const it of snapshot) {
      const ctrl = (it as any)._controller as AbortController | undefined
      if (ctrl) ctrl.abort()
    }
    for (const it of snapshot) {
      await remove(it).catch(() => undefined)
    }
    items.value = []
  }

  return {
    items,
    uploadOne,
    retry,
    cancel,
    remove,
    clearAll
  }
}

/**
 * POST a file to a Cloudinary signed-upload URL with progress + abort.
 * Returns the final URL and (best-effort) dimensions.
 */
function uploadToCloudinary(
  uploadUrl: string,
  file: File,
  sign: {
    cloud_name: string
    api_key: string
    timestamp: string
    signature: string
    folder: string
    public_id: string
    upload_preset?: string
  },
  signal: AbortSignal | undefined,
  onProgress: (pct: number) => void
): Promise<{ url: string; width: number; height: number; publicId: string; format: string }> {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest()

    if (signal) {
      if (signal.aborted) {
        reject(makeAbortError())
        return
      }
      signal.addEventListener('abort', () => {
        xhr.abort()
        reject(makeAbortError())
      })
    }

    xhr.upload.onprogress = (e) => {
      if (e.lengthComputable) {
        onProgress(Math.round((e.loaded / e.total) * 100))
      }
    }

    xhr.onload = () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        try {
          const resp = JSON.parse(xhr.responseText)
          const url: string = resp.secure_url || resp.url
          const width: number = resp.width || 0
          const height: number = resp.height || 0
          // Cloudinary returns the FULL public_id including the folder
          // prefix (e.g. "dev/0/product/<uuid>"). The sign endpoint only
          // hands out the bare UUID, but the backend's RegisterAsset
          // requires the folder-prefixed form for its tenant/category
          // guard. Falling back to `${sign.folder}/${sign.public_id}` is
          // safe because we signed that exact path.
          const publicId: string = resp.public_id || `${sign.folder}/${sign.public_id}`
          const format: string = resp.format || ''
          onProgress(100)
          resolve({ url, width, height, publicId, format })
        } catch (err: any) {
          reject(new Error(`Invalid Cloudinary response: ${err?.message || err}`))
        }
      } else {
        reject(new Error(`Cloudinary upload failed: ${xhr.status} ${xhr.statusText}`))
      }
    }

    xhr.onerror = () => reject(new Error('Network error during Cloudinary upload'))
    xhr.onabort = () => reject(makeAbortError())

    const form = new FormData()
    form.append('file', file)
    form.append('api_key', sign.api_key)
    form.append('timestamp', String(sign.timestamp))
    form.append('signature', sign.signature)
    form.append('folder', sign.folder)
    form.append('public_id', sign.public_id)
    if (sign.upload_preset) form.append('upload_preset', sign.upload_preset)

    xhr.open('POST', uploadUrl, true)
    xhr.send(form)
  })
}

function makeAbortError(): Error {
  const e: any = new Error('Upload cancelled')
  e.name = 'AbortError'
  return e
}
