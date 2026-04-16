package shipping_calculator

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/domain/shipping"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/colinrs/shopjoy/pkg/contextx"
	"github.com/shopspring/decimal"

	"github.com/zeromicro/go-zero/core/logx"
)

type CalculateShippingFeeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCalculateShippingFeeLogic(ctx context.Context, svcCtx *svc.ServiceContext) CalculateShippingFeeLogic {
	return CalculateShippingFeeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CalculateShippingFeeLogic) CalculateShippingFee(req *types.CalculateShippingFeeReq) (resp *types.CalculateShippingFeeResp, err error) {
	// Get tenant ID from context
	tenantID, ok := contextx.GetTenantID(l.ctx)
	if !ok && !contextx.IsPlatformAdmin(l.ctx) {
		return nil, code.ErrUnauthorized
	}
	if contextx.IsPlatformAdmin(l.ctx) {
		tenantID = 0
	}

	// Validate input
	if len(req.Items) == 0 {
		return nil, code.ErrShippingCalcItemsRequired
	}
	if req.Address.CityCode == "" {
		return nil, code.ErrShippingCalcAddressRequired
	}

	// Validate each item
	for _, item := range req.Items {
		if item.Quantity <= 0 {
			return nil, code.ErrShippingCalcInvalidQuantity
		}
		if item.Weight <= 0 {
			return nil, code.ErrShippingCalcInvalidWeight
		}
		if item.Price == "" || parseAmount(item.Price).IsNegative() {
			return nil, code.ErrShippingCalcInvalidPrice
		}
	}

	// Convert calculator items to calculate items
	items := make([]shipping.CalculateItem, 0, len(req.Items))
	var orderAmount decimal.Decimal
	var totalWeight int
	var totalQuantity int

	for _, item := range req.Items {
		price := parseAmount(item.Price)
		items = append(items, shipping.CalculateItem{
			ProductID: item.ProductID,
			SKUID:     item.SKUID,
			Quantity:  item.Quantity,
			Weight:    item.Weight,
			Price:     price,
		})
		orderAmount = orderAmount.Add(price.Mul(decimal.NewFromInt(int64(item.Quantity))))
		totalWeight += item.Weight * item.Quantity
		totalQuantity += item.Quantity
	}

	// Find template using priority: Product > Category > Default
	template, zone := l.findTemplateForItems(tenantID, req.Address.CityCode, req.Items)
	if template == nil || zone == nil {
		return nil, code.ErrShippingTemplateNotFound
	}

	// Calculate shipping fee
	shippingFee := zone.CalculateFee(items, orderAmount, totalQuantity)

	// Build response
	return &types.CalculateShippingFeeResp{
		ShippingFee:  formatAmount(shippingFee),
		Currency:     "CNY",
		TemplateID:   int64(template.ID),
		TemplateName: template.Name,
		ZoneName:     zone.Name,
		FeeDetail: types.FeeCalculationDetail{
			FeeType:          string(zone.FeeType),
			FirstUnit:        zone.FirstUnit,
			FirstFee:         formatAmount(zone.FirstFee),
			AdditionalUnit:   zone.AdditionalUnit,
			AdditionalFee:    formatAmount(zone.AdditionalFee),
			CalculatedWeight: totalWeight,
			CalculatedUnits:  totalQuantity,
		},
	}, nil
}

// findTemplateForItems finds the appropriate template and zone using priority: Product > Category > Default
func (l *CalculateShippingFeeLogic) findTemplateForItems(tenantID int64, cityCode string, reqItems []types.CalculatorItem) (*shipping.ShippingTemplate, *shipping.ShippingZone) {
	// Priority 1: Check for product-specific template
	for _, item := range reqItems {
		mapping, err := l.svcCtx.ShippingRepo.FindMappingByTarget(l.ctx, l.svcCtx.DB, shipping.TargetTypeProduct, item.ProductID)
		if err == nil && mapping != nil {
			template, err := l.svcCtx.ShippingRepo.FindByID(l.ctx, l.svcCtx.DB, tenantID, mapping.TemplateID)
			if err == nil && template != nil && template.IsActive {
				zone := l.findZoneForCity(int64(template.ID), cityCode)
				if zone != nil {
					return template, zone
				}
			}
		}
	}

	// Priority 2: Check for category-specific template
	// Note: This would require looking up product categories from ProductRepo
	// For now, we skip this step as products would need category info

	// Priority 3: Use default template
	defaultTemplate, err := l.svcCtx.ShippingRepo.FindDefault(l.ctx, l.svcCtx.DB, tenantID)
	if err == nil && defaultTemplate != nil {
		zone := l.findZoneForCity(int64(defaultTemplate.ID), cityCode)
		if zone != nil {
			return defaultTemplate, zone
		}
	}

	return nil, nil
}

// findZoneForCity finds a zone matching the city code, or returns the first zone if no match
func (l *CalculateShippingFeeLogic) findZoneForCity(templateID int64, cityCode string) *shipping.ShippingZone {
	// Try to find zone matching the city code
	if cityCode != "" {
		zones, err := l.svcCtx.ShippingRepo.FindZoneByCityCode(l.ctx, l.svcCtx.DB, templateID, cityCode)
		if err == nil && len(zones) > 0 {
			return zones[0]
		}
	}

	// Fall back to first zone of the template
	zones, err := l.svcCtx.ShippingRepo.FindZonesByTemplateID(l.ctx, l.svcCtx.DB, templateID)
	if err == nil && len(zones) > 0 {
		return zones[0]
	}

	return nil
}

// parseAmount converts string amount to int64 cents
func parseAmount(s string) decimal.Decimal {
	if s == "" {
		return decimal.Zero
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero
	}
	return d
}

// formatAmount converts decimal.Decimal to string
func formatAmount(amount decimal.Decimal) string {
	return amount.StringFixed(2)
}
