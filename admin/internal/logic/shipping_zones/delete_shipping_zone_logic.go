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
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return err
	}

	// Find existing zone
	zone, err := l.svcCtx.ShippingRepo.FindZoneByID(l.ctx, l.svcCtx.DB, req.ID)
	if err != nil {
		return err
	}

	// Verify zone belongs to tenant
	if zone.TenantID != tenantID {
		return code.ErrShippingZoneNotFound
	}

	// Delete zone (this will also delete zone_regions)
	return l.svcCtx.ShippingRepo.DeleteZone(l.ctx, l.svcCtx.DB, req.ID)
}