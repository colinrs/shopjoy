package fulfillment_orders

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

type CancelOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) CancelOrderLogic {
	return CancelOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelOrderLogic) CancelOrder(req *types.CancelOrderReq) (resp *types.CancelOrderResp, err error) {
	// Get tenant ID from context
	tenantIDRaw, _ := contextx.GetTenantID(l.ctx)
	tenantID := shared.TenantID(tenantIDRaw)

	// Get the order
	order, err := l.svcCtx.OrderRepo.FindByID(l.ctx, l.svcCtx.DB, tenantID, req.ID)
	if err != nil {
		return nil, err
	}

	// Validate order can be cancelled
	// Only pending_payment orders can be cancelled by admin
	if order.Status != fulfillment.OrderStatusPendingPayment {
		return nil, code.ErrOrderCannotCancel
	}

	// Update order status
	now := time.Now().UTC()
	order.Status = fulfillment.OrderStatusCancelled
	order.CancelledAt = &now
	order.Remark = req.Reason

	// Save the order
	err = l.svcCtx.OrderRepo.UpdateWithVersion(l.ctx, l.svcCtx.DB, order)
	if err != nil {
		return nil, err
	}

	// Restore inventory (if needed)
	// This would typically be done through an event or a separate service
	// For now, we'll just log it
	l.Logger.Infof("Order %s cancelled, inventory should be restored", order.OrderNo)

	return &types.CancelOrderResp{
		OrderID:     order.ID,
		OrderNo:     order.OrderNo,
		Status:      string(order.Status),
		StatusText:  order.Status.Text(),
		CancelledAt: now.Format(time.RFC3339),
		Reason:      req.Reason,
	}, nil
}