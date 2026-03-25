# 图片上传接口与存储抽象实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现通用图片上传 HTTP 接口及存储抽象，支持本地存储并预留云存储扩展点

**Architecture:** 使用 go-zero 框架，通过 .api 文件定义 HTTP 接口，storage 接口抽象存储层，本地存储实现为默认方案

**Tech Stack:** Go, go-zero, multipart/form-data, 依赖注入

---

## 文件结构

```
admin/
├── desc/
│   └── upload.api                 # API 定义文件
├── internal/
│   ├── handler/
│   │   └── upload_handler.go       # HTTP 处理器（自动生成，需修改）
│   ├── logic/
│   │   └── upload/
│   │       └── upload_logic.go    # 上传业务逻辑
│   ├── types/
│   │   └── types.go               # 自动生成
│   └── infrastructure/
│       └── storage/
│           ├── storage.go         # 存储接口定义
│           ├── local.go           # 本地存储实现
│           └── factory.go         # 存储工厂
pkg/code/code.go                    # 错误码定义
```

---

## Task 1: 定义 API 接口和错误码

**Files:**
- Create: `admin/desc/upload.api`
- Modify: `pkg/code/code.go`

- [ ] **Step 1: 创建 upload.api 文件**

```go
syntax = "v1"

info (
	title:   "Upload API"
	desc:    "通用图片上传接口"
	version: "v1"
)

type (
	UploadRequest {
		// 图片文件
		File *multipart.FileHeader `json:"file,optional"`
		// 业务分类：product/banner/avatar
		Category string `json:"category,optional"`
	}

	UploadResponse {
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

	DeleteUploadReq {
		ID string `path:"id"`
	}
)

@server (
	group:      uploads
	middleware: AuthMiddleware
)
service admin-api {
	@doc "上传图片"
	@handler UploadHandler
	post /api/v1/uploads (UploadRequest) returns (UploadResponse)

	@doc "删除图片"
	@handler DeleteUploadHandler
	delete /api/v1/uploads/:id returns
}
```

- [ ] **Step 2: 添加错误码到 pkg/code/code.go**

在 Err 定义处添加：
```go
// Upload Module (20xxx)
ErrUploadUnsupportedFileType = &Err{HTTPCode: http.StatusBadRequest, Code: 200001, Msg: "unsupported file type"}
ErrUploadFileSizeExceeded    = &Err{HTTPCode: http.StatusBadRequest, Code: 200002, Msg: "file size exceeded"}
ErrUploadInvalidCategory     = &Err{HTTPCode: http.StatusBadRequest, Code: 200003, Msg: "invalid category"}
ErrUploadFailed              = &Err{HTTPCode: http.StatusInternalServerError, Code: 200004, Msg: "upload failed"}
ErrUploadNotFound            = &Err{HTTPCode: http.StatusNotFound, Code: 200005, Msg: "file not found"}
```

- [ ] **Step 3: 运行 make api 生成代码**

```bash
cd admin && make api
```

- [ ] **Step 4: Commit**

```bash
git add admin/desc/upload.api pkg/code/code.go
git commit -m "feat(upload): add upload API definition and error codes"
```

---

## Task 2: 实现存储接口

**Files:**
- Create: `admin/internal/infrastructure/storage/storage.go`
- Create: `admin/internal/infrastructure/storage/local.go`
- Create: `admin/internal/infrastructure/storage/factory.go`

- [ ] **Step 1: 创建 storage.go 接口定义**

```go
package storage

import (
	"context"
	"mime/multipart"
	"time"
)

// Category 图片分类
type Category string

const (
	CategoryProduct Category = "product"
	CategoryBanner  Category = "banner"
	CategoryAvatar  Category = "avatar"
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

- [ ] **Step 2: 创建 local.go 本地存储实现**

```go
package storage

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	maxFileSize       = 5 * 1024 * 1024 // 5MB
	allowedExtensions = ".jpg,.jpeg,.png,.gif,.webp"
	uploadDir         = "./uploads"
)

