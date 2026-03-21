package brands

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
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
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Find brand
	brand, err := l.svcCtx.BrandRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, code.ErrBrandNotFound
	}

	// Toggle page
	brand.TogglePage(req.Enabled)
	brand.Audit.UpdatedAt = time.Now()

	if err := l.svcCtx.BrandRepo.Update(l.ctx, l.svcCtx.DB, brand); err != nil {
		return nil, err
	}

	// Get product count
	productCount, _ := l.svcCtx.BrandRepo.GetProductCount(l.ctx, l.svcCtx.DB, brand.ID)

	return toBrandDetailResp(brand, productCount), nil
}
