package shipping_calculator

import (
	"context"
	"strconv"

	"github.com/colinrs/shopjoy/admin/internal/domain/shipping"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
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

// parsedItem holds the converted (int64, decimal) form of a wire CalculatorItem
// so we don't repeat strconv.ParseInt at every call site.
type parsedItem struct {
	productID int64
	skuID     int64
	quantity  int
	weight    int
	price     decimal.Decimal
}

func (l *CalculateShippingFeeLogic) CalculateShippingFee(req *types.CalculateShippingFeeReq) (resp *types.CalculateShippingFeeResp, err error) {
	// Validate input
	if len(req.Items) == 0 {
		return nil, code.ErrShippingCalcItemsRequired
	}
	if req.Address.CityCode == "" {
		return nil, code.ErrShippingCalcAddressRequired
	}

	// Validate each item + convert wire string IDs to int64 in one pass.
	parsed := make([]parsedItem, 0, len(req.Items))
	items := make([]shipping.CalculateItem, 0, len(req.Items))
	var orderAmount decimal.Decimal
	var totalWeight int
	var totalQuantity int

	for _, item := range req.Items {
		if item.Quantity <= 0 {
			return nil, code.ErrShippingCalcInvalidQuantity
		}
		if item.Weight <= 0 {
			return nil, code.ErrShippingCalcInvalidWeight
		}
		price := parseAmount(item.Price)
		if item.Price == "" || price.IsNegative() {
			return nil, code.ErrShippingCalcInvalidPrice
		}

		// wire.ProductID is a JSON string of int64 — must parse before use.
		productID, perr := strconv.ParseInt(item.ProductID, 10, 64)
		if perr != nil || productID <= 0 {
			return nil, code.ErrSharedInvalidParam
		}
		// wire.SKUID is optional; default to 0 if empty.
		var skuID int64
		if item.SKUID != "" {
			skuID, perr = strconv.ParseInt(item.SKUID, 10, 64)
			if perr != nil || skuID < 0 {
				return nil, code.ErrSharedInvalidParam
			}
		}

		pi := parsedItem{
			productID: productID,
			skuID:     skuID,
			quantity:  item.Quantity,
			weight:    item.Weight,
			price:     price,
		}
		parsed = append(parsed, pi)
		items = append(items, shipping.CalculateItem{
			ProductID: pi.productID,
			SKUID:     pi.skuID,
			Quantity:  pi.quantity,
			Weight:    pi.weight,
			Price:     pi.price,
		})
		orderAmount = orderAmount.Add(pi.price.Mul(decimal.NewFromInt(int64(pi.quantity))))
		totalWeight += pi.weight * pi.quantity
		totalQuantity += pi.quantity
	}

	// Find template using priority: Product > Default (market-aware)
	template, zone := l.findTemplateForItems(req.MarketID, req.Address.CityCode, parsed)
	if template == nil || zone == nil {
		return nil, code.ErrShippingTemplateNotFound
	}

	// Calculate shipping fee
	shippingFee := zone.CalculateFee(items, orderAmount, totalQuantity)

	// Build response
	return &types.CalculateShippingFeeResp{
		ShippingFee:  formatAmount(shippingFee),
		Currency:     template.Currency,
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

// findTemplateForItems finds the appropriate template and zone using priority: Product > Default.
// marketID scopes the default-template lookup via FindDefaultByMarket (with fallback to marketID=0).
func (l *CalculateShippingFeeLogic) findTemplateForItems(marketID int64, cityCode string, items []parsedItem) (*shipping.ShippingTemplate, *shipping.ShippingZone) {
	// Priority 1: Check for product-specific template
	for _, item := range items {
		mapping, err := l.svcCtx.ShippingRepo.FindMappingByTarget(l.ctx, l.svcCtx.DB, shipping.TargetTypeProduct, item.productID)
		if err == nil && mapping != nil {
			template, err := l.svcCtx.ShippingRepo.FindByID(l.ctx, l.svcCtx.DB, mapping.TemplateID)
			if err == nil && template != nil && template.IsActive {
				zone := l.findZoneForCity(int64(template.ID), cityCode)
				if zone != nil {
					return template, zone
				}
			}
		}
	}

	// Priority 2: Use market-scoped default template (falls back to marketID=0).
	defaultTemplate, err := l.svcCtx.ShippingRepo.FindDefaultByMarket(l.ctx, l.svcCtx.DB, marketID)
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
		zones, err := l.svcCtx.ShippingRepo.FindZoneByCityCode(l.ctx, l.svcCtx.DB, cityCode)
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

// parseAmount converts string amount to decimal.Decimal; empty/invalid → 0.
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

// formatAmount converts decimal.Decimal to 2-decimal string.
func formatAmount(amount decimal.Decimal) string {
	return amount.StringFixed(2)
}
