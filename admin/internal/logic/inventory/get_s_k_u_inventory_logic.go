package inventory

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetSKUInventoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSKUInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetSKUInventoryLogic {
	return GetSKUInventoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSKUInventoryLogic) GetSKUInventory(req *types.GetSKUInventoryReq) (resp *types.SKUInventoryResp, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Get warehouse inventory for the SKU
	warehouseInventories, err := l.svcCtx.WarehouseInventoryRepo.FindBySKU(
		l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.SKUCode,
	)
	if err != nil {
		return nil, err
	}

	// Calculate totals and build warehouse list
	totalAvailable := 0
	totalLocked := 0
	warehouseItems := make([]*types.WarehouseInventoryItemResp, 0, len(warehouseInventories))

	for _, wi := range warehouseInventories {
		totalAvailable += wi.AvailableStock
		totalLocked += wi.LockedStock

		// Get warehouse name
		warehouse, _ := l.svcCtx.WarehouseRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), wi.WarehouseID)
		warehouseName := ""
		if warehouse != nil {
			warehouseName = warehouse.Name
		}

		warehouseItems = append(warehouseItems, &types.WarehouseInventoryItemResp{
			WarehouseID:    wi.WarehouseID,
			WarehouseName:  warehouseName,
			AvailableStock: wi.AvailableStock,
			LockedStock:    wi.LockedStock,
		})
	}

	return &types.SKUInventoryResp{
		SKUCode:        req.SKUCode,
		ProductID:      0, // Would need to fetch from SKU table
		TotalStock:     totalAvailable + totalLocked,
		AvailableStock: totalAvailable,
		LockedStock:    totalLocked,
		SafetyStock:    0,     // Would need to fetch from SKU table
		IsLowStock:     false, // Would need to compare with safety stock
		Warehouses:     warehouseItems,
	}, nil
}
