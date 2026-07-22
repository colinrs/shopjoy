package shipping_calculator

import (
	"context"
	"strconv"

	shippingapp "github.com/colinrs/shopjoy/admin/internal/application/shipping"
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

// parsedItem holds the converted (int64, decimal) form of a wire CalculatorItem
// so we don't repeat strconv.ParseInt at every call site.
type parsedItem struct {
	productID int64
	skuID     int64
	quantity  int
	weight    int
	length    int
	width     int
	height    int
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
	if req.MarketID == 0 {
		return nil, code.ErrShippingCalcMarketRequired
	}

	// Validate each item + convert wire string IDs to int64 in one pass.
	parsed := make([]parsedItem, 0, len(req.Items))
	items := make([]shipping.CalculateItem, 0, len(req.Items))
	var orderAmount decimal.Decimal
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
			length:    item.Length,
			width:     item.Width,
			height:    item.Height,
			price:     price,
		}
		parsed = append(parsed, pi)
		items = append(items, shipping.CalculateItem{
			ProductID: pi.productID,
			SKUID:     pi.skuID,
			Quantity:  pi.quantity,
			Weight:    pi.weight,
			Length:    pi.length,
			Width:     pi.width,
			Height:    pi.height,
			Price:     pi.price,
		})
		orderAmount = orderAmount.Add(pi.price.Mul(decimal.NewFromInt(int64(pi.quantity))))
		totalQuantity += pi.quantity
	}

	// Find template using priority: Product > Default (market-aware, tenant-scoped)
	tenantID, _ := contextx.GetTenantID(l.ctx)
	template, zone := l.findTemplateForItems(tenantID, req.MarketID, req.Address, parsed)

	// ZoneName is locale-aware when Accept-Language is injected into ctx by the
	// handler (see handler header-injection block). Empty/missing header falls
	// back to zone.Name inside ResolveZoneName.
	acceptLanguage := getAcceptLanguage(l.ctx)
	if template == nil || zone == nil {
		return nil, code.ErrShippingTemplateNotFound
	}

	// Calculate shipping fee
	shippingFee := zone.CalculateFee(items, orderAmount, totalQuantity)

	// Surcharges: remote-area + fuel are computed off the base shipping fee and
	// added on top before tax is derived, so tax applies to the surcharged fee.
	surcharge := shipping.CalculateSurcharges(zone, shipping.SurchargeInput{
		BaseFee:    shippingFee,
		PostalCode: req.Address.PostalCode,
	})
	shippingFee = shippingFee.Add(surcharge.Total)

	// Tax + Total: when the zone is taxable, derive tax from the shipping fee.
	// Non-taxable zones (or a zero rate) yield tax=0 and total==shippingFee.
	var tax, total decimal.Decimal
	if zone.Taxable {
		tax, total = shipping.CalculateTax(shippingFee, zone.TaxRate, zone.TaxIncluded)
	} else {
		total = shippingFee
	}

	// Carrier code from the template (fallback "standard").
	carrierCode := resolveCarrierCode(template.CarrierCode)

	// EstimatedDays from the carrier registry, keyed by the template's carrier
	// code. Unknown carriers (registry miss, future-proofing) fall back to 5
	// days so callers always receive a usable SLA instead of 0.
	estimatedDays := resolveEstimatedDays(l.svcCtx.CarrierRegistry, carrierCode, req.Address.CountryCode)

	// Fee-detail weight breakdown: CalculatedWeight is the chargeable weight
	// (max of real vs. volumetric), VolumetricWeight is the accumulated
	// volumetric weight for debug display. Uses the zone's divisor.
	divisor := zone.VolumetricDivisor
	if divisor <= 0 {
		divisor = shipping.DefaultVolumetricDivisor
	}
	chargeableWeight, volumetricWeight := sumWeights(items, divisor)

	// Build response
	return &types.CalculateShippingFeeResp{
		ShippingFee:      formatAmount(shippingFee),
		Tax:              formatAmount(tax),
		Total:            formatAmount(total),
		Currency:         template.Currency,
		PriceIncludesTax: zone.TaxIncluded,
		TemplateID:       int64(template.ID),
		TemplateName:     template.Name,
		ZoneName:         shippingapp.ResolveZoneName(zone, acceptLanguage),
		CarrierCode:      carrierCode,
		EstimatedDays:    estimatedDays,
		FeeDetail: types.FeeCalculationDetail{
			FeeType:          string(zone.FeeType),
			FirstUnit:        zone.FirstUnit,
			FirstFee:         formatAmount(zone.FirstFee),
			AdditionalUnit:   zone.AdditionalUnit,
			AdditionalFee:    formatAmount(zone.AdditionalFee),
			CalculatedWeight: chargeableWeight,
			VolumetricWeight: volumetricWeight,
			CalculatedUnits:  totalQuantity,
			AppliedSurcharge: formatAmount(surcharge.Total),
			AppliedTax:       formatAmount(tax),
		},
	}, nil
}

