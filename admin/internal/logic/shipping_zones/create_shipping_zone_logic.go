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

func (l *CreateShippingZoneLogic) CreateShippingZone(req *types.CreateShippingZoneReq) (resp *types.ShippingZoneDetail, err error) {
	tenantID, err := contextx.MustGetTenantIDForLogic(l.ctx)
	if err != nil {
		l.Logger.Errorf("failed to get tenant ID: %v", err)
		return nil, err
	}

	// Verify template exists and belongs to tenant
	template, err := l.svcCtx.ShippingRepo.FindByID(l.ctx, l.svcCtx.DB, tenantID, req.TemplateID)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, code.ErrShippingTemplateNotFound
	}

	// Parse fee amounts from string to cents
	firstFee := parseAmount(req.FirstFee)
	additionalFee := parseAmount(req.AdditionalFee)
	freeThresholdAmount := parseAmount(req.FreeThresholdAmount)

	// Create zone entity
	zone := &shipping.ShippingZone{
		TenantID:            tenantID,
		TemplateID:          req.TemplateID,
		Name:                req.Name,
		Regions:             req.Regions,
		FeeType:             shipping.FeeType(req.FeeType),
		FirstUnit:           req.FirstUnit,
		FirstFee:            firstFee,
		AdditionalUnit:      req.AdditionalUnit,
		AdditionalFee:       additionalFee,
		FreeThresholdAmount: freeThresholdAmount,
		FreeThresholdCount:  req.FreeThresholdCount,
		Sort:                req.Sort,
	}

	// Validate zone
	if err := zone.Validate(); err != nil {
		return nil, err
	}

	// Save zone
	if err := l.svcCtx.ShippingRepo.CreateZone(l.ctx, l.svcCtx.DB, zone); err != nil {
		return nil, err
	}

	return &types.ShippingZoneDetail{
		ID:                  int64(zone.ID),
		TemplateID:          zone.TemplateID,
		Name:                zone.Name,
		Regions:             zone.Regions,
		FeeType:             string(zone.FeeType),
		FirstUnit:           zone.FirstUnit,
		FirstFee:            formatAmount(zone.FirstFee),
		AdditionalUnit:      zone.AdditionalUnit,
		AdditionalFee:       formatAmount(zone.AdditionalFee),
		FreeThresholdAmount: formatAmount(zone.FreeThresholdAmount),
		FreeThresholdCount:  zone.FreeThresholdCount,
		Sort:                zone.Sort,
	}, nil
}

// parseAmount converts string amount to decimal.Decimal
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