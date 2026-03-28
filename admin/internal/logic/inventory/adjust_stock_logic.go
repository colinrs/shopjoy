package inventory

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type AdjustStockLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdjustStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) AdjustStockLogic {
	return AdjustStockLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdjustStockLogic) AdjustStock(req *types.AdjustStockReq) (resp *types.CreateWarehouseResp, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)
	userID, _ := contextx.GetUserID(l.ctx)

	// Verify warehouse exists
	warehouse, err := l.svcCtx.WarehouseRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.WarehouseID)
	if err != nil {
		return nil, err
	}
	if warehouse == nil {
		return nil, code.ErrInventoryWarehouseNotFound
	}

	// Find existing inventory record
	inventory, err := l.svcCtx.WarehouseInventoryRepo.FindBySKUAndWarehouse(
		l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.SKUCode, req.WarehouseID,
	)
	if err != nil {
		return nil, err
	}

	beforeStock := 0
	if inventory == nil {
		// Create new inventory record
		id, err := l.svcCtx.IDGen.NextID(l.ctx)
		if err != nil {
			return nil, err
		}
		newStock := req.Quantity
		if newStock < 0 {
			newStock = 0
		}
		inventory = &product.WarehouseInventory{
			Model:          application.Model{ID: id, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
			TenantID:       shared.TenantID(tenantID),
			SKUCode:        req.SKUCode,
			WarehouseID:    req.WarehouseID,
			AvailableStock: newStock,
			LockedStock:    0,
		}
		if err := l.svcCtx.WarehouseInventoryRepo.Create(l.ctx, l.svcCtx.DB, inventory); err != nil {
			return nil, err
		}
	} else {
		// Update existing inventory
		beforeStock = inventory.AvailableStock
		newStock := inventory.AvailableStock + req.Quantity
		if newStock < 0 {
			return nil, code.ErrInventoryInsufficientStock
		}
		inventory.AvailableStock = newStock
		inventory.Audit.UpdatedAt = time.Now().UTC()
		if err := l.svcCtx.WarehouseInventoryRepo.Update(l.ctx, l.svcCtx.DB, inventory); err != nil {
			return nil, err
		}
	}

	// Create inventory log
	logID, _ := l.svcCtx.IDGen.NextID(l.ctx)
	log := &product.InventoryLog{
		Model:          application.Model{ID: logID},
		TenantID:       shared.TenantID(tenantID),
		SKUCode:        req.SKUCode,
		ProductID:      0, // Would need to fetch from SKU
		WarehouseID:    req.WarehouseID,
		ChangeType:     product.InventoryChangeAdjustment,
		ChangeQuantity: req.Quantity,
		BeforeStock:    beforeStock,
		AfterStock:     inventory.AvailableStock,
		Remark:         req.Remark,
		OperatorID:     userID,
	}
	if err := l.svcCtx.InventoryLogRepo.Create(l.ctx, l.svcCtx.DB, log); err != nil {
		return nil, err
	}

	return &types.CreateWarehouseResp{
		ID: 0,
	}, nil
}
