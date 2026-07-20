# Cloudinary Image Storage — Design Spec

- Date: 2026-07-11
- Status: Draft — pending user review
- Layer: Cross-cutting (admin backend + shop-admin frontend)
- Supersedes: not directly; extends `docs/superpowers/specs/2026-03-26-upload-storage-design.md`

## Goal

Replace the current local-only image upload pipeline in `admin` with a multi-backend-capable storage layer whose new driver is **Cloudinary** (direct browser-to-Cloudinary signed upload). Existing `local` storage remains as a fallback / dev option. Every uploaded image is recorded in a new `media_assets` table so deletes and metadata queries are reliable. Frontend ships a reusable `ImageUploader` component.

## Non-Goals (YAGNI)

- Image transformation / variants / smart cropping on the backend (Cloudinary can do it; we just store the URL).
- Archival period before hard delete.
- OSS / S3 driver implementations (placeholders already exist in `factory.go`; not in this spec).
- Dedicated media-management admin UI page.
- Asset dedup / perceptual hashing.

## Decisions Locked

| # | Decision | Why |
|---|---|---|
| D1 | Storage layer remains a Go interface (`storage.Storage`); add a `cloudinary` driver | Keeps existing factory multi-backend design; local keeps working. |
| D2 | Frontend uploads directly to Cloudinary using backend-issued signature | Lower admin bandwidth, faster UX, matches Cloudinary recommendation. |
| D3 | `Storage.Save` is refactored to take an `io.Reader` (`AssetDraft`); `RegisterAsset` is added for direct-upload-finished assets | Single interface fits both proxy and direct-upload flows. |
| D4 | Backend signs uploads using `api_secret` from config (`Storage.Cloudinary.APISecret`); secret never leaves server | Standard Cloudinary signed-upload security. |
| D5 | New table `media_assets` (provider, public_id, url, mime, size, w/h, format, category, tenant_id, created_by) — soft delete | Single source of truth for asset metadata; supports safe delete + tenant isolation. |
| D6 | Cloudinary folder = `{env}/{tenant_id}/{category}` | Multi-tenant isolation; per-environment clarity. |
| D7 | Frontend ships `components/ImageUploader.vue` + `composables/useImageUpload.ts` + `components/ImageCropperDialog.vue`; single and multi mode supported | Replaces duplicated `<el-upload>` instances in product / banner / avatar pages. |

---

## 1. Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                       shop-admin (Vue 3)                        │
│  ImageUploader.vue  ←── useImageUpload.ts  ←── api/upload.ts     │
│       │                                                          │
│       │ POST /api/v1/uploads/sign                                │
│       │ POST /api/v1/uploads/confirm                             │
│       │ DELETE /api/v1/uploads/:id                              │
│       │ GET  /api/v1/uploads/:id                                │
│       ▼                                                          │
│      Cloudinary  (signed direct upload, no admin hop)           │
└─────────────────────────────┬───────────────────────────────────┘
                              │ JSON (no file bytes)
┌─────────────────────────────▼───────────────────────────────────┐
│                    admin API (go-zero)                          │
│                                                                 │
│  handler/uploads/{sign,upload,confirm,delete,get}_handler.go    │
│  logic/uploads/{...}.go                                         │
│  application/upload/service.go  (NEW)                           │
│                                                                 │
│  domain/media/entity.go (NEW), domain/media/repository.go       │
│  infrastructure/persistence/media_asset_repo.go (NEW)           │
│                                                                 │
│  infrastructure/storage/{storage.go,local.go,cloudinary.go,      │
│                          factory.go}  ←── interface refactor   │
└─────────────────────────────┬───────────────────────────────────┘
                              │
                  ┌───────────┴────────────┐
              local disk             Cloudinary SDK
              (./uploads)           (cloud_name account)
```

Layering: handler → logic → application/upload service → infrastructure/storage driver → media.Repository → GORM → `media_assets` table.

---

## 2. Storage Interface & Cloudinary Driver

### 2.1 New types (in `admin/internal/infrastructure/storage/storage.go`)

```go
type Category string
// keep existing: CategoryProduct / CategoryBanner / CategoryAvatar

