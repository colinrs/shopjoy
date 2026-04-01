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

type DeleteBrandLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteBrandLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteBrandLogic {
	return DeleteBrandLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteBrandLogic) DeleteBrand(req *types.GetBrandReq) (resp *types.CreateBrandResp, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	// Find brand
	brand, err := l.svcCtx.BrandRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, code.ErrBrandNotFound
	}

	// Check if brand has products
	productCount, err := l.svcCtx.BrandRepo.GetProductCount(l.ctx, l.svcCtx.DB, req.ID)
	if err != nil {
		return nil, err
	}
	if productCount > 0 {
		return nil, code.ErrBrandHasProducts
	}

	// Delete brand
	if err := l.svcCtx.BrandRepo.Delete(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID); err != nil {
		return nil, err
	}

	return &types.CreateBrandResp{
		ID: req.ID,
	}, nil
}
