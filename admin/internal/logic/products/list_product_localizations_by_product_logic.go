package products

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
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
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Platform admin can access all tenant data
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Find localizations by product ID
	localizations, err := l.svcCtx.ProductLocalizationRepo.FindByProductID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ProductID)
	if err != nil {
		return nil, err
	}

	// Convert to response
	list := make([]*types.ProductLocalizationResp, len(localizations))
	for i, loc := range localizations {
		list[i] = &types.ProductLocalizationResp{
			ID:           loc.ID,
			ProductID:    loc.ProductID,
			LanguageCode: loc.LanguageCode,
			Name:         loc.Name,
			Description:  loc.Description,
			CreatedAt:    loc.AuditInfo.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    loc.AuditInfo.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &types.ListProductLocalizationsResp{
		List:  list,
		Total: int64(len(list)),
	}, nil
}