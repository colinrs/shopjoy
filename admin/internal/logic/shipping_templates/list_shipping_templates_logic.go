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

	// C3 fix: tenant scope is required to prevent cross-tenant data leak.
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Market-aware + tenant-aware list: marketID=0 means "all markets within
	// this tenant"; the query never crosses tenant boundaries.
	templates, total, err := l.svcCtx.ShippingRepo.FindListByMarket(
		l.ctx, l.svcCtx.DB,
		tenantID,
		req.MarketID,
		req.Name,
		req.IsActive,
		req.Page, req.PageSize,
	)
	if err != nil {
		return nil, err
	}

	// Batch the per-template zone counts in one GROUP BY query.
	// Important/N+1 fix: previously this loop called
	// CountZonesByTemplateID once per template (page_size extra queries).
	// Now we collect all template IDs first and let the repo issue a single
	// SELECT … GROUP BY template_id query — constant query count regardless
	// of page size. An empty page short-circuits inside the repo (no SQL).
	var ids []int64
	if len(templates) > 0 {
		ids = make([]int64, len(templates))
		for i, t := range templates {
			ids[i] = int64(t.ID)
		}
	}
	zoneCounts, zErr := l.svcCtx.ShippingRepo.CountZonesByTemplateIDs(l.ctx, l.svcCtx.DB, ids)
	if zErr != nil {
		return nil, zErr
	}

	// Build response. ProductCount/CategoryCount were removed from the wire
	// (ShippingTemplateListItem no longer carries them), so we don't compute
	// them — only ZoneCount is required.
	list := make([]*types.ShippingTemplateListItem, 0, len(templates))
	for _, t := range templates {
		// Templates with zero zones are absent from the map; map lookup
		// returns the zero value (0) which is what we want for ZoneCount.
		zoneCount := zoneCounts[int64(t.ID)]

		// ─── entity → response field map ───
		//   entity.ID          → resp.ID
		//   entity.TenantID    → resp.TenantID
		//   entity.MarketID    → resp.MarketID
		//   entity.Currency    → resp.Currency
		//   entity.CarrierCode → resp.CarrierCode
		//   entity.WarehouseID → resp.WarehouseID
		//   entity.Name        → resp.Name
		//   entity.IsDefault   → resp.IsDefault
		//   entity.IsActive    → resp.IsActive
		//   repo count         → resp.ZoneCount
		//   entity.CreatedAt   → resp.CreatedAt
		list = append(list, &types.ShippingTemplateListItem{
			ID:          int64(t.ID),
			TenantID:    t.TenantID,
			MarketID:    t.MarketID,
			Currency:    t.Currency,
			CarrierCode: t.CarrierCode,
			WarehouseID: t.WarehouseID,
			Name:        t.Name,
			IsDefault:   t.IsDefault,
			IsActive:    t.IsActive,
			ZoneCount:   int(zoneCount),
			CreatedAt:   t.CreatedAt.Format(time.RFC3339),
		})
	}

	return &types.ListShippingTemplatesResp{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
