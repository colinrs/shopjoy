# 运费配置国际化改造 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将运费配置从中国单一市场改造为支持 US/EU/SEA/JP+KR 多市场的国际化方案，涵盖多币种、多语言、多税费模型、多物流商规则引擎。

**Architecture:**
- 后端在 Go 领域层增加 `Market`、`Currency`、`Warehouse`、`Carrier` 等核心概念；仓储层 schema 演进式迁移；`fee_engine` 抽象通用规则引擎承载不同物流商计费模型
- 前端 TypeScript 类型与后端同步重写；`useCurrency` composable 统一币种展示；ZoneConfigForm 增加国家选择与体积计费配置
- 数据库采用 additive 迁移 + 主 schema.sql 同步合并模式
- API 允许 breaking change，前端一次性同步升级

**Tech Stack:**
- 后端: Go 1.22+, go-zero, GORM, decimal.Decimal, pkg/code, application.Model
- 前端: Vue 3, TypeScript, Element Plus, Pinia, vue-i18n
- 数据库: MySQL 8, JSONB (regions, name_i18n, tax_rules)
- 区域数据: ISO 3166-1/2 国家与行政区划标准

## Global Constraints

- 所有金额字段：DB 用 `DECIMAL(19,4)`，API 层 `string`，domain 用 `decimal.Decimal`
- 所有时间字段：`time.Time`，DB 用 `TIMESTAMP`，统一 UTC
- 所有错误：使用 `pkg/code` 中已定义的 `ErrXxx`，禁止 `errors.New()`
- 所有 Model 实体：嵌入 `application.Model`（含 ID/CreatedAt/UpdatedAt/DeletedAt）
- 所有 DB 表名：小写+下划线，所有列 NOT NULL 或显式 default
- 所有前端 TypeScript：与后端枚举 1:1 严格对齐，使用字面量联合类型
- 所有前端 API 调用：`useErrorHandler` composable 统一处理错误
- 所有币种：ISO 4217 三字母代码（USD/EUR/SGD/JPY/KRW/CNY）
- 所有国家：ISO 3166-1 alpha-2 代码（US/DE/SG/JP/KR/...）
- 所有迁移：文件名 `{YYYYMMDD}{seq}_{action}_{object}.sql`，同时合并到 `sql/fulfillment/schema.sql`
- 每个 Task 末尾 `git commit`，commit message 遵循 Conventional Commits

---

## 文件总览

### 新增文件

```
admin/internal/domain/shipping/
├── fee_engine.go                    # 通用物流商规则引擎
├── weight_unit.go                   # 重量单位转换
├── zone_matcher.go                  # 多级区域匹配
├── tax_calculator.go                # 税费计算
└── carrier.go                       # 物流商抽象接口

admin/internal/domain/warehouse/
├── entity.go                        # 仓库实体
└── repository.go                    # 仓库仓储接口

admin/internal/infrastructure/persistence/
├── warehouse_repository.go          # 仓库仓储实现
└── region_repository.go             # 区域仓储实现

admin/internal/application/shipping/
└── i18n_name_resolver.go            # 国际化名称解析

admin/internal/logic/
├── warehouses/                      # 仓库 CRUD
└── carrier_configs/                 # 物流商配置 CRUD

admin/internal/handler/             # 由 go-zero 自动生成
├── warehouses/
└── carrier_configs/

admin/desc/
├── shipping.api                     # 重写：所有 type
└── warehouses.api                   # 新增：仓库 API

sql/fulfillment/migrations/
├── 2026072001_alter_templates_market.sql
├── 2026072002_alter_zones_currency.sql
├── 2026072003_alter_zones_dimensions.sql
├── 2026072004_alter_zones_tax.sql
├── 2026072005_alter_zones_surcharge.sql
├── 2026072006_alter_zones_name_i18n.sql
├── 2026072007_alter_regions_country.sql
├── 2026072008_create_carriers.sql
├── 2026072009_create_warehouses.sql
└── 2026072010_seed_regions_overseas.sql

shop-admin/src/api/
├── warehouse.ts                     # 新增
└── carrier.ts                       # 新增

shop-admin/src/composables/
└── useCurrency.ts                   # 新增：币种格式化

shop-admin/src/views/shipping/components/
├── RegionSelector.vue               # 重写：支持国家树
├── ZoneConfigForm.vue               # 重写：增加维度/税/附加费
└── CarrierRuleEditor.vue            # 新增：物流商规则编辑
```

### 修改文件

```
admin/internal/domain/shipping/entity.go         # 增加字段：Currency/MarketID/WarehouseID 等
admin/internal/domain/shipping/repository.go     # 增加方法：ListByMarket/FindByCityAndCountry 等
admin/internal/infrastructure/persistence/shipping_repository.go
admin/internal/application/shipping/service.go
admin/internal/logic/shipping_calculator/calculate_shipping_fee_logic.go
admin/internal/logic/shipping_zones/*.go
admin/internal/logic/shipping_templates/*.go

shop-admin/src/api/shipping.ts                   # 重写：所有 type
shop-admin/src/views/shipping/index.vue
shop-admin/src/views/shipping/[id]/index.vue
shop-admin/src/views/shipping/calculator/index.vue

sql/fulfillment/schema.sql                       # 合并所有 migrations
```

---

## Phase 1: Foundation — Currency 字段 + Market 隔离 (P0-1, P0-5)

> **目标：** 解决最严重的硬编码 `Currency: "CNY"` 问题，引入 Market 隔离，让一个租户可以为不同市场配置独立模板与币种。

### Task 1.1: 扩展 ShippingTemplate 实体（加 MarketID + Currency）

**Files:**
- Modify: `admin/internal/domain/shipping/entity.go:64-70`

**Interfaces:**
- Consumes: 无
- Produces: `ShippingTemplate{MarketID int64, Currency string}` 字段

- [ ] **Step 1: 写失败测试**

```go
// admin/internal/domain/shipping/entity_test.go
package shipping

import "testing"

func TestShippingTemplate_MarketIsolation(t *testing.T) {
    tmpl := &ShippingTemplate{
        TenantID:  1,
        MarketID:  2, // US market
        Currency:  "USD",
        Name:      "US Standard",
        IsDefault: true,
        IsActive:  true,
    }
    if tmpl.MarketID != 2 {
        t.Errorf("expected MarketID=2, got %d", tmpl.MarketID)
    }
    if tmpl.Currency != "USD" {
        t.Errorf("expected Currency=USD, got %s", tmpl.Currency)
    }
}
```

- [ ] **Step 2: 运行测试，确认失败**

```bash
cd admin && go test ./internal/domain/shipping/... -run TestShippingTemplate_MarketIsolation -v
```

Expected: FAIL — `MarketID` 字段不存在。

- [ ] **Step 3: 修改实体**

```go
// admin/internal/domain/shipping/entity.go
// ShippingTemplate 运费模板实体
type ShippingTemplate struct {
    application.Model
    TenantID  int64  `gorm:"column:tenant_id;not null;index"`
    MarketID  int64  `gorm:"column:market_id;not null;default:0;index:idx_market_default"` // 0=全市场通用
    Currency  string `gorm:"column:currency;size:3;not null;default:'CNY'"`
    Name      string `gorm:"column:name;size:100;not null"`
    IsDefault bool   `gorm:"column:is_default;not null;default:false;index:idx_market_default"`
    IsActive  bool   `gorm:"column:is_active;not null;default:true"`
    CarrierCode string `gorm:"column:carrier_code;size:50;not null;default:'standard'"` // P2 预留
    WarehouseID  int64  `gorm:"column:warehouse_id;not null;default:0;index"`             // P2 预留
}
```

- [ ] **Step 4: 运行测试，确认通过**

```bash
cd admin && go test ./internal/domain/shipping/... -run TestShippingTemplate_MarketIsolation -v
```

Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add admin/internal/domain/shipping/entity.go admin/internal/domain/shipping/entity_test.go
git commit -m "feat(shipping): add MarketID and Currency fields to ShippingTemplate"
```

---

### Task 1.2: 扩展 ShippingZone 实体（加 Currency + Taxable 字段）

**Files:**
- Modify: `admin/internal/domain/shipping/entity.go:106-120`

- [ ] **Step 1: 写失败测试**

```go
func TestShippingZone_Taxable(t *testing.T) {
    zone := &ShippingZone{
        TenantID: 1,
        Currency: "EUR",
        FeeType:  FeeTypeByWeight,
        Taxable:  true,
    }
    if !zone.Taxable {
        t.Error("expected Taxable=true")
    }
    if zone.Currency != "EUR" {
        t.Errorf("expected Currency=EUR, got %s", zone.Currency)
    }
}
```

- [ ] **Step 2: 运行测试，确认失败**

```bash
cd admin && go test ./internal/domain/shipping/... -run TestShippingZone_Taxable -v
```

- [ ] **Step 3: 修改实体**

```go
// ShippingZone 配送区域实体
type ShippingZone struct {
    application.Model
    TenantID            int64           `gorm:"column:tenant_id;not null;index"`
    TemplateID          int64           `gorm:"column:template_id;not null;index"`
    MarketID            int64           `gorm:"column:market_id;not null;default:0;index"` // 冗余便于查询
    Currency            string          `gorm:"column:currency;size:3;not null;default:'CNY'"`
    Name                string          `gorm:"column:name;size:100;not null"`
    NameI18n            StringI18n      `gorm:"column:name_i18n;type:jsonb"`              // P1-10
    Regions             Regions         `gorm:"column:regions;type:jsonb;not null"`
    FeeType             FeeType         `gorm:"column:fee_type;size:20;not null"`
    FirstUnit           int             `gorm:"column:first_unit;not null;default:1"`
    FirstFee            decimal.Decimal `gorm:"column:first_fee;type:decimal(19,4);not null;default:0"`
    AdditionalUnit      int             `gorm:"column:additional_unit;not null;default:1"`
    AdditionalFee       decimal.Decimal `gorm:"column:additional_fee;type:decimal(19,4);not null;default:0"`
    FreeThresholdAmount decimal.Decimal `gorm:"column:free_threshold_amount;type:decimal(19,4);not null;default:0"`
    FreeThresholdCount  int             `gorm:"column:free_threshold_count;not null;default:0"`
    Taxable             bool            `gorm:"column:taxable;not null;default:false"`     // P1-6
    TaxRate             decimal.Decimal `gorm:"column:tax_rate;type:decimal(5,4);not null;default:0"` // P1-6
    TaxIncluded         bool            `gorm:"column:tax_included;not null;default:false"` // P1-6
    IossApplicable      bool            `gorm:"column:ioss_applicable;not null;default:false"` // P1-6
    // P1-7 偏远地区附加费
    RemoteSurcharge     decimal.Decimal `gorm:"column:remote_surcharge;type:decimal(19,4);not null;default:0"`
    RemoteZipPatterns   StringArray     `gorm:"column:remote_zip_patterns;type:jsonb"`
    // P1-8 燃油附加费
    FuelSurchargePct    decimal.Decimal `gorm:"column:fuel_surcharge_pct;type:decimal(5,4);not null;default:0"`
    // P1-9 体积计费
    VolumetricDivisor   int             `gorm:"column:volumetric_divisor;not null;default:5000"` // cm³/kg
    Sort                int             `gorm:"column:sort;not null;default:0"`
}
```

> **说明：** 一次性引入所有 P1 字段（task 1.2 同时为后续 Phase 做准备），新增的 `StringI18n`、`StringArray` 类型见 Task 1.4。

- [ ] **Step 4: 运行测试，确认通过**

```bash
cd admin && go test ./internal/domain/shipping/... -v
```

- [ ] **Step 5: 提交**

```bash
git add admin/internal/domain/shipping/entity.go admin/internal/domain/shipping/entity_test.go
git commit -m "feat(shipping): add Currency/Taxable/i18n/surcharge fields to ShippingZone"
```

---

### Task 1.3: 新增 StringI18n / StringArray 值对象

**Files:**
- Create: `admin/internal/domain/shipping/value_objects.go`

- [ ] **Step 1: 写失败测试**

```go
// admin/internal/domain/shipping/value_objects_test.go
package shipping

import (
    "encoding/json"
    "testing"
)

func TestStringI18n_Marshal(t *testing.T) {
    i18n := StringI18n{"zh-CN": "华东", "en-US": "East China", "ja-JP": "華東"}
    data, err := json.Marshal(i18n)
    if err != nil {
        t.Fatal(err)
    }
    var got StringI18n
    if err := json.Unmarshal(data, &got); err != nil {
        t.Fatal(err)
    }
    if got["en-US"] != "East China" {
        t.Errorf("expected East China, got %s", got["en-US"])
    }
}

func TestStringArray_Contains(t *testing.T) {
    arr := StringArray{"110000", "310000", "US-CA"}
    if !arr.Contains("110000") {
        t.Error("expected contains 110000")
    }
    if arr.Contains("999999") {
        t.Error("expected not contains 999999")
    }
}
```

- [ ] **Step 2: 运行测试，确认失败**

```bash
cd admin && go test ./internal/domain/shipping/... -run "TestStringI18n|TestStringArray" -v
```

- [ ] **Step 3: 实现值对象**

```go
// admin/internal/domain/shipping/value_objects.go
package shipping

import (
    "database/sql/driver"
    "encoding/json"
    "errors"
    "fmt"
)

// StringI18n 多语言字符串映射（如 {"en-US": "Standard", "ja-JP": "通常"}）
type StringI18n map[string]string

// Value 实现 driver.Valuer
func (s StringI18n) Value() (driver.Value, error) {
    if s == nil {
        return nil, nil
    }
    return json.Marshal(s)
}

// Scan 实现 sql.Scanner
func (s *StringI18n) Scan(src any) error {
    if src == nil {
        *s = nil
        return nil
    }
    var data []byte
    switch v := src.(type) {
    case []byte:
        data = v
    case string:
        data = []byte(v)
    default:
        return fmt.Errorf("StringI18n: cannot scan %T", src)
    }
    return json.Unmarshal(data, s)
}

