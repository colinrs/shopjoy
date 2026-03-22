package refunds

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRefundLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetRefundLogic {
	return GetRefundLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRefundLogic) GetRefund(req *types.GetRefundReq) (resp *types.RefundDetailResp, err error) {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	refundResp, err := l.svcCtx.RefundApp.GetRefund(l.ctx, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}

	return convertRefundToDetailResp(refundResp), nil
}