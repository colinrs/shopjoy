# Cross-Border Product Management - Phase 1: Core Infrastructure

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Establish the core data models and APIs for multi-market product management (Market entity, Product compliance fields, ProductMarket association).

**Architecture:** DDD pattern following existing codebase structure. Domain layer defines entities and repository interfaces. Infrastructure layer implements persistence. API layer uses go-zero handlers/logic pattern.

**Tech Stack:** Go 1.24, go-zero, GORM, MySQL, github.com/shopspring/decimal

---

## File Structure

```
admin/
├── migrations/
│   └── 004_multi_market_support.sql     # NEW: Database migration
├── internal/
│   ├── domain/
│   │   ├── market/
│   │   │   └── entity.go                 # NEW: Market entity + repository interface
│   │   └── product/
│   │       └── entity.go                 # MODIFY: Add compliance fields
│   ├── infrastructure/
│   │   └── persistence/
│   │       ├── market_repository.go      # NEW: Market repository implementation
│   │       └── product_repository.go     # MODIFY: Update for new fields
│   └── types/
│       └── types.go                      # AUTO-GEN: API types
├── desc/
│   └── market.api                        # NEW: Market API definition
└── etc/
    └── admin-api.yaml                    # No changes needed

pkg/
└── code/
    └── code.go                           # MODIFY: Add Market error codes
```

---

## Task 1: Database Migration

**Files:**
- Create: `migrations/004_multi_market_support.sql`

- [ ] **Step 1: Create migration file with tables**

```sql
-- migrations/004_multi_market_support.sql

-- Markets table: Market configuration entity
CREATE TABLE IF NOT EXISTS `markets` (
    `id`               bigint(20)      NOT NULL AUTO_INCREMENT COMMENT 'Market ID',
    `tenant_id`        bigint(20)      NOT NULL DEFAULT 0 COMMENT 'Tenant ID, 0 = global markets',
    `code`             varchar(10)     NOT NULL COMMENT 'Market code: US, UK, DE, FR, AU',
    `name`             varchar(64)     NOT NULL COMMENT 'Market name: United States',
    `currency`         varchar(10)     NOT NULL COMMENT 'Currency: USD, GBP, EUR, AUD',
    `default_language` varchar(10)     NOT NULL DEFAULT 'en' COMMENT 'Default language: en, de, fr',
    `flag`             varchar(32)     NULL COMMENT 'Flag emoji or image URL',
    `is_active`        tinyint(1)      NOT NULL DEFAULT 1 COMMENT '1 = active, 0 = inactive',
    `is_default`       tinyint(1)      NOT NULL DEFAULT 0 COMMENT '1 = primary market',
    `tax_rules`        json            NULL COMMENT 'Tax configuration: VAT, GST, IOSS',
    `created_at`       datetime        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`       datetime        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at`       datetime        NULL COMMENT 'Soft delete timestamp',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant_code` (`tenant_id`, `code`),
    KEY `idx_is_active` (`is_active`),
    KEY `idx_is_default` (`is_default`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Markets table';

-- Product markets table: Market-specific product data
CREATE TABLE IF NOT EXISTS `product_markets` (
    `id`                    bigint(20)      NOT NULL AUTO_INCREMENT COMMENT 'Record ID',
    `tenant_id`             bigint(20)      NOT NULL COMMENT 'Tenant ID',
    `product_id`            bigint(20)      NOT NULL COMMENT 'Product ID',
    `variant_id`            bigint(20)      NULL COMMENT 'Variant ID, NULL for base product',
    `market_id`             bigint(20)      NOT NULL COMMENT 'Market ID',
    `is_enabled`            tinyint(1)      NOT NULL DEFAULT 0 COMMENT 'Product visible in this market',
    `status_override`       tinyint(4)      NULL COMMENT 'Override product status per market',
    `price`                 decimal(19,4)   NOT NULL DEFAULT 0.0000 COMMENT 'Market-specific price',
    `compare_at_price`      decimal(19,4)   NULL COMMENT 'Compare at price for sales',
    `stock_alert_threshold` int(11)         NOT NULL DEFAULT 0 COMMENT 'Low stock alert threshold',
    `published_at`          datetime        NULL COMMENT 'Published timestamp in this market',
    `created_at`            datetime        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`            datetime        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_product_variant_market` (`product_id`, `variant_id`, `market_id`),
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_market_id` (`market_id`),
    KEY `idx_is_enabled` (`is_enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Product-Market association table';

-- Add compliance fields to products table
ALTER TABLE `products`
    ADD COLUMN `sku` varchar(64) NULL COMMENT 'SKU code' AFTER `name`,
    ADD COLUMN `brand` varchar(64) NULL COMMENT 'Brand name' AFTER `category_id`,
    ADD COLUMN `tags` json NULL COMMENT 'Product tags' AFTER `brand`,
    ADD COLUMN `images` json NULL COMMENT 'Product images' AFTER `tags`,
    ADD COLUMN `is_matrix_product` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'Has variants' AFTER `status`,
    ADD COLUMN `hs_code` varchar(20) NULL COMMENT 'Harmonized System Code' AFTER `is_matrix_product`,
    ADD COLUMN `coo` varchar(10) NULL COMMENT 'Country of Origin' AFTER `hs_code`,
    ADD COLUMN `weight` decimal(10,2) NULL COMMENT 'Weight' AFTER `coo`,
    ADD COLUMN `weight_unit` varchar(10) NULL DEFAULT 'g' COMMENT 'Weight unit: g, kg' AFTER `weight`,
    ADD COLUMN `length` decimal(10,2) NULL COMMENT 'Package length (cm)' AFTER `weight_unit`,
    ADD COLUMN `width` decimal(10,2) NULL COMMENT 'Package width (cm)' AFTER `length`,
    ADD COLUMN `height` decimal(10,2) NULL COMMENT 'Package height (cm)' AFTER `width`,
    ADD COLUMN `dangerous_goods` json NULL COMMENT 'Dangerous goods flags' AFTER `height`;

-- Add unique index on SKU
ALTER TABLE `products` ADD UNIQUE KEY `uk_sku` (`sku`);

-- Insert default markets for MVP
INSERT INTO `markets` (`tenant_id`, `code`, `name`, `currency`, `default_language`, `flag`, `is_active`, `is_default`, `tax_rules`) VALUES
(0, 'US', 'United States', 'USD', 'en', '🇺🇸', 1, 1, '{"IncludeTax": false}'),
(0, 'UK', 'United Kingdom', 'GBP', 'en', '🇬🇧', 1, 0, '{"VATRate": "20", "IncludeTax": true}'),
(0, 'DE', 'Germany', 'EUR', 'de', '🇩🇪', 1, 0, '{"VATRate": "19", "IOSSEnabled": true, "IncludeTax": true}'),
(0, 'FR', 'France', 'EUR', 'fr', '🇫🇷', 1, 0, '{"VATRate": "20", "IOSSEnabled": true, "IncludeTax": true}'),
(0, 'AU', 'Australia', 'AUD', 'en', '🇦🇺', 1, 0, '{"GSTRate": "10", "IncludeTax": true}')
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);
```

- [ ] **Step 2: Verify migration file syntax**

Run: `cat migrations/004_multi_market_support.sql | head -20`
Expected: File content displayed correctly

- [ ] **Step 3: Commit migration**

```bash
git add migrations/004_multi_market_support.sql
git commit -m "feat(db): add multi-market support migration

- Add markets table for market configuration
- Add product_markets table for market-specific product data
- Add compliance fields to products (hs_code, coo, weight, dimensions)
- Insert default markets: US, UK, DE, FR, AU"
```

---

## Task 2: Market Domain Entity

**Files:**
- Create: `admin/internal/domain/market/entity.go`

- [ ] **Step 1: Create market domain package**

```go
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
	VATRate     decimal.Decimal `json:"vat_rate"`      // 增值税率 (EU)
	GSTRate     decimal.Decimal `json:"gst_rate"`      // 商品及服务税税率 (AU)
	IOSSEnabled bool            `json:"ioss_enabled"`  // 低值商品进口方案
	IncludeTax  bool            `json:"include_tax"`   // 价格是否含税显示
}

