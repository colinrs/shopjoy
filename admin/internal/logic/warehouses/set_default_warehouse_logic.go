package warehouses

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type SetDefaultWarehouseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetDefaultWarehouseLogic(ctx context.Context, svcCtx *svc.ServiceContext) SetDefaultWarehouseLogic {
	return SetDefaultWarehouseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetDefaultWarehouseLogic) SetDefaultWarehouse(req *types.SetDefaultWarehouseReq) (resp *types.WarehouseDetailResp, err error) {

	// Resolve and authorize: rejects missing TenantID and cross-tenant access.
	warehouse, err := requireWarehouseOwnership(l.ctx, l.svcCtx, req.ID)
	if err != nil {
		return nil, err
	}

	// Already default
	if warehouse.IsDefault {
		return toWarehouseDetailResp(warehouse), nil
	}

	// Set as default
	warehouse.IsDefault = true
	warehouse.Model.UpdatedAt = time.Now().UTC()

	if err := l.svcCtx.WarehouseRepo.Update(l.ctx, l.svcCtx.DB, warehouse); err != nil {
		return nil, err
	}

	return toWarehouseDetailResp(warehouse), nil
}
