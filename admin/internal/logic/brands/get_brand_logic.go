package brands

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
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

	brand, err := l.svcCtx.BrandRepo.FindByID(l.ctx, l.svcCtx.DB, req.ID)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, code.ErrBrandNotFound
	}

	// Get product count
	productCount, err := l.svcCtx.BrandRepo.GetProductCount(l.ctx, l.svcCtx.DB, brand.ID)
	if err != nil {
		l.Logger.Errorf("failed to get product count for brand %d: %v", brand.ID, err)
		productCount = 0
	}

	return toBrandDetailResp(brand, productCount), nil
}
