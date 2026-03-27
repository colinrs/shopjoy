package inventory

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSKUStockLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSKUStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateSKUStockLogic {
	return UpdateSKUStockLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSKUStockLogic) UpdateSKUStock(req *types.UpdateSKUStockReq) (resp *types.CreateWarehouseResp, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)
	userID, _ := contextx.GetUserID(l.ctx)

	// If warehouse ID is specified, update that warehouse's inventory
	if req.WarehouseID > 0 {
		// Find existing inventory record
		inventory, err := l.svcCtx.WarehouseInventoryRepo.FindBySKUAndWarehouse(
			l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.SKUCode, req.WarehouseID,
		)
		if err != nil {
			return nil, err
		}

		if inventory == nil {
			// Create new inventory record
			id, err := l.svcCtx.IDGen.NextID(l.ctx)
			if err != nil {
				return nil, err
			}
			inventory = &product.WarehouseInventory{
				ID:             id,
				TenantID:       shared.TenantID(tenantID),
				SKUCode:        req.SKUCode,
				WarehouseID:    req.WarehouseID,
				AvailableStock: req.AvailableStock,
				LockedStock:    0,
				Audit: shared.AuditInfo{
					CreatedAt: time.Now().UTC(),
					UpdatedAt: time.Now().UTC(),
				},
			}
			if err := l.svcCtx.WarehouseInventoryRepo.Create(l.ctx, l.svcCtx.DB, inventory); err != nil {
				return nil, err
			}
		} else {
			// Update existing inventory
			inventory.AvailableStock = req.AvailableStock
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
			ChangeType:     product.InventoryChangeManual,
			ChangeQuantity: req.AvailableStock,
			BeforeStock:    0, // Would need to track
			AfterStock:     req.AvailableStock,
			Remark:         req.Remark,
			OperatorID:     userID,
		}
		if err := l.svcCtx.InventoryLogRepo.Create(l.ctx, l.svcCtx.DB, log); err != nil {
			return nil, err
		}
	}

	return &types.CreateWarehouseResp{
		ID: 0,
	}, nil
}
