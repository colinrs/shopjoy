package products

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

type GetProductLocalizationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductLocalizationLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetProductLocalizationLogic {
	return GetProductLocalizationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductLocalizationLogic) GetProductLocalization(req *types.GetProductLocalizationReq) (resp *types.ProductLocalizationResp, err error) {
	// Get tenant ID from context
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok && !contextx.IsPlatformAdmin(l.ctx) {
		return nil, code.ErrUnauthorized
	}
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Find localization
	localization, err := l.svcCtx.ProductLocalizationRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}

	return &types.ProductLocalizationResp{
		ID:           localization.ID,
		ProductID:    localization.ProductID,
		LanguageCode: localization.LanguageCode,
		Name:         localization.Name,
		Description:  localization.Description,
		CreatedAt:    localization.AuditInfo.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    localization.AuditInfo.UpdatedAt.Format(time.RFC3339),
	}, nil
}