type Asset struct {
    ID        string    // snowflake
    PublicID  string    // Cloudinary public_id or local path
    URL       string    // public-facing URL
    Filename  string
    Size      int64
    MimeType  string
    Width     int
    Height    int
    Format    string
    Category  Category
    Provider  string    // "local" | "cloudinary"
    TenantID  int64
    CreatedBy int64
    CreatedAt time.Time
}

type AssetDraft struct {
    Filename  string
    Reader    io.Reader
    MimeType  string
    Category  Category
    TenantID  int64
    CreatedBy int64
}

type RemoteAsset struct {
    PublicID   string
    URL        string
    Filename   string
    Size       int64
    MimeType   string
    Width      int
    Height     int
    Format     string
    Category   Category
    TenantID   int64
    CreatedBy  int64
}

type Storage interface {
    Save(ctx context.Context, draft AssetDraft) (*Asset, error)
    RegisterAsset(ctx context.Context, remote RemoteAsset) (*Asset, error)
    Delete(ctx context.Context, id string) error
    Get(ctx context.Context, id string) (*Asset, error)
    GetURL(ctx context.Context, id string) (string, error)
}

type Signer interface {
    Sign(ctx context.Context, p SignParams) (Signature, error)
}

type SignParams struct {
    Category  Category
    TenantID  int64
    Filename  string
    Timestamp int64
    Folder    string
}

type Signature struct {
    CloudName    string `json:"cloud_name"`
    APIKey       string `json:"api_key"`
    Timestamp    string `json:"timestamp"`
    Signature    string `json:"signature"`
    Folder       string `json:"folder"`
    PublicID     string `json:"public_id"`
    UploadPreset string `json:"upload_preset,omitempty"`
}
```

### 2.2 `cloudinaryStorage`

```go
type cloudinaryStorage struct {
    cloud *cloudinary.Cloudinary
    cfg   CloudinaryConfig
    repo  media.Repository
    idGen snowflake.Snowflake
}

func (s *cloudinaryStorage) Save(ctx, draft) (*Asset, error)            // optional: server-side upload path
func (s *cloudinaryStorage) RegisterAsset(ctx, remote) (*Asset, error) // primary path after direct upload
func (s *cloudinaryStorage) Delete(ctx, id) error                      // → Destroy API + DB soft delete
func (s *cloudinaryStorage) Get(ctx, id) (*Asset, error)
func (s *cloudinaryStorage) GetURL(ctx, id) (string, error)            // returns stored URL
func (s *cloudinaryStorage) Sign(ctx, p) (Signature, error)            // SHA1 over canonical params
```

- `Sign` uses canonical `params` sorted alphabetically, joined `k=v&k=v`, hashed with `sha1(params + api_secret)`.
- `Delete` first calls `uploader.Destroy{ PublicID }`; on failure logs and continues with DB soft delete (best-effort destructive action, but recovery is impossible once DB row is gone, so Destroy failure is logged with high severity — see §6).
- `RegisterAsset` writes a new row in `media_assets`; the unique key `(provider, public_id)` rejects duplicates with `code.ErrUploadConfirmFailed`.

### 2.3 `localStorage` adapter

- `Save(ctx, draft)` rewrites the existing `Save` to use `io.Reader`; the existing danger checks stay in the upload logic caller.
- `RegisterAsset` for local: derive `URL = /uploads/{category}/{yyyy}/{mm}/{dd}/{id}.{ext}`, write row.
- `Delete` and `Get`/`GetURL` use the same path lookup but go through `media.Repository` for metadata.

### 2.4 Factory

```go
const (
    StorageTypeLocal      StorageType = "local"
    StorageTypeCloudinary StorageType = "cloudinary"
    // StorageTypeOSS / StorageTypeS3 reserved but not implemented
)

type Config struct {
    Type       StorageType
    Local      LocalConfig
    Cloudinary CloudinaryConfig
}

type CloudinaryConfig struct {
    CloudName    string
    APIKey       string
    APISecret    string
    UploadPreset string
    Environment  string
    Secure       bool
}

