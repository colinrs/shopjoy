package products

import (
	"context"
	"strconv"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListSKUsByProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListSKUsByProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListSKUsByProductLogic {
	return ListSKUsByProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListSKUsByProductLogic) ListSKUsByProduct(req *types.ListSKUsByProductReq) (resp *types.ListSKUsResp, err error) {
	// Get tenant ID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all tenant data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Find SKUs by product ID
	skus, err := l.svcCtx.SKURepo.FindByProductID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ProductID)
	if err != nil {
		return nil, err
	}

	// Convert to response
	list := make([]*types.SKUDetailResp, len(skus))
	for i, sku := range skus {
		list[i] = &types.SKUDetailResp{
			ID:             sku.ID,
			ProductID:      sku.ProductID,
			Code:           sku.Code,
			Price:          sku.Price.Amount,
			Currency:       sku.Price.Currency,
			Stock:          sku.Stock,
			AvailableStock: sku.AvailableStock,
			LockedStock:    sku.LockedStock,
			SafetyStock:    sku.SafetyStock,
			PreSaleEnabled: sku.PreSaleEnabled,
			Attributes:     sku.Attributes,
			Status:         strconv.Itoa(int(sku.Status)),
			IsLowStock:     sku.IsLowStock(),
			CreatedAt:      sku.Audit.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      sku.Audit.UpdatedAt.Format(time.RFC3339),
		}
	}

	return &types.ListSKUsResp{
		List:  list,
		Total: int64(len(list)),
	}, nil
}
