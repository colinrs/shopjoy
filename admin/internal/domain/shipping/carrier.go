package shipping

import (
	"context"

	"github.com/shopspring/decimal"
)

// Carrier 承运商抽象。不同承运商可提供不同的运费报价与时效估算策略。
// StandardCarrier 是默认实现（基于 ShippingZone 配置计算）。
type Carrier interface {
	// Code 返回承运商唯一代码（如 "standard"），与 ShippingTemplate.CarrierCode 对应。
	Code() string
	// Name 返回承运商展示名称。
	Name() string
	// Quote 依据请求给出运费报价（基础费 + 附加费 + 税）。
	Quote(ctx context.Context, req QuoteRequest) (*QuoteResult, error)
	// EstimateDays 依据目的国家估算配送时效（天）。
	EstimateDays(destinationCountry string) int
}

// QuoteRequest 运费报价请求。
type QuoteRequest struct {
	TemplateID  int64
	Zone        *ShippingZone
	Items       []CalculateItem
	OrderAmount decimal.Decimal
	ItemCount   int
	Address     MatchInput
}

// QuoteResult 运费报价结果。Total = BaseFee + Surcharges.Total + Tax（税外加时）。
type QuoteResult struct {
	BaseFee       decimal.Decimal
	Surcharges    SurchargeBreakdown
	Tax           decimal.Decimal
	Total         decimal.Decimal
	Currency      string
	EstimatedDays int
	Weight        int // 计费重（克）
}
