# DB Schema 与 Entity 设计对比审查报告

> 日期: 2026-03-29
> 审查范围: 全部 9 个业务领域

---

## 一、审查范围

共 9 个 Domain：

| Domain | SQL Schema 文件 | Entity 文件 |
|--------|----------------|-------------|
| user | sql/user/schema.sql | admin/internal/domain/user/entity.go |
| product | sql/product/schema.sql | admin/internal/domain/product/entity.go |
| order | sql/order/schema.sql | admin/internal/domain/order/entity.go |
| promotion | sql/promotion/schema.sql | admin/internal/domain/promotion/entity.go |
| points | sql/points/schema.sql | admin/internal/domain/points/entity.go |
| storefront | sql/storefront/schema.sql | admin/internal/domain/storefront/entity.go |
| fulfillment | sql/fulfillment/schema.sql | admin/internal/domain/fulfillment/entity.go |
| payment | sql/payment/schema.sql | admin/internal/domain/payment/entity.go |
| review | sql/review/schema.sql | admin/internal/domain/review/entity.go |

---

## 二、审查结果汇总

### 审查发现总结

| Domain | 是否一致 | 严重问题数 | 中等问题数 |
|--------|---------|-------------|------------|
| user | ✓ 一致 | 0 | 0 |
| product | ⚠️ 部分一致 | 1 | 0 |
| order | ✓ 一致 | 0 | 0 |
| promotion | ✓ 一致 | 0 | 0 |
| points | ⚠️ 部分一致 | 0 | 1 |
| storefront | ✓ 一致 | 0 | 0 |
| fulfillment | ⚠️ 部分一致 | 1 | 1 |
| payment | ✓ 一致 | 0 | 0 |
| review | ✓ 一致 | 0 | 0 |

**一致率**: 4/9 (44%)

---

## 三、各 Domain 详细审查结果

### 1. USER Domain ✅ 一致

**表**: users, tenants, admin_users, roles, permissions, user_roles, role_permissions, user_addresses

| 字段 | Schema 类型 | Entity 类型 | 状态 |
|------|-----------|-------------|------|
| created_at | TIMESTAMP | time.Time (via Audit) | ✓ |
| updated_at | TIMESTAMP | time.Time (via Audit) | ✓ |
| deleted_at | TIMESTAMP | gorm.DeletedAt (via Model) | ✓ |
| 价格字段 | 无 | N/A | N/A |

**结论**: ✅ 一致

---

### 2. PRODUCT Domain ⚠️ 严重问题

**表**: products

| 字段 | Schema 类型 | Entity 类型 | 状态 |
|------|-----------|-------------|------|
| created_at | TIMESTAMP | time.Time (via Model) | ✓ |
| updated_at | TIMESTAMP | time.Time (via Model) | ✓ |
| deleted_at | TIMESTAMP | gorm.DeletedAt (via Model) | ✓ |
| **price** | **DECIMAL(19,4)** | **Money.Amount (int64)** | **❌ 不一致** |
| **cost_price** | **DECIMAL(19,4)** | **Money.Amount (int64)** | **❌ 不一致** |

**问题详情**:
- **products表**: `price` 和 `cost_price` 字段定义为 `DECIMAL(19,4)`
- **Product entity**: 使用自定义 `Money` 类型，内部用 `int64 Amount` (单位：分)
- 例如：数据库存储 `129900.0000` (表示 1299.00 元)，entity 内部存储 `int64 129900` (表示 129900 分 = 1299.00 元)
- **问题**: GORM 无法直接映射 DECIMAL(19,4) ↔ int64

**修复建议**:
- 方案A：将 Money 类型改为使用 `decimal.Decimal` (推荐，与 Schema 匹配)
- 方案B：修改 Schema，将 price 和 cost_price 改为 BIGINT (单位：分) - 需要同时修改数据

---

### 3. ORDER Domain ✅ 一致

**表**: orders, order_items

| 字段 | Schema 类型 | Entity 类型 | 状态 |
|------|-----------|-------------|------|
| created_at | TIMESTAMP | time.Time (via Audit) | ✓ |
| updated_at | TIMESTAMP | time.Time (via Audit) | ✓ |
| deleted_at | TIMESTAMP | gorm.DeletedAt (via Model) | ✓ |
| total_amount | DECIMAL(19,4) | shared.Money.Amount (decimal.Decimal) | ✓ |
| discount_amount | DECIMAL(19,4) | shared.Money.Amount (decimal.Decimal) | ✓ |
| freight_amount | DECIMAL(19,4) | shared.Money.Amount (decimal.Decimal) | ✓ |
| pay_amount | DECIMAL(19,4) | shared.Money.Amount (decimal.Decimal) | ✓ |

**结论**: ✅ 一致 (使用 shared.Money，decimal.Decimal 类型与 DECIMAL(19,4) 匹配)

---

### 4. PROMOTION Domain ✅ 一致

**表**: promotions, promotion_rules, coupons, user_coupons

