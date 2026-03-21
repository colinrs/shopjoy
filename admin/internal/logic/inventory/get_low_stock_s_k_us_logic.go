package inventory

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetLowStockSKUsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLowStockSKUsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetLowStockSKUsLogic {
	return GetLowStockSKUsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLowStockSKUsLogic) GetLowStockSKUs(req *types.GetLowStockSKUsReq) (resp *types.ListLowStockSKUsResp, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Get low stock SKUs
	skus, err := l.svcCtx.WarehouseInventoryRepo.FindBySKU(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), "")
	if err != nil {
		return nil, err
	}

	// Filter and build response
	list := make([]*types.LowStockSKUResp, 0)
	for _, sku := range skus {
		// Check if low stock (would need safety_stock from SKU table)
		// For now, just return all SKUs with stock < 10 as example
		if sku.AvailableStock < 10 {
			list = append(list, &types.LowStockSKUResp{
				SKUCode:        sku.SKUCode,
				ProductID:      0,  // Would need to fetch from SKU
				ProductName:    "", // Would need to fetch from product
				AvailableStock: sku.AvailableStock,
				SafetyStock:    10, // Default safety stock
			})
		}
	}

	// Apply pagination
	total := int64(len(list))
	start := (req.Page - 1) * req.PageSize
	end := start + req.PageSize
	if start > len(list) {
		list = []*types.LowStockSKUResp{}
	} else if end > len(list) {
		list = list[start:]
	} else {
		list = list[start:end]
	}

	return &types.ListLowStockSKUsResp{
		List:  list,
		Total: total,
	}, nil
}
