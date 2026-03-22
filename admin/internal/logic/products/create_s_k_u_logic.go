package products

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
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

	// Create SKU entity
	sku := &product.SKU{
		ID:             id,
		TenantID:       shared.TenantID(tenantID),
		ProductID:      req.ProductID,
		Code:           req.Code,
		Price:          shared.Money{Amount: req.Price, Currency: currency},
		Stock:          req.Stock,
		AvailableStock: req.Stock,
		LockedStock:    0,
		SafetyStock:    req.SafetyStock,
		PreSaleEnabled: req.PreSaleEnabled,
		Attributes:     req.Attributes,
		Status:         shared.StatusEnabled,
		Audit: shared.AuditInfo{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if err := l.svcCtx.SKURepo.Create(l.ctx, l.svcCtx.DB, sku); err != nil {
		return nil, err
	}

	return &types.CreateSKUResp{
		ID: id,
	}, nil
}