var (
	allowedMimes = map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
)

type localStorage struct {
	basePath string
	redis    *redis.Redis
}

// NewLocalStorage 创建本地存储实例
func NewLocalStorage(redis *redis.Redis) Storage {
	return &localStorage{
		basePath: uploadDir,
		redis:    redis,
	}
}

func (s *localStorage) Save(ctx context.Context, file *multipart.FileHeader, category Category) (*FileInfo, error) {
	// 验证文件大小
	if file.Size > maxFileSize {
		return nil, fmt.Errorf("file size exceeds 5MB")
	}

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !strings.Contains(allowedExtensions, ext) {
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}

	// 生成唯一ID
	id := fmt.Sprintf("img_%s", uuid.New().String()[:12])

	// 构建存储路径: /uploads/{category}/{year}/{month}/{day}/
	now := time.Now()
	dir := filepath.Join(s.basePath, string(category), fmt.Sprintf("%d", now.Year()), fmt.Sprintf("%02d", now.Month()), fmt.Sprintf("%02d", now.Day()))

	// 创建目录
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create directory failed: %w", err)
	}

	// 文件完整路径
	filePath := filepath.Join(dir, id+ext)

	// 打开上传文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("open file failed: %w", err)
	}
	defer src.Close()

	// 获取图片尺寸
	width, height := 0, 0
	if ext != ".gif" {
		img, _, err := image.DecodeConfig(src)
		if err == nil {
			width, height = img.Width, img.Height
		}
		// 重新打开文件
		src.Seek(0, 0)
	}

	// 保存文件
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("create file failed: %w", err)
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(src); err != nil {
		return nil, fmt.Errorf("save file failed: %w", err)
	}

	// 构建相对路径
	relPath := fmt.Sprintf("/uploads/%s/%d/%02d/%02d/%s%s", category, now.Year(), now.Month(), now.Day(), id, ext)

	fileInfo := &FileInfo{
		ID:        id,
		Name:      file.Filename,
		Path:      relPath,
		Size:      file.Size,
		MimeType:  file.Header.Get("Content-Type"),
		Width:     width,
		Height:    height,
		Category:  category,
		CreatedAt: now,
	}

	// 可以选择将文件信息存入 Redis 便于查询
	// 这里简化处理，仅返回文件信息

	return fileInfo, nil
}

func (s *localStorage) Delete(ctx context.Context, id string) error {
	// 从 Redis 或数据库获取文件路径
	// 这里简化处理，构建可能的路径进行删除
	// 实际项目中应该存储文件元数据

	// 尝试删除常见目录下的文件
	now := time.Now()
	for _, category := range []Category{CategoryProduct, CategoryBanner, CategoryAvatar} {
		for day := now.Day() - 7; day <= now.Day()+1; day++ {
			for _, ext := range []string{".jpg", ".jpeg", ".png", ".gif", ".webp"} {
				path := filepath.Join(s.basePath, string(category), fmt.Sprintf("%d", now.Year()), fmt.Sprintf("%02d", now.Month()), fmt.Sprintf("%02d", day), id+ext)
				if _, err := os.Stat(path); err == nil {
					return os.Remove(path)
				}
			}
		}
	}

	return fmt.Errorf("file not found")
}

func (s *localStorage) GetURL(ctx context.Context, id string) (string, error) {
	// 简化实现：返回相对路径
	// 实际项目中应该从元数据存储中获取完整路径
	now := time.Now()
	for _, category := range []Category{CategoryProduct, CategoryBanner, CategoryAvatar} {
		for day := now.Day() - 7; day <= now.Day()+1; day++ {
			for _, ext := range []string{".jpg", ".jpeg", ".png", ".gif", ".webp"} {
				path := fmt.Sprintf("/uploads/%s/%d/%02d/%02d/%s%s", category, now.Year(), now.Month(), now.Day(), id, ext)
				fullPath := filepath.Join(s.basePath, string(category), fmt.Sprintf("%d", now.Year()), fmt.Sprintf("%02d", now.Month()), fmt.Sprintf("%02d", day), id+ext)
				if _, err := os.Stat(fullPath); err == nil {
					return path, nil
				}
			}
		}
	}

	return "", fmt.Errorf("file not found")
}

