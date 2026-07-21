package shipping_templates

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type SetDefaultTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetDefaultTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) SetDefaultTemplateLogic {
	return SetDefaultTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetDefaultTemplateLogic) SetDefaultTemplate(req *types.SetDefaultTemplateReq) (resp *types.ShippingTemplateDetailResp, err error) {

	// C3 fix: tenant scope required to avoid clearing another tenant's defaults.
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Verify template exists and belongs to tenant
	template, err := l.svcCtx.ShippingRepo.FindByID(l.ctx, l.svcCtx.DB, req.ID)
	if err != nil {
		return nil, err
	}

	// Use transaction to ensure atomicity
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// Unset all defaults within this template's (tenant, market) first.
		// (tenant_id, market_id, is_default=true) is a per-tenant per-market
		// partial unique index, so we must scope the unset to the same
		// tenant_id + market_id we're about to set.
		if err := l.svcCtx.ShippingRepo.UnsetAllDefaultByMarket(l.ctx, tx, tenantID, template.MarketID); err != nil {
			return err
		}

		// Set new default
		if err := l.svcCtx.ShippingRepo.SetDefault(l.ctx, tx, req.ID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Update the template object
	template.IsDefault = true

	// Get zones and mappings for response
	template, zones, mappings, err := l.svcCtx.ShippingRepo.FindByIDWithDetails(l.ctx, l.svcCtx.DB, req.ID)
	if err != nil {
		return nil, err
	}

	return &types.ShippingTemplateDetailResp{
		ID:          int64(template.ID),
		TenantID:    template.TenantID,
		MarketID:    template.MarketID,
		Currency:    template.Currency,
		CarrierCode: template.CarrierCode,
		WarehouseID: template.WarehouseID,
		Name:        template.Name,
		IsDefault:   template.IsDefault,
		IsActive:    template.IsActive,
		Zones:       buildZoneDetails(zones),
		Mappings:    buildMappingDetails(mappings),
		CreatedAt:   template.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   template.UpdatedAt.Format(time.RFC3339),
	}, nil
}
