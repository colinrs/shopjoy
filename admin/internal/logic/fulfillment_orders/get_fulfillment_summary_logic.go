package fulfillment_orders

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFulfillmentSummaryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFulfillmentSummaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetFulfillmentSummaryLogic {
	return GetFulfillmentSummaryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFulfillmentSummaryLogic) GetFulfillmentSummary(req *types.GetFulfillmentSummaryReq) (resp *types.FulfillmentSummaryResp, err error) {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Get fulfillment summary
	summary, err := l.svcCtx.OrderFulfillmentApp.GetFulfillmentSummary(l.ctx, shared.TenantID(tenantID))
	if err != nil {
		return nil, err
	}

	return &types.FulfillmentSummaryResp{
		PendingShipment: summary.PendingShipment,
		PartialShipped:  summary.PartialShipped,
		Shipped:         summary.Shipped,
		Delivered:       summary.Delivered,
		PendingRefund:   summary.PendingRefund,
		Refunding:       summary.Refunding,
		TotalOrders:     summary.TotalOrders,
	}, nil
}