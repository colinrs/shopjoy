package coupons

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeactivateCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeactivateCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeactivateCouponLogic {
	return &DeactivateCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DeactivateCoupon flips Status → StatusPaused via the unified app.
func (l *DeactivateCouponLogic) DeactivateCoupon(req *types.DeactivateCouponReq) (resp *types.CouponDetailResp, err error) {
	p, err := l.svcCtx.PromotionApp.Deactivate(l.ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return convertPromotionToCouponResp(p), nil
}
