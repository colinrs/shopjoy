# 图片上传接口与存储抽象设计

**日期：** 2026-03-26
**状态：** 已批准

## 1. 概述

设计一套通用的图片上传 HTTP 接口，以及对存储层的抽象接口。初期实现本地存储，后续可灵活扩展到 OSS、S3 等云存储。

## 2. API 接口设计

### 2.1 上传图片

**接口路径：** `POST /api/v1/uploads`

| 项目 | 内容 |
|-----|------|
| Content-Type | `multipart/form-data` |
| 认证 | 需要（AuthMiddleware） |

**Request 参数：**

| 字段 | 类型 | 必填 | 说明 |
|-----|------|-----|------|
| file | binary | 是 | 图片文件 |
| category | string | 是 | 业务场景（product/banner/avatar 等） |

**Response：**

```json
{
  "id": "img_abc123",
  "url": "/uploads/product/2026/03/26/img_abc123.jpg",
  "filename": "original_name.jpg",
  "category": "product",
  "size": 1024000,
  "mime_type": "image/jpeg",
  "width": 800,
  "height": 600,
  "created_at": "2026-03-26T10:30:00Z"
}
```

### 2.2 删除图片

**接口路径：** `DELETE /api/v1/uploads/{id}`

| 项目 | 内容 |
|-----|------|
| 认证 | 需要（AuthMiddleware） |

**Response：** 无返回内容（204 No Content）

### 2.3 错误码

| 状态码 | 错误码 | 说明 |
|-------|-------|------|
| 400 | 200001 | 文件类型不支持 |
| 400 | 200002 | 文件大小超限 |
| 400 | 200003 | category 无效 |
| 401 | 130001 | 未授权 |
| 404 | 200005 | 文件不存在 |
| 500 | 200004 | 上传失败 |

## 3. 存储抽象接口设计

### 3.1 接口定义

**文件位置：** `admin/internal/infrastructure/storage/storage.go`

```go
// Category 图片分类
type Category string

const (
    CategoryProduct Category = "product"
    CategoryBanner  Category = "banner"
    CategoryAvatar   Category = "avatar"
)

// FileInfo 文件信息
type FileInfo struct {
    ID        string    // 文件唯一ID
    Name      string    // 原始文件名
    Path      string    // 存储路径
    Size      int64     // 文件大小
    MimeType  string    // MIME类型
    Width     int       // 图片宽度
    Height    int       // 图片高度
    Category  Category  // 业务分类
    CreatedAt time.Time // 上传时间
}

// Storage 存储接口
type Storage interface {
    // Save 保存文件，返回文件信息
    Save(ctx context.Context, file *multipart.FileHeader, category Category) (*FileInfo, error)

    // Delete 删除文件
    Delete(ctx context.Context, id string) error

    // GetURL 获取文件访问URL
    GetURL(ctx context.Context, id string) (string, error)

    // Get 获取文件信息
    Get(ctx context.Context, id string) (*FileInfo, error)
}
```

## 4. 本地存储实现

### 4.1 实现细节

**文件位置：** `admin/internal/infrastructure/storage/local.go`

| 项目 | 内容 |
|-----|------|
| 存储路径 | `/uploads/{category}/{year}/{month}/{day}/` |
| 文件命名 | `{uuid}.{ext}` |
| 支持格式 | jpg, jpeg, png, gif, webp |
| 大小限制 | 5MB |
| 静态服务 | `/uploads/*` 路径映射到本地存储目录 |

### 4.2 图片处理

- 上传时自动获取图片宽高
- 保持原始文件扩展名

## 5. 目录结构

```
admin/internal/
├── logic/
│   └── upload/
│       └── upload_logic.go          # 业务逻辑
└── infrastructure/
    └── storage/
        ├── storage.go               # 接口定义
        ├── local.go                 # 本地存储实现
        └── factory.go               # 存储工厂
```

## 6. 扩展点

| 扩展方向 | 说明 |
|---------|------|
| OSS | 实现阿里云 OSS Storage |
| S3 | 实现 AWS S3 Storage |
| CDN | GetURL 返回 CDN 域名 |

## 7. 验收标准

- [ ] 上传接口支持 multipart/form-data
- [ ] 校验文件类型（仅图片）
- [ ] 校验文件大小（最大 5MB）
- [ ] 支持 category 分类
- [ ] 返回文件 ID + URL
- [ ] Storage 接口可替换实现
- [ ] 本地存储正常工作