package coupons

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteCouponLogic {
	return DeleteCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DeleteCoupon removes the unified Promotion row (kind=COUPON). The
// app layer also deletes the associated rules in the same call.
func (l *DeleteCouponLogic) DeleteCoupon(req *types.DeleteCouponReq) (resp *types.CreateCouponResp, err error) {
	if err := l.svcCtx.PromotionApp.Delete(l.ctx, req.ID); err != nil {
		return nil, err
	}
	return &types.CreateCouponResp{ID: req.ID}, nil
}
