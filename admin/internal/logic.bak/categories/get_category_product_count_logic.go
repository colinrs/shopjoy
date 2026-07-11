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

type GetCategoryProductCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCategoryProductCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetCategoryProductCountLogic {
	return GetCategoryProductCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCategoryProductCountLogic) GetCategoryProductCount(req *types.GetCategoryProductCountReq) (resp *types.GetCategoryProductCountResp, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	// Verify category exists
	category, err := l.svcCtx.CategoryRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, code.ErrCategoryNotFound
	}

	// Get product count
	count, err := l.svcCtx.CategoryRepo.GetProductCount(l.ctx, l.svcCtx.DB, req.ID)
	if err != nil {
		return nil, err
	}

	return &types.GetCategoryProductCountResp{
		Count: count,
	}, nil
}
