package refunds

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type RejectRefundLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRejectRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) RejectRefundLogic {
	return RejectRefundLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RejectRefundLogic) RejectRefund(req *types.RejectRefundReq) (resp *types.RefundDetailResp, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	// Get current user ID for updated_by
	userID, _ := contextx.GetUserID(l.ctx)

	refundResp, err := l.svcCtx.RefundApp.RejectRefund(l.ctx, shared.TenantID(tenantID), req.ID, req.RejectReason, userID)
	if err != nil {
		return nil, err
	}

	return convertRefundToDetailResp(refundResp), nil
}