// Get 按 locale 取值，回退到 fallback，再回退到第一个非空值
func (s StringI18n) Get(locale, fallback string) string {
    if v, ok := s[locale]; ok && v != "" {
        return v
    }
    if v, ok := s[fallback]; ok && v != "" {
        return v
    }
    for _, v := range s {
        if v != "" {
            return v
        }
    }
    return ""
}

// StringArray 字符串数组
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
    if a == nil {
        return nil, nil
    }
    return json.Marshal(a)
}

func (a *StringArray) Scan(src any) error {
    if src == nil {
        *a = nil
        return nil
    }
    var data []byte
    switch v := src.(type) {
    case []byte:
        data = v
    case string:
        data = []byte(v)
    default:
        return errors.New("StringArray: unsupported scan type")
    }
    return json.Unmarshal(data, a)
}

func (a StringArray) Contains(s string) bool {
    for _, v := range a {
        if v == s {
            return true
        }
    }
    return false
}
```

- [ ] **Step 4: 运行测试，确认通过**

```bash
cd admin && go test ./internal/domain/shipping/... -run "TestStringI18n|TestStringArray" -v
```

- [ ] **Step 5: 提交**

```bash
git add admin/internal/domain/shipping/value_objects.go admin/internal/domain/shipping/value_objects_test.go
git commit -m "feat(shipping): add StringI18n and StringArray value objects"
```

---

### Task 1.4: 编写数据库迁移（合并 P0 + P1 全部字段）

**Files:**
- Create: `sql/fulfillment/migrations/2026072001_alter_shipping_intl.sql`

- [ ] **Step 1: 编写迁移 SQL**

```sql
-- sql/fulfillment/migrations/2026072001_alter_shipping_intl.sql
-- 一次性为 shipping_templates 和 shipping_zones 增加国际化所需全部字段

ALTER TABLE shipping_templates
    ADD COLUMN market_id      BIGINT       NOT NULL DEFAULT 0 COMMENT '市场ID，0=全市场通用',
    ADD COLUMN currency       VARCHAR(3)   NOT NULL DEFAULT 'CNY' COMMENT 'ISO 4217',
    ADD COLUMN carrier_code   VARCHAR(50)  NOT NULL DEFAULT 'standard' COMMENT '物流商代码',
    ADD COLUMN warehouse_id   BIGINT       NOT NULL DEFAULT 0 COMMENT '发货仓库ID',
    ADD INDEX idx_market_default (market_id, is_default);

ALTER TABLE shipping_zones
    ADD COLUMN market_id              BIGINT       NOT NULL DEFAULT 0 COMMENT '市场ID',
    ADD COLUMN currency               VARCHAR(3)   NOT NULL DEFAULT 'CNY',
    ADD COLUMN name_i18n              JSON         NULL COMMENT '多语言名称',
    ADD COLUMN taxable                TINYINT(1)   NOT NULL DEFAULT 0 COMMENT '是否计税',
    ADD COLUMN tax_rate               DECIMAL(5,4) NOT NULL DEFAULT 0 COMMENT '税率 0.0000-1.0000',
    ADD COLUMN tax_included           TINYINT(1)   NOT NULL DEFAULT 0 COMMENT '价格含税',
    ADD COLUMN ioss_applicable        TINYINT(1)   NOT NULL DEFAULT 0 COMMENT '欧盟IOSS申报',
    ADD COLUMN remote_surcharge       DECIMAL(19,4) NOT NULL DEFAULT 0 COMMENT '偏远地区附加费',
    ADD COLUMN remote_zip_patterns    JSON         NULL COMMENT '偏远邮编正则',
    ADD COLUMN fuel_surcharge_pct     DECIMAL(5,4) NOT NULL DEFAULT 0 COMMENT '燃油附加费%',
    ADD COLUMN volumetric_divisor     INT          NOT NULL DEFAULT 5000 COMMENT '体积重除数 cm³/kg',
    ADD INDEX idx_zone_market (market_id);
```

- [ ] **Step 2: 在本地 DB 应用迁移并验证**

```bash
mysql -h 192.168.0.100 -u root shopjoy < sql/fulfillment/migrations/2026072001_alter_shipping_intl.sql
mysql -h 192.168.0.100 -u root shopjoy -e "SHOW FULL COLUMNS FROM shipping_zones WHERE Field IN ('currency','taxable','name_i18n')"
```

Expected: 三行输出，列注释正确。

- [ ] **Step 3: 同步合并到 schema.sql**

将上述 ALTER 追加到 `sql/fulfillment/schema.sql` 的 `shipping_zones` 与 `shipping_templates` 表定义末尾。

- [ ] **Step 4: 提交**

```bash
git add sql/fulfillment/migrations/2026072001_alter_shipping_intl.sql sql/fulfillment/schema.sql
git commit -m "feat(shipping): add internationalization columns (currency/market/tax/i18n/surcharge)"
```

---

### Task 1.5: 仓储层增加 MarketID 过滤方法

**Files:**
- Modify: `admin/internal/domain/shipping/repository.go`
- Modify: `admin/internal/infrastructure/persistence/shipping_repository.go`

- [ ] **Step 1: 写失败测试**

```go
// admin/internal/infrastructure/persistence/shipping_repository_test.go (使用 testfixtures 或 sqlite)
package persistence

import (
    "context"
    "testing"

    "github.com/colinrs/shopjoy/admin/internal/domain/shipping"
)

func TestShippingRepo_FindDefaultByMarket(t *testing.T) {
    // 假设已 seed: tmpl-1 (market=0, is_default=true), tmpl-2 (market=2, is_default=true)
    repo := NewShippingTemplateRepository()

    got, err := repo.FindDefaultByMarket(context.Background(), testDB, 2)
    if err != nil {
        t.Fatal(err)
    }
    if got.ID != tmpl2ID {
        t.Errorf("expected tmpl-2, got %d", got.ID)
    }
}
```

- [ ] **Step 2: 在仓储接口定义新方法**

```go
// admin/internal/domain/shipping/repository.go
type ShippingTemplateRepository interface {
    // ... 现有方法
    FindDefaultByMarket(ctx context.Context, db *gorm.DB, marketID int64) (*ShippingTemplate, error)
    FindListByMarket(ctx context.Context, db *gorm.DB, marketID int64, name string, isActive *bool, page, pageSize int) ([]*ShippingTemplate, int64, error)
}
```

- [ ] **Step 3: 实现新方法**

```go
// FindDefaultByMarket 查找指定市场的默认模板；找不到则回退到 market_id=0 的全市场默认
func (r *shippingTemplateRepo) FindDefaultByMarket(ctx context.Context, db *gorm.DB, marketID int64) (*shipping.ShippingTemplate, error) {
    var tmpl shipping.ShippingTemplate
    // 1. 优先匹配 marketID
    err := db.WithContext(ctx).
        Where("market_id = ? AND is_default = ? AND is_active = ?", marketID, true, true).
        First(&tmpl).Error
    if err == nil {
        return &tmpl, nil
    }
    if !errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, err
    }
    // 2. 回退到全市场默认
    err = db.WithContext(ctx).
        Where("market_id = ? AND is_default = ? AND is_active = ?", 0, true, true).
        First(&tmpl).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, code.ErrShippingTemplateNotFound
        }
        return nil, err
    }
    return &tmpl, nil
}

// FindListByMarket 按 market 过滤列表（marketID=0 返回全部）
func (r *shippingTemplateRepo) FindListByMarket(ctx context.Context, db *gorm.DB, marketID int64, name string, isActive *bool, page, pageSize int) ([]*shipping.ShippingTemplate, int64, error) {
    var templates []*shipping.ShippingTemplate
    var total int64
    q := db.WithContext(ctx).Model(&shipping.ShippingTemplate{})
    if marketID > 0 {
        q = q.Where("market_id = ?", marketID)
    }
    if name != "" {
        q = q.Where("name LIKE ?", "%"+name+"%")
    }
    if isActive != nil {
        q = q.Where("is_active = ?", *isActive)
    }
    if err := q.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    err := q.Order("is_default DESC, id DESC").
        Offset((page - 1) * pageSize).
        Limit(pageSize).
        Find(&templates).Error
    return templates, total, err
}
```

> UnsetAllDefault 同步改造：增加 `market_id` 过滤，只清同一市场的默认。

- [ ] **Step 4: 运行测试，确认通过**

```bash
cd admin && go test ./internal/infrastructure/persistence/... -run TestShippingRepo -v
```

- [ ] **Step 5: 提交**

```bash
git add admin/internal/domain/shipping/repository.go admin/internal/infrastructure/persistence/shipping_repository.go admin/internal/infrastructure/persistence/shipping_repository_test.go
git commit -m "feat(shipping): add market-aware template lookup with fallback"
```

---

### Task 1.6: 修改 API 定义文件（重写所有 shipping 相关 type）

**Files:**
- Modify: `admin/desc/shipping.api`

- [ ] **Step 1: 在 `.api` 文件中重写所有结构**

将以下完整 `.api` 片段替换 `admin/desc/shipping.api` 中对应 section（保持路由与 handler 名称不变）：

```go
// ============================================================
// Shipping Templates
// ============================================================

type ShippingTemplate {
    Id           int64  `json:"id,string"`
    TenantID     int64  `json:"tenant_id,string"`
    MarketID     int64  `json:"market_id,string"`     // 0=全市场通用
    Currency     string `json:"currency"`             // ISO 4217
    CarrierCode  string `json:"carrier_code"`         // 物流商代码
    WarehouseID  int64  `json:"warehouse_id,string"`  // 仓库ID
    Name         string `json:"name"`
    IsDefault    bool   `json:"is_default"`
    IsActive     bool   `json:"is_active"`
    ZoneCount    int    `json:"zone_count"`
    CreatedAt    string `json:"created_at"`           // RFC3339
}

type CreateShippingTemplateReq {
    MarketID    int64  `json:"market_id,string,optional"` // 默认 0
    Currency    string `json:"currency,optional"`         // 默认 CNY
    CarrierCode string `json:"carrier_code,optional"`     // 默认 standard
    WarehouseID int64  `json:"warehouse_id,string,optional"`
    Name        string `json:"name"`
    IsDefault   bool   `json:"is_default,optional"`
}

type UpdateShippingTemplateReq {
    Id          int64  `path:"id"`
    Name        string `json:"name,optional"`
    IsActive    *bool  `json:"is_active,optional"`
    CarrierCode string `json:"carrier_code,optional"`
    WarehouseID int64  `json:"warehouse_id,string,optional"`
}

type ListShippingTemplatesReq {
    MarketID  int64  `form:"market_id,string,optional"`
    Name      string `form:"name,optional"`
    IsActive  *bool  `form:"is_active,optional"`
    Page      int    `form:"page,default=1"`
    PageSize  int    `form:"page_size,default=20"`
}

// ============================================================
// Shipping Zones
// ============================================================

type ShippingZone {
    Id                  int64    `json:"id,string"`
    TenantID            int64    `json:"tenant_id,string"`
    TemplateID          int64    `json:"template_id,string"`
    MarketID            int64    `json:"market_id,string"`
    Currency            string   `json:"currency"`
    Name                string   `json:"name"`
    NameI18n            []NameI18nEntry `json:"name_i18n,optional"` // P1-10
    Regions             []string `json:"regions"`
    FeeType             string   `json:"fee_type"`  // fixed|by_count|by_weight|by_volume|free
    FirstUnit           int      `json:"first_unit"`
    FirstFee            string   `json:"first_fee"`            // decimal string
    AdditionalUnit      int      `json:"additional_unit"`
    AdditionalFee       string   `json:"additional_fee"`
    FreeThresholdAmount string   `json:"free_threshold_amount"`
    FreeThresholdCount  int      `json:"free_threshold_count"`
    Taxable             bool     `json:"taxable"`              // P1-6
    TaxRate             string   `json:"tax_rate"`             // P1-6
    TaxIncluded         bool     `json:"tax_included"`         // P1-6
    IossApplicable      bool     `json:"ioss_applicable"`      // P1-6
    RemoteSurcharge     string   `json:"remote_surcharge"`     // P1-7
    RemoteZipPatterns   []string `json:"remote_zip_patterns"`  // P1-7
    FuelSurchargePct    string   `json:"fuel_surcharge_pct"`   // P1-8
    VolumetricDivisor   int      `json:"volumetric_divisor"`   // P1-9
    Sort                int      `json:"sort"`
}

type NameI18nEntry {
    Locale string `json:"locale"`  // en-US, ja-JP, ...
    Name   string `json:"name"`
}

type CreateShippingZoneReq {
    TemplateID          int64    `json:"template_id,string"`
    Currency            string   `json:"currency,optional"`
    Name                string   `json:"name"`
    NameI18n            []NameI18nEntry `json:"name_i18n,optional"`
    Regions             []string `json:"regions"`
    FeeType             string   `json:"fee_type"`
    FirstUnit           int      `json:"first_unit"`
    FirstFee            string   `json:"first_fee"`
    AdditionalUnit      int      `json:"additional_unit"`
    AdditionalFee       string   `json:"additional_fee"`
    FreeThresholdAmount string   `json:"free_threshold_amount,optional"`
    FreeThresholdCount  int      `json:"free_threshold_count,optional"`
    Taxable             bool     `json:"taxable,optional"`
    TaxRate             string   `json:"tax_rate,optional"`
    TaxIncluded         bool     `json:"tax_included,optional"`
    IossApplicable      bool     `json:"ioss_applicable,optional"`
    RemoteSurcharge     string   `json:"remote_surcharge,optional"`
    RemoteZipPatterns   []string `json:"remote_zip_patterns,optional"`
    FuelSurchargePct    string   `json:"fuel_surcharge_pct,optional"`
    VolumetricDivisor   int      `json:"volumetric_divisor,optional"`
    Sort                int      `json:"sort,optional"`
}

