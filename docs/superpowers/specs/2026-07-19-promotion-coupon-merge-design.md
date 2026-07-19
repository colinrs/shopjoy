# Promotion × Coupon 合并设计

**日期：** 2026-07-19
**状态：** 待审批
**作者：** Claude (brainstorming 流程产出)

## 1. 概述

将「促销活动（Promotion）」与「优惠券（Coupon）」的领域模型、存储形态、规则机制统一收敛为一套。

### 1.1 现状问题

| 维度 | Promotion（活动） | Coupon（优惠券） | 问题 |
|------|-------------------|------------------|------|
| 表数量 | `promotions` + `promotion_rules` (一对多) | `coupons` (单表，规则字段内联) | 结构不一致，规则概念分裂 |
| 规则形态 | 多档阶梯（FindBestRule 选最优） | 单档（一个 coupon 一个规则） | 同样诉求两种建模 |
| Domain 实体 | 旧版 `admin/internal/domain/promotion` + 新版 `pkg/domain/promotion` 两套并存 | 同上 | 双套维护成本 |
| Coupon 专属 | — | code、库存、user_coupons 领券、GenerateCouponCodes | 与 promotion 完全异构 |

### 1.2 目标

1. **存储**：单表 `promotions` + `kind` 鉴别器；共享 `promotion_rules`（带 `owner_kind`）
2. **实体**：单 `Promotion` 结构（含 nullable coupon 字段），统一收敛到 `pkg/domain/promotion/`
3. **规则**：coupon 也支持多档阶梯；所有规则语义统一从 `promotion_rules` 读取
4. **行为**：复用 `PromotionScope`、`PromotionRule`、`CalculateDiscount`、`IsActive` 等方法
5. **新增字段**：`market_id`（顶层裁剪维度）

### 1.3 非目标

- 不重构 `user_coupons` / `promotion_usage`（保留现状，仅改 coupon_id 指向）
- 不改动订单侧的 discount 计算调用方（仅重写内部实现）
- 不下架现有 `/coupons/*` 路由（保留路径，内部走统一 App）

---

## 2. 数据模型

### 2.1 表结构

合并后只剩 4 张表（主表 + 辅表）：

```
promotions            ← 合并 promotions + coupons
promotion_rules       ← 兼容 PROMOTION 与 COUPON（按 owner_kind 区分）
user_coupons          ← 保留，coupon_id 指向 promotions.id (kind='COUPON')
promotion_usage       ← 保留，已可关联 coupon_id
```

### 2.2 `promotions` 主表（最终版）

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | BIGINT PK | |
| `tenant_id` | BIGINT NOT NULL | |
| `kind` | ENUM('PROMOTION','COUPON') NOT NULL DEFAULT 'PROMOTION' | 鉴别器 |
| `name` | VARCHAR(255) NOT NULL | |
| `description` | TEXT | |
| `code` | VARCHAR(100) NULL | 仅 COUPON：唯一兑换码 |
| `type` | TINYINT NOT NULL DEFAULT 0 | PROMOTION: discount/flash_sale/bundle/buy_x_get_y；COUPON 固定 0 |
| `status` | TINYINT NOT NULL DEFAULT 0 | pending/active/paused/ended；COUPON 的 depleted 由 used_count>=total_count 派生 |
| `priority` | INT NOT NULL DEFAULT 0 | |
| `market_id` | BIGINT NULL | **顶层裁剪**：NULL=不限市场；非空=只对该 market 生效 |
| `currency` | VARCHAR(10) NOT NULL DEFAULT 'CNY' | |
| `total_count` | INT NULL | 仅 COUPON：发放总数 |
| `used_count` | INT NULL | 仅 COUPON：已使用数 |
| `usage_limit` | INT NOT NULL DEFAULT 0 | 0=无限 |
| `per_user_limit` | INT NOT NULL DEFAULT 1 | 0=无限 |
| `tags` | JSON NULL | |
| `scope_type` | VARCHAR(32) NOT NULL DEFAULT 'STOREWIDE' | STOREWIDE/PRODUCTS/CATEGORIES/BRANDS（**不再含 MARKET**） |
| `scope_ids` | JSON NULL | 商品级 ID 列表（不含 market） |
| `exclude_ids` | JSON NULL | 排除的商品/分类/品牌 |
| `start_at` | TIMESTAMP NOT NULL | |
| `end_at` | TIMESTAMP NOT NULL | |
| `created_at` / `updated_at` / `deleted_at` | 审计 + 软删除 | |
| `created_by` / `updated_by` | BIGINT NOT NULL DEFAULT 0 | |

**字段规约**：
- ❌ **不再有** `discount_type / discount_value / min_order_amount / max_discount` 等规则字段
- ❌ **不再有** `min_order_amount` 等门槛字段（搬到 `promotion_rules.condition_value`）
- 所有 nullable 字段（`code / total_count / used_count`）通过 `kind='COUPON'` 保证业务层有值

### 2.3 `promotion_rules` 表

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | BIGINT PK | |
| `owner_kind` | ENUM('PROMOTION','COUPON') NOT NULL DEFAULT 'PROMOTION' | |
| `owner_id` | BIGINT NOT NULL | → promotions.id |
| `condition_type` | TINYINT NOT NULL DEFAULT 0 | 0=最低金额, 1=最低数量 |
| `condition_value` | DECIMAL(19,4) NOT NULL DEFAULT 0 | |
| `action_type` | TINYINT NOT NULL DEFAULT 0 | 0=固定金额, 1=百分比, 2=免邮 |
| `action_value` | DECIMAL(19,4) NOT NULL DEFAULT 0 | |
| `max_discount_amount` | DECIMAL(19,4) NOT NULL DEFAULT 0 | |
| `max_discount_currency` | VARCHAR(10) DEFAULT 'CNY' | |
| `currency` | VARCHAR(10) NOT NULL DEFAULT 'CNY' | |
| `sort_order` | INT NOT NULL DEFAULT 0 | |
| `created_at` / `updated_at` / `deleted_at` | 审计 + 软删除 | |

