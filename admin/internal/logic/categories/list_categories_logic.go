package categories

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
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

	var categories []*productCategory
	if req.ParentID > 0 {
		categories, err = l.svcCtx.CategoryRepo.FindByParentID(l.ctx, l.svcCtx.DB, req.ParentID)
	} else {
		categories, err = l.svcCtx.CategoryRepo.FindAll(l.ctx, l.svcCtx.DB)
	}
	if err != nil {
		return nil, err
	}

	list := make([]*types.CategoryDetailResp, 0, len(categories))
	for _, c := range categories {
		productCount, err := l.svcCtx.CategoryRepo.GetProductCount(l.ctx, l.svcCtx.DB, c.ID)
		if err != nil {
			l.Logger.Errorf("failed to get product count for category %d: %v", c.ID, err)
			productCount = 0
		}
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
			CreatedAt:      c.Model.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      c.Model.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &types.ListCategoryResp{
		List: list,
	}, nil
}

// Alias to avoid import cycle
type productCategory = product.Category
