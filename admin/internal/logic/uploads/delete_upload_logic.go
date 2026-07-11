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
	// Look up the asset first so we can enforce tenant scoping.
	asset, err := l.svcCtx.Storage.Get(l.ctx, req.ID)
	if err != nil {
		return code.ErrUploadNotFound
	}

	// Resolve tenant from auth context; platform admins may have tenantID=0.
	tenantID, _ := contextx.GetTenantID(l.ctx)
	if tenantID != 0 && asset.TenantID != tenantID {
		l.Logger.Errorf("delete upload: cross-tenant access denied: request_tenant=%d asset_tenant=%d asset_id=%s",
			tenantID, asset.TenantID, req.ID)
		return code.ErrUploadCrossTenantAccess
	}

	if err := l.svcCtx.Storage.Delete(l.ctx, req.ID); err != nil {
		return code.ErrUploadNotFound
	}
	return nil
}