| 字段 | Schema 类型 | Entity 类型 | 状态 |
|------|-----------|-------------|------|
| created_at | TIMESTAMP | time.Time (via Model) | ✓ |
| updated_at | TIMESTAMP | time.Time (via Model) | ✓ |
| deleted_at | TIMESTAMP | gorm.DeletedAt (via Model) | ✓ |
| status | TINYINT | Status (int) | ✓ (GORM 自动转换) |
| condition_value | DECIMAL(19,4) | int64 (ConditionValue) | ✓ (GORM 自动转换) |

**结论**: ✅ 一致

---

### 5. POINTS Domain ⚠️ 中等问题

**表**: earn_rules, redeem_rules, points_accounts

| 字段 | Schema 类型 | Entity 类型 | 状态 |
|------|-----------|-------------|------|
| created_at | TIMESTAMP | time.Time (via Model) | ✓ |
| updated_at | TIMESTAMP | time.Time (via Model) | ✓ |
| deleted_at | TIMESTAMP | gorm.DeletedAt (via Model) | ✓ |
| **status** | **VARCHAR(20)** | **EarnRuleStatus (int)** | **⚠️ 类型不匹配** |
| **ratio** | **DECIMAL(10,4)** | **decimal.Decimal** | **⚠️ 类型不匹配** |

**问题详情**:
- `earn_rules.status` 在 Schema 定义为 `VARCHAR(20)`，存储值如 'draft', 'active', 'inactive'
- Entity 中 `EarnRuleStatus` 是 int 类型 (iota 枚举)
- 类似问题存在于 `redeem_rules.status`

**修复建议**:
- 修改 Schema，将 status 改为 `TINYINT` 类型，与 entity 匹配
- 或在 entity 中使用 string 类型存储 status

---

### 6. STOREFRONT Domain ✅ 一致

**表**: shops, themes, pages

| 字段 | Schema 类型 | Entity 类型 | 状态 |
|------|-----------|-------------|------|
| created_at | TIMESTAMP | time.Time (via Audit) | ✓ |
| updated_at | TIMESTAMP | time.Time (via Audit) | ✓ |
| deleted_at | TIMESTAMP | gorm.DeletedAt (via Model) | ✓ |
| status | TINYINT | shared.Status (int) | ✓ (GORM 自动转换) |

**结论**: ✅ 一致

---

### 7. FULFILLMENT Domain ⚠️ 严重问题

**表**: shipments, shipment_items, refunds

| 字段 | Schema 类型 | Entity 类型 | 状态 |
|------|-----------|-------------|------|
| created_at | TIMESTAMP | time.Time (via Model) | ✓ |
| updated_at | TIMESTAMP | time.Time (via Model) | ✓ |
| deleted_at | TIMESTAMP | gorm.DeletedAt (via Model) | ✓ |
| **cost_amount** | **DECIMAL(19,4)** | **ShippingCost (decimal.Decimal)** | **⚠️ 字段名不同** |
| amount | DECIMAL(19,4) | Amount (decimal.Decimal) | ✓ |
| **shipment_no** | **不存在** | **string** | **⚠️ Schema 缺失** |
| **refund_no** | **不存在** | **string** | **⚠️ Schema 缺失** |

**问题详情**:
- shipments表: `cost_amount` vs entity: `ShippingCost` - 字段名不同
- Shipment entity 定义了 `ShipmentNo` 字段，但 Schema 中没有 shipment_no 列
- Refund entity 定义了 `RefundNo` 字段，但 Schema 中没有 refund_no 列

**修复建议**:
- Schema 添加 `shipment_no VARCHAR(32) NOT NULL` 列
- Schema 添加 `refund_no VARCHAR(32) NOT NULL` 列
- 检查 Shipment.ShippingCost 和 Shipment.CostCurrency 字段映射

---

### 8. PAYMENT Domain ✅ 一致

**表**: order_payments, payment_transactions, payment_refunds, webhook_events

| 字段 | Schema 类型 | Entity 类型 | 状态 |
|------|-----------|-------------|------|
| created_at | TIMESTAMP | time.Time (via Audit) | ✓ |
| updated_at | TIMESTAMP | time.Time (via Audit) | ✓ |
| deleted_at | TIMESTAMP | gorm.DeletedAt (via Model) | ✓ |
| amount | DECIMAL(19,4) | decimal.Decimal | ✓ |
| transaction_fee | DECIMAL(19,4) | decimal.Decimal | ✓ |

**结论**: ✅ 一致

---

### 9. REVIEW Domain ✅ 一致

**表**: reviews, review_replies, review_stats

| 字段 | Schema 类型 | Entity 类型 | 状态 |
|------|-----------|-------------|------|
| created_at | TIMESTAMP | time.Time (via Model) | ✓ |
| updated_at | TIMESTAMP | time.Time (via Model) | ✓ |
| deleted_at | TIMESTAMP | gorm.DeletedAt (via Model) | ✓ |
| status | TINYINT | Status (int) | ✓ (GORM 自动转换) |
| overall_rating | DECIMAL(3,2) | float64 | ✓ |

