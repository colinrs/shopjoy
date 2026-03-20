package products

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TakeOffSaleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTakeOffSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) TakeOffSaleLogic {
	return TakeOffSaleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TakeOffSaleLogic) TakeOffSale(req *types.TakeOffSaleReq) (resp *types.ProductDetailResp, err error) {
	productResp, err := l.svcCtx.ProductService.TakeOffSale(l.ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return convertToProductDetailResp(productResp), nil
}