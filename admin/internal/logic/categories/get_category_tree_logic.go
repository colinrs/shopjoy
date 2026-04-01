package categories

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetCategoryTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCategoryTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetCategoryTreeLogic {
	return GetCategoryTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCategoryTreeLogic) GetCategoryTree(req *types.CategoryTreeReq) (resp []*types.CategoryTreeResp, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)

	categories, err := l.svcCtx.CategoryRepo.FindTree(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID))
	if err != nil {
		return nil, err
	}

	// Build a map of product counts for all categories
	productCountMap := make(map[int64]int64)
	for _, c := range categories {
		count, err := l.svcCtx.CategoryRepo.GetProductCount(l.ctx, l.svcCtx.DB, c.ID)
		if err == nil {
			productCountMap[c.ID] = count
		}
	}

	return buildCategoryTree(categories, 0, productCountMap), nil
}

func buildCategoryTree(categories []*product.Category, parentID int64, productCountMap map[int64]int64) []*types.CategoryTreeResp {
	var result []*types.CategoryTreeResp
	for _, c := range categories {
		if c.ParentID == parentID {
			node := &types.CategoryTreeResp{
				ID:             c.ID,
				ParentID:       c.ParentID,
				Name:           c.Name,
				Code:           c.Code,
				Level:          c.Level,
				Sort:           c.Sort,
				Icon:           c.Icon,
				Image:          c.Image,
				SeoTitle:       c.SeoTitle,
				SeoDescription: c.SeoDescription,
				Status:         int8(c.Status),
				ProductCount:   productCountMap[c.ID],
				CreatedAt:      c.Audit.CreatedAt.Format(time.RFC3339),
				UpdatedAt:      c.Audit.UpdatedAt.Format(time.RFC3339),
				Children:       buildCategoryTree(categories, c.ID, productCountMap),
			}
			result = append(result, node)
		}
	}
	return result
}
