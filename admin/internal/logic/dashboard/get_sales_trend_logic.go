package dashboard

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSalesTrendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSalesTrendLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetSalesTrendLogic {
	return GetSalesTrendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSalesTrendLogic) GetSalesTrend(req *types.SalesTrendRequest) (resp *types.SalesTrendResponse, err error) {
	helper := NewDashboardHelper(l.ctx, l.svcCtx)
	tenantID, ok := helper.GetTenantID()
	if !ok {
		return nil, code.ErrUnauthorized
	}
	return helper.GetSalesTrend(tenantID, req.Period)
}
