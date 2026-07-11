package payments

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTransactionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetTransactionLogic {
	return GetTransactionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTransactionLogic) GetTransaction(req *types.GetTransactionReq) (resp *types.TransactionResp, err error) {
	// Get tenant ID from context

	// Get transaction from service
	txn, err := l.svcCtx.PaymentService.GetTransaction(l.ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &types.TransactionResp{
		ID:                   txn.ID,
		TransactionID:        txn.TransactionID,
		OrderID:              txn.OrderID,
		OrderNo:              txn.OrderNo,
		PaymentMethod:        txn.PaymentMethod,
		PaymentMethodText:    txn.PaymentMethodText,
		ChannelTransactionID: txn.ChannelTransactionID,
		Amount:               txn.Amount,
		Currency:             txn.Currency,
		TransactionFee:       txn.TransactionFee,
		Status:               txn.Status,
		StatusText:           txn.StatusText,
		CreatedAt:            txn.CreatedAt,
		PaidAt:               txn.PaidAt,
	}, nil
}