func NewStorage(cfg Config, repo media.Repository, idGen snowflake.Snowflake) (Storage, error)
```

`NewStorage` returns an error (not panic) on missing config; the wiring in `svc/service_context.go` converts to one of the existing `code.Err*Upload*` errors.

---

## 3. HTTP API & Flow

### 3.1 Endpoints (added to `admin/desc/admin.api`)

```api
@server (
    group:      uploads
    middleware: AuthMiddleware
)
service admin-api {

    @handler UploadSignHandler
    post /api/v1/uploads/sign (UploadSignRequest) returns (UploadSignResponse)

    @handler UploadHandler
    post /api/v1/uploads (UploadRequest) returns (UploadResponse)

    @handler UploadConfirmHandler
    post /api/v1/uploads/confirm (UploadConfirmRequest) returns (UploadConfirmResponse)

    @handler GetUploadHandler
    get /api/v1/uploads/:id (GetUploadReq) returns (UploadResponse)

    @handler DeleteUploadHandler
    delete /api/v1/uploads/:id (DeleteUploadReq) returns (DeleteUploadResp)
}
```

### 3.2 Types (added to `admin/internal/types/types.go`)

```go
type UploadSignRequest struct {
    Category string `form:"category,optional"`
    Filename string `form:"filename,optional"`
    MimeType string `form:"mime_type,optional"`
}

type UploadSignResponse struct {
    CloudName    string `json:"cloud_name"`
    APIKey       string `json:"api_key"`
    Timestamp    string `json:"timestamp"`
    Signature    string `json:"signature"`
    Folder       string `json:"folder"`
    PublicID     string `json:"public_id"`
    UploadPreset string `json:"upload_preset,omitempty"`
    AssetID      string `json:"asset_id"`     // backend-generated snowflake
    UploadURL    string `json:"upload_url"`  // Cloudinary upload endpoint
}

type UploadConfirmRequest struct {
    AssetID      string `json:"asset_id,optional"`
    PublicID     string `json:"public_id"`
    URL          string `json:"url"`
    Filename     string `json:"filename"`
    Size         int64  `json:"size"`
    MimeType     string `json:"mime_type"`
    Width        int    `json:"width"`
    Height       int    `json:"height"`
    Format       string `json:"format"`
    Category     string `json:"category"`
    ChecksumEtag string `json:"checksum_etag,optional"`
}

type UploadConfirmResponse struct {
    ID        string `json:"id"`
    URL       string `json:"url"`
    Filename  string `json:"filename"`
    Category  string `json:"category"`
    Size      int64  `json:"size"`
    MimeType  string `json:"mime_type"`
    Width     int    `json:"width"`
    Height    int    `json:"height"`
    CreatedAt string `json:"created_at"`
}

type GetUploadReq   struct { ID string `path:"id"` }
type DeleteUploadReq struct { ID string `path:"id"` }
type DeleteUploadResp struct { ID string `json:"id"` }
```

`UploadRequest` / `UploadResponse` keep current shape.

### 3.3 Direct-upload sequence

```
Frontend            admin API            Cloudinary         media_assets
   │                    │                     │                   │
   │ sign(category, fn) │                     │                   │
   ├───────────────────▶│ Sign: authz + sign  │                   │
   │                    ├────────┐            │                   │
   │                    │ asset_id,           │                   │
   │                    │ folder, public_id,  │                   │
   │                    │ ts, signature       │                   │
   │◀──── {…sig + asset_id}┤                    │                   │
   │ upload (XHR, formdata with sig + file)    │                   │
   ├────────────────────────────────────────────▶│                   │
   │◀──── { secure_url, public_id, bytes, w/h } ┤                   │
   │ confirm(public_id, secure_url, w, h, …)   │                   │
   ├───────────────────▶│ Confirm: validate    │                   │
   │                    ├──────────────────────────────────────INSERT│
   │◀──── { id, url, … }┤                                          │