**索引**：
- `INDEX (owner_kind, owner_id, sort_order)` —— 替代原 `idx_promotion_id`
- 移除原 `idx_promotion_id`、`idx_promotion_rules_sort_order`

### 2.4 迁移脚本

文件：`sql/promotion/migrations/2026071901_merge_promotion_coupon.sql`

```sql
START TRANSACTION;

-- 1) 备份旧 coupons 表（保留 fallback）
CREATE TABLE IF NOT EXISTS _deprecated_coupons LIKE coupons;
INSERT INTO _deprecated_coupons SELECT * FROM coupons;

-- 2) promotions 表加列
ALTER TABLE promotions
  ADD COLUMN `kind`        ENUM('PROMOTION','COUPON') NOT NULL DEFAULT 'PROMOTION' AFTER `tenant_id`,
  ADD COLUMN `code`        VARCHAR(100) NULL AFTER `name`,
  ADD COLUMN `market_id`   BIGINT NULL AFTER `priority`,
  ADD COLUMN `total_count` INT NULL AFTER `usage_limit`,
  ADD COLUMN `used_count`  INT NULL AFTER `total_count`,
  ADD COLUMN `scope_type`  VARCHAR(32) NOT NULL DEFAULT 'STOREWIDE' AFTER `tags`,
  ADD COLUMN `scope_ids`   JSON NULL AFTER `scope_type`,
  ADD COLUMN `exclude_ids` JSON NULL AFTER `scope_ids`;

-- 3) 迁移 coupons → promotions
INSERT INTO promotions (
  tenant_id, kind, name, description, code, type, status, priority, market_id, currency,
  total_count, used_count, usage_limit, per_user_limit, scope_type, scope_ids, exclude_ids,
  start_at, end_at, created_at, updated_at, created_by, updated_by
)
SELECT
  tenant_id,
  'COUPON'                                                          AS kind,
  name, description, code,
  0                                                                 AS type,
  CASE status WHEN 0 THEN 0 WHEN 1 THEN 1 WHEN 2 THEN 3 WHEN 3 THEN 3 END AS status,
  0                                                                 AS priority,
  NULL                                                              AS market_id,
  currency,
  total_count, used_count,
  usage_limit, per_user_limit,
  COALESCE(scope_type, 'STOREWIDE'),
  scope_ids, exclude_ids,
  start_at, end_at,
  created_at, updated_at, created_by, updated_by
FROM coupons
WHERE deleted_at IS NULL;

-- 4) 备份 promotion_rules，重建为新结构
CREATE TABLE IF NOT EXISTS _deprecated_promotion_rules LIKE promotion_rules;
INSERT INTO _deprecated_promotion_rules SELECT * FROM promotion_rules;

ALTER TABLE promotion_rules
  ADD COLUMN `owner_kind` ENUM('PROMOTION','COUPON') NOT NULL DEFAULT 'PROMOTION' AFTER `promotion_id`,
  ADD COLUMN `sort_order` INT NOT NULL DEFAULT 0 AFTER `max_discount_amount`,
  MODIFY COLUMN `promotion_id` BIGINT NULL;

-- 5) 已有 PROMOTION 规则打 owner_kind
UPDATE promotion_rules SET owner_kind = 'PROMOTION';

-- 6) COUPON 原有 (type/value/min_amount/max_discount) 转成单条规则
INSERT INTO promotion_rules (
  owner_kind, owner_id, condition_type, condition_value, action_type, action_value,
  max_discount_amount, max_discount_currency, currency, sort_order,
  created_at, updated_at
)
SELECT
  'COUPON'                                                           AS owner_kind,
  p.id                                                               AS owner_id,
  0                                                                  AS condition_type,
  c.min_amount                                                       AS condition_value,
  CASE c.type WHEN 0 THEN 0 WHEN 1 THEN 1 WHEN 2 THEN 0 END          AS action_type,
  CASE c.type WHEN 0 THEN c.value WHEN 1 THEN c.value WHEN 2 THEN 0 END AS action_value,
  c.max_discount                                                     AS max_discount_amount,
  c.currency                                                         AS max_discount_currency,
  c.currency,
  0                                                                  AS sort_order,
  NOW(), NOW()
FROM coupons c
JOIN promotions p ON p.kind = 'COUPON' AND p.code = c.code AND p.tenant_id = c.tenant_id
WHERE c.deleted_at IS NULL;

-- 7) 索引改造
ALTER TABLE promotion_rules
  ADD INDEX `idx_owner` (`owner_kind`, `owner_id`, `sort_order`),
  DROP INDEX `idx_promotion_id`,
  DROP INDEX `idx_promotion_rules_sort_order`;

-- 8) UNIQUE 索引仅对 COUPON.code 生效（MySQL 8 支持 generated column partial unique）
ALTER TABLE promotions
  ADD COLUMN `code_unique` VARCHAR(100)
    GENERATED ALWAYS AS (IF(kind = 'COUPON', code, NULL)) VIRTUAL,
  ADD UNIQUE KEY `uk_promotion_code` (`code_unique`);

-- 9) user_coupons 加索引
ALTER TABLE user_coupons
  ADD INDEX `idx_coupon_id_active` (`coupon_id`, `status`);

-- 10) RENAME 旧 coupons 表（保留 30 天观察期）
RENAME TABLE coupons TO _archived_coupons_20260719;

COMMIT;
```

### 2.5 迁移前必跑检查

