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
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	// Verify template exists and belongs to tenant
	template, err := l.svcCtx.ShippingRepo.FindByID(l.ctx, l.svcCtx.DB, tenantID, req.ID)
	if err != nil {
		return nil, err
	}

	// Use transaction to ensure atomicity
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// Unset all defaults first
		if err := l.svcCtx.ShippingRepo.UnsetAllDefault(l.ctx, tx, tenantID); err != nil {
			return err
		}

		// Set new default
		if err := l.svcCtx.ShippingRepo.SetDefault(l.ctx, tx, tenantID, req.ID); err != nil {
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
	template, zones, mappings, err := l.svcCtx.ShippingRepo.FindByIDWithDetails(l.ctx, l.svcCtx.DB, tenantID, req.ID)
	if err != nil {
		return nil, err
	}

	return &types.ShippingTemplateDetailResp{
		ID:        int64(template.ID),
		Name:      template.Name,
		IsDefault: template.IsDefault,
		IsActive:  template.IsActive,
		Zones:     buildZoneDetails(zones),
		Mappings:  buildMappingDetails(mappings),
		CreatedAt: template.CreatedAt.Format(time.RFC3339),
		UpdatedAt: template.UpdatedAt.Format(time.RFC3339),
	}, nil
}
