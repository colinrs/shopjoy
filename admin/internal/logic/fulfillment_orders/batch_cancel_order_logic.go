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

type BatchCancelOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) BatchCancelOrderLogic {
	return BatchCancelOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchCancelOrderLogic) BatchCancelOrder(req *types.BatchCancelOrderReq) (resp *types.BatchCancelOrderResp, err error) {
	// Get tenant ID from context
	tenantIDRaw, _ := contextx.GetTenantID(l.ctx)
	tenantID := shared.TenantID(tenantIDRaw)

	var success []int64
	var failed []types.BatchCancelFail

	for _, orderID := range req.OrderIDs {
		failEntry := types.BatchCancelFail{
			OrderID: orderID,
		}

		// Get the order
		order, err := l.svcCtx.OrderRepo.FindByID(l.ctx, l.svcCtx.DB, tenantID, orderID)
		if err != nil {
			failEntry.Code = code.ErrOrderNotFound.Code
			failEntry.Message = code.ErrOrderNotFound.Msg
			failed = append(failed, failEntry)
			continue
		}

		// Validate order can be cancelled
		// Only pending_payment orders can be cancelled by admin
		if order.Status != fulfillment.OrderStatusPendingPayment {
			failEntry.Code = code.ErrOrderCannotCancel.Code
			failEntry.Message = code.ErrOrderCannotCancel.Msg
			failed = append(failed, failEntry)
			continue
		}

		// Update order status
		now := time.Now().UTC()
		order.Status = fulfillment.OrderStatusCancelled
		order.CancelledAt = &now
		order.Remark = req.Reason

		err = l.svcCtx.OrderRepo.UpdateWithVersion(l.ctx, l.svcCtx.DB, order)
		if err != nil {
			failEntry.Code = code.ErrInternalServer.Code
			failEntry.Message = err.Error()
			failed = append(failed, failEntry)
			continue
		}

		// Log inventory restoration
		l.Logger.Infof("Order %s cancelled, inventory should be restored", order.OrderNo)

		success = append(success, orderID)
	}

	resp = &types.BatchCancelOrderResp{
		Success: success,
		Failed:  failed,
	}

	return resp, nil
}