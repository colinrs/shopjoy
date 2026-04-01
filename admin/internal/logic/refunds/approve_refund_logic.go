package refunds

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApproveRefundLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApproveRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) ApproveRefundLogic {
	return ApproveRefundLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApproveRefundLogic) ApproveRefund(req *types.ApproveRefundReq) (resp *types.RefundDetailResp, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	// Get current user ID for approved_by
	userID, _ := contextx.GetUserID(l.ctx)

	refundResp, err := l.svcCtx.RefundApp.ApproveRefund(l.ctx, shared.TenantID(tenantID), req.ID, userID)
	if err != nil {
		return nil, err
	}

	return convertRefundToDetailResp(refundResp), nil
}