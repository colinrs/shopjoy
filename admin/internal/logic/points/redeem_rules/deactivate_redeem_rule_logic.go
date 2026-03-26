package redeem_rules

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeactivateRedeemRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeactivateRedeemRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeactivateRedeemRuleLogic {
	return DeactivateRedeemRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeactivateRedeemRuleLogic) DeactivateRedeemRule(req *types.DeactivateRedeemRuleReq) (resp *types.RedeemRule, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)
	userID, _ := contextx.GetUserID(l.ctx)

	if err := l.svcCtx.PointsService.DeactivateRedeemRule(l.ctx, shared.TenantID(tenantID), req.ID, userID); err != nil {
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
		CreatedAt:      rule.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      rule.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}