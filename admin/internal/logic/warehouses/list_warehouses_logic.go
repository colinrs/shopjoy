package warehouses

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListWarehousesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListWarehousesLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListWarehousesLogic {
	return ListWarehousesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListWarehousesLogic) ListWarehouses(req *types.ListWarehouseReq) (resp *types.ListWarehouseResp, err error) {

	// Resolve TenantID from ctx (injected by AuthMiddleware). A missing or zero
	// TenantID must not fall back to listing every tenant's warehouses, so we
	// reject it up front rather than returning cross-tenant rows.
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok || tenantID == 0 {
		return nil, code.ErrTenantNotFound
	}

	warehouses, err := l.svcCtx.WarehouseRepo.FindByTenant(l.ctx, l.svcCtx.DB, tenantID)
	if err != nil {
		return nil, err
	}

	list := make([]*types.WarehouseDetailResp, 0, len(warehouses))
	for _, w := range warehouses {
		list = append(list, toWarehouseDetailResp(w))
	}

	return &types.ListWarehouseResp{
		List: list,
	}, nil
}
