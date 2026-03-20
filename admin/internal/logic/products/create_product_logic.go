package products

import (
	"context"

	appProduct "github.com/colinrs/shopjoy/admin/internal/application/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateProductLogic {
	return CreateProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateProductLogic) CreateProduct(req *types.CreateProductReq) (resp *types.CreateProductResp, err error) {
	createReq := appProduct.CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Currency:    req.Currency,
		CostPrice:   req.CostPrice,
		CategoryID:  req.CategoryID,
	}

	productResp, err := l.svcCtx.ProductService.CreateProduct(l.ctx, createReq)
	if err != nil {
		return nil, err
	}

	return &types.CreateProductResp{
		ID: productResp.ID,
	}, nil
}