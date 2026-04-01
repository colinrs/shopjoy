package shipping_templates

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/shipping"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/shopspring/decimal"

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
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	// Find template with details
	template, zones, mappings, err := l.svcCtx.ShippingRepo.FindByIDWithDetails(l.ctx, l.svcCtx.DB, tenantID, req.ID)
	if err != nil {
		return nil, err
	}

	// Build response
	return &types.ShippingTemplateDetailResp{
		ID:        int64(template.ID),
		Name:      template.Name,
		IsDefault: template.IsDefault,
		IsActive:  template.IsActive,
		Zones:     buildZoneDetails(zones),
		Mappings:  buildMappingDetails(mappings),
		CreatedAt: template.CreatedAt.Format(time.RFC3339),
		UpdatedAt: template.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func buildZoneDetails(zones []*shipping.ShippingZone) []*types.ShippingZoneDetail {
	result := make([]*types.ShippingZoneDetail, 0, len(zones))
	for _, z := range zones {
		result = append(result, &types.ShippingZoneDetail{
			ID:                  int64(z.ID),
			TemplateID:          z.TemplateID,
			Name:                z.Name,
			Regions:             z.Regions,
			FeeType:             string(z.FeeType),
			FirstUnit:           z.FirstUnit,
			FirstFee:            formatAmount(z.FirstFee),
			AdditionalUnit:      z.AdditionalUnit,
			AdditionalFee:       formatAmount(z.AdditionalFee),
			FreeThresholdAmount: formatAmount(z.FreeThresholdAmount),
			FreeThresholdCount:  z.FreeThresholdCount,
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

// formatAmount converts decimal.Decimal to string with 2 decimal places (in yuan)
func formatAmount(amount decimal.Decimal) string {
	return amount.StringFixed(2)
}