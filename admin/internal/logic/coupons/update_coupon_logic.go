package coupons

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateCouponLogic {
	return UpdateCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCouponLogic) UpdateCoupon(req *types.UpdateCouponReq) (resp *types.CouponDetailResp, err error) {
	// todo: add your logic here and delete this line

	return
}
