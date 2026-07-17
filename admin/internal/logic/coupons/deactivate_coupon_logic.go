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

type DeactivateCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 停用优惠券
func NewDeactivateCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeactivateCouponLogic {
	return &DeactivateCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeactivateCouponLogic) DeactivateCoupon(req *types.DeactivateCouponReq) (resp *types.CouponDetailResp, err error) {
	if err := l.svcCtx.CouponApp.DeactivateCoupon(l.ctx, req.ID); err != nil {
		return nil, err
	}

	couponResp, err := l.svcCtx.CouponApp.GetCoupon(l.ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return convertCouponToDetailResp(couponResp), nil
}

var _ = apppromotion.CouponResponse{}
