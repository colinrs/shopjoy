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

	return buildCategoryTree(categories, 0), nil
}

func buildCategoryTree(categories []*product.Category, parentID int64) []*types.CategoryTreeResp {
	var result []*types.CategoryTreeResp
	for _, c := range categories {
		if c.ParentID == parentID {
			productCount, _ := getCategoryProductCount(c.ID)
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
				ProductCount:   productCount,
				CreatedAt:      c.Audit.CreatedAt.Format(time.RFC3339),
				UpdatedAt:      c.Audit.UpdatedAt.Format(time.RFC3339),
				Children:       buildCategoryTree(categories, c.ID),
			}
			result = append(result, node)
		}
	}
	return result
}

// Helper function placeholder - actual implementation would need access to repository
func getCategoryProductCount(categoryID int64) (int64, error) {
	// This is a placeholder - in actual implementation, this would call the repository
	return 0, nil
}
