package shipments

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetShipmentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetShipmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetShipmentLogic {
	return GetShipmentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetShipmentLogic) GetShipment(req *types.GetShipmentReq) (resp *types.ShipmentDetailResp, err error) {
	// Get tenantID from context

	// Get shipment
	shipmentResp, err := l.svcCtx.ShipmentApp.GetShipment(l.ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return toShipmentDetailResp(shipmentResp), nil
}