type UpdateShippingZoneReq {
    Id                  int64    `path:"id"`
    Name                string   `json:"name,optional"`
    NameI18n            []NameI18nEntry `json:"name_i18n,optional"`
    Regions             []string `json:"regions,optional"`
    FeeType             string   `json:"fee_type,optional"`
    FirstUnit           int      `json:"first_unit,optional"`
    FirstFee            string   `json:"first_fee,optional"`
    AdditionalUnit      int      `json:"additional_unit,optional"`
    AdditionalFee       string   `json:"additional_fee,optional"`
    FreeThresholdAmount string   `json:"free_threshold_amount,optional"`
    FreeThresholdCount  int      `json:"free_threshold_count,optional"`
    Taxable             bool     `json:"taxable,optional"`
    TaxRate             string   `json:"tax_rate,optional"`
    TaxIncluded         bool     `json:"tax_included,optional"`
    IossApplicable      bool     `json:"ioss_applicable,optional"`
    RemoteSurcharge     string   `json:"remote_surcharge,optional"`
    RemoteZipPatterns   []string `json:"remote_zip_patterns,optional"`
    FuelSurchargePct    string   `json:"fuel_surcharge_pct,optional"`
    VolumetricDivisor   int      `json:"volumetric_divisor,optional"`
    Sort                int      `json:"sort,optional"`
}

// ============================================================
// Calculator
// ============================================================

type CalculatorAddress {
    CountryCode  string `json:"country_code"`   // ISO 3166-1 alpha-2 (REQUIRED)
    ProvinceCode string `json:"province_code,optional"`  // state/province
    CityCode     string `json:"city_code"`              // 兼容旧中国城市码
    DistrictCode string `json:"district_code,optional"`
    PostalCode   string `json:"postal_code,optional"`   // 邮编（P1-7 偏远判定）
}

type CalculatorItem {
    ProductID string `json:"product_id"`
    SKUID     string `json:"sku_id,optional"`
    Quantity  int    `json:"quantity"`
    Weight    int    `json:"weight"`           // 单位：克（统一）
    Length    int    `json:"length,optional"`  // mm（P1-9 体积）
    Width     int    `json:"width,optional"`   // mm
    Height    int    `json:"height,optional"`  // mm
    Price     string `json:"price"`            // decimal string
}

type CalculateShippingFeeReq {
    MarketID int64             `json:"market_id,string"`  // 多市场必填
    Address  CalculatorAddress `json:"address"`
    Items    []CalculatorItem  `json:"items"`
}

type FeeCalculationDetail {
    FeeType          string `json:"fee_type"`
    FirstUnit        int    `json:"first_unit"`
    FirstFee         string `json:"first_fee"`
    AdditionalUnit   int    `json:"additional_unit"`
    AdditionalFee    string `json:"additional_fee"`
    CalculatedWeight int    `json:"calculated_weight"`  // 计费重（实重 or 体积重）
    VolumetricWeight int    `json:"volumetric_weight,optional"`
    CalculatedUnits  int    `json:"calculated_units"`
    AppliedSurcharge string `json:"applied_surcharge,optional"` // 总附加费
    AppliedTax       string `json:"applied_tax,optional"`       // 税费
}

type CalculateShippingFeeResp {
    ShippingFee       string               `json:"shipping_fee"`           // 不含税基础运费
    Tax               string               `json:"tax,optional"`           // 税费
    Total             string               `json:"total"`                  // 含税总价
    Currency          string               `json:"currency"`
    PriceIncludesTax  bool                 `json:"price_includes_tax"`     // 展示口径
    TemplateID        int64                `json:"template_id,string"`
    TemplateName      string               `json:"template_name"`
    ZoneName          string               `json:"zone_name"`
    CarrierCode       string               `json:"carrier_code"`
    EstimatedDays     int                  `json:"estimated_days,optional"` // P2-12
    FeeDetail         FeeCalculationDetail `json:"fee_detail"`
}
```

- [ ] **Step 2: 重新生成代码**

```bash
cd admin && make api
```

Expected: `internal/types/types.go` 与 `internal/handler/routes.go` 重新生成；旧的 `Logic` 文件保留为下一步手写实现。

- [ ] **Step 3: 编译验证（不要 commit，运行测试看错误）**

```bash
cd admin && go build ./...
```

Expected: 编译错误（因为旧的 logic 不再匹配新签名），按错误清单手写更新 logic。

- [ ] **Step 4: 提交生成的代码**

```bash
cd admin && git add desc/shipping.api internal/types/types.go internal/handler/
git commit -m "feat(shipping): regenerate API types for internationalization"
```

---

### Task 1.7: 更新 logic 层 — 模板 CRUD

**Files:**
- Modify: `admin/internal/logic/shipping_templates/create_shipping_template_logic.go`
- Modify: `admin/internal/logic/shipping_templates/update_shipping_template_logic.go`
- Modify: `admin/internal/logic/shipping_templates/list_shipping_templates_logic.go`
- Modify: `admin/internal/logic/shipping_templates/get_shipping_template_logic.go`

- [ ] **Step 1: 更新 Create**

```go
// create_shipping_template_logic.go
func (l *CreateShippingTemplateLogic) CreateShippingTemplate(req *types.CreateShippingTemplateReq) (*types.CreateShippingTemplateResp, error) {
    // 校验
    if req.Name == "" {
        return nil, code.ErrShippingTemplateNameRequired
    }
    if req.Currency != "" && !isValidCurrency(req.Currency) {
        return nil, code.ErrShippingTemplateInvalidCurrency
    }

    template := &shipping.ShippingTemplate{
        TenantID:    tenantID(l.ctx),
        MarketID:    req.MarketID,
        Currency:    defaultCurrency(req.Currency),
        CarrierCode: defaultCarrierCode(req.CarrierCode),
        WarehouseID: req.WarehouseID,
        Name:        req.Name,
        IsDefault:   req.IsDefault,
        IsActive:    true,
    }

    if err := l.svcCtx.ShippingRepo.Create(l.ctx, l.svcCtx.DB, template); err != nil {
        return nil, err
    }
    return &types.CreateShippingTemplateResp{ID: int64(template.ID), Name: template.Name}, nil
}

func isValidCurrency(c string) bool {
    switch c {
    case "CNY", "USD", "EUR", "GBP", "JPY", "KRW", "SGD", "MYR", "THB", "IDR", "PHP", "VND":
        return true
    }
    return false
}

func defaultCurrency(c string) string {
    if c == "" { return "CNY" }
    return c
}

func defaultCarrierCode(c string) string {
    if c == "" { return "standard" }
    return c
}
```

- [ ] **Step 2: 更新 List — 增加 market_id 过滤**

```go
// list_shipping_templates_logic.go
func (l *ListShippingTemplatesLogic) ListShippingTemplates(req *types.ListShippingTemplatesReq) (*types.ListShippingTemplatesResp, error) {
    templates, total, err := l.svcCtx.ShippingRepo.FindListByMarket(
        l.ctx, l.svcCtx.DB,
        req.MarketID,
        req.Name,
        req.IsActive,
        req.Page, req.PageSize,
    )
    if err != nil {
        return nil, err
    }
    // ... 转换为 resp
}
```

- [ ] **Step 3: 增加新错误码**

在 `pkg/code/code.go` 中追加：
```go
// Shipping Module (120xxx)
var (
    // ... 现有
    ErrShippingTemplateInvalidCurrency  = &Err{HTTPCode: 400, Code: 120016, Msg: "invalid currency code"}
    ErrShippingTemplateMarketMismatch   = &Err{HTTPCode: 400, Code: 120017, Msg: "template market does not match zone"}
    ErrShippingZoneInvalidFeeTypeV2     = &Err{HTTPCode: 400, Code: 120018, Msg: "invalid fee type for international zone"}
)
```

- [ ] **Step 4: 编译 + 运行所有 shipping 测试**

```bash
cd admin && make build && go test ./internal/logic/shipping_templates/... ./internal/domain/shipping/... -v
```

- [ ] **Step 5: 提交**

```bash
cd admin && git add .
git commit -m "feat(shipping-templates): wire MarketID and Currency through CRUD"
```

---

### Task 1.8: 更新 logic 层 — Zone CRUD（一次性包含 P1 字段）

**Files:**
- Modify: `admin/internal/logic/shipping_zones/create_shipping_zone_logic.go`
- Modify: `admin/internal/logic/shipping_zones/update_shipping_zone_logic.go`

- [ ] **Step 1: 写失败测试**

```go
// admin/internal/logic/shipping_zones/create_shipping_zone_logic_test.go
func TestCreateShippingZone_WithTax(t *testing.T) {
    // 模拟带 taxable/tax_rate 的请求
    req := &types.CreateShippingZoneReq{
        TemplateID: 100,
        Currency:   "EUR",
        Name:       "Germany",
        Regions:    []string{"DE"},
        FeeType:    "by_weight",
        FirstUnit:  1000,
        FirstFee:   "5.00",
        AdditionalUnit: 500,
        AdditionalFee:  "2.00",
        Taxable:    true,
        TaxRate:    "0.19",
        TaxIncluded: false,
    }
    // 验证 CreateShippingZone 写入正确
}
```

- [ ] **Step 2: 更新 Create 逻辑**

```go
func (l *CreateShippingZoneLogic) CreateShippingZone(req *types.CreateShippingZoneReq) (*types.CreateShippingZoneResp, error) {
    if req.TemplateID == 0 {
        return nil, code.ErrShippingZoneTemplateRequired
    }

    // 校验 fee_type 支持 by_volume (P1-9)
    feeType := shipping.FeeType(req.FeeType)
    if !feeType.IsValidV2() { // 新校验方法包含 by_volume
        return nil, code.ErrShippingZoneInvalidFeeType
    }

    // 校验币种与税率范围
    if req.Taxable {
        rate, err := decimal.NewFromString(req.TaxRate)
        if err != nil || rate.IsNegative() || rate.GreaterThan(decimal.NewFromInt(1)) {
            return nil, code.ErrShippingZoneInvalidTaxRate
        }
    }

    zone := &shipping.ShippingZone{
        TenantID:            tenantID(l.ctx),
        TemplateID:          req.TemplateID,
        Currency:            defaultCurrency(req.Currency),
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
        VolumetricDivisor:   defaultInt(req.VolumetricDivisor, 5000),
        Sort:                req.Sort,
    }

    if err := zone.Validate(); err != nil {
        return nil, err
    }

    if err := l.svcCtx.ShippingRepo.CreateZone(l.ctx, l.svcCtx.DB, zone); err != nil {
        return nil, err
    }
    return &types.CreateShippingZoneResp{ID: int64(zone.ID)}, nil
}

func toStringI18n(entries []types.NameI18nEntry) shipping.StringI18n {
    if len(entries) == 0 { return nil }
    out := shipping.StringI18n{}
    for _, e := range entries {
        if e.Locale != "" && e.Name != "" {
            out[e.Locale] = e.Name
        }
    }
    return out
}
```

- [ ] **Step 3: 更新 Validate 实体方法（包含 by_volume 与新字段）**

```go
// entity.go
func (z *ShippingZone) Validate() error {
    if z.Name == "" { return code.ErrShippingZoneNameRequired }
    if len(z.Regions) == 0 { return code.ErrShippingZoneRegionsRequired }
    if !z.FeeType.IsValidV2() { return code.ErrShippingZoneInvalidFeeType }
    if z.FeeType != FeeTypeFree {
        if z.FirstUnit <= 0 || z.FirstFee.IsNegative() ||
           z.AdditionalUnit <= 0 || z.AdditionalFee.IsNegative() {
            return code.ErrShippingZoneFeeConfigRequired
        }
    }
    // by_volume 校验：必须设置 VolumetricDivisor
    if z.FeeType == FeeTypeByVolume && z.VolumetricDivisor <= 0 {
        return code.ErrShippingZoneFeeConfigRequired
    }
    return nil
}

func (t FeeType) IsValidV2() bool {
    switch t {
    case FeeTypeFixed, FeeTypeByCount, FeeTypeByWeight, FeeTypeByVolume, FeeTypeFree:
        return true
    }
    return false
}

const (
    FeeTypeByVolume FeeType = "by_volume"
)
```

- [ ] **Step 4: Update 逻辑同步包含全部字段**

参考 Create，将每个字段 `if req.Xxxx != zero { zone.Xxxx = req.Xxxx }` 写入。

- [ ] **Step 5: 编译 + 测试**

```bash
cd admin && make build && go test ./internal/logic/shipping_zones/... -v
```

- [ ] **Step 6: 提交**

```bash
cd admin && git add .
git commit -m "feat(shipping-zones): wire currency/tax/i18n/surcharge/dimensions through CRUD"
```

---

### Task 1.9: 前端类型与 API 同步

**Files:**
- Modify: `shop-admin/src/api/shipping.ts`

- [ ] **Step 1: 同步更新 TypeScript 类型**

```typescript
// shop-admin/src/api/shipping.ts
export interface ShippingTemplate {
  id: string
  tenant_id: string
  market_id: string        // 0 = 全市场
  currency: string         // ISO 4217
  carrier_code: string
  warehouse_id: string
  name: string
  is_default: boolean
  is_active: boolean
  zone_count: number
  created_at: string
}

export interface NameI18nEntry {
  locale: string           // 'en-US', 'ja-JP', 'zh-CN' ...
  name: string
}

