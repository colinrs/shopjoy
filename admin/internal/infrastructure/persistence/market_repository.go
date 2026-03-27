package persistence

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/colinrs/shopjoy/admin/internal/domain/market"
	"github.com/colinrs/shopjoy/pkg/application"
	"github.com/colinrs/shopjoy/pkg/code"
	"gorm.io/gorm"
)

type marketRepo struct{}

func NewMarketRepository() market.Repository {
	return &marketRepo{}
}

type marketModel struct {
	application.Model
	TenantID        int64  `gorm:"column:tenant_id;not null;default:0"`
	Code            string `gorm:"column:code;size:10;not null;uniqueIndex:uk_tenant_code"`
	Name            string `gorm:"column:name;size:64;not null"`
	Currency        string `gorm:"column:currency;size:10;not null"`
	DefaultLanguage string `gorm:"column:default_language;size:10;not null;default:'en'"`
	Flag            string `gorm:"column:flag;size:32"`
	IsActive        bool   `gorm:"column:is_active;not null;default:true"`
	IsDefault       bool   `gorm:"column:is_default;not null;default:false"`
	TaxRules        string `gorm:"column:tax_rules;type:json"`
}

func (marketModel) TableName() string {
	return "markets"
}

func (m *marketModel) toEntity() *market.Market {
	var taxRules market.TaxConfig
	if m.TaxRules != "" {
		json.Unmarshal([]byte(m.TaxRules), &taxRules)
	}

	entity := &market.Market{
		TenantID:        m.TenantID,
		Code:            m.Code,
		Name:            m.Name,
		Currency:        m.Currency,
		DefaultLanguage: m.DefaultLanguage,
		Flag:            m.Flag,
		IsActive:        m.IsActive,
		IsDefault:       m.IsDefault,
		TaxRules:        taxRules,
	}
	entity.ID = m.ID
	entity.CreatedAt = m.CreatedAt
	entity.UpdatedAt = m.UpdatedAt
	return entity
}

func fromMarketEntity(m *market.Market) *marketModel {
	taxRulesJSON, _ := json.Marshal(m.TaxRules)

	model := &marketModel{
		TenantID:        m.TenantID,
		Code:            m.Code,
		Name:            m.Name,
		Currency:        m.Currency,
		DefaultLanguage: m.DefaultLanguage,
		Flag:            m.Flag,
		IsActive:        m.IsActive,
		IsDefault:       m.IsDefault,
		TaxRules:        string(taxRulesJSON),
	}
	model.ID = m.ID
	model.CreatedAt = m.CreatedAt
	model.UpdatedAt = m.UpdatedAt
	return model
}

func (r *marketRepo) Create(ctx context.Context, db *gorm.DB, m *market.Market) error {
	model := fromMarketEntity(m)
	return db.WithContext(ctx).Create(model).Error
}

func (r *marketRepo) Update(ctx context.Context, db *gorm.DB, m *market.Market) error {
	model := fromMarketEntity(m)
	return db.WithContext(ctx).Save(model).Error
}

func (r *marketRepo) Delete(ctx context.Context, db *gorm.DB, id int64) error {
	result := db.WithContext(ctx).Delete(&marketModel{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return code.ErrMarketNotFound
	}
	return nil
}

func (r *marketRepo) FindByID(ctx context.Context, db *gorm.DB, id int64) (*market.Market, error) {
	var model marketModel
	err := db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrMarketNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *marketRepo) FindByIDs(ctx context.Context, db *gorm.DB, ids []int64) ([]*market.Market, error) {
	if len(ids) == 0 {
		return []*market.Market{}, nil
	}

	var models []marketModel
	err := db.WithContext(ctx).Where("id IN ?", ids).Find(&models).Error
	if err != nil {
		return nil, err
	}

	markets := make([]*market.Market, len(models))
	for i, m := range models {
		markets[i] = m.toEntity()
	}
	return markets, nil
}

func (r *marketRepo) FindByCode(ctx context.Context, db *gorm.DB, codeStr string) (*market.Market, error) {
	var model marketModel
	err := db.WithContext(ctx).Where("code = ?", codeStr).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrMarketNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *marketRepo) FindAll(ctx context.Context, db *gorm.DB) ([]*market.Market, error) {
	var models []marketModel
	err := db.WithContext(ctx).Order("is_default DESC, code ASC").Find(&models).Error
	if err != nil {
		return nil, err
	}

	markets := make([]*market.Market, len(models))
	for i, m := range models {
		markets[i] = m.toEntity()
	}
	return markets, nil
}

func (r *marketRepo) FindActive(ctx context.Context, db *gorm.DB) ([]*market.Market, error) {
	var models []marketModel
	err := db.WithContext(ctx).Where("is_active = ?", true).Order("is_default DESC, code ASC").Find(&models).Error
	if err != nil {
		return nil, err
	}

	markets := make([]*market.Market, len(models))
	for i, m := range models {
		markets[i] = m.toEntity()
	}
	return markets, nil
}

func (r *marketRepo) FindDefault(ctx context.Context, db *gorm.DB) (*market.Market, error) {
	var model marketModel
	err := db.WithContext(ctx).Where("is_default = ?", true).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrMarketNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *marketRepo) ClearDefault(ctx context.Context, db *gorm.DB, tenantID int64) error {
	return db.WithContext(ctx).Model(&marketModel{}).
		Where("tenant_id = ? AND is_default = ?", tenantID, true).
		Update("is_default", false).Error
}