// Market 市场实体
type Market struct {
	ID              int64      // 市场ID
	TenantID        int64      // 租户ID
	Code            string     // 市场代码: US, UK, DE, FR, AU
	Name            string     // 市场名称
	Currency        string     // 货币: USD, GBP, EUR, AUD
	DefaultLanguage string     // 默认语言
	Flag            string     // 旗帜图标
	IsActive        bool       // 是否启用
	IsDefault       bool       // 是否主市场
	TaxRules        TaxConfig  `gorm:"type:json"` // 税务配置
	CreatedAt       time.Time  // 创建时间
	UpdatedAt       time.Time  // 更新时间
}

// TableName 表名
func (m *Market) TableName() string {
	return "markets"
}

// NewMarket 创建新市场
func NewMarket(code, name, currency, defaultLanguage string) (*Market, error) {
	if code == "" {
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
		Code:            code,
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
	FindByCode(ctx context.Context, db *gorm.DB, code string) (*Market, error)
	FindAll(ctx context.Context, db *gorm.DB) ([]*Market, error)
	FindActive(ctx context.Context, db *gorm.DB) ([]*Market, error)
	FindDefault(ctx context.Context, db *gorm.DB) (*Market, error)
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
```

- [ ] **Step 2: Verify file compiles**

Run: `cd admin && go build ./internal/domain/market/...`
Expected: No errors

- [ ] **Step 3: Commit**

```bash
git add admin/internal/domain/market/entity.go
git commit -m "feat(domain): add Market entity with repository interface

- Define Market entity with tax configuration
- Add business methods: Activate, Deactivate, SetAsDefault
- Define Repository interface for market CRUD"
```

---

## Task 3: Market Error Codes

**Files:**
- Modify: `pkg/code/code.go`

- [ ] **Step 1: Add market error codes**

Add after the existing error definitions (around line 130):

```go
	// ErrMarketNotFound ==================== Market Module (150xxx) ====================
	ErrMarketNotFound         = &Err{HTTPCode: http.StatusNotFound, Code: 150001, Msg: "market not found"}
	ErrMarketDuplicate        = &Err{HTTPCode: http.StatusConflict, Code: 150002, Msg: "duplicate market code"}
	ErrMarketCodeRequired     = &Err{HTTPCode: http.StatusBadRequest, Code: 150003, Msg: "market code is required"}
	ErrMarketNameRequired     = &Err{HTTPCode: http.StatusBadRequest, Code: 150004, Msg: "market name is required"}
	ErrMarketCurrencyRequired = &Err{HTTPCode: http.StatusBadRequest, Code: 150005, Msg: "market currency is required"}
	ErrMarketInactive         = &Err{HTTPCode: http.StatusBadRequest, Code: 150006, Msg: "market is inactive"}
	ErrMarketAlreadyDefault   = &Err{HTTPCode: http.StatusBadRequest, Code: 150007, Msg: "market is already default"}
	ErrMarketCannotDelete     = &Err{HTTPCode: http.StatusBadRequest, Code: 150008, Msg: "cannot delete default market"}

	// ErrProductMarketNotFound ==================== ProductMarket Module (160xxx) ====================
	ErrProductMarketNotFound      = &Err{HTTPCode: http.StatusNotFound, Code: 160001, Msg: "product market not found"}
	ErrProductMarketAlreadyExists = &Err{HTTPCode: http.StatusConflict, Code: 160002, Msg: "product already in market"}
	ErrProductMarketPriceRequired = &Err{HTTPCode: http.StatusBadRequest, Code: 160003, Msg: "price is required for market"}
```

- [ ] **Step 2: Verify file compiles**

Run: `cd pkg/code && go build .`
Expected: No errors

- [ ] **Step 3: Commit**

```bash
git add pkg/code/code.go
git commit -m "feat(code): add error codes for Market and ProductMarket modules

- Market errors: 150xxx range
- ProductMarket errors: 160xxx range"
```

---

## Task 4: Market Repository Implementation

**Files:**
- Create: `admin/internal/infrastructure/persistence/market_repository.go`

- [ ] **Step 1: Create market repository implementation**

```go
package persistence

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/colinrs/shopjoy/admin/internal/domain/market"
	"github.com/colinrs/shopjoy/pkg/code"
	"gorm.io/gorm"
)

type marketRepo struct{}

func NewMarketRepository() market.Repository {
	return &marketRepo{}
}

type marketModel struct {
	ID              int64          `gorm:"column:id;primaryKey;autoIncrement"`
	TenantID        int64          `gorm:"column:tenant_id;not null;default:0"`
	Code            string         `gorm:"column:code;size:10;not null;uniqueIndex:uk_tenant_code"`
	Name            string         `gorm:"column:name;size:64;not null"`
	Currency        string         `gorm:"column:currency;size:10;not null"`
	DefaultLanguage string         `gorm:"column:default_language;size:10;not null;default:'en'"`
	Flag            string         `gorm:"column:flag;size:32"`
	IsActive        bool           `gorm:"column:is_active;not null;default:true"`
	IsDefault       bool           `gorm:"column:is_default;not null;default:false"`
	TaxRules        string         `gorm:"column:tax_rules;type:json"`
	CreatedAt       time.Time      `gorm:"column:created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (marketModel) TableName() string {
	return "markets"
}

func (m *marketModel) toEntity() *market.Market {
	var taxRules market.TaxConfig
	if m.TaxRules != "" {
		json.Unmarshal([]byte(m.TaxRules), &taxRules)
	}

	return &market.Market{
		ID:              m.ID,
		TenantID:        m.TenantID,
		Code:            m.Code,
		Name:            m.Name,
		Currency:        m.Currency,
		DefaultLanguage: m.DefaultLanguage,
		Flag:            m.Flag,
		IsActive:        m.IsActive,
		IsDefault:       m.IsDefault,
		TaxRules:        taxRules,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func fromMarketEntity(m *market.Market) *marketModel {
	taxRulesJSON, _ := json.Marshal(m.TaxRules)

	return &marketModel{
		ID:              m.ID,
		TenantID:        m.TenantID,
		Code:            m.Code,
		Name:            m.Name,
		Currency:        m.Currency,
		DefaultLanguage: m.DefaultLanguage,
		Flag:            m.Flag,
		IsActive:        m.IsActive,
		IsDefault:       m.IsDefault,
		TaxRules:        string(taxRulesJSON),
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
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
```

- [ ] **Step 2: Add missing import for time**

The file needs `time` import. Add to imports:

```go
import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/market"
	"github.com/colinrs/shopjoy/pkg/code"
	"gorm.io/gorm"
)
```

- [ ] **Step 3: Verify file compiles**

Run: `cd admin && go build ./internal/infrastructure/persistence/...`
Expected: No errors

- [ ] **Step 4: Commit**

```bash
git add admin/internal/infrastructure/persistence/market_repository.go
git commit -m "feat(infra): implement Market repository

- Implement all CRUD operations
- Add FindByCode, FindActive, FindDefault methods
- Handle soft delete with gorm.DeletedAt"
```

---

## Task 5: Market API Definition

**Files:**
- Create: `admin/desc/market.api`

- [ ] **Step 1: Create market API definition**

```go
syntax = "v1"

info (
	title:   "Market API"
	desc:    "市场管理相关接口"
	version: "v1"
)

type (
	// Tax configuration
	TaxConfig {
		VatRate     string `json:"vat_rate,optional"`
		GstRate     string `json:"gst_rate,optional"`
		IossEnabled bool   `json:"ioss_enabled,optional"`
		IncludeTax  bool   `json:"include_tax,optional"`
	}

	// Market responses
	MarketResponse {
		ID              int64      `json:"id"`
		Code            string     `json:"code"`
		Name            string     `json:"name"`
		Currency        string     `json:"currency"`
		DefaultLanguage string     `json:"default_language"`
		Flag            string     `json:"flag"`
		IsActive        bool       `json:"is_active"`
		IsDefault       bool       `json:"is_default"`
		TaxRules        TaxConfig  `json:"tax_rules"`
		CreatedAt       string     `json:"created_at"`
		UpdatedAt       string     `json:"updated_at"`
	}

	// Create market request
	CreateMarketReq {
		Code            string    `json:"code"`                     // US, UK, DE, FR, AU
		Name            string    `json:"name"`                     // United States
		Currency        string    `json:"currency"`                 // USD, GBP, EUR, AUD
		DefaultLanguage string    `json:"default_language,optional"`
		Flag            string    `json:"flag,optional"`
		TaxRules        TaxConfig `json:"tax_rules,optional"`
	}

	// Update market request
	UpdateMarketReq {
		ID              int64     `path:"id"`
		Name            string    `json:"name,optional"`
		IsActive        *bool     `json:"is_active,optional"`
		TaxRules        TaxConfig `json:"tax_rules,optional"`
	}

	// Get market request
	GetMarketReq {
		ID int64 `path:"id"`
	}

	// List markets response
	ListMarketsResp {
		List []*MarketResponse `json:"list"`
		Total int64             `json:"total"`
	}
)

@server (
	group:      markets
	middleware: AuthMiddleware
)
service admin-api {
	@doc "创建市场"
	@handler CreateMarketHandler
	post /api/v1/markets (CreateMarketReq) returns (MarketResponse)

	@doc "更新市场"
	@handler UpdateMarketHandler
	put /api/v1/markets/:id (UpdateMarketReq) returns (MarketResponse)

	@doc "获取市场详情"
	@handler GetMarketHandler
	get /api/v1/markets/:id (GetMarketReq) returns (MarketResponse)

	@doc "获取市场列表"
	@handler ListMarketsHandler
	get /api/v1/markets returns (ListMarketsResp)

	@doc "删除市场"
	@handler DeleteMarketHandler
	delete /api/v1/markets/:id (GetMarketReq) returns (interface{})
}
```

- [ ] **Step 2: Generate API code**

Run: `cd admin && make api`
Expected: Code generated successfully

- [ ] **Step 3: Commit**

```bash
git add admin/desc/market.api admin/internal/types/types.go admin/internal/handler/markets/ admin/internal/handler/routes.go
git commit -m "feat(api): add Market API definition and generated code

- CRUD endpoints for market management
- TaxConfig embedded in market responses"
```

---

## Task 6: Market Handlers and Logic

**Files:**
- Modify: `admin/internal/handler/markets/*.go` (auto-generated, scaffold)
- Create/Modify: `admin/internal/logic/markets/*.go`

- [ ] **Step 1: Create market logic files**

The handlers are auto-generated. Create the logic implementations.

`admin/internal/logic/markets/create_market_logic.go`:

```go
package markets

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/domain/market"
	"github.com/colinrs/shopjoy/admin/internal/infrastructure/persistence"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
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
		ID:              m.ID,
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
		CreatedAt: m.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: m.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
```

Add import for decimal:
```go
import "github.com/shopspring/decimal"
```

`admin/internal/logic/markets/list_markets_logic.go`:

```go
package markets

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/infrastructure/persistence"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListMarketsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMarketsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMarketsLogic {
	return &ListMarketsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMarketsLogic) ListMarkets() (resp *types.ListMarketsResp, err error) {
	repo := persistence.NewMarketRepository()
	markets, err := repo.FindAll(l.ctx, l.svcCtx.DB)
	if err != nil {
		return nil, err
	}

	list := make([]*types.MarketResponse, len(markets))
	for i, m := range markets {
		list[i] = toMarketResponse(m)
	}

	return &types.ListMarketsResp{
		List:  list,
		Total: int64(len(list)),
	}, nil
}
```

`admin/internal/logic/markets/get_market_logic.go`:

```go
package markets

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/infrastructure/persistence"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetMarketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMarketLogic {
	return &GetMarketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMarketLogic) GetMarket(req *types.GetMarketReq) (resp *types.MarketResponse, err error) {
	repo := persistence.NewMarketRepository()
	m, err := repo.FindByID(l.ctx, l.svcCtx.DB, req.ID)
	if err != nil {
		return nil, err
	}
	return toMarketResponse(m), nil
}
```

`admin/internal/logic/markets/update_market_logic.go`:

```go
package markets

import (
	"context"

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
```

Add import: `import "github.com/colinrs/shopjoy/admin/internal/domain/market"`

`admin/internal/logic/markets/delete_market_logic.go`:

```go
package markets

import (
	"context"

	"github.com/colinrs/shopjoy/admin/internal/infrastructure/persistence"
	"github.com/colinrs/shopjoy/admin/internal/svc"
	"github.com/colinrs/shopjoy/admin/internal/types"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMarketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMarketLogic {
	return &DeleteMarketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMarketLogic) DeleteMarket(req *types.GetMarketReq) (resp interface{}, err error) {
	repo := persistence.NewMarketRepository()
	m, err := repo.FindByID(l.ctx, l.svcCtx.DB, req.ID)
	if err != nil {
		return nil, err
	}

	// Cannot delete default market
	if m.IsDefault {
		return nil, code.ErrMarketCannotDelete
	}

	if err := repo.Delete(l.ctx, l.svcCtx.DB, req.ID); err != nil {
		return nil, err
	}

	return nil, nil
}
```

- [ ] **Step 2: Verify build**

Run: `cd admin && make build`
Expected: Build successful

- [ ] **Step 3: Commit**

```bash
git add admin/internal/logic/markets/
git commit -m "feat(logic): implement Market CRUD logic

- Create, Update, Get, Delete, List operations
- Prevent deletion of default market
- Map domain entities to API responses"
```

---

## Task 7: ProductMarket Domain Entity

**Files:**
- Create: `admin/internal/domain/product/product_market.go`

- [ ] **Step 1: Create ProductMarket entity in product domain**

```go
// Package product 商品领域层
package product

import (
	"context"
	"time"

	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ProductMarket 商品-市场关联实体
type ProductMarket struct {
	ID                 int64           // 关联ID
	TenantID           int64           // 租户ID
	ProductID          int64           // 商品ID
	VariantID          *int64          // 变体ID，NULL表示基础商品
	MarketID           int64           // 市场ID
	IsEnabled          bool            // 是否在该市场可见
	StatusOverride     *Status         // 状态覆盖
	Price              decimal.Decimal // 市场专属价格
	CompareAtPrice     *decimal.Decimal // 对比价格
	StockAlertThreshold int            // 库存预警阈值
	PublishedAt        *time.Time      // 发布时间
	CreatedAt          time.Time       // 创建时间
	UpdatedAt          time.Time       // 更新时间
}

// TableName 表名
func (pm *ProductMarket) TableName() string {
	return "product_markets"
}

// NewProductMarket 创建商品市场关联
func NewProductMarket(productID, marketID int64) (*ProductMarket, error) {
	if productID <= 0 {
		return nil, code.ErrProductInvalidID
	}
	if marketID <= 0 {
		return nil, code.ErrMarketNotFound
	}

	now := time.Now()
	return &ProductMarket{
		ProductID: productID,
		MarketID:  marketID,
		IsEnabled: false, // Default to disabled, requires price setup
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Enable 启用市场可见性
func (pm *ProductMarket) Enable() {
	pm.IsEnabled = true
	pm.UpdatedAt = time.Now()
}

// Disable 禁用市场可见性
func (pm *ProductMarket) Disable() {
	pm.IsEnabled = false
	pm.UpdatedAt = time.Now()
}

// SetPrice 设置价格
func (pm *ProductMarket) SetPrice(price decimal.Decimal) error {
	if price.LessThanOrEqual(decimal.Zero) {
		return code.ErrProductMarketPriceRequired
	}
	pm.Price = price
	pm.UpdatedAt = time.Now()
	return nil
}

// Publish 发布到市场
func (pm *ProductMarket) Publish() error {
	if pm.Price.IsZero() {
		return code.ErrProductMarketPriceRequired
	}
	pm.Enable()
	now := time.Now()
	pm.PublishedAt = &now
	return nil
}

// ProductMarketRepository 商品市场仓储接口
type ProductMarketRepository interface {
	Create(ctx context.Context, db *gorm.DB, pm *ProductMarket) error
	Update(ctx context.Context, db *gorm.DB, pm *ProductMarket) error
	Delete(ctx context.Context, db *gorm.DB, id int64) error
	FindByID(ctx context.Context, db *gorm.DB, id int64) (*ProductMarket, error)
	FindByProductAndMarket(ctx context.Context, db *gorm.DB, productID, marketID int64, variantID *int64) (*ProductMarket, error)
	FindByProductID(ctx context.Context, db *gorm.DB, productID int64) ([]*ProductMarket, error)
	FindByMarketID(ctx context.Context, db *gorm.DB, marketID int64, query Query) ([]*ProductMarket, int64, error)
	BatchCreate(ctx context.Context, db *gorm.DB, pms []*ProductMarket) error
}
```

- [ ] **Step 2: Verify file compiles**

Run: `cd admin && go build ./internal/domain/product/...`
Expected: No errors

- [ ] **Step 3: Commit**

```bash
git add admin/internal/domain/product/product_market.go
git commit -m "feat(domain): add ProductMarket entity for product-market association

- Enable/Disable market visibility
- SetPrice and Publish methods
- ProductMarketRepository interface"
```

---

## Task 8: ProductMarket Repository Implementation

**Files:**
- Create: `admin/internal/infrastructure/persistence/product_market_repository.go`

- [ ] **Step 1: Create ProductMarket repository**

```go
package persistence

import (
	"context"
	"errors"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/pkg/code"
	"gorm.io/gorm"
)

type productMarketRepo struct{}

func NewProductMarketRepository() product.ProductMarketRepository {
	return &productMarketRepo{}
}

type productMarketModel struct {
	ID                  int64          `gorm:"column:id;primaryKey;autoIncrement"`
	TenantID            int64          `gorm:"column:tenant_id;not null"`
	ProductID           int64          `gorm:"column:product_id;not null;uniqueIndex:uk_product_variant_market"`
	VariantID           *int64         `gorm:"column:variant_id;uniqueIndex:uk_product_variant_market"`
	MarketID            int64          `gorm:"column:market_id;not null;uniqueIndex:uk_product_variant_market"`
	IsEnabled           bool           `gorm:"column:is_enabled;not null;default:false"`
	StatusOverride      *int           `gorm:"column:status_override"`
	Price               string         `gorm:"column:price;type:decimal(19,4);not null"`
	CompareAtPrice      *string        `gorm:"column:compare_at_price;type:decimal(19,4)"`
	StockAlertThreshold int            `gorm:"column:stock_alert_threshold;not null;default:0"`
	PublishedAt         *time.Time     `gorm:"column:published_at"`
	CreatedAt           time.Time      `gorm:"column:created_at"`
	UpdatedAt           time.Time      `gorm:"column:updated_at"`
}

func (productMarketModel) TableName() string {
	return "product_markets"
}

func (m *productMarketModel) toEntity() *product.ProductMarket {
	pm := &product.ProductMarket{
		ID:                  m.ID,
		TenantID:            m.TenantID,
		ProductID:           m.ProductID,
		VariantID:           m.VariantID,
		MarketID:            m.MarketID,
		IsEnabled:           m.IsEnabled,
		StockAlertThreshold: m.StockAlertThreshold,
		PublishedAt:         m.PublishedAt,
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
	}

	pm.Price, _ = decimal.NewFromString(m.Price)
	if m.CompareAtPrice != nil {
		cap, _ := decimal.NewFromString(*m.CompareAtPrice)
		pm.CompareAtPrice = &cap
	}
	if m.StatusOverride != nil {
		status := product.Status(*m.StatusOverride)
		pm.StatusOverride = &status
	}

	return pm
}

func fromProductMarketEntity(pm *product.ProductMarket) *productMarketModel {
	m := &productMarketModel{
		ID:                  pm.ID,
		TenantID:            pm.TenantID,
		ProductID:           pm.ProductID,
		VariantID:           pm.VariantID,
		MarketID:            pm.MarketID,
		IsEnabled:           pm.IsEnabled,
		StockAlertThreshold: pm.StockAlertThreshold,
		PublishedAt:         pm.PublishedAt,
		CreatedAt:           pm.CreatedAt,
		UpdatedAt:           pm.UpdatedAt,
	}

	m.Price = pm.Price.String()
	if pm.CompareAtPrice != nil {
		cap := pm.CompareAtPrice.String()
		m.CompareAtPrice = &cap
	}
	if pm.StatusOverride != nil {
		status := int(*pm.StatusOverride)
		m.StatusOverride = &status
	}

	return m
}

func (r *productMarketRepo) Create(ctx context.Context, db *gorm.DB, pm *product.ProductMarket) error {
	model := fromProductMarketEntity(pm)
	return db.WithContext(ctx).Create(model).Error
}

func (r *productMarketRepo) Update(ctx context.Context, db *gorm.DB, pm *product.ProductMarket) error {
	model := fromProductMarketEntity(pm)
	return db.WithContext(ctx).Save(model).Error
}

func (r *productMarketRepo) Delete(ctx context.Context, db *gorm.DB, id int64) error {
	result := db.WithContext(ctx).Delete(&productMarketModel{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return code.ErrProductMarketNotFound
	}
	return nil
}

func (r *productMarketRepo) FindByID(ctx context.Context, db *gorm.DB, id int64) (*product.ProductMarket, error) {
	var model productMarketModel
	err := db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrProductMarketNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *productMarketRepo) FindByProductAndMarket(ctx context.Context, db *gorm.DB, productID, marketID int64, variantID *int64) (*product.ProductMarket, error) {
	var model productMarketModel
	query := db.WithContext(ctx).Where("product_id = ? AND market_id = ?", productID, marketID)
	if variantID != nil {
		query = query.Where("variant_id = ?", *variantID)
	} else {
		query = query.Where("variant_id IS NULL")
	}

	err := query.First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, code.ErrProductMarketNotFound
		}
		return nil, err
	}
	return model.toEntity(), nil
}

func (r *productMarketRepo) FindByProductID(ctx context.Context, db *gorm.DB, productID int64) ([]*product.ProductMarket, error) {
	var models []productMarketModel
	err := db.WithContext(ctx).Where("product_id = ?", productID).Find(&models).Error
	if err != nil {
		return nil, err
	}

	pms := make([]*product.ProductMarket, len(models))
	for i, m := range models {
		pms[i] = m.toEntity()
	}
	return pms, nil
}

func (r *productMarketRepo) FindByMarketID(ctx context.Context, db *gorm.DB, marketID int64, query product.Query) ([]*product.ProductMarket, int64, error) {
	dbQuery := db.WithContext(ctx).Model(&productMarketModel{}).Where("market_id = ?", marketID)

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []productMarketModel
	err := dbQuery.Order("created_at DESC").
		Offset(query.Offset()).
		Limit(query.Limit()).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	pms := make([]*product.ProductMarket, len(models))
	for i, m := range models {
		pms[i] = m.toEntity()
	}
	return pms, total, nil
}

func (r *productMarketRepo) BatchCreate(ctx context.Context, db *gorm.DB, pms []*product.ProductMarket) error {
	models := make([]*productMarketModel, len(pms))
	for i, pm := range pms {
		models[i] = fromProductMarketEntity(pm)
	}
	return db.WithContext(ctx).Create(&models).Error
}
```

Add imports:
```go
import (
	"context"
	"errors"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)
```

- [ ] **Step 2: Verify build**

Run: `cd admin && make build`
Expected: Build successful

- [ ] **Step 3: Commit**

```bash
git add admin/internal/infrastructure/persistence/product_market_repository.go
git commit -m "feat(infra): implement ProductMarket repository

- CRUD operations for product-market associations
- FindByProductID, FindByMarketID queries
- BatchCreate for bulk operations"
```

---

## Task 9: Update Product Entity with Compliance Fields

**Files:**
- Modify: `admin/internal/domain/product/entity.go`
- Modify: `admin/internal/infrastructure/persistence/product_repository.go`

- [ ] **Step 1: Add compliance fields to Product entity**

Add Dimensions struct after Money struct (around line 100):

```go
// Dimensions 尺寸（值对象）
type Dimensions struct {
	Length decimal.Decimal
	Width  decimal.Decimal
	Height decimal.Decimal
	Unit   string // cm
}
```

Update Product struct to add new fields. Replace the entire Product struct:

```go
// Product 商品实体
type Product struct {
	ID              int64          // 商品ID
	SKU             string         // SKU代码
	Name            string         // 商品名称
	Description     string         // 商品描述
	Price           Money          `gorm:"embedded"` // 售价
	CostPrice       Money          `gorm:"embedded"` // 成本价
	Stock           int            // 库存
	Status          Status         // 状态
	CategoryID      int64          // 分类ID
	Brand           string         // 品牌
	Tags            []string       `gorm:"type:json"` // 标签
	Images          []string       `gorm:"type:json"` // 图片列表
	IsMatrixProduct bool           // 是否有变体

	// Compliance fields (cross-border)
	HSCode         string          // HS编码
	COO            string          // 原产国
	Weight         decimal.Decimal // 重量
	WeightUnit     string          // 重量单位
	Dimensions     Dimensions      `gorm:"embedded"` // 尺寸
	DangerousGoods []string        `gorm:"type:json"` // 危险品标识

	CreatedAt       time.Time      // 创建时间
	UpdatedAt       time.Time      // 更新时间
}
```

Add new constructor and methods after the existing NewProduct function:

```go
// NewProductWithCompliance 创建带合规信息的商品
func NewProductWithCompliance(id int64, name, description, sku string, price Money, categoryID int64) (*Product, error) {
	if name == "" {
		return nil, code.ErrProductEmptyName
	}
	if price.Amount <= 0 {
		return nil, code.ErrProductInvalidPrice
	}
	if id <= 0 {
		return nil, code.ErrProductInvalidID
	}

	now := time.Now()
	return &Product{
		ID:          id,
		SKU:         sku,
		Name:        name,
		Description: description,
		Price:       price,
		Status:      StatusDraft,
		CategoryID:  categoryID,
		WeightUnit:  "g",
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// SetCompliance 设置合规信息
func (p *Product) SetCompliance(hsCode, coo string, weight decimal.Decimal, weightUnit string, dims Dimensions) {
	p.HSCode = hsCode
	p.COO = coo
	p.Weight = weight
	p.WeightUnit = weightUnit
	p.Dimensions = dims
	p.UpdatedAt = time.Now()
}

// HasComplianceInfo 检查是否有合规信息
func (p *Product) HasComplianceInfo() bool {
	return p.HSCode != "" && p.COO != "" && !p.Weight.IsZero()
}

// IsDangerousGoods 检查是否为危险品
func (p *Product) IsDangerousGoods() bool {
	return len(p.DangerousGoods) > 0
}
```

Add import for decimal:
```go
import "github.com/shopspring/decimal"
```

- [ ] **Step 2: Update productModel and conversion methods in repository**

Replace the entire `productModel` struct and conversion methods in `admin/internal/infrastructure/persistence/product_repository.go`:

```go
type productModel struct {
	ID              int64  `gorm:"column:id;primaryKey;autoIncrement:false"`
	SKU             string `gorm:"column:sku;size:64;uniqueIndex"`
	Name            string `gorm:"column:name;size:200;not null;index"`
	Description     string `gorm:"column:description;type:text"`
	Price           int64  `gorm:"column:price;not null"`
	CostPrice       int64  `gorm:"column:cost_price"`
	Currency        string `gorm:"column:currency;size:10;default:'CNY'"`
	Stock           int    `gorm:"column:stock;default:0"`
	Status          int    `gorm:"column:status;default:0;index"`
	CategoryID      int64  `gorm:"column:category_id;index"`
	Brand           string `gorm:"column:brand;size:64"`
	Tags            string `gorm:"column:tags;type:json"`
	Images          string `gorm:"column:images;type:json"`
	IsMatrixProduct bool   `gorm:"column:is_matrix_product;default:false"`
	HSCode          string `gorm:"column:hs_code;size:20"`
	COO             string `gorm:"column:coo;size:10"`
	Weight          string `gorm:"column:weight;type:decimal(10,2)"`
	WeightUnit      string `gorm:"column:weight_unit;size:10;default:'g'"`
	Length          string `gorm:"column:length;type:decimal(10,2)"`
	Width           string `gorm:"column:width;type:decimal(10,2)"`
	Height          string `gorm:"column:height;type:decimal(10,2)"`
	DangerousGoods  string `gorm:"column:dangerous_goods;type:json"`
	CreatedAt       int64  `gorm:"column:created_at"`
	UpdatedAt       int64  `gorm:"column:updated_at"`
}

func (productModel) TableName() string {
	return "products"
}

func (m *productModel) toEntity() *product.Product {
	// Parse JSON arrays
	var tags, images, dangerousGoods []string
	if m.Tags != "" {
		json.Unmarshal([]byte(m.Tags), &tags)
	}
	if m.Images != "" {
		json.Unmarshal([]byte(m.Images), &images)
	}
	if m.DangerousGoods != "" {
		json.Unmarshal([]byte(m.DangerousGoods), &dangerousGoods)
	}

	// Parse decimals
	weight, _ := decimal.NewFromString(m.Weight)
	length, _ := decimal.NewFromString(m.Length)
	width, _ := decimal.NewFromString(m.Width)
	height, _ := decimal.NewFromString(m.Height)

	return &product.Product{
		ID:              m.ID,
		SKU:             m.SKU,
		Name:            m.Name,
		Description:     m.Description,
		Price:           product.NewMoney(m.Price, m.Currency),
		CostPrice:       product.NewMoney(m.CostPrice, m.Currency),
		Stock:           m.Stock,
		Status:          product.Status(m.Status),
		CategoryID:      m.CategoryID,
		Brand:           m.Brand,
		Tags:            tags,
		Images:          images,
		IsMatrixProduct: m.IsMatrixProduct,
		HSCode:          m.HSCode,
		COO:             m.COO,
		Weight:          weight,
		WeightUnit:      m.WeightUnit,
		Dimensions: product.Dimensions{
			Length: length,
			Width:  width,
			Height: height,
			Unit:   "cm",
		},
		DangerousGoods: dangerousGoods,
		CreatedAt:      time.Unix(m.CreatedAt, 0),
		UpdatedAt:      time.Unix(m.UpdatedAt, 0),
	}
}

func fromEntity(p *product.Product) *productModel {
	// Serialize JSON arrays
	tagsJSON, _ := json.Marshal(p.Tags)
	imagesJSON, _ := json.Marshal(p.Images)
	dangerousGoodsJSON, _ := json.Marshal(p.DangerousGoods)

	return &productModel{
		ID:              p.ID,
		SKU:             p.SKU,
		Name:            p.Name,
		Description:     p.Description,
		Price:           p.Price.Amount,
		CostPrice:       p.CostPrice.Amount,
		Currency:        p.Price.Currency,
		Stock:           p.Stock,
		Status:          int(p.Status),
		CategoryID:      p.CategoryID,
		Brand:           p.Brand,
		Tags:            string(tagsJSON),
		Images:          string(imagesJSON),
		IsMatrixProduct: p.IsMatrixProduct,
		HSCode:          p.HSCode,
		COO:             p.COO,
		Weight:          p.Weight.String(),
		WeightUnit:      p.WeightUnit,
		Length:          p.Dimensions.Length.String(),
		Width:           p.Dimensions.Width.String(),
		Height:          p.Dimensions.Height.String(),
		DangerousGoods:  string(dangerousGoodsJSON),
		CreatedAt:       p.CreatedAt.Unix(),
		UpdatedAt:       p.UpdatedAt.Unix(),
	}
}
```

Add imports at the top of the file:
```go
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/colinrs/shopjoy/admin/internal/domain/product"
	"github.com/colinrs/shopjoy/pkg/code"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)
```

Update the Update method to include new fields. Replace the Update function:

```go
func (r *productRepo) Update(ctx context.Context, db *gorm.DB, p *product.Product) error {
	model := fromEntity(p)
	return db.WithContext(ctx).
		Where("id = ? AND status != ?", p.ID, product.StatusDeleted).
		Updates(map[string]interface{}{
			"sku":               model.SKU,
			"name":              model.Name,
			"description":       model.Description,
			"price":             model.Price,
			"cost_price":        model.CostPrice,
			"currency":          model.Currency,
			"stock":             model.Stock,
			"status":            model.Status,
			"category_id":       model.CategoryID,
			"brand":             model.Brand,
			"tags":              model.Tags,
			"images":            model.Images,
			"is_matrix_product": model.IsMatrixProduct,
			"hs_code":           model.HSCode,
			"coo":               model.COO,
			"weight":            model.Weight,
			"weight_unit":       model.WeightUnit,
			"length":            model.Length,
			"width":             model.Width,
			"height":            model.Height,
			"dangerous_goods":   model.DangerousGoods,
			"updated_at":        model.UpdatedAt,
		}).Error
}
```

- [ ] **Step 3: Verify build**

Run: `cd admin && make build`
Expected: Build successful

- [ ] **Step 4: Commit**

```bash
git add admin/internal/domain/product/entity.go admin/internal/infrastructure/persistence/product_repository.go
git commit -m "feat(product): add compliance fields to Product entity

- SKU, Brand, Tags, Images
- Compliance: HSCode, COO, Weight, Dimensions
- DangerousGoods array
- HasComplianceInfo method for validation
- Complete repository mapping for all new fields"
```

---

## Task 10: Final Build Verification

- [ ] **Step 1: Run full build**

Run: `cd admin && make build`
Expected: Build successful, outputs `bin/admin`

- [ ] **Step 2: Run tests (if any)**

Run: `cd admin && go test ./...`
Expected: All tests pass

- [ ] **Step 3: Final commit summary**

```bash
git status
git log --oneline -10
```

---

## Summary

Phase 1 establishes the core infrastructure for multi-market product management:

1. **Database** - `markets`, `product_markets` tables + product compliance fields
2. **Market Entity** - Domain model with tax configuration
3. **Market API** - CRUD endpoints for market management
4. **ProductMarket Entity** - Product-market association with pricing
5. **Product Updates** - Compliance fields (HS Code, COO, dimensions)

**Next Phase:** Product Management UI (list page with market filter, detail page tabs, Push to Market flow)