package pages

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/colinrs/shopjoy/pkg/domain/shared"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListPagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListPagesLogic {
	return ListPagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPagesLogic) ListPages(req *types.ListPagesRequest) (resp *types.ListPagesResponse, err error) {
	tenantID, _ := contextx.GetTenantID(l.ctx)

	result, err := l.svcCtx.PageService.ListPages(l.ctx, shared.TenantID(tenantID), req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	items := make([]*types.PageListItem, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, &types.PageListItem{
			ID:          p.ID,
			PageType:    p.PageType,
			Name:        p.Name,
			Slug:        p.Slug,
			IsPublished: p.IsPublished,
			Version:     p.Version,
		})
	}

	return &types.ListPagesResponse{
		Pages:    items,
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.PageSize,
	}, nil
}