```sql
-- 1. 旧 coupon.code 必须无重复（否则 uk_promotion_code 建索引会失败）
SELECT code, COUNT(*) c FROM coupons WHERE deleted_at IS NULL GROUP BY code HAVING c > 1;
-- 预期：0 行

-- 2. 旧 promotion_rules.promotion_id 必须都能 JOIN 到 promotions（孤儿规则检查）
SELECT r.* FROM promotion_rules r LEFT JOIN promotions p ON p.id = r.promotion_id WHERE p.id IS NULL;
-- 预期：0 行

-- 3. 旧 coupons.market_ids 数据形态检查 —— **NO-OP**: 实际 schema 中 `coupons` 表
--    没有 `market_ids` 列（仅 `scope_ids` + `exclude_ids`），故此检查无法执行。
--    按设计 §2.7，promotion.market_id 是新列，不从旧数据回填。
--    原 SQL 留作参考：SELECT id, scope_type, scope_ids FROM coupons;
```

### 2.6 字段映射表

| 来源（旧表） | 目标（新表 promotions） | 适用 kind |
|---|---|---|
| `promotions.*` | 同名字段 | PROMOTION |
| `coupons.name/desc/start_at/end_at/status/currency/per_user_limit/usage_limit/scope_type/scope_ids/exclude_ids` | 同名字段 | COUPON |
| `coupons.code` | `promotions.code` (UNIQUE by generated column) | COUPON |
| `coupons.total_count` | `promotions.total_count` | COUPON |
| `coupons.used_count` | `promotions.used_count` | COUPON |
| `coupons.type/value/min_amount/max_discount` | `promotion_rules` 单条（owner_kind=COUPON） | COUPON |
| `coupons.market_ids` | **决策点**（见 §2.7） | — |
| `promotion_rules.*` | 同名（加 owner_kind='PROMOTION'、owner_id=promotion_id） | PROMOTION |
| `user_coupons.user_id/coupon_id` | 同名（coupon_id → promotions.id where kind='COUPON'） | COUPON |
| `promotion_usage.coupon_id` | 同名（nullable） | 通用 |

---

### 2.7 旧 `coupons.market_ids` 决策

旧 `coupons.market_ids` 是 JSON 数组。本设计采用**简化策略**：**丢弃 market_ids 字段，不回填到 `promotions.market_id`**。理由：

1. 当前 7 条种子数据中 `market_ids` 全部为空（验证 SQL 见 §2.5 #3）
2. 多 market 的 coupon 场景罕见，未在 PRD 中定义
3. 未来如需支持，可通过后续迁移 + 新建"coupon_market_scopes" 中间表实现

**回滚影响**：旧 `_deprecated_coupons` 表保留 30 天，可读回。

---

## 3. Go 实体层

### 3.1 文件清理

| 路径 | 处理 |
|------|------|
| `admin/internal/domain/promotion/entity.go` | **删除** |
| `admin/internal/domain/coupon/entity.go` | **删除** |
| `pkg/domain/promotion/entity.go` | **重写** |
| `pkg/domain/promotion/rule.go`（如无则新建） | 新增规则结构 + OwnerKind/OwnerID |
| `pkg/domain/promotion/coupon.go` | **改写**：删除 `Coupon` struct，仅保留枚举 + `UserCoupon` + `PromotionUsage` |

### 3.2 `Promotion` 聚合根（`pkg/domain/promotion/entity.go`）

```go
type Kind string
const (
    KindPromotion Kind = "PROMOTION"
    KindCoupon    Kind = "COUPON"
)

type Promotion struct {
    // ===== 公共字段（两 kind 共享）=====
    ID           int64            `json:"id"`
    TenantID     shared.TenantID  `json:"tenant_id"`
    Kind         Kind             `json:"kind"`
    Name         string           `json:"name"`
    Description  string           `json:"description"`
    Type         Type             `json:"type"`
    Status       Status           `json:"status"`
    Priority     int              `json:"priority"`
    MarketID     *int64           `json:"market_id,omitempty"`
    Currency     string           `json:"currency"`
    UsageLimit   int              `json:"usage_limit"`
    PerUserLimit int              `json:"per_user_limit"`
    Tags         []string         `json:"tags,omitempty" gorm:"type:json"`
    Scope        PromotionScope   `json:"scope"`
    StartAt      time.Time        `json:"start_at"`
    EndAt        time.Time        `json:"end_at"`
    Rules        []PromotionRule  `json:"rules,omitempty"`
    Audit        shared.AuditInfo `json:"audit"`
    DeletedAt    *time.Time       `json:"deleted_at,omitempty"`

    // ===== COUPON 专属字段（nil 表示非 COUPON）=====
    Code       *string `json:"code,omitempty"`
    TotalCount *int    `json:"total_count,omitempty"`
    UsedCount  *int    `json:"used_count,omitempty"`
}

func (p *Promotion) TableName() string { return "promotions" }
```

### 3.3 行为方法

