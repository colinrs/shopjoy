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

	// Build response. ProductCount/CategoryCount were removed from the wire
	// (ShippingTemplateListItem no longer carries them), so we don't compute
	// them — only ZoneCount is required.
	list := make([]*types.ShippingTemplateListItem, 0, len(templates))
	for _, t := range templates {
		zoneCount, zErr := l.svcCtx.ShippingRepo.CountZonesByTemplateID(l.ctx, l.svcCtx.DB, int64(t.ID))
		if zErr != nil {
			return nil, zErr
		}

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
