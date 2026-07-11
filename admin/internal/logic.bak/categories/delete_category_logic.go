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

type DeleteCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteCategoryLogic {
	return DeleteCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCategoryLogic) DeleteCategory(req *types.GetCategoryReq) (resp *types.CreateCategoryResp, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	// Find category
	category, err := l.svcCtx.CategoryRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, code.ErrCategoryNotFound
	}

	// Check if category has children
	children, err := l.svcCtx.CategoryRepo.FindByParentID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}
	if len(children) > 0 {
		return nil, code.ErrCategoryHasChildren
	}

	// Check if category has products
	productCount, err := l.svcCtx.CategoryRepo.GetProductCount(l.ctx, l.svcCtx.DB, req.ID)
	if err != nil {
		return nil, err
	}
	if productCount > 0 {
		return nil, code.ErrCategoryHasProducts
	}

	// Delete category
	if err := l.svcCtx.CategoryRepo.Delete(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID); err != nil {
		return nil, err
	}

	return &types.CreateCategoryResp{
		ID: req.ID,
	}, nil
}