export interface ShippingZone {
  id: string
  tenant_id: string
  template_id: string
  market_id: string
  currency: string
  name: string
  name_i18n?: NameI18nEntry[]
  regions: string[]
  fee_type: 'fixed' | 'by_count' | 'by_weight' | 'by_volume' | 'free'
  first_unit: number
  first_fee: string
  additional_unit: number
  additional_fee: string
  free_threshold_amount: string
  free_threshold_count: number
  taxable: boolean
  tax_rate: string
  tax_included: boolean
  ioss_applicable: boolean
  remote_surcharge: string
  remote_zip_patterns: string[]
  fuel_surcharge_pct: string
  volumetric_divisor: number
  sort: number
}

export interface CalculatorAddress {
  country_code: string         // ISO 3166-1 alpha-2 (REQUIRED)
  province_code?: string
  city_code: string
  district_code?: string
  postal_code?: string
}

export interface CalculatorItem {
  product_id: string
  sku_id?: string
  quantity: number
  weight: number              // 克
  length?: number             // mm
  width?: number              // mm
  height?: number             // mm
  price: string
}

export interface CalculateShippingFeeReq {
  market_id: string           // REQUIRED
  address: CalculatorAddress
  items: CalculatorItem[]
}

export interface FeeCalculationDetail {
  fee_type: string
  first_unit: number
  first_fee: string
  additional_unit: number
  additional_fee: string
  calculated_weight: number
  volumetric_weight?: number
  calculated_units: number
  applied_surcharge?: string
  applied_tax?: string
}

export interface CalculateShippingFeeResp {
  shipping_fee: string
  tax?: string
  total: string
  currency: string
  price_includes_tax: boolean
  template_id: string
  template_name: string
  zone_name: string
  carrier_code: string
  estimated_days?: number
  fee_detail: FeeCalculationDetail
}

// 新增 useCurrency composable 在 Task 1.10
```

- [ ] **Step 2: 提交**

```bash
cd shop-admin && git add src/api/shipping.ts
git commit -m "feat(shipping): sync TypeScript types with backend"
```

---

### Task 1.10: 新增 useCurrency composable

**Files:**
- Create: `shop-admin/src/composables/useCurrency.ts`

- [ ] **Step 1: 实现 composable**

```typescript
// shop-admin/src/composables/useCurrency.ts
import { computed, type Ref } from 'vue'
import { useI18n } from 'vue-i18n'

export interface CurrencyConfig {
  code: string
  symbol: string
  decimals: number          // JPY=0, USD=2, KRW=0
  position: 'before' | 'after'
}

const CURRENCY_CONFIG: Record<string, CurrencyConfig> = {
  CNY: { code: 'CNY', symbol: '¥',  decimals: 2, position: 'before' },
  USD: { code: 'USD', symbol: '$',  decimals: 2, position: 'before' },
  EUR: { code: 'EUR', symbol: '€',  decimals: 2, position: 'before' },
  GBP: { code: 'GBP', symbol: '£',  decimals: 2, position: 'before' },
  JPY: { code: 'JPY', symbol: '¥',  decimals: 0, position: 'before' },
  KRW: { code: 'KRW', symbol: '₩',  decimals: 0, position: 'before' },
  SGD: { code: 'SGD', symbol: 'S$', decimals: 2, position: 'before' },
  MYR: { code: 'MYR', symbol: 'RM', decimals: 2, position: 'before' },
  THB: { code: 'THB', symbol: '฿',  decimals: 2, position: 'before' },
  IDR: { code: 'IDR', symbol: 'Rp', decimals: 0, position: 'before' },
  PHP: { code: 'PHP', symbol: '₱',  decimals: 2, position: 'before' },
  VND: { code: 'VND', symbol: '₫',  decimals: 0, position: 'after' },
}

export function useCurrency(currencyRef: Ref<string>) {
  const { locale } = useI18n()

  const config = computed(() =>
    CURRENCY_CONFIG[currencyRef.value] || CURRENCY_CONFIG.CNY
  )

  /**
   * 格式化金额：
   * - input: "1.99" -> "$1.99"
   * - JPY/KRW/IDR/VND 无小数
   * - VND 符号在右侧
   */
  const format = (amount: string | number): string => {
    const num = typeof amount === 'string' ? parseFloat(amount) : amount
    if (isNaN(num)) return ''
    const fixed = num.toFixed(config.value.decimals)
    const formatted = config.value.decimals > 0
      ? Number(fixed).toLocaleString(locale.value, {
          minimumFractionDigits: config.value.decimals,
          maximumFractionDigits: config.value.decimals,
        })
      : Number(fixed).toLocaleString(locale.value)
    return config.value.position === 'before'
      ? `${config.value.symbol}${formatted}`
      : `${formatted}${config.value.symbol}`
  }

  return { config, format }
}
```

- [ ] **Step 2: 写一个最小测试**

```typescript
// shop-admin/src/composables/useCurrency.test.ts
import { describe, expect, it } from 'vitest'
import { ref } from 'vue'
import { useCurrency } from './useCurrency'

describe('useCurrency', () => {
  it('formats USD with 2 decimals', () => {
    const { format } = useCurrency(ref('USD'))
    expect(format('1.99')).toBe('$1.99')
    expect(format('1234.5')).toBe('$1,234.50')
  })

  it('formats JPY with 0 decimals', () => {
    const { format } = useCurrency(ref('JPY'))
    expect(format('1234.56')).toBe('¥1,235')
  })

  it('formats VND with symbol after', () => {
    const { format } = useCurrency(ref('VND'))
    expect(format('100000')).toBe('100,000₫')
  })
})
```

- [ ] **Step 3: 提交**

```bash
cd shop-admin && git add src/composables/useCurrency.ts src/composables/useCurrency.test.ts
git commit -m "feat(currency): add useCurrency composable for multi-currency display"
```

---

### Task 1.11: 更新运费计算逻辑（去掉 CNY 硬编码）

**Files:**
- Modify: `admin/internal/logic/shipping_calculator/calculate_shipping_fee_logic.go`

- [ ] **Step 1: 修改响应构造**

```go
// 删掉 Currency: "CNY" 硬编码
return &types.CalculateShippingFeeResp{
    ShippingFee:  formatAmount(shippingFee),
    Currency:     zone.Currency,  // ← 取 zone 自己的币种
    // ...
}, nil
```

- [ ] **Step 2: 增加 marketID 参数校验**

```go
if req.MarketID == 0 {
    return nil, code.ErrShippingCalcMarketRequired  // 新错误码
}
```

- [ ] **Step 3: 仓储查询按 market 过滤**

```go
// findTemplateForItems 改为 findTemplateForItems(marketID, ...)
func (l *CalculateShippingFeeLogic) findTemplateForItems(marketID int64, cityCode string, reqItems []types.CalculatorItem) (*shipping.ShippingTemplate, *shipping.ShippingZone) {
    // Priority 1: 商品级（带 market 过滤）
    for _, item := range reqItems {
        mapping, err := l.svcCtx.ShippingRepo.FindMappingByTarget(l.ctx, l.svcCtx.DB, shipping.TargetTypeProduct, item.ProductID)
        if err == nil && mapping != nil {
            template, err := l.svcCtx.ShippingRepo.FindByID(l.ctx, l.svcCtx.DB, mapping.TemplateID)
            if err == nil && template != nil && template.IsActive &&
               (template.MarketID == marketID || template.MarketID == 0) {
                zone := l.findZoneForCity(int64(template.ID), cityCode, marketID)
                if zone != nil { return template, zone }
            }
        }
    }
    // Priority 3: 默认模板（按 market）
    defaultTemplate, err := l.svcCtx.ShippingRepo.FindDefaultByMarket(l.ctx, l.svcCtx.DB, marketID)
    if err == nil && defaultTemplate != nil {
        zone := l.findZoneForCity(int64(defaultTemplate.ID), cityCode, marketID)
        if zone != nil { return defaultTemplate, zone }
    }
    return nil, nil
}
```

- [ ] **Step 4: 增加新错误码**

```go
ErrShippingCalcMarketRequired = &Err{HTTPCode: 400, Code: 120019, Msg: "market_id is required"}
```

- [ ] **Step 5: 编译 + 测试**

```bash
cd admin && make build && go test ./internal/logic/shipping_calculator/... -v
```

- [ ] **Step 6: 提交**

```bash
cd admin && git add .
git commit -m "feat(shipping-calc): remove CNY hardcode, add market_id filtering"
```

---

### Phase 1 完成检查

- [ ] 所有 amount 字段类型正确（DB decimal / API string / domain decimal.Decimal）
- [ ] `findTemplateForItems` 按 market 过滤，回退到 market_id=0
- [ ] `pkg/code` 中所有新错误码已注册
- [ ] 前端 `useCurrency` composable 通过单元测试
- [ ] `cd admin && make build && go test ./...` 全绿
- [ ] 文档：更新 `docs/reference/error-codes.md` 中 shipping 错误码段

---

## Phase 2: 国家代码 + 区域匹配 + 海外数据导入 (P0-2, P0-3)

> **目标：** Address 增加 country_code；区域匹配支持国家级；导入 US/EU/SEA/JP+KR 区域数据。

### Task 2.1: 扩展区域表（增加 country_code）

**Files:**
- Create: `sql/fulfillment/migrations/2026072002_alter_regions_country.sql`

- [ ] **Step 1: 编写迁移**

```sql
-- sql/fulfillment/migrations/2026072002_alter_regions_country.sql
ALTER TABLE regions
    ADD COLUMN country_code VARCHAR(2) NOT NULL DEFAULT 'CN' COMMENT 'ISO 3166-1 alpha-2',
    ADD COLUMN postal_pattern VARCHAR(255) NULL COMMENT '邮编正则模式',
    ADD INDEX idx_country (country_code),
    ADD INDEX idx_country_parent (country_code, parent_code);
```

- [ ] **Step 2: 应用 + 合并到 schema.sql**

```bash
mysql -h 192.168.0.100 -u root shopjoy < sql/fulfillment/migrations/2026072002_alter_regions_country.sql
```

合并到 `sql/fulfillment/schema.sql` 的 `regions` 表定义。

- [ ] **Step 3: 提交**

```bash
git add sql/fulfillment/migrations/2026072002_alter_regions_country.sql sql/fulfillment/schema.sql
git commit -m "feat(regions): add country_code and postal_pattern columns"
```

---

### Task 2.2: 区域匹配器（支持多级）

**Files:**
- Create: `admin/internal/domain/shipping/zone_matcher.go`
- Create: `admin/internal/domain/shipping/zone_matcher_test.go`

- [ ] **Step 1: 写失败测试**

```go
package shipping

import (
    "testing"
)

func TestZoneMatcher_MatchByCountry(t *testing.T) {
    zones := []*ShippingZone{
        {ID: 1, Name: "US", Regions: Regions{"US"}, Sort: 10},
        {ID: 2, Name: "EU-DE", Regions: Regions{"DE"}, Sort: 5},
        {ID: 3, Name: "EU-FR", Regions: Regions{"FR"}, Sort: 5},
    }
    matcher := NewZoneMatcher(zones)

    got := matcher.Match(MatchInput{CountryCode: "US"})
    if got == nil || got.ID != 1 {
        t.Errorf("expected US zone, got %v", got)
    }
}

func TestZoneMatcher_MatchByProvince(t *testing.T) {
    zones := []*ShippingZone{
        {ID: 1, Name: "US-CA", Regions: Regions{"US-CA"}, Sort: 5},
        {ID: 2, Name: "US-TX", Regions: Regions{"US-TX"}, Sort: 5},
        {ID: 3, Name: "US-Default", Regions: Regions{"US"}, Sort: 100},
    }
    matcher := NewZoneMatcher(zones)

    got := matcher.Match(MatchInput{CountryCode: "US", ProvinceCode: "US-CA"})
    if got == nil || got.ID != 1 {
        t.Errorf("expected US-CA, got %v", got)
    }
}

func TestZoneMatcher_FallbackToCountry(t *testing.T) {
    zones := []*ShippingZone{
        {ID: 1, Name: "US-Default", Regions: Regions{"US"}, Sort: 100},
        {ID: 2, Name: "US-FL", Regions: Regions{"US-FL"}, Sort: 5},
    }
    matcher := NewZoneMatcher(zones)

    got := matcher.Match(MatchInput{CountryCode: "US", ProvinceCode: "US-WY"})
    if got == nil || got.ID != 1 {
        t.Errorf("expected US-Default fallback, got %v", got)
    }
}

func TestZoneMatcher_NoMatch(t *testing.T) {
    zones := []*ShippingZone{
        {ID: 1, Name: "US", Regions: Regions{"US"}, Sort: 100},
    }
    matcher := NewZoneMatcher(zones)
    got := matcher.Match(MatchInput{CountryCode: "JP"})
    if got != nil {
        t.Errorf("expected no match for JP, got %v", got)
    }
}
```

- [ ] **Step 2: 实现匹配器**

```go
// admin/internal/domain/shipping/zone_matcher.go
package shipping

// MatchInput 区域匹配输入
type MatchInput struct {
    CountryCode  string  // ISO 3166-1 alpha-2
    ProvinceCode string  // ISO 3166-2 或旧城市码
    CityCode     string  // 兼容旧中国城市码
    PostalCode   string
}

// ZoneMatcher 多级区域匹配器（province > city > country fallback）
type ZoneMatcher struct {
    zones []*ShippingZone
}

func NewZoneMatcher(zones []*ShippingZone) *ZoneMatcher {
    return &ZoneMatcher{zones: zones}
}

