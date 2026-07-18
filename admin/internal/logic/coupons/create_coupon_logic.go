package coupons

import (
	"context"
	"time"

	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateCouponLogic {
	return CreateCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCouponLogic) CreateCoupon(req *types.CreateCouponReq) (resp *types.CreateCouponResp, err error) {

	// Parse time
	startAt, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, err
	}
	endAt, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, err
	}

	createReq := apppromotion.CreateCouponRequest{
		Name:         req.Name,
		Code:         req.Code,
		Description:  req.Description,
		Type:         mapCouponType(req.Type),
		Value:        parseMoneyToDecimal(req.DiscountValue),
		MinAmount:    parseMoneyToDecimal(req.MinOrderAmount),
		MaxDiscount:  parseMoneyToDecimal(req.MaxDiscount),
		TotalCount:   req.UsageLimit,
		PerUserLimit: req.PerUserLimit,
		Scope:        buildCouponScope(req.ScopeType),
		StartAt:      startAt,
		EndAt:        endAt,
	}

	couponResp, err := l.svcCtx.CouponApp.CreateCoupon(l.ctx, createReq)
	if err != nil {
		return nil, err
	}

	return &types.CreateCouponResp{
		ID: couponResp.ID,
	}, nil
}
