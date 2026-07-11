package products

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteProductLogic {
	return DeleteProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteProductLogic) DeleteProduct(req *types.GetProductReq) (resp *types.CreateProductResp, err error) {
	// Get tenant ID from context

	// Delete product
	if err := l.svcCtx.ProductService.DeleteProduct(l.ctx, req.ID); err != nil {
		return nil, err
	}

	return &types.CreateProductResp{ID: req.ID}, nil
}
