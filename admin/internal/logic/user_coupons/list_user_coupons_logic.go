package user_coupons

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserCouponsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListUserCouponsLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListUserCouponsLogic {
	return ListUserCouponsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUserCouponsLogic) ListUserCoupons(req *types.ListUserCouponsReq) (resp *types.ListUserCouponsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
