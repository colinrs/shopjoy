package dashboard

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTopProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTopProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetTopProductsLogic {
	return GetTopProductsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTopProductsLogic) GetTopProducts(req *types.TopProductsRequest) (resp *types.TopProductsResponse, err error) {
	helper := NewDashboardHelper(l.ctx, l.svcCtx)
	tenantID := helper.GetTenantID()

	limit := req.Limit
	if limit <= 0 || limit > 20 {
		limit = 5
	}

	return helper.GetTopProducts(tenantID, limit, req.Period)
}