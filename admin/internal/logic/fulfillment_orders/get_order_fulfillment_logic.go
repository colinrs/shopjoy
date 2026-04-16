package fulfillment_orders

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
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
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok && !contextx.IsPlatformAdmin(l.ctx) {
		return nil, code.ErrUnauthorized
	}
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Get order fulfillment detail
	detail, err := l.svcCtx.OrderFulfillmentApp.GetOrderFulfillment(l.ctx, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}

	return toOrderFulfillmentDetailResp(detail), nil
}
