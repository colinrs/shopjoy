package redeem_rules

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRedeemRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRedeemRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetRedeemRuleLogic {
	return GetRedeemRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRedeemRuleLogic) GetRedeemRule(req *types.GetRedeemRuleReq) (resp *types.RedeemRule, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
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
