package redemptions

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRedemptionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRedemptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetRedemptionLogic {
	return GetRedemptionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRedemptionLogic) GetRedemption(req *types.GetRedemptionReq) (resp *types.PointsRedemption, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	redemption, err := l.svcCtx.PointsService.GetRedemption(l.ctx, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}

	return &types.PointsRedemption{
		ID:           redemption.ID,
		UserID:       redemption.UserID,
		RedeemRuleID: redemption.RedeemRuleID,
		CouponID:     redemption.CouponID,
		CouponName:   "",
		UserCouponID: redemption.UserCouponID,
		PointsUsed:   redemption.PointsUsed,
		Status:       redemption.Status,
		CreatedAt:    redemption.CreatedAt.Format(time.RFC3339),
		CompletedAt:  formatTimePtrFromTime(redemption.CompletedAt),
	}, nil
}