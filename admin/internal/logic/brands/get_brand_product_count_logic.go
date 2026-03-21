package brands

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetBrandProductCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBrandProductCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetBrandProductCountLogic {
	return GetBrandProductCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBrandProductCountLogic) GetBrandProductCount(req *types.GetBrandProductCountReq) (resp *types.GetBrandProductCountResp, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Verify brand exists
	brand, err := l.svcCtx.BrandRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, code.ErrBrandNotFound
	}

	// Get product count
	count, err := l.svcCtx.BrandRepo.GetProductCount(l.ctx, l.svcCtx.DB, req.ID)
	if err != nil {
		return nil, err
	}

	return &types.GetBrandProductCountResp{
		Count: count,
	}, nil
}
