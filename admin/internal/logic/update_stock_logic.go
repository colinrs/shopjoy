// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	productApp "github.com/colinrs/shopjoy/admin/internal/application/product"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateStockLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新库存
func NewUpdateStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStockLogic {
	return &UpdateStockLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateStockLogic) UpdateStock(req *types.UpdateStockReq) (resp *types.ProductDetailResp, err error) {
	appReq := productApp.UpdateStockRequest{
		ID:       req.ID,
		Quantity: req.Quantity,
	}

	if err := l.svcCtx.ProductService.UpdateStock(l.ctx, appReq); err != nil {
		return nil, err
	}

	result, err := l.svcCtx.ProductService.GetProduct(l.ctx, req.ID)
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
