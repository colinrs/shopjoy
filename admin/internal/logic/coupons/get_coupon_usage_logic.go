package coupons

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

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

	// Note: This would require a dedicated method in the app service
	// For now, return empty list

	return &types.ListCouponUsageResp{
		List:     []*types.CouponUsageResp{},
		Total:    0,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
