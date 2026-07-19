package coupons

import (
	"context"
	"time"

	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	pkgpromotion "github.com/colinrs/shopjoy/pkg/domain/promotion"
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

// CreateCoupon folds the wire coupon request into a single
// PromotionRule (one rule per coupon for the current form shape) and
// routes through PromotionApp.Create with Kind=KindCoupon.
//
// Per handoff I2, COUPON rows may carry usage_limit=0 ("unlimited")
// after the data migration. We honor that and don't validate the
// limit here — the app layer handles empty TotalCount (=nil) cleanly.
func (l *CreateCouponLogic) CreateCoupon(req *types.CreateCouponReq) (resp *types.CreateCouponResp, err error) {
	startAt, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, err
	}
	endAt, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, err
	}

	code := req.Code
	total := req.UsageLimit

	rule := pkgpromotion.PromotionRule{
		ConditionType: pkgpromotion.ConditionMinAmount,
		ActionType:    mapCouponActionType(req.Type),
	}
	if req.MinOrderAmount != "" {
		rule.ConditionValue = parseMoney(req.MinOrderAmount)
	}
	if req.DiscountValue != "" {
		rule.ActionValue = parseMoney(req.DiscountValue)
	}
	if req.MaxDiscount != "" {
		rule.MaxDiscount = parseMoney(req.MaxDiscount)
	}

	createReq := &apppromotion.CreatePromotionRequest{
		TenantID:     shared.TenantID(tenantID(l.ctx)),
		Kind:         pkgpromotion.KindCoupon,
		Name:         req.Name,
		Description:  req.Description,
		Code:         &code,
		Type:         pkgpromotion.TypeDiscount,
		Currency:     defaultCurrency(""),
		TotalCount:   &total,
		UsageLimit:   req.UsageLimit,
		PerUserLimit: req.PerUserLimit,
		Scope:        buildCouponScope(""),
		StartAt:      startAt,
		EndAt:        endAt,
		Rules:        []pkgpromotion.PromotionRule{rule},
		ActorID:      actorID(l.ctx),
	}

	p, err := l.svcCtx.PromotionApp.Create(l.ctx, createReq)
	if err != nil {
		return nil, err
	}

	return &types.CreateCouponResp{ID: p.ID}, nil
}