package warehouses

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetWarehouseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWarehouseLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetWarehouseLogic {
	return GetWarehouseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWarehouseLogic) GetWarehouse(req *types.GetWarehouseReq) (resp *types.WarehouseDetailResp, err error) {

	// Resolve and authorize: rejects missing TenantID and cross-tenant access.
	warehouse, err := requireWarehouseOwnership(l.ctx, l.svcCtx, req.ID)
	if err != nil {
		return nil, err
	}

	return toWarehouseDetailResp(warehouse), nil
}