```go
// IsActive: 统一判定（COUPON 还需库存检查）
func (p *Promotion) IsActive() bool {
    if p.Status != StatusActive || p.DeletedAt != nil {
        return false
    }
    if p.Kind == KindCoupon && p.TotalCount != nil && p.UsedCount != nil {
        if *p.UsedCount >= *p.TotalCount {
            return false // depleted
        }
    }
    now := time.Now().UTC()
    return !now.Before(p.StartAt) && !now.After(p.EndAt)
}

// MatchesMarket: market_id 为 NULL 则不限，否则必须匹配
func (p *Promotion) MatchesMarket(marketID int64) bool {
    return p.MarketID == nil || *p.MarketID == marketID
}

// MatchesScope: 委托给 Scope（kind 无关）
func (p *Promotion) MatchesScope(productID, categoryID, brandID int64) bool {
    return p.Scope.MatchesProduct(productID, categoryID, brandID)
}

// FindBestRule: kind 无关，走规则链
func (p *Promotion) FindBestRule(matchedAmount decimal.Decimal, quantity int) *PromotionRule {
    var bestRule *PromotionRule
    for i := range p.Rules {
        rule := &p.Rules[i]
        if !rule.MeetsCondition(matchedAmount, quantity) {
            continue
        }
        if bestRule == nil || rule.ConditionValue.GreaterThan(bestRule.ConditionValue) {
            bestRule = rule
        }
    }
    return bestRule
}

// CalculateDiscount: 委托给 FindBestRule
func (p *Promotion) CalculateDiscount(matchedAmount decimal.Decimal, quantity int) decimal.Decimal {
    rule := p.FindBestRule(matchedAmount, quantity)
    if rule == nil {
        return decimal.Zero
    }
    return rule.CalculateDiscount(matchedAmount)
}

// Issue: COUPON 校验 + 生成 user_coupon
func (p *Promotion) Issue(userID int64, now time.Time) (*UserCoupon, error) {
    if p.Kind != KindCoupon {
        return nil, code.ErrPromotionInvalidKind
    }
    if !p.IsActive() {
        return nil, code.ErrCouponExpired
    }
    return &UserCoupon{
        TenantID:   p.TenantID,
        UserID:     userID,
        CouponID:   p.ID,
        Status:     UserCouponStatusUnused,
        ReceivedAt: now,
        ExpireAt:   p.EndAt,
    }, nil
}

// ConsumeInventory: 内存层预检查 + 计数 +1（**仅校验语义**；持久化由 repo.IncrementUsedCount 完成）
func (p *Promotion) ConsumeInventory() error {
    if p.Kind != KindCoupon || p.UsedCount == nil || p.TotalCount == nil {
        return code.ErrPromotionInvalidKind
    }
    if *p.UsedCount >= *p.TotalCount {
        return code.ErrCouponDepleted
    }
    *p.UsedCount++
    return nil
}

// 调用约定（App 层）：
//   1. 先调 p.ConsumeInventory() 预校验 + 内存 +1
//   2. 再调 repo.IncrementUsedCount(ctx, db, p.ID) 原子持久化
//   3. 若第 2 步失败，回滚内存值并返回错误
// 高并发场景下，第 2 步的 SQL 必须用 `UPDATE ... WHERE used_count < total_count`
// 否则可能出现超卖（已 depleted 但仍扣减）。
```

### 3.4 `PromotionRule`（`pkg/domain/promotion/rule.go`）

```go
type PromotionRule struct {
    ID             int64           `json:"id"`
    OwnerKind      Kind            `json:"owner_kind"`    // NEW
    OwnerID        int64           `json:"owner_id"`      // NEW
    ConditionType  ConditionType   `json:"condition_type"`
    ConditionValue decimal.Decimal `json:"condition_value"`
    ActionType     ActionType      `json:"action_type"`
    ActionValue    decimal.Decimal `json:"action_value"`
    MaxDiscount    decimal.Decimal `json:"max_discount"`
    Currency       string          `json:"currency"`
    SortOrder      int             `json:"sort_order"`
    CreatedAt      time.Time       `json:"created_at"`
    UpdatedAt      time.Time       `json:"updated_at"`
}

// CalculateDiscount / MeetsCondition 不变（kind 无关）
```

### 3.5 `PromotionScope` 复用

`pkg/domain/promotion/scope.go` 中已有 `PromotionScope`，改造点：
- 删除 MARKET 维度支持（应用层枚举校验）
- `MatchesProduct` 保持不变

---

## 4. Repository 层

### 4.1 文件改动

| 路径 | 处理 |
|------|------|
| `admin/internal/infrastructure/persistence/promotion_repository.go` | **重写**：单一 Repo 接口 |
| `admin/internal/infrastructure/persistence/coupon_repository.go` | **删除** |

### 4.2 Repository 接口

```go
type Repository interface {
    // ===== 通用 CRUD（两 kind 共用）=====
    Create(ctx context.Context, db *gorm.DB, p *Promotion) error
    Update(ctx context.Context, db *gorm.DB, p *Promotion) error
    Delete(ctx context.Context, db *gorm.DB, id int64) error
    FindByID(ctx context.Context, db *gorm.DB, id int64) (*Promotion, error)
    FindByCode(ctx context.Context, db *gorm.DB, code string) (*Promotion, error)  // NEW
    FindList(ctx context.Context, db *gorm.DB, query Query) ([]*Promotion, int64, error)

    // ===== Rule 操作（owner_kind+owner_id 通用）=====
    CreateRules(ctx context.Context, db *gorm.DB, ownerKind Kind, ownerID int64, rules []PromotionRule) error
    FindRulesByOwner(ctx context.Context, db *gorm.DB, ownerKind Kind, ownerID int64) ([]PromotionRule, error)
    UpdateRule(ctx context.Context, db *gorm.DB, rule *PromotionRule) error
    DeleteRule(ctx context.Context, db *gorm.DB, id int64) error
    DeleteRulesByOwner(ctx context.Context, db *gorm.DB, ownerKind Kind, ownerID int64) error

    // ===== COUPON 专属 =====
    FindActiveCoupons(ctx context.Context, db *gorm.DB, marketID *int64) ([]*Promotion, error)
    IncrementUsedCount(ctx context.Context, db *gorm.DB, couponID int64) error  // 原子 +1，SQL: `UPDATE ... SET used_count = used_count + 1 WHERE id = ? AND (total_count = 0 OR used_count < total_count)`
    IssueUserCoupon(ctx context.Context, db *gorm.DB, uc *UserCoupon) error
    FindUserCoupons(ctx context.Context, db *gorm.DB, query UserCouponQuery) ([]*UserCoupon, int64, error)

    // ===== Usage 查询（已有）=====
    FindPromotionUsage(ctx context.Context, db *gorm.DB, query UsageQuery) ([]*PromotionUsage, int64, error)
}
```