```

### 3.4 Defence-in-depth

| Rule | Where |
|---|---|
| `public_id` prefix must equal `folder` returned at sign | `UploadConfirmLogic` |
| `timestamp` ±5min window | `UploadSignLogic` (return error if request older than window) |
| `mime_type` ∈ whitelist (`image/jpeg, image/png, image/gif, image/webp`) | `UploadConfirmLogic` |
| magic byte validation (already in `UploadLogic`) | kept, server-side path only |
| Replay protection: `(asset_id, public_id)` recorded in Redis with TTL 10min; if sign is reused → return same signature but mark as already-paid | `UploadSignLogic` (optional v2) |
| Path traversal prevention on `public_id` (reject `..`, `//`) | `UploadSignLogic` (set `public_id = uuid`, no user input) |

### 3.5 Logic files

| Logic | Responsibilities |
|---|---|
| `UploadSignLogic.Sign` | authz → load tenant/user from ctx → build `SignParams` → `storage.Signer.Sign` → return |
| `UploadLogic.Upload` | existing magic-byte checks (file reader) → `storage.Save(draft{Reader: file})` |
| `UploadConfirmLogic.Confirm` | authz → path-prefix check → repo unique-key check → `storage.RegisterAsset` |
| `GetUploadLogic.Get` | authz → `storage.Get(id)` |
| `DeleteUploadLogic.Delete` | authz → `storage.Delete(id)` (tenant-scoped: reject cross-tenant) |

### 3.6 Frontend

#### 3.6.1 `shop-admin/src/api/upload.ts`

```typescript
export interface UploadSignResponse {
  cloud_name: string; api_key: string; timestamp: string;
  signature: string; folder: string; public_id: string;
  upload_preset?: string; asset_id: string; upload_url: string;
}
export interface UploadConfirmRequest {
  asset_id?: string; public_id: string; url: string;
  filename: string; size: number; mime_type: string;
  width: number; height: number; format: string; category: string;
}
export const signImage = (p: { category: string; filename: string; mime_type: string }) =>
  request<UploadSignResponse>({ url: '/api/v1/uploads/sign', method: 'post', params: p })
export const confirmImage = (b: UploadConfirmRequest) =>
  request<UploadResponse>({ url: '/api/v1/uploads/confirm', method: 'post', data: b })
export const deleteImage = (id: string) =>
  request<void>({ url: `/api/v1/uploads/${id}`, method: 'delete' })
export const getImage = (id: string) =>
  request<UploadResponse>({ url: `/api/v1/uploads/${id}`, method: 'get' })
```

#### 3.6.2 `shop-admin/src/composables/useImageUpload.ts`

```typescript
export type UploadStatus = 'pending' | 'signing' | 'uploading' | 'confirming' | 'done' | 'error'
export interface UploadItem {
  uid: string
  status: UploadStatus
  progress: number   // 0..100
  url?: string
  assetId?: string
  publicId?: string
  error?: string
  raw?: File
  preview?: string   // object URL
}
export function useImageUpload(opts: { category: string; maxSize?: number; concurrency?: number }) {
  const items = ref<UploadItem[]>([])
  const { category, maxSize = 5 * 1024 * 1024, concurrency = 3 } = opts
  async function uploadOne(file: File): Promise<UploadItem> { /* sign → xhr → confirm */ }
  async function remove(item: UploadItem): Promise<void>    { /* DELETE + revokeObjectURL */ }
  function clearAll(): void                                  { /* revoke all + reset */ }
  return { items, uploadOne, remove, clearAll }
}
```

Concurrency: simple semaphore limiting parallel `uploadOne` to `concurrency`.

#### 3.6.3 `shop-admin/src/components/ImageUploader.vue`

Modes:
- `v-model: value` is `string` ⇒ single-file mode.
- `v-model: value` is `string[]` ⇒ multi-file mode (driven by `multiple` prop).

Props:

```ts
interface Props {
  modelValue: string | string[]
  category: 'product' | 'banner' | 'avatar' | string
  multiple?: boolean
  max?: number
  maxSize?: number
  accept?: string
  compress?: boolean
  crop?: boolean
  cropRatio?: number
}
```

Behaviour:

| Concern | Implementation |
|---|---|
| state machine | `pending → signing → uploading → confirming → done` (any failure → `error`, click to retry) |
| client-side compression | `browser-image-compression` (>1MB auto-compress) |
| cropping | modal `ImageCropperDialog.vue` using `cropperjs` |
| drag-reorder | `vuedraggable` (only in multi mode) |
| preview | `URL.createObjectURL` (revoked on remove) |
| progress | `XMLHttpRequest.upload.onprogress` |
| retry | click on item in `error` state re-invokes `uploadOne(item.raw)` |
| cancel | `AbortController` for in-flight XHR; on cancel, call `POST /uploads/cleanup?asset_id=…` (admin endpoint added later) |

CORS: direct upload to Cloudinary needs none of our CORS — admin sign endpoint sets only the public fields (`cloud_name, api_key, signature, folder, public_id, timestamp`) — no cookies forwarded.

#### 3.6.4 File layout (frontend)

```
shop-admin/src/
├── api/upload.ts
├── composables/useImageUpload.ts
├── components/
│   ├── ImageUploader.vue          (main)
│   └── ImageCropperDialog.vue     (used when crop=true)
└── views/
    ├── products/form.vue          ← replaces <el-upload> with ImageUploader
    ├── banners/form.vue           ← replaces <el-upload> with ImageUploader
    └── admin_users/profile.vue    ← avatar: crop=true, ratio=1
```

#### 3.6.5 Pages that change

| Page | Before | After |
|---|---|---|
| Product create/edit | `<el-upload>` (multi) | `<ImageUploader v-model:value="form.images" category="product" multiple :max="9" :compress="true" />` |
| Banner create/edit | `<el-upload>` (single) | `<ImageUploader v-model:value="form.image" category="banner" :compress="true" />` |
| Admin profile avatar | file input | `<ImageUploader v-model:value="form.avatar" category="avatar" :crop="true" :crop-ratio="1" :compress="false" />` |

---

## 4. `media_assets` Table & Repository

### 4.1 Schema (merge into `sql/admin/schema.sql`)

```sql
CREATE TABLE media_assets (
    id              BIGINT          PRIMARY KEY,
    public_id       VARCHAR(255)    NOT NULL,
    url             VARCHAR(1024)   NOT NULL,
    filename        VARCHAR(255)    NOT NULL DEFAULT '',
    size_bytes      BIGINT          NOT NULL DEFAULT 0,
    mime_type       VARCHAR(64)     NOT NULL DEFAULT '',
    width           INT             NOT NULL DEFAULT 0,
    height          INT             NOT NULL DEFAULT 0,
    format          VARCHAR(32)     NOT NULL DEFAULT '',
    category        VARCHAR(32)     NOT NULL DEFAULT 'common',
    provider        VARCHAR(16)     NOT NULL DEFAULT 'local',
    tenant_id       BIGINT          NOT NULL,
    created_by      BIGINT          NOT NULL DEFAULT 0,
    created_at      TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP       NULL,

    CONSTRAINT uk_provider_public UNIQUE (provider, public_id),
    INDEX idx_tenant_category (tenant_id, category, deleted_at),
    INDEX idx_created_at (created_at)
);
```

### 4.2 Domain entity (`admin/internal/domain/media/entity.go`)

```go
type Asset struct {
    application.Model
    PublicID  string
    URL       string
    Filename  string
    SizeBytes int64
    MimeType  string
    Width     int
    Height    int
    Format    string
    Category  string
    Provider  string
    TenantID  int64
    CreatedBy int64
}
func (*Asset) TableName() string { return "media_assets" }
```

### 4.3 Repository

```go
type Repository interface {
    Insert(ctx context.Context, a *Asset) error
    FindByID(ctx context.Context, id int64) (*Asset, error)
    FindByPublicID(ctx context.Context, provider, publicID string) (*Asset, error)
    SoftDelete(ctx context.Context, id int64) error
    ListByTenant(ctx context.Context, tenantID int64, category string, page int) ([]*Asset, int64, error)
}
```

Implementation: GORM-based in `infrastructure/persistence/media_asset_repo.go`; soft delete uses GORM's `deleted_at` automatically.

### 4.4 Migrations

Per project rule "SQL Schema Consolidation: 每个领域只保留一个 schema.sql", we **append** the table above to `sql/admin/schema.sql`. No `migrations/` is added for this table.

