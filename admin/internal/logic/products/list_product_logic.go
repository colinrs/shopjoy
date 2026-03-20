package products

import (
	"context"

	appProduct "github.com/colinrs/shopjoy/admin/internal/application/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListProductLogic {
	return ListProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListProductLogic) ListProduct(req *types.ListProductReq) (resp *types.ListProductResp, err error) {
	queryReq := appProduct.QueryProductRequest{
		Name:       req.Name,
		CategoryID: req.CategoryID,
		Status:     req.Status,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}

	if req.MinPrice > 0 {
		minPrice := req.MinPrice
		queryReq.MinPrice = &minPrice
	}
	if req.MaxPrice > 0 {
		maxPrice := req.MaxPrice
		queryReq.MaxPrice = &maxPrice
	}

	listResp, err := l.svcCtx.ProductService.GetProductList(l.ctx, queryReq)
	if err != nil {
		return nil, err
	}

	list := make([]*types.ProductDetailResp, len(listResp.List))
	for i, p := range listResp.List {
		list[i] = convertToProductDetailResp(p)
	}

	return &types.ListProductResp{
		List:     list,
		Total:    listResp.Total,
		Page:     listResp.Page,
		PageSize: listResp.PageSize,
	}, nil
}