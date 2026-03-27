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

type SetBrandMarketVisibilityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetBrandMarketVisibilityLogic(ctx context.Context, svcCtx *svc.ServiceContext) SetBrandMarketVisibilityLogic {
	return SetBrandMarketVisibilityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetBrandMarketVisibilityLogic) SetBrandMarketVisibility(req *types.SetBrandMarketVisibilityReq) (resp *types.CreateBrandResp, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Verify brand exists
	brand, err := l.svcCtx.BrandRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.BrandID)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, code.ErrBrandNotFound
	}

	// Delete existing market visibility for this brand
	if err := l.svcCtx.BrandMarketRepo.DeleteByBrand(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.BrandID); err != nil {
		return nil, err
	}

	// Create new market visibility entries
	if len(req.MarketIDs) > 0 {
		items := make([]*product.BrandMarket, 0, len(req.MarketIDs))
		now := time.Now().UTC()
		for _, marketID := range req.MarketIDs {
			id, _ := l.svcCtx.IDGen.NextID(l.ctx)
			items = append(items, &product.BrandMarket{
				ID:        id,
				TenantID:  shared.TenantID(tenantID),
				BrandID:   req.BrandID,
				MarketID:  marketID,
				IsVisible: req.Visible,
				Audit: shared.AuditInfo{
					CreatedAt: now,
					UpdatedAt: now,
				},
			})
		}
		if err := l.svcCtx.BrandMarketRepo.BatchCreate(l.ctx, l.svcCtx.DB, items); err != nil {
			return nil, err
		}
	}

	return &types.CreateBrandResp{
		ID: req.BrandID,
	}, nil
}
