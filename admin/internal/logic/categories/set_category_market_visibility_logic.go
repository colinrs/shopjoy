package categories

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type SetCategoryMarketVisibilityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetCategoryMarketVisibilityLogic(ctx context.Context, svcCtx *svc.ServiceContext) SetCategoryMarketVisibilityLogic {
	return SetCategoryMarketVisibilityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetCategoryMarketVisibilityLogic) SetCategoryMarketVisibility(req *types.SetCategoryMarketVisibilityReq) (resp *types.CreateCategoryResp, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Verify category exists
	category, err := l.svcCtx.CategoryRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.CategoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, code.ErrCategoryNotFound
	}

	// Delete existing market visibility for this category
	if err := l.svcCtx.CategoryMarketRepo.DeleteByCategory(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.CategoryID); err != nil {
		return nil, err
	}

	// Create new market visibility entries
	if len(req.MarketIDs) > 0 {
		items := make([]*product.CategoryMarket, 0, len(req.MarketIDs))
		now := time.Now().UTC()
		for _, marketID := range req.MarketIDs {
			id, _ := l.svcCtx.IDGen.NextID(l.ctx)
			items = append(items, &product.CategoryMarket{
				Model:      application.Model{ID: id, CreatedAt: now, UpdatedAt: now},
				TenantID:   shared.TenantID(tenantID),
				CategoryID: req.CategoryID,
				MarketID:   marketID,
				IsVisible:  req.Visible,
			})
		}
		if err := l.svcCtx.CategoryMarketRepo.BatchCreate(l.ctx, l.svcCtx.DB, items); err != nil {
			return nil, err
		}
	}

	return &types.CreateCategoryResp{
		ID: req.CategoryID,
	}, nil
}
