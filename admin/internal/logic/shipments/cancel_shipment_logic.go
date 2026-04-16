package shipments

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/fulfillment"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelShipmentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelShipmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) CancelShipmentLogic {
	return CancelShipmentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelShipmentLogic) CancelShipment(req *types.CancelShipmentReq) (resp *types.CancelShipmentResp, err error) {
	// Get tenantID from context
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok && !contextx.IsPlatformAdmin(l.ctx) {
		return nil, code.ErrUnauthorized
	}
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Get user ID from context for audit
	userID, _ := contextx.GetUserID(l.ctx)

	// First get the shipment to validate and capture info
	shipment, err := l.svcCtx.ShipmentApp.GetShipment(l.ctx, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}

	// Capture shipment info for response before transaction
	shipmentNo := shipment.ShipmentNo

	// Cancel the shipment using the application layer
	_, err = l.svcCtx.ShipmentApp.CancelShipment(l.ctx, shared.TenantID(tenantID), userID, req.ID, req.Reason)
	if err != nil {
		return nil, err
	}

	return &types.CancelShipmentResp{
		ID:          req.ID,
		ShipmentNo:  shipmentNo,
		Status:      fulfillment.ShipmentStatusCancelled.String(),
		StatusText:  fulfillment.ShipmentStatusCancelled.String(),
		CancelledAt: time.Now().UTC().Format(time.RFC3339),
		Reason:      req.Reason,
	}, nil
}
