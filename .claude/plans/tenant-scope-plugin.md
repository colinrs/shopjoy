# Tenant Scope Plugin 实施计划

> 基于技术方案 `docs/cross-cutting/tech-design/2026-07-10-tenant-scope-plugin-design.md`

## 决策确认

| 决策项 | 选择 |
|--------|------|
| 实施路径 | 直接方案 D（注册表 + AST 改写） |
| Context 统一 | 本次统一（合并 `pkg/contextx` + `pkg/tenant`，删除 `pkg/domain/shared` 死代码） |
| 适用端 | 只做 Admin 端 |
| Create 自动填充 | 一并实现（Query + Create 都自动化） |

## 代码现状

| 维度 | 数据 |
|------|------|
| Repository 手工隔离 | 20 个文件，117 处 `if tenantID != 0` |
| Tenant Context 实现 | 3 套（`pkg/contextx` 223 文件、`pkg/tenant` 21 文件、`pkg/domain/shared` 死代码） |
| 租户隔离表 | 48 张表有 `tenant_id` 列 |
| 系统级表 | 6 张（tenants, permissions, user_roles, role_permissions, carriers, refund_reasons） |
| 子表（继承父表租户） | 5 张（order_items, nav_items, shop_business_hours 等） |
| GORM 版本 | v1.31.1 |
| DB 初始化 | `pkg/infra/db.go`，已有 GormMetricsPlugin 插件模式 |

---

## 阶段一：统一 Tenant Context（1天）

### 目标
将 3 套 tenant context 合并为 1 套，消除双写。

### 步骤

#### 1.1 增强 `pkg/contextx/context.go` 为唯一 tenant context 来源

当前 `pkg/contextx` 已经是 Admin 端的事实标准（223 文件使用），但返回 `int64`。需要：
- 新增 `GetTenantIDValueObject(ctx) (shared.TenantID, bool)` —— 返回 `shared.TenantID` 类型，供 Plugin 和 infrastructure 层使用
- 保持现有 `GetTenantID(ctx) (int64, bool)` 不变，不破坏现有调用

```go
// pkg/contextx/context.go 新增
func GetTenantIDValueObject(ctx context.Context) (shared.TenantID, bool) {
    id, ok := GetTenantID(ctx)
    return shared.TenantID(id), ok
}
```

#### 1.2 迁移 `pkg/tenant` 的调用者到 `pkg/contextx`

`pkg/tenant` 被 21 个文件使用，主要在：
- `pkg/application/base.go` —— `GetTenantID()` / `MustGetTenantID()` 包装函数
- `pkg/infra/` 下的 infrastructure 代码
- `shop/` 端代码（本次不动）

对 Admin 端的调用者：
- `pkg/application/base.go`：改为调用 `contextx.GetTenantIDValueObject()`
- `admin/internal/infrastructure/persistence/` 下的 repository：改为从 ctx 直接获取，不再需要参数传入

#### 1.3 删除 `pkg/domain/shared/context.go` 中的死代码

`GetTenantIDFromContext` / `SetTenantIDInContext` 有 0 个调用者，直接删除。

#### 1.4 清理 `admin/internal/middleware/auth_middleware.go` 双写

```go
// 当前（双写）
ctx = contextx.SetTenantID(ctx, tenantID)
ctx = tenant.WithContext(ctx, shared.TenantID(tenantID))

// 改为（单写）
ctx = contextx.SetTenantID(ctx, tenantID)
// Plugin 通过 contextx.GetTenantIDValueObject() 读取
```

### 影响范围
- `pkg/contextx/context.go` —— 新增 1 个函数
- `pkg/application/base.go` —— 改 2 个函数
- `admin/internal/middleware/auth_middleware.go` —— 删 1 行
- `pkg/domain/shared/context.go` —— 删除死代码函数（文件保留，其他内容不动）

---

## 阶段二：实现 TenantScopePlugin（2-3天）

### 步骤

#### 2.1 创建表名注册表 `pkg/infra/tenant_tables.go`

```go
package infra

// tenantScopedTables 记录所有需要租户隔离的表
// 新增表时必须在此注册，配合 linter 检查遗漏
var tenantScopedTables = map[string]bool{
    // User Domain
    "users": true, "admin_users": true, "roles": true,
    "user_addresses": true, "user_operation_logs": true,
    // Product Domain
    "categories": true, "brands": true, "products": true, "skus": true,
    "product_markets": true, "product_localizations": true,
    "category_markets": true, "brand_markets": true,
    "warehouses": true, "warehouse_inventories": true,
    "inventory_logs": true, "markets": true,
    // Order Domain
    "orders": true, "carts": true, "cart_items": true,
    // Payment Domain
    "payments": true, "order_payments": true,
    "payment_transactions": true, "payment_refunds": true,
    "webhook_events": true,
    // Fulfillment Domain
    "shipments": true, "shipment_items": true, "refunds": true,
    "shipping_templates": true, "shipping_zones": true,
    "shipping_template_mappings": true,
    // Promotion Domain
    "promotions": true, "promotion_rules": true,
    "promotion_usage": true, "coupons": true, "user_coupons": true,
    // Points Domain
    "earn_rules": true, "redeem_rules": true,
    "points_accounts": true, "points_transactions": true,
    "points_redemptions": true,
    // Review Domain
    "reviews": true, "review_replies": true, "review_stats": true,
    // Storefront Domain
    "shops": true, "themes": true, "pages": true,
    "navigations": true, "decorations": true,
    "page_versions": true, "seo_configs": true,
    "theme_audit_logs": true,
    // Shop Domain
    "shop_settings": true,
}

func IsTenantScopedTable(tableName string) bool {
    return tenantScopedTables[tableName]
}
```

#### 2.2 创建 TenantScopePlugin `pkg/infra/tenant_scope_plugin.go`

核心结构：
- `Name()` 返回 `"tenant_scope"`
- `Initialize(db *gorm.DB)` 注册 4 个 callback：Query、Delete、Update、Create
- `onQuery` —— 对 Query/Delete/Update 注入 `WHERE tenant_id = ?`
- `onCreate` —— 对 Create 自动 `SetColumn("tenant_id", ?)`
- `onRaw` —— 对 Raw SQL 调用 AST 改写

```go
type TenantScopePlugin struct{}

func (p *TenantScopePlugin) Name() string { return "tenant_scope" }

func (p *TenantScopePlugin) Initialize(db *gorm.DB) error {
    db.Callback().Query().Before("gorm:query").Register("tenant_scope:before_query", p.onQuery)
    db.Callback().Delete().Before("gorm:delete").Register("tenant_scope:before_delete", p.onQuery)
    db.Callback().Update().Before("gorm:update").Register("tenant_scope:before_update", p.onQuery)
    db.Callback().Create().Before("gorm:create").Register("tenant_scope:before_create", p.onCreate)
    db.Callback().Raw().Before("gorm:raw").Register("tenant_scope:before_raw", p.onRaw)
    return nil
}
```

#### 2.3 实现 SQL AST 改写 `pkg/infra/sql_rewrite.go`

引入依赖：`github.com/pingcap/tidb/pkg/parser`

核心逻辑：
1. 用 TiDB parser 解析 SQL 为 AST
2. 遍历 AST 中的 `SELECT`/`UPDATE`/`DELETE` 语句
3. 对 FROM/JOIN 中引用的表，检查是否在注册表中
4. 对注册表中的表注入 `WHERE tenant_id = ?` / `AND tenant_id = ?`
5. 对 INSERT 语句，检查目标表是否在注册表中，自动添加 `tenant_id` 列
6. Restore 回 SQL 字符串

#### 2.4 实现 SkipTenantScope 跳过机制 `pkg/infra/tenant_scope_skip.go`

```go
type skipTenantScopeKeyType struct{}
var skipTenantScopeKey = skipTenantScopeKeyType{}

func SkipTenantScope(ctx context.Context) context.Context {
    return context.WithValue(ctx, skipTenantScopeKey, true)
}

func IsSkipTenantScope(ctx context.Context) bool {
    v, _ := ctx.Value(skipTenantScopeKey).(bool)
    return v
}
```

#### 2.5 注册 Plugin 到 DB 初始化 `pkg/infra/db.go`

```go
func Database(mysqlConfig *DBConfig) (*gorm.DB, error) {
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Plugins: map[string]gorm.Plugin{
            metricsName: NewGormMetricsPlugin(WithDataBaseName(mysqlConfig.Database)),
            "tenant_scope": &TenantScopePlugin{},  // 新增
        },
    })
    // ...
}
```

### 新增文件清单

| 文件 | 用途 | 预估行数 |
|------|------|---------|
| `pkg/infra/tenant_tables.go` | 表名注册表 | ~80 |
| `pkg/infra/tenant_scope_plugin.go` | GORM Plugin 主体 | ~150 |
| `pkg/infra/sql_rewrite.go` | SQL AST 改写 | ~200 |
| `pkg/infra/tenant_scope_skip.go` | Skip 机制 | ~20 |
| `pkg/infra/tenant_scope_plugin_test.go` | 单元测试 | ~300 |
| `pkg/infra/sql_rewrite_test.go` | AST 改写测试 | ~200 |

---

## 阶段三：集成测试（1天）

### 测试场景

| 场景 | 验证内容 |
|------|---------|
| ORM Query 自动注入 | `db.Find(&products)` 生成的 SQL 包含 `WHERE tenant_id = ?` |
| ORM Create 自动填充 | `db.Create(&product)` 自动设置 `tenant_id` 列 |
| Raw SQL 自动改写 | `db.Raw("SELECT * FROM products")` 被改写为带 `tenant_id` 条件 |
| 子查询覆盖 | `WHERE id IN (SELECT ...)` 中的子查询也被处理 |
| JOIN 多表 | `orders JOIN users` 两张表都注入 `tenant_id` 条件 |
| 平台管理员跳过 | `tenantID == 0` 时不注入任何条件 |
| SkipTenantScope 显式跳过 | `SkipTenantScope(ctx)` 后不注入条件 |
| 非租户表不注入 | 查询 `tenants` 表时不注入 `tenant_id` |

### 测试方法
- 使用 SQLite 内存数据库（go.mod 已有 `gorm.io/driver/sqlite`）
- 对比改造前后生成的 SQL 是否语义一致

---

## 阶段四：Repository 改造（3-5天）

### 改造策略

**逐步迁移，每个 Repository 独立改造，改造后运行现有测试确认行为不变。**

改造模式（每个方法）：
1. 去掉 `tenantID shared.TenantID` 参数
2. 去掉 `if tenantID != 0 { query = query.Where(...) }` 逻辑
3. 去掉调用方传入的 `tenantID` 参数

### 改造顺序（按文件影响大小）

| 优先级 | 文件 | 处数 | 说明 |
|--------|------|------|------|
| P0 | `sku_repository.go` | 3 | 最少，先验证流程 |
| P0 | `user_operation_log_repository.go` | 1 | 最少 |
| P1 | `payment_repository.go` | 4 | |
| P1 | `payment_transaction_repository.go` | 4 | |
| P1 | `shipment_item_repository.go` | 4 | |
| P1 | `product_repository.go` | 6 | |
| P1 | `order_repository.go` | 6 | |
| P2 | `coupon_repository.go` | 5 | |
| P2 | `promotion_repository.go` | 5 | |
| P2 | `payment_refund_repository.go` | 5 | |
| P2 | `review_repository.go` | 5 | |
| P2 | `shipment_repository.go` | 7 | |
| P2 | `user_repository.go` | 7 | |
| P3 | `refund_repository.go` | 8 | |
| P3 | `points/earn_rule_repo.go` | 7 | |
| P3 | `points/redeem_rule_repo.go` | 8 | |
| P3 | `points_repository.go` | 29 | 最多，最后改造 |

### 同步改造 Logic 层

Repository 去掉 `tenantID` 参数后，调用方（`admin/internal/logic/` 和 `admin/internal/application/`）也需要同步去掉传参。这些文件的改动是机械性的：删除 `tenantID` 变量声明和传参。

---

## 阶段五：清理与文档（0.5天）

| 任务 | 说明 |
|------|------|
| 删除 `pkg/tenant` 中不再需要的 context 函数 | 保留 `Tenant` 实体和 `Repository` 接口，只删 context 相关 |
| 更新 AGENTS.md | 添加 Repository 规范：不再需要手动传 `tenantID` |
| 添加 linter 规则 | 检查 SQL schema 中有 `tenant_id` 的表是否都在注册表中 |
| 更新 `.claude/rules/` | 更新 Repository 相关规则 |

---

## 风险与缓解

| 风险 | 缓解措施 |
|------|---------|
| AST 改写引入 SQL 注入 | 只使用 AST 节点的值，不拼接用户输入；用参数化查询 |
| TiDB parser 不兼容某些 MySQL 语法 | 降级为不改写，记录日志；大部分 ORM 生成的 SQL 格式固定 |
| 性能开销（每次 SQL 解析 AST） | SQL 解析通常 < 1ms；可加开关在非生产环境关闭 |
| 改造遗漏导致数据泄漏 | 先上线 Plugin（自动注入），再逐步删除手工代码；两层防线共存 |
| 子表（如 order_items）无 tenant_id | 子表通过 JOIN 父表时，父表的 tenant_id 条件已保证隔离 |

---

## 依赖变更

`go.mod` 新增：
```
github.com/pingcap/tidb/pkg/parser  # SQL AST 解析
```

---

## 时间线

| 阶段 | 内容 | 工期 |
|------|------|------|
| 一 | 统一 Tenant Context | 1 天 |
| 二 | 实现 TenantScopePlugin + AST 改写 | 2-3 天 |
| 三 | 集成测试 | 1 天 |
| 四 | Repository 改造（20 文件，117 处） | 3-5 天 |
| 五 | 清理与文档 | 0.5 天 |
| **合计** | | **7.5-10.5 天** |
