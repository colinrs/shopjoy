package shipping_templates

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"

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
	// Get tenant ID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)
	if tenantID == 0 {
		return nil, code.ErrUnauthorized
	}

	// Find existing template
	template, err := l.svcCtx.ShippingRepo.FindByID(l.ctx, l.svcCtx.DB, tenantID, req.ID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		template.Name = req.Name
	}
	if req.IsActive != nil {
		template.IsActive = *req.IsActive
	}

	// Save changes
	if err := l.svcCtx.ShippingRepo.Update(l.ctx, l.svcCtx.DB, template); err != nil {
		return nil, err
	}

	// Get zones and mappings for response
	template, zones, mappings, err := l.svcCtx.ShippingRepo.FindByIDWithDetails(l.ctx, l.svcCtx.DB, tenantID, req.ID)
	if err != nil {
		return nil, err
	}

	return &types.ShippingTemplateDetailResp{
		ID:        template.ID,
		Name:      template.Name,
		IsDefault: template.IsDefault,
		IsActive:  template.IsActive,
		Zones:     buildZoneDetails(zones),
		Mappings:  buildMappingDetails(mappings),
		CreatedAt: template.CreatedAt.Format(time.RFC3339),
		UpdatedAt: template.UpdatedAt.Format(time.RFC3339),
	}, nil
}