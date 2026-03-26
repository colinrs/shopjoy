package stats

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetExpiringPointsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetExpiringPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetExpiringPointsLogic {
	return GetExpiringPointsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetExpiringPointsLogic) GetExpiringPoints(req *types.GetExpiringPointsReq) (resp *types.ExpiringPointsResp, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	days := req.Days
	if days <= 0 {
		days = 30
	}
	if days > 365 {
		days = 365
	}

	expiringData, totalPoints, err := l.svcCtx.PointsService.GetExpiringPoints(l.ctx, shared.TenantID(tenantID), days)
	if err != nil {
		return nil, err
	}

	list := make([]types.ExpiringPoints, len(expiringData))
	for i, d := range expiringData {
		list[i] = types.ExpiringPoints{
			Date:      d.Date,
			Points:    d.Points,
			UserCount: d.UserCount,
		}
	}

	return &types.ExpiringPointsResp{
		List:        list,
		TotalPoints: totalPoints,
	}, nil
}