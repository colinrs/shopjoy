# API 时间与金额类型规范修复设计

## 1. 概述

### 1.1 背景

根据项目规范（`rules/golang/time.md` 和 `rules/golang/price.md`），API 定义文件中：
- **时间字段**：必须使用 `string` 类型，格式为 RFC3339
- **金额字段**：必须使用 `string` 类型，表示 **yuan（元）**，而非 cents

在审查 `admin/desc/` 目录下的 API 定义文件时，发现以下问题：
1. `storefront.api` 中 `ChangedAt` 字段使用了 `int64` 而非 `string`
2. `fulfillment.api` 和 `shipping.api` 中金额字段的注释说明使用 "cents"，与规范不符

### 1.2 修复目标

1. 修正 `storefront.api` 中 `ChangedAt` 的类型定义
2. 修正金额字段注释，统一为 "yuan" 或移除误导性注释
3. 全链路跟进，确保类型转换正确

---

## 2. 问题详情

### 2.1 问题一：storefront.api - ChangedAt 类型错误

**文件**: `admin/desc/storefront.api`
**位置**: `CurrentThemeResponse` 结构体第 43 行

```go
// 当前（错误）
CurrentThemeResponse {
    Theme     *ThemeItem     `json:"theme"`
    Config    ThemeConfigDTO `json:"config"`
    ChangedAt int64          `json:"changed_at,omitempty"`    // ❌ 错误：int64 类型
    ChangedBy int64          `json:"changed_by,omitempty"`
}
```

```go
// 修复后
CurrentThemeResponse {
    Theme     *ThemeItem     `json:"theme"`
    Config    ThemeConfigDTO `json:"config"`
    ChangedAt string        `json:"changed_at,omitempty"`     // ✅ 正确：string 类型 (RFC3339)
    ChangedBy int64          `json:"changed_by,omitempty"`
}
```

**影响范围**:
- `admin/desc/storefront.api` - API 定义
- `admin/internal/types/types.go` - 自动生成代码
- `admin/internal/logic/storefront/` - 业务逻辑层
- `admin/internal/domain/storefront/` - 领域层（需检查 entity）

### 2.2 问题二：fulfillment.api - AdjustAmount 注释错误

**文件**: `admin/desc/fulfillment.api`
**位置**: `AdjustOrderPriceReq` 结构体第 378 行

```go
// 当前（错误注释）
AdjustOrderPriceReq {
    ID           int64  `path:"id"`
    AdjustAmount string `json:"adjust_amount"` // Positive = increase, Negative = decrease (in cents) ❌
    Reason       string `json:"reason"`
}
```

```go
// 修复后（正确注释）
AdjustOrderPriceReq {
    ID           int64  `path:"id"`
    AdjustAmount string `json:"adjust_amount"` // Positive = increase, Negative = decrease (in yuan) ✅
    Reason       string `json:"reason"`
}
```

### 2.3 问题三：shipping.api - 运费字段注释错误

**文件**: `admin/desc/shipping.api`
**位置**: 多处 `FirstFee` 和 `AdditionalFee` 字段

| 行号 | 结构体 | 字段 | 当前注释 | 修复后注释 |
|------|--------|------|----------|------------|
| 65 | `ShippingZoneDetail` | `FirstFee` | `(cents)` | `(yuan)` |
| 67 | `ShippingZoneDetail` | `AdditionalFee` | `(cents)` | `(yuan)` |
| 79 | `CreateShippingZoneReq` | `FirstFee` | `(cents)` | `(yuan)` |
| 81 | `CreateShippingZoneReq` | `AdditionalFee` | `(cents)` | `(yuan)` |
| 93 | `UpdateShippingZoneReq` | `FirstFee` | `(cents)` | `(yuan)` |
| 95 | `UpdateShippingZoneReq` | `AdditionalFee` | `(cents)` | `(yuan)` |
| 164 | `FeeCalculationDetail` | `FirstFee` | `(cents)` | `(yuan)` |
| 166 | `FeeCalculationDetail` | `AdditionalFee` | `(cents)` | `(yuan)` |

---

## 3. 实施步骤

### Phase 1: 修复 storefront.api ChangedAt 类型

#### Step 1.1: 修改 API 定义文件

**文件**: `admin/desc/storefront.api`

```go
// 修改 CurrentThemeResponse 结构体
CurrentThemeResponse {
    Theme     *ThemeItem     `json:"theme"`
    Config    ThemeConfigDTO `json:"config"`
    ChangedAt string        `json:"changed_at,omitempty"`    // int64 → string
    ChangedBy int64          `json:"changed_by,omitempty"`
}
```

#### Step 1.2: 重新生成 API 代码

```bash
cd admin && make api
```

#### Step 1.3: 检查并更新 logic 层

**需要检查的文件**:
- `admin/internal/logic/storefront/get_current_theme_logic.go`
- `admin/internal/logic/storefront/switch_theme_logic.go`
- `admin/internal/logic/storefront/update_theme_config_logic.go`

