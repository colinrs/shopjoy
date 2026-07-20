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
// Wire-level design note: UpdateShippingZoneReq uses value-type (non-pointer)
// optional fields, so we apply zero-value detection. This means callers
// cannot explicitly set a field to its zero value (e.g. Taxable=false after
// it was true). If that becomes required, switch the wire type to pointer
// fields. This task preserves the brief's pattern: `if req.X != zero {...}`.
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
//   wire.Taxable             → entity.Taxable              (true only — false is silent no-op)
//   wire.TaxRate             → entity.TaxRate              (parseAmount)
//   wire.TaxIncluded         → entity.TaxIncluded          (true only — false is silent no-op)
//   wire.IossApplicable      → entity.IossApplicable       (true only — false is silent no-op)
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
	// Taxable/TaxIncluded/IossApplicable: zero-value (false) is a no-op.
	// Callers cannot disable these via this endpoint unless the wire type
	// is changed to pointer-bool. Documented limitation.
	if req.Taxable {
		zone.Taxable = true
	}
	if req.TaxRate != "" {
		zone.TaxRate = parseAmount(req.TaxRate)
	}
	if req.TaxIncluded {
		zone.TaxIncluded = true
	}
	if req.IossApplicable {
		zone.IossApplicable = true
	}
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
