package uploads

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/colinrs/shopjoy/admin/internal/infrastructure/storage"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	maxFileSize       = 5 * 1024 * 1024 // 5MB
	allowedExtensions = ".jpg,.jpeg,.png,.gif,.webp"
)

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

func (l *UploadLogic) Upload(req *types.UploadRequest) (resp *types.UploadResponse, err error) {
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
	fileInfo, err := l.svcCtx.Storage.Save(l.ctx, req.File, category)
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
		CreatedAt: fileInfo.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
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