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

type CreateRedeemRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateRedeemRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateRedeemRuleLogic {
	return CreateRedeemRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRedeemRuleLogic) CreateRedeemRule(req *types.CreateRedeemRuleReq) (resp *types.RedeemRule, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)
	userID, _ := contextx.GetUserID(l.ctx)

	createReq := apppoints.CreateRedeemRuleRequest{
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
			createReq.StartAt = &t
		}
	}
	if req.EndAt != "" {
		t, err := time.Parse(time.RFC3339, req.EndAt)
		if err == nil {
			createReq.EndAt = &t
		}
	}

	rule, err := l.svcCtx.PointsService.CreateRedeemRule(l.ctx, shared.TenantID(tenantID), createReq, userID)
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
