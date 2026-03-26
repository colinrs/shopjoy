// Package market 市场领域层
package market

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// TaxConfig 税务配置
type TaxConfig struct {
	VATRate     decimal.Decimal `json:"VATRate"`     // 增值税率 (EU)
	GSTRate     decimal.Decimal `json:"GSTRate"`     // 商品及服务税税率 (AU)
	IOSSEnabled bool            `json:"IOSSEnabled"` // 低值商品进口方案
	IncludeTax  bool            `json:"IncludeTax"`  // 价格是否含税显示
}

// Market 市场实体
type Market struct {
	ID              int64          // 市场ID
	TenantID        int64          // 租户ID
	Code            string         // 市场代码: US, UK, DE, FR, AU
	Name            string         // 市场名称
	Currency        string         // 货币: USD, GBP, EUR, AUD
	DefaultLanguage string         // 默认语言
	Flag            string         // 旗帜图标
	IsActive        bool           // 是否启用
	IsDefault       bool           // 是否主市场
	TaxRules        TaxConfig      `gorm:"type:json"` // 税务配置
	CreatedAt       time.Time      // 创建时间
	UpdatedAt       time.Time      // 更新时间
	DeletedAt       *int64        // 软删除时间
}

// TableName 表名
func (m *Market) TableName() string {
	return "markets"
}

// NewMarket 创建新市场
func NewMarket(marketCode, name, currency, defaultLanguage string) (*Market, error) {
	if marketCode == "" {
		return nil, code.ErrMarketCodeRequired
	}
	if name == "" {
		return nil, code.ErrMarketNameRequired
	}
	if currency == "" {
		return nil, code.ErrMarketCurrencyRequired
	}

	now := time.Now()
	return &Market{
		Code:            marketCode,
		Name:            name,
		Currency:        currency,
		DefaultLanguage: defaultLanguage,
		IsActive:        true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

// Activate 激活市场
func (m *Market) Activate() {
	m.IsActive = true
	m.UpdatedAt = time.Now()
}

// Deactivate 停用市场
func (m *Market) Deactivate() {
	m.IsActive = false
	m.UpdatedAt = time.Now()
}

// SetAsDefault 设置为主市场
func (m *Market) SetAsDefault() {
	m.IsDefault = true
	m.UpdatedAt = time.Now()
}

// Repository 市场仓储接口
type Repository interface {
	Create(ctx context.Context, db *gorm.DB, market *Market) error
	Update(ctx context.Context, db *gorm.DB, market *Market) error
	Delete(ctx context.Context, db *gorm.DB, id int64) error
	FindByID(ctx context.Context, db *gorm.DB, id int64) (*Market, error)
	FindByIDs(ctx context.Context, db *gorm.DB, ids []int64) ([]*Market, error)
	FindByCode(ctx context.Context, db *gorm.DB, code string) (*Market, error)
	FindAll(ctx context.Context, db *gorm.DB) ([]*Market, error)
	FindActive(ctx context.Context, db *gorm.DB) ([]*Market, error)
	FindDefault(ctx context.Context, db *gorm.DB) (*Market, error)
	ClearDefault(ctx context.Context, db *gorm.DB, tenantID int64) error
}

// Query 查询条件
type Query struct {
	Code     string
	IsActive *bool
	Page     int
	PageSize int
}

// Offset 计算偏移量
func (q Query) Offset() int {
	if q.Page < 1 {
		q.Page = 1
	}
	return (q.Page - 1) * q.PageSize
}

// Limit 返回限制数
func (q Query) Limit() int {
	if q.PageSize < 1 || q.PageSize > 100 {
		q.PageSize = 20
	}
	return q.PageSize
}