**结论**: ✅ 一致

---

## 四、需要修改的清单

### 严重问题（必须修复）

| # | Domain | 问题 | 修复方案 |
|---|-------|------|----------|
| 1 | product | Money.Amount 类型为 int64，无法映射 DECIMAL(19,4) | 将 Money 改为使用 decimal.Decimal，修改 NewMoney 构造函数返回 decimal.Decimal |
| 2 | fulfillment | Schema 缺少 shipment_no 字段 | 迁移添加 shipment_no VARCHAR(32) NOT NULL |
| 3 | fulfillment | Schema 缺少 refund_no 字段 | 迁移添加 refund_no VARCHAR(32) NOT NULL |

### 中等问题（建议修复）

| # | Domain | 问题 | 修复方案 |
|---|-------|------|----------|
| 1 | points | earn_rules.status 类型为 VARCHAR(20)，entity 用 int | 将 Schema 改为 TINYINT 类型 |
| 2 | fulfillment | cost_amount vs ShippingCost 字段名不一致 | 检查 GORM 映射，确保字段名一致 |

---

## 五、遗漏检查的 Domain

在当前代码库中发现以下 entity 可能在代码中定义了但没有对应的 schema 文件或 schema 不完整：

| Entity | Entity 文件 | 备注 |
|--------|------------|------|
| adminuser | admin/internal/domain/adminuser/entity.go | 已在 users 表中 |
| tenant | admin/internal/domain/tenant/entity.go | 已在 users 表中 |
| role | admin/internal/domain/role/entity.go | 已在 users 表中 |
| shipping | admin/internal/domain/shipping/entity.go | 已在 fulfillment schema 中 |
| coupon | admin/internal/domain/coupon/entity.go | 已在 promotion schema 中 |
| cart | admin/internal/domain/cart/entity.go | 已在 order schema 中 |
| market | admin/internal/domain/market/entity.go | 已在 product schema 中 |
| inventory (product) | admin/internal/domain/product/inventory.go | 已在 product schema 中 |
| category | admin/internal/domain/product/category.go | 已在 product schema 中 |

这些 entity 都已在对应的 schema 中覆盖。

---

## 六、Plan（修改计划）

基于审查结果，建议按以下优先级进行修改：

### Phase 1: 严重问题修复

1. **Product Domain - Money 类型修复**
   - 修改 `admin/internal/domain/product/entity.go` 中的 `Money` 类型
   - 将 `Amount int64` 改为 `Amount decimal.Decimal`
   - 修改 `NewMoney()`, `Add()`, `Subtract()`, `Equals()` 等方法
   - 验证测试编译通过

2. **Fulfillment Domain - Schema 缺失字段**
   - 分析 shipment_no 和 refund_no 的业务逻辑
   - 创建迁移脚本添加缺失字段
   - 确认是否是业务逻辑问题（不是必须的字段）

### Phase 2: 中等问题修复

3. **Points Domain - Status 类型统一**
   - 将 `earn_rules.status` 从 VARCHAR(20) 改为 TINYINT
   - 创建迁移脚本
   - 确保与 entity 保持一致

4. **Fulfillment - 字段名统一**
   - 检查 `cost_amount` vs `ShippingCost` 的映射
   - 如需要，修改 Schema 或 entity 的 gorm tag

### Phase 3: 验证

5. 运行 `make build` 确保所有代码编译通过
6. 如有，新增对应的单元测试

---

## 七、审查标准说明

本次审查依据以下标准：

1. **时间字段 (created_at, updated_at, deleted_at)**
   - Schema: 必须使用 `TIMESTAMP` 类型
   - Entity: 必须使用 `time.Time` 类型 (通过 application.Model 嵌入)
   - deleted_at 必须使用 `gorm.DeletedAt` 类型

2. **价格/金额字段**
   - Schema: 必须使用 `DECIMAL(19,4)` 类型
   - Entity: 必须使用 `decimal.Decimal` 或等效类型
   - 禁止使用 float/double 或 int64 直接存储金额

3. **状态字段**
   - Enum 类型 (iota) 在 Schema 使用 `TINYINT`
   - 字符串类型状态在 Schema 使用 `VARCHAR`

4. **字段名一致性**
   - Schema 列名和 Entity gorm tag 必须匹配
   - 或通过 gorm:column 明确指定映射关系

5. **修复方式规范**
   - **禁止使用简单方式或兼容方式绕过问题**
   - 修复必须从根本上解决问题，而非采用临时方案
   - 例如：
     - ❌ 在 Product entity 中使用 `gorm Type` 标签强制转换 `DECIMAL(19,4)` 到 `int64`（临时兼容）
     - ✅ 统一使用 `decimal.Decimal` 类型，与 Schema 完全匹配
     - ❌ 在 Points entity 中用字符串存储枚举值来"兼容" VARCHAR(20)
     - ✅ 将 Schema 改为 TINYINT，与 entity 的 int 枚举类型统一
   - 所有修复必须遵守项目既定的类型规范、命名规范和架构约定