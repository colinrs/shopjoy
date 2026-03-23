package fulfillment_orders

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrderRemarkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateOrderRemarkLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateOrderRemarkLogic {
	return UpdateOrderRemarkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateOrderRemarkLogic) UpdateOrderRemark(req *types.UpdateOrderRemarkReq) (resp *types.UpdateOrderRemarkResp, err error) {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Truncate remark to 500 characters (using rune count for UTF-8 safety)
	remark := truncateString(req.Remark, 500)

	// Update order remark
	err = l.svcCtx.OrderFulfillmentApp.UpdateOrderRemark(l.ctx, shared.TenantID(tenantID), req.ID, remark)
	if err != nil {
		return nil, err
	}

	return &types.UpdateOrderRemarkResp{
		OrderID: req.ID,
	}, nil
}