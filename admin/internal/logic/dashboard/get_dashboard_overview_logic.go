package dashboard

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDashboardOverviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDashboardOverviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetDashboardOverviewLogic {
	return GetDashboardOverviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDashboardOverviewLogic) GetDashboardOverview(req *types.DashboardOverviewRequest) (resp *types.DashboardOverviewResponse, err error) {
	helper := NewDashboardHelper(l.ctx, l.svcCtx)
	tenantID := helper.GetTenantID()
	return helper.GetOverview(tenantID)
}
