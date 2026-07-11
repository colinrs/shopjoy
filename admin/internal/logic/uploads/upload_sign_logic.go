// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package uploads

import (
	"context"
	"strconv"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/infrastructure/storage"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadSignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 申请 Cloudinary 签名（前端直传第一步）
func NewUploadSignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadSignLogic {
	return &UploadSignLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadSignLogic) UploadSign(req *types.UploadSignRequest) (resp *types.UploadSignResponse, err error) {
	// Default category when caller did not specify one.
	cat := storage.Category(req.Category)
	if cat == "" {
		cat = storage.CategoryProduct
	}

	// Resolve tenant/user from auth context. Reject only when there is no
	// authenticated user at all (forged request); platform admins legitimately
	// have tenantID=0 and must still be able to sign.
	tenantID, _ := contextx.GetTenantID(l.ctx)
	userID, _ := contextx.GetUserID(l.ctx)
	if userID == 0 {
		l.Logger.Errorf("upload sign rejected: missing auth context (user_id=%d, tenant_id=%d)", userID, tenantID)
		return nil, code.ErrUploadCrossTenantAccess
	}

	// Storage may not be configured (e.g., local dev with localStorage only),
	// in which case signing is not supported.
	signer, ok := l.svcCtx.Storage.(storage.Signer)
	if !ok || signer == nil {
		return nil, code.ErrUploadProviderError
	}

	sig, err := signer.Sign(l.ctx, storage.SignParams{
		Category:  cat,
		TenantID:  tenantID,
		Filename:  req.Filename,
		Timestamp: time.Now().UTC().Unix(),
	})
	if err != nil {
		l.Logger.Errorf("upload sign failed: user_id=%d tenant_id=%d category=%s err=%v",
			userID, tenantID, cat, err)
		return nil, code.ErrUploadSignFailed
	}

	// Pre-allocate asset ID so the browser can include it on the confirm step.
	assetID, err := l.svcCtx.IDGen.NextID(l.ctx)
	if err != nil {
		l.Logger.Errorf("upload sign: generate asset id failed: user_id=%d err=%v", userID, err)
		return nil, code.ErrUploadSignFailed
	}

	return &types.UploadSignResponse{
		CloudName:    sig.CloudName,
		APIKey:       sig.APIKey,
		Timestamp:    sig.Timestamp,
		Signature:    sig.Signature,
		Folder:       sig.Folder,
		PublicID:     sig.PublicID,
		UploadPreset: sig.UploadPreset,
		AssetID:      strconv.FormatInt(assetID, 10),
		UploadURL:    "https://api.cloudinary.com/v1_1/" + sig.CloudName + "/image/upload",
	}, nil
}