package coupons

import (
	"context"

	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

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
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	queryReq := apppromotion.QueryCouponRequest{
		Name:     req.Name,
		Type:     mapCouponType(req.Type),
		Status:   mapCouponStatusToInt(req.Status),
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	listResp, err := l.svcCtx.CouponApp.ListCoupons(l.ctx, shared.TenantID(tenantID), queryReq)
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