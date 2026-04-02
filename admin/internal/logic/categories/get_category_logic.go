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

type GetCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetCategoryLogic {
	return GetCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCategoryLogic) GetCategory(req *types.GetCategoryReq) (resp *types.CategoryDetailResp, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	category, err := l.svcCtx.CategoryRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, code.ErrCategoryNotFound
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
		Status:         int8(category.Status), // #nosec G115 // status values are small (tinyint range)
		ProductCount:   productCount,
		CreatedAt:      category.Audit.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      category.Audit.UpdatedAt.Format(time.RFC3339),
	}, nil
}
