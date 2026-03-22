package coupons

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCouponUsageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCouponUsageLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetCouponUsageLogic {
	return GetCouponUsageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCouponUsageLogic) GetCouponUsage(req *types.GetCouponUsageReq) (resp *types.ListCouponUsageResp, err error) {
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Note: This would require a dedicated method in the app service
	// For now, return empty list
	_ = shared.TenantID(tenantID)

	return &types.ListCouponUsageResp{
		List:     []*types.CouponUsageResp{},
		Total:    0,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}