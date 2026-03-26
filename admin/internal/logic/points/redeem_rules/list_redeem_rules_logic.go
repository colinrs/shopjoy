package redeem_rules

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/points"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRedeemRulesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRedeemRulesLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListRedeemRulesLogic {
	return ListRedeemRulesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRedeemRulesLogic) ListRedeemRules(req *types.ListRedeemRulesReq) (resp *types.ListRedeemRulesResp, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	query := points.RedeemRuleQuery{
		PageQuery: shared.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		Name: req.Name,
	}

	if req.Status != "" {
		switch req.Status {
		case "active":
			query.Status = points.RedeemRuleStatusActive
		case "inactive":
			query.Status = points.RedeemRuleStatusInactive
		}
	}

	rules, total, stats, err := l.svcCtx.PointsService.ListRedeemRules(l.ctx, shared.TenantID(tenantID), query)
	if err != nil {
		return nil, err
	}

	list := make([]*types.RedeemRule, len(rules))
	for i, r := range rules {
		list[i] = &types.RedeemRule{
			ID:             r.ID,
			Name:           r.Name,
			Description:    r.Description,
			CouponID:       r.CouponID,
			CouponName:     "", // Would need to fetch from coupon service
			PointsRequired: r.PointsRequired,
			TotalStock:     r.TotalStock,
			UsedStock:      r.UsedStock,
			PerUserLimit:   r.PerUserLimit,
			Status:         r.Status,
			StartAt:        formatTimePtrFromTime(r.StartAt),
			EndAt:          formatTimePtrFromTime(r.EndAt),
			CreatedAt:      r.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:      r.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &types.ListRedeemRulesResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Stats: types.RedeemRulesStats{
			Total:         stats.Total,
			Active:        stats.Active,
			TotalRedeemed: stats.TotalRedeemed,
		},
	}, nil
}

func formatTimePtrFromTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}