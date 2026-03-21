package brands

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetBrandLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBrandLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetBrandLogic {
	return GetBrandLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBrandLogic) GetBrand(req *types.GetBrandReq) (resp *types.BrandDetailResp, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)

	brand, err := l.svcCtx.BrandRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, code.ErrBrandNotFound
	}

	// Get product count
	productCount, _ := l.svcCtx.BrandRepo.GetProductCount(l.ctx, l.svcCtx.DB, brand.ID)

	return toBrandDetailResp(brand, productCount), nil
}
