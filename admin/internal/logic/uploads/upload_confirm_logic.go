// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package uploads

import (
	"context"
	"strings"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/infrastructure/storage"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadConfirmLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 前端直传后回调确认入库
func NewUploadConfirmLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadConfirmLogic {
	return &UploadConfirmLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// allowedImageMIME is the set of MIME types accepted on the confirm step.
// The browser-uploaded asset must be a still image.
var allowedImageMIME = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
	"image/gif":  {},
	"image/webp": {},
}

func (l *UploadConfirmLogic) UploadConfirm(req *types.UploadConfirmRequest) (resp *types.UploadConfirmResponse, err error) {
	// 1. MIME whitelist — we only register still images, never arbitrary blobs.
	if _, ok := allowedImageMIME[req.MimeType]; !ok {
		return nil, code.ErrUploadUnsupportedFileType
	}
	// 2. Reject empty identifiers — a successful Cloudinary upload always returns
	//    both a public_id and a URL, so an empty value means the browser is
	//    mis-reporting or tampering.
	if req.PublicID == "" || req.URL == "" {
		return nil, code.ErrUploadConfirmFailed
	}
	// 3. Path-traversal guard — public_id is later used as a folder key; reject
	//    anything that escapes its intended bucket.
	if strings.Contains(req.PublicID, "..") || strings.Contains(req.PublicID, `\`) {
		return nil, code.ErrUploadConfirmFailed
	}

	// 4. Resolve tenant/user from auth context. Platform admins may have
	//    tenantID=0, which downstream storage adapters treat as the
	//    "platform" bucket.
	tenantID, _ := contextx.GetTenantID(l.ctx)
	userID, _ := contextx.GetUserID(l.ctx)

	asset, err := l.svcCtx.Storage.RegisterAsset(l.ctx, storage.RemoteAsset{
		PublicID:  req.PublicID,
		URL:       req.URL,
		Filename:  req.Filename,
		Size:      req.Size,
		MimeType:  req.MimeType,
		Width:     req.Width,
		Height:    req.Height,
		Format:    req.Format,
		Category:  storage.Category(req.Category),
		TenantID:  tenantID,
		CreatedBy: userID,
	})
	if err != nil {
		l.Logger.Errorf("upload confirm failed: user_id=%d tenant_id=%d public_id=%s err=%v",
			userID, tenantID, req.PublicID, err)
		return nil, code.ErrUploadConfirmFailed
	}

	return &types.UploadConfirmResponse{
		ID:        asset.ID,
		URL:       asset.URL,
		Filename:  asset.Filename,
		Category:  string(asset.Category),
		Size:      asset.Size,
		MimeType:  asset.MimeType,
		Width:     asset.Width,
		Height:    asset.Height,
		CreatedAt: asset.CreatedAt.Format(time.RFC3339),
	}, nil
}