// Match 按优先级匹配：精确 > 城市 > 国家
func (m *ZoneMatcher) Match(in MatchInput) *ShippingZone {
    // 按 sort 升序遍历，第一个匹配胜出
    sorted := make([]*ShippingZone, len(m.zones))
    copy(sorted, m.zones)
    sort.Slice(sorted, func(i, j int) bool {
        return sorted[i].Sort < sorted[j].Sort
    })

    var countryFallback *ShippingZone
    for _, z := range sorted {
        if !z.Regions.Contains(in.ProvinceCode) && !z.Regions.Contains(in.CityCode) {
            // 不是精确匹配
            if z.Regions.Contains(in.CountryCode) {
                if countryFallback == nil || z.Sort < countryFallback.Sort {
                    countryFallback = z
                }
            }
            continue
        }
        return z
    }
    return countryFallback
}
```

- [ ] **Step 3: 测试通过**

```bash
cd admin && go test ./internal/domain/shipping/... -run TestZoneMatcher -v
```

- [ ] **Step 4: 提交**

```bash
cd admin && git add internal/domain/shipping/zone_matcher.go internal/domain/shipping/zone_matcher_test.go
git commit -m "feat(shipping): add multi-level zone matcher (country/province/city)"
```

---

### Task 2.3: 重构 findZoneForCity 使用新匹配器

**Files:**
- Modify: `admin/internal/logic/shipping_calculator/calculate_shipping_fee_logic.go`

- [ ] **Step 1: 修改 findZoneForCity**

```go
func (l *CalculateShippingFeeLogic) findZoneForCity(templateID int64, req *types.CalculateShippingFeeReq) *shipping.ShippingZone {
    zones, err := l.svcCtx.ShippingRepo.FindZonesByTemplateID(l.ctx, l.svcCtx.DB, templateID)
    if err != nil {
        return nil
    }
    matcher := shipping.NewZoneMatcher(zones)
    return matcher.Match(shipping.MatchInput{
        CountryCode:  req.Address.CountryCode,
        ProvinceCode: req.Address.ProvinceCode,
        CityCode:     req.Address.CityCode,
        PostalCode:   req.Address.PostalCode,
    })
}
```

- [ ] **Step 2: 编译 + 测试**

```bash
cd admin && make build && go test ./internal/logic/shipping_calculator/... -v
```

- [ ] **Step 3: 提交**

```bash
cd admin && git add internal/logic/shipping_calculator/
git commit -m "feat(shipping-calc): use multi-level zone matcher with country fallback"
```

---

### Task 2.4: 海外区域数据导入脚本

**Files:**
- Create: `sql/fulfillment/migrations/2026072003_seed_regions_overseas.sql`

- [ ] **Step 1: 编写种子数据 — 美国**

```sql
-- sql/fulfillment/migrations/2026072003_seed_regions_overseas.sql
-- 海外区域种子数据：US, EU(主要国家), JP, KR, SEA

-- ============================================================
-- United States
-- ============================================================
INSERT INTO regions (code, name, level, parent_code, country_code, postal_pattern, sort) VALUES
('US', 'United States', 1, '', 'US', '^[0-9]{5}(-[0-9]{4})?$', 0);

INSERT INTO regions (code, name, level, parent_code, country_code, sort) VALUES
('US-AL', 'Alabama', 2, 'US', 'US', 1),
('US-AK', 'Alaska', 2, 'US', 'US', 2),
('US-AZ', 'Arizona', 2, 'US', 'US', 3),
-- ... 50 州全部
('US-WY', 'Wyoming', 2, 'US', 'US', 50);

-- ============================================================
-- European Union 主要国家
-- ============================================================
INSERT INTO regions (code, name, level, parent_code, country_code, postal_pattern, sort) VALUES
('DE', 'Germany', 1, '', 'DE', '^[0-9]{5}$', 0),
('FR', 'France', 1, '', 'FR', '^[0-9]{5}$', 0),
('IT', 'Italy', 1, '', 'IT', '^[0-9]{5}$', 0),
('ES', 'Spain', 1, '', 'ES', '^[0-9]{5}$', 0),
('NL', 'Netherlands', 1, '', 'NL', '^[0-9]{4}\\s?[A-Z]{2}$', 0);

-- ============================================================
-- Japan / Korea
-- ============================================================
INSERT INTO regions (code, name, level, parent_code, country_code, postal_pattern, sort) VALUES
('JP', 'Japan', 1, '', 'JP', '^[0-9]{3}-[0-9]{4}$', 0);

INSERT INTO regions (code, name, level, parent_code, country_code, sort) VALUES
('JP-13', 'Tokyo', 2, 'JP', 'JP', 1),
('JP-27', 'Osaka', 2, 'JP', 'JP', 2),
-- ... 47 都道府県
('JP-47', 'Okinawa', 2, 'JP', 'JP', 47);

INSERT INTO regions (code, name, level, parent_code, country_code, postal_pattern, sort) VALUES
('KR', 'South Korea', 1, '', 'KR', '^[0-9]{5}$', 0);

INSERT INTO regions (code, name, level, parent_code, country_code, sort) VALUES
('KR-11', 'Seoul', 2, 'KR', 'KR', 1),
('KR-26', 'Busan', 2, 'KR', 'KR', 2),
-- ... 17 시/도
('KR-50', 'Jeju', 2, 'KR', 'KR', 17);

-- ============================================================
-- Southeast Asia
-- ============================================================
INSERT INTO regions (code, name, level, parent_code, country_code, postal_pattern, sort) VALUES
('SG', 'Singapore', 1, '', 'SG', '^[0-9]{6}$', 0),
('MY', 'Malaysia', 1, '', 'MY', '^[0-9]{5}$', 0),
('TH', 'Thailand', 1, '', 'TH', '^[0-9]{5}$', 0),
('PH', 'Philippines', 1, '', 'PH', '^[0-9]{4}$', 0),
('ID', 'Indonesia', 1, '', 'ID', '^[0-9]{5}$', 0),
('VN', 'Vietnam', 1, '', 'VN', '^[0-9]{6}$', 0);
```

> **说明：** 完整 50 州 / 47 都道府県 / 17 시/도 的列表为了可读性在此省略；实际 SQL 文件中需补全。建议用脚本生成而非手写。

- [ ] **Step 2: 应用迁移**

```bash
mysql -h 192.168.0.100 -u root shopjoy < sql/fulfillment/migrations/2026072003_seed_regions_overseas.sql
mysql -h 192.168.0.100 -u root shopjoy -e "SELECT country_code, COUNT(*) FROM regions GROUP BY country_code"
```

Expected: US=51 (1 root + 50 states), JP=48, KR=18, SEA=6 等

- [ ] **Step 3: 提交**

```bash
git add sql/fulfillment/migrations/2026072003_seed_regions_overseas.sql
git commit -m "feat(regions): seed overseas region data (US/EU/JP/KR/SEA)"
```

---

### Task 2.5: 区域 API 增加 country_code 过滤

**Files:**
- Modify: `admin/desc/regions.api`
- Modify: `admin/internal/logic/regions/list_regions_logic.go`

- [ ] **Step 1: 重写 regions API**

```go
// admin/desc/regions.api
type ListRegionsReq {
    CountryCode  string `form:"country_code,optional"`  // ISO 3166-1 alpha-2
    ParentCode   string `form:"parent_code,optional"`
    Level        int    `form:"level,optional"`
}

type RegionItem {
    Code         string        `json:"code"`
    Name         string        `json:"name"`
    Level        int           `json:"level"`
    ParentCode   string        `json:"parent_code"`
    CountryCode  string        `json:"country_code"`
    PostalPattern string       `json:"postal_pattern,optional"`
    Children     []*RegionItem `json:"children,optional"`
}
```

- [ ] **Step 2: 重写 ListRegionsLogic**

```go
func (l *ListRegionsLogic) ListRegions(req *types.ListRegionsReq) (*types.ListRegionsResp, error) {
    regions, err := l.svcCtx.RegionRepo.FindTree(l.ctx, l.svcCtx.DB, regionQuery{
        CountryCode: req.CountryCode,
        ParentCode:  req.ParentCode,
        Level:       req.Level,
    })
    if err != nil { return nil, err }
    return &types.ListRegionsResp{List: toRegionItems(regions)}, nil
}
```

- [ ] **Step 3: 重新生成 + 编译**

```bash
cd admin && make api && make build
```

- [ ] **Step 4: 提交**

```bash
cd admin && git add .
git commit -m "feat(regions): add country_code filter to region API"
```

---

### Phase 2 完成检查

- [ ] 海外区域数据已 seed 到 DB
- [ ] `MatchInput` 支持 country/province/city/postal_code 多级匹配
- [ ] ListRegions API 支持 country_code 过滤
- [ ] 计算逻辑使用新匹配器
- [ ] 前端地址选择器 Task 3.x 同步更新（见 Phase 3）

---

## Phase 3: 重量单位 + 体积计费 (P0-4, P1-9)

> **目标：** 重量统一为克存储，SKU 增加长宽高，FeeType 新增 `by_volume`，CalculateFee 支持体积重 vs 实重取大。

### Task 3.1: 重量单位转换器

**Files:**
- Create: `admin/internal/domain/shipping/weight_unit.go`
- Create: `admin/internal/domain/shipping/weight_unit_test.go`

- [ ] **Step 1: 写失败测试**

```go
package shipping

import "testing"

func TestWeightConverter_ToGrams(t *testing.T) {
    cases := []struct{
        unit   WeightUnit
        value  float64
        expect int
    }{
        {"g",  500,    500},
        {"kg", 1.5,    1500},
        {"lb", 2.2046, 1000},  // 1 lb = 453.592g
        {"oz", 35.274, 999},   // 1 oz = 28.349g
    }
    for _, c := range cases {
        got := c.unit.ToGrams(c.value)
        if abs(got - c.expect) > 5 {
            t.Errorf("%v.ToGrams(%v) = %d, want ~%d", c.unit, c.value, got, c.expect)
        }
    }
}

func TestWeightConverter_FromGrams(t *testing.T) {
    if got := WeightUnitLb.FromGrams(1000); abs(got - 2.2046) > 0.001 {
        t.Errorf("FromGrams = %v, want ~2.2046", got)
    }
}
```

- [ ] **Step 2: 实现**

```go
// admin/internal/domain/shipping/weight_unit.go
package shipping

type WeightUnit string

const (
    WeightUnitG  WeightUnit = "g"
    WeightUnitKg WeightUnit = "kg"
    WeightUnitLb WeightUnit = "lb"
    WeightUnitOz WeightUnit = "oz"
)

func (u WeightUnit) IsValid() bool {
    switch u {
    case WeightUnitG, WeightUnitKg, WeightUnitLb, WeightUnitOz:
        return true
    }
    return false
}

// ToGrams 转为克（int）
func (u WeightUnit) ToGrams(value float64) int {
    switch u {
    case WeightUnitG:  return int(value + 0.5)
    case WeightUnitKg: return int(value*1000 + 0.5)
    case WeightUnitLb: return int(value*453.59237 + 0.5)
    case WeightUnitOz: return int(value*28.349523125 + 0.5)
    }
    return int(value + 0.5)
}

// FromGrams 克转其他单位（float）
func (u WeightUnit) FromGrams(grams int) float64 {
    switch u {
    case WeightUnitG:  return float64(grams)
    case WeightUnitKg: return float64(grams) / 1000
    case WeightUnitLb: return float64(grams) / 453.59237
    case WeightUnitOz: return float64(grams) / 28.349523125
    }
    return float64(grams)
}
```

- [ ] **Step 3: 测试通过**

```bash
cd admin && go test ./internal/domain/shipping/... -run TestWeightConverter -v
```

- [ ] **Step 4: 提交**

```bash
cd admin && git add internal/domain/shipping/weight_unit.go internal/domain/shipping/weight_unit_test.go
git commit -m "feat(shipping): add weight unit converter (g/kg/lb/oz)"
```

---

### Task 3.2: 扩展 SKU 实体（加长宽高）

**Files:**
- Modify: `pkg/domain/product/entity.go`（如存在）或 `admin/internal/domain/product/entity.go`

- [ ] **Step 1: 检查现有 SKU 实体**

```bash
grep -n "type SKU" admin/internal/domain/product/entity.go
```

- [ ] **Step 2: 增加字段**

```go
type SKU struct {
    // ... 现有字段
    WeightG  int `gorm:"column:weight_g;not null;default:0"` // 统一存克
    LengthMm int `gorm:"column:length_mm;not null;default:0"`
    WidthMm  int `gorm:"column:width_mm;not null;default:0"`
    HeightMm int `gorm:"column:height_mm;not null;default:0"`
    WeightUnit string `gorm:"column:weight_unit;size:8;not null;default:'g'"` // 用户输入单位
}
```

- [ ] **Step 3: 写迁移**

```sql
-- sql/product/migrations/2026072001_alter_skus_dimensions.sql
ALTER TABLE skus
    ADD COLUMN weight_g     INT NOT NULL DEFAULT 0 COMMENT '重量(克)',
    ADD COLUMN length_mm    INT NOT NULL DEFAULT 0 COMMENT '长(mm)',
    ADD COLUMN width_mm     INT NOT NULL DEFAULT 0 COMMENT '宽(mm)',
    ADD COLUMN height_mm    INT NOT NULL DEFAULT 0 COMMENT '高(mm)',
    ADD COLUMN weight_unit  VARCHAR(8) NOT NULL DEFAULT 'g' COMMENT '用户输入单位 g/kg/lb/oz';
```

- [ ] **Step 4: 应用 + 提交**

```bash
mysql -h 192.168.0.100 -u root shopjoy < sql/product/migrations/2026072001_alter_skus_dimensions.sql
git add sql/product/migrations/2026072001_alter_skus_dimensions.sql admin/internal/domain/product/entity.go sql/product/schema.sql
git commit -m "feat(sku): add weight and dimensions fields"
```

---

### Task 3.3: CalculateFee 支持 by_volume

**Files:**
- Modify: `admin/internal/domain/shipping/entity.go`（CalculateFee 方法）

- [ ] **Step 1: 扩展 CalculateItem**

```go
type CalculateItem struct {
    ProductID int64
    SKUID     int64
    Quantity  int
    Weight    int             // 克（实重）
    Length    int             // mm
    Width     int             // mm
    Height    int             // mm
    Price     decimal.Decimal
}

