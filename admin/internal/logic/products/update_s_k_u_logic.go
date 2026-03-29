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

type UpdateSKULogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSKULogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateSKULogic {
	return UpdateSKULogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSKULogic) UpdateSKU(req *types.UpdateSKUReq) (resp *types.SKUDetailResp, err error) {
	// Get tenant ID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all tenant data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Find SKU
	sku, err := l.svcCtx.SKURepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Code != "" {
		sku.Code = req.Code
	}
	if req.Price != "" {
		currency := req.Currency
		if currency == "" {
			currency = sku.Price.Currency
		}
		// Parse price from string (yuan)
		priceAmount, err := shared.ParseMoneyFromString(req.Price)
		if err != nil {
			return nil, err
		}
		sku.Price = shared.NewMoney(priceAmount, currency)
	}
	if req.Stock > 0 {
		diff := req.Stock - sku.Stock
		sku.Stock = req.Stock
		sku.AvailableStock += diff
		if sku.AvailableStock < 0 {
			sku.AvailableStock = 0
		}
	}
	if req.SafetyStock > 0 {
		sku.SafetyStock = req.SafetyStock
	}
	sku.PreSaleEnabled = req.PreSaleEnabled
	if req.Attributes != nil {
		sku.Attributes = req.Attributes
	}
	sku.Audit.UpdatedAt = time.Now().UTC()

	// Save
	if err := l.svcCtx.SKURepo.Update(l.ctx, l.svcCtx.DB, sku); err != nil {
		return nil, err
	}

	return &types.SKUDetailResp{
		ID:             sku.ID,
		ProductID:      sku.ProductID,
		Code:           sku.Code,
		Price:          shared.FormatMoneyToStringOnly(sku.Price.Amount),
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
	}, nil
}
