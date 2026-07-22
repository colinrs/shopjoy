# Phase 1 Patch Implementation Plan (Batch 1.5)

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Fix 4 Critical + 6 Important gaps discovered by Phase 1 final whole-branch review. Make Phase 1 actually production-ready.

**Architecture:** Focused patch — modify existing files (no new abstractions). Each Task = 1 commit. Backend + frontend fixes interleaved by concern, not by layer.

**Tech Stack:** Go 1.22+, go-zero, GORM, decimal.Decimal, Vue 3, TypeScript, Element Plus

## Global Constraints

- Same as Phase 1 main plan: decimal/decimal/string 三段式; pkg/code errors; application.Model; TDD; commit粒度
- Don't introduce new dependencies
- Don't add new public APIs (use existing types)
- Don't restructure code — minimal surgical changes
- This plan's task order is dependency-driven: C3 first (security), then C1+C2 (calculation), then C4 (currency), then UI forms

## File Structure

### Modify

```
admin/internal/domain/shipping/entity.go                  # C1: CalculateFee by_volume
admin/internal/logic/shipping_calculator/...              # C1+C2: volume + tax/total/carrier/ETA
admin/internal/infrastructure/persistence/shipping_repository.go  # C3: tenant_id
admin/internal/logic/shipping_zones/create_*.go           # C4: currency check
admin/internal/logic/shipping_zones/update_*.go           # Important: bool pointer
admin/internal/logic/shipping_templates/list_*.go        # Important: N+1 fix
admin/desc/shipping.api                                   # Important: bool pointer wire
admin/internal/types/types.go                             # regenerated

shop-admin/src/api/shipping.ts                            # Important: CreateZoneRequest补字段
shop-admin/src/views/shipping/components/ZoneConfigForm.vue  # Important: 补P1字段
shop-admin/src/views/shipping/calculator/index.vue       # Important: 市场/国家UI

pkg/code/code.go                                          # C3/C4 新错误码（如需）

admin/internal/domain/shipping/entity_test.go            # C1: by_volume 测试
admin/internal/logic/shipping_zones/*_test.go             # C4/Important 测试
admin/internal/infrastructure/persistence/shipping_repository_test.go  # C3 跨租户测试
admin/internal/logic/shipping_calculator/*_test.go        # C1+C2 测试
```

### No new files

---

## Task 1: C3 — 仓储层增加 tenant_id 过滤（最优先，安全级别）

**Files:**
- Modify: `admin/internal/infrastructure/persistence/shipping_repository.go:124-251`
- Modify: `admin/internal/infrastructure/persistence/shipping_repository_test.go`

**Interfaces:**
- 改函数签名: `FindDefaultByMarket(ctx, db, marketID)` → `FindDefaultByMarket(ctx, db, tenantID, marketID)`
- 改函数签名: `FindListByMarket(ctx, db, marketID, name, isActive, page, pageSize)` → `FindListByMarket(ctx, db, tenantID, marketID, ...)`
- 改函数签名: `UnsetAllDefaultByMarket(ctx, db, marketID)` → `UnsetAllDefaultByMarket(ctx, db, tenantID, marketID)`
- 更新所有调用方（logic 层）

- [ ] **Step 1: 写失败测试**

在 `shipping_repository_test.go` 新增测试：
```go
func TestShippingRepo_FindDefaultByMarket_TenantIsolation(t *testing.T) {
    // Seed: tenant=1 market=2 default template; tenant=2 market=2 default template
    // 调用 FindDefaultByMarket(ctx, db, 1, 2) 应返回 tenant=1 的模板
    // 调用 FindDefaultByMarket(ctx, db, 2, 2) 应返回 tenant=2 的模板
}
```

- [ ] **Step 2: 跑测试确认失败**

```bash
cd admin && go test ./internal/infrastructure/persistence/... -run TestShippingRepo_FindDefaultByMarket_TenantIsolation -v
```

- [ ] **Step 3: 修改接口 + 实现**

