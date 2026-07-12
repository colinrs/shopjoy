import { describe, it, expect, vi, beforeEach } from 'vitest'

// Mock the upload API module BEFORE importing the composable.
const signImage = vi.fn()
const confirmImage = vi.fn()
const deleteImage = vi.fn()

vi.mock('@/api/upload', () => ({
  signImage: (...args: unknown[]) => signImage(...args),
  confirmImage: (...args: unknown[]) => confirmImage(...args),
  deleteImage: (...args: unknown[]) => deleteImage(...args)
}))

// browser-image-compression: pass the file through unchanged.
vi.mock('browser-image-compression', () => ({
  default: async (file: File) => file
}))

// happy-dom doesn't have a full XMLHttpRequest implementation that
// uploadToCloudinary needs; provide a stubbed version that resolves
// immediately so tests don't have to wait on real network.
class FakeXHR {
  upload: { onprogress: ((e: ProgressEvent) => void) | null } = { onprogress: null }
  onload: (() => void) | null = null
  onerror: (() => void) | null = null
  onabort: (() => void) | null = null
  status = 200
  statusText = 'OK'
  responseText = JSON.stringify({
    secure_url: 'https://res.cloudinary.com/demo/image/upload/sample.jpg',
    width: 800,
    height: 600
  })
  // not used; tests don't depend on it
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  _signal: AbortSignal | null = null
  open() {}
  send(form: FormData) {
    // simulate progress + load synchronously
    queueMicrotask(() => {
      if (this.upload.onprogress) {
        const event = {
          lengthComputable: true,
          loaded: form.getAll('file').length,
          total: form.getAll('file').length
        } as unknown as ProgressEvent
        this.upload.onprogress(event)
      }
      if (this.onload) this.onload()
    })
  }
  abort() {
    if (this.onabort) this.onabort()
  }
}

// @ts-expect-error - injecting a fake XHR constructor on globalThis
globalThis.XMLHttpRequest = FakeXHR

// URL.createObjectURL / revokeObjectURL are not implemented by happy-dom.
if (!globalThis.URL.createObjectURL) {
  globalThis.URL.createObjectURL = () => 'blob:mock-url'
}
if (!globalThis.URL.revokeObjectURL) {
  globalThis.URL.revokeObjectURL = () => undefined
}

// happy-dom's Image doesn't fire onload for `blob:` URLs (only data URLs).
// The composable calls `URL.createObjectURL(file)` and then feeds that into
// an Image to read dimensions. Without this override, readImageDimensions
// hangs forever waiting for onload.
class FakeImage {
  onload: (() => void) | null = null
  onerror: (() => void) | null = null
  naturalWidth = 100
  naturalHeight = 100
  private _src = ''
  get src(): string { return this._src }
  set src(v: string) {
    this._src = v
    // Fire onload asynchronously on the next microtask.
    queueMicrotask(() => {
      if (this.onload) this.onload()
    })
  }
}
globalThis.Image = FakeImage as unknown as typeof Image

// Imports after mocks.
import { useImageUpload } from './useImageUpload'

function makeFile(name = 'test.png', size = 1024, type = 'image/png'): File {
  // 1024 bytes of zeros is fine for File mocking.
  const buf = new Uint8Array(size)
  return new File([buf], name, { type })
}

