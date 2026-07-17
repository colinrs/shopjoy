// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package coupons

import (
	"context"

	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ActivateCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 激活优惠券
func NewActivateCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActivateCouponLogic {
	return &ActivateCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ActivateCouponLogic) ActivateCoupon(req *types.ActivateCouponReq) (resp *types.CouponDetailResp, err error) {
	if err := l.svcCtx.CouponApp.ActivateCoupon(l.ctx, req.ID); err != nil {
		return nil, err
	}

	couponResp, err := l.svcCtx.CouponApp.GetCoupon(l.ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return convertCouponToDetailResp(couponResp), nil
}

var _ = apppromotion.CouponResponse{}
