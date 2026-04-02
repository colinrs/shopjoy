package shipping_mappings

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/domain/shipping"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTemplateMappingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateTemplateMappingLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateTemplateMappingLogic {
	return CreateTemplateMappingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTemplateMappingLogic) CreateTemplateMapping(req *types.CreateTemplateMappingReq) (resp *types.TemplateMappingDetail, err error) {
	// Get tenant ID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)
	if tenantID == 0 {
		return nil, code.ErrUnauthorized
	}

	// Verify template exists and belongs to tenant
	template, err := l.svcCtx.ShippingRepo.FindByID(l.ctx, l.svcCtx.DB, tenantID, req.TemplateID)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, code.ErrShippingTemplateNotFound
	}

	// Validate target type
	targetType := shipping.TargetType(req.TargetType)
	if !targetType.IsValid() {
		return nil, code.ErrShippingMappingInvalidTarget
	}

	// Check if mapping already exists
	existing, err := l.svcCtx.ShippingRepo.FindMappingByTarget(l.ctx, l.svcCtx.DB, targetType, req.TargetID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		// Update existing mapping to point to new template
		existing.TemplateID = req.TemplateID
		if err := l.svcCtx.ShippingRepo.UpdateMapping(l.ctx, l.svcCtx.DB, existing); err != nil {
			return nil, err
		}
		return &types.TemplateMappingDetail{
			ID:         int64(existing.ID),
			TemplateID: existing.TemplateID,
			TargetType: string(existing.TargetType),
			TargetID:   existing.TargetID,
		}, nil
	}

	// Create new mapping
	mapping := &shipping.ShippingTemplateMapping{
		TenantID:   tenantID,
		TemplateID: req.TemplateID,
		TargetType: targetType,
		TargetID:   req.TargetID,
	}

	// Validate mapping
	if err := mapping.Validate(); err != nil {
		return nil, err
	}

	// Save mapping
	if err := l.svcCtx.ShippingRepo.CreateMapping(l.ctx, l.svcCtx.DB, mapping); err != nil {
		return nil, err
	}

	return &types.TemplateMappingDetail{
		ID:         int64(mapping.ID),
		TemplateID: mapping.TemplateID,
		TargetType: string(mapping.TargetType),
		TargetID:   mapping.TargetID,
	}, nil
}