func (s *localStorage) Get(ctx context.Context, id string) (*FileInfo, error) {
	// 简化实现：返回 URL
	url, err := s.GetURL(ctx, id)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		ID:        id,
		Path:      url,
		CreatedAt: time.Now(),
	}, nil
}
```

- [ ] **Step 3: 创建 factory.go 存储工厂**

```go
package storage

import "github.com/zeromicro/go-zero/core/stores/redis"

// StorageType 存储类型
type StorageType string

const (
	StorageTypeLocal StorageType = "local"
	StorageTypeOSS   StorageType = "oss"
	StorageTypeS3    StorageType = "s3"
)

// Config 存储配置
type Config struct {
	Type   StorageType
	Local  LocalConfig
	OSS    OSSConfig
	S3     S3Config
}

// LocalConfig 本地存储配置
type LocalConfig struct {
	BasePath string
}

// OSSConfig OSS 配置
type OSSConfig struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	Bucket          string
}

// S3Config S3 配置
type S3Config struct {
	Endpoint   string
	Region     string
	AccessKey  string
	SecretKey  string
	Bucket     string
}

// NewStorage 创建存储实例
func NewStorage(cfg Config, redis *redis.Redis) (Storage, error) {
	switch cfg.Type {
	case StorageTypeLocal:
		localCfg := cfg.Local
		if localCfg.BasePath == "" {
			localCfg.BasePath = "./uploads"
		}
		return &localStorage{
			basePath: localCfg.BasePath,
			redis:    redis,
		}, nil
	// 后续可扩展 OSS、S3
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.Type)
	}
}
```

- [ ] **Step 4: Commit**

```bash
git add admin/internal/infrastructure/storage/
git commit -m "feat(storage): implement storage interface and local storage"
```

---

## Task 3: 实现上传业务逻辑

**Files:**
- Create: `admin/internal/logic/upload/upload_logic.go`
- Modify: `admin/internal/handler/upload_handler.go`

- [ ] **Step 1: 创建 upload_logic.go 业务逻辑**

```go
package upload

import (
	"context"
	"mime/multipart"
	"path/filepath"
	"strings"

	"admin/internal/infrastructure/storage"
	"admin/internal/svc"
	"admin/internal/types"
	"pkg/code"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	maxFileSize       = 5 * 1024 * 1024 // 5MB
	allowedExtensions = ".jpg,.jpeg,.png,.gif,.webp"
)

type UploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	storage storage.Storage
}

func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogic {
	// 从 svcCtx 获取 storage 实例
	return &UploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		storage: svcCtx.Storage,
	}
}

