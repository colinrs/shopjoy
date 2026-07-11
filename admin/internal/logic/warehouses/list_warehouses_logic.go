package warehouses

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
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

	warehouses, err := l.svcCtx.WarehouseRepo.FindAll(l.ctx, l.svcCtx.DB)
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
