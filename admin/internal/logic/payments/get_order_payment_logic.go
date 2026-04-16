package payments

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderPaymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderPaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetOrderPaymentLogic {
	return GetOrderPaymentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderPaymentLogic) GetOrderPayment(req *types.GetOrderPaymentReq) (resp *types.OrderPaymentResp, err error) {
	// Get tenant ID from context
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok && !contextx.IsPlatformAdmin(l.ctx) {
		return nil, code.ErrUnauthorized
	}
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Get order payment from service
	payment, err := l.svcCtx.PaymentService.GetOrderPayment(l.ctx, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}

	// Convert refunds
	refunds := make([]*types.PaymentRefundResp, len(payment.Refunds))
	for i, refund := range payment.Refunds {
		refunds[i] = &types.PaymentRefundResp{
			ID:              refund.ID,
			RefundNo:        refund.RefundNo,
			ChannelRefundID: refund.ChannelRefundID,
			Amount:          refund.Amount,
			Currency:        refund.Currency,
			Status:          refund.Status,
			StatusText:      refund.StatusText,
			ReasonType:      refund.ReasonType,
			Reason:          refund.Reason,
			RefundedAt:      refund.RefundedAt,
			CreatedAt:       refund.CreatedAt,
		}
	}

	return &types.OrderPaymentResp{
		PaymentID:         payment.PaymentID,
		PaymentNo:         payment.PaymentNo,
		PaymentMethod:     payment.PaymentMethod,
		PaymentMethodText: payment.PaymentMethodText,
		ChannelIntentID:   payment.ChannelIntentID,
		ChannelPaymentID:  payment.ChannelPaymentID,
		Amount:            payment.Amount,
		Currency:          payment.Currency,
		TransactionFee:    payment.TransactionFee,
		FeeCurrency:       payment.FeeCurrency,
		Status:            payment.Status,
		StatusText:        payment.StatusText,
		PaidAt:            payment.PaidAt,
		RefundedAmount:    payment.RefundedAmount,
		Refunds:           refunds,
	}, nil
}