**变更要点**:
- 如果 logic 中将 `ChangedAt` 作为 `int64`（Unix timestamp）处理，需要修改为 `time.Time`
- 转换为 RFC3339 字符串输出给 API

**示例修改**:

```go
// 修改前
type CurrentThemeResp struct {
    Theme     *ThemeItem    `json:"theme"`
    Config    ThemeConfigDTO `json:"config"`
    ChangedAt int64         `json:"changed_at,omitempty"`
    ChangedBy int64         `json:"changed_by,omitempty"`
}

// 修改后
type CurrentThemeResp struct {
    Theme     *ThemeItem    `json:"theme"`
    Config    ThemeConfigDTO `json:"config"`
    ChangedAt string        `json:"changed_at,omitempty"` // RFC3339 format
    ChangedBy int64         `json:"changed_by,omitempty"`
}

// 转换逻辑
func toCurrentThemeResponse(theme *Theme, config *ThemeConfig, changedAt time.Time) *CurrentThemeResp {
    return &CurrentThemeResp{
        Theme:     toThemeItem(theme),
        Config:    toThemeConfigDTO(config),
        ChangedAt: changedAt.Format(time.RFC3339), // int64 → RFC3339 string
        ChangedBy: config.ChangedBy,
    }
}
```

#### Step 1.4: 检查 domain entity 层

**需要检查的文件**:
- `admin/internal/domain/storefront/` 下的 entity 文件

**变更要点**:
- 如果 entity 中 `ChangedAt` 使用 `int64`，需要改为 `time.Time`
- Repository 层确保 `time.Time` 与数据库 `TIMESTAMP` 正确转换

#### Step 1.5: 编译验证

```bash
cd admin && make build
```

---

### Phase 2: 修正金额字段注释

#### Step 2.1: 修正 fulfillment.api

**文件**: `admin/desc/fulfillment.api`

修改 `AdjustOrderPriceReq` 结构体注释：

```go
// 修改前
AdjustAmount string `json:"adjust_amount"` // Positive = increase, Negative = decrease (in cents)

// 修改后
AdjustAmount string `json:"adjust_amount"` // Positive = increase, Negative = decrease (in yuan)
```

#### Step 2.2: 修正 shipping.api

**文件**: `admin/desc/shipping.api`

修改所有 `FirstFee` 和 `AdditionalFee` 字段的注释（共 8 处）：

```go
// 修改前
FirstFee string `json:"first_fee"` // First fee (cents)

// 修改后
FirstFee string `json:"first_fee"` // First fee (yuan)
```

#### Step 2.3: 重新生成 API 代码

```bash
cd admin && make api
```

#### Step 2.4: 确认业务逻辑

**需要确认的文件**:
- `admin/internal/logic/fulfillment/adjdust_order_price_logic.go`
- `admin/internal/logic/shipping/` 下的运费计算逻辑

**确认要点**:
- 业务逻辑层是否按 yuan 处理金额
- 如果业务逻辑按 cents 处理，需要在 API 层做单位转换

#### Step 2.5: 编译验证

```bash
cd admin && make build
```

---

## 4. 变更文件清单

### 4.1 API 定义文件（需手动修改）

| 文件 | 修改内容 |
|------|----------|
| `admin/desc/storefront.api` | `ChangedAt` 类型 `int64` → `string` |
| `admin/desc/fulfillment.api` | `AdjustAmount` 注释 `(cents)` → `(yuan)` |
| `admin/desc/shipping.api` | `FirstFee`/`AdditionalFee` 注释 `(cents)` → `(yuan)` (8处) |

### 4.2 自动生成文件（make api）

| 文件 | 说明 |
|------|------|
| `admin/internal/types/types.go` | API 类型定义自动重新生成 |

### 4.3 Logic 层文件（需检查并可能修改）

| 文件 | 检查项 |
|------|--------|
| `admin/internal/logic/storefront/get_current_theme_logic.go` | `ChangedAt` 类型转换 |
| `admin/internal/logic/storefront/switch_theme_logic.go` | `ChangedAt` 类型转换 |
| `admin/internal/logic/storefront/update_theme_config_logic.go` | `ChangedAt` 类型转换 |
| `admin/internal/logic/fulfillment/adjdust_order_price_logic.go` | 确认金额单位 |
| `admin/internal/logic/shipping/` | 确认运费金额单位 |

### 4.4 Domain 层文件（需检查）

| 文件 | 检查项 |
|------|--------|
| `admin/internal/domain/storefront/` | entity `ChangedAt` 类型 |

---

## 5. 验证清单

- [ ] `admin/desc/storefront.api` 中 `ChangedAt` 已改为 `string` 类型
- [ ] `admin/desc/fulfillment.api` 中 `AdjustAmount` 注释已修正
- [ ] `admin/desc/shipping.api` 中 `FirstFee`/`AdditionalFee` 注释已修正（8处）
- [ ] `make api` 成功执行，无错误
- [ ] `make build` 成功编译，无错误
- [ ] Logic 层已正确处理 `time.Time` ↔ RFC3339 string 转换
- [ ] Domain entity 层 `ChangedAt` 类型一致

---

## 6. 风险评估

