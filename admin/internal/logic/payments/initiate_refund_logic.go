package payments

import (
	"context"

	appPayment "github.com/colinrs/shopjoy/admin/internal/application/payment"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type InitiateRefundLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInitiateRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) InitiateRefundLogic {
	return InitiateRefundLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InitiateRefundLogic) InitiateRefund(req *types.InitiateRefundReq) (resp *types.InitiateRefundResp, err error) {
	// Get tenant ID and admin ID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)
	if tenantID == 0 {
		tenantID = 0
	}
	adminID := contextx.GetCurrentUserID(l.ctx)

	// Build request
	appReq := appPayment.InitiateRefundRequest{
		OrderID:        req.ID,
		IdempotencyKey: req.IdempotencyKey,
		Amount:         req.Amount,
		ReasonType:     req.ReasonType,
		Reason:         req.Reason,
	}

	// Initiate refund through service
	result, err := l.svcCtx.PaymentService.InitiateRefund(l.ctx, shared.TenantID(tenantID), adminID, appReq)
	if err != nil {
		return nil, err
	}

	return &types.InitiateRefundResp{
		RefundID:        result.RefundID,
		RefundNo:        result.RefundNo,
		Amount:          result.Amount,
		Currency:        result.Currency,
		Status:          result.Status,
		StatusText:      result.StatusText,
		ChannelRefundID: result.ChannelRefundID,
	}, nil
}