---

## 5. Configuration, Secrets & Error Codes

### 5.1 `config.Config` (extend `admin/internal/config/config.go`)

```go
type StorageConf struct {
    Type       string                `json:"type,default=local"`
    Local      LocalStorageConf      `json:"local,optional"`
    Cloudinary CloudinaryStorageConf `json:"cloudinary,optional"`
}
type LocalStorageConf      struct { BasePath string `json:"base_path,default=./uploads"` }
type CloudinaryStorageConf struct {
    CloudName    string `json:"cloud_name"`
    APIKey       string `json:"api_key"`
    APISecret    string `json:"api_secret"`
    UploadPreset string `json:"upload_preset,optional"`
    Environment  string `json:"environment,default=dev"`
    Secure       bool   `json:"secure,default=true"`
}
```

### 5.2 YAML example (`admin/etc/admin-api.yaml`)

```yaml
Storage:
  Type: cloudinary
  Cloudinary:
    CloudName: ${CLOUDINARY_CLOUD_NAME}
    APIKey:    ${CLOUDINARY_API_KEY}
    APISecret: ${CLOUDINARY_API_SECRET}
    Environment: dev
    Secure: true
```

`api_secret` is loaded via env-only; the load step returns error (no default) if missing.

### 5.3 Error codes (`pkg/code/code.go` increments — Upload module is `240xxx`)

| Code | Symbol | Use |
|---|---|---|
| 240001 (existing) | `ErrUploadUnsupportedFileType` | magic byte / mime reject |
| 240002 (existing) | `ErrUploadFileSizeExceeded` | size > max |
| 240003 (existing) | `ErrUploadInvalidCategory` | category unknown |
| 240004 (existing) | `ErrUploadFailed` | generic |
| 240005 (existing) | `ErrUploadNotFound` | asset not found |
| **240006 (new)** | `ErrUploadSignFailed` | signing error / missing config |
| **240007 (new)** | `ErrUploadConfirmFailed` | path-prefix mismatch / duplicate / expired |
| **240008 (new)** | `ErrUploadProviderError` | Cloudinary 5xx / network |
| **240009 (new)** | `ErrUploadCrossTenantAccess` | admin queries asset in another tenant |

### 5.4 Security checklist

- `api_secret` only on server; never returned to frontend.
- Signature params are alphabetically sorted and signed with `sha1(params + api_secret)`.
- `public_id` is server-generated (`uuid.New().String()`) — never user-controlled.
- `folder` always prefixed with `s.cfg.Environment` and `tenant_id`.
- Confirm path enforces `strings.HasPrefix(req.PublicID, sig.Folder)`.
- Sign endpoint rate-limited (`5 req/s/user`, future enhancement).
- Magic byte validation enforced for any file path that still touches admin bytes (server-side proxy upload keeps existing checks).

---

## 6. Testing, Deployment, Rollback

### 6.1 Test matrix

| Type | Scope | Tool |
|---|---|---|
| Unit | `cloudinary.Sign` matches Cloudinary's documentation examples | `go test` |
| Unit | `cloudinaryStorage.RegisterAsset` rejects duplicate `(provider, public_id)` | `go test` + sqlmock |
| Unit | local storage `Save` returns `Asset` with `provider="local"` and writes row | `go test` + sqlmock |
| Integration | `POST /uploads/sign` authz → signature | `httptest` |
| Integration | `POST /uploads/confirm` path prefix check | `httptest` |
| Integration | `DELETE /uploads/:id` calls Destroy + soft-delete; destroy failure logged but doesn't fail the API | `httptest` + Cloudinary mock |
| Integration | `DELETE /uploads/:id` cross-tenant access returns `ErrUploadCrossTenantAccess` | `httptest` |
| E2E | Sign → direct upload (mock Cloudinary) → confirm → get → delete | happy path |
| Concurrency | Two parallel confirms for the same `public_id` only one succeeds | `httptest` |
| Frontend | `useImageUpload` transitions through every state; cancellation aborts XHR | Vitest + happy-dom |

### 6.2 Regression checklist

