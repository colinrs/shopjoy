package shipping_templates

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateShippingTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateShippingTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateShippingTemplateLogic {
	return UpdateShippingTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateShippingTemplateLogic) UpdateShippingTemplate(req *types.UpdateShippingTemplateReq) (resp *types.ShippingTemplateDetailResp, err error) {

	// Find existing template
	template, err := l.svcCtx.ShippingRepo.FindByID(l.ctx, l.svcCtx.DB, req.ID)
	if err != nil {
		return nil, err
	}

	// ─── wire → entity field map (anti-silent-drop guard) ───
	// Every UpdateShippingTemplateReq field must be conditionally assigned below.
	//   wire.Name         → entity.Name        (only if non-empty)
	//   wire.IsActive     → entity.IsActive    (only if pointer != nil)
	//   wire.CarrierCode  → entity.CarrierCode (only if non-empty)
	//   wire.WarehouseID  → entity.WarehouseID (only if non-zero)
	if req.Name != "" {
		template.Name = req.Name
	}
	if req.IsActive != nil {
		template.IsActive = *req.IsActive
	}
	if req.CarrierCode != "" {
		template.CarrierCode = req.CarrierCode
	}
	if req.WarehouseID != 0 {
		template.WarehouseID = req.WarehouseID
	}

	// Save changes
	if err := l.svcCtx.ShippingRepo.Update(l.ctx, l.svcCtx.DB, template); err != nil {
		return nil, err
	}

	// Get zones and mappings for response
	template, zones, mappings, err := l.svcCtx.ShippingRepo.FindByIDWithDetails(l.ctx, l.svcCtx.DB, req.ID)
	if err != nil {
		return nil, err
	}

	// ─── entity → response field map ───
	// All response fields on ShippingTemplateDetailResp must be populated below.
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
