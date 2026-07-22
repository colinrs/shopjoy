package warehouses

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
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

	// Resolve TenantID from ctx. UpdateWarehouse never rewrites TenantID —
	// the row's tenant_id stays whatever it was when it was created. We only
	// use the caller's TenantID here to enforce isolation: a tenant must not
	// be able to update another tenant's warehouse.
	//
	// FIX-2: a missing or zero TenantID is no longer rejected here. Platform
	// admins must be able to manage warehouses they own; the GORM tenant
	// middleware is responsible for the broader platform-scope invariant.
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Find existing warehouse
	warehouse, err := l.svcCtx.WarehouseRepo.FindByID(l.ctx, l.svcCtx.DB, req.ID)
	if err != nil {
		return nil, err
	}
	if warehouse == nil {
		return nil, code.ErrInventoryWarehouseNotFound
	}

	// Defense in depth: refuse cross-tenant updates from ordinary users.
	// FIX-2: platform admins (tenantID == 0) bypass this check so they can
	// update any warehouse; the GORM tenant middleware enforces the broader
	// platform-scope invariant.
	if tenantID != 0 && int64(warehouse.TenantID) != tenantID {
		return nil, code.ErrInventoryWarehouseNotFound
	}

	// Update fields
	warehouse.Name = req.Name
	warehouse.Country = req.Country
	warehouse.Address = req.Address
	warehouse.IsDefault = req.IsDefault
	warehouse.Model.UpdatedAt = time.Now().UTC()

	if err := l.svcCtx.WarehouseRepo.Update(l.ctx, l.svcCtx.DB, warehouse); err != nil {
		return nil, err
	}

	return toWarehouseDetailResp(warehouse), nil
}