// VolumetricWeight 计算体积重（克）：(L*W*H mm³) / VolumetricDivisor cm³/kg
// 默认除数 5000（国际快递常用）
func (item CalculateItem) VolumetricWeight(divisor int) int {
    if item.Length == 0 || item.Width == 0 || item.Height == 0 {
        return 0
    }
    // mm³ → cm³: / 1000
    volumeCm3 := (item.Length * item.Width * item.Height) / 1000
    return (volumeCm3 * 1000) / divisor  // → 克
}

// ChargeableWeight 取实重与体积重的最大值
func (item CalculateItem) ChargeableWeight(divisor int) int {
    volumetric := item.VolumetricWeight(divisor)
    if volumetric > item.Weight {
        return volumetric
    }
    return item.Weight
}
```

- [ ] **Step 2: 重写 CalculateFee**

```go
func (z *ShippingZone) CalculateFee(items []CalculateItem, orderAmount decimal.Decimal, itemCount int) (decimal.Decimal, int) {
    // 免运费条件（同前）
    if z.FreeThresholdAmount.IsPositive() && orderAmount.GreaterThanOrEqual(z.FreeThresholdAmount) {
        return decimal.Zero, 0
    }
    if z.FreeThresholdCount > 0 && itemCount >= z.FreeThresholdCount {
        return decimal.Zero, 0
    }
    if z.FeeType == FeeTypeFree {
        return decimal.Zero, 0
    }
    if z.FeeType == FeeTypeFixed {
        return z.FirstFee, 0
    }

    var totalUnits int
    divisor := z.VolumetricDivisor
    if divisor <= 0 { divisor = 5000 }

    switch z.FeeType {
    case FeeTypeByCount:
        for _, item := range items { totalUnits += item.Quantity }
    case FeeTypeByWeight, FeeTypeByVolume:
        for _, item := range items {
            cw := item.ChargeableWeight(divisor)
            totalUnits += cw * item.Quantity
        }
    }

    if totalUnits <= z.FirstUnit {
        return z.FirstFee, totalUnits
    }
    additionalUnits := (totalUnits - z.FirstUnit + z.AdditionalUnit - 1) / z.AdditionalUnit
    return z.FirstFee.Add(z.AdditionalFee.Mul(decimal.NewFromInt(int64(additionalUnits)))), totalUnits
}
```

- [ ] **Step 3: 更新调用方**

```go
// calculate_shipping_fee_logic.go
shippingFee, totalUnits := zone.CalculateFee(items, orderAmount, totalQuantity)
return &types.CalculateShippingFeeResp{
    // ...
    FeeDetail: types.FeeCalculationDetail{
        CalculatedWeight: totalChargeableWeight,  // 见下方
        VolumetricWeight: totalVolumetricWeight,
        CalculatedUnits:  totalUnits,
        // ...
    },
}, nil
```

- [ ] **Step 4: 测试**

```bash
cd admin && go test ./internal/domain/shipping/... -v
```

- [ ] **Step 5: 提交**

```bash
cd admin && git add .
git commit -m "feat(shipping): support volumetric weight and by_volume fee type"
```

---

### Task 3.4: 前端 ZoneConfigForm 增加体积与单位配置

**Files:**
- Modify: `shop-admin/src/views/shipping/components/ZoneConfigForm.vue`

- [ ] **Step 1: 扩展 form state**

```typescript
interface ZoneForm extends CreateZoneRequest {
  enable_amount_threshold: boolean
  enable_count_threshold: boolean
  enable_remote_surcharge: boolean     // 新增
  enable_fuel_surcharge: boolean       // 新增
  enable_tax: boolean                  // 新增
}

const form = reactive<ZoneForm>({
  // ... 现有字段
  fee_type: 'fixed',
  volumetric_divisor: 5000,             // 新增
  remote_surcharge: '0',
  remote_zip_patterns: [],
  fuel_surcharge_pct: '0',
  taxable: false,
  tax_rate: '0',
  tax_included: false,
  ioss_applicable: false,
  enable_remote_surcharge: false,
  enable_fuel_surcharge: false,
  enable_tax: false,
})
```

- [ ] **Step 2: 模板中增加 fee_type 选项**

```vue
<el-select v-model="form.fee_type">
  <el-option :label="t('shipping.feeTypeFixed')" value="fixed" />
  <el-option :label="t('shipping.feeTypeByCount')" value="by_count" />
  <el-option :label="t('shipping.feeTypeByWeight')" value="by_weight" />
  <el-option :label="t('shipping.feeTypeByVolume')" value="by_volume" />
  <el-option :label="t('shipping.feeTypeFree')" value="free" />
</el-select>
```

- [ ] **Step 3: 体积计费配置（fee_type=by_volume 时显示）**

```vue
<el-form-item v-if="form.fee_type === 'by_volume'" :label="t('shipping.volumetricDivisor')">
  <el-input-number v-model="form.volumetric_divisor" :min="1000" :max="10000" />
  <span class="form-hint">{{ t('shipping.volumetricDivisorHint') }}</span>
</el-form-item>
```

- [ ] **Step 4: 提交时转换**

```typescript
const data: CreateZoneRequest = {
  // ... 现有
  volumetric_divisor: form.fee_type === 'by_volume' ? form.volumetric_divisor : undefined,
  // ...
}
```

- [ ] **Step 5: 提交**

```bash
cd shop-admin && git add src/views/shipping/components/ZoneConfigForm.vue
git commit -m "feat(shipping): add volumetric fee config to ZoneConfigForm"
```

---

### Phase 3 完成检查

- [ ] WeightUnit 支持 g/kg/lb/oz 互转
- [ ] SKU 实体含 weight_g / length_mm / width_mm / height_mm
- [ ] CalculateFee 支持 by_volume 类型，使用 ChargeableWeight
- [ ] ZoneConfigForm 支持体积计费配置
- [ ] `go test ./...` 全绿
- [ ] `cd shop-admin && pnpm tsc --noEmit` 无类型错误

---

## Phase 4: 税费与 IOSS (P1-6)

> **目标：** ShippingZone 集成 Market.tax_rules；计算结果输出税与总价；EU IOSS 标记。

### Task 4.1: 税费计算器

**Files:**
- Create: `admin/internal/domain/shipping/tax_calculator.go`
- Create: `admin/internal/domain/shipping/tax_calculator_test.go`

- [ ] **Step 1: 写失败测试**

```go
package shipping

import (
    "testing"
    "github.com/shopspring/decimal"
)

func TestTaxCalculator_ExcludedFromPrice(t *testing.T) {
    fee := decimal.NewFromFloat(10.00)
    rate := decimal.RequireFromString("0.19")
    tax, total := CalculateTax(fee, rate, false)
    if !tax.Equal(decimal.RequireFromString("1.90")) {
        t.Errorf("expected 1.90 tax, got %s", tax)
    }
    if !total.Equal(decimal.RequireFromString("11.90")) {
        t.Errorf("expected 11.90 total, got %s", total)
    }
}

func TestTaxCalculator_IncludedInPrice(t *testing.T) {
    fee := decimal.RequireFromString("11.90")
    rate := decimal.RequireFromString("0.19")
    tax, total := CalculateTax(fee, rate, true)
    if !total.Equal(fee) {
        t.Errorf("included: total should equal fee, got %s", total)
    }
    // tax = 11.90 - 11.90/1.19 = 1.90
    expected := decimal.RequireFromString("1.90")
    if tax.Sub(expected).Abs().GreaterThan(decimal.RequireFromString("0.01")) {
        t.Errorf("expected ~1.90 tax, got %s", tax)
    }
}

func TestTaxCalculator_ZeroRate(t *testing.T) {
    fee := decimal.NewFromInt(10)
    tax, total := CalculateTax(fee, decimal.Zero, false)
    if !tax.IsZero() || !total.Equal(fee) {
        t.Errorf("expected zero tax, got tax=%s total=%s", tax, total)
    }
}
```

- [ ] **Step 2: 实现**

```go
// admin/internal/domain/shipping/tax_calculator.go
package shipping

import "github.com/shopspring/decimal"

// CalculateTax 计算税费
// - includedInPrice=true: 输入的 fee 已含税，拆出税费
// - includedInPrice=false: 输入的 fee 不含税，加上税费
func CalculateTax(fee, rate decimal.Decimal, includedInPrice bool) (tax, total decimal.Decimal) {
    if rate.IsZero() || fee.IsZero() {
        return decimal.Zero, fee
    }
    if includedInPrice {
        // fee = total; tax = total - total/(1+rate)
        divisor := decimal.NewFromInt(1).Add(rate)
        net := fee.Div(divisor)
        return fee.Sub(net).Round(2), fee
    }
    tax = fee.Mul(rate).Round(2)
    return tax, fee.Add(tax)
}
```

- [ ] **Step 3: 测试通过**

```bash
cd admin && go test ./internal/domain/shipping/... -run TestTaxCalculator -v
```

- [ ] **Step 4: 提交**

```bash
cd admin && git add internal/domain/shipping/tax_calculator.go internal/domain/shipping/tax_calculator_test.go
git commit -m "feat(shipping): add tax calculator (inclusive/exclusive)"
```

---

### Task 4.2: 集成税费到 CalculateShippingFee

**Files:**
- Modify: `admin/internal/logic/shipping_calculator/calculate_shipping_fee_logic.go`

- [ ] **Step 1: 修改计算逻辑**

```go
// 计算基础运费后
shippingFee, totalUnits := zone.CalculateFee(items, orderAmount, totalQuantity)

// 应用附加费（P1-7/8）
shippingFee, surcharges := applySurcharges(shippingFee, zone, req.Address.PostalCode, items)

// 计算税费
var tax decimal.Decimal
var total decimal.Decimal
if zone.Taxable {
    tax, total = shipping.CalculateTax(shippingFee, zone.TaxRate, zone.TaxIncluded)
} else {
    total = shippingFee
}

return &types.CalculateShippingFeeResp{
    ShippingFee:      formatAmount(shippingFee),
    Tax:              formatAmount(tax),
    Total:            formatAmount(total),
    Currency:         zone.Currency,
    PriceIncludesTax: zone.TaxIncluded,
    FeeDetail: types.FeeCalculationDetail{
        // ...
        AppliedSurcharge: formatAmount(surcharges),
        AppliedTax:       formatAmount(tax),
    },
}, nil
```

- [ ] **Step 2: 测试**

```bash
cd admin && go test ./internal/logic/shipping_calculator/... -v
```

- [ ] **Step 3: 提交**

```bash
cd admin && git add .
git commit -m "feat(shipping-calc): integrate tax and total into response"
```

---

### Phase 4 完成检查

- [ ] CalculateTax 支持含税/不含税两种模式
- [ ] CalculateShippingFeeResp 输出 tax + total
- [ ] IOSS 标记保留在 zone 但不影响计费逻辑（仅作为标记供后续订单流程使用）
- [ ] 前端 calculator 页面显示含税/不含税切换（见 Phase 6）

---

## Phase 5: 偏远地区与燃油附加费 (P1-7, P1-8)

> **目标：** 按邮编正则匹配偏远地区并加附加费；按比例加燃油附加费。

### Task 5.1: 附加费计算器

**Files:**
- Create: `admin/internal/domain/shipping/surcharge_calculator.go`
- Create: `admin/internal/domain/shipping/surcharge_calculator_test.go`

- [ ] **Step 1: 写失败测试**

```go
package shipping

import (
    "testing"
    "github.com/shopspring/decimal"
)

func TestSurchargeCalculator_RemoteAreaMatch(t *testing.T) {
    zone := &ShippingZone{
        RemoteSurcharge:   decimal.RequireFromString("50.00"),
        RemoteZipPatterns: StringArray{"^9[0-9]{4}$"},  // USPS 阿拉斯加 9xxxx
        FuelSurchargePct:  decimal.RequireFromString("0.15"),
    }
    in := SurchargeInput{
        BaseFee:    decimal.RequireFromString("20.00"),
        PostalCode: "99501",  // Anchorage, AK
    }
    surcharge := CalculateSurcharges(zone, in)
    if !surcharge.Remote.Equal(decimal.RequireFromString("50.00")) {
        t.Errorf("expected 50 remote, got %s", surcharge.Remote)
    }
    // fuel = 20 * 0.15 = 3.00
    if !surcharge.Fuel.Equal(decimal.RequireFromString("3.00")) {
        t.Errorf("expected 3 fuel, got %s", surcharge.Fuel)
    }
}

func TestSurchargeCalculator_NoMatch(t *testing.T) {
    zone := &ShippingZone{
        RemoteSurcharge:   decimal.RequireFromString("50.00"),
        RemoteZipPatterns: StringArray{"^9[0-9]{4}$"},
    }
    in := SurchargeInput{
        BaseFee:    decimal.RequireFromString("20.00"),
        PostalCode: "10001",  // NY
    }
    surcharge := CalculateSurcharges(zone, in)
    if !surcharge.Remote.IsZero() {
        t.Errorf("expected 0 remote for non-remote, got %s", surcharge.Remote)
    }
}

func TestSurchargeCalculator_EmptyPatterns(t *testing.T) {
    zone := &ShippingZone{RemoteSurcharge: decimal.NewFromInt(50)}
    surcharge := CalculateSurcharges(zone, SurchargeInput{BaseFee: decimal.NewFromInt(10)})
    if !surcharge.Remote.IsZero() {
        t.Errorf("expected 0 remote when patterns empty, got %s", surcharge.Remote)
    }
}
```

- [ ] **Step 2: 实现**

```go
// admin/internal/domain/shipping/surcharge_calculator.go
package shipping

import (
    "regexp"
    "github.com/shopspring/decimal"
)

type SurchargeInput struct {
    BaseFee    decimal.Decimal
    PostalCode string
    Weight     int  // 克
}

