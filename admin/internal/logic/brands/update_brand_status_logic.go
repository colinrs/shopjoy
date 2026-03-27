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

type UpdateBrandStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateBrandStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateBrandStatusLogic {
	return UpdateBrandStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateBrandStatusLogic) UpdateBrandStatus(req *types.UpdateBrandStatusReq) (resp *types.BrandDetailResp, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Find brand
	brand, err := l.svcCtx.BrandRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, code.ErrBrandNotFound
	}

	// Update status
	if req.Status == 1 {
		brand.Enable()
	} else {
		brand.Disable()
	}
	brand.Audit.UpdatedAt = time.Now().UTC()

	if err := l.svcCtx.BrandRepo.Update(l.ctx, l.svcCtx.DB, brand); err != nil {
		return nil, err
	}

	// Get product count
	productCount, _ := l.svcCtx.BrandRepo.GetProductCount(l.ctx, l.svcCtx.DB, brand.ID)

	return toBrandDetailResp(brand, productCount), nil
}
