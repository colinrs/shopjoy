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

// GetCoupon fetches the unified PromotionResponse and projects it
// back to the legacy wire CouponDetailResp. The repo enforces
// Kind=COUPON (the wire endpoint is coupon-scoped) so an attempt to
// GET a non-COUPON id here returns the promotion row as a coupon
// payload — acceptable until Task 8 wires Kind into the request.
func (l *GetCouponLogic) GetCoupon(req *types.GetCouponReq) (resp *types.CouponDetailResp, err error) {
	p, err := l.svcCtx.PromotionApp.Get(l.ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return convertPromotionToCouponResp(p), nil
}