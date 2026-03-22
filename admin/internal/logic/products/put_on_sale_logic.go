package products

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type PutOnSaleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPutOnSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) PutOnSaleLogic {
	return PutOnSaleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PutOnSaleLogic) PutOnSale(req *types.PutOnSaleReq) (resp *types.ProductDetailResp, err error) {
	// 从 context 获取 tenantID
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// 平台管理员设置 tenantID = 0 以访问所有数据
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	productResp, err := l.svcCtx.ProductService.PutOnSale(l.ctx, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}

	return convertToProductDetailResp(productResp), nil
}
