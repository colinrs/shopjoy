package products

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSKULogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSKULogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteSKULogic {
	return DeleteSKULogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSKULogic) DeleteSKU(req *types.GetSKUReq) (resp *types.CreateSKUResp, err error) {
	// Get tenant ID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Delete SKU
	if err := l.svcCtx.SKURepo.Delete(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID); err != nil {
		return nil, err
	}

	return &types.CreateSKUResp{
		ID: req.ID,
	}, nil
}
