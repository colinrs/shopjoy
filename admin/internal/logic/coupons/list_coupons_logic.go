package coupons

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCouponsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListCouponsLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListCouponsLogic {
	return ListCouponsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListCouponsLogic) ListCoupons(req *types.ListCouponsReq) (resp *types.ListCouponsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
