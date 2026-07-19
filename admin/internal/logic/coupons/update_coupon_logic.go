package coupons

import (
	"context"
	"time"

	apppromotion "github.com/colinrs/shopjoy/admin/internal/application/promotion"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	pkgpromotion "github.com/colinrs/shopjoy/pkg/domain/promotion"

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
	rules := []pkgpromotion.PromotionRule{rule}

	updateReq := &apppromotion.UpdatePromotionRequest{
		ID:           req.ID,
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
		Rules:        &rules,
		ActorID:      actorID(l.ctx),
	}

	p, err := l.svcCtx.PromotionApp.Update(l.ctx, updateReq)
	if err != nil {
		return nil, err
	}

	return convertPromotionToCouponResp(p), nil
}
