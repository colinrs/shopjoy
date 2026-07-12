package uploads

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUploadLogic {
	return &DeleteUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUploadLogic) DeleteUpload(req *types.DeleteUploadReq) error {
	// Resolve tenant from auth context; platform admins may have tenantID=0.
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Single-lookup delete: Storage.DeleteByTenant runs one DB UPDATE
	// `WHERE id = ? AND tenant_id = ?` and reports RowsAffected == 0 as
	// ErrMediaAssetNotFound (collapsing missing-row and cross-tenant into
	// one signal — no IDOR existence leak). No pre-fetch needed.
	if err := l.svcCtx.Storage.DeleteByTenant(l.ctx, req.ID, tenantID); err != nil {
		l.Logger.Errorf("delete upload: storage.DeleteByTenant failed: id=%s tenant_id=%d err=%v",
			req.ID, tenantID, err)
		return code.ErrUploadNotFound
	}
	return nil
}
