package shipping_zones

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/domain/shipping"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateShippingZoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateShippingZoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateShippingZoneLogic {
	return UpdateShippingZoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateShippingZoneLogic) UpdateShippingZone(req *types.UpdateShippingZoneReq) (resp *types.ShippingZoneDetail, err error) {
	// Get tenant ID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)
	if tenantID == 0 {
		return nil, code.ErrUnauthorized
	}

	// Find existing zone
	zone, err := l.svcCtx.ShippingRepo.FindZoneByID(l.ctx, l.svcCtx.DB, req.ID)
	if err != nil {
		return nil, err
	}

	// Verify zone belongs to tenant
	if zone.TenantID != tenantID {
		return nil, code.ErrShippingZoneNotFound
	}

	// Update fields if provided
	if req.Name != "" {
		zone.Name = req.Name
	}
	if req.Regions != nil {
		zone.Regions = req.Regions
	}
	if req.FeeType != "" {
		zone.FeeType = shipping.FeeType(req.FeeType)
	}
	if req.FirstUnit != 0 {
		zone.FirstUnit = req.FirstUnit
	}
	if req.FirstFee != "" {
		zone.FirstFee = parseAmount(req.FirstFee)
	}
	if req.AdditionalUnit != 0 {
		zone.AdditionalUnit = req.AdditionalUnit
	}
	if req.AdditionalFee != "" {
		zone.AdditionalFee = parseAmount(req.AdditionalFee)
	}
	if req.FreeThresholdAmount != "" {
		zone.FreeThresholdAmount = parseAmount(req.FreeThresholdAmount)
	}
	if req.FreeThresholdCount != 0 {
		zone.FreeThresholdCount = req.FreeThresholdCount
	}
	if req.Sort != 0 {
		zone.Sort = req.Sort
	}

	// Validate zone
	if err := zone.Validate(); err != nil {
		return nil, err
	}

	// Save changes
	if err := l.svcCtx.ShippingRepo.UpdateZone(l.ctx, l.svcCtx.DB, zone); err != nil {
		return nil, err
	}

	return &types.ShippingZoneDetail{
		ID:                  int64(zone.ID),
		TemplateID:          zone.TemplateID,
		Name:                zone.Name,
		Regions:             zone.Regions,
		FeeType:             string(zone.FeeType),
		FirstUnit:           zone.FirstUnit,
		FirstFee:            formatAmount(zone.FirstFee),
		AdditionalUnit:      zone.AdditionalUnit,
		AdditionalFee:       formatAmount(zone.AdditionalFee),
		FreeThresholdAmount: formatAmount(zone.FreeThresholdAmount),
		FreeThresholdCount:  zone.FreeThresholdCount,
		Sort:                zone.Sort,
	}, nil
}