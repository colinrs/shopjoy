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
	// todo: add your logic here and delete this line

	return
}
