package shipping_templates

import (
	"context"
	"strconv"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/shipping"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"

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
	// Get tenant ID from context
	tenantID, _ := contextx.GetTenantID(l.ctx)
	if tenantID == 0 {
		return nil, code.ErrUnauthorized
	}

	// Find template with details
	template, zones, mappings, err := l.svcCtx.ShippingRepo.FindByIDWithDetails(l.ctx, l.svcCtx.DB, tenantID, req.ID)
	if err != nil {
		return nil, err
	}

	// Build response
	return &types.ShippingTemplateDetailResp{
		ID:        template.ID,
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
			ID:                  z.ID,
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
			ID:         m.ID,
			TemplateID: m.TemplateID,
			TargetType: string(m.TargetType),
			TargetID:   m.TargetID,
		})
	}
	return result
}

// Helper to convert int64 cents to string
func formatAmount(cents int64) string {
	return strconv.FormatFloat(float64(cents)/100, 'f', 2, 64)
}