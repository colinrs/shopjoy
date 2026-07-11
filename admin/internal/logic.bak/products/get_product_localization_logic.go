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
	if !ok {
		return nil, code.ErrUnauthorized
	}

	// Find localization
	localization, err := l.svcCtx.ProductLocalizationRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}

	return toProductLocalizationResp(localization), nil
}
