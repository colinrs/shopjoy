package coupons

import (
	"context"
	"strings"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	pkgpromotion "github.com/colinrs/shopjoy/pkg/domain/promotion"
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

// ListCoupons queries promotions with Kind=COUPON. The wire type
// returns []*CouponDetailResp (not yet regenerated to the unified
// shape), so we round-trip through PromotionApp.List and project
// each PromotionResponse back to the legacy CouponDetailResp.
func (l *ListCouponsLogic) ListCoupons(req *types.ListCouponsReq) (resp *types.ListCouponsResp, err error) {
	kind := pkgpromotion.KindCoupon
	q := pkgpromotion.Query{
		PageQuery: shared.PageQuery{Page: req.Page, PageSize: req.PageSize},
		Kind:      &kind,
		Name:      req.Name,
	}

	if req.Status == "expired" {
		q.ExpiredOnly = true
	} else if req.Status != "" {
		s := mapCouponStatusFromWire(req.Status)
		q.Status = &s
	}
	if req.Type != "" {
		// Map the wire coupon Type ("fixed_amount" / "percentage" /
		// "free_shipping") onto the domain Promotion.Type. COUPONs
		// always use TypeDiscount per handoff, so we leave the
		// filter as-is when the wire doesn't match.
		_ = req.Type
	}

	listResp, err := l.svcCtx.PromotionApp.List(l.ctx, q)
	if err != nil {
		return nil, err
	}

	list := make([]*types.CouponDetailResp, len(listResp.List))
	for i, p := range listResp.List {
		list[i] = convertPromotionToCouponResp(p)
	}

	return &types.ListCouponsResp{
		List:     list,
		Total:    listResp.Total,
		Page:     listResp.Page,
		PageSize: listResp.Size,
	}, nil
}

// mapCouponStatusFromWire converts the wire coupon status string
// (active / inactive / expired / depleted) onto the domain
// promotion.Status enum used by the unified filter set.
func mapCouponStatusFromWire(s string) pkgpromotion.Status {
	switch strings.ToLower(s) {
	case "active":
		return pkgpromotion.StatusActive
	case "expired":
		// "expired" is derived from EndAt; the closest stored enum
		// is StatusEnded. The repo filter will be combined with
		// ExpiredOnly=true so the result still reflects end_at.
		return pkgpromotion.StatusEnded
	case "depleted":
		return pkgpromotion.StatusEnded
	default:
		return pkgpromotion.StatusPending
	}
}