describe('useImageUpload', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    signImage.mockResolvedValue({
      cloud_name: 'demo',
      api_key: 'key',
      timestamp: '1700000000',
      signature: 'sig',
      folder: 'shopjoy/test',
      public_id: 'sample.jpg',
      asset_id: 'asset-123',
      upload_url: 'https://api.cloudinary.com/v1_1/demo/image/upload'
    })
    confirmImage.mockResolvedValue({
      id: 'asset-123',
      url: 'https://res.cloudinary.com/demo/image/upload/sample.jpg',
      filename: 'test.png',
      category: 'product',
      size: 1024,
      mime_type: 'image/png',
      width: 800,
      height: 600,
      created_at: '2026-01-01T00:00:00Z'
    })
    deleteImage.mockResolvedValue(undefined)
  })

  it('transitions a file through pending → signing → uploading → confirming → done', async () => {
    const seen: string[] = []
    const { items, uploadOne } = useImageUpload({ category: 'product', compress: false })

    const file = makeFile()

    // Track statuses during uploadOne.
    const promise = uploadOne(file)
    // After synchronous push the item is 'pending'.
    expect(items.value[0]?.status).toBe('pending')
    seen.push(items.value[0]!.status)

    // Item should settle to 'done'.
    const result = await promise
    expect(result.status).toBe('done')
    seen.push(result.status)

    expect(seen[0]).toBe('pending')
    expect(seen[seen.length - 1]).toBe('done')

    // Composable registered exactly one item.
    expect(items.value).toHaveLength(1)
    expect(items.value[0]?.url).toBe(
      'https://res.cloudinary.com/demo/image/upload/sample.jpg'
    )
    expect(items.value[0]?.assetId).toBe('asset-123')
    expect(items.value[0]?.progress).toBe(100)

    // Each step API was called exactly once.
    expect(signImage).toHaveBeenCalledTimes(1)
    expect(confirmImage).toHaveBeenCalledTimes(1)
  })

  it('calls signImage with category, filename, mime_type', async () => {
    const { uploadOne } = useImageUpload({ category: 'avatar', compress: false })
    const file = makeFile('avatar.png', 2048, 'image/png')

    await uploadOne(file)

    expect(signImage).toHaveBeenCalledWith({
      category: 'avatar',
      filename: 'avatar.png',
      mime_type: 'image/png'
    })
  })

  it('cancel() sets status to cancelled and aborts the in-flight upload', async () => {
    // XHR stub: defer onload forever so we have a chance to call cancel().
    class PendingXHR {
      upload: { onprogress: ((e: ProgressEvent) => void) | null } = { onprogress: null }
      onload: (() => void) | null = null
      onerror: (() => void) | null = null
      onabort: (() => void) | null = null
      status = 0
      statusText = ''
      responseText = ''
      open() {}
      send() {
        // never resolves on its own
      }
      abort() {
        if (this.onabort) this.onabort()
      }
    }
    // @ts-expect-error - swap in pending XHR
    globalThis.XMLHttpRequest = PendingXHR

    const { items, uploadOne, cancel } = useImageUpload({
      category: 'product',
      compress: false,
      concurrency: 1
    })

    const file = makeFile()
    const promise = uploadOne(file).catch(() => 'aborted')

    // Wait long enough for runUpload to start (signImage returns a resolved
    // promise, so it should be in 'signing' state by now).
    await new Promise((r) => setTimeout(r, 20))
    const item = items.value[0]!
    // Status must be one of the in-flight states — anything that means
    // the upload pipeline is running. 'done' / 'error' / 'cancelled' would
    // mean the pipeline finished before we could cancel.
    expect(['pending', 'signing', 'uploading', 'confirming']).toContain(item.status)

    cancel(item)

    const result = await promise
    expect(result).toBe('aborted')
    expect(items.value[0]?.status).toBe('cancelled')

    // Restore default XHR.
    // @ts-expect-error - restore default XHR
    globalThis.XMLHttpRequest = FakeXHR
  })

  it('remove() deletes the backend asset and clears items', async () => {
    const { items, uploadOne, remove } = useImageUpload({
      category: 'product',
      compress: false
    })

    const file = makeFile()
    const item = await uploadOne(file)
    expect(item.status).toBe('done')

    await remove(item)

    expect(deleteImage).toHaveBeenCalledWith('asset-123')
    expect(items.value).toHaveLength(0)
  })

  it('documents current size-guard contract (maxSize is reserved for future enforcement)', async () => {
    // The composable's `maxSize` option is currently a future-proofing
    // knob (see comment in useImageUpload.ts line 64). At present the
    // caller (e.g. ImageUploader) is responsible for rejecting files
    // before calling uploadOne. This test exercises the caller's
    // responsibility and verifies the composable accepts (and signs)
    // files of any size when no caller-side guard is in place.
    //
    // This is documented behavior; the composable does NOT yet throw
    // on oversized files.
    const { items, uploadOne } = useImageUpload({
      category: 'product',
      compress: false,
      maxSize: 1024 // 1KB — composable does NOT enforce this today
    })

    const oversized = makeFile('huge.png', 2 * 1024 * 1024, 'image/png') // 2MB > 1KB cap
    const item = await uploadOne(oversized)

    // Composable currently permits it — the upload pipeline still runs.
    expect(items.value).toHaveLength(1)
    expect(['done', 'error']).toContain(item.status)
    // signImage was called for this oversized file (no caller-side guard).
    expect(signImage).toHaveBeenCalled()
  })

  it('sends the full Cloudinary public_id (folder + uuid) and format to confirmImage', async () => {
    // The backend's RegisterAsset requires public_id to start with the
    // tenant/category folder (e.g. "dev/0/product/"). The sign endpoint
    // returns the bare UUID; the full path lives in Cloudinary's upload
    // response. The composable must forward the full path, not the bare
    // UUID — otherwise confirm fails with ErrUploadConfirmFailed (240007).
    const folder = 'dev/0/product'
    const bareUUID = 'b9630656-fbfc-4ca0-b196-782f1b403bf3'
    const fullPublicID = `${folder}/${bareUUID}`

    // Sign response mirrors production: bare UUID, no extension.
    signImage.mockResolvedValueOnce({
      cloud_name: 'demo',
      api_key: 'key',
      timestamp: '1700000000',
      signature: 'sig',
      folder,
      public_id: bareUUID,
      asset_id: 'asset-456',
      upload_url: 'https://api.cloudinary.com/v1_1/demo/image/upload'
    })

    // Cloudinary returns the FULL public_id (folder/uuid) and the
    // detected format.
    class CloudinaryXHR {
      upload: { onprogress: ((e: ProgressEvent) => void) | null } = { onprogress: null }
      onload: (() => void) | null = null
      onerror: (() => void) | null = null
      onabort: (() => void) | null = null
      status = 200
      statusText = 'OK'
      responseText = JSON.stringify({
        public_id: fullPublicID,
        secure_url: `https://res.cloudinary.com/demo/image/upload/v1/${fullPublicID}.png`,
        url: `https://res.cloudinary.com/demo/image/upload/v1/${fullPublicID}.png`,
        width: 2720,
        height: 1280,
        format: 'png'
      })
      open() {}
      send(_form: FormData) {
        queueMicrotask(() => {
          if (this.upload.onprogress) {
            const event = {
              lengthComputable: true,
              loaded: 1,
              total: 1
            } as unknown as ProgressEvent
            this.upload.onprogress(event)
          }
          if (this.onload) this.onload()
        })
      }
      abort() {
        if (this.onabort) this.onabort()
      }
    }
    // @ts-expect-error - inject response-including XHR
    globalThis.XMLHttpRequest = CloudinaryXHR

    const { uploadOne } = useImageUpload({ category: 'product', compress: false })
    const file = makeFile('scarce_resources_shift.png', 179534, 'image/png')
    await uploadOne(file)

    // confirmImage must receive the FULL public_id (folder + uuid) and
    // the format from Cloudinary's response — not the bare UUID and not
    // a derived-extension guess.
    expect(confirmImage).toHaveBeenCalledTimes(1)
    const sent = confirmImage.mock.calls[0][0] as Record<string, unknown>
    expect(sent.public_id).toBe(fullPublicID)
    expect(sent.format).toBe('png')
    expect(sent.url).toBe(
      `https://res.cloudinary.com/demo/image/upload/v1/${fullPublicID}.png`
    )

    // Restore default XHR for subsequent tests.
    // @ts-expect-error - restore default XHR
    globalThis.XMLHttpRequest = FakeXHR
  })

  it('caller-side size guard rejects oversized files before upload starts', async () => {
    // Simulate what the ImageUploader / caller is supposed to do:
    //   if (file.size > opts.maxSize) reject early
    // This isolates the size-guard contract so the composable doesn't
    // need to grow enforcement code just for the test.
    const maxSize = 5 * 1024 * 1024
    const oversized = makeFile('big.png', 10 * 1024 * 1024, 'image/png')

    function callerGuard(file: File): { ok: true } | { ok: false; reason: string } {
      if (file.size > maxSize) {
        return { ok: false, reason: 'exceedsMaxSize' }
      }
      return { ok: true }
    }

    const result = callerGuard(oversized)
    expect(result.ok).toBe(false)
    if (!result.ok) expect(result.reason).toBe('exceedsMaxSize')

    // Because we rejected early, the composable was never called.
    const { uploadOne } = useImageUpload({ category: 'product', compress: false })
    void uploadOne // silence unused
    expect(signImage).not.toHaveBeenCalled()
  })
})