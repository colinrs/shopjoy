package stats

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTopUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTopUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetTopUsersLogic {
	return GetTopUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTopUsersLogic) GetTopUsers(req *types.GetTopUsersReq) (resp *types.TopUsersResp, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Calculate time range based on period
	var startTime, endTime time.Time
	now := time.Now().UTC()
	endTime = now

	switch req.Period {
	case "7d":
		startTime = now.AddDate(0, 0, -7)
	case "30d":
		startTime = now.AddDate(0, 0, -30)
	case "90d":
		startTime = now.AddDate(0, 0, -90)
	case "1y":
		startTime = now.AddDate(-1, 0, 0)
	default:
		startTime = now.AddDate(0, 0, -7)
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	topUsers, err := l.svcCtx.PointsService.GetTopUsers(l.ctx, shared.TenantID(tenantID), startTime, endTime, limit)
	if err != nil {
		return nil, err
	}

	users := make([]types.TopUser, len(topUsers))
	for i, u := range topUsers {
		users[i] = types.TopUser{
			UserID:       u.UserID,
			UserEmail:    "", // Would need to fetch from user service
			PointsEarned: u.PointsEarned,
			CreatedAt:    time.Now().Format("2006-01-02 15:04:05"),
		}
	}

	return &types.TopUsersResp{
		List:   users,
		Period: req.Period,
	}, nil
}