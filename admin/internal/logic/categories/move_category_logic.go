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

	// Find new parent
	if req.NewParentID > 0 {
		newParent, err := l.svcCtx.CategoryRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.NewParentID)
		if err != nil {
			return nil, err
		}
		if newParent == nil {
			return nil, code.ErrCategoryNotFound
		}
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
		CreatedAt:      category.Audit.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      category.Audit.UpdatedAt.Format(time.RFC3339),
	}, nil
}
