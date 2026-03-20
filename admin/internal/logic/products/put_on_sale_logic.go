package products

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

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
	productResp, err := l.svcCtx.ProductService.PutOnSale(l.ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return convertToProductDetailResp(productResp), nil
}