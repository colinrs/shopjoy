package products

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetProductLogic {
	return GetProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductLogic) GetProduct(req *types.GetProductReq) (resp *types.ProductDetailResp, err error) {
	// 从 context 获取 tenantID

	productResp, err := l.svcCtx.ProductService.GetProduct(l.ctx, req.ID)
	if err != nil {
		return nil, err
	}

	resp = convertToProductDetailResp(productResp)
	resp.CategoryPath = buildCategoryPath(l.ctx, l.svcCtx.DB, productResp.CategoryID)
	return resp, nil
}
