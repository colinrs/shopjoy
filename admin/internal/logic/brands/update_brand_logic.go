package brands

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateBrandLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateBrandLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateBrandLogic {
	return UpdateBrandLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateBrandLogic) UpdateBrand(req *types.UpdateBrandReq) (resp *types.BrandDetailResp, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	// Find existing brand
	brand, err := l.svcCtx.BrandRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, code.ErrBrandNotFound
	}

	// Check for duplicate name if name changed
	if req.Name != brand.Name {
		existing, err := l.svcCtx.BrandRepo.FindByName(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.Name)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != req.ID {
			return nil, code.ErrBrandDuplicate
		}
	}

	// Update fields
	brand.Name = req.Name
	brand.Logo = req.Logo
	brand.Description = req.Description
	brand.Website = req.Website
	brand.TrademarkNumber = req.TrademarkNumber
	brand.TrademarkCountry = req.TrademarkCountry
	brand.EnablePage = req.EnablePage
	brand.Sort = req.Sort
	brand.Audit.UpdatedAt = time.Now().UTC()

	if err := l.svcCtx.BrandRepo.Update(l.ctx, l.svcCtx.DB, brand); err != nil {
		return nil, err
	}

	// Get product count
	productCount, _ := l.svcCtx.BrandRepo.GetProductCount(l.ctx, l.svcCtx.DB, brand.ID)

	return toBrandDetailResp(brand, productCount), nil
}

func toBrandDetailResp(brand *product.Brand, productCount int64) *types.BrandDetailResp {
	return &types.BrandDetailResp{
		ID:               brand.ID,
		Name:             brand.Name,
		Logo:             brand.Logo,
		Description:      brand.Description,
		Website:          brand.Website,
		TrademarkNumber:  brand.TrademarkNumber,
		TrademarkCountry: brand.TrademarkCountry,
		EnablePage:       brand.EnablePage,
		Sort:             brand.Sort,
		Status:           int8(brand.Status),
		ProductCount:     productCount,
		CreatedAt:        brand.Audit.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        brand.Audit.UpdatedAt.Format(time.RFC3339),
	}
}