- [ ] Existing `localStorage` unit tests pass after interface refactor.
- [ ] Existing magic-byte validation in `UploadLogic` still fires for the proxy-upload path.
- [ ] `media_assets.uk_provider_public` enforces uniqueness (DB-level test).

### 6.3 Deployment phases

| Phase | Description | Gate |
|---|---|---|
| 0 | Interface refactor + `localStorage` adapter + new tests. Production stays on `Type: local`. | All unit + integration tests pass; smoke check. |
| 1 | Staging: switch `Type: cloudinary`, `Environment: staging`. Manual sign/upload/confirm. | Visual + DB inspection. |
| 2 | Production: switch `Type: cloudinary`, `Environment: prod`. | Feature flag gated; canary 10% → 100%. |

### 6.4 Rollback

- **Config-level**: change YAML `Storage.Type=local`. Old `media_assets` rows with `provider='cloudinary'` stay; new writes go to local disk.
- **DB-level**: never delete the `media_assets` table on rollback.
- **Code-level**: keep this PR atomic & revertible.

### 6.5 Risks and mitigations

| Risk | Mitigation |
|---|---|
| Cloudinary quota exhaustion | Feature flag `Storage.FallbackToLocalOnError` (`true` triggers local fallback after 3 consecutive Cloudinary 5xx). |
| Signature replay | `timestamp` window + `asset_id` recorded in Redis (10 min TTL) — collision returns same signature but commit in v2. |
| Direct-upload timeout | Frontend XHR `timeout: 30s`; admin endpoint `/uploads/cleanup` deletes orphan `asset_id` after 10 min (future). |
| Hard delete on Cloudinary fails | log with `ERROR` severity; local DB row still soft-deleted; asset may resurface visually (acceptable for v1). |
| Multi-tenant cross-access | Tenant check in `Get/Delete` logic; tests cover this. |

### 6.6 Observability

- Metrics (`pkg/observability`):
  - `cloudinary_sign_total` (counter)
  - `cloudinary_destroy_fail_total` (counter)
  - `media_assets_count` (gauge by tenant)
  - `uploads_inflight` (gauge)
- Logs: structured (`tenant_id, category, public_id, asset_id, latency_ms, provider`).
- Alert: Cloudinary 5xx > 0.5% sustained 5 min → DingTalk.

### 6.7 Documentation artefacts

- This file: `docs/superpowers/specs/2026-07-11-cloudinary-image-storage-design.md`
- Append to `docs/reference/error-codes.md` (sections 140006–140009).
- Append to `docs/cross-cutting/api/api-reference.md` (sign / confirm / get endpoints).
- Note in `docs/reference/database-overview.md` (new table).

---

## Open Questions

None at the time of writing. Will be appended if discovered during review.

---

## Implementation Order (preview for `writing-plans` phase)

1. Add `media_assets` table to `sql/admin/schema.sql`.
2. Create `domain/media` entity + repository interface.
3. Implement `persistence/media_asset_repo.go` (GORM).
4. Refactor `infrastructure/storage/storage.go` (new types: `Asset`, `AssetDraft`, `RemoteAsset`, `Storage`, `Signer`, etc.).
5. Adapt `localStorage` to new interface.
6. Add `infrastructure/storage/cloudinary.go` (driver + `Sign`).
7. Update `factory.go` to dispatch `cloudinary`.
8. Wire `svc/service_context.go` to inject `media.Repository` + `idGen` to `NewStorage`.
9. Add `config.Storage` and YAML schema.
10. `.api` + `make api`: generate handlers/logic stubs for `UploadSignHandler` / `UploadConfirmHandler` / `GetUploadHandler`.
11. Implement logics.
12. Add new error codes to `pkg/code/code.go`.
13. Frontend: `api/upload.ts`, `composables/useImageUpload.ts`, `components/ImageUploader.vue`, `components/ImageCropperDialog.vue`.
14. Replace usages in `views/products/form.vue`, `views/banners/form.vue`, `views/admin_users/profile.vue`.
15. Tests per §6.1.
16. Run `make build` (per project rule).
17. Run code review (`/review` or `simplify`) before merge.
