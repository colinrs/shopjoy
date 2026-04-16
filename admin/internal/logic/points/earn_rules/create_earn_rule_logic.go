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

type CreateEarnRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateEarnRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateEarnRuleLogic {
	return CreateEarnRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateEarnRuleLogic) CreateEarnRule(req *types.CreateEarnRuleReq) (resp *types.EarnRule, err error) {
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok && !contextx.IsPlatformAdmin(l.ctx) {
		return nil, code.ErrUnauthorized
	}
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}
	userID, _ := contextx.GetUserID(l.ctx)

	createReq := apppoints.CreateEarnRuleRequest{
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
			createReq.StartAt = &t
		}
	}
	if req.EndAt != "" {
		t, err := time.Parse(time.RFC3339, req.EndAt)
		if err == nil {
			createReq.EndAt = &t
		}
	}

	rule, err := l.svcCtx.PointsService.CreateEarnRule(l.ctx, shared.TenantID(tenantID), createReq, userID)
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

func formatTimePtrToStr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}