| 风险 | 等级 | 缓解措施 |
|------|------|----------|
| `ChangedAt` 类型变更可能导致前端显示异常 | 中 | 前端需同步适配 RFC3339 字符串解析 |
| 金额单位从 cents 改为 yuan 需确认业务逻辑 | 高 | 先确认业务逻辑层实际处理单位再修改 |
| 自动生成的 types.go 可能覆盖手动修改 | 低 | 使用 `make api` 而非手动修改 types.go |

---

## 7. 后续跟进

本次修复仅涉及 admin 服务。如 shop 服务有相同问题，需在 `shop/desc/` 目录下执行相同审查和修复流程。

---

## 8. 参考文档

- [时间处理规范](../../rules/golang/time.md)
- [价格处理规范](../../rules/golang/price.md)
- [API 定义修改流程](../../guides/developer-guide.md)

---

## 9. 实施进度

> 开始时间: 2026-03-29

### Phase 1: 修复 storefront.api ChangedAt 类型
**状态**: ✅ 完成

**执行者**: subagent (voltagent-core-dev:backend-developer)
**Agent ID**: aac9550ce23802b78

**任务清单**:
- [x] 修改 `admin/desc/storefront.api` 中 `ChangedAt` 类型 `int64` → `string`
- [x] 运行 `make api` 重新生成 types.go
- [x] 检查并更新 logic 层 `ChangedAt` 类型转换
- [x] 检查 domain entity 层 `ChangedAt` 类型
- [x] 运行 `make build` 验证

**修改文件**:
| 文件 | 修改内容 |
|------|---------|
| `admin/desc/storefront.api` | `ChangedAt` 从 `int64` 改为 `string` (RFC3339 format) |
| `admin/internal/application/storefront/service.go` | `CurrentThemeDTO.ChangedAt` 从 `int64` 改为 `string`；添加 `ChangedAt`/`ChangedBy` 字段填充 |
| `admin/internal/types/types.go` | 自动生成，包含正确的 `ChangedAt string` 类型 |

**验证结果**:
- ✅ `make api` 成功
- ✅ `make build` 编译通过

**完成时间**: 2026-03-29

---

### Phase 2: 修正金额字段注释
**状态**: ✅ 完成

**执行者**: subagent (voltagent-core-dev:backend-developer)
**Agent ID**: ae557a6ccb96f5cda

**任务清单**:
- [x] 修正 `fulfillment.api` `AdjustAmount` 注释
- [x] 修正 `shipping.api` `FirstFee`/`AdditionalFee` 注释
- [x] 运行 `make api` 重新生成代码
- [x] 确认业务逻辑层金额处理单位
- [x] 运行 `make build` 验证

**修改文件**:
| 文件 | 修改内容 |
|------|---------|
| `admin/desc/fulfillment.api` | `AdjustAmount` 注释 `(in cents)` → `(in yuan)` |
| `admin/desc/shipping.api` | `FirstFee`/`AdditionalFee` 注释 `(cents)` → `(yuan)` |

**业务逻辑确认**:
- ✅ `fulfillment_orders/adjust_order_price_logic.go` - 使用 yuan
- ✅ `fulfillment_orders/helper.go` - 使用 yuan
- ✅ `shipping_templates/get_shipping_template_logic.go` - 使用 yuan（注意：代码中有一处过时注释 "int64 cents to string" 需后续清理）

**验证结果**:
- ✅ `make api` 成功
- ✅ `make build` 编译通过

**完成时间**: 2026-03-29

---

### Phase 1 + Phase 2 完成汇总

**总状态**: ✅ 全部完成

| 项目 | 状态 | 完成时间 |
|------|------|----------|
| Phase 1: ChangedAt 类型修复 | ✅ 完成 | 2026-03-29 |
| Phase 2: 金额注释修正 | ✅ 完成 | 2026-03-29 |
| 最终验证 (`make build`) | ✅ 通过 | 2026-03-29 |
| 文档更新 | ✅ 完成 | 2026-03-29 |

---

## 10. 完成摘要

### 修改文件汇总

| 文件 | 修改内容 |
|------|---------|
| `admin/desc/storefront.api` | `ChangedAt` 从 `int64` 改为 `string` (RFC3339) |
| `admin/desc/fulfillment.api` | `AdjustAmount` 注释改为 `(in yuan)` |
| `admin/desc/shipping.api` | `FirstFee`/`AdditionalFee` 注释改为 `(yuan)` |
| `admin/internal/application/storefront/service.go` | `CurrentThemeDTO.ChangedAt` 类型和填充逻辑更新 |
| `admin/internal/types/types.go` | 自动重新生成 |

### 遗留问题（已解决）

1. ~~**`admin/internal/logic/shipping_templates/get_shipping_template_logic.go` 第91行**~~
   - ~~过时注释：`// Helper to convert int64 cents to string`~~
   - ~~实际代码按 yuan 处理，注释需更新为 `(yuan)`~~
   - ✅ **已修复**: 注释更新为 `// formatAmount converts decimal.Decimal to string with 2 decimal places (in yuan)`
