package earn_rules

import (
	"context"
	"time"

	apppoints "github.com/colinrs/shopjoy/admin/internal/application/points"
	"github.com/colinrs/shopjoy/admin/internal/domain/points"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/shopspring/decimal"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateEarnRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateEarnRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateEarnRuleLogic {
	return UpdateEarnRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateEarnRuleLogic) UpdateEarnRule(req *types.UpdateEarnRuleReq) (resp *types.EarnRule, err error) {
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok && !contextx.IsPlatformAdmin(l.ctx) {
		return nil, code.ErrUnauthorized
	}
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}
	userID, _ := contextx.GetUserID(l.ctx)

	updateReq := apppoints.UpdateEarnRuleRequest{
		ID:               req.ID,
		Name:             req.Name,
		Description:      req.Description,
		Scenario:         points.EarnScenario(req.Scenario),
		CalculationType:  points.CalculationType(req.CalculationType),
		FixedPoints:      req.FixedPoints,
		Ratio:            decimal.RequireFromString(req.Ratio),
		Tiers:            parseTiers(req.Tiers),
		ConditionType:    points.ConditionType(req.ConditionType),
		ConditionValue:   req.ConditionValue,
		ExpirationMonths: req.ExpirationMonths,
		Priority:         req.Priority,
	}

	if req.StartAt != "" {
		t, err := time.Parse(time.RFC3339, req.StartAt)
		if err == nil {
			updateReq.StartAt = &t
		}
	}
	if req.EndAt != "" {
		t, err := time.Parse(time.RFC3339, req.EndAt)
		if err == nil {
			updateReq.EndAt = &t
		}
	}

	rule, err := l.svcCtx.PointsService.UpdateEarnRule(l.ctx, shared.TenantID(tenantID), updateReq, userID)
	if err != nil {
		return nil, err
	}

	return &types.EarnRule{
		ID:               rule.ID,
		Name:             rule.Name,
		Description:      rule.Description,
		Scenario:         rule.Scenario,
		CalculationType:  rule.CalculationType,
		FixedPoints:      rule.FixedPoints,
		Ratio:            rule.Ratio.String(),
		Tiers:            convertTiers(rule.Tiers),
		ConditionType:    rule.ConditionType,
		ConditionValue:   rule.ConditionValue,
		ExpirationMonths: rule.ExpirationMonths,
		Status:           rule.Status,
		Priority:         rule.Priority,
		StartAt:          formatTimePtrToStr(rule.StartAt),
		EndAt:            formatTimePtrToStr(rule.EndAt),
		CreatedAt:        rule.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        rule.UpdatedAt.Format(time.RFC3339),
	}, nil
}
