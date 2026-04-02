package pages

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteDecorationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteDecorationLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteDecorationLogic {
	return DeleteDecorationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteDecorationLogic) DeleteDecoration(req *types.DeleteDecorationRequest) error {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return err
	}

	return l.svcCtx.DecorationService.DeleteDecoration(l.ctx, shared.TenantID(tenantID), req.ID)
}
