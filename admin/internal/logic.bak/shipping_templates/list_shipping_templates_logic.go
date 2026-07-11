package shipping_templates

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListShippingTemplatesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListShippingTemplatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListShippingTemplatesLogic {
	return ListShippingTemplatesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListShippingTemplatesLogic) ListShippingTemplates(req *types.ListShippingTemplatesReq) (resp *types.ListShippingTemplatesResp, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	// Find templates with stats (single query with subqueries)
	results, total, err := l.svcCtx.ShippingRepo.FindListWithStats(
		l.ctx, l.svcCtx.DB, tenantID,
		req.Name, req.IsActive,
		req.Page, req.PageSize,
	)
	if err != nil {
		return nil, err
	}

	// Build response
	list := make([]*types.ShippingTemplateListItem, 0, len(results))
	for _, t := range results {
		list = append(list, &types.ShippingTemplateListItem{
			ID:            int64(t.ID),
			Name:          t.Name,
			IsDefault:     t.IsDefault,
			IsActive:      t.IsActive,
			ZoneCount:     int(t.ZoneCount),
			ProductCount:  int(t.ProductCount),
			CategoryCount: int(t.CategoryCount),
			CreatedAt:     t.CreatedAt.Format(time.RFC3339),
		})
	}

	return &types.ListShippingTemplatesResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