func (l *UploadLogic) Upload(req *types.UploadRequest) (*types.UploadResponse, error) {
	// 验证文件
	if req.File == nil {
		return nil, code.ErrUploadFailed
	}

	// 验证文件大小
	if req.File.Size > maxFileSize {
		return nil, code.ErrUploadFileSizeExceeded
	}

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(req.File.Filename))
	if !strings.Contains(allowedExtensions, ext) {
		return nil, code.ErrUploadUnsupportedFileType
	}

	// 验证 category
	category := storage.Category(req.Category)
	if !isValidCategory(category) {
		return nil, code.ErrUploadInvalidCategory
	}

	// 保存文件
	fileInfo, err := l.storage.Save(l.ctx, req.File, category)
	if err != nil {
		return nil, code.ErrUploadFailed
	}

	// 获取访问 URL
	url, err := l.storage.GetURL(l.ctx, fileInfo.ID)
	if err != nil {
		return nil, code.ErrUploadFailed
	}

	return &types.UploadResponse{
		ID:        fileInfo.ID,
		URL:       url,
		Filename:  fileInfo.Name,
		Category:  string(fileInfo.Category),
		Size:      fileInfo.Size,
		MimeType:  fileInfo.MimeType,
		Width:     fileInfo.Width,
		Height:    fileInfo.Height,
		CreatedAt: fileInfo.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (l *UploadLogic) Delete(req *types.DeleteUploadReq) error {
	err := l.storage.Delete(l.ctx, req.ID)
	if err != nil {
		return code.ErrUploadNotFound
	}
	return nil
}

func isValidCategory(category storage.Category) bool {
	switch category {
	case storage.CategoryProduct, storage.CategoryBanner, storage.CategoryAvatar:
		return true
	default:
		// 允许自定义 category
		return category != ""
	}
}
```

- [ ] **Step 2: 修改 handler 调用逻辑**

生成的 handler 代码结构如下，需要修改调用 UploadLogic：

```go
// admin/internal/handler/upload_handler.go
// 生成的代码大致如下，根据实际情况修改

func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析 form 数据
		file, header, err := r.FormFile("file")
		if err != nil {
			// 处理错误
		}
		defer file.Close()

		category := r.FormValue("category")

		req := &types.UploadRequest{
			File:     header, // 需要转换为 multipart.FileHeader
			Category: category,
		}

		l := upload.NewUploadLogic(r.Context(), svcCtx)
		resp, err := l.Upload(req)
		if err != nil {
			// 处理错误
		}

		// 返回响应
	}
}
```

- [ ] **Step 3: 修改 svcCtx 添加 Storage 依赖**

在 `admin/internal/svc/servicecontext.go` 中添加：

```go
type ServiceContext struct {
	Config       config.Config
	AdminUserSvc *adminuser.ServiceContext
	Storage      storage.Storage // 新增
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 storage
	storageCfg := storage.Config{
		Type: storage.StorageTypeLocal,
		Local: storage.LocalConfig{
			BasePath: "./uploads",
		},
	}
	// 需要传入 redis，可从 config 获取
	s, err := storage.NewStorage(storageCfg, nil)
	if err != nil {
		// 处理错误或使用 mock
	}

	return &ServiceContext{
		Config:       c,
		AdminUserSvc: adminuser.NewServiceContext(c),
		Storage:      s,
	}
}
```

- [ ] **Step 4: Commit**

```bash
git add admin/internal/logic/upload/ admin/internal/handler/upload_handler.go admin/internal/svc/
git commit -m "feat(upload): implement upload business logic"
```

---

## Task 4: 配置静态文件服务

**Files:**
- Modify: `admin/internal/config/config.go`
- Modify: `admin/etc/admin-api.yaml`

- [ ] **Step 1: 在 admin-api.yaml 中添加静态文件配置**

```yaml
Static:
  Root: ./uploads
  Prefix: /uploads
```

- [ ] **Step 2: 更新路由配置添加静态文件服务**

在 server 启动时添加静态文件中间件：

```go
// 在 main.go 中
server := rest.MustNewServer(c.Config.RestConf)
server.Use staticmiddleware.StaticsMiddleware(staticmiddleware.StaticsConf{
    Root: c.Config.Static.Root,
    Prefix: c.Config.Static.Prefix,
})
```

- [ ] **Step 3: Commit**

```bash
git add admin/
git commit -m "feat(upload): add static file serving for uploads"
```

---

## Task 5: 验证构建

- [ ] **Step 1: 运行 make build 验证编译**

```bash
cd admin && make build
```

- [ ] **Step 2: Commit**

```bash
git commit -m "chore: verify build passes"
```

---

## 验收清单

- [ ] API 接口定义完成
- [ ] 错误码已添加
- [ ] Storage 接口定义完成
- [ ] 本地存储实现完成
- [ ] 存储工厂创建完成
- [ ] 业务逻辑实现完成
- [ ] 静态文件服务配置完成
- [ ] 编译通过