package coupons

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCouponUsageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCouponUsageLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetCouponUsageLogic {
	return GetCouponUsageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCouponUsageLogic) GetCouponUsage(req *types.GetCouponUsageReq) (resp *types.ListCouponUsageResp, err error) {
	// todo: add your logic here and delete this line

	return
}
