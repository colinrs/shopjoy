package categories

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

type MoveCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMoveCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) MoveCategoryLogic {
	return MoveCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MoveCategoryLogic) MoveCategory(req *types.MoveCategoryReq) (resp *types.CategoryDetailResp, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Find category
	category, err := l.svcCtx.CategoryRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, code.ErrCategoryNotFound
	}

	// Calculate new level
	newLevel := 1
	if req.NewParentID > 0 {
		// Cannot move to self
		if req.NewParentID == req.ID {
			return nil, code.ErrCategoryInvalid
		}
		newParent, err := l.svcCtx.CategoryRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.NewParentID)
		if err != nil {
			return nil, err
		}
		if newParent == nil {
			return nil, code.ErrCategoryNotFound
		}
		newLevel = newParent.Level + 1
	}

	// Validate new level doesn't exceed max (3 levels)
	// Also need to check if category has children - they would go even deeper
	children, err := l.svcCtx.CategoryRepo.FindByParentID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}
	maxChildDepth := 0
	if len(children) > 0 {
		maxChildDepth = l.getMaxChildDepth(children, 1)
	}
	// New level + max child depth should not exceed 3
	if newLevel+maxChildDepth > 3 {
		return nil, code.ErrCategoryMaxLevelExceeded
	}

	// Move category
	if err := l.svcCtx.CategoryRepo.Move(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID, req.NewParentID); err != nil {
		return nil, err
	}

	// Get updated category
	category, err = l.svcCtx.CategoryRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
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
		CreatedAt:      time.Unix(category.Audit.CreatedAt, 0).Format(time.RFC3339),
		UpdatedAt:      time.Unix(category.Audit.UpdatedAt, 0).Format(time.RFC3339),
	}, nil
}

// getMaxChildDepth recursively calculates the maximum depth of child categories
func (l *MoveCategoryLogic) getMaxChildDepth(categories []*product.Category, currentDepth int) int {
	maxDepth := currentDepth
	for _, cat := range categories {
		children, err := l.svcCtx.CategoryRepo.FindByParentID(l.ctx, l.svcCtx.DB, cat.TenantID, cat.ID)
		if err == nil && len(children) > 0 {
			childDepth := l.getMaxChildDepth(children, currentDepth+1)
			if childDepth > maxDepth {
				maxDepth = childDepth
			}
		}
	}
	return maxDepth
}
