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