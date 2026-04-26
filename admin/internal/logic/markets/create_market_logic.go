package markets

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/market"
	"github.com/colinrs/shopjoy/admin/internal/infrastructure/persistence"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMarketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMarketLogic {
	return &CreateMarketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateMarketLogic) CreateMarket(req *types.CreateMarketReq) (resp *types.MarketResponse, err error) {
	// Create domain entity
	defaultLang := req.DefaultLanguage
	if defaultLang == "" {
		defaultLang = "en"
	}

	m, err := market.NewMarket(req.Code, req.Name, req.Currency, defaultLang)
	if err != nil {
		return nil, err
	}

	// Set optional fields
	m.Flag = req.Flag
	m.TaxRules = market.TaxConfig{
		VATRate:     parseDecimal(req.TaxRules.VatRate),
		GSTRate:     parseDecimal(req.TaxRules.GstRate),
		IOSSEnabled: req.TaxRules.IossEnabled,
		IncludeTax:  req.TaxRules.IncludeTax,
	}

	// Persist
	repo := persistence.NewMarketRepository()

	// Check if market with same code already exists
	existing, err := repo.FindByCode(l.ctx, l.svcCtx.DB, req.Code)
	if err == nil && existing != nil {
		return nil, code.ErrMarketDuplicate
	}

	if err := repo.Create(l.ctx, l.svcCtx.DB, m); err != nil {
		return nil, code.ErrMarketDuplicate
	}

	return toMarketResponse(m), nil
}

func parseDecimal(s string) decimal.Decimal {
	if s == "" {
		return decimal.Zero
	}
	d, _ := decimal.NewFromString(s)
	return d
}

func toMarketResponse(m *market.Market) *types.MarketResponse {
	return &types.MarketResponse{
		ID:              int64(m.ID),
		Code:            m.Code,
		Name:            m.Name,
		Currency:        m.Currency,
		DefaultLanguage: m.DefaultLanguage,
		Flag:            m.Flag,
		IsActive:        m.IsActive,
		IsDefault:       m.IsDefault,
		TaxRules: types.TaxConfig{
			VatRate:     m.TaxRules.VATRate.String(),
			GstRate:     m.TaxRules.GSTRate.String(),
			IossEnabled: m.TaxRules.IOSSEnabled,
			IncludeTax:  m.TaxRules.IncludeTax,
		},
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
		UpdatedAt: m.UpdatedAt.Format(time.RFC3339),
	}
}
