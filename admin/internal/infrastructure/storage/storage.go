// admin/internal/infrastructure/storage/storage.go
package storage

import (
	"context"
	"io"
	"mime/multipart"
	"time"
)

// Category 图片分类。
type Category string

const (
	CategoryProduct Category = "product"
	CategoryBanner  Category = "banner"
	CategoryAvatar  Category = "avatar"
)

// Asset backend-managed asset metadata.
type Asset struct {
	ID        string
	PublicID  string
	URL       string
	Filename  string
	Size      int64
	MimeType  string
	Width     int
	Height    int
	Format    string
	Category  Category
	Provider  string // "local" | "cloudinary"
	TenantID  int64
	CreatedBy int64
	CreatedAt time.Time
}

// AssetDraft Save() input — caller provides reader.
type AssetDraft struct {
	Filename  string
	Reader    io.Reader
	MimeType  string
	Category  Category
	TenantID  int64
	CreatedBy int64
}

// RemoteAsset frontend finished direct upload and reports back.
type RemoteAsset struct {
	PublicID  string
	URL       string
	Filename  string
	Size      int64
	MimeType  string
	Width     int
	Height    int
	Format    string
	Category  Category
	TenantID  int64
	CreatedBy int64
}

// Storage 资源存储接口。
type Storage interface {
	Save(ctx context.Context, draft AssetDraft) (*Asset, error)
	RegisterAsset(ctx context.Context, remote RemoteAsset) (*Asset, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*Asset, error)
	GetURL(ctx context.Context, id string) (string, error)
}

// Signer 远端签名能力（仅 cloudinaryStorage 实现）。
type Signer interface {
	Sign(ctx context.Context, p SignParams) (Signature, error)
}

// SignParams is the requested signature context.
type SignParams struct {
	Category  Category
	TenantID  int64
	Filename  string
	Timestamp int64
}

// Signature is the public-safe bundle returned to the browser.
type Signature struct {
	CloudName    string `json:"cloud_name"`
	APIKey       string `json:"api_key"`
	Timestamp    string `json:"timestamp"`
	Signature    string `json:"signature"`
	Folder       string `json:"folder"`
	PublicID     string `json:"public_id"`
	UploadPreset string `json:"upload_preset,omitempty"`
}

// FileHeaderCompat is a small adapter so caller-side logic keeps working
// with multipart headers without forcing storage to depend on multipart package elsewhere.
//
//	The existing UploadLogic reads the file and detects MIME; it then opens
//	the file and passes through to Storage.Save with a reader. The compat
//	helper below turns a *multipart.FileHeader into an AssetDraft quickly.
func FromMultipart(file *multipart.FileHeader, category Category, tenantID, createdBy int64) (AssetDraft, io.ReadCloser, error) {
	src, err := file.Open()
	if err != nil {
		return AssetDraft{}, nil, err
	}
	return AssetDraft{
		Filename:  file.Filename,
		Reader:    src,
		MimeType:  file.Header.Get("Content-Type"),
		Category:  category,
		TenantID:  tenantID,
		CreatedBy: createdBy,
	}, src, nil
}
