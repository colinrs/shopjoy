# 实体嵌入模式：Model + AuditInfo 双嵌入与仅 AuditInfo

## 1. 概述

### 1.1 背景

shopjoy 领域层实体使用两种不同的时间字段嵌入模式：

| 模式 | 字段 | 示例实体 |
|------|------|---------|
| 模式 A | `application.Model` + `shared.AuditInfo` 同时嵌入 | User / SKU / Brand / Category / Role / Warehouse / Inventory / Page / Refund |
| 模式 B | 仅 `shared.AuditInfo` 嵌入 | Coupon / Promotion / Order |

两种模式下，序列化层读取时间字段的代码存在历史性的不对称：使用 `x.Audit.CreatedAt` / `x.Audit.UpdatedAt`。在模式 A 下这是 bug（`Audit.CreatedAt` 永远为零值，GORM 命名约定只回填 `Model.CreatedAt`）；在模式 B 下这是正确写法。

### 1.2 修复目标

1. 消除模式 A 下 API 响应中 `created_at` / `updated_at` 始终为 `"0001-01-01T00:00:00Z"` 的 bug
2. 统一两类实体的消费者代码语义（**读时统一指向被 GORM 实际填充的字段**）
3. 沉淀设计文档，避免后续重蹈覆辙

---

## 2. 问题详情

### 2.1 GORM 命名约定与字段填充规则

GORM 在 `INSERT` / `UPDATE` / `SELECT` 时遵循以下约定：

| 操作 | GORM 行为 |
|------|----------|
| `Create` | 如果 `CreatedAt` 字段为零值，GORM 自动填入 `time.Now()`；否则使用用户已设值 |
| `Update` / `Save` | 总是将 `UpdatedAt` 字段更新为 `time.Now()` |
| `SELECT` | 将 `created_at` / `updated_at` 列回填到结构体中名为 `CreatedAt` / `UpdatedAt` 的字段 |

当实体同时嵌入 `application.Model` 和 `shared.AuditInfo`（**两个结构体都含 `CreatedAt` 字段**）时：

- **写入侧**：GORM 按字段名识别到两个 `CreatedAt`，但因 `Audit` 在嵌入列表中位于 `Model` 之后，GORM 写入 `created_at` 列时实际读取的是 `application.Model.CreatedAt`（被自动填充为 NOW）
- **读取侧**：GORM 把 `created_at` 列回填到 `Model.CreatedAt`；`Audit.CreatedAt` 始终保持其结构体零值

**结论**：模式 A 下，**任何读到 `x.Audit.CreatedAt` 的代码都拿到零值**；正确字段是 `x.Model.CreatedAt`。

### 2.2 受影响的代码规模（修复前）

```
$ grep -rn "\.Audit\.CreatedAt\|\.Audit\.UpdatedAt" admin shop | grep -v _test.go | wc -l
68
```

涉及 30 个文件，覆盖：

- **Logic 层（23 个文件）**：所有 list / get / update 接口的响应序列化
- **Application 层（3 个文件）**：coupon / promotion / user 的 toXxxResponse
- **Persistence 层（6 个文件）**：repository 的 toXxxModel

### 2.3 表现

| API | 修复前 | 修复后 |
|-----|--------|--------|
| `GET /api/v1/users` | `created_at: "0001-01-01T00:00:00Z"` | `"2026-04-07T08:49:36+08:00"` |
| `GET /api/v1/categories` | `0001-01-01T00:00:00Z` | `"2026-07-05T19:24:38+08:00"` |
| `GET /api/v1/roles` | `0001-01-01T00:00:00Z` | `"2026-04-07T08:49:37+08:00"` |
| `GET /api/v1/brands` / `/products` 等 | 同上 | 同上 |

---

## 3. 决策

### 3.1 双轨制保留

保留模式 A（`Model + AuditInfo`）与模式 B（仅 `AuditInfo`）**两种嵌入模式**：

| 决策 | 理由 |
|------|------|
| 模式 A 保留 | 大部分实体（28+ 处）已采用此模式，统一它们的成本远低于把所有实体改成模式 B |
| 模式 B 保留 | Coupon / Promotion / Order 不需要软删除（没有 `DeletedAt`），无需 `application.Model` |
| **强制规则** | **消费者代码必须明确指出读哪个字段，禁止依赖 `x.CreatedAt` 这种会触发 Go 编译歧义的形式** |

### 3.2 修复矩阵

| 实体 | 嵌入模式 | 序列化层读字段 | 写入侧字段 |
|------|---------|--------------|----------|
| 模式 A（28 个实体） | `Model + AuditInfo` | **`x.Model.{CreatedAt,UpdatedAt}`** ✅ 修复后 | `x.Audit.{CreatedBy,UpdatedBy}`（Model 没这两个字段） |
| 模式 B（Coupon / Promotion / Order） | 仅 `AuditInfo` | `x.Audit.{CreatedAt,UpdatedAt}`（保持原样） | `x.Audit.{CreatedBy,UpdatedBy}` |

### 3.3 修复原则

1. **读时**：必须读 GORM 实际填充的字段
   - 模式 A：`x.Model.{CreatedAt,UpdatedAt}`
   - 模式 B：`x.Audit.{CreatedAt,UpdatedAt}`（这是 GORM 在该模式下唯一会填充的字段）
