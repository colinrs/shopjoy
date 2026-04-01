package shipping_zones

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReorderZonesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReorderZonesLogic(ctx context.Context, svcCtx *svc.ServiceContext) ReorderZonesLogic {
	return ReorderZonesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReorderZonesLogic) ReorderZones(req *types.ReorderZonesReq) error {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return err
	}

	// Verify template exists and belongs to tenant
	template, err := l.svcCtx.ShippingRepo.FindByID(l.ctx, l.svcCtx.DB, tenantID, req.TemplateID)
	if err != nil {
		return err
	}
	if template == nil {
		return code.ErrShippingTemplateNotFound
	}

	// Reorder zones
	return l.svcCtx.ShippingRepo.ReorderZones(l.ctx, l.svcCtx.DB, req.TemplateID, req.ZoneIDs)
}