```go
// ShippingTemplateRepository 接口
type ShippingTemplateRepository interface {
    // ... 现有
    FindDefaultByMarket(ctx context.Context, db *gorm.DB, tenantID, marketID int64) (*ShippingTemplate, error)
    FindListByMarket(ctx context.Context, db *gorm.DB, tenantID, marketID int64, name string, isActive *bool, page, pageSize int) ([]*ShippingTemplate, int64, error)
    UnsetAllDefaultByMarket(ctx context.Context, db *gorm.DB, tenantID, marketID int64) error
}

// 实现：所有 WHERE 子句加 tenant_id
func (r *shippingTemplateRepo) FindDefaultByMarket(ctx, db, tenantID, marketID int64) (*ShippingTemplate, error) {
    var tmpl ShippingTemplate
    err := db.WithContext(ctx).
        Where("tenant_id = ? AND market_id = ? AND is_default = ? AND is_active = ?", tenantID, marketID, true, true).
        First(&tmpl).Error
    if err == nil { return &tmpl, nil }
    if !errors.Is(err, gorm.ErrRecordNotFound) { return nil, err }
    // fallback to market_id=0
    err = db.WithContext(ctx).
        Where("tenant_id = ? AND market_id = ? AND is_default = ? AND is_active = ?", tenantID, 0, true, true).
        First(&tmpl).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) { return nil, code.ErrShippingTemplateNotFound }
        return nil, err
    }
    return &tmpl, nil
}
```

类似处理 `FindListByMarket` 和 `UnsetAllDefaultByMarket`。

- [ ] **Step 4: 更新调用方**

```bash
cd admin && grep -rn "FindDefaultByMarket\|FindListByMarket\|UnsetAllDefaultByMarket" --include="*.go"
```

修改 logic 层调用：
- `calculate_shipping_fee_logic.go` — `FindDefaultByMarket(ctx, db, tenantID, marketID)` 提取 tenantID from ctx
- `list_shipping_templates_logic.go` — 同上
- `set_default_template_logic.go` — 同上
- `create_shipping_template_logic.go` — 同上

- [ ] **Step 5: 跑测试 + build**

```bash
cd admin && go build ./... && go test ./internal/infrastructure/persistence/... ./internal/logic/shipping_templates/... ./internal/logic/shipping_calculator/... -v
```

- [ ] **Step 6: commit**

```bash
git add admin/internal/infrastructure/persistence/shipping_repository.go admin/internal/infrastructure/persistence/shipping_repository_test.go admin/internal/logic/shipping_templates/ admin/internal/logic/shipping_calculator/
git commit -m "fix(shipping): add tenant_id filter to market-aware repo methods"
```

---

## Task 2: C1 — CalculateFee 支持 by_volume + 传入 length/width/height

**Files:**
- Modify: `admin/internal/domain/shipping/entity.go:CalcFee` + `CalculateItem` + 加 `VolumetricWeight` 方法
- Modify: `admin/internal/logic/shipping_calculator/calculate_shipping_fee_logic.go:89-99`（传入 length/width/height 到 CalculateItem）
- Modify: `admin/internal/domain/shipping/entity_test.go`

- [ ] **Step 1: 写失败测试**

```go
func TestShippingZone_CalculateFee_ByVolume(t *testing.T) {
    zone := &ShippingZone{
        FeeType:           FeeTypeByVolume,
        FirstUnit:         1000,  // 首件 1000g
        FirstFee:          decimal.RequireFromString("10.00"),
        AdditionalUnit:    500,
        AdditionalFee:     decimal.RequireFromString("3.00"),
        VolumetricDivisor: 5000,
    }
    // 商品: 长宽高 200mm × 200mm × 200mm → 体积 = 8,000,000 mm³ = 8000 cm³ → 体积重 1600g
    items := []CalculateItem{
        {ProductID: 1, Quantity: 1, Weight: 500, Length: 200, Width: 200, Height: 200, Price: decimal.NewFromInt(20)},
    }
    fee, _ := zone.CalculateFee(items, decimal.NewFromInt(20), 1)
    // 实重 500g, 体积重 1600g, 取大 = 1600g
    // 超过首件 1000g 600g, 600/500=1.2, 向上取整 = 2 续件
    // fee = 10 + 3*2 = 16
    if !fee.Equal(decimal.RequireFromString("16.00")) {
        t.Errorf("expected 16.00, got %s", fee)
    }
}
```

- [ ] **Step 2: 跑测试确认失败**

```bash
cd admin && go test ./internal/domain/shipping/... -run TestShippingZone_CalculateFee_ByVolume -v
```

- [ ] **Step 3: 实现 CalculateItem.VolumetricWeight + CalculateFee by_volume**

```go
// CalculateItem 加字段
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

// VolumetricWeight 计算体积重（克）：L*W*H / VolumetricDivisor (cm³/kg → g)
func (item CalculateItem) VolumetricWeight(divisor int) int {
    if item.Length <= 0 || item.Width <= 0 || item.Height <= 0 || divisor <= 0 {
        return 0
    }
    volumeCm3 := (item.Length * item.Width * item.Height) / 1000
    return (volumeCm3 * 1000) / divisor
}

// ChargeableWeight 取实重与体积重的较大值
func (item CalculateItem) ChargeableWeight(divisor int) int {
    volumetric := item.VolumetricWeight(divisor)
    if volumetric > item.Weight {
        return volumetric
    }
    return item.Weight
}
```

