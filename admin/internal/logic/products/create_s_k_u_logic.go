package products

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

type CreateSKULogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateSKULogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateSKULogic {
	return CreateSKULogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSKULogic) CreateSKU(req *types.CreateSKUReq) (resp *types.CreateSKUResp, err error) {
	// Get tenant ID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Generate SKU ID
	id, err := l.svcCtx.IDGen.NextID(l.ctx)
	if err != nil {
		return nil, err
	}

	// Set default currency
	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}

	// Determine SKU code: user-provided or auto-generated
	var skuCode string
	if req.Code != "" {
		// User provided code - use as-is for backward compatibility
		skuCode = req.Code
	} else {
		// Auto-generate SKU code
		// Fetch product to get its SKU prefix
		productSKUPrefix := ""
		prod, err := l.svcCtx.ProductRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ProductID)
		if err == nil && prod != nil {
			productSKUPrefix = prod.SKUPrefix
		}

		// Tenant SKU prefix is not available (no tenant repo in admin service)
		// Use empty prefix as fallback
		tenantSKUPrefix := ""

		skuCode, err = l.svcCtx.SKUGenerator.GenerateWithRetry(
			tenantID,
			tenantSKUPrefix,
			productSKUPrefix,
			3,
		)
		if err != nil {
			return nil, err
		}
	}

	// Create SKU entity
	// Parse price from string (yuan)
	priceAmount, err := shared.ParseMoneyFromString(req.Price)
	if err != nil {
		return nil, err
	}
	sku := &product.SKU{
		Model:          application.Model{ID: id, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
		TenantID:       shared.TenantID(tenantID),
		ProductID:      req.ProductID,
		Code:           skuCode,
		Price:          shared.NewMoney(priceAmount, currency),
		Stock:          req.Stock,
		AvailableStock: req.Stock,
		LockedStock:    0,
		SafetyStock:    req.SafetyStock,
		PreSaleEnabled: req.PreSaleEnabled,
		Attributes:     req.Attributes,
		Status:         shared.StatusEnabled,
	}

	if err := l.svcCtx.SKURepo.Create(l.ctx, l.svcCtx.DB, sku); err != nil {
		return nil, err
	}

	return &types.CreateSKUResp{
		ID: id,
	}, nil
}