// resolveCarrierCode returns the template carrier code, falling back to
// "standard" when the template leaves it empty.
func resolveCarrierCode(carrierCode string) string {
	if carrierCode == "" {
		return "standard"
	}
	return carrierCode
}

// resolveEstimatedDays looks up the carrier by code in the registry and asks
// it to estimate transit days for the destination country. When the registry
// is nil, when the carrier code is unknown, or when the carrier lookup errors
// silently, it falls back to a 5-day default so downstream code never sees 0.
//
// 5 is chosen as a neutral SLA between StandardCarrier's "regional" tier
// (CN=3, JP/KR/SEA=5) and its "long-haul" tier (US/EU=7, default=10).
func resolveEstimatedDays(registry *shipping.CarrierRegistry, carrierCode, countryCode string) int {
	const fallbackDays = 5
	if registry == nil {
		return fallbackDays
	}
	carrier, ok := registry.Get(carrierCode)
	if !ok || carrier == nil {
		return fallbackDays
	}
	return carrier.EstimateDays(countryCode)
}

// sumWeights accumulates the chargeable weight (max of real vs. volumetric,
// per item, multiplied by quantity) and the volumetric weight across all items
// using the given volumetric divisor. Both are returned in grams.
func sumWeights(items []shipping.CalculateItem, divisor int) (chargeable, volumetric int) {
	for _, item := range items {
		chargeable += item.ChargeableWeight(divisor) * item.Quantity
		volumetric += item.VolumetricWeight(divisor) * item.Quantity
	}
	return chargeable, volumetric
}

// findTemplateForItems finds the appropriate template and zone using priority: Product > Default.
// marketID scopes the default-template lookup via FindDefaultByMarket (with fallback to marketID=0).
// tenantID is REQUIRED (C3 fix) — without it, the storefront could see another tenant's defaults.
func (l *CalculateShippingFeeLogic) findTemplateForItems(tenantID, marketID int64, address types.CalculatorAddress, items []parsedItem) (*shipping.ShippingTemplate, *shipping.ShippingZone) {
	// Priority 1: Check for product-specific template
	for _, item := range items {
		mapping, err := l.svcCtx.ShippingRepo.FindMappingByTarget(l.ctx, l.svcCtx.DB, shipping.TargetTypeProduct, item.productID)
		if err == nil && mapping != nil {
			template, err := l.svcCtx.ShippingRepo.FindByID(l.ctx, l.svcCtx.DB, mapping.TemplateID)
			if err == nil && template != nil && template.IsActive {
				zone := l.findZoneForAddress(int64(template.ID), address)
				if zone != nil {
					return template, zone
				}
			}
		}
	}

	// Priority 2: Use tenant-scoped + market-scoped default template
	// (falls back to marketID=0 within the same tenant).
	defaultTemplate, err := l.svcCtx.ShippingRepo.FindDefaultByMarket(l.ctx, l.svcCtx.DB, tenantID, marketID)
	if err == nil && defaultTemplate != nil {
		zone := l.findZoneForAddress(int64(defaultTemplate.ID), address)
		if zone != nil {
			return defaultTemplate, zone
		}
	}

	return nil, nil
}

// findZoneForAddress picks the best zone for a destination address using the
// multi-level ZoneMatcher: exact (province/city) > country fallback. The
// matcher sorts zones by Sort ascending, so the first exact match wins; if
// none, the lowest-sort zone whose regions include the country code wins.
// Returns nil when no zone matches and the template has no zones at all.
func (l *CalculateShippingFeeLogic) findZoneForAddress(templateID int64, address types.CalculatorAddress) *shipping.ShippingZone {
	zones, err := l.svcCtx.ShippingRepo.FindZonesByTemplateID(l.ctx, l.svcCtx.DB, templateID)
	if err != nil || len(zones) == 0 {
		return nil
	}
	matcher := shipping.NewZoneMatcher(zones)
	return matcher.Match(shipping.MatchInput{
		CountryCode:  address.CountryCode,
		ProvinceCode: address.ProvinceCode,
		CityCode:     address.CityCode,
		PostalCode:   address.PostalCode,
	})
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

// acceptLanguageKey is the local ctx key under which the handler injects the
// raw HTTP Accept-Language header value (e.g. "en-US,en;q=0.9,zh;q=0.8").
// We store the whole header string verbatim — ResolveZoneName takes the first
// BCP-47 tag from it via language.Parse. Storing verbatim means q-values and
// trailing tags are ignored for now, matching what most locale resolvers do
// with Accept-Language in practice.
type ctxKey string

const acceptLanguageKey ctxKey = "accept-language"

// getAcceptLanguage returns the Accept-Language header that the handler
// injected into ctx, or "" when absent. Empty string is fine — the resolver
// treats it as "no locale signal" and falls back to zone.Name.
func getAcceptLanguage(ctx context.Context) string {
	if v, ok := ctx.Value(acceptLanguageKey).(string); ok {
		return v
	}
	return ""
}

// formatAmount converts decimal.Decimal to 2-decimal string.
func formatAmount(amount decimal.Decimal) string {
	return amount.StringFixed(2)
}
