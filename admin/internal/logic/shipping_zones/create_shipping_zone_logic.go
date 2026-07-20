package shipping_zones

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

type CreateShippingZoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateShippingZoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateShippingZoneLogic {
	return CreateShippingZoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateShippingZone wires ALL CreateShippingZoneReq fields into the
// shipping.ShippingZone entity. The Task 1.7 review found that previous
// logic silently dropped 9 new fields (Currency / NameI18n / Taxable /
// TaxRate / TaxIncluded / IossApplicable / RemoteSurcharge /
// RemoteZipPatterns / FuelSurchargePct / VolumetricDivisor) — this
// implementation explicitly maps every field and documents the wire→entity
// mapping in a comment block.
//
// ─── wire → entity field map (anti-silent-drop guard) ───
//   wire.TemplateID          → entity.TemplateID           (required, path)
//   wire.Currency            → entity.Currency             (optional, default "CNY")
//   wire.Name                → entity.Name                 (required)
//   wire.NameI18n ([]Entry)  → entity.NameI18n (StringI18n) via toStringI18n()
//   wire.Regions             → entity.Regions              (required)
//   wire.FeeType             → entity.FeeType              (required, validated IsValidV2)
//   wire.FirstUnit           → entity.FirstUnit            (int)
//   wire.FirstFee            → entity.FirstFee             (string → decimal via parseAmount)
//   wire.AdditionalUnit      → entity.AdditionalUnit       (int)
//   wire.AdditionalFee       → entity.AdditionalFee        (string → decimal via parseAmount)
//   wire.FreeThresholdAmount → entity.FreeThresholdAmount  (string → decimal via parseAmount)
//   wire.FreeThresholdCount  → entity.FreeThresholdCount   (int)
//   wire.Taxable             → entity.Taxable              (bool, P1-6)
//   wire.TaxRate             → entity.TaxRate              (string → decimal via parseAmount)
//   wire.TaxIncluded         → entity.TaxIncluded          (bool, P1-6)
//   wire.IossApplicable      → entity.IossApplicable       (bool, P1-6)
//   wire.RemoteSurcharge     → entity.RemoteSurcharge      (string → decimal via parseAmount, P1-7)
//   wire.RemoteZipPatterns   → entity.RemoteZipPatterns    ([]string → StringArray, P1-7)
//   wire.FuelSurchargePct    → entity.FuelSurchargePct     (string → decimal via parseAmount, P1-8)
//   wire.VolumetricDivisor   → entity.VolumetricDivisor    (int, P1-9; default 5000)
//   wire.Sort                → entity.Sort                 (int)
//   from ctx                 → entity.TenantID
func (l *CreateShippingZoneLogic) CreateShippingZone(req *types.CreateShippingZoneReq) (resp *types.ShippingZoneDetail, err error) {
	// Resolve tenantID from context (REQUIRED for multi-tenancy).
	tenantID, _ := contextx.GetTenantID(l.ctx)

	// Validate fee_type (must accept P1-9 by_volume).
	feeType := shipping.FeeType(req.FeeType)
	if !feeType.IsValidV2() {
		return nil, code.ErrShippingZoneInvalidFeeType
	}

	// Validate taxable/tax_rate (P1-6): rate in [0, 1] (stored as decimal(5,4)).
	if req.Taxable && req.TaxRate != "" {
		rate, err := decimal.NewFromString(req.TaxRate)
		if err != nil || rate.IsNegative() || rate.GreaterThan(decimal.NewFromInt(1)) {
			return nil, code.ErrShippingZoneFeeConfigRequired
		}
	}

	// Verify template exists (and inherit Currency/MarketID from template
	// when not supplied — wire zone Currency/MarketID optional).
	var (
		templateCurrency = defaultCurrency(req.Currency)
		marketID         int64
	)
	tpl, tplErr := l.svcCtx.ShippingRepo.FindByID(l.ctx, l.svcCtx.DB, req.TemplateID)
	if tplErr != nil {
		return nil, tplErr
	}
	if tpl != nil {
		marketID = tpl.MarketID
		// If caller didn't specify currency, inherit from template.
		if req.Currency == "" {
			templateCurrency = tpl.Currency
		}
	}

	// Build zone entity with ALL 9 new fields explicitly mapped.
	zone := &shipping.ShippingZone{
		TenantID:            tenantID,
		TemplateID:          req.TemplateID,
		MarketID:            marketID,
		Currency:            templateCurrency,
		Name:                req.Name,
		NameI18n:            toStringI18n(req.NameI18n),
		Regions:             shipping.Regions(req.Regions),
		FeeType:             feeType,
		FirstUnit:           req.FirstUnit,
		FirstFee:            parseAmount(req.FirstFee),
		AdditionalUnit:      req.AdditionalUnit,
		AdditionalFee:       parseAmount(req.AdditionalFee),
		FreeThresholdAmount: parseAmount(req.FreeThresholdAmount),
		FreeThresholdCount:  req.FreeThresholdCount,
		Taxable:             req.Taxable,
		TaxRate:             parseAmount(req.TaxRate),
		TaxIncluded:         req.TaxIncluded,
		IossApplicable:      req.IossApplicable,
		RemoteSurcharge:     parseAmount(req.RemoteSurcharge),
		RemoteZipPatterns:   shipping.StringArray(req.RemoteZipPatterns),
		FuelSurchargePct:    parseAmount(req.FuelSurchargePct),
		VolumetricDivisor:   defaultInt(req.VolumetricDivisor, shipping.DefaultVolumetricDivisor),
		Sort:                req.Sort,
	}

	// Validate zone (entity-level; uses IsValidV2 and by_volume VolumetricDivisor check).
	if err := zone.Validate(); err != nil {
		return nil, err
	}

	// Save zone.
	if err := l.svcCtx.ShippingRepo.CreateZone(l.ctx, l.svcCtx.DB, zone); err != nil {
		return nil, err
	}

	// ─── entity → response field map (must include all 22 fields) ───
	return &types.ShippingZoneDetail{
		ID:                  int64(zone.ID),
		TenantID:            zone.TenantID,
		TemplateID:          zone.TemplateID,
		MarketID:            zone.MarketID,
		Currency:            zone.Currency,
		Name:                zone.Name,
		NameI18n:            fromStringI18n(zone.NameI18n),
		Regions:             zone.Regions,
		FeeType:             string(zone.FeeType),
		FirstUnit:           zone.FirstUnit,
		FirstFee:            formatAmount(zone.FirstFee),
		AdditionalUnit:      zone.AdditionalUnit,
		AdditionalFee:       formatAmount(zone.AdditionalFee),
		FreeThresholdAmount: formatAmount(zone.FreeThresholdAmount),
		FreeThresholdCount:  zone.FreeThresholdCount,
		Taxable:             zone.Taxable,
		TaxRate:             formatAmount(zone.TaxRate),
		TaxIncluded:         zone.TaxIncluded,
		IossApplicable:      zone.IossApplicable,
		RemoteSurcharge:     formatAmount(zone.RemoteSurcharge),
		RemoteZipPatterns:   zone.RemoteZipPatterns,
		FuelSurchargePct:    formatAmount(zone.FuelSurchargePct),
		VolumetricDivisor:   zone.VolumetricDivisor,
		Sort:                zone.Sort,
	}, nil
}

// parseAmount converts a string amount to decimal.Decimal (empty → 0).
// Errors are silently coerced to 0; callers needing validation must do so
// explicitly (see TaxRate range check in Create).
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

// formatAmount converts decimal.Decimal to a fixed 2-decimal-place string.
func formatAmount(amount decimal.Decimal) string {
	return amount.StringFixed(2)
}

// defaultCurrency returns the wire-supplied currency or "CNY" if empty.
func defaultCurrency(c string) string {
	if c == "" {
		return "CNY"
	}
	return c
}

// defaultInt returns v if non-zero, otherwise def.
func defaultInt(v, def int) int {
	if v == 0 {
		return def
	}
	return v
}

// toStringI18n converts a wire []NameI18nEntry (array) into an entity
// StringI18n (map[locale]name). Empty/invalid entries are dropped.
func toStringI18n(entries []types.NameI18nEntry) shipping.StringI18n {
	if len(entries) == 0 {
		return nil
	}
	out := make(shipping.StringI18n, len(entries))
	for _, e := range entries {
		if e.Locale != "" && e.Name != "" {
			out[e.Locale] = e.Name
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// fromStringI18n converts an entity StringI18n (map) back into a wire
// []NameI18nEntry (array). Used by response builders.
func fromStringI18n(s shipping.StringI18n) []types.NameI18nEntry {
	if len(s) == 0 {
		return nil
	}
	out := make([]types.NameI18nEntry, 0, len(s))
	for locale, name := range s {
		out = append(out, types.NameI18nEntry{Locale: locale, Name: name})
	}
	return out
}
