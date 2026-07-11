package brands

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListBrandsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListBrandsLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListBrandsLogic {
	return ListBrandsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListBrandsLogic) ListBrands(req *types.ListBrandReq) (resp *types.ListBrandResp, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	query := product.BrandQuery{
		PageQuery: shared.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		Name:   req.Name,
		Status: shared.Status(req.Status),
	}

	brands, total, err := l.svcCtx.BrandRepo.FindList(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), query)
	if err != nil {
		return nil, err
	}

	list := make([]types.BrandDetailResp, 0, len(brands))
	for _, b := range brands {
		productCount, err := l.svcCtx.BrandRepo.GetProductCount(l.ctx, l.svcCtx.DB, b.ID)
		if err != nil {
			l.Logger.Errorf("failed to get product count for brand %d: %v", b.ID, err)
			productCount = 0
		}
		list = append(list, *toBrandDetailResp(b, productCount))
	}

	return &types.ListBrandResp{
		List:  list,
		Total: total,
	}, nil
}