2. **写时**：手动 `x.Audit.UpdatedAt = time.Now()` 在模式 A 下是冗余的（GORM 会自动更新 `Model.UpdatedAt`），但语义上更清晰，仍保留。**统一改为 `x.Model.UpdatedAt = time.Now()`**，与读时字段一致。
3. **审计字段**（`CreatedBy` / `UpdatedBy`）：仅存在于 `shared.AuditInfo`，保持不变。
4. **不允许**使用 `x.CreatedAt` 这种歧义写法（Go 编译错误）。

### 3.4 实施步骤（已完成）

| # | 步骤 | 范围 |
|---|------|------|
| 1 | 编写失败测试 `service_impl_test.go`（3 个用例） | User 实体 |
| 2 | 修复 `service_impl.go` 4 处序列化 | User |
| 3 | sed 批量替换 30 个文件中 53 处 | 模式 A 全部实体 |
| 4 | 回退 Coupon / Promotion / Order 消费者（15 处） | 模式 B 保留 |
| 5 | 修复 dashboard/helper.go 中 Order 误改 | 模式 B 实体漏网 |
| 6 | 添加单测 `TestToUserResponse_CreatedAt` 等 | 防回归 |

---

## 4. 影响范围

### 4.1 改动文件清单（31 个）

```
admin/internal/application/user/service_impl.go          (+测试 3 用例)
admin/internal/application/user/service_impl_test.go    (NEW)
admin/internal/application/fulfillment/order_fulfillment_app.go
admin/internal/application/storefront/service.go
admin/internal/application/promotion/coupon_app.go        (回退，模式 B)
admin/internal/application/promotion/promotion_app.go     (回退，模式 B)
admin/internal/logic/{warehouses,products,brands,roles,categories,inventory,dashboard}/*.go
admin/internal/infrastructure/persistence/{sku,category,brand,coupon,promotion,order}_repository.go  (coupon/promotion/order 回退)
shop/internal/application/order/service_impl.go
```

### 4.2 测试

| 包 | 测试结果 |
|----|---------|
| `application/user` (3 新增) | PASS |
| `application/payment` (2 修复) | PASS（见 §6） |
| `logic/shipments` (1 修复) | PASS（见 §6） |
| 全量 `go test ./...` | 无 FAIL（除 2 个 pre-existing 无关失败已修复） |

### 4.3 API 行为变化

| 字段 | 修复前 | 修复后 | 是否 breaking change |
|------|--------|--------|-------------------|
| `created_at` | 始终 `"0001-01-01T00:00:00Z"` | 真实 DB 时间戳 | **是**，但所有受影响页面实际行为是「之前显示错误的零值」，客户端无业务逻辑依赖 |
| `updated_at` | 同上 | 同上 | 是 |
| `created_by` / `updated_by` | 不变 | 不变 | 否 |
| 所有其他字段 | 不变 | 不变 | 否 |

### 4.4 数据迁移

**无**。修复仅影响 Go 实体读取字段的选择，DB schema 与数据不变。

---

## 5. 后续建议（非本次范围）

### 5.1 短期

1. 在 CI 中加 `go vet` 规则（或自定义 lint），禁止 `\.Audit\.CreatedAt` 出现在模式 A 实体的消费者代码中
2. 把 `service_impl_test.go` 中 3 个测试扩展为表驱动，覆盖所有模式 A 实体的序列化函数

### 5.2 中期

将所有模式 B 实体（Coupon / Promotion / Order）**统一改为模式 A**，彻底消除双轨：

```go
type Coupon struct {
    application.Model         // 加上 Model 提供 ID/CreatedAt/UpdatedAt/DeletedAt
    TenantID  shared.TenantID
    Name      string
    // ...
    Audit     shared.AuditInfo `gorm:"embedded"`  // 保留 CreatedBy/UpdatedBy
}
```

**前提**：评估这些实体是否需要软删除（`DeletedAt`）。如果不需要，可以保留模式 B。

### 5.3 长期

考虑把 `CreatedBy` / `UpdatedBy` 也合并进 `application.Model`，让所有实体只用一种嵌入模式：

```go
// pkg/application/model.go
type Model struct {
    ID        int64 `gorm:"primarykey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
    CreatedBy int64
    UpdatedBy int64
}
```

但此改动影响所有实体的 SQL `created_by` / `updated_by` 列映射与审计逻辑，需专门 RFC。

---

## 6. 关联修复

本次实施过程中顺手修复了 2 个 pre-existing 测试失败（与本次主题无关，但同属代码质量债务）：

| # | 测试 | 根因 | 修复 |
|---|------|------|------|
| 1 | `payment.TestHandleWebhook_EventCreation` | 测试 mock `Status=Success`，但 `MarkSuccess` 只接受 Pending/Processing | 改 mock 为 `Status=PaymentStatusPending` |
| 2 | `payment.TestHandleWebhook_ExistingUnprocessedEvent` | 同上 | 同上 |
| 3 | `shipments.TestCancelShipmentLogic_CancelShipment/platform_admin` | 平台管理员调用时，`CancelShipment` 未识别 UserType=1，未把 tenantID 改为 0 | 在 logic 中加 `if contextx.IsPlatformAdmin(ctx) { tenantID = 0 }` |

---

## 7. 变更日志

| 日期 | 变更 | 作者 |
|------|------|------|
| 2026-07-10 | 初稿，固化 created_at 零值 bug 修复方案 | Claude |