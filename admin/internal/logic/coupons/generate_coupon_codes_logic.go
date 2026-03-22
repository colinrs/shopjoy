package coupons

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateCouponCodesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateCouponCodesLogic(ctx context.Context, svcCtx *svc.ServiceContext) GenerateCouponCodesLogic {
	return GenerateCouponCodesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateCouponCodesLogic) GenerateCouponCodes(req *types.GenerateCouponCodesReq) (resp *types.GenerateCouponCodesResp, err error) {
	// todo: add your logic here and delete this line

	return
}