修改 `CalculateFee`：
```go
func (z *ShippingZone) CalculateFee(items []CalculateItem, orderAmount decimal.Decimal, itemCount int) (decimal.Decimal, int) {
    // ... 现有免运费检查 ...
    
    var totalUnits int
    divisor := z.VolumetricDivisor
    if divisor <= 0 { divisor = 5000 }
    
    switch z.FeeType {
    case FeeTypeByCount:
        for _, item := range items { totalUnits += item.Quantity }
    case FeeTypeByWeight, FeeTypeByVolume:  // ← by_volume 在此分支
        for _, item := range items {
            cw := item.ChargeableWeight(divisor)
            totalUnits += cw * item.Quantity
        }
    }
    
    // ... 续费计算 ...
}
```

- [ ] **Step 4: calculator logic 传入 length/width/height**

```go
// calculate_shipping_fee_logic.go
items := make([]shipping.CalculateItem, 0, len(req.Items))
for _, item := range req.Items {
    items = append(items, shipping.CalculateItem{
        ProductID: productID,
        SKUID:     skuID,
        Quantity:  item.Quantity,
        Weight:    item.Weight,
        Length:    item.Length,  // ← 新增
        Width:     item.Width,   // ← 新增
        Height:    item.Height,  // ← 新增
        Price:     price,
    })
}
```

- [ ] **Step 5: 跑测试 + commit**

```bash
cd admin && go build ./... && go test ./internal/domain/shipping/... ./internal/logic/shipping_calculator/... -v
git add admin/internal/domain/shipping/entity.go admin/internal/domain/shipping/entity_test.go admin/internal/logic/shipping_calculator/calculate_shipping_fee_logic.go
git commit -m "fix(shipping): implement by_volume billing with volumetric weight"
```

---

## Task 3: C2 — 计算器响应补全 Tax/Total/CarrierCode/EstimatedDays

**Files:**
- Modify: `admin/internal/logic/shipping_calculator/calculate_shipping_fee_logic.go:107-129`

- [ ] **Step 1: 写测试**

```go
// admin/internal/logic/shipping_calculator/calculate_shipping_fee_logic_test.go (新建)
package shipping_calculator

import (
    "context"
    "testing"
    "github.com/colinrs/shopjoy/admin/internal/svc"
    "github.com/colinrs/shopjoy/admin/internal/types"
    "github.com/stretchr/testify/assert"
)

func TestCalculateShippingFee_FillsTaxAndTotal(t *testing.T) {
    // Mock svcCtx with zone + template
    // zone: taxable=true, tax_rate=0.19, tax_included=false
    // Expected: shipping_fee + tax = total; price_includes_tax=false
}

func TestCalculateShippingFee_CarrierCode(t *testing.T) {
    // template.carrier_code = "standard"
    // Expected: response.CarrierCode = "standard"
}

func TestCalculateShippingFee_EstimatedDays(t *testing.T) {
    // Mock StandardCarrier.EstimateDays
    // Expected: response.EstimatedDays = N
}
```

- [ ] **Step 2: 跑测试确认失败**

- [ ] **Step 3: 实现 Tax/Total/CarrierCode/EstimatedDays**

```go
// 1. 计算基础费（含 volume + surcharge）
shippingFee, totalUnits := zone.CalculateFee(items, orderAmount, totalQuantity)

// 2. 计算税费
var tax decimal.Decimal
if zone.Taxable {
    tax, total = shipping.CalculateTax(shippingFee, zone.TaxRate, zone.TaxIncluded)
} else {
    total = shippingFee
}

// 3. Carrier code（从 template）
carrierCode := template.CarrierCode
if carrierCode == "" { carrierCode = "standard" }

// 4. Estimated days
estimatedDays := 0
if c, ok := l.svcCtx.CarrierRegistry.Get(carrierCode); ok {
    estimatedDays = c.EstimateDays(req.Address.CountryCode)
}

return &types.CalculateShippingFeeResp{
    ShippingFee:      formatAmount(shippingFee),
    Tax:              formatAmount(tax),
    Total:            formatAmount(total),
    Currency:         template.Currency,
    PriceIncludesTax: zone.TaxIncluded,
    CarrierCode:      carrierCode,
    EstimatedDays:    int64(estimatedDays),
    // ... 其他现有字段
}, nil
```

