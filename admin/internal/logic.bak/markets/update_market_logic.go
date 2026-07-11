package markets

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/domain/market"
	"github.com/colinrs/shopjoy/admin/internal/infrastructure/persistence"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMarketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMarketLogic {
	return &UpdateMarketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMarketLogic) UpdateMarket(req *types.UpdateMarketReq) (resp *types.MarketResponse, err error) {
	repo := persistence.NewMarketRepository()
	m, err := repo.FindByID(l.ctx, l.svcCtx.DB, req.ID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		m.Name = req.Name
	}
	if req.IsActive != nil {
		if *req.IsActive {
			m.Activate()
		} else {
			m.Deactivate()
		}
	}
	if req.IsDefault != nil && *req.IsDefault {
		// Clear existing default market first
		if err := repo.ClearDefault(l.ctx, l.svcCtx.DB, m.TenantID); err != nil {
			return nil, err
		}
		m.SetAsDefault()
	}
	if req.TaxRules.VatRate != "" || req.TaxRules.GstRate != "" {
		m.TaxRules = market.TaxConfig{
			VATRate:     parseDecimal(req.TaxRules.VatRate),
			GSTRate:     parseDecimal(req.TaxRules.GstRate),
			IOSSEnabled: req.TaxRules.IossEnabled,
			IncludeTax:  req.TaxRules.IncludeTax,
		}
	}

	if err := repo.Update(l.ctx, l.svcCtx.DB, m); err != nil {
		return nil, err
	}

	return toMarketResponse(m), nil
}
