// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	productApp "github.com/colinrs/shopjoy/shop/internal/application/product"

	"github.com/colinrs/shopjoy/shop/internal/svc"
	"github.com/colinrs/shopjoy/shop/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取商品列表
func NewListProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListProductLogic {
	return &ListProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListProductLogic) ListProduct(req *types.ListProductReq) (resp *types.ListProductResp, err error) {
	var minPrice, maxPrice *int64
	if req.MinPrice > 0 {
		minPrice = &req.MinPrice
	}
	if req.MaxPrice > 0 {
		maxPrice = &req.MaxPrice
	}

	appReq := productApp.QueryProductRequest{
		Name:       req.Name,
		CategoryID: req.CategoryID,
		Status:     req.Status,
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}

	result, err := l.svcCtx.ProductService.GetProductList(l.ctx, appReq)
	if err != nil {
		return nil, err
	}

	list := make([]*types.ProductDetailResp, len(result.List))
	for i, item := range result.List {
		list[i] = &types.ProductDetailResp{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			Currency:    item.Currency,
			CostPrice:   item.CostPrice,
			Stock:       item.Stock,
			Status:      item.Status,
			CategoryID:  item.CategoryID,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		}
	}

	return &types.ListProductResp{
		List:     list,
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.PageSize,
	}, nil
}
