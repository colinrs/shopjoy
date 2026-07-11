package versions

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
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
	userID := contextx.GetCurrentUserID(l.ctx)
	return l.svcCtx.VersionService.RestoreVersion(l.ctx, req.PageID, req.Version, userID)
}
