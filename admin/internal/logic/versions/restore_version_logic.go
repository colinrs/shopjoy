package versions

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type RestoreVersionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRestoreVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) RestoreVersionLogic {
	return RestoreVersionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RestoreVersionLogic) RestoreVersion(req *types.RestoreVersionRequest) error {
	tenantID, _ := contextx.GetTenantID(l.ctx)
	userID, _ := contextx.GetUserID(l.ctx)

	return l.svcCtx.VersionService.RestoreVersion(l.ctx, shared.TenantID(tenantID), req.PageID, req.Version, userID)
}
