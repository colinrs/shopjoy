package shipping

import (
	"context"
	"strings"

	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/shopspring/decimal"
)

// StandardCarrier 默认承运商实现。运费完全依据 ShippingZone 配置计算：
// 基础运费（zone.CalculateFee）+ 偏远/燃油附加费 + 税（若 zone 计税）。
type StandardCarrier struct{}

// Code 返回承运商代码。
func (StandardCarrier) Code() string { return "standard" }

// Name 返回承运商展示名称。
func (StandardCarrier) Name() string { return "Standard Shipping" }

// Quote 计算报价：base = zone.CalculateFee → surcharge（基于 base + 邮编）→ tax（若 Taxable）。
// Total = base + surcharge.Total + tax。zone 为空返回 ErrShippingZoneNotFound。
func (c StandardCarrier) Quote(_ context.Context, req QuoteRequest) (*QuoteResult, error) {
	zone := req.Zone
	if zone == nil {
		return nil, code.ErrShippingZoneNotFound
	}

	// 1. 基础运费
	base := zone.CalculateFee(req.Items, req.OrderAmount, req.ItemCount)

	// 2. 附加费（偏远 + 燃油），以基础费为计算依据
	surcharges := CalculateSurcharges(zone, SurchargeInput{
		BaseFee:    base,
		PostalCode: req.Address.PostalCode,
	})

	feeBeforeTax := base.Add(surcharges.Total)

	// 3. 税（仅当 zone 计税）
	var tax decimal.Decimal
	total := feeBeforeTax
	if zone.Taxable {
		tax, total = CalculateTax(feeBeforeTax, zone.TaxRate, zone.TaxIncluded)
	}

	// 4. 计费重量（取所有 item 计费重之和）
	divisor := zone.VolumetricDivisor
	if divisor <= 0 {
		divisor = DefaultVolumetricDivisor
	}
	var weight int
	for _, item := range req.Items {
		weight += item.ChargeableWeight(divisor) * item.Quantity
	}

	return &QuoteResult{
		BaseFee:       base,
		Surcharges:    surcharges,
		Tax:           tax,
		Total:         total,
		Currency:      zone.Currency,
		EstimatedDays: c.EstimateDays(req.Address.CountryCode),
		Weight:        weight,
	}, nil
}

// EstimateDays 依据目的国家分档估算配送时效（简化实现）。
//
//	CN                    → 3
//	US / EU 主要国家       → 7
//	JP / KR / 东南亚       → 5
//	其他                  → 10
func (StandardCarrier) EstimateDays(destinationCountry string) int {
	switch strings.ToUpper(destinationCountry) {
	case "CN":
		return 3
	case "US", "DE", "FR", "GB", "IT", "ES", "NL":
		return 7
	case "JP", "KR", "SG", "MY", "TH", "VN", "ID", "PH":
		return 5
	default:
		return 10
	}
}

// CarrierRegistry 承运商注册表，按 code 管理可用承运商。
type CarrierRegistry struct {
	carriers map[string]Carrier
}

// NewCarrierRegistry 创建注册表，默认注册 StandardCarrier。
func NewCarrierRegistry() *CarrierRegistry {
	r := &CarrierRegistry{carriers: map[string]Carrier{}}
	r.Register(StandardCarrier{})
	return r
}

// Register 注册承运商（按 Code 索引，重复注册覆盖）。
func (r *CarrierRegistry) Register(c Carrier) {
	r.carriers[c.Code()] = c
}

// Get 按 code 取承运商，未注册返回 (nil, false)。
func (r *CarrierRegistry) Get(code string) (Carrier, bool) {
	c, ok := r.carriers[code]
	return c, ok
}
