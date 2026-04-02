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

type ListCategoriesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListCategoriesLogic {
	return ListCategoriesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListCategoriesLogic) ListCategories(req *types.ListCategoryReq) (resp *types.ListCategoryResp, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	var categories []*productCategory
	if req.ParentID > 0 {
		categories, err = l.svcCtx.CategoryRepo.FindByParentID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ParentID)
	} else {
		categories, err = l.svcCtx.CategoryRepo.FindAll(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID))
	}
	if err != nil {
		return nil, err
	}

	list := make([]*types.CategoryDetailResp, 0, len(categories))
	for _, c := range categories {
		productCount, _ := l.svcCtx.CategoryRepo.GetProductCount(l.ctx, l.svcCtx.DB, c.ID)
		list = append(list, &types.CategoryDetailResp{
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
			Status:         int8(c.Status), // #nosec G115 // status values are small (tinyint range)
			ProductCount:   productCount,
			CreatedAt:      c.Audit.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      c.Audit.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &types.ListCategoryResp{
		List: list,
	}, nil
}

// Alias to avoid import cycle
type productCategory = product.Category
