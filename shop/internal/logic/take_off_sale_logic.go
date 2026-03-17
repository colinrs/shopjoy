// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"github.com/colinrs/shopjoy/shop/internal/svc"
	"github.com/colinrs/shopjoy/shop/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TakeOffSaleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 下架商品
func NewTakeOffSaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TakeOffSaleLogic {
	return &TakeOffSaleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TakeOffSaleLogic) TakeOffSale(req *types.TakeOffSaleReq) (resp *types.ProductDetailResp, err error) {
	result, err := l.svcCtx.ProductService.TakeOffSale(l.ctx, req.ID)
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
