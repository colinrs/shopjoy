package warehouses

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateWarehouseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateWarehouseLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateWarehouseLogic {
	return UpdateWarehouseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateWarehouseLogic) UpdateWarehouse(req *types.UpdateWarehouseReq) (resp *types.WarehouseDetailResp, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	// Find existing warehouse
	warehouse, err := l.svcCtx.WarehouseRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}
	if warehouse == nil {
		return nil, code.ErrInventoryWarehouseNotFound
	}

	// Update fields
	warehouse.Name = req.Name
	warehouse.Country = req.Country
	warehouse.Address = req.Address
	warehouse.IsDefault = req.IsDefault
	warehouse.Audit.UpdatedAt = time.Now().UTC()

	if err := l.svcCtx.WarehouseRepo.Update(l.ctx, l.svcCtx.DB, warehouse); err != nil {
		return nil, err
	}

	return toWarehouseDetailResp(warehouse), nil
}
