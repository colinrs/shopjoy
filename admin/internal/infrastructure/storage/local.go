package storage

import (
	"context"
	"fmt"
	"image"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	maxFileSize       = 5 * 1024 * 1024 // 5MB
	allowedExtensions = ".jpg,.jpeg,.png,.gif,.webp"
	uploadDir         = "./uploads"
)

type localStorage struct {
	basePath string
}

// NewLocalStorage 创建本地存储实例
func NewLocalStorage() Storage {
	return &localStorage{
		basePath: uploadDir,
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
	now := time.Now().UTC()
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
	if ext != ".gif" && ext != ".webp" {
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

	return fileInfo, nil
}

func (s *localStorage) Delete(ctx context.Context, id string) error {
	// 尝试删除常见目录下的文件
	now := time.Now().UTC()
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
	now := time.Now().UTC()
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
	url, err := s.GetURL(ctx, id)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		ID:        id,
		Path:      url,
		CreatedAt: time.Now().UTC(),
	}, nil
}
