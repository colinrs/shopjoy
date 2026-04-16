package earn_rules

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/points"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListEarnRulesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListEarnRulesLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListEarnRulesLogic {
	return ListEarnRulesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListEarnRulesLogic) ListEarnRules(req *types.ListEarnRulesReq) (resp *types.ListEarnRulesResp, err error) {
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok && !contextx.IsPlatformAdmin(l.ctx) {
		return nil, code.ErrUnauthorized
	}
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	query := points.EarnRuleQuery{
		PageQuery: shared.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		Name:            req.Name,
		Scenario:        points.EarnScenario(req.Scenario),
		CalculationType: points.CalculationType(req.CalculationType),
	}

	if req.Status != "" {
		switch req.Status {
		case "draft":
			query.Status = points.EarnRuleStatusDraft
		case "active":
			query.Status = points.EarnRuleStatusActive
		case "inactive":
			query.Status = points.EarnRuleStatusInactive
		}
	}

	rules, total, stats, err := l.svcCtx.PointsService.ListEarnRules(l.ctx, shared.TenantID(tenantID), query)
	if err != nil {
		return nil, err
	}

	list := make([]*types.EarnRule, len(rules))
	for i, r := range rules {
		list[i] = &types.EarnRule{
			ID:               r.ID,
			Name:             r.Name,
			Description:      r.Description,
			Scenario:         r.Scenario,
			CalculationType:  r.CalculationType,
			FixedPoints:      r.FixedPoints,
			Ratio:            r.Ratio.String(),
			Tiers:            convertTiers(r.Tiers),
			ConditionType:    r.ConditionType,
			ConditionValue:   r.ConditionValue,
			ExpirationMonths: r.ExpirationMonths,
			Status:           r.Status,
			Priority:         r.Priority,
			StartAt:          formatTimePtrToStr(r.StartAt),
			EndAt:            formatTimePtrToStr(r.EndAt),
			CreatedAt:        r.CreatedAt.Format(time.RFC3339),
			UpdatedAt:        r.UpdatedAt.Format(time.RFC3339),
		}
	}

	return &types.ListEarnRulesResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Stats: types.EarnRulesStats{
			Total:  stats.Total,
			Active: stats.Active,
		},
	}, nil
}

func convertTiers(tiers points.TierConfigs) []*types.TierConfig {
	if tiers == nil {
		return nil
	}
	result := make([]*types.TierConfig, len(tiers))
	for i, t := range tiers {
		result[i] = &types.TierConfig{
			Threshold: t.Threshold,
			Ratio:     t.Ratio.String(),
		}
	}
	return result
}

func parseTiers(tiers []*types.TierConfig) points.TierConfigs {
	if tiers == nil {
		return nil
	}
	result := make(points.TierConfigs, len(tiers))
	for i, t := range tiers {
		ratio, _ := decimal.NewFromString(t.Ratio)
		result[i] = points.TierConfig{
			Threshold: t.Threshold,
			Ratio:     ratio,
		}
	}
	return result
}