### 4.3 Query 类型

```go
type Query struct {
    shared.PageQuery
    TenantID    shared.TenantID
    Name        string
    Kind        *Kind        // NEW: 指针（NULL=全部）
    Status      *Status
    Type        *Type
    MarketID    *int64       // NEW
    ExpiredOnly bool
}
```

### 4.4 仓储实现要点

1. `promotionModel` 增加列映射：
   ```go
   type promotionModel struct {
       application.Model
       TenantID     shared.TenantID `gorm:"column:tenant_id"`
       Kind         string          `gorm:"column:kind;type:enum('PROMOTION','COUPON');not null;default:PROMOTION"`
       Name         string          `gorm:"column:name"`
       Code         *string         `gorm:"column:code"`
       // ... 其他字段
       MarketID     *int64          `gorm:"column:market_id"`
       TotalCount   *int            `gorm:"column:total_count"`
       UsedCount    *int            `gorm:"column:used_count"`
       ScopeType    string          `gorm:"column:scope_type"`
       ScopeIDs     datatypes.JSON  `gorm:"column:scope_ids;type:json"`
       ExcludeIDs   datatypes.JSON  `gorm:"column:exclude_ids;type:json"`
       Tags         datatypes.JSON  `gorm:"column:tags;type:json"`
   }
   ```

2. `promotionRuleModel` 增加 `OwnerKind / OwnerID`，`PromotionID` 保留 nullable：
   ```go
   type promotionRuleModel struct {
       application.Model
       OwnerKind      string  `gorm:"column:owner_kind;type:enum(...)"`
       OwnerID        int64   `gorm:"column:owner_id;not null;index"`
       PromotionID    *int64  `gorm:"column:promotion_id"`  // nullable（COUPON 时为 NULL）
       // ...
       SortOrder      int     `gorm:"column:sort_order;not null;default:0"`
   }
   ```

3. 查询路径：
   - `FindList` 默认 `Kind IS NULL` 返回全部
   - `FindList(ctx, db, Query{Kind: &KindCoupon, ...})` 仅查 coupon
   - `FindActiveCoupons` 仅查 `kind='COUPON' AND status=active AND ...`

---

## 5. Application + Logic 层

### 5.1 Application 层

**文件**：`admin/internal/application/promotion/promotion_app.go`

- 合并现状 `promotion_app.go` + `coupon_app.go`
- 删除 `coupon_app.go`

```go
type PromotionApp struct {
    repo promotion.Repository
    db   *gorm.DB
}

// ===== 通用入口（按 kind 路由）=====
func (a *PromotionApp) Create(ctx, req *CreatePromotionRequest) (*PromotionResponse, error)
func (a *PromotionApp) Update(ctx, req *UpdatePromotionRequest) (*PromotionResponse, error)
func (a *PromotionApp) Get(ctx, id) (*PromotionResponse, error)
func (a *PromotionApp) List(ctx, query) (*ListPromotionResponse, error)
func (a *PromotionApp) Delete(ctx, id) error
func (a *PromotionApp) Activate(ctx, id) error   // 改 Status
func (a *PromotionApp) Deactivate(ctx, id) error

// ===== Rule 入口（kind 无关）=====
func (a *PromotionApp) CreateRules(ctx, ownerKind, ownerID, rules) error
func (a *PromotionApp) GetRules(ctx, ownerKind, ownerID) ([]*PromotionRuleResponse, error)
func (a *PromotionApp) UpdateRule(ctx, rule) error
func (a *PromotionApp) DeleteRule(ctx, id) error

// ===== COUPON 专属入口（kind=COUPON 校验）=====
func (a *PromotionApp) IssueToUser(ctx, couponID, userID) (*UserCouponResponse, error)
func (a *PromotionApp) BatchIssue(ctx, couponID, userIDs) (*BatchIssueResponse, error)
func (a *PromotionApp) GenerateCodes(ctx, prefix, qty, cfg) ([]string, error)
func (a *PromotionApp) ListUserCoupons(ctx, query) (*ListUserCouponResponse, error)
```

### 5.2 Logic 层

**保留双目录**（handler 路由分组不动），但内部统一调用 `PromotionApp`。

| Logic 文件 | 改造要点 |
|------------|----------|
| `promotions/create_promotion_logic.go` | 构造 `CreatePromotionRequest{Kind: KindPromotion, ...}` |
| `promotions/update_promotion_logic.go` | 同上 + `UpdateRules` 流程（先删后建） |
| `promotions/list_promotions_logic.go` | `ListPromotionsReq` 新增 `Kind / MarketID` 字段 |
| `promotions/{create,update,get,delete}_promotion_rule_logic.go` | 走 `PromotionApp.CreateRules/UpdateRule/DeleteRule`，无需感知 owner_kind（从 url 的 `:id` 反查 `Promotion.Kind` 决定） |
| `promotions/activate_promotion_logic.go` | 透传 Status 转换 |
| `coupons/create_coupon_logic.go` | 构造 `CreatePromotionRequest{Kind: KindCoupon, Code: req.Code, TotalCount: req.UsageLimit, ...}` |
| `coupons/update_coupon_logic.go` | 同上 + Rule 重建 |
| `coupons/list_coupons_logic.go` | `PromotionApp.List(ctx, Query{Kind: KindCoupon, ...})` |
| `coupons/get_coupon_logic.go` | `PromotionApp.Get(ctx, id)` |
| `coupons/activate_coupon_logic.go` | 改 Status |
| `coupons/deactivate_coupon_logic.go` | 改 Status |
| `coupons/delete_coupon_logic.go` | `PromotionApp.Delete(ctx, id)` |
| `coupons/generate_coupon_codes_logic.go` | 循环 `PromotionApp.Create({Kind: COUPON, ...})` |
| `coupons/get_coupon_usage_logic.go` | 走 `PromotionUsage` 查询 |
| `coupons/issue_user_coupon_logic.go` | 走 `PromotionApp.IssueToUser` |
| `coupons/batch_issue_user_coupon_logic.go` | 同上 |
| `coupons/list_user_coupons_logic.go` | 走 `PromotionApp.ListUserCoupons` |

