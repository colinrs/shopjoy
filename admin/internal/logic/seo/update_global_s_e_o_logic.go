package seo

import (
	"context"

	appStorefront "github.com/colinrs/shopjoy/admin/internal/application/storefront"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateGlobalSEOLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateGlobalSEOLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateGlobalSEOLogic {
	return UpdateGlobalSEOLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateGlobalSEOLogic) UpdateGlobalSEO(req *types.UpdateSEOConfigRequest) error {

	config := appStorefront.SEOConfigDTO{
		Title:       req.Title,
		Description: req.Description,
		Keywords:    req.Keywords,
	}

	return l.svcCtx.SEOService.UpdateGlobalSEO(l.ctx, config)
}
