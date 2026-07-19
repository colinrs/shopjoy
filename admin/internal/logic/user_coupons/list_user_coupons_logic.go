package user_coupons

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	pkgpromotion "github.com/colinrs/shopjoy/pkg/domain/promotion"

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
	status := mapUserCouponStatus(req.Status)

	var statusPtr *pkgpromotion.UserCouponStatus
	if req.Status != "" {
		s := status
		statusPtr = &s
	}

	var userIDPtr *int64
	if req.UserID != 0 {
		uid := req.UserID
		userIDPtr = &uid
	}

	q := pkgpromotion.UserCouponQuery{
		UserID: userIDPtr,
		Status: statusPtr,
	}

	listResp, err := l.svcCtx.PromotionApp.ListUserCoupons(l.ctx, q)
	if err != nil {
		return nil, err
	}

	list := make([]*types.UserCouponDetailResp, len(listResp.List))
	for i, uc := range listResp.List {
		list[i] = convertUserCouponToDetailResp(uc)
	}

	return &types.ListUserCouponsResp{
		List:  list,
		Total: listResp.Total,
		Page:  req.Page,
		PageSize: req.PageSize,
	}, nil
}