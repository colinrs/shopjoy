package users

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/colinrs/shopjoy/pkg/tenant"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserStatsEnhancedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserStatsEnhancedLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetUserStatsEnhancedLogic {
	return GetUserStatsEnhancedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserStatsEnhancedLogic) GetUserStatsEnhanced(req *types.UserStatsRequest) (resp *types.UserStatsEnhancedResponse, err error) {
	tenantID, ok := tenant.FromContext(l.ctx)
	if !ok {
		tenantID = shared.TenantID(1) // 默认租户
	}

	stats, err := l.svcCtx.UserService.GetUserStats(l.ctx, tenantID)
	if err != nil {
		return nil, err
	}

	return &types.UserStatsEnhancedResponse{
		TotalUsers:     stats.TotalUsers,
		ActiveUsers:    stats.ActiveUsers,
		SuspendedUsers: stats.SuspendedUsers,
		NewUsersToday:  stats.NewUsersToday,
	}, nil
}
