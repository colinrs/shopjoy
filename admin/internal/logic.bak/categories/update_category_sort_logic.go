package categories

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCategorySortLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCategorySortLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateCategorySortLogic {
	return UpdateCategorySortLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCategorySortLogic) UpdateCategorySort(req *types.UpdateCategorySortReq) (resp *types.CreateCategoryResp, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	// Convert to domain type
	sorts := make([]product.CategorySort, 0, len(req.Sorts))
	for _, s := range req.Sorts {
		sorts = append(sorts, product.CategorySort{
			ID:   s.ID,
			Sort: s.Sort,
		})
	}

	if err := l.svcCtx.CategoryRepo.UpdateSort(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), sorts); err != nil {
		return nil, err
	}

	return &types.CreateCategoryResp{
		ID: 0,
	}, nil
}
