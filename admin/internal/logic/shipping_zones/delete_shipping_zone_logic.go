package shipping_zones

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteShippingZoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteShippingZoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteShippingZoneLogic {
	return DeleteShippingZoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteShippingZoneLogic) DeleteShippingZone(req *types.DeleteShippingZoneReq) error {
	// Get tenant ID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)
	if tenantID == 0 {
		return code.ErrUnauthorized
	}

	// Find existing zone
	zone, err := l.svcCtx.ShippingRepo.FindZoneByID(l.ctx, l.svcCtx.DB, req.ID)
	if err != nil {
		return err
	}

	// Verify template belongs to tenant
	template, err := l.svcCtx.ShippingRepo.FindByID(l.ctx, l.svcCtx.DB, tenantID, zone.TemplateID)
	if err != nil {
		return code.ErrShippingTemplateNotFound
	}
	if template == nil {
		return code.ErrShippingTemplateNotFound
	}

	// Delete zone
	return l.svcCtx.ShippingRepo.DeleteZone(l.ctx, l.svcCtx.DB, req.ID)
}