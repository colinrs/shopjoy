package shipments

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateShipmentStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateShipmentStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateShipmentStatusLogic {
	return UpdateShipmentStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateShipmentStatusLogic) UpdateShipmentStatus(req *types.UpdateShipmentStatusReq) (resp *types.ShipmentDetailResp, err error) {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Get user ID from context
	userID := contextx.GetCurrentUserID(l.ctx)

	// Update shipment status
	shipmentResp, err := l.svcCtx.ShipmentApp.UpdateShipmentStatus(l.ctx, shared.TenantID(tenantID), userID, req.ID, fulfillment.ShipmentStatus(req.Status))
	if err != nil {
		return nil, err
	}

	return toShipmentDetailResp(shipmentResp), nil
}