package products

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProductLocalizationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProductLocalizationLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateProductLocalizationLogic {
	return UpdateProductLocalizationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProductLocalizationLogic) UpdateProductLocalization(req *types.UpdateProductLocalizationReq) (resp *types.ProductLocalizationResp, err error) {
	// Get tenant ID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Find existing localization
	localization, err := l.svcCtx.ProductLocalizationRepo.FindByID(l.ctx, l.svcCtx.DB, shared.TenantID(tenantID), req.ID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		localization.Name = req.Name
	}
	if req.Description != "" {
		localization.Description = req.Description
	}
	localization.AuditInfo.UpdatedAt = time.Now().UTC()

	// Save
	if err := l.svcCtx.ProductLocalizationRepo.Update(l.ctx, l.svcCtx.DB, localization); err != nil {
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