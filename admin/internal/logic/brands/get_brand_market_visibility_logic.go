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

type GetBrandMarketVisibilityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBrandMarketVisibilityLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetBrandMarketVisibilityLogic {
	return GetBrandMarketVisibilityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBrandMarketVisibilityLogic) GetBrandMarketVisibility(req *types.GetBrandMarketVisibilityReq) (resp *types.BrandMarketVisibilityResp, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	// Verify brand exists
	brand, err := l.svcCtx.BrandRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.BrandID)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, code.ErrBrandNotFound
	}

	// Get market visibility
	markets, err := l.svcCtx.BrandMarketRepo.FindByBrand(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.BrandID)
	if err != nil {
		return nil, err
	}

	marketItems := make([]types.BrandMarketItemResp, 0, len(markets))
	for _, m := range markets {
		marketItems = append(marketItems, types.BrandMarketItemResp{
			MarketID:  m.MarketID,
			IsVisible: m.IsVisible,
		})
	}

	return &types.BrandMarketVisibilityResp{
		BrandID: req.BrandID,
		Markets: marketItems,
	}, nil
}
