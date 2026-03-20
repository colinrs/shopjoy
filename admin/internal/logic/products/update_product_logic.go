package products

import (
	"context"

	appProduct "github.com/colinrs/shopjoy/admin/internal/application/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateProductLogic {
	return UpdateProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProductLogic) UpdateProduct(req *types.UpdateProductReq) (resp *types.ProductDetailResp, err error) {
	updateReq := appProduct.UpdateProductRequest{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Currency:    req.Currency,
		CategoryID:  req.CategoryID,
	}

	productResp, err := l.svcCtx.ProductService.UpdateProduct(l.ctx, updateReq)
	if err != nil {
		return nil, err
	}

	return convertToProductDetailResp(productResp), nil
}

func convertToProductDetailResp(p *appProduct.ProductResponse) *types.ProductDetailResp {
	return &types.ProductDetailResp{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Currency:    p.Currency,
		CostPrice:   p.CostPrice,
		Stock:       p.Stock,
		Status:      p.Status,
		CategoryID:  p.CategoryID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}