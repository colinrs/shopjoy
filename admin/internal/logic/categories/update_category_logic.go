package categories

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateCategoryLogic {
	return UpdateCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCategoryLogic) UpdateCategory(req *types.UpdateCategoryReq) (resp *types.CategoryDetailResp, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Find existing category
	category, err := l.svcCtx.CategoryRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, code.ErrCategoryNotFound
	}

	// Check for duplicate code if changed
	if req.Code != "" && req.Code != category.Code {
		existing, err := l.svcCtx.CategoryRepo.FindByCode(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.Code)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != req.ID {
			return nil, code.ErrCategoryDuplicate
		}
	}

	// Update fields
	category.Name = req.Name
	category.Code = req.Code
	category.Sort = req.Sort
	category.Icon = req.Icon
	category.Image = req.Image
	category.SeoTitle = req.SeoTitle
	category.SeoDescription = req.SeoDescription
	category.Audit.UpdatedAt = time.Now()

	if err := l.svcCtx.CategoryRepo.Update(l.ctx, l.svcCtx.DB, category); err != nil {
		return nil, err
	}

	// Get product count
	productCount, _ := l.svcCtx.CategoryRepo.GetProductCount(l.ctx, l.svcCtx.DB, category.ID)

	return &types.CategoryDetailResp{
		ID:             category.ID,
		ParentID:       category.ParentID,
		Name:           category.Name,
		Code:           category.Code,
		Level:          category.Level,
		Sort:           category.Sort,
		Icon:           category.Icon,
		Image:          category.Image,
		SeoTitle:       category.SeoTitle,
		SeoDescription: category.SeoDescription,
		Status:         int8(category.Status),
		ProductCount:   productCount,
		CreatedAt:      category.Audit.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      category.Audit.UpdatedAt.Format(time.RFC3339),
	}, nil
}
