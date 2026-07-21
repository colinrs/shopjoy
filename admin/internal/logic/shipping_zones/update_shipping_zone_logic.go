package shipping_zones

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/domain/shipping"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"

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

// UpdateShippingZone wires ALL UpdateShippingZoneReq fields into the
// existing shipping.ShippingZone entity. Each field is conditionally
// assigned — only updates that the caller explicitly provided.
//
// Wire-level design note: UpdateShippingZoneReq uses pointer fields for
// the 3 boolean toggles (Taxable, TaxIncluded, IossApplicable). The old
// non-pointer `bool` form silently dropped a caller-side `false` when the
// current value was `true` (because `if req.Taxable` only fires on true).
// Switching to `*bool` lets the caller explicitly set these to false.
//
// ─── wire → entity field map (anti-silent-drop guard) ───
//   wire.Name                → entity.Name                 (if != "")
//   wire.NameI18n ([]Entry)  → entity.NameI18n (StringI18n) via toStringI18n() (if != nil)
//   wire.Regions             → entity.Regions              (if != nil)
//   wire.FeeType             → entity.FeeType              (if != "")
//   wire.FirstUnit           → entity.FirstUnit            (if != 0)
//   wire.FirstFee            → entity.FirstFee             (parseAmount)
//   wire.AdditionalUnit      → entity.AdditionalUnit       (if != 0)
//   wire.AdditionalFee       → entity.AdditionalFee        (parseAmount)
//   wire.FreeThresholdAmount → entity.FreeThresholdAmount  (parseAmount)
//   wire.FreeThresholdCount  → entity.FreeThresholdCount   (if != 0)
//   wire.Taxable             → entity.Taxable              (if *bool != nil — explicit false OK)
//   wire.TaxRate             → entity.TaxRate              (parseAmount)
//   wire.TaxIncluded         → entity.TaxIncluded          (if *bool != nil — explicit false OK)
//   wire.IossApplicable      → entity.IossApplicable       (if *bool != nil — explicit false OK)
//   wire.RemoteSurcharge     → entity.RemoteSurcharge      (parseAmount)
//   wire.RemoteZipPatterns   → entity.RemoteZipPatterns    ([]string → StringArray, if != nil)
//   wire.FuelSurchargePct    → entity.FuelSurchargePct     (parseAmount)
//   wire.VolumetricDivisor   → entity.VolumetricDivisor    (if != 0)
//   wire.Sort                → entity.Sort                 (if != 0)
func (l *UpdateShippingZoneLogic) UpdateShippingZone(req *types.UpdateShippingZoneReq) (resp *types.ShippingZoneDetail, err error) {
	// Find existing zone.
	zone, err := l.svcCtx.ShippingRepo.FindZoneByID(l.ctx, l.svcCtx.DB, req.ID)
	if err != nil {
		return nil, err
	}
	if zone == nil {
		return nil, code.ErrShippingZoneNotFound
	}

	// ─── conditional wire → entity field updates ───
	if req.Name != "" {
		zone.Name = req.Name
	}
	if req.NameI18n != nil {
		zone.NameI18n = toStringI18n(req.NameI18n)
	}
	if req.Regions != nil {
		zone.Regions = shipping.Regions(req.Regions)
	}
	if req.FeeType != "" {
		newFeeType := shipping.FeeType(req.FeeType)
		if !newFeeType.IsValidV2() {
			return nil, code.ErrShippingZoneInvalidFeeType
		}
		zone.FeeType = newFeeType
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
	// Taxable/TaxIncluded/IossApplicable: pointer fields on the wire so that
	// callers can EXPLICITLY set these to false. Old non-pointer bool fields
	// silently dropped a `false` when the current value was `true` because
	// `if req.Taxable` only fired on `true`. Pointer form: nil=skip, non-nil
	// dereference wins. This is the Important/bool-pointer fix and is
	// implemented as a pure helper (applyBoolPtrOverride) so the explicit-
	// false contract can be unit-tested without a DB.
	zone.Taxable = applyBoolPtrOverride(zone.Taxable, req.Taxable)
	if req.TaxRate != "" {
		zone.TaxRate = parseAmount(req.TaxRate)
	}
	zone.TaxIncluded = applyBoolPtrOverride(zone.TaxIncluded, req.TaxIncluded)
	zone.IossApplicable = applyBoolPtrOverride(zone.IossApplicable, req.IossApplicable)
	if req.RemoteSurcharge != "" {
		zone.RemoteSurcharge = parseAmount(req.RemoteSurcharge)
	}
	if req.RemoteZipPatterns != nil {
		zone.RemoteZipPatterns = shipping.StringArray(req.RemoteZipPatterns)
	}
	if req.FuelSurchargePct != "" {
		zone.FuelSurchargePct = parseAmount(req.FuelSurchargePct)
	}
	if req.VolumetricDivisor != 0 {
		zone.VolumetricDivisor = req.VolumetricDivisor
	}
	if req.Sort != 0 {
		zone.Sort = req.Sort
	}

	// Validate zone (entity-level; uses IsValidV2 and by_volume VolumetricDivisor check).
	if err := zone.Validate(); err != nil {
		return nil, err
	}

	// Save changes.
	if err := l.svcCtx.ShippingRepo.UpdateZone(l.ctx, l.svcCtx.DB, zone); err != nil {
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

// applyBoolPtrOverride returns the new bool value for an updatable field,
// preserving a *bool → bool wire contract:
//   - nil pointer → keep current (caller did not supply the field).
//   - non-nil pointer → caller explicitly chose the dereferenced value.
//
// The Important/bool-pointer fix changes three tax/fulfillment toggles on
// UpdateShippingZoneReq from `bool` to `*bool`. Old form: callers could not
// explicitly disable a flag that was currently true (because `if req.X`
// only fired on true, dropping false-on-true). New form: caller owns the
// value, including the explicit false. This helper centralises that
// contract so it can be unit-tested without a DB.
func applyBoolPtrOverride(current bool, override *bool) bool {
	if override == nil {
		return current
	}
	return *override
}
