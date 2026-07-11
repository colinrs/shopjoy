package users

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetUserStatsLogic {
	return GetUserStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserStatsLogic) GetUserStats(req *types.UserStatsRequest) (resp *types.UserStatsResponse, err error) {
	stats, err := l.svcCtx.UserService.GetStats(l.ctx)
	if err != nil {
		return nil, err
	}

	return &types.UserStatsResponse{
		Total:     stats.Total,
		Active:    stats.Active,
		Suspended: stats.Suspended,
		NewToday:  stats.NewToday,
	}, nil
}
