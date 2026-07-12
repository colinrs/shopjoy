// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package uploads

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取图片元数据
func NewGetUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUploadLogic {
	return &GetUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUploadLogic) GetUpload(req *types.GetUploadReq) (resp *types.UploadResponse, err error) {
	// Note: req.ID is the internal snowflake primary key (string-encoded),
	// not the Cloudinary public_id. Storage.Get uses this as the lookup key.
	// We then check that the returned asset's tenant_id matches the auth
	// context to prevent IDOR across tenants.
	asset, err := l.svcCtx.Storage.Get(l.ctx, req.ID)
	if err != nil {
		l.Logger.Errorf("get upload: storage.Get failed: id=%s err=%v", req.ID, err)
		return nil, code.ErrUploadNotFound
	}

	// Cross-tenant guard. Platform admins may have tenantID=0, which is
	// allowed to access any tenant's asset.
	tenantID, _ := contextx.GetTenantID(l.ctx)
	if tenantID != 0 && asset.TenantID != tenantID {
		l.Logger.Errorf("get upload: cross-tenant access denied: request_tenant=%d asset_tenant=%d asset_id=%s",
			tenantID, asset.TenantID, req.ID)
		return nil, code.ErrUploadCrossTenantAccess
	}

	return &types.UploadResponse{
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
