package warehouses

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
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
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	warehouses, err := l.svcCtx.WarehouseRepo.FindAll(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID))
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
