package redeem_rules

import (
	"context"
	"time"

	apppoints "github.com/colinrs/shopjoy/admin/internal/application/points"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRedeemRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRedeemRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateRedeemRuleLogic {
	return UpdateRedeemRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRedeemRuleLogic) UpdateRedeemRule(req *types.UpdateRedeemRuleReq) (resp *types.RedeemRule, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)
	userID, _ := contextx.GetUserID(l.ctx)

	updateReq := apppoints.UpdateRedeemRuleRequest{
		ID:             req.ID,
		Name:           req.Name,
		Description:    req.Description,
		CouponID:       req.CouponID,
		PointsRequired: req.PointsRequired,
		TotalStock:     req.TotalStock,
		PerUserLimit:   req.PerUserLimit,
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

	rule, err := l.svcCtx.PointsService.UpdateRedeemRule(l.ctx, shared.TenantID(tenantID), updateReq, userID)
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