type SurchargeBreakdown struct {
    Remote decimal.Decimal
    Fuel   decimal.Decimal
    Total  decimal.Decimal
}

// CalculateSurcharges 计算附加费
func CalculateSurcharges(zone *ShippingZone, in SurchargeInput) SurchargeBreakdown {
    var b SurchargeBreakdown

    // 偏远地区
    if zone.RemoteSurcharge.IsPositive() && len(zone.RemoteZipPatterns) > 0 && in.PostalCode != "" {
        for _, pattern := range zone.RemoteZipPatterns {
            re, err := regexp.Compile(pattern)
            if err != nil { continue }
            if re.MatchString(in.PostalCode) {
                b.Remote = zone.RemoteSurcharge
                break
            }
        }
    }

    // 燃油附加费（百分比）
    if zone.FuelSurchargePct.IsPositive() {
        b.Fuel = in.BaseFee.Mul(zone.FuelSurchargePct).Round(2)
    }

    b.Total = b.Remote.Add(b.Fuel)
    return b
}
```

- [ ] **Step 3: 测试**

```bash
cd admin && go test ./internal/domain/shipping/... -run TestSurchargeCalculator -v
```

- [ ] **Step 4: 提交**

```bash
cd admin && git add internal/domain/shipping/surcharge_calculator.go internal/domain/shipping/surcharge_calculator_test.go
git commit -m "feat(shipping): add surcharge calculator (remote area + fuel)"
```

---

### Task 5.2: 前端 ZoneConfigForm 增加附加费配置

**Files:**
- Modify: `shop-admin/src/views/shipping/components/ZoneConfigForm.vue`

- [ ] **Step 1: 模板增加区块**

```vue
<el-divider content-position="left">{{ t('shipping.surchargeSection') }}</el-divider>

<el-form-item :label="t('shipping.enableRemoteSurcharge')">
  <el-switch v-model="form.enable_remote_surcharge" />
</el-form-item>

<template v-if="form.enable_remote_surcharge">
  <el-form-item :label="t('shipping.remoteSurcharge')">
    <el-input v-model="form.remote_surcharge" placeholder="0.00">
      <template #prepend>{{ currencySymbol }}</template>
    </el-input>
  </el-form-item>
  <el-form-item :label="t('shipping.remoteZipPatterns')">
    <el-input
      v-model="zipPatternInput"
      type="textarea"
      :placeholder="t('shipping.zipPatternHint')"
      @blur="addZipPattern"
    />
    <div class="pattern-list">
      <el-tag v-for="(p, i) in form.remote_zip_patterns" :key="i" closable @close="removeZipPattern(i)">
        {{ p }}
      </el-tag>
    </div>
  </el-form-item>
</template>

<el-form-item :label="t('shipping.enableFuelSurcharge')">
  <el-switch v-model="form.enable_fuel_surcharge" />
</el-form-item>

<template v-if="form.enable_fuel_surcharge">
  <el-form-item :label="t('shipping.fuelSurchargePct')">
    <el-input-number v-model="fuelPctNumber" :min="0" :max="100" :step="0.5" />
    <span class="form-hint">%</span>
  </el-form-item>
</template>
```

- [ ] **Step 2: 提交**

```bash
cd shop-admin && git add src/views/shipping/components/ZoneConfigForm.vue
git commit -m "feat(shipping): add surcharge config (remote area + fuel) to form"
```

---

### Phase 5 完成检查

- [ ] `CalculateSurcharges` 单元测试覆盖正常/偏远/无邮编三种场景
- [ ] 计算响应 `applied_surcharge` 字段填充正确
- [ ] 前端表单可录入邮编正则和燃油百分比

---

## Phase 6: 多语言 Zone 名称 (P1-10)

> **目标：** ShippingZone.NameI18n 字段生效，前端按 locale 展示对应名称。

### Task 6.1: 国际化名称解析器

**Files:**
- Create: `admin/internal/application/shipping/i18n_name_resolver.go`

- [ ] **Step 1: 实现**

```go
package shipping

import (
    "golang.org/x/text/language"
    "golang.org/x/text/language/display"
)

// ResolveZoneName 按 Accept-Language 解析 zone 显示名
func ResolveZoneName(zone *ShippingZone, acceptLanguage string) string {
    if len(zone.NameI18n) == 0 {
        return zone.Name  // fallback
    }
    // 1. 精确匹配 locale
    if v, ok := zone.NameI18n[acceptLanguage]; ok && v != "" {
        return v
    }
    // 2. 仅匹配语言部分（如 en-US → en）
    if tag, err := language.Parse(acceptLanguage); err == nil {
        base, _ := tag.Base()
        if base != language.Und {
            for locale, name := range zone.NameI18n {
                if t, err := language.Parse(locale); err == nil {
                    if b, _ := t.Base(); b == base && name != "" {
                        return name
                    }
                }
            }
        }
    }
    // 3. 取第一个非空
    for _, v := range zone.NameI18n {
        if v != "" { return v }
    }
    return zone.Name
}
```

- [ ] **Step 2: 写测试**

```go
func TestResolveZoneName(t *testing.T) {
    zone := &ShippingZone{
        Name: "华东",
        NameI18n: StringI18n{
            "en-US": "East China",
            "ja-JP": "華東",
        },
    }
    cases := []struct{ lang, want string }{
        {"en-US", "East China"},
        {"ja-JP", "華東"},
        {"en-GB", "East China"},  // 匹配语言部分
        {"fr-FR", "East China"},  // fallback 第一个
        {"", "华东"},               // 无 Accept-Language 用 Name
    }
    for _, c := range cases {
        if got := ResolveZoneName(zone, c.lang); got != c.want {
            t.Errorf("lang=%s got=%s want=%s", c.lang, got, c.want)
        }
    }
}
```

- [ ] **Step 3: 测试通过**

```bash
cd admin && go test ./internal/application/shipping/... -v
```

- [ ] **Step 4: 集成到 API 响应**

在 List/Get zone 的 logic 中调用 `ResolveZoneName(zone, acceptLang)` 并填充到响应的 `DisplayName` 字段。

- [ ] **Step 5: 提交**

```bash
cd admin && git add internal/application/shipping/
git commit -m "feat(shipping): add i18n zone name resolver with Accept-Language"
```

---

### Task 6.2: 前端增加多语言名称编辑

**Files:**
- Modify: `shop-admin/src/views/shipping/components/ZoneConfigForm.vue`

- [ ] **Step 1: 名称多语言编辑器**

```vue
<el-form-item :label="t('shipping.nameI18n')">
  <el-table :data="i18nNameRows" border>
    <el-table-column :label="t('common.locale')" width="180">
      <template #default="{ row }">
        <el-select v-model="row.locale" filterable>
          <el-option v-for="l in availableLocales" :key="l" :label="l" :value="l" />
        </el-select>
      </template>
    </el-table-column>
    <el-table-column :label="t('common.name')">
      <template #default="{ row }">
        <el-input v-model="row.name" />
      </template>
    </el-table-column>
    <el-table-column width="60">
      <template #default="{ $index }">
        <el-button link @click="i18nNameRows.splice($index, 1)">
          <el-icon><Delete /></el-icon>
        </el-button>
      </template>
    </el-table-column>
  </el-table>
  <el-button @click="addI18nRow">{{ t('shipping.addLocale') }}</el-button>
</el-form-item>
```

- [ ] **Step 2: 提交时转换**

```typescript
const data: CreateZoneRequest = {
  // ...
  name_i18n: i18nNameRows.value
    .filter(r => r.locale && r.name)
    .map(r => ({ locale: r.locale, name: r.name })),
}
```

- [ ] **Step 3: 提交**

```bash
cd shop-admin && git add src/views/shipping/components/ZoneConfigForm.vue
git commit -m "feat(shipping): add multi-language zone name editor"
```

---

### Phase 6 完成检查

- [ ] 后端按 Accept-Language 返回 DisplayName
- [ ] 前端可录入多语言名称
- [ ] 多语言 fallback 链路测试通过

---

## Phase 7: 物流商通用规则引擎 (P2-11)

> **目标：** 抽象 `Carrier` 接口，标准/经济/特快/海运等不同物流商走同一套规则引擎。

### Task 7.1: 物流商抽象接口

**Files:**
- Create: `admin/internal/domain/shipping/carrier.go`
- Create: `admin/internal/domain/shipping/fee_engine.go`

- [ ] **Step 1: 定义接口**

```go
// admin/internal/domain/shipping/carrier.go
package shipping

import (
    "context"
    "github.com/shopspring/decimal"
)

// Carrier 物流商抽象接口
type Carrier interface {
    Code() string
    Name() string

    // Quote 计算运费（按物流商自家规则）
    Quote(ctx context.Context, req QuoteRequest) (*QuoteResult, error)

    // EstimateDays 估算送达天数（按 destination country 分档）
    EstimateDays(destinationCountry string) int
}

// QuoteRequest 报价请求
type QuoteRequest struct {
    TemplateID int64
    Zone       *ShippingZone
    Items      []CalculateItem
    OrderAmount decimal.Decimal
    ItemCount  int
    Address    MatchInput
}

// QuoteResult 报价结果
type QuoteResult struct {
    BaseFee      decimal.Decimal
    Surcharges   SurchargeBreakdown
    Tax          decimal.Decimal
    Total        decimal.Decimal
    Currency     string
    EstimatedDays int
    Weight       int  // 计费重
}
```

- [ ] **Step 2: 定义默认实现 StandardCarrier**

```go
// admin/internal/domain/shipping/fee_engine.go
package shipping

import (
    "context"
    "github.com/shopspring/decimal"
)

// StandardCarrier 标准物流商（内置规则，与 zone 配置一致）
type StandardCarrier struct{}

func (StandardCarrier) Code() string { return "standard" }
func (StandardCarrier) Name() string { return "Standard Shipping" }

func (StandardCarrier) Quote(_ context.Context, req QuoteRequest) (*QuoteResult, error) {
    fee, units := req.Zone.CalculateFee(req.Items, req.OrderAmount, req.ItemCount)
    surcharge := CalculateSurcharges(req.Zone, SurchargeInput{
        BaseFee: fee,
        PostalCode: req.Address.PostalCode,
    })
    fee = fee.Add(surcharge.Total)

    var tax decimal.Decimal
    if req.Zone.Taxable {
        tax, _ = CalculateTax(fee, req.Zone.TaxRate, req.Zone.TaxIncluded)
    }
    return &QuoteResult{
        BaseFee: req.Zone.FirstFee,
        Surcharges: surcharge,
        Tax: tax,
        Total: fee.Add(tax),
        Currency: req.Zone.Currency,
        Weight: units,
    }, nil
}

func (StandardCarrier) EstimateDays(_ string) int { return 5 }

// CarrierRegistry 物流商注册中心
type CarrierRegistry struct {
    carriers map[string]Carrier
}

func NewCarrierRegistry() *CarrierRegistry {
    r := &CarrierRegistry{carriers: map[string]Carrier{}}
    r.Register(StandardCarrier{})
    return r
}

func (r *CarrierRegistry) Register(c Carrier) {
    r.carriers[c.Code()] = c
}

func (r *CarrierRegistry) Get(code string) (Carrier, bool) {
    c, ok := r.carriers[code]
    return c, ok
}
```

- [ ] **Step 3: 写测试**

```go
func TestStandardCarrier_Quote(t *testing.T) {
    carrier := StandardCarrier{}
    zone := &ShippingZone{
        FeeType: FeeTypeByWeight, FirstUnit: 1000, FirstFee: decimal.NewFromInt(5),
        AdditionalUnit: 500, AdditionalFee: decimal.NewFromInt(2),
        Currency: "USD",
    }
    result, err := carrier.Quote(context.Background(), QuoteRequest{
        Zone: zone,
        Items: []CalculateItem{{Quantity: 2, Weight: 800}},  // 总重 1600g
        OrderAmount: decimal.NewFromInt(100),
        ItemCount: 2,
    })
    if err != nil { t.Fatal(err) }
    if result.Total.IntPart() != 9 {  // 5 + (1600-1000)/500 * 2 = 5 + 4 = 9
        t.Errorf("expected 9, got %s", result.Total)
    }
}

func TestCarrierRegistry(t *testing.T) {
    r := NewCarrierRegistry()
    c, ok := r.Get("standard")
    if !ok { t.Fatal("standard not registered") }
    if c.Code() != "standard" { t.Errorf("got %s", c.Code()) }
}
```

- [ ] **Step 4: 提交**

```bash
cd admin && git add internal/domain/shipping/carrier.go internal/domain/shipping/fee_engine.go internal/domain/shipping/carrier_test.go
git commit -m "feat(shipping): add generic Carrier interface and StandardCarrier implementation"
```

---

### Task 7.2: 在 svc 中注入 CarrierRegistry

**Files:**
- Modify: `admin/internal/svc/service_context.go`
- Modify: `admin/internal/logic/shipping_calculator/calculate_shipping_fee_logic.go`

- [ ] **Step 1: ServiceContext 增加 Registry**

```go
type ServiceContext struct {
    // ...
    CarrierRegistry *shipping.CarrierRegistry
}

