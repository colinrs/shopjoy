# GORM Tenant Scope Plugin 技术方案

> 状态：草案 | 日期：2026-07-10

## 1. 背景与问题

当前项目中，租户隔离逻辑 100% 手工实现。每个 Repository 方法都需要：

1. 接收 `tenantID shared.TenantID` 参数
2. 判断 `if tenantID != 0` 后手动追加 `WHERE tenant_id = ?`

**现状统计：**

| 文件 | 出现次数 |
|------|---------|
| `product_repository.go` | 6 |
| `shipment_repository.go` | 7 |
| `payment_repository.go` | 4 |
| `payment_transaction_repository.go` | 4 |
| `sku_repository.go` | 4 |
| `user_repository.go` | 2 |
| **合计** | **27+** |

**问题：**
- 大量重复代码，违反 DRY
- 容易遗漏 — 某个查询忘了加判断就会造成数据泄漏
- Repository 方法签名冗长（每个方法都要传 `tenantID`）

## 2. 目标

1. **消除重复**：Repository 方法不再需要手动写 `if tenantID != 0` 逻辑
2. **防止遗漏**：新增表/查询时，不容易忘记租户过滤
3. **平台管理员兼容**：`tenantID == 0` 时自动跳过过滤（保持现有语义）
4. **平滑迁移**：可以逐步迁移，新旧方式可以共存
5. **覆盖全面**：ORM 查询、Raw SQL、子查询、JOIN 都能处理

---

## 3. 方案总览

| 方案 | 核心思路 | 复杂度 | 覆盖率 |
|------|---------|--------|--------|
| **A. 表名白名单（声明式）** | 维护表名 map，GORM callback 按表名注入 | 低 | 80% |
| **B. 查询实例显式标记（命令式）** | 每次查询前调用 `db.WithTenantScope(ctx)` | 低 | 100%（靠人） |
| **C. Model 接口标记** | Model 实现 `TenantScoped` 接口，callback 运行时检查 | 中 | 70% |
| **D. 表名注册表 + SQL AST 改写** | 注册表 + SQL parser 改写，参考分表中间件 | 高 | 99% |

---

## 4. 方案 A：表名白名单（声明式）

### 4.1 思路

维护一个 `map[string]bool`，记录哪些表需要租户隔离。GORM callback 中检查当前操作的表名是否在 map 中，如果是则注入 `WHERE tenant_id = ?`。

### 4.2 Demo

```go
// pkg/infra/tenant_tables.go
var tenantScopedTables = map[string]bool{
    "products":    true,
    "orders":      true,
    "payments":    true,
    "tenants":     false, // 系统级表，不需要
}

// pkg/infra/tenant_scope_plugin.go
type TenantScopePlugin struct{}

func (p *TenantScopePlugin) Initialize(db *gorm.DB) error {
    db.Callback().Query().Before("gorm:query").Register("tenant_scope:query", p.onQuery)
    db.Callback().Delete().Before("gorm:delete").Register("tenant_scope:delete", p.onQuery)
    db.Callback().Update().Before("gorm:update").Register("tenant_scope:update", p.onQuery)
    return nil
}

func (p *TenantScopePlugin) onQuery(db *gorm.DB) {
    tenantID, ok := tenant.FromContext(db.Statement.Context)
    if !ok || tenantID == 0 {
        return
    }

    tableName := db.Statement.Table
    if tableName == "" && db.Statement.Model != nil {
        tableName = db.Statement.Schema.Table
    }

    if !tenantScopedTables[tableName] {
        return
    }

    db.Statement.AddClause(clause.Where{
        Exprs: []clause.Expression{
            clause.Eq{Column: "tenant_id", Value: tenantID.Int64()},
        },
    })
}
```

### 4.3 优点

- **实现简单**：一个 map + 一个 callback，代码量少
- **集中管理**：所有需要租户隔离的表在一个文件里，一目了然
- **无侵入**：Model 不需要任何改造

### 4.4 缺陷

| 缺陷 | 严重程度 | 说明 |
|------|---------|------|
| Raw SQL 绕过 | **高** | `db.Raw("SELECT * FROM products")` 时 `Statement.Table` 可能为空 |
| 子查询不覆盖 | 中 | `WHERE id IN (SELECT ...)` 中的子查询不会被处理 |
| 需要维护 map | 低 | 新增表时需要记得加到 map 中（可 linter 检查） |

### 4.5 Repository 改造示例

