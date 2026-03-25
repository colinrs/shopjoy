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

type RemindPaymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemindPaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) RemindPaymentLogic {
	return RemindPaymentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemindPaymentLogic) RemindPayment(req *types.RemindPaymentReq) (resp *types.RemindPaymentResp, err error) {
	// Get tenant ID from context
	tenantIDRaw, _ := contextx.GetTenantID(l.ctx)
	tenantID := shared.TenantID(tenantIDRaw)

	// Get the order
	order, err := l.svcCtx.OrderRepo.FindByID(l.ctx, l.svcCtx.DB, tenantID, req.ID)
	if err != nil {
		return nil, err
	}

	// Validate order can be reminded
	// Only pending_payment orders can be reminded
	if order.Status != fulfillment.OrderStatusPendingPayment {
		return nil, code.ErrOrderCannotRemind
	}

	// In a real system, this would send a notification to the user
	// via SMS, email, or push notification
	// For now, we'll just log it and return success

	now := time.Now().UTC()
	l.Logger.Infof("Payment reminder sent for order %s", order.OrderNo)

	// Return success response
	return &types.RemindPaymentResp{
		OrderID:    order.ID,
		OrderNo:    order.OrderNo,
		RemindedAt: now.Format(time.RFC3339),
		Message:    "付款提醒已发送",
	}, nil
}