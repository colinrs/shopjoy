package warehouses

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteWarehouseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteWarehouseLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteWarehouseLogic {
	return DeleteWarehouseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteWarehouseLogic) DeleteWarehouse(req *types.GetWarehouseReq) (resp *types.CreateWarehouseResp, err error) {

	// Find warehouse
	warehouse, err := l.svcCtx.WarehouseRepo.FindByID(l.ctx, l.svcCtx.DB, req.ID)
	if err != nil {
		return nil, err
	}
	if warehouse == nil {
		return nil, code.ErrInventoryWarehouseNotFound
	}

	// Cannot delete default warehouse
	if warehouse.IsDefault {
		return nil, code.ErrInventoryCannotDeleteDefault
	}

	// Delete warehouse
	if err := l.svcCtx.WarehouseRepo.Delete(l.ctx, l.svcCtx.DB, req.ID); err != nil {
		return nil, err
	}

	return &types.CreateWarehouseResp{
		ID: req.ID,
	}, nil
}
