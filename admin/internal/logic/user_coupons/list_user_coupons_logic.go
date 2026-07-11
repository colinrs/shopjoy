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

	// Map status string to domain type
	status := mapUserCouponStatus(req.Status)

	listResp, err := l.svcCtx.CouponApp.ListUserCoupons(
		l.ctx,
		req.UserID,
		status,
		req.Page,
		req.PageSize,
	)
	if err != nil {
		return nil, err
	}

	list := make([]*types.UserCouponDetailResp, len(listResp.List))
	for i, uc := range listResp.List {
		list[i] = convertUserCouponToDetailResp(uc)
	}

	return &types.ListUserCouponsResp{
		List:     list,
		Total:    listResp.Total,
		Page:     listResp.Page,
		PageSize: listResp.PageSize,
	}, nil
}
