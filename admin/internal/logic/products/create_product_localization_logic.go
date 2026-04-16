package products

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

type CreateProductLocalizationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateProductLocalizationLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateProductLocalizationLogic {
	return CreateProductLocalizationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateProductLocalizationLogic) CreateProductLocalization(req *types.CreateProductLocalizationReq) (resp *types.CreateProductResp, err error) {
	// Get tenant ID from context
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok && !contextx.IsPlatformAdmin(l.ctx) {
		return nil, code.ErrUnauthorized
	}
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Generate ID
	id, err := l.svcCtx.IDGen.NextID(l.ctx)
	if err != nil {
		return nil, err
	}

	// Create localization entity
	localization := &product.ProductLocalization{
		Model:        application.Model{ID: id, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
		TenantID:     shared.TenantID(tenantID),
		ProductID:    req.ProductID,
		LanguageCode: req.LanguageCode,
		Name:         req.Name,
		Description:  req.Description,
	}

	if err := l.svcCtx.ProductLocalizationRepo.Create(l.ctx, l.svcCtx.DB, localization); err != nil {
		return nil, err
	}

	return &types.CreateProductResp{ID: id}, nil
}
