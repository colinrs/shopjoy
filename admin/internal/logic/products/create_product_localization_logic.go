package products

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
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
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Generate ID
	id, err := l.svcCtx.IDGen.NextID(l.ctx)
	if err != nil {
		return nil, err
	}

	// Create localization entity
	localization := &product.ProductLocalization{
		ID:           id,
		TenantID:     shared.TenantID(tenantID),
		ProductID:    req.ProductID,
		LanguageCode: req.LanguageCode,
		Name:         req.Name,
		Description:  req.Description,
		AuditInfo: shared.AuditInfo{
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
	}

	if err := l.svcCtx.ProductLocalizationRepo.Create(l.ctx, l.svcCtx.DB, localization); err != nil {
		return nil, err
	}

	return &types.CreateProductResp{ID: id}, nil
}