package shipping_templates

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/domain/shipping"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateShippingTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateShippingTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateShippingTemplateLogic {
	return CreateShippingTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateShippingTemplateLogic) CreateShippingTemplate(req *types.CreateShippingTemplateReq) (resp *types.CreateShippingTemplateResp, err error) {
	// Validate
	if req.Name == "" {
		return nil, code.ErrShippingTemplateNameRequired
	}
	if req.Currency != "" && !isValidCurrency(req.Currency) {
		return nil, code.ErrShippingTemplateInvalidCurrency
	}

	// Resolve tenantID from context (REQUIRED for multi-tenancy)
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// If a warehouse is specified, validate ownership: must exist and belong to the same tenant.
	// Warehouse has no MarketID field, so market-level validation is skipped here.
	if req.WarehouseID > 0 {
		warehouse, err := l.svcCtx.WarehouseRepo.FindByID(l.ctx, l.svcCtx.DB, req.WarehouseID)
		if err != nil {
			return nil, err
		}
		if warehouse == nil || warehouse.TenantID.Int64() != tenantID {
			return nil, code.ErrShippingTemplateWarehouseMismatch
		}
	}

	// Create template entity.
	// ─── wire → entity field map (anti-silent-drop guard) ───
	//   wire.Name         → entity.Name        (required)
	//   wire.IsDefault    → entity.IsDefault   (bool, optional)
	//   wire.MarketID     → entity.MarketID    (int64, default 0 = 全市场通用)
	//   wire.Currency     → entity.Currency    (default "CNY")
	//   wire.CarrierCode  → entity.CarrierCode (default "standard")
	//   wire.WarehouseID  → entity.WarehouseID (int64, default 0)
	//   hardcoded         → entity.IsActive    (true)
	//   from ctx          → entity.TenantID
	template := &shipping.ShippingTemplate{
		TenantID:    tenantID,
		MarketID:    req.MarketID,
		Currency:    defaultCurrency(req.Currency),
		CarrierCode: defaultCarrierCode(req.CarrierCode),
		WarehouseID: req.WarehouseID,
		Name:        req.Name,
		IsDefault:   req.IsDefault,
		IsActive:    true,
	}

	// If setting as default, unset other defaults in the same (tenant, market) first.
	// This enforces the (tenant_id, market_id, is_default=true) unique partial index.
	if req.IsDefault {
		if err := l.svcCtx.ShippingRepo.UnsetAllDefaultByMarket(l.ctx, l.svcCtx.DB, tenantID, req.MarketID); err != nil {
			return nil, err
		}
	}

	// Save template
	if err := l.svcCtx.ShippingRepo.Create(l.ctx, l.svcCtx.DB, template); err != nil {
		return nil, err
	}

	return &types.CreateShippingTemplateResp{
		ID:   int64(template.ID),
		Name: template.Name,
	}, nil
}

// isValidCurrency checks whether c is one of the supported ISO 4217 currency
// codes the platform accepts for shipping templates.
func isValidCurrency(c string) bool {
	switch c {
	case "CNY", "USD", "EUR", "GBP", "JPY", "KRW", "SGD", "MYR", "THB", "IDR", "PHP", "VND":
		return true
	}
	return false
}

// defaultCurrency returns the wire-supplied currency or "CNY" if empty.
func defaultCurrency(c string) string {
	if c == "" {
		return "CNY"
	}
	return c
}

// defaultCarrierCode returns the wire-supplied carrier code or "standard".
func defaultCarrierCode(c string) string {
	if c == "" {
		return "standard"
	}
	return c
}
