package coupons

import (
	"context"

	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
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

	queryReq := apppromotion.QueryCouponRequest{
		Name:     req.Name,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	// Only set the typed filters when the caller actually supplied a value.
	// mapCouponType("") and mapCouponStatusToInt("") both return zero-valued
	// enums (CouponTypeFixedAmount / CouponStatusInactive) which are valid
	// filter values, so we cannot rely on zero-detection downstream.
	if req.Type != "" {
		t := mapCouponType(req.Type)
		queryReq.Type = &t
	}
	if req.Status != "" {
		s := mapCouponStatusToInt(req.Status)
		queryReq.Status = &s
	}

	listResp, err := l.svcCtx.CouponApp.ListCoupons(l.ctx, queryReq)
	if err != nil {
		return nil, err
	}

	list := make([]*types.CouponDetailResp, len(listResp.List))
	for i, c := range listResp.List {
		list[i] = convertCouponToDetailResp(c)
	}

	return &types.ListCouponsResp{
		List:     list,
		Total:    listResp.Total,
		Page:     listResp.Page,
		PageSize: listResp.PageSize,
	}, nil
}
