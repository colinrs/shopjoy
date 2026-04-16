package shipments

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderShipmentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderShipmentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetOrderShipmentsLogic {
	return GetOrderShipmentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderShipmentsLogic) GetOrderShipments(req *types.GetOrderShipmentsReq) (resp *types.ListOrderShipmentsResp, err error) {
	// Get tenantID from context
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok && !contextx.IsPlatformAdmin(l.ctx) {
		return nil, code.ErrUnauthorized
	}
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Get order shipments
	shipments, err := l.svcCtx.ShipmentApp.GetOrderShipments(l.ctx, shared.TenantID(tenantID), req.OrderID)
	if err != nil {
		return nil, err
	}

	// Build response
	resp = &types.ListOrderShipmentsResp{
		List: make([]*types.ShipmentDetailResp, len(shipments)),
	}

	for i, s := range shipments {
		resp.List[i] = toShipmentDetailResp(s)
	}

	return resp, nil
}
