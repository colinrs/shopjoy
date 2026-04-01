package shipping_templates

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteShippingTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteShippingTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteShippingTemplateLogic {
	return DeleteShippingTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteShippingTemplateLogic) DeleteShippingTemplate(req *types.DeleteShippingTemplateReq) error {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return err
	}

	// Find existing template
	template, err := l.svcCtx.ShippingRepo.FindByID(l.ctx, l.svcCtx.DB, tenantID, req.ID)
	if err != nil {
		return err
	}

	// Check if can delete
	if err := template.CanDelete(); err != nil {
		return err
	}

	// Delete template with cascade (zones and mappings) in transaction
	return l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// Delete all zones (which also deletes zone_regions)
		zones, err := l.svcCtx.ShippingRepo.FindZonesByTemplateID(l.ctx, tx, req.ID)
		if err != nil {
			return err
		}
		for _, zone := range zones {
			if err := l.svcCtx.ShippingRepo.DeleteZone(l.ctx, tx, int64(zone.ID)); err != nil {
				return err
			}
		}

		// Delete all mappings
		mappings, err := l.svcCtx.ShippingRepo.FindMappingsByTemplateID(l.ctx, tx, req.ID)
		if err != nil {
			return err
		}
		for _, mapping := range mappings {
			if err := l.svcCtx.ShippingRepo.DeleteMapping(l.ctx, tx, int64(mapping.ID)); err != nil {
				return err
			}
		}

		// Finally delete the template
		return l.svcCtx.ShippingRepo.Delete(l.ctx, tx, tenantID, req.ID)
	})
}