- [ ] **Step 4: 跑测试 + commit**

```bash
cd admin && go build ./... && go test ./internal/logic/shipping_calculator/... -v
git add admin/internal/logic/shipping_calculator/
git commit -m "fix(shipping-calc): populate Tax/Total/CarrierCode/EstimatedDays in response"
```

---

## Task 4: C4 — Zone currency 与 Template currency 一致性校验

**Files:**
- Modify: `admin/internal/logic/shipping_zones/create_shipping_zone_logic.go:79-102`
- Modify: `admin/internal/logic/shipping_zones/update_shipping_zone_logic.go`（如有 currency 变更）

- [ ] **Step 1: 写测试**

```go
func TestCreateShippingZone_CurrencyMismatch(t *testing.T) {
    // template.Currency = USD
    // req.Currency = EUR
    // Expected: ErrShippingTemplateMarketMismatch
}
```

- [ ] **Step 2: 跑测试确认失败**

- [ ] **Step 3: 实现一致性校验**

```go
// create_shipping_zone_logic.go: 在 CreateZone 前
template, err := l.svcCtx.ShippingRepo.FindByID(l.ctx, l.svcCtx.DB, req.TemplateID)
if err != nil { return nil, err }
if template.TenantID != tenantID(l.ctx) { return nil, code.ErrShippingTemplateNotFound }

// currency 一致性
zoneCurrency := defaultCurrency(req.Currency)
if template.Currency != "" && zoneCurrency != "" && template.Currency != zoneCurrency {
    return nil, code.ErrShippingTemplateMarketMismatch
}
```

- [ ] **Step 4: 跑测试 + commit**

```bash
cd admin && go test ./internal/logic/shipping_zones/... -v
git add admin/internal/logic/shipping_zones/create_shipping_zone_logic.go admin/internal/logic/shipping_zones/create_shipping_zone_logic_test.go
git commit -m "fix(shipping): enforce currency consistency between template and zone"
```

---

## Task 5: Important — Bool 字段用 Pointer + N+1 修复 + TaxRate 严格解析（批量）

**Files:**
- Modify: `admin/desc/shipping.api` —— UpdateReq 字段改 `*bool`
- Modify: `admin/internal/logic/shipping_zones/update_shipping_zone_logic.go` —— 接收 `*bool`
- Modify: `admin/internal/logic/shipping_templates/list_shipping_templates_logic.go` —— 批量 CountZonesByTemplateID
- Modify: `admin/internal/logic/shipping_zones/create_shipping_zone_logic.go` —— TaxRate 严格校验

- [ ] **Step 1: 修改 .api（Taxable/IossApplicable 改 pointer）**

```go
// admin/desc/shipping.api
type UpdateShippingZoneReq {
    Id                  int64    `path:"id"`
    // ...
    Taxable             *bool    `json:"taxable,optional"`
    TaxIncluded         *bool    `json:"tax_included,optional"`
    IossApplicable      *bool    `json:"ioss_applicable,optional"`
    // ...
}
```

- [ ] **Step 2: regen + 修补 update logic**

```bash
cd admin && make api
```

```go
// update_shipping_zone_logic.go
if req.Taxable != nil {
    zone.Taxable = *req.Taxable
}
if req.TaxIncluded != nil {
    zone.TaxIncluded = *req.TaxIncluded
}
// ... 同 Taxable/IossApplicable
```

- [ ] **Step 3: 批量 CountZonesByTemplateID**

新增仓储方法：
```go
func (r *shippingTemplateRepo) CountZonesByTemplateIDs(ctx context.Context, db *gorm.DB, ids []int64) (map[int64]int64, error) {
    type result struct {
        TemplateID int64
        Count      int64
    }
    var results []result
    err := db.WithContext(ctx).
        Model(&shipping.ShippingZone{}).
        Select("template_id, COUNT(*) as count").
        Where("template_id IN ?", ids).
        Group("template_id").
        Find(&results).Error
    if err != nil { return nil, err }
    
    out := make(map[int64]int64, len(results))
    for _, r := range results {
        out[r.TemplateID] = r.Count
    }
    return out, nil
}
```

list logic 调用：
```go
templates, total, err := l.svcCtx.ShippingRepo.FindListByMarket(...)
ids := make([]int64, len(templates))
for i, t := range templates { ids[i] = int64(t.ID) }
counts, _ := l.svcCtx.ShippingRepo.CountZonesByTemplateIDs(ctx, db, ids)
// 填到响应
```

- [ ] **Step 4: TaxRate 严格解析**

