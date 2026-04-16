package inventory

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
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
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok && !contextx.IsPlatformAdmin(l.ctx) {
		return nil, code.ErrUnauthorized
	}
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Find low stock SKUs using SafetyStock threshold
	skus, total, err := l.svcCtx.SKURepo.FindLowStock(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	if len(skus) == 0 {
		return &types.ListLowStockSKUsResp{
			List:  []*types.LowStockSKUResp{},
			Total: 0,
		}, nil
	}

	// Collect product IDs
	productIDs := make([]int64, 0, len(skus))
	productIDSet := make(map[int64]struct{})
	for _, sku := range skus {
		if _, exists := productIDSet[sku.ProductID]; !exists {
			productIDs = append(productIDs, sku.ProductID)
			productIDSet[sku.ProductID] = struct{}{}
		}
	}

	// Fetch products for names
	products, err := l.svcCtx.ProductRepo.FindByIDs(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), productIDs)
	if err != nil {
		return nil, err
	}

	// Build product name map
	productNames := make(map[int64]string)
	for _, p := range products {
		productNames[int64(p.ID)] = p.Name
	}

	// Build response
	list := make([]*types.LowStockSKUResp, 0, len(skus))
	for _, sku := range skus {
		productName := productNames[sku.ProductID]
		if productName == "" {
			productName = "Unknown Product"
		}
		list = append(list, &types.LowStockSKUResp{
			SKUCode:        sku.Code,
			ProductID:      sku.ProductID,
			ProductName:    productName,
			AvailableStock: sku.AvailableStock,
			SafetyStock:    sku.SafetyStock,
		})
	}

	return &types.ListLowStockSKUsResp{
		List:  list,
		Total: total,
	}, nil
}
