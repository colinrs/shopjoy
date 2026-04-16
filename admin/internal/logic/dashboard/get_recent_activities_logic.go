package dashboard

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRecentActivitiesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRecentActivitiesLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetRecentActivitiesLogic {
	return GetRecentActivitiesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRecentActivitiesLogic) GetRecentActivities(req *types.RecentActivitiesRequest) (resp *types.RecentActivitiesResponse, err error) {
	helper := NewDashboardHelper(l.ctx, l.svcCtx)
	tenantID, ok := helper.GetTenantID()
	if !ok {
		return nil, code.ErrUnauthorized
	}

	limit := req.Limit
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	return helper.GetRecentActivities(tenantID, limit)
}
