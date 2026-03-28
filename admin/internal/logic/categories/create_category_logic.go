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

type CreateCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateCategoryLogic {
	return CreateCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCategoryLogic) CreateCategory(req *types.CreateCategoryReq) (resp *types.CreateCategoryResp, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Check for duplicate code if provided
	if req.Code != "" {
		existing, err := l.svcCtx.CategoryRepo.FindByCode(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.Code)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, code.ErrCategoryDuplicate
		}
	}

	// Calculate level based on parent
	level := 1
	if req.ParentID > 0 {
		parent, err := l.svcCtx.CategoryRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ParentID)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			return nil, code.ErrCategoryNotFound
		}
		level = parent.Level + 1
		// Validate max level (max 3 levels)
		if level > 3 {
			return nil, code.ErrCategoryMaxLevelExceeded
		}
	}

	// Generate ID
	id, err := l.svcCtx.IDGen.NextID(l.ctx)
	if err != nil {
		return nil, err
	}

	category := &product.Category{
		Model:          application.Model{ID: id, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
		TenantID:       shared.TenantID(tenantID),
		ParentID:       req.ParentID,
		Name:           req.Name,
		Code:           req.Code,
		Level:          level,
		Sort:           req.Sort,
		Icon:           req.Icon,
		Image:          req.Image,
		SeoTitle:       req.SeoTitle,
		SeoDescription: req.SeoDescription,
		Status:         product.CategoryStatusEnabled,
	}

	if err := l.svcCtx.CategoryRepo.Create(l.ctx, l.svcCtx.DB, category); err != nil {
		return nil, err
	}

	return &types.CreateCategoryResp{
		ID: category.Model.ID,
	}, nil
}
