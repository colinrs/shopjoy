// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"github.com/colinrs/shopjoy/shop/internal/svc"
	"github.com/colinrs/shopjoy/shop/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PutOnSaleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 上架商品
func NewPutOnSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PutOnSaleLogic {
	return &PutOnSaleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PutOnSaleLogic) PutOnSale(req *types.PutOnSaleReq) (resp *types.ProductDetailResp, err error) {
	result, err := l.svcCtx.ProductService.PutOnSale(l.ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &types.ProductDetailResp{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description,
		Price:       result.Price,
		Currency:    result.Currency,
		CostPrice:   result.CostPrice,
		Stock:       result.Stock,
		Status:      result.Status,
		CategoryID:  result.CategoryID,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}, nil
}
