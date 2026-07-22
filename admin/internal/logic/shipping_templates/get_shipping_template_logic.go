package shipping_templates

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/shipping"
	"github.com/colinrs/shopjoy/admin/internal/logic/shipping_zones"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetShippingTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetShippingTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetShippingTemplateLogic {
	return GetShippingTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetShippingTemplateLogic) GetShippingTemplate(req *types.GetShippingTemplateReq) (resp *types.ShippingTemplateDetailResp, err error) {

	// Find template with details
	template, zones, mappings, err := l.svcCtx.ShippingRepo.FindByIDWithDetails(l.ctx, l.svcCtx.DB, req.ID)
	if err != nil {
		return nil, err
	}

	// ─── entity → response field map ───
	// All ShippingTemplateDetailResp fields must be populated below.
	return &types.ShippingTemplateDetailResp{
		ID:          int64(template.ID),
		TenantID:    template.TenantID,
		MarketID:    template.MarketID,
		Currency:    template.Currency,
		CarrierCode: template.CarrierCode,
		WarehouseID: template.WarehouseID,
		Name:        template.Name,
		IsDefault:   template.IsDefault,
		IsActive:    template.IsActive,
		Zones:       buildZoneDetails(zones),
		Mappings:    buildMappingDetails(mappings),
		CreatedAt:   template.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   template.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// buildZoneDetails converts entity ShippingZone → wire ShippingZoneDetail.
// Must populate ALL 22 ShippingZoneDetail fields (Task 1.7 review found
// this function silently dropped 10 new fields added in Tasks 1.2/1.6).
//
// ─── entity → response field map (must include all 22 fields) ───
//
//	entity.ID                  → resp.ID
//	entity.TenantID            → resp.TenantID            (P1-2)
//	entity.TemplateID          → resp.TemplateID
//	entity.MarketID            → resp.MarketID            (P1-5)
//	entity.Currency            → resp.Currency            (P1-2)
//	entity.Name                → resp.Name
//	entity.NameI18n            → resp.NameI18n            (P1-10, via fromStringI18n)
//	entity.Regions             → resp.Regions
//	entity.FeeType             → resp.FeeType
//	entity.FirstUnit           → resp.FirstUnit
//	entity.FirstFee            → resp.FirstFee            (decimal → string)
//	entity.AdditionalUnit      → resp.AdditionalUnit
//	entity.AdditionalFee       → resp.AdditionalFee       (decimal → string)
//	entity.FreeThresholdAmount → resp.FreeThresholdAmount (decimal → string)
//	entity.FreeThresholdCount  → resp.FreeThresholdCount
//	entity.Taxable             → resp.Taxable             (P1-6)
//	entity.TaxRate             → resp.TaxRate             (P1-6, decimal → string)
//	entity.TaxIncluded         → resp.TaxIncluded         (P1-6)
//	entity.IossApplicable      → resp.IossApplicable      (P1-6)
//	entity.RemoteSurcharge     → resp.RemoteSurcharge     (P1-7, decimal → string)
//	entity.RemoteZipPatterns   → resp.RemoteZipPatterns   (P1-7)
//	entity.FuelSurchargePct    → resp.FuelSurchargePct    (P1-8, decimal → string)
//	entity.VolumetricDivisor   → resp.VolumetricDivisor   (P1-9)
//	entity.Sort                → resp.Sort
func buildZoneDetails(zones []*shipping.ShippingZone) []*types.ShippingZoneDetail {
	result := make([]*types.ShippingZoneDetail, 0, len(zones))
	for _, z := range zones {
		result = append(result, &types.ShippingZoneDetail{
			ID:                  int64(z.ID),
			TenantID:            z.TenantID,
			TemplateID:          z.TemplateID,
			MarketID:            z.MarketID,
			Currency:            z.Currency,
			Name:                z.Name,
			NameI18n:            shipping_zones.FromStringI18n(z.NameI18n),
			Regions:             z.Regions,
			FeeType:             string(z.FeeType),
			FirstUnit:           z.FirstUnit,
			FirstFee:            utils.FormatAmount(z.FirstFee),
			AdditionalUnit:      z.AdditionalUnit,
			AdditionalFee:       utils.FormatAmount(z.AdditionalFee),
			FreeThresholdAmount: utils.FormatAmount(z.FreeThresholdAmount),
			FreeThresholdCount:  z.FreeThresholdCount,
			Taxable:             z.Taxable,
			TaxRate:             utils.FormatAmount(z.TaxRate),
			TaxIncluded:         z.TaxIncluded,
			IossApplicable:      z.IossApplicable,
			RemoteSurcharge:     utils.FormatAmount(z.RemoteSurcharge),
			RemoteZipPatterns:   z.RemoteZipPatterns,
			FuelSurchargePct:    utils.FormatAmount(z.FuelSurchargePct),
			VolumetricDivisor:   z.VolumetricDivisor,
			Sort:                z.Sort,
		})
	}
	return result
}

func buildMappingDetails(mappings []*shipping.ShippingTemplateMapping) []*types.TemplateMappingDetail {
	result := make([]*types.TemplateMappingDetail, 0, len(mappings))
	for _, m := range mappings {
		result = append(result, &types.TemplateMappingDetail{
			ID:         int64(m.ID),
			TemplateID: m.TemplateID,
			TargetType: string(m.TargetType),
			TargetID:   m.TargetID,
		})
	}
	return result
}
