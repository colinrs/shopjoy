package earn_rules

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type ActivateEarnRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewActivateEarnRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) ActivateEarnRuleLogic {
	return ActivateEarnRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ActivateEarnRuleLogic) ActivateEarnRule(req *types.ActivateEarnRuleReq) (resp *types.EarnRule, err error) {
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok && !contextx.IsPlatformAdmin(l.ctx) {
		return nil, code.ErrUnauthorized
	}
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}
	userID, _ := contextx.GetUserID(l.ctx)

	if err := l.svcCtx.PointsService.ActivateEarnRule(l.ctx, shared.TenantID(tenantID), req.ID, userID); err != nil {
		return nil, err
	}

	rule, err := l.svcCtx.PointsService.GetEarnRule(l.ctx, shared.TenantID(tenantID), req.ID)
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
