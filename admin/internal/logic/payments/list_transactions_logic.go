package payments

import (
	"context"
	"time"

	appPayment "github.com/colinrs/shopjoy/admin/internal/application/payment"
	"github.com/colinrs/shopjoy/admin/internal/domain/payment"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTransactionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListTransactionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListTransactionsLogic {
	return ListTransactionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTransactionsLogic) ListTransactions(req *types.ListTransactionsReq) (resp *types.ListTransactionsResp, err error) {
	// Get tenant ID from context
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok && !contextx.IsPlatformAdmin(l.ctx) {
		return nil, code.ErrUnauthorized
	}
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Build request
	appReq := appPayment.ListTransactionsRequest{
		Page:          req.Page,
		PageSize:      req.PageSize,
		TransactionID: req.TransactionID,
		PaymentMethod: payment.PaymentMethod(req.PaymentMethod),
		Status:        payment.TransactionStatus(req.Status),
	}

	// Parse time filters
	if req.StartTime != "" {
		startTime, err := time.Parse(time.RFC3339, req.StartTime)
		if err == nil {
			appReq.StartTime = startTime
		}
	}
	if req.EndTime != "" {
		endTime, err := time.Parse(time.RFC3339, req.EndTime)
		if err == nil {
			appReq.EndTime = endTime
		}
	}

	// Get transactions from service
	result, err := l.svcCtx.PaymentService.ListTransactions(l.ctx, shared.TenantID(tenantID), appReq)
	if err != nil {
		return nil, err
	}

	// Convert response
	list := make([]*types.TransactionResp, len(result.List))
	for i, txn := range result.List {
		list[i] = &types.TransactionResp{
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
		}
	}

	return &types.ListTransactionsResp{
		List:     list,
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.PageSize,
		Stats: types.TransactionStats{
			Success: result.Stats.Success,
			Pending: result.Stats.Pending,
			Failed:  result.Stats.Failed,
		},
	}, nil
}
