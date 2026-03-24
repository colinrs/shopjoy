package seo

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListPageSEOLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPageSEOLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListPageSEOLogic {
	return ListPageSEOLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPageSEOLogic) ListPageSEO() (resp *types.ListPageSEOConfigsResponse, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)

	configs, err := l.svcCtx.SEOService.ListPageSEO(l.ctx, shared.TenantID(tenantID))
	if err != nil {
		return nil, err
	}

	items := make([]*types.PageSEOConfigResponse, 0, len(configs))
	for _, c := range configs {
		items = append(items, &types.PageSEOConfigResponse{
			PageType: c.PageType,
			PageID:   c.PageID,
			SEOConfigDTO: types.SEOConfigDTO{
				Title:       c.Config.Title,
				Description: c.Config.Description,
				Keywords:    c.Config.Keywords,
			},
		})
	}

	return &types.ListPageSEOConfigsResponse{
		Configs: items,
	}, nil
}