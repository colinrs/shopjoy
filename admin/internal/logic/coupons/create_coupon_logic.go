package coupons

import (
	"context"
	"time"

	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

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
	// Get tenantID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

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
		Value:        parseMoneyToInt64(req.DiscountValue),
		MinAmount:    parseMoneyToInt64(req.MinOrderAmount),
		MaxDiscount:  parseMoneyToInt64(req.MaxDiscount),
		TotalCount:   req.UsageLimit,
		PerUserLimit: req.PerUserLimit,
		StartAt:      startAt,
		EndAt:        endAt,
	}

	couponResp, err := l.svcCtx.CouponApp.CreateCoupon(l.ctx, shared.TenantID(tenantID), createReq)
	if err != nil {
		return nil, err
	}

	return &types.CreateCouponResp{
		ID: couponResp.ID,
	}, nil
}