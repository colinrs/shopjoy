package brands

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type ToggleBrandPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewToggleBrandPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) ToggleBrandPageLogic {
	return ToggleBrandPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ToggleBrandPageLogic) ToggleBrandPage(req *types.ToggleBrandPageReq) (resp *types.BrandDetailResp, err error) {

	// Find brand
	brand, err := l.svcCtx.BrandRepo.FindByID(l.ctx, l.svcCtx.DB, req.ID)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, code.ErrBrandNotFound
	}

	// Toggle page
	brand.TogglePage(req.Enabled)
	brand.Model.UpdatedAt = time.Now().UTC()

	if err := l.svcCtx.BrandRepo.Update(l.ctx, l.svcCtx.DB, brand); err != nil {
		return nil, err
	}

	// Get product count
	productCount, err := l.svcCtx.BrandRepo.GetProductCount(l.ctx, l.svcCtx.DB, brand.ID)
	if err != nil {
		l.Logger.Errorf("failed to get product count for brand %d: %v", brand.ID, err)
		productCount = 0
	}

	return toBrandDetailResp(brand, productCount), nil
}
