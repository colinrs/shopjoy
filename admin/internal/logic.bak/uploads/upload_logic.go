package uploads

import (
	"context"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/infrastructure/storage"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	maxFileSize       = 5 * 1024 * 1024 // 5MB
	maxHeaderSize     = 512             // bytes to read for magic number detection
	allowedExtensions = ".jpg,.jpeg,.png,.gif,.webp"
)

// allowedMimes 允许的 MIME 类型及其对应的 magic bytes 签名
var allowedMimes = map[string][]byte{
	"image/jpeg": {0xFF, 0xD8, 0xFF},
	"image/png":  {0x89, 0x50, 0x4E, 0x47},
	"image/gif":  {0x47, 0x49, 0x46, 0x38},
	"image/webp": {0x57, 0x45, 0x42, 0x50},
}

// magicSignatures 危险文件的 magic bytes (可执行文件等)
var dangerousSignatures = [][]byte{
	{0x4D, 0x5A},                   // Windows EXE
	{0x7F, 0x45, 0x4C, 0x46},       // Linux ELF
	{0x50, 0x4B, 0x03, 0x04},       // ZIP (可能包含恶意脚本)
	{0xCA, 0xFE, 0xBA, 0xBE},       // macOS Mach-O
	{0x23, 0x21},                   // Shell script (#!)
	{0x3C, 0x3F, 0x70, 0x68, 0x70}, // PHP <?php
}

type UploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadLogic) Upload(req *types.UploadRequest, file multipart.File, header *multipart.FileHeader) (resp *types.UploadResponse, err error) {
	// 验证文件
	if file == nil {
		return nil, code.ErrUploadFailed
	}

	// 验证文件大小
	if header.Size > maxFileSize {
		return nil, code.ErrUploadFileSizeExceeded
	}

	// 验证文件扩展名
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !strings.Contains(allowedExtensions, ext) {
		return nil, code.ErrUploadUnsupportedFileType
	}

	// 读取文件头部用于 magic byte 验证
	headerBytes := make([]byte, maxHeaderSize)
	n, err := file.Read(headerBytes)
	if err != nil && err != io.EOF {
		return nil, code.ErrUploadFailed
	}
	headerBytes = headerBytes[:n]

	// 验证 magic bytes
	mimeType := detectMimeType(headerBytes)
	if mimeType == "" {
		return nil, code.ErrUploadUnsupportedFileType
	}

	// 验证是否为危险文件类型
	if isDangerousFile(headerBytes) {
		l.Logger.Errorf("dangerous file upload attempt: %s", header.Filename)
		return nil, code.ErrUploadUnsupportedFileType
	}

	// 将文件指针重置到开头
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, code.ErrUploadFailed
	}

	// 验证 category
	category := storage.Category(req.Category)
	if !isValidCategory(category) {
		return nil, code.ErrUploadInvalidCategory
	}

	// 保存文件
	fileInfo, err := l.svcCtx.Storage.Save(l.ctx, header, category)
	if err != nil {
		return nil, code.ErrUploadFailed
	}

	// 获取访问 URL
	url, err := l.svcCtx.Storage.GetURL(l.ctx, fileInfo.ID)
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
		CreatedAt: fileInfo.CreatedAt.Format(time.RFC3339),
	}, nil
}

// detectMimeType 通过 magic bytes 检测文件类型
func detectMimeType(data []byte) string {
	for mime, signature := range allowedMimes {
		if hasPrefix(data, signature) {
			return mime
		}
	}
	return ""
}

// isDangerousFile 检测危险文件类型
func isDangerousFile(data []byte) bool {
	for _, sig := range dangerousSignatures {
		if hasPrefix(data, sig) {
			return true
		}
	}
	return false
}

// hasPrefix 检查数据是否以给定签名开头
func hasPrefix(data []byte, prefix []byte) bool {
	if len(data) < len(prefix) {
		return false
	}
	for i := range prefix {
		if data[i] != prefix[i] {
			return false
		}
	}
	return true
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