```go
// create_shipping_zone_logic.go
if req.Taxable {
    rate, err := decimal.NewFromString(req.TaxRate)
    if err != nil {
        return nil, code.ErrShippingZoneInvalidTaxRate  // 新错误码
    }
    if rate.IsNegative() || rate.GreaterThan(decimal.NewFromInt(1)) {
        return nil, code.ErrShippingZoneInvalidTaxRate
    }
}
```

- [ ] **Step 5: 跑测试 + commit**

```bash
cd admin && go build ./... && go test ./...
git add admin/desc/shipping.api admin/internal/types/types.go admin/internal/logic/shipping_zones/ admin/internal/logic/shipping_templates/ admin/internal/infrastructure/persistence/shipping_repository.go
git commit -m "fix(shipping): pointer bools for updates + batch zone count + strict tax rate"
```

---

## Task 6: Important — ZoneConfigForm 补 P1 字段 + Calculator 视图 UI

**Files:**
- Modify: `shop-admin/src/views/shipping/components/ZoneConfigForm.vue`
- Modify: `shop-admin/src/api/shipping.ts`（确认 CreateZoneRequest 字段完整）
- Modify: `shop-admin/src/views/shipping/calculator/index.vue`

- [ ] **Step 1: 检查 CreateZoneRequest 字段**

```typescript
// shop-admin/src/api/shipping.ts
export interface CreateZoneRequest {
  name: string
  name_i18n?: NameI18nEntry[]
  regions: string[]
  fee_type: 'fixed' | 'by_count' | 'by_weight' | 'by_volume' | 'free'
  first_unit: number
  first_fee: string
  additional_unit: number
  additional_fee: string
  free_threshold_amount?: string
  free_threshold_count?: number
  taxable?: boolean
  tax_rate?: string
  tax_included?: boolean
  ioss_applicable?: boolean
  remote_surcharge?: string
  remote_zip_patterns?: string[]
  fuel_surcharge_pct?: string
  volumetric_divisor?: number
  sort?: number
}
```

（对照后端 `.api` 确认完整）

- [ ] **Step 2: ZoneConfigForm 增加字段**

按 Phase 1 main plan Task 3.4 + 5.2 + 6.2 实施（已存在 brief 但未实际完成）：
- fee_type 加 by_volume
- volumetric_divisor（by_volume 时显示）
- taxable + tax_rate + tax_included + ioss_applicable 开关
- remote_surcharge + remote_zip_patterns
- fuel_surcharge_pct
- name_i18n 多语言编辑器

- [ ] **Step 3: Calculator 视图加市场选择 UI**

```vue
<el-form-item :label="t('shipping.market')">
  <el-select v-model="selectedMarketId">
    <el-option v-for="m in marketOptions" :key="m.id" :label="m.name" :value="m.id" />
  </el-select>
</el-form-item>
<el-form-item :label="t('shipping.countryCode')">
  <el-input v-model="address.country_code" placeholder="US, DE, JP, ..." />
</el-form-item>
```

- [ ] **Step 4: 跑测试 + commit**

```bash
cd shop-admin && pnpm vue-tsc --noEmit && pnpm build
git add shop-admin/src/
git commit -m "fix(shipping-ui): add P1 fields to ZoneConfigForm and market/country UI to calculator"
```

---

## Task 7: 收尾回归

```bash
cd admin && make build && go test ./...
cd shop-admin && pnpm vue-tsc --noEmit && pnpm build && pnpm test
```

全部绿才视为 Patch 完成。

更新 `docs/reference/error-codes.md` 添加 `230308 | 400 | market_id is required` 一行（如未加）。

---

## 验收清单

- [ ] C1: by_volume 区域按体积重计费正确
- [ ] C2: CalculateShippingFeeResp Tax/Total/CarrierCode/EstimatedDays 全部填充
- [ ] C3: FindDefaultByMarket/FindListByMarket/UnsetAllDefaultByMarket 都按 tenant_id 过滤 + sqlmock 跨租户测试
- [ ] C4: zone currency 与 template currency 不一致返回 ErrShippingTemplateMarketMismatch
- [ ] Important: bool 字段 update 可显式设 false
- [ ] Important: 列表无 N+1
- [ ] Important: TaxRate 解析失败不静默归零
- [ ] Important: ZoneConfigForm 可录入 P1 全部字段
- [ ] Important: Calculator 视图有市场/国家选择 UI
- [ ] `make build` + `pnpm build` + 全测试绿
- [ ] error-codes.md 文档同步
- [ ] 最终全分支 review 重跑通过