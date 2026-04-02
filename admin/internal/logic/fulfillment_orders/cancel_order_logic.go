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
	"gorm.io/gorm"
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

	// Get user ID from context for audit
	userID, _ := contextx.GetUserID(l.ctx)

	// First get the order to validate and capture order info
	order, err := l.svcCtx.OrderRepo.FindByID(l.ctx, l.svcCtx.DB, tenantID, req.ID)
	if err != nil {
		return nil, err
	}

	// Validate order can be cancelled
	// Only pending_payment orders can be cancelled by admin
	if order.Status != fulfillment.OrderStatusPendingPayment {
		return nil, code.ErrOrderCannotCancel
	}

	// Capture order info for response before transaction
	orderNo := order.OrderNo

	// Use transaction for cancel operation
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// Update order status
		now := time.Now().UTC()
		order.Status = fulfillment.OrderStatusCancelled
		order.CancelledAt = &now
		order.CancelledBy = userID
		order.Remark = req.Reason

		// Save the order
		return l.svcCtx.OrderRepo.UpdateWithVersion(l.ctx, tx, order)
	})
	if err != nil {
		return nil, err
	}

	// Restore inventory (if needed)
	// This would typically be done through an event or a separate service
	// For now, we'll just log it
	l.Logger.Infof("Order %s cancelled, inventory should be restored", orderNo)

	return &types.CancelOrderResp{
		OrderID:     order.ID,
		OrderNo:     orderNo,
		Status:      string(fulfillment.OrderStatusCancelled),
		StatusText:  fulfillment.OrderStatusCancelled.Text(),
		CancelledAt: time.Now().UTC().Format(time.RFC3339),
		Reason:      req.Reason,
	}, nil
}
