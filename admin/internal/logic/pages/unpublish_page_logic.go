package pages

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type UnpublishPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnpublishPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) UnpublishPageLogic {
	return UnpublishPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnpublishPageLogic) UnpublishPage(req *types.UnpublishPageRequest) error {
	tenantID, _ := contextx.GetTenantID(l.ctx)

	return l.svcCtx.PageService.UnpublishPage(l.ctx, shared.TenantID(tenantID), req.ID)
}