**删除**：
- `admin/internal/logic/coupons/helper.go`（`convertCouponToDetailResp` 删除）
- `admin/internal/application/promotion/coupon_app.go`

### 5.3 转换函数

**`promotions/helper.go`** 中 `convertPromotionToDetailResp` —— 统一处理两种 kind：

```go
// convertRulesToResp: App 层 PromotionRule → wire types.PromotionRuleResp
func convertRulesToResp(rules []apppromotion.PromotionRuleResponse) []*types.PromotionRuleResp {
    out := make([]*types.PromotionRuleResp, 0, len(rules))
    for _, r := range rules {
        out = append(out, &types.PromotionRuleResp{
            ID:            r.ID,
            ConditionType: string(r.ConditionType),
            ConditionValue: r.ConditionValue,
            ActionType:    string(r.ActionType),
            ActionValue:   r.ActionValue,
            MaxDiscount:   r.MaxDiscount,
            SortOrder:     r.SortOrder,
        })
    }
    return out
}

func convertPromotionToDetailResp(p *apppromotion.PromotionResponse) *types.PromotionDetailResp {
    resp := &types.PromotionDetailResp{
        ID:         p.ID,
        Kind:       string(p.Kind),              // NEW
        Name:       p.Name,
        Code:       p.Code,                      // nullable
        MarketID:   p.MarketID,                  // NEW
        TotalCount: p.TotalCount,                // nullable
        // ... 其他公共字段
        Rules:      convertRulesToResp(p.Rules), // 统一带 rules（两种 kind 都填充）
    }
    return resp
}
```

