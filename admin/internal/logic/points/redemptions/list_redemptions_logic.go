package redemptions

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

type ListRedemptionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRedemptionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListRedemptionsLogic {
	return ListRedemptionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRedemptionsLogic) ListRedemptions(req *types.ListRedemptionsReq) (resp *types.ListRedemptionsResp, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	query := points.PointsRedemptionQuery{
		PageQuery: shared.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		UserID: req.UserID,
	}

	if req.Status != "" {
		switch req.Status {
		case "pending":
			query.Status = points.RedemptionStatusPending
		case "completed":
			query.Status = points.RedemptionStatusCompleted
		case "cancelled":
			query.Status = points.RedemptionStatusCancelled
		}
	}

	if req.StartTime != "" {
		t, err := time.Parse(time.RFC3339, req.StartTime)
		if err == nil {
			query.StartTime = &t
		}
	}
	if req.EndTime != "" {
		t, err := time.Parse(time.RFC3339, req.EndTime)
		if err == nil {
			query.EndTime = &t
		}
	}

	redemptions, total, err := l.svcCtx.PointsService.ListRedemptions(l.ctx, shared.TenantID(tenantID), query)
	if err != nil {
		return nil, err
	}

	list := make([]*types.PointsRedemption, len(redemptions))
	for i, r := range redemptions {
		list[i] = &types.PointsRedemption{
			ID:           r.ID,
			UserID:       r.UserID,
			RedeemRuleID: r.RedeemRuleID,
			CouponID:     r.CouponID,
			CouponName:   "",
			UserCouponID: r.UserCouponID,
			PointsUsed:   r.PointsUsed,
			Status:       r.Status,
			CreatedAt:    r.CreatedAt.Format(time.RFC3339),
			CompletedAt:  formatTimePtrFromTime(r.CompletedAt),
		}
	}

	return &types.ListRedemptionsResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func formatTimePtrFromTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}