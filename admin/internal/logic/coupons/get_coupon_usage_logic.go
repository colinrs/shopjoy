package coupons

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	pkgpromotion "github.com/colinrs/shopjoy/pkg/domain/promotion"
	"github.com/shopspring/decimal"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCouponUsageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCouponUsageLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetCouponUsageLogic {
	return GetCouponUsageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetCouponUsage routes through PromotionApp.FindPromotionUsage.
// The wire response shape (CouponUsageResp) keeps the legacy
// pre-merge field names; we project PromotionUsageResponse onto it.
func (l *GetCouponUsageLogic) GetCouponUsage(req *types.GetCouponUsageReq) (resp *types.ListCouponUsageResp, err error) {
	couponID := req.ID
	q := pkgpromotion.UsageQuery{
		Page:     req.Page,
		Size:     req.PageSize,
		CouponID: &couponID,
	}
	usageResp, err := l.svcCtx.PromotionApp.FindPromotionUsage(l.ctx, q)
	if err != nil {
		return nil, err
	}

	list := make([]*types.CouponUsageResp, len(usageResp.List))
	for i, u := range usageResp.List {
		out := &types.CouponUsageResp{
			ID:             u.ID,
			UserID:         u.UserID,
			OrderID:        u.OrderID,
			DiscountAmount: formatMoney(u.DiscountAmount),
		}
		if u.CouponID != nil {
			out.CouponID = *u.CouponID
		} else {
			out.CouponID = req.ID
		}
		// CouponUsageResp carries a "UsedAt" string; fall back to
		// CreatedAt when the usage row didn't record a used_at.
		usedAt := u.CreatedAt
		out.UsedAt = usedAt.Format(time.RFC3339)
		list[i] = out
	}

	// Avoid "declared and not used" lint when decimal is otherwise
	// unused — keep the import wired so future helpers (e.g.
	// zero-discount checks) can drop in cleanly.
	_ = decimal.Zero

	return &types.ListCouponUsageResp{
		List:     list,
		Total:    usageResp.Total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
