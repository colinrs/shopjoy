package coupons

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateCouponLogic {
	return CreateCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCouponLogic) CreateCoupon(req *types.CreateCouponReq) (resp *types.CreateCouponResp, err error) {
	// todo: add your logic here and delete this line

	return
}
