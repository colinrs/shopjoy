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

type AdjustOrderPriceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdjustOrderPriceLogic(ctx context.Context, svcCtx *svc.ServiceContext) AdjustOrderPriceLogic {
	return AdjustOrderPriceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdjustOrderPriceLogic) AdjustOrderPrice(req *types.AdjustOrderPriceReq) (resp *types.AdjustOrderPriceResp, err error) {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Get user ID from context
	userID := contextx.GetCurrentUserID(l.ctx)

	// Validate adjust_amount is provided
	if req.AdjustAmount == "" {
		return nil, code.ErrParam
	}

	// Parse adjust amount from string
	adjustAmount, err := parseMoneyToDecimal(req.AdjustAmount)
	if err != nil {
		return nil, err
	}

	// Validate reason is provided
	if req.Reason == "" {
		return nil, code.ErrOrderAdjustReasonRequired
	}

	// Truncate reason to 200 characters (using rune count for UTF-8 safety)
	reason := truncateString(req.Reason, 200)

	// Adjust order price
	result, err := l.svcCtx.OrderFulfillmentApp.AdjustOrderPrice(l.ctx, shared.TenantID(tenantID), userID, req.ID, adjustAmount, reason)
	if err != nil {
		return nil, err
	}

	return &types.AdjustOrderPriceResp{
		OrderID:        result.OrderID,
		OriginalAmount: formatDecimal(result.OriginalAmount),
		AdjustAmount:   formatDecimal(result.AdjustAmount),
		NewPayAmount:   formatDecimal(result.NewPayAmount),
		AdjustReason:   result.AdjustReason,
		AdjustedAt:     result.AdjustedAt,
	}, nil
}