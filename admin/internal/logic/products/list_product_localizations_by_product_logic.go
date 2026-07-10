package products

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListProductLocalizationsByProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListProductLocalizationsByProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListProductLocalizationsByProductLogic {
	return ListProductLocalizationsByProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListProductLocalizationsByProductLogic) ListProductLocalizationsByProduct(req *types.ListProductLocalizationsReq) (resp *types.ListProductLocalizationsResp, err error) {
	// Get tenant ID from context
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok {
		return nil, code.ErrUnauthorized
	}

	// Find localizations by product ID
	localizations, err := l.svcCtx.ProductLocalizationRepo.FindByProductID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ProductID)
	if err != nil {
		return nil, err
	}

	// Convert to response
	list := make([]*types.ProductLocalizationResp, len(localizations))
	for i, loc := range localizations {
		list[i] = toProductLocalizationResp(loc)
	}

	return &types.ListProductLocalizationsResp{
		List:  list,
		Total: int64(len(list)),
	}, nil
}
