package coupons

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ActivateCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewActivateCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActivateCouponLogic {
	return &ActivateCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ActivateCoupon routes through the unified PromotionApp. Activate
// flips Status → StatusActive (the app layer rejects expired
// end-at values, returning code.ErrPromotionExpired).
func (l *ActivateCouponLogic) ActivateCoupon(req *types.ActivateCouponReq) (resp *types.CouponDetailResp, err error) {
	p, err := l.svcCtx.PromotionApp.Activate(l.ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return convertPromotionToCouponResp(p), nil
}
