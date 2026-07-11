package coupons

import (
	"context"
	"time"

	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateCouponLogic {
	return UpdateCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCouponLogic) UpdateCoupon(req *types.UpdateCouponReq) (resp *types.CouponDetailResp, err error) {

	// Parse time
	startAt, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, err
	}
	endAt, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, err
	}

	updateReq := apppromotion.UpdateCouponRequest{
		ID:           req.ID,
		Name:         req.Name,
		Description:  req.Description,
		MinAmount:    parseMoneyToDecimal(req.MinOrderAmount),
		MaxDiscount:  parseMoneyToDecimal(req.MaxDiscount),
		TotalCount:   req.UsageLimit,
		PerUserLimit: req.PerUserLimit,
		StartAt:      startAt,
		EndAt:        endAt,
	}

	couponResp, err := l.svcCtx.CouponApp.UpdateCoupon(l.ctx, updateReq)
	if err != nil {
		return nil, err
	}

	return convertCouponToDetailResp(couponResp), nil
}
