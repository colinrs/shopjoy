package categories

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetCategoryMarketVisibilityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCategoryMarketVisibilityLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetCategoryMarketVisibilityLogic {
	return GetCategoryMarketVisibilityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCategoryMarketVisibilityLogic) GetCategoryMarketVisibility(req *types.GetCategoryMarketVisibilityReq) (resp *types.CategoryMarketVisibilityResp, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	// Verify category exists
	category, err := l.svcCtx.CategoryRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.CategoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, code.ErrCategoryNotFound
	}

	// Get market visibility
	markets, err := l.svcCtx.CategoryMarketRepo.FindByCategory(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.CategoryID)
	if err != nil {
		return nil, err
	}

	marketItems := make([]types.CategoryMarketItemResp, 0, len(markets))
	for _, m := range markets {
		marketItems = append(marketItems, types.CategoryMarketItemResp{
			MarketID:  m.MarketID,
			IsVisible: m.IsVisible,
		})
	}

	return &types.CategoryMarketVisibilityResp{
		CategoryID: req.CategoryID,
		Markets:    marketItems,
	}, nil
}
