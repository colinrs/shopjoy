package warehouses

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
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

	// Resolve TenantID from ctx (injected by AuthMiddleware).
	//
	// FIX-2: a missing or zero TenantID no longer short-circuits with
	// code.ErrTenantNotFound — platform admins must be able to list
	// warehouses they own. The repo's FindByTenant still scopes by tenant_id,
	// so an ordinary user with a non-zero TenantID only ever sees their own
	// warehouses; a platform admin (TenantID == 0) gets whatever rows match
	// tenant_id = 0 today, and the GORM tenant middleware is responsible for
	// expanding that to platform-wide listings once it ships.
	tenantID, _ := contextx.GetTenantID(l.ctx)

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
