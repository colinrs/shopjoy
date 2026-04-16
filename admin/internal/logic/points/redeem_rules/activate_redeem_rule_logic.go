package redeem_rules

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

type ActivateRedeemRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewActivateRedeemRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) ActivateRedeemRuleLogic {
	return ActivateRedeemRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ActivateRedeemRuleLogic) ActivateRedeemRule(req *types.ActivateRedeemRuleReq) (resp *types.RedeemRule, err error) {
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok && !contextx.IsPlatformAdmin(l.ctx) {
		return nil, code.ErrUnauthorized
	}
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}
	userID, _ := contextx.GetUserID(l.ctx)

	if err := l.svcCtx.PointsService.ActivateRedeemRule(l.ctx, shared.TenantID(tenantID), req.ID, userID); err != nil {
		return nil, err
	}

	rule, err := l.svcCtx.PointsService.GetRedeemRule(l.ctx, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}

	return &types.RedeemRule{
		ID:             rule.ID,
		Name:           rule.Name,
		Description:    rule.Description,
		CouponID:       rule.CouponID,
		CouponName:     "",
		PointsRequired: rule.PointsRequired,
		TotalStock:     rule.TotalStock,
		UsedStock:      rule.UsedStock,
		PerUserLimit:   rule.PerUserLimit,
		Status:         rule.Status,
		StartAt:        formatTimePtrFromTime(rule.StartAt),
		EndAt:          formatTimePtrFromTime(rule.EndAt),
		CreatedAt:      rule.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      rule.UpdatedAt.Format(time.RFC3339),
	}, nil
}