func NewServiceContext(c config.Config) *ServiceContext {
    return &ServiceContext{
        // ...
        CarrierRegistry: shipping.NewCarrierRegistry(),
    }
}
```

- [ ] **Step 2: Calculator 改用 Registry**

```go
func (l *CalculateShippingFeeLogic) CalculateShippingFee(req *types.CalculateShippingFeeReq) (*types.CalculateShippingFeeResp, error) {
    // ... 找到 template + zone 后
    carrier, ok := l.svcCtx.CarrierRegistry.Get(template.CarrierCode)
    if !ok {
        carrier = l.svcCtx.CarrierRegistry.Get("standard")  // fallback
    }
    result, err := carrier.Quote(l.ctx, shipping.QuoteRequest{
        TemplateID:  int64(template.ID),
        Zone:        zone,
        Items:       items,
        OrderAmount: orderAmount,
        ItemCount:   totalQuantity,
        Address:     /* MatchInput */,
    })
    if err != nil { return nil, err }

    return &types.CalculateShippingFeeResp{
        ShippingFee:      formatAmount(result.BaseFee.Add(result.Surcharges.Total)),
        Tax:              formatAmount(result.Tax),
        Total:            formatAmount(result.Total),
        Currency:         result.Currency,
        PriceIncludesTax: zone.TaxIncluded,
        TemplateID:       int64(template.ID),
        TemplateName:     template.Name,
        ZoneName:         zone.Name,
        CarrierCode:      carrier.Code(),
        EstimatedDays:    carrier.EstimateDays(req.Address.CountryCode),
        FeeDetail:        /* ... */,
    }, nil
}
```

- [ ] **Step 3: 测试 + 提交**

```bash
cd admin && make build && go test ./internal/logic/shipping_calculator/... -v
cd admin && git add .
git commit -m "feat(shipping-calc): integrate CarrierRegistry into calculator"
```

---

### Phase 7 完成检查

- [ ] `Carrier` 接口定义清晰，可扩展（特快/海运/海外仓专用等）
- [ ] `StandardCarrier` 通过单元测试
- [ ] 计算响应 `carrier_code` + `estimated_days` 字段填充

---

## Phase 8: 仓库选择 + 履约时效 (P2-12, P2-13)

> **目标：** 单仓起步 — 每个 market 绑定一个 Warehouse；Zone 配置中可指定发货仓库；响应返回估算时效。

### Task 8.1: Warehouse 实体

**Files:**
- Create: `admin/internal/domain/warehouse/entity.go`
- Create: `admin/internal/domain/warehouse/repository.go`
- Create: `sql/fulfillment/migrations/2026072004_create_warehouses.sql`

- [ ] **Step 1: 迁移**

```sql
CREATE TABLE warehouses (
    id              BIGINT       PRIMARY KEY,
    tenant_id       BIGINT       NOT NULL,
    market_id       BIGINT       NOT NULL,
    code            VARCHAR(50)  NOT NULL,
    name            VARCHAR(100) NOT NULL,
    country_code    VARCHAR(2)   NOT NULL,
    city_code       VARCHAR(20)  NOT NULL,
    address         VARCHAR(255) NOT NULL,
    contact_phone   VARCHAR(50),
    is_active       TINYINT(1)   NOT NULL DEFAULT 1,
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP    NULL,
    UNIQUE KEY uk_tenant_code (tenant_id, code, deleted_at),
    INDEX idx_market (market_id, is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='海外仓/发货仓库';
```

- [ ] **Step 2: 实体**

```go
package warehouse

import "github.com/colinrs/shopjoy/pkg/application"

type Warehouse struct {
    application.Model
    TenantID     int64  `gorm:"column:tenant_id;not null;index"`
    MarketID     int64  `gorm:"column:market_id;not null;index"`
    Code         string `gorm:"column:code;size:50;not null"`
    Name         string `gorm:"column:name;size:100;not null"`
    CountryCode  string `gorm:"column:country_code;size:2;not null"`
    CityCode     string `gorm:"column:city_code;size:20;not null"`
    Address      string `gorm:"column:address;size:255;not null"`
    ContactPhone string `gorm:"column:contact_phone;size:50"`
    IsActive     bool   `gorm:"column:is_active;not null;default:true"`
}

func (w *Warehouse) TableName() string { return "warehouses" }
```

- [ ] **Step 3: 仓储接口 + 实现**

```go
package warehouse

import (
    "context"
    "gorm.io/gorm"
)

type Repository interface {
    Create(ctx context.Context, db *gorm.DB, w *Warehouse) error
    Update(ctx context.Context, db *gorm.DB, w *Warehouse) error
    Delete(ctx context.Context, db *gorm.DB, id int64) error
    FindByID(ctx context.Context, db *gorm.DB, id int64) (*Warehouse, error)
    FindByMarket(ctx context.Context, db *gorm.DB, marketID int64) ([]*Warehouse, error)
}
```

（实现省略，遵循现有 `admin/internal/infrastructure/persistence/` 模式）

- [ ] **Step 4: 应用 + 提交**

```bash
mysql -h 192.168.0.100 -u root shopjoy < sql/fulfillment/migrations/2026072004_create_warehouses.sql
git add sql/fulfillment/migrations/2026072004_create_warehouses.sql admin/internal/domain/warehouse/ sql/fulfillment/schema.sql
git commit -m "feat(warehouse): add Warehouse entity, repository, and migration"
```

---

### Task 8.2: Warehouse API + Logic

**Files:**
- Create: `admin/desc/warehouses.api`
- Create: `admin/internal/logic/warehouses/*.go`（4 个 CRUD 文件）

- [ ] **Step 1: API 定义**

```go
// admin/desc/warehouses.api
service admin-api {
    @handler ListWarehousesHandler
    get /api/v1/warehouses (ListWarehousesReq) returns (ListWarehousesResp)

    @handler GetWarehouseHandler
    get /api/v1/warehouses/:id (GetWarehouseReq) returns (Warehouse)

    @handler CreateWarehouseHandler
    post /api/v1/warehouses (CreateWarehouseReq) returns (CreateWarehouseResp)

    @handler UpdateWarehouseHandler
    put /api/v1/warehouses/:id (UpdateWarehouseReq) returns (Warehouse)

    @handler DeleteWarehouseHandler
    delete /api/v1/warehouses/:id (DeleteWarehouseReq) returns (DeleteWarehouseResp)
}

type Warehouse {
    Id            int64  `json:"id,string"`
    TenantID      int64  `json:"tenant_id,string"`
    MarketID      int64  `json:"market_id,string"`
    Code          string `json:"code"`
    Name          string `json:"name"`
    CountryCode   string `json:"country_code"`
    CityCode      string `json:"city_code"`
    Address       string `json:"address"`
    ContactPhone  string `json:"contact_phone,optional"`
    IsActive      bool   `json:"is_active"`
    CreatedAt     string `json:"created_at"`
}

type ListWarehousesReq {
    MarketID int64 `form:"market_id,string,optional"`
    Page     int   `form:"page,default=1"`
    PageSize int   `form:"page_size,default=20"`
}
type ListWarehousesResp { List []Warehouse `json:"list"` Total int64 `json:"total"` }
// ... 其他类型省略，按现有模式
```

- [ ] **Step 2: 重新生成 + 实现 logic**

```bash
cd admin && make api
# 实现 create_shipping_zone_logic.go / update_/list_/get_/delete_ 五个文件
```

- [ ] **Step 3: 编译 + 测试**

```bash
cd admin && make build && go test ./internal/logic/warehouses/... -v
```

- [ ] **Step 4: 提交**

```bash
cd admin && git add .
git commit -m "feat(warehouse): add warehouse CRUD API and logic"
```

---

### Task 8.3: ShippingTemplate 绑定 Warehouse

- [ ] **Step 1: 在 CreateShippingTemplate / Update 逻辑中校验 warehouse_id 存在**

```go
if req.WarehouseID > 0 {
    wh, err := l.svcCtx.WarehouseRepo.FindByID(l.ctx, l.svcCtx.DB, req.WarehouseID)
    if err != nil { return nil, err }
    if wh.MarketID != req.MarketID {
        return nil, code.ErrShippingTemplateMarketMismatch
    }
}
```

- [ ] **Step 2: 在 CalculateShippingFee 响应中返回仓库信息**

```go
// 如果 template 有 WarehouseID，加载并填充到响应
var warehouse *warehouse.Warehouse
if template.WarehouseID > 0 {
    warehouse, _ = l.svcCtx.WarehouseRepo.FindByID(l.ctx, l.svcCtx.DB, template.WarehouseID)
}
return &types.CalculateShippingFeeResp{
    // ...
    Warehouse: toWarehouseInfo(warehouse),  // 包含 code/name/country/city
}, nil
```

- [ ] **Step 3: 前端模板列表显示仓库**

```vue
<el-table-column :label="t('shipping.warehouse')">
  <template #default="{ row }">
    {{ row.warehouse_name || '-' }}
  </template>
</el-table-column>
```

- [ ] **Step 4: 提交**

```bash
git add admin/internal/logic/shipping_templates/ admin/internal/logic/shipping_calculator/ shop-admin/src/views/shipping/index.vue
git commit -m "feat(shipping): bind template to warehouse and surface in responses"
```

---

### Phase 8 完成检查

- [ ] Warehouses CRUD API 可用
- [ ] ShippingTemplate 可绑定 warehouse_id（按 market 校验一致性）
- [ ] CalculateShippingFee 响应包含仓库信息

---

## 全局收尾

### Task 9.1: 错误码文档同步

**Files:**
- Modify: `docs/reference/error-codes.md`

- [ ] **Step 1: 追加 Shipping 模块新错误码段**

按现有格式补全 120013 ~ 120019 的所有新增错误码。

- [ ] **Step 2: 提交**

```bash
git add docs/reference/error-codes.md
git commit -m "docs(shipping): document new error codes for internationalization"
```

---

### Task 9.2: 端到端集成测试

**Files:**
- Create: `admin/internal/logic/shipping_calculator/calculate_integration_test.go`

- [ ] **Step 1: 编写端到端测试**

```go
// TestE2E_US_EUR_JP_FullFlow 跨多市场端到端
func TestE2E_US_EUR_JP_FullFlow(t *testing.T) {
    // 1. 创建 US 模板 + US 区域
    // 2. 创建 EU 模板 + DE 区域（带税）
    // 3. 创建 JP 模板 + JP-13 (Tokyo) 区域
    // 4. 绑定商品到 US 模板
    // 5. 计算三种地址的运费
    // 6. 验证币种、税费、远程附加费
}
```

- [ ] **Step 2: 测试通过**

```bash
cd admin && go test ./internal/logic/shipping_calculator/... -v
```

- [ ] **Step 3: 提交**

```bash
cd admin && git add .
git commit -m "test(shipping): add end-to-end multi-market calculation test"
```

---

### Task 9.3: 全量回归

```bash
cd admin && make build && go test ./... -v
cd shop-admin && pnpm tsc --noEmit && pnpm test
```

**必须全绿。**

---

### Task 9.4: 更新 CLAUDE.md / AGENTS.md

- [ ] 在 `docs/domains/fulfillment/` 下新增 `2026-07-20-shipping-internationalization-prd.md`
- [ ] 在 `docs/reference/database-overview.md` 中补充新增表与列
- [ ] 提交文档同步

```bash
git add docs/
git commit -m "docs(shipping): document internationalization feature"
```

---

## 验收清单 (Definition of Done)

### P0 必须满足

- [ ] ShippingZone 不再硬编码 CNY，币种来自 zone.Currency 字段
- [ ] CalculateShippingFeeReq 必须传 market_id
- [ ] 仓储层 FindDefaultByMarket 支持 market 隔离 + fallback
- [ ] CalculatorAddress.country_code 必填
- [ ] ZoneMatcher 支持 country/province/city 多级匹配
- [ ] 海外区域数据已 seed（US/EU/JP/KR/SEA）
- [ ] WeightUnit 支持 g/kg/lb/oz 互转
- [ ] `cd admin && make build && go test ./...` 全绿

### P1 必须满足

- [ ] CalculateTax 支持含税/不含税，输出 tax + total
- [ ] CalculateSurcharges 支持偏远邮编正则 + 燃油百分比
- [ ] ShippingZone 增加 name_i18n JSONB
- [ ] 后端按 Accept-Language 返回 DisplayName
- [ ] FeeType 增加 by_volume，使用 ChargeableWeight
- [ ] SKU 实体增加 weight_g / length_mm / width_mm / height_mm
- [ ] ZoneConfigForm 支持体积计费 + 偏远 + 燃油 + 税 + 多语言名称
- [ ] useCurrency composable 全币种格式化通过单元测试

### P2 必须满足

- [ ] Carrier 接口 + StandardCarrier + CarrierRegistry
- [ ] CalculateShippingFee 响应 carrier_code + estimated_days
- [ ] Warehouse 实体 + CRUD API
- [ ] ShippingTemplate.WarehouseID 校验 market 一致性
- [ ] 文档同步：PRD / 数据库总览 / 错误码
- [ ] 端到端集成测试通过

---

## 风险与权衡

| 风险 | 影响 | 缓解 |
|------|------|------|
| 14 项同步改造导致单次 commit 巨大 | 代码审查困难 | 每 Phase 独立提交，每个 Task 单独 commit |
| 海外区域数据 seed 不全 | 上线后部分国家无法匹配 | 先做 US（最大）+ JP（业务硬性需求），EU/SEA 渐进补 |
| 体积计费精度 | L*W*H 用 int 可能溢出 | 用 int64 mm³，最终再除；测试覆盖大件场景 |
| 税率正则 | 用户录入 `0.19` vs `19` | 后端统一接受 0-1 范围 decimal，文档明示 |
| 多语言名称 fallback 链路 | 缺失 locale 时显示生硬 | 三级 fallback（精确 → 语言 → 第一个非空）|
| `ShippingTemplate.MarketID=0` 兼容性 | 旧模板未指定市场 | 视为全市场通用；findTemplateForItems 回退到 market_id=0 |
| CarrierRegistry 全局单例 vs 每租户 | 误用其他租户物流商 | 后续按 tenant 注入；本期单例够用 |

---

## 后续可扩展点（本期不做）

- **A/B 测试**：运费模板版本化，按用户分桶
- **物流商真实 API 接入**：DHL/FedEx/UPS 报价 API 适配器（实现 `Carrier` 接口即可接入）
- **多仓智能路由**：库存 + 距离综合选仓
- **DDP/DDU 模式切换**：关税承担方
- **区域自定义导入**：运营自行上传 CSV/Excel 而非代码 seed