package products

import (
	"context"

	appProduct "github.com/colinrs/shopjoy/admin/internal/application/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

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
	stockReq := appProduct.UpdateStockRequest{
		ID:       req.ID,
		Quantity: req.Quantity,
	}

	if err := l.svcCtx.ProductService.UpdateStock(l.ctx, stockReq); err != nil {
		return nil, err
	}

	productResp, err := l.svcCtx.ProductService.GetProduct(l.ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return convertToProductDetailResp(productResp), nil
}
