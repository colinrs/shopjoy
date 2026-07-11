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
	// 从 context 获取 tenantID

	// 解析价格字符串（单位：元）
	price, err := appProduct.ToDomainMoneyFromString(req.Price, req.Currency)
	if err != nil {
		return nil, err
	}
	costPrice, err := appProduct.ToDomainMoneyFromString(req.CostPrice, req.Currency)
	if err != nil {
		return nil, err
	}

	createReq := appProduct.CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       price.Amount,
		Currency:    req.Currency,
		CostPrice:   costPrice.Amount,
		CategoryID:  req.CategoryID,
		Stock:       req.Stock,
		SKU:         req.SKU,
	}

	productResp, err := l.svcCtx.ProductService.CreateProduct(l.ctx, createReq)
	if err != nil {
		return nil, err
	}

	return &types.CreateProductResp{
		ID: productResp.ID,
	}, nil
}
