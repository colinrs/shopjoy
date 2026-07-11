package coupons

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetCouponLogic {
	return GetCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCouponLogic) GetCoupon(req *types.GetCouponReq) (resp *types.CouponDetailResp, err error) {

	couponResp, err := l.svcCtx.CouponApp.GetCoupon(l.ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return convertCouponToDetailResp(couponResp), nil
}
