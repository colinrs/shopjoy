package seo

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	appStorefront "github.com/colinrs/shopjoy/admin/internal/application/storefront"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePageSEOLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePageSEOLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdatePageSEOLogic {
	return UpdatePageSEOLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePageSEOLogic) UpdatePageSEO(req *types.UpdatePageSEORequest) error {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return err
	}

	config := appStorefront.SEOConfigDTO{
		Title:       req.Title,
		Description: req.Description,
		Keywords:    req.Keywords,
	}

	return l.svcCtx.SEOService.UpdatePageSEO(l.ctx, shared.TenantID(tenantID), req.PageType, nil, config)
}