```go
// 改造前
func (r *productRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*product.Product, error) {
    query := db.WithContext(ctx).Where("deleted_at IS NULL")
    if tenantID != 0 {
        query = query.Where("tenant_id = ?", tenantID.Int64())
    }
    var model productModel
    err := query.First(&model, id).Error
    // ...
}

// 改造后
func (r *productRepo) FindByID(ctx context.Context, db *gorm.DB, id int64) (*product.Product, error) {
    var model productModel
    err := db.WithContext(ctx).Where("deleted_at IS NULL").First(&model, id).Error
    // ...
}
```

---

## 5. 方案 B：查询实例显式标记（命令式）

### 5.1 思路

提供一个 `WithTenantScope(ctx)` 方法，每个需要租户隔离的查询前显式调用。不依赖表名或 Model，完全由开发者控制。

### 5.2 Demo

```go
// pkg/infra/tenant_scope.go
type tenantScopeKeyType struct{}
var tenantScopeKey = tenantScopeKeyType{}

type TenantScope struct {
    TenantID shared.TenantID
}

func WithTenantScope(ctx context.Context, tenantID shared.TenantID) context.Context {
    return context.WithValue(ctx, tenantScopeKey, &TenantScope{TenantID: tenantID})
}

func GetTenantScope(ctx context.Context) (*TenantScope, bool) {
    ts, ok := ctx.Value(tenantScopeKey).(*TenantScope)
    return ts, ok
}

// pkg/infra/tenant_scope_plugin.go
func (p *TenantScopePlugin) onQuery(db *gorm.DB) {
    ts, ok := GetTenantScope(db.Statement.Context)
    if !ok || ts.TenantID == 0 {
        return
    }

    db.Statement.AddClause(clause.Where{
        Exprs: []clause.Expression{
            clause.Eq{Column: "tenant_id", Value: ts.TenantID.Int64()},
        },
    })
}
```

### 5.3 Repository 改造示例

```go
// 改造前
func (r *productRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*product.Product, error) {
    query := db.WithContext(ctx).Where("deleted_at IS NULL")
    if tenantID != 0 {
        query = query.Where("tenant_id = ?", tenantID.Int64())
    }
    var model productModel
    err := query.First(&model, id).Error
    // ...
}

// 改造后 — 去掉 if 判断，但需调用 WithTenantScope
func (r *productRepo) FindByID(ctx context.Context, db *gorm.DB, id int64) (*product.Product, error) {
    var model productModel
    err := db.WithContext(WithTenantScope(ctx, tenant.MustFromContext(ctx))).
        Where("deleted_at IS NULL").
        First(&model, id).Error
    // ...
}
```

### 5.4 优点

- **精确控制**：开发者明确知道哪些查询有租户过滤
- **无表名依赖**：不关心表名，只关心「这个查询需要租户隔离」
- **实现简单**：一个 context key + 一个 callback

### 5.5 缺陷

| 缺陷 | 严重程度 | 说明 |
|------|---------|------|
| 仍然容易遗漏 | **高** | 忘了调用 `WithTenantScope` 就没有过滤 — 和忘写 `if` 本质一样 |
| 代码量没减少 | 中 | 每个查询前都要调用，只是从 `if tenantID != 0` 变成了 `WithTenantScope` |
| Raw SQL 绕过 | 中 | `db.Raw()` 不经过 GORM callback |

---

## 6. 方案 C：Model 接口标记

### 6.1 思路

定义 `TenantScoped` 接口，需要租户隔离的 Model 实现该接口。GORM callback 中通过 type assertion 检查 Model 是否实现了接口。

### 6.2 Demo

```go
// pkg/infra/tenant_scope.go
type TenantScoped interface {
    IsTenantScoped() bool
}

// Model 改造
type productModel struct {
    application.Model
    TenantID int64  `gorm:"column:tenant_id;not null;index"`
    Name     string `gorm:"column:name;size:200;not null"`
}
func (productModel) IsTenantScoped() bool { return true }

// Plugin
func (p *TenantScopePlugin) onQuery(db *gorm.DB) {
    tenantID, ok := tenant.FromContext(db.Statement.Context)
    if !ok || tenantID == 0 {
        return
    }

    // 运行时检查 Model 是否实现了 TenantScoped 接口
    if db.Statement.Model == nil {
        return
    }
    ts, ok := db.Statement.Model.(TenantScoped)
    if !ok || !ts.IsTenantScoped() {
        return
    }

    db.Statement.AddClause(clause.Where{
        Exprs: []clause.Expression{
            clause.Eq{Column: "tenant_id", Value: tenantID.Int64()},
        },
    })
}
```