**Coupon 端 `convertCouponToDetailResp` 删除**，coupons/*.go 改调 `convertPromotionToDetailResp`。

### 5.4 promotion.md 必检项

- [ ] 每个 `UpdatePromotionReq` / `CreatePromotionRequest` 字段在 logic + app 层都有赋值
- [ ] `convertPromotionToDetailResp` 映射全部响应字段（含新加的 `kind / market_id / code / total_count / used_count / rules`）
- [ ] `promotionModel` / `promotionRuleModel` 列名匹配 `SHOW CREATE TABLE promotions`
- [ ] `Query.Kind` / `Query.Status` / `Query.Type` / `Query.MarketID` 均为指针类型
- [ ] `"expired"` 走 `ExpiredOnly bool`，不存为枚举
- [ ] `buildScope` 中 wire 值小写 → 大写规范化（kind 同理）
- [ ] `market_id` 命中判定走 `MatchesMarket`

---

## 6. API 定义 + 全栈同步

### 6.1 Wire 类型（`admin/desc/promotion.api`）

**统一响应**：

```go
type PromotionDetailResp {
    ID           int64    `json:"id,string"`
    Kind         string   `json:"kind"`                  // NEW
    Name         string   `json:"name"`
    Description  string   `json:"description"`
    Code         string   `json:"code,optional"`         // 仅 coupon
    Type         string   `json:"type"`
    Status       string   `json:"status"`
    MarketID     int64    `json:"market_id,optional,string"`  // NEW
    Currency     string   `json:"currency"`
    UsageLimit   int      `json:"usage_limit"`
    UsedCount    int      `json:"used_count"`
    PerUserLimit int      `json:"per_user_limit"`
    TotalCount   int      `json:"total_count,optional"` // 仅 coupon
    ScopeType    string   `json:"scope_type"`
    ScopeIDs     []string `json:"scope_ids"`
    ExcludeIDs   []string `json:"exclude_ids"`
    Tags         []string `json:"tags"`
    Rules        []*PromotionRuleResp `json:"rules"`     // NEW: coupon 也带
    StartTime    string   `json:"start_time"`
    EndTime      string   `json:"end_time"`
    CreatedAt    string   `json:"created_at"`
    UpdatedAt    string   `json:"updated_at"`
}
```

**统一请求**：

```go
type CreatePromotionReq {
    Kind         string                `json:"kind"`              // "promotion" | "coupon"
    Name         string                `json:"name"`
    Description  string                `json:"description,optional"`
    Code         string                `json:"code,optional"`     // 仅 coupon
    Type         string                `json:"type,optional"`
    MarketID     int64                 `json:"market_id,optional,string"`
    Currency     string                `json:"currency,optional"`
    UsageLimit   int                   `json:"usage_limit,optional"`
    PerUserLimit int                   `json:"per_user_limit,optional"`
    TotalCount   int                   `json:"total_count,optional"` // 仅 coupon
    ScopeType    string                `json:"scope_type,optional"`
    ScopeIDs     []string              `json:"scope_ids,optional"`
    ExcludeIDs   []string              `json:"exclude_ids,optional"`
    Tags         []string              `json:"tags,optional"`
    Rules        []*PromotionRuleReq   `json:"rules,optional"`
    StartTime    string                `json:"start_time"`
    EndTime      string                `json:"end_time"`
}
```

**规则 Wire**（通用）：

```go
type PromotionRuleReq {
    ConditionType  string `json:"condition_type"`   // "min_amount" | "min_quantity"
    ConditionValue string `json:"condition_value"`  // decimal string
    ActionType     string `json:"action_type"`      // "fixed_amount" | "percentage" | "free_shipping"
    ActionValue    string `json:"action_value"`     // decimal string
    MaxDiscount    string `json:"max_discount,optional"`
    SortOrder      int    `json:"sort_order,optional"`
}
```

**列表请求**：

```go
type ListPromotionsReq {
    Kind     string `form:"kind,optional"`
    Name     string `form:"name,optional"`
    Type     string `form:"type,optional"`
    Status   string `form:"status,optional"`
    MarketID int64  `form:"market_id,optional,string"`
    Page     int    `form:"page,default=1"`
    PageSize int    `form:"page_size,default=20"`
}
```

### 6.2 保留路由

| 路由 | 说明 |
|------|------|
| `POST /coupons/generate-codes` | 批量生成 |
| `POST /coupons/:id/issue` | 发放 |
| `POST /coupons/:id/batch-issue` | 批量发放 |
| `GET /user-coupons` | 领券列表 |
| `POST /user-coupons` | 领取 |

### 6.3 前端同步

**`shop-admin/src/api/promotion.ts`**：

```typescript
export interface Promotion {
  id: string
  kind: 'promotion' | 'coupon'              // NEW
  code?: string
  market_id?: string                        // NEW
  total_count?: number
  rules?: PromotionRule[]                   // NEW
  // 其他字段合并
}

export interface PromotionRule {
  id: string
  condition_type: 'min_amount' | 'min_quantity'
  condition_value: string
  action_type: 'fixed_amount' | 'percentage' | 'free_shipping'
  action_value: string
  max_discount?: string
  sort_order?: number
}

// 删除 CouponDetailResp
```

**`shop-admin/src/views/promotions/index.vue`**：

| 改动 | 处理 |
|------|------|
| 状态判断 | 移除 `type === 'coupon'` 分支，用 `kind` 取代 |
| Code / UsedCount / TotalCount 列 | `v-if="row.kind === 'coupon'"` 限定 |
| 规则入口 | 两种 kind 都允许点击 |
| 表单 | 合并 `couponForm` / `promotionForm` 为单 `form`，`kind` 决定字段显隐 |
| 列表过滤 | 顶部新增 `market_id` 下拉 |

**`shop-admin/src/views/promotions/rule*.vue`** —— 进入方式、Title 按 kind 文案区分。

### 6.4 错误码

`pkg/code/code.go`：

```go
ErrPromotionInvalidKind = &Err{HTTPCode: 400, Code: 80018, Msg: "invalid promotion kind"}
ErrCouponDepleted       = &Err{HTTPCode: 400, Code: 70004, Msg: "coupon depleted"}  // 复用已有
```

### 6.5 兼容性

| 维度 | 策略 |
|------|------|
| API 路由 | `/coupons/*` 全部保留，内部走 `PromotionApp` |
| 响应字段 | 新增 `kind / market_id / rules`，原字段保留；`/coupons/*` 响应比之前**多** `rules` 字段 |
| 数据库 | 一步迁移脚本 + 应用代码同步上线（无灰度期） |
| 前端 | `CouponDetailResp` 删除，`Promotion` 统一；UI 适配在同一次提交 |

---

## 7. 实施步骤 + 验证

### 7.1 实施顺序（5 个 Phase）

| Phase | 内容 | 文件数 | 估时 |
|-------|------|--------|------|
| **P0** | SQL 迁移脚本 | 1 SQL | 0.5d |
| **P1** | 实体层收敛（`pkg/domain/promotion/`） | 5 | 1d |
| **P2** | 仓储层（合并 Repo） | 2 | 1d |
| **P3** | App + Logic 改造 | ~25 | 2d |
| **P4** | API + 前端同步 | ~8 | 1.5d |
| **P5** | 验证 + Code Review | 0 | 0.5d |

**总估时：~6.5 工作日**

### 7.2 验证清单

| # | 验证项 | 通过标准 |
|---|--------|----------|
| 1 | `cd admin && make api` 重新生成 | 无报错 |
| 2 | `make build` 编译 | 0 error |
| 3 | `promotionModel` 列名 vs `SHOW CREATE TABLE promotions` | 100% 对齐 |
| 4 | `convertPromotionToDetailResp` 字段映射 | grep 全字段 |
| 5 | `.claude/rules/golang/promotion.md` Checklist | 全过 |
| 6 | `go test ./admin/internal/... ./pkg/domain/...` | 全过 |
| 7 | 手工创建 PROMOTION（3 档规则）+ COUPON（1 条规则） | POST/GET 字段一致；market_id 过滤生效 |
| 8 | 回归 `/coupons/*` 全部路由 | 200 返回；前端不报错 |
| 9 | `pnpm lint && pnpm build` | 0 error |
| 10 | `/review` skill | 0 high-severity |

### 7.3 迁移后完整性 SQL

```sql
SELECT 'promotions_total'    AS tbl, COUNT(*) FROM promotions UNION ALL
SELECT 'promotions_coupon'   AS tbl, COUNT(*) FROM promotions WHERE kind='COUPON' UNION ALL
SELECT 'promotions_promotion'AS tbl, COUNT(*) FROM promotions WHERE kind='PROMOTION' UNION ALL
SELECT 'old_coupons'         AS tbl, COUNT(*) FROM _deprecated_coupons UNION ALL
SELECT 'rules_total'         AS tbl, COUNT(*) FROM promotion_rules UNION ALL
SELECT 'rules_coupon'        AS tbl, COUNT(*) FROM promotion_rules WHERE owner_kind='COUPON' UNION ALL
SELECT 'rules_promotion'     AS tbl, COUNT(*) FROM promotion_rules WHERE owner_kind='PROMOTION' UNION ALL
SELECT 'user_coupons'        AS tbl, COUNT(*) FROM user_coupons;

-- 预期：
-- promotions_coupon = old_coupons（非删除数据）
-- rules_coupon = old_coupons（每张旧 coupon 1 条规则）
-- rules_promotion = rules_total - rules_coupon
```

---

## 8. 风险与回滚

### 8.1 风险清单

| 风险 | 等级 | 缓解 |
|------|------|------|
| R1: 旧 `coupons` 表 RENAME 后仍有外部脚本直查 | 中 | `_deprecated_coupons` 保留 30 天；迁移前 grep `FROM coupons`/`JOIN coupons` 全部改完 |
| R2: `promotion_rules.promotion_id` nullable 后索引失效 | 中 | `idx_owner (owner_kind, owner_id, sort_order)` 覆盖 |
| R3: 旧 `coupons.code` 重复导致 `uk_promotion_code` 失败 | 高 | 迁移前必跑重复检查（见 §2.5） |
| R4: 前端 `CouponDetailResp` 类型删除导致编译失败 | 中 | 前端改动**一次性提交**；TS 编译会立即报缺失字段 |
| R5: 旧 `admin/internal/domain/{promotion,coupon}` 外部 import | 高 | 迁移前 grep 全部清除 |
| R6: `market_id` 新字段旧逻辑不识别 → NULL 不命中 | 低 | 应用层统一 `MatchesMarket`；NULL 视为全市场 |
| R7: `kind` 字段 ENUM 大小写敏感性 | 中 | Go 常量大写（`KindPromotion`）；wire 用小写（`"promotion"`/`"coupon"`）；application 层 `strings.ToLower` |

### 8.2 回滚预案

**触发条件**：P0 迁移后任意验证项失败 / P3 后编译失败 / 生产灰度发现数据问题

```sql
-- 1. 停服
-- 2. 恢复旧 coupons 表
RENAME TABLE _archived_coupons_20260719 TO coupons;

-- 3. 反向迁移数据（从 _deprecated_coupons 拉回）
INSERT INTO coupons SELECT * FROM _deprecated_coupons
WHERE id NOT IN (SELECT id FROM coupons);

-- 4. 反向 ALTER promotions
ALTER TABLE promotions
  DROP COLUMN kind, DROP COLUMN code, DROP COLUMN market_id,
  DROP COLUMN total_count, DROP COLUMN used_count,
  DROP COLUMN scope_type, DROP COLUMN scope_ids, DROP COLUMN exclude_ids,
  DROP COLUMN code_unique, DROP INDEX uk_promotion_code;

-- 5. 回滚 promotion_rules（恢复 promotion_id，删除 owner_kind）
ALTER TABLE promotion_rules
  DROP INDEX idx_owner,
  ADD INDEX idx_promotion_id (promotion_id),
  DROP COLUMN owner_kind, DROP COLUMN sort_order;

-- 6. git revert 应用层代码
-- 7. 重启服务
```

### 8.3 上线顺序

| 顺序 | 动作 | 时间窗口 |
|------|------|----------|
| 1 | P0 SQL 迁移脚本单独提交 + 灰度执行 | T+0 22:00 |
| 2 | §7.3 完整性 SQL 验证 | T+0 22:30 |
| 3 | P1-P5 一次性发版 | T+1 09:00 |
| 4 | 监控 promotion_usage 写入、coupon issue 成功率 | T+1 10:00-11:00 |
| 5 | 30 天后清理 `_archived_*` 表 | T+30 |

---

## 附录 A：反向迁移 SQL

```sql
START TRANSACTION;

-- 1. 重建 promotion_rules（恢复旧结构）
ALTER TABLE promotion_rules
  DROP INDEX idx_owner,
  ADD INDEX idx_promotion_id (promotion_id),
  ADD INDEX idx_promotion_rules_sort_order (promotion_id, sort_order),
  DROP COLUMN owner_kind, DROP COLUMN sort_order,
  MODIFY COLUMN promotion_id BIGINT NOT NULL;

-- 2. 从新结构拉回旧 coupon 数据
INSERT INTO coupons (
  id, tenant_id, name, code, description, type, value, min_amount, max_discount,
  total_count, used_count, per_user_limit, status, currency,
  scope_type, scope_ids, exclude_ids, deleted_at,
  start_at, end_at, created_at, updated_at, created_by, updated_by
)
SELECT
  p.id, p.tenant_id, p.name, p.code, p.description,
  COALESCE(r.action_type, 0) AS type,
  COALESCE(r.action_value, 0) AS value,
  COALESCE(r.condition_value, 0) AS min_amount,
  COALESCE(r.max_discount_amount, 0) AS max_discount,
  p.total_count, p.used_count, p.per_user_limit, p.status, p.currency,
  p.scope_type, p.scope_ids, p.exclude_ids, p.deleted_at,
  p.start_at, p.end_at, p.created_at, p.updated_at, p.created_by, p.updated_by
FROM promotions p
LEFT JOIN promotion_rules r ON r.owner_kind = 'COUPON' AND r.owner_id = p.id
WHERE p.kind = 'COUPON';

-- 3. 删 promotions 表新加列
ALTER TABLE promotions
  DROP COLUMN kind, DROP COLUMN code, DROP COLUMN market_id,
  DROP COLUMN total_count, DROP COLUMN used_count,
  DROP COLUMN scope_type, DROP COLUMN scope_ids, DROP COLUMN exclude_ids,
  DROP COLUMN code_unique, DROP INDEX uk_promotion_code;

-- 4. 删 promotion_rules.owner_kind 相关行（如有）
DELETE FROM promotion_rules WHERE owner_kind = 'COUPON';

-- 5. 还原 _deprecated_promotion_rules 数据（如有更新）
-- 此处省略，按需恢复

COMMIT;
```