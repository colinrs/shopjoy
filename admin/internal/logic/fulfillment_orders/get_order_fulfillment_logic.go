package fulfillment_orders

import (
	"context"
	"fmt"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderFulfillmentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderFulfillmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetOrderFulfillmentLogic {
	return GetOrderFulfillmentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderFulfillmentLogic) GetOrderFulfillment(req *types.GetOrderFulfillmentReq) (resp *types.OrderFulfillmentDetailResp, err error) {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// For now, we use the ID as the order ID string
	// In a real implementation, you would look up the order by ID
	orderID := fmt.Sprintf("%d", req.ID)

	// Get order fulfillment detail
	detail, err := l.svcCtx.OrderFulfillmentApp.GetOrderFulfillment(l.ctx, shared.TenantID(tenantID), orderID)
	if err != nil {
		return nil, err
	}

	return toOrderFulfillmentDetailResp(detail), nil
}