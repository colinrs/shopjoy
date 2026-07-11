package warehouses

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
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

	// Find warehouse
	warehouse, err := l.svcCtx.WarehouseRepo.FindByID(l.ctx, l.svcCtx.DB, req.ID)
	if err != nil {
		return nil, err
	}
	if warehouse == nil {
		return nil, code.ErrInventoryWarehouseNotFound
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
