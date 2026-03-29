package products

import (
	"context"

	appProduct "github.com/colinrs/shopjoy/admin/internal/application/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

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
	// 从 context 获取 tenantID
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// 平台管理员设置 tenantID = 0 以访问所有数据
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// 解析价格字符串（单位：元）
	price, err := appProduct.ToDomainMoneyFromString(req.Price, req.Currency)
	if err != nil {
		return nil, err
	}

	updateReq := appProduct.UpdateProductRequest{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Price:       price.Amount,
		Currency:    req.Currency,
		CategoryID:  req.CategoryID,
	}

	productResp, err := l.svcCtx.ProductService.UpdateProduct(l.ctx, shared.TenantID(tenantID), updateReq)
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
		Price:       p.Price.String(),
		Currency:    p.Currency,
		CostPrice:   p.CostPrice.String(),
		Stock:       p.Stock,
		Status:      p.Status,
		CategoryID:  p.CategoryID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