### 6.3 优点

- **与 Model 绑定**：租户隔离的声明和实体放在一起，符合 DDD 思想
- **无表名维护**：不需要额外的 map 文件
- **ORM 查询覆盖**：GORM 的 Query/Delete/Update callback 都能处理

### 6.4 缺陷

| 缺陷 | 严重程度 | 说明 |
|------|---------|------|
| 运行时发现，非编译期 | **高** | 忘加 `IsTenantScoped()` 不报编译错，静默不过滤 — 和"遗漏 `if`"本质一样 |
| 隐式行为，调试困难 | 中 | 开发者写 `db.Find(&products)` 看不到任何 tenant 过滤痕迹 |
| Raw SQL 完全绕过 | **高** | `db.Raw("SELECT * FROM products")` 不经过 Model 检查 |
| 子查询无法覆盖 | 中 | `WHERE id IN (SELECT ...)` 子查询不被拦截 |
| 接口散落各处 | 低 | 想知道哪些表有租户隔离，需要逐个文件检查 |

---

## 7. 方案 D：表名注册表 + SQL AST 改写（推荐）

### 7.1 思路

参考分表中间件（[ShardingSphere](https://shardingsphere.apache.org/)、[Vitess](https://vitess.io/)）的做法：

1. 维护一个 **表名注册表**（集中的 map）
2. 在 GORM callback 中拦截 SQL
3. 通过 **SQL AST 解析** 找到所有表引用
4. 对注册表中的表自动注入 `WHERE tenant_id = ?`
5. 子查询、JOIN、Raw SQL 全覆盖

```
SQL → GORM callback 拦截 → AST 解析 → 遍历表引用
  → 注册表中的表？ → 注入 WHERE tenant_id = ?
  → 不在注册表？ → 跳过
```

### 7.2 Demo

**表名注册表：**
```go
// pkg/infra/tenant_tables.go
var tenantScopedTables = map[string]bool{
    "products":           true,
    "product_markets":    true,
    "skus":               true,
    "orders":             true,
    "order_items":        true,
    "payments":           true,
    "shipments":          true,
    "users":              true,
    "admin_users":        true,
    "roles":              true,
    "promotions":         true,
    "coupons":            true,
    "carts":              true,
    "reviews":            true,
    "tenants":            false, // 系统级表
    "markets":            false,
}

func IsTenantScopedTable(tableName string) bool {
    return tenantScopedTables[tableName]
}
```

**GORM Plugin：**
```go
// pkg/infra/tenant_scope_plugin.go
type TenantScopePlugin struct{}

func (p *TenantScopePlugin) Initialize(db *gorm.DB) error {
    db.Callback().Query().Before("gorm:query").Register("tenant_scope:query", p.onQuery)
    db.Callback().Delete().Before("gorm:delete").Register("tenant_scope:delete", p.onQuery)
    db.Callback().Update().Before("gorm:update").Register("tenant_scope:update", p.onQuery)
    db.Callback().Raw().Before("gorm:raw").Register("tenant_scope:raw", p.onRaw)
    return nil
}

func (p *TenantScopePlugin) onQuery(db *gorm.DB) {
    tenantID, ok := tenant.FromContext(db.Statement.Context)
    if !ok || tenantID == 0 {
        return
    }

    tableName := db.Statement.Table
    if tableName == "" && db.Statement.Model != nil {
        tableName = db.Statement.Schema.Table
    }

    if !IsTenantScopedTable(tableName) {
        return
    }

    db.Statement.AddClause(clause.Where{
        Exprs: []clause.Expression{
            clause.Eq{Column: "tenant_id", Value: tenantID.Int64()},
        },
    })
}

// onRaw 处理 db.Raw() — 通过 SQL 改写
func (p *TenantScopePlugin) onRaw(db *gorm.DB) {
    tenantID, ok := tenant.FromContext(db.Statement.Context)
    if !ok || tenantID == 0 {
        return
    }

    sql := db.Statement.SQL.String()
    rewritten, err := RewriteSQLWithTenantScope(sql, tenantID.Int64())
    if err != nil {
        // 解析失败，记录日志，不改写（安全降级）
        return
    }
    db.Statement.SQL.Reset()
    db.Statement.SQL.WriteString(rewritten)
}
```

**SQL AST 改写（使用 [pingcap/tidb/pkg/parser](https://github.com/pingcap/tidb/tree/master/pkg/parser)）：**
```go
// pkg/infra/sql_rewrite.go
import (
    "github.com/pingcap/tidb/pkg/parser"
    "github.com/pingcap/tidb/pkg/parser/ast"
    "github.com/pingcap/tidb/pkg/parser/format"
    _ "github.com/pingcap/tidb/pkg/parser/test_driver"
)

func RewriteSQLWithTenantScope(sql string, tenantID int64) (string, error) {
    p := parser.New()
    stmts, _, err := p.ParseSQL(sql)
    if err != nil {
        return "", err
    }

    for _, stmt := range stmts {
        ast.Walk(&tenantScopeRewriter{tenantID: tenantID}, stmt)
    }

    var buf strings.Builder
    restoreFlags := format.RestoreStringSingleQuotes | format.RestoreNameBackQuotes
    for _, stmt := range stmts {
        if err := stmt.Restore(format.NewRestoreCtx(restoreFlags, &buf)); err != nil {
            return "", err
        }
        buf.WriteString(";")
    }
    return buf.String(), nil
}

type tenantScopeRewriter struct {
    tenantID int64
}

func (r *tenantScopeRewriter) Enter(in ast.Node) (ast.Node, bool) {
    switch node := in.(type) {
    case *ast.SelectStmt:
        r.injectIntoSelect(node)
    case *ast.DeleteStmt:
        r.injectIntoDelete(node)
    case *ast.UpdateStmt:
        r.injectIntoUpdate(node)
    }
    return in, false
}

func (r *tenantScopeRewriter) Leave(in ast.Node) (ast.Node, bool) {
    return in, true
}

func (r *tenantScopeRewriter) injectIntoSelect(sel *ast.SelectStmt) {
    if sel.From == nil || sel.From.TableRefs == nil {
        return
    }
    // 遍历 FROM/JOIN 中的所有表，对注册表中的表注入 WHERE tenant_id = ?
    // ... 具体实现省略
}
```

### 7.3 优点

- **覆盖全面**：ORM 查询、Raw SQL、子查询、JOIN 统一处理
- **集中管理**：注册表一个文件看全貌
- **可审计**：注册表就是「哪些表需要租户隔离」的权威清单
- **可扩展**：注册表 + linter 可以自动检查遗漏

### 7.4 缺陷

| 缺陷 | 严重程度 | 说明 |
|------|---------|------|
| 实现复杂度高 | 中 | 需要引入 SQL parser 依赖，AST 改写逻辑较复杂 |
| 性能开销 | 低 | 每次 SQL 执行前都要解析 AST（但 SQL 解析通常 < 1ms） |
| 需要维护注册表 | 低 | 新增表时需要加到注册表（可 linter 检查） |

### 7.5 Repository 改造示例

```go
// 改造前
func (r *productRepo) FindByID(ctx context.Context, db *gorm.DB, tenantID shared.TenantID, id int64) (*product.Product, error) {
    query := db.WithContext(ctx).Where("deleted_at IS NULL")
    if tenantID != 0 {
        query = query.Where("tenant_id = ?", tenantID.Int64())
    }
    var model productModel
    err := query.First(&model, id).Error
    // ...
}

// 改造后
func (r *productRepo) FindByID(ctx context.Context, db *gorm.DB, id int64) (*product.Product, error) {
    var model productModel
    err := db.WithContext(ctx).Where("deleted_at IS NULL").First(&model, id).Error
    // ...
}
```

---

## 8. 全方位对比

### 8.1 功能对比

| 能力 | 现状（手工） | A. 表名白名单 | B. 显式标记 | C. Model 接口 | D. 注册表+AST |
|------|------------|--------------|------------|--------------|--------------|
| ORM 查询 | ✅ 手动 | ✅ 自动 | ⚠️ 靠人调用 | ✅ 自动 | ✅ 自动 |
| Raw SQL | ✅ 手动 | ❌ 绕过 | ⚠️ 靠人调用 | ❌ 绕过 | ✅ 自动 |
| 子查询 | ✅ 手动 | ❌ 不覆盖 | ⚠️ 靠人调用 | ❌ 不覆盖 | ✅ 自动 |
| JOIN 多表 | ✅ 手动 | ⚠️ 只主表 | ⚠️ 靠人调用 | ⚠️ 只主表 | ✅ 全表 |
| Create 自动填 | ❌ 手动 | ❌ 手动 | ❌ 手动 | ❌ 手动 | ✅ 自动 |

### 8.2 安全性对比

| 安全维度 | 现状（手工） | A. 表名白名单 | B. 显式标记 | C. Model 接口 | D. 注册表+AST |
|---------|------------|--------------|------------|--------------|--------------|
| 遗漏风险 | **高** — 忘写 `if` | **低** — 集中注册 | **高** — 忘调用 | **高** — 忘实现接口 | **低** — 集中注册 |
| 可审计性 | 低 — 散落各处 | **高** — 一个文件 | 低 — 散落各处 | 中 — 接口分散 | **高** — 一个文件 |
| 调试难度 | 低 — 显式可见 | 中 — 隐式注入 | 低 — 显式调用 | **高** — 隐式行为 | 中 — 隐式注入 |
| Raw SQL 泄漏 | 无 — 手动处理 | **有** — 绕过 | 无 — 手动处理 | **有** — 绕过 | **无** — AST 改写 |

### 8.3 开发体验对比

| 体验维度 | 现状（手工） | A. 表名白名单 | B. 显式标记 | C. Model 接口 | D. 注册表+AST |
|---------|------------|--------------|------------|--------------|--------------|
| 新增查询 | 写 `if` + `Where` | 不改代码 | 写 `WithTenantScope` | 不改代码 | 不改代码 |
| 新增表 | 写 `if` + `Where` | 注册表加一行 | 写 `WithTenantScope` | Model 加接口 | 注册表加一行 |
| 方法签名 | 多 `tenantID` 参数 | 少参数 | 少参数 | 少参数 | 少参数 |
| 代码重复 | **27+ 处** | **0 处** | **27+ 处**（换形式） | **27+ 处**（换形式） | **0 处** |

### 8.4 实现复杂度对比

| 维度 | A. 表名白名单 | B. 显式标记 | C. Model 接口 | D. 注册表+AST |
|------|--------------|------------|--------------|--------------|
| 代码量 | ~100 行 | ~50 行 | ~150 行 | ~400 行 |
| 外部依赖 | 无 | 无 | 无 | `pingcap/tidb/pkg/parser` |
| 实现难度 | 低 | 低 | 中 | 高 |
| 测试复杂度 | 低 | 低 | 中 | 高 |

---

## 9. 推荐方案与决策

### 9.1 推荐：方案 D（注册表 + AST 改写）

**理由：**

1. **安全性最高** — 唯一能覆盖 Raw SQL 和子查询的方案
2. **重复代码归零** — 27+ 处 `if tenantID != 0` 全部消除
3. **注册表可审计** — 一个文件看全貌，配合 linter 防遗漏
4. **参考成熟实践** — ShardingSphere、Vitess 都用这个模式

### 9.2 渐进式路径

如果方案 D 实现成本太高，可以分步走：

```
第一步：方案 A（表名白名单）
  → 覆盖 ORM 查询，解决 80% 问题
  → 实现快（1-2天）

第二步：升级到方案 D（加 AST 改写）
  → 覆盖 Raw SQL、子查询
  → 引入 pingcap/tidb/pkg/parser
  → 在方案 A 基础上增加 ~200 行代码
```

### 9.3 不推荐的方案

| 方案 | 不推荐理由 |
|------|-----------|
| B. 显式标记 | 遗漏风险和现状一样，代码量也没减少多少 |
| C. Model 接口 | 运行时发现遗漏，和忘写 `if` 本质一样；Raw SQL 绕过 |

---

## 10. 特殊场景处理（方案 D）

### 10.1 平台管理员查询特定租户

平台管理员通过 `X-Tenant-ID` header 指定要查看的租户时，auth middleware 已经将该 tenantID 写入 ctx。Plugin 自然生效：

```go
// auth_middleware.go 现有逻辑
if adminType == 1 && r.Header.Get("X-Tenant-ID") != "" {
    tenantID = parseHeader(r.Header.Get("X-Tenant-ID"))
    ctx = tenant.WithContext(ctx, shared.TenantID(tenantID))
}
// → Plugin 会自动注入 WHERE tenant_id = 指定租户
```

### 10.2 平台管理员查看全量数据

当 `tenantID == 0` 时，`tenant.FromContext(ctx)` 返回 0，Plugin 直接跳过。

### 10.3 跳过 Plugin（显式全租户查询）

```go
// 定义
func SkipTenantScope(ctx context.Context) context.Context {
    return context.WithValue(ctx, skipTenantScopeKey, true)
}

// 使用
db.WithContext(infra.SkipTenantScope(ctx)).Find(&allProducts)
```

### 10.4 JOIN 多表查询

AST 改写会遍历所有表引用，自动给每个 tenant-scoped 表加条件：

```go
// 自动改写为：
// SELECT * FROM orders o
// LEFT JOIN users u ON o.user_id = u.id
// WHERE o.tenant_id = 1 AND u.tenant_id = 1
db.WithContext(ctx).
    Table("orders").
    Joins("LEFT JOIN users ON orders.user_id = users.id").
    Find(&results)
```

### 10.5 Create 操作自动填充

```go
// Plugin 的 onCreate callback 自动设置 tenant_id
func (p *TenantScopePlugin) onCreate(db *gorm.DB) {
    tenantID, ok := tenant.FromContext(db.Statement.Context)
    if !ok || tenantID == 0 {
        return
    }
    tableName := db.Statement.Table
    if !IsTenantScopedTable(tableName) {
        return
    }
    db.Statement.SetColumn("tenant_id", tenantID.Int64())
}

// Repository 可以不再手动设置 tenant_id
func (r *productRepo) Create(ctx context.Context, db *gorm.DB, p *product.Product) error {
    model := toProductModel(p)
    // model.TenantID = ... // 不再需要，Plugin 自动设置
    return db.WithContext(ctx).Create(&model).Error
}
```

---

## 11. 安全保障

### 11.1 多层防线

Plugin 是**便利层**，不是唯一安全层：

| 层级 | 措施 |
|------|------|
| **应用层** | Tenant Scope Plugin 自动注入 |
| **注册表** | 集中管理，linter 检查遗漏 |
| **数据库层** | 关键表可启用 Row Level Security（MySQL 不原生支持，需应用层模拟） |
| **Review** | Code Review 时检查注册表是否完整 |

### 11.2 仍有的局限

| 场景 | 处理方式 |
|------|---------|
| 存储过程/触发器 | 超出 Plugin 能力范围，需 DB 层保障 |
| 多租户数据迁移脚本 | 使用 `SkipTenantScope(ctx)` 显式跳过 |
| 非 GORM 的数据库操作 | 不经过 Plugin，需手动处理 |

---

## 12. 迁移计划

### 阶段一：基础设施（1-2天）

- [ ] 建立 `tenant_scoped_tables` 注册表（`pkg/infra/tenant_tables.go`）
- [ ] 实现 `TenantScopePlugin`（GORM callback 注入）
- [ ] 实现 `SkipTenantScope`
- [ ] 实现 SQL 改写（`RewriteSQLWithTenantScope`）
- [ ] 注册 Plugin 到 `db.go`
- [ ] 编写单元测试

### 阶段二：验证（1天）

- [ ] 写集成测试：验证 ORM 查询自动注入
- [ ] 写集成测试：验证 Raw SQL 自动改写
- [ ] 写集成测试：验证平台管理员跳过
- [ ] 写集成测试：验证 `SkipTenantScope` 显式跳过
- [ ] 对比改造前后 SQL 输出，确认行为一致

### 阶段三：Repository 改造（3-5天，可逐步）

- [ ] 改造 `product_repository.go`（6 处）
- [ ] 改造 `shipment_repository.go`（7 处）
- [ ] 改造 `payment_repository.go`（4 处）
- [ ] 改造 `payment_transaction_repository.go`（4 处）
- [ ] 改造 `sku_repository.go`（4 处）
- [ ] 改造 `user_repository.go`（2 处）

每个 Repository 改造后：
1. 去掉方法中的 `tenantID` 参数
2. 去掉 `if tenantID != 0` 逻辑
3. 运行现有测试确认行为不变

### 阶段四：清理（1天）

- [ ] 更新 AGENTS.md 中的 Repository 规范
- [ ] 添加 linter 规则检查注册表完整性

---

## 13. 待讨论

1. **是否采用渐进式路径？** 先上方案 A（快速覆盖 ORM 查询），再升级到方案 D（加 AST 改写）
2. **是否同时统一 TenantID context？** 当前有 3 套 context key（`pkg/tenant`、`pkg/contextx`、`pkg/domain/shared`），建议统一为一套
3. **Shop 端是否也用？** Shop 端目前无条件过滤 tenant_id，可以也用 Plugin 统一
4. **注册表遗漏怎么防？** 建议写一个 test：扫描 SQL schema 中有 `tenant_id` 字段的表，对比注册表，差异则报警
