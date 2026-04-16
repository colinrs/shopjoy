package products

import (
	"context"

	appProduct "github.com/colinrs/shopjoy/admin/internal/application/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateStockLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateStockLogic {
	return UpdateStockLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateStockLogic) UpdateStock(req *types.UpdateStockReq) (resp *types.ProductDetailResp, err error) {
	// 从 context 获取 tenantID
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok && !contextx.IsPlatformAdmin(l.ctx) {
		return nil, code.ErrUnauthorized
	}
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	stockReq := appProduct.UpdateStockRequest{
		ID:       req.ID,
		Quantity: req.Quantity,
	}

	if err := l.svcCtx.ProductService.UpdateStock(l.ctx, shared.TenantID(tenantID), stockReq); err != nil {
		return nil, err
	}

	productResp, err := l.svcCtx.ProductService.GetProduct(l.ctx, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}

	return convertToProductDetailResp(productResp), nil
}
