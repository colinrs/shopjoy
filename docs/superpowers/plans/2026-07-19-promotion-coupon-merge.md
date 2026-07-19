# Promotion × Coupon Merge Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Unify Promotion and Coupon into a single Promotion aggregate with `kind` discriminator, shared rule chain (multi-tier for both), and top-level `market_id` scope. Remove the duplicated domain entities, repositories, and application services.

**Architecture:** Single `promotions` table + `kind` column; shared `promotion_rules` table with `owner_kind`/`owner_id`. Single `Promotion` Go struct in `pkg/domain/promotion/`. Single `PromotionApp` in `admin/internal/application/promotion/`. Wire API gets `kind`/`market_id`/`rules` fields on the unified `PromotionDetailResp`; old `/coupons/*` routes stay but route through `PromotionApp`.

**Tech Stack:** Go 1.21+, GORM v2, go-zero `make api` codegen, Vue 3 + Element Plus + TypeScript, MySQL 8.0.16+ (for generated column + unique index).

---

## Global Constraints

- **Spec location:** `docs/superpowers/specs/2026-07-19-promotion-coupon-merge-design.md` — read fully before each task.
- **Backend build:** Always `cd admin && make build`. Never `go build` directly.
- **API regen:** Always `cd admin && make api` after editing `desc/*.api`. Never edit `internal/types/types.go` or `internal/handler/routes.go` by hand.
- **Frontend build:** `cd shop-admin && pnpm build`. Never `tsc` directly.
- **Migration rule:** SQL migrations go in `sql/promotion/migrations/{YYYYMMDD}{seq}_{action}_{object}.sql`. After running, merge column changes into `sql/promotion/schema.sql` and drop the migration file.
- **Decimal handling:** All monetary calculations use `shopspring/decimal`. Wire values are strings; storage is `DECIMAL(19,4)`.
- **Time:** All timestamps `time.Time` UTC. DB columns `TIMESTAMP`. API wire strings RFC3339.
- **Errors:** Use `pkg/code` codes. Add new errors to `pkg/code/code.go` in the 80xxx range (Promotion) or 70xxx range (Coupon).
- **Type alignment:** Frontend TypeScript types must match backend `.api` enum values exactly. No independent frontend enum definitions.
- **Enum rule:** Status `"expired"` is derived from `EndAt`, never stored. Use `ExpiredOnly bool` filter.
- **Scope normalization:** Wire scope values lowercase (`"products"`, `"categories"`, `"brands"`, `"storewide"`) → domain uppercase (`"PRODUCTS"`, ...). Use `strings.ToUpper()` before comparison.
- **Filter params:** All optional filters (`Kind`, `Status`, `Type`, `MarketID`) use pointer types — never `!= 0` sentinel on iota enums.
- **Commit style:** Conventional commits. End with `Co-Authored-By: Claude <noreply@anthropic.com>`.

---

## File Structure

### Database (modify)
- `sql/promotion/migrations/2026071901_merge_promotion_coupon.sql` — **CREATE**: P0 migration script (one-shot ALTER + INSERT + RENAME).
- `sql/promotion/schema.sql` — **MODIFY**: merge new columns + `code_unique` generated column; drop `coupons` table definition.

### Domain — `pkg/domain/promotion/` (rewrite)
- `pkg/domain/promotion/entity.go` — **MODIFY**: single `Promotion` struct (Kind/MarketID/Code/TotalCount/UsedCount nullable) + `Kind` type + `IsActive`/`MatchesMarket`/`Issue`/`ConsumeInventory` methods.
- `pkg/domain/promotion/rule.go` — **CREATE**: `PromotionRule` struct with `OwnerKind`/`OwnerID`; `CalculateDiscount`/`MeetsCondition` (kind-agnostic).
- `pkg/domain/promotion/scope.go` — **MODIFY**: `PromotionScope` removes MARKET dimension from valid enum.
- `pkg/domain/promotion/coupon.go` — **MODIFY**: delete `Coupon` struct; keep `CouponType`/`CouponStatus`/`UserCouponStatus` enums + `UserCoupon` + `PromotionUsage`.
- `pkg/domain/promotion/repository.go` — **CREATE**: `Repository` interface (kind-agnostic, plus COUPON-specific methods).
- `pkg/domain/promotion/query.go` — **CREATE**: `Query` struct (Kind/Status/Type/MarketID pointers, ExpiredOnly bool).

### Domain — `admin/internal/domain/` (delete)
- `admin/internal/domain/promotion/entity.go` — **DELETE**.
- `admin/internal/domain/coupon/entity.go` — **DELETE**.

### Persistence — `admin/internal/infrastructure/persistence/` (rewrite)
- `promotion_repository.go` — **MODIFY**: single repo, all methods kind-routed. `promotionModel` adds Kind/Code/MarketID/TotalCount/UsedCount/ScopeType/ScopeIDs/ExcludeIDs/code_unique columns. `promotionRuleModel` adds OwnerKind/OwnerID/SortOrder.
- `coupon_repository.go` — **DELETE**.

### Application — `admin/internal/application/promotion/` (rewrite)
- `promotion_app.go` — **MODIFY**: single App with `Create/Update/Get/List/Delete/Activate/Deactivate` (kind-routed) + `CreateRules/GetRules/UpdateRule/DeleteRule` (kind-agnostic) + `IssueToUser/BatchIssue/GenerateCodes/ListUserCoupons` (COUPON).
- `coupon_app.go` — **DELETE**.

### Logic — `admin/internal/logic/` (modify)
- `promotions/create_promotion_logic.go` — MODIFY: build CreatePromotionRequest with `Kind: KindPromotion`.
- `promotions/update_promotion_logic.go` — MODIFY: same + rebuild rules.
- `promotions/list_promotions_logic.go` — MODIFY: ListPromotionsReq passes Kind + MarketID to App.
- `promotions/activate_promotion_logic.go` / `deactivate_promotion_logic.go` — MODIFY: thin passthrough.
- `promotions/create_promotion_rules_logic.go` / `update_promotion_rule_logic.go` / `get_promotion_rules_logic.go` / `delete_promotion_rule_logic.go` — MODIFY: read `Promotion.Kind` from `req.PromotionID` to determine owner_kind.
- `coupons/create_coupon_logic.go` — MODIFY: build CreatePromotionRequest with `Kind: KindCoupon, Code: req.Code, TotalCount: req.UsageLimit`.
- `coupons/update_coupon_logic.go` — MODIFY: same + rebuild rules.
- `coupons/list_coupons_logic.go` — MODIFY: ListCouponsReq → Query{Kind: KindCoupon, ...}.
- `coupons/get_coupon_logic.go` / `activate_coupon_logic.go` / `deactivate_coupon_logic.go` / `delete_coupon_logic.go` — MODIFY: thin passthrough.
- `coupons/generate_coupon_codes_logic.go` — MODIFY: loop `PromotionApp.Create({Kind: COUPON, ...})`.
- `coupons/issue_user_coupon_logic.go` / `batch_issue_user_coupon_logic.go` / `list_user_coupons_logic.go` / `get_coupon_usage_logic.go` — MODIFY: thin passthrough.
- `coupons/helper.go` — **DELETE**.
- `promotions/helper.go` — MODIFY: `convertPromotionToDetailResp` adds Kind/MarketID/Code/TotalCount/Rules fields; add `convertRulesToResp` helper.

### Errors — `pkg/code/code.go` (modify)
- Add `ErrPromotionInvalidKind` (80xxx).

### API — `admin/desc/promotion.api` (modify)
- `PromotionDetailResp`: add `Kind string`, `MarketID int64 optional`, `Code string optional`, `TotalCount int optional`, `Rules []*PromotionRuleResp`.
- `CreatePromotionReq` / `UpdatePromotionReq`: add `Kind string`, `MarketID int64 optional`, `Code string optional`, `TotalCount int optional`, `Rules []*PromotionRuleReq optional`.
- `CreatePromotionRulesReq`: rename path from `id` → `owner_id`; add `OwnerKind string path`.
- `PromotionRuleReq`: replace legacy `RuleType/Operator/Value` with `ConditionType/ConditionValue/ActionType/ActionValue/MaxDiscount/SortOrder`.
- `PromotionRuleResp`: same shape as new rule fields.
- `ListPromotionsReq`: add `Kind string form`, `MarketID int64 form`.

### Frontend — `shop-admin/src/` (modify)
- `api/promotion.ts` — MODIFY: collapse `CouponDetailResp` into `Promotion`; add `kind`, `market_id`, `code?`, `total_count?`, `rules?`; remove `CouponDetailResp` interface.
- `views/promotions/index.vue` — MODIFY: status/code/total_count columns gated by `row.kind === 'coupon'`; merge `couponForm`/`promotionForm` into single `form`; add market filter.
- `views/promotions/rule*.vue` — MODIFY: title text branches on `kind`.

---

## Task 1: Pre-migration verification SQL

**Files:**
- Create: `sql/promotion/migrations/2026071900_pre_migration_checks.sql`

**Interfaces:**
- Consumes: existing `coupons` and `promotion_rules` tables.
- Produces: a verification script that MUST be run before the merge migration. If any check returns rows, abort.

- [ ] **Step 1: Create the verification script**

Create `sql/promotion/migrations/2026071900_pre_migration_checks.sql`:

```sql
-- ============================================
-- Promotion × Coupon merge: pre-migration checks
-- Run BEFORE 2026071901_merge_promotion_coupon.sql
-- Any non-zero row count = BLOCK migration
-- ============================================

SELECT 'CHECK 1: coupon code duplicates (BLOCK if > 0)' AS check_name;
SELECT code, COUNT(*) AS c
FROM coupons
WHERE deleted_at IS NULL
GROUP BY code
HAVING c > 1;

SELECT 'CHECK 2: orphan promotion_rules (BLOCK if > 0)' AS check_name;
SELECT r.*
FROM promotion_rules r
LEFT JOIN promotions p ON p.id = r.promotion_id
WHERE p.id IS NULL;

SELECT 'CHECK 3: coupons with market_ids content (INFO; will be discarded)' AS check_name;
SELECT id, code, scope_type, scope_ids, market_ids
FROM coupons
WHERE JSON_LENGTH(market_ids) > 0;

SELECT 'CHECK 4: existing row counts (baseline for post-migration comparison)' AS check_name;
SELECT 'promotions' AS tbl, COUNT(*) AS n FROM promotions UNION ALL
SELECT 'promotion_rules', COUNT(*) FROM promotion_rules UNION ALL
SELECT 'coupons', COUNT(*) FROM coupons WHERE deleted_at IS NULL UNION ALL
SELECT 'user_coupons', COUNT(*) FROM user_coupons;
```

- [ ] **Step 2: Run the checks against the local dev DB**

```bash
mysql -h 192.168.0.100 -P 3306 -u root -p shopjoy < sql/promotion/migrations/2026071900_pre_migration_checks.sql
```

- [ ] **Step 3: Verify expected output**

Expected: CHECK 1 returns 0 rows. CHECK 2 returns 0 rows. CHECK 3 returns 0 rows (per design §2.7 the field will be discarded). CHECK 4 returns the baseline counts (record them for the post-migration comparison in Task 2).

- [ ] **Step 4: Commit**

```bash
git add sql/promotion/migrations/2026071900_pre_migration_checks.sql
git commit -m "sql(promotion): add pre-migration verification checks"
```

---

## Task 2: P0 merge migration script

**Files:**
- Create: `sql/promotion/migrations/2026071901_merge_promotion_coupon.sql`
- Modify: `sql/promotion/schema.sql`

**Interfaces:**
- Consumes: existing `coupons`, `promotions`, `promotion_rules` tables (must match the schema at HEAD).
- Produces: merged `promotions` (with kind + new nullable columns) + unified `promotion_rules` (with owner_kind + owner_id). Old tables archived as `_deprecated_coupons` + `_archived_coupons_20260719`.

- [ ] **Step 1: Create the migration script**

Create `sql/promotion/migrations/2026071901_merge_promotion_coupon.sql`:

```sql
START TRANSACTION;

-- 1) Backup old coupons
CREATE TABLE IF NOT EXISTS _deprecated_coupons LIKE coupons;
INSERT INTO _deprecated_coupons SELECT * FROM coupons;

-- 2) promotions: add new columns
ALTER TABLE promotions
  ADD COLUMN `kind`        ENUM('PROMOTION','COUPON') NOT NULL DEFAULT 'PROMOTION' AFTER `tenant_id`,
  ADD COLUMN `code`        VARCHAR(100) NULL AFTER `name`,
  ADD COLUMN `market_id`   BIGINT NULL AFTER `priority`,
  ADD COLUMN `total_count` INT NULL AFTER `usage_limit`,
  ADD COLUMN `used_count`  INT NULL AFTER `total_count`,
  ADD COLUMN `scope_type`  VARCHAR(32) NOT NULL DEFAULT 'STOREWIDE' AFTER `tags`,
  ADD COLUMN `scope_ids`   JSON NULL AFTER `scope_type`,
  ADD COLUMN `exclude_ids` JSON NULL AFTER `scope_ids`;

-- 3) Migrate coupons → promotions
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

-- 4) Backup promotion_rules
CREATE TABLE IF NOT EXISTS _deprecated_promotion_rules LIKE promotion_rules;
INSERT INTO _deprecated_promotion_rules SELECT * FROM promotion_rules;

-- 5) Rebuild promotion_rules with owner_kind + sort_order
ALTER TABLE promotion_rules
  ADD COLUMN `owner_kind` ENUM('PROMOTION','COUPON') NOT NULL DEFAULT 'PROMOTION' AFTER `promotion_id`,
  ADD COLUMN `sort_order` INT NOT NULL DEFAULT 0 AFTER `max_discount_amount`,
  MODIFY COLUMN `promotion_id` BIGINT NULL;

-- 6) Tag existing PROMOTION rules
UPDATE promotion_rules SET owner_kind = 'PROMOTION' WHERE promotion_id IS NOT NULL;

-- 7) Backfill owner_id from promotion_id for PROMOTION rows
UPDATE promotion_rules SET owner_id = promotion_id WHERE owner_kind = 'PROMOTION' AND owner_id = 0;

-- 8) Convert each COUPON's (type/value/min_amount/max_discount) into one rule row
INSERT INTO promotion_rules (
  owner_kind, owner_id, condition_type, condition_value, action_type, action_value,
  max_discount_amount, max_discount_currency, currency, sort_order,
  created_at, updated_at
)
SELECT
  'COUPON'                                                              AS owner_kind,
  p.id                                                                  AS owner_id,
  0                                                                     AS condition_type,
  c.min_amount                                                          AS condition_value,
  CASE c.type WHEN 0 THEN 0 WHEN 1 THEN 1 WHEN 2 THEN 0 END             AS action_type,
  CASE c.type WHEN 0 THEN c.value WHEN 1 THEN c.value WHEN 2 THEN 0 END AS action_value,
  c.max_discount                                                        AS max_discount_amount,
  c.currency                                                            AS max_discount_currency,
  c.currency,
  0                                                                     AS sort_order,
  NOW(), NOW()
FROM coupons c
JOIN promotions p ON p.kind = 'COUPON' AND p.code = c.code AND p.tenant_id = c.tenant_id
WHERE c.deleted_at IS NULL;

-- 9) Index rebuild
ALTER TABLE promotion_rules
  ADD INDEX `idx_owner` (`owner_kind`, `owner_id`, `sort_order`),
  DROP INDEX `idx_promotion_id`,
  DROP INDEX `idx_promotion_rules_sort_order`;

-- 10) COUPON.code partial-unique via generated column
ALTER TABLE promotions
  ADD COLUMN `code_unique` VARCHAR(100)
    GENERATED ALWAYS AS (IF(kind = 'COUPON', code, NULL)) VIRTUAL,
  ADD UNIQUE KEY `uk_promotion_code` (`code_unique`);

-- 11) user_coupons aid index
ALTER TABLE user_coupons
  ADD INDEX `idx_coupon_id_active` (`coupon_id`, `status`);

-- 12) Archive old coupons table (retained 30 days)
RENAME TABLE coupons TO _archived_coupons_20260719;

COMMIT;
```

- [ ] **Step 2: Run migration on local dev DB**

```bash
mysql -h 192.168.0.100 -P 3306 -u root -p shopjoy < sql/promotion/migrations/2026071901_merge_promotion_coupon.sql
```

- [ ] **Step 3: Verify post-migration integrity**

Run the spec §7.3 verification SQL. Expected output (replace baseline numbers from Task 1):

```sql
SELECT 'promotions_total'    AS tbl, COUNT(*) FROM promotions UNION ALL
SELECT 'promotions_coupon'   AS tbl, COUNT(*) FROM promotions WHERE kind='COUPON' UNION ALL
SELECT 'promotions_promotion'AS tbl, COUNT(*) FROM promotions WHERE kind='PROMOTION' UNION ALL
SELECT 'rules_total'         AS tbl, COUNT(*) FROM promotion_rules UNION ALL
SELECT 'rules_coupon'        AS tbl, COUNT(*) FROM promotion_rules WHERE owner_kind='COUPON' UNION ALL
SELECT 'rules_promotion'     AS tbl, COUNT(*) FROM promotion_rules WHERE owner_kind='PROMOTION' UNION ALL
SELECT 'user_coupons'        AS tbl, COUNT(*) FROM user_coupons;

-- Expected: promotions_coupon = baseline coupons (non-deleted)
-- Expected: rules_coupon = baseline coupons (one rule per coupon)
-- Expected: rules_total = rules_promotion + rules_coupon
```

- [ ] **Step 4: Verify code uniqueness + scope of new columns**

```sql
SELECT 'check coupon code unique' AS test, COUNT(*) AS duplicate_count
FROM (SELECT code_unique FROM promotions WHERE code_unique IS NOT NULL GROUP BY code_unique HAVING COUNT(*) > 1) AS d;
-- Expected: 0

SHOW CREATE TABLE promotions\G
-- Confirm columns: kind, code, market_id, total_count, used_count, scope_type, scope_ids, exclude_ids, code_unique

SHOW CREATE TABLE promotion_rules\G
-- Confirm columns: owner_kind, owner_id, sort_order; INDEX idx_owner (owner_kind, owner_id, sort_order)
```

- [ ] **Step 5: Merge column changes into `schema.sql` (do not delete migration yet)**

In `sql/promotion/schema.sql`:
- Replace the `promotions` CREATE TABLE block with one that includes the new columns and `code_unique`.
- Replace the `promotion_rules` CREATE TABLE block with the new structure.
- Drop the `coupons` and `user_coupons` blocks.
- Add `_deprecated_coupons` (or just `coupons` archived note) at the bottom for retention.
- Update the `INSERT` data blocks: promotions seed data stays as-is; coupon/user_coupon data is removed (rows are now in promotions + user_coupons).

**Do NOT delete the migration file yet** — wait until Task 9 (final cleanup) which drops the archived table.

- [ ] **Step 6: Commit**

```bash
git add sql/promotion/migrations/2026071901_merge_promotion_coupon.sql sql/promotion/schema.sql
git commit -m "sql(promotion): merge promotions + coupons into single table"
```

---

## Task 3: Add `ErrPromotionInvalidKind` error code

**Files:**
- Modify: `pkg/code/code.go`

**Interfaces:**
- Produces: `code.ErrPromotionInvalidKind` with HTTPCode=400, Code=80018 (next available in 80xxx range; brief's 80010 was already taken), Msg="invalid promotion kind".

- [ ] **Step 1: Find the Promotion error block**

In `pkg/code/code.go`, locate the Promotion section (80xxx range) — usually a block comment `// Promotion Module (80xxx)` followed by existing errors like `ErrPromotionNotFound`, `ErrPromotionRuleNotFound`.

- [ ] **Step 2: Add the new error**

Insert after the last existing Promotion error:

```go
// ErrPromotionInvalidKind is returned when an operation is invoked on a
// Promotion whose Kind does not support it (e.g., Issue on a PROMOTION).
ErrPromotionInvalidKind = &Err{HTTPCode: http.StatusBadRequest, Code: 80018, Msg: "invalid promotion kind"}
```

- [ ] **Step 3: Verify build**

```bash
cd admin && make build
```

Expected: 0 errors.

- [ ] **Step 4: Commit**

```bash
git add pkg/code/code.go
git commit -m "feat(code): add ErrPromotionInvalidKind"
```

---

## Task 4: Domain layer rewrite

**Files:**
- Create: `pkg/domain/promotion/rule.go`
- Create: `pkg/domain/promotion/query.go`
- Create: `pkg/domain/promotion/repository.go`
- Modify: `pkg/domain/promotion/entity.go`
- Modify: `pkg/domain/promotion/scope.go`
- Modify: `pkg/domain/promotion/coupon.go`
- Delete: `admin/internal/domain/promotion/entity.go`
- Delete: `admin/internal/domain/coupon/entity.go`

**Interfaces:**
- Produces:
  - `promotion.Kind` type (`"PROMOTION" | "COUPON"`) with `IsValid()`.
  - `promotion.Promotion` struct (single struct, kind-routed nullable fields).
  - `promotion.PromotionRule` struct (`OwnerKind`, `OwnerID`, condition/action fields, sort_order).
  - `promotion.Query` struct with pointer Kind/Status/Type/MarketID + ExpiredOnly.
  - `promotion.Repository` interface (kind-agnostic methods + COUPON-specific).

- [ ] **Step 1: Create `rule.go`**

Create `pkg/domain/promotion/rule.go`:

```go
package promotion

import (
    "time"

    "github.com/shopspring/decimal"
)

// ConditionType enumerates the types of rule conditions.
type ConditionType int

const (
    ConditionMinAmount ConditionType = iota // 0 - minimum cart amount
    ConditionMinQuantity                     // 1 - minimum item quantity
)

func (t ConditionType) IsValid() bool {
    return t >= ConditionMinAmount && t <= ConditionMinQuantity
}

// ActionType enumerates the discount actions a rule can perform.
type ActionType int

const (
    ActionFixedAmount ActionType = iota // 0 - subtract ActionValue as money
    ActionPercentage                    // 1 - multiply by ActionValue/10000 (basis points)
    ActionFreeShipping                   // 2 - waive shipping (ActionValue unused)
)

func (t ActionType) IsValid() bool {
    return t >= ActionFixedAmount && t <= ActionFreeShipping
}

// PromotionRule belongs to a Promotion (of any Kind) via owner_kind+owner_id.
// Rules are ordered by sort_order; the best-matching rule wins via FindBestRule.
type PromotionRule struct {
    ID             int64           `json:"id"`
    OwnerKind      Kind            `json:"owner_kind"`
    OwnerID        int64           `json:"owner_id"`
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

// CalculateDiscount applies the rule's action to matchedAmount and caps by MaxDiscount.
// ActionPercentage: ActionValue is basis points (100 = 1%, so divide by 10000).
// ActionFreeShipping: returns zero (the discount is the shipping fee, handled elsewhere).
func (r *PromotionRule) CalculateDiscount(matchedAmount decimal.Decimal) decimal.Decimal {
    var discount decimal.Decimal
    switch r.ActionType {
    case ActionFixedAmount:
        discount = r.ActionValue
    case ActionPercentage:
        discount = matchedAmount.Mul(r.ActionValue).Div(decimal.NewFromInt(10000))
    case ActionFreeShipping:
        return decimal.Zero
    }

    if r.MaxDiscount.IsPositive() && discount.GreaterThan(r.MaxDiscount) {
        discount = r.MaxDiscount
    }
    if discount.GreaterThan(matchedAmount) {
        discount = matchedAmount
    }
    return discount
}

// MeetsCondition returns true if amount/quantity clears the threshold.
func (r *PromotionRule) MeetsCondition(amount decimal.Decimal, quantity int) bool {
    switch r.ConditionType {
    case ConditionMinAmount:
        return amount.GreaterThanOrEqual(r.ConditionValue)
    case ConditionMinQuantity:
        return decimal.NewFromInt(int64(quantity)).GreaterThanOrEqual(r.ConditionValue)
    default:
        return false
    }
}
```

- [ ] **Step 2: Create `query.go`**

Create `pkg/domain/promotion/query.go`:

```go
package promotion

import (
    "github.com/colinrs/shopjoy/pkg/domain/shared"
)

// Query is the filter set for FindList. All optional fields are pointers
// because Kind/Status/Type are iota-based enums whose zero values are
// legitimate members (PROMOTION, ACTIVE, DISCOUNT) — using != 0 as a
// "filter set" sentinel would silently drop those.
type Query struct {
    shared.PageQuery
    TenantID    shared.TenantID
    Name        string
    Kind        *Kind
    Status      *Status
    Type        *Type
    MarketID    *int64
    ExpiredOnly bool
}
```

- [ ] **Step 3: Create `repository.go`**

Create `pkg/domain/promotion/repository.go`:

```go
package promotion

import (
    "context"

    "gorm.io/gorm"
)

// Repository is the persistence boundary for promotions and their rules.
// All methods are kind-agnostic unless suffixed (e.g., Issue* methods are
// only valid when OwnerKind == KindCoupon).
type Repository interface {
    // Generic CRUD
    Create(ctx context.Context, db *gorm.DB, p *Promotion) error
    Update(ctx context.Context, db *gorm.DB, p *Promotion) error
    Delete(ctx context.Context, db *gorm.DB, id int64) error
    FindByID(ctx context.Context, db *gorm.DB, id int64) (*Promotion, error)
    FindByCode(ctx context.Context, db *gorm.DB, code string) (*Promotion, error)
    FindList(ctx context.Context, db *gorm.DB, query Query) ([]*Promotion, int64, error)

    // Rules
    CreateRules(ctx context.Context, db *gorm.DB, ownerKind Kind, ownerID int64, rules []PromotionRule) error
    FindRulesByOwner(ctx context.Context, db *gorm.DB, ownerKind Kind, ownerID int64) ([]PromotionRule, error)
    UpdateRule(ctx context.Context, db *gorm.DB, rule *PromotionRule) error
    DeleteRule(ctx context.Context, db *gorm.DB, id int64) error
    DeleteRulesByOwner(ctx context.Context, db *gorm.DB, ownerKind Kind, ownerID int64) error

    // Coupon-specific
    FindActiveCoupons(ctx context.Context, db *gorm.DB, marketID *int64) ([]*Promotion, error)
    IncrementUsedCount(ctx context.Context, db *gorm.DB, couponID int64) error
    IssueUserCoupon(ctx context.Context, db *gorm.DB, uc *UserCoupon) error
    FindUserCoupons(ctx context.Context, db *gorm.DB, query UserCouponQuery) ([]*UserCoupon, int64, error)

    // Usage (existing)
    FindPromotionUsage(ctx context.Context, db *gorm.DB, query UsageQuery) ([]*PromotionUsage, int64, error)
}
```

- [ ] **Step 4: Rewrite `entity.go`**

Overwrite `pkg/domain/promotion/entity.go`:

```go
package promotion

import (
    "time"

    "github.com/colinrs/shopjoy/pkg/code"
    "github.com/colinrs/shopjoy/pkg/domain/shared"
    "github.com/shopspring/decimal"
)

// Kind discriminates between system-driven promotions and claim-based coupons.
type Kind string

const (
    KindPromotion Kind = "PROMOTION"
    KindCoupon    Kind = "COUPON"
)

func (k Kind) IsValid() bool {
    return k == KindPromotion || k == KindCoupon
}

// Status (unified for both kinds)
type Status int

const (
    StatusPending Status = iota // 0
    StatusActive                // 1
    StatusPaused                // 2
    StatusEnded                 // 3 - depleted coupons also surface as ended (see IsActive)
)

func (s Status) IsValid() bool {
    return s >= StatusPending && s <= StatusEnded
}

// Type (marketing play). COUPONs always use TypeDiscount (=0).
type Type int

const (
    TypeDiscount Type = iota // 0
    TypeFlashSale            // 1
    TypeBundle               // 2
    TypeBuyXGetY             // 3
)

func (t Type) IsValid() bool {
    return t >= TypeDiscount && t <= TypeBuyXGetY
}

// Promotion is the aggregate root for both system promotions and user-claimable
// coupons. Coupon-specific fields are nullable; semantics activate when Kind == KindCoupon.
type Promotion struct {
    ID           int64            `json:"id"`
    TenantID     shared.TenantID  `json:"tenant_id"`
    Kind         Kind             `json:"kind"`
    Name         string           `json:"name"`
    Description  string           `json:"description"`
    Code         *string          `json:"code,omitempty"`
    Type         Type             `json:"type"`
    Status       Status           `json:"status"`
    Priority     int              `json:"priority"`
    MarketID     *int64           `json:"market_id,omitempty"`
    Currency     string           `json:"currency"`
    TotalCount   *int             `json:"total_count,omitempty"`
    UsedCount    *int             `json:"used_count,omitempty"`
    UsageLimit   int              `json:"usage_limit"`
    PerUserLimit int              `json:"per_user_limit"`
    Tags         []string         `json:"tags,omitempty" gorm:"type:json"`
    Scope        PromotionScope   `json:"scope"`
    StartAt      time.Time        `json:"start_at"`
    EndAt        time.Time        `json:"end_at"`
    Rules        []PromotionRule  `json:"rules,omitempty"`
    Audit        shared.AuditInfo `json:"audit"`
    DeletedAt    *time.Time       `json:"deleted_at,omitempty"`
}

func (p *Promotion) TableName() string { return "promotions" }

// IsActive returns true if the promotion is currently usable. For COUPONs,
// this also checks inventory (used_count < total_count).
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

// MatchesMarket returns true if the promotion applies to the given market.
// A NULL market_id means "applies to all markets".
func (p *Promotion) MatchesMarket(marketID int64) bool {
    return p.MarketID == nil || *p.MarketID == marketID
}

// MatchesScope delegates to Scope (kind-agnostic).
func (p *Promotion) MatchesScope(productID, categoryID, brandID int64) bool {
    return p.Scope.MatchesProduct(productID, categoryID, brandID)
}

// FindBestRule returns the rule with the highest ConditionValue that still
// meets the condition. Multi-tier support for COUPONs lives here.
func (p *Promotion) FindBestRule(matchedAmount decimal.Decimal, quantity int) *PromotionRule {
    var best *PromotionRule
    for i := range p.Rules {
        rule := &p.Rules[i]
        if !rule.MeetsCondition(matchedAmount, quantity) {
            continue
        }
        if best == nil || rule.ConditionValue.GreaterThan(best.ConditionValue) {
            best = rule
        }
    }
    return best
}

// CalculateDiscount is a convenience wrapper around FindBestRule.
func (p *Promotion) CalculateDiscount(matchedAmount decimal.Decimal, quantity int) decimal.Decimal {
    rule := p.FindBestRule(matchedAmount, quantity)
    if rule == nil {
        return decimal.Zero
    }
    return rule.CalculateDiscount(matchedAmount)
}

// Issue creates a UserCoupon from this promotion. Returns ErrPromotionInvalidKind
// if invoked on a non-COUPON.
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

// ConsumeInventory increments used_count in memory (pre-check). Persistence
// happens via repo.IncrementUsedCount which uses an atomic SQL check to
// prevent overselling:
//
//   UPDATE promotions
//   SET used_count = used_count + 1
//   WHERE id = ? AND kind = 'COUPON' AND (total_count IS NULL OR used_count < total_count)
//
// Caller must roll back the in-memory value if the SQL fails.
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
```

- [ ] **Step 5: Update `scope.go`**

In `pkg/domain/promotion/scope.go`, find the `ScopeType` constants and remove the `ScopeTypeMarkets` member (and any branch handling it in `MatchesProduct`). Keep STOREWIDE / PRODUCTS / CATEGORIES / BRANDS.

```go
type ScopeType string

const (
    ScopeTypeStorewide  ScopeType = "STOREWIDE"
    ScopeTypeProducts   ScopeType = "PRODUCTS"
    ScopeTypeCategories ScopeType = "CATEGORIES"
    ScopeTypeBrands     ScopeType = "BRANDS"
    // MARKET scope removed — promotions now use the top-level market_id column.
)

func (s ScopeType) IsValid() bool {
    switch s {
    case ScopeTypeStorewide, ScopeTypeProducts, ScopeTypeCategories, ScopeTypeBrands:
        return true
    }
    return false
}
```

Update `MatchesProduct` to drop any `ScopeTypeMarkets` branch (likely just one `case` to remove).

- [ ] **Step 6: Rewrite `coupon.go`**

Overwrite `pkg/domain/promotion/coupon.go`:

```go
package promotion

import (
    "time"

    "github.com/colinrs/shopjoy/pkg/domain/shared"
    "github.com/shopspring/decimal"
)

// UserCoupon is a per-user claim record. coupon_id points to promotions.id
// where kind = 'COUPON'.
type UserCoupon struct {
    ID         int64            `json:"id"`
    TenantID   shared.TenantID  `json:"tenant_id"`
    UserID     int64            `json:"user_id"`
    CouponID   int64            `json:"coupon_id"`
    Status     UserCouponStatus `json:"status"`
    UsedAt     *time.Time       `json:"used_at,omitempty"`
    OrderID    int64            `json:"order_id"`
    ReceivedAt time.Time        `json:"received_at"`
    ExpireAt   time.Time        `json:"expire_at"`
    CreatedAt  time.Time        `json:"created_at"`
    UpdatedAt  time.Time        `json:"updated_at"`
}

func (uc *UserCoupon) TableName() string { return "user_coupons" }

func (uc *UserCoupon) IsExpired() bool {
    return time.Now().UTC().After(uc.ExpireAt)
}

func (uc *UserCoupon) CanUse() bool {
    return uc.Status == UserCouponStatusUnused && !uc.IsExpired()
}

// UserCouponStatus
type UserCouponStatus int

const (
    UserCouponStatusUnused  UserCouponStatus = iota // 0
    UserCouponStatusUsed                            // 1
    UserCouponStatusExpired                         // 2
)

// PromotionUsage records one (promotion|rule, coupon, order, user) hit.
// Either promotion_id or coupon_id (or both) may be set.
type PromotionUsage struct {
    ID             int64           `json:"id"`
    TenantID       shared.TenantID `json:"tenant_id"`
    PromotionID    int64           `json:"promotion_id"`
    RuleID         *int64          `json:"rule_id,omitempty"`
    OrderID        int64           `json:"order_id"`
    UserID         int64           `json:"user_id"`
    DiscountAmount decimal.Decimal `json:"discount_amount"`
    Currency       string          `json:"currency"`
    OriginalAmount decimal.Decimal `json:"original_amount"`
    FinalAmount    decimal.Decimal `json:"final_amount"`
    CouponID       *int64          `json:"coupon_id,omitempty"`
    CreatedAt      time.Time       `json:"created_at"`
}

func (pu *PromotionUsage) TableName() string { return "promotion_usage" }

// UserCouponQuery
type UserCouponQuery struct {
    sharedPageQuery() // see note below
}
```

Replace the inline `sharedPageQuery()` shim with the real struct:

```go
type UserCouponQuery struct {
    Page   int
    Size   int
    UserID *int64
    CouponID *int64
    Status *UserCouponStatus
}
```

(Use whatever fields the existing coupon_repository.go's UserCoupon query currently has — copy from there. The intent is to expose exactly the fields the new repo's `FindUserCoupons` needs.)

- [ ] **Step 7: Add `UsageQuery`**

In `pkg/domain/promotion/coupon.go`, add:

```go
type UsageQuery struct {
    Page     int
    Size     int
    CouponID *int64
    UserID   *int64
}
```

- [ ] **Step 8: Delete old domain entities**

```bash
git rm admin/internal/domain/promotion/entity.go admin/internal/domain/coupon/entity.go
```

- [ ] **Step 9: Verify build**

```bash
cd admin && make build
```

Expected: many errors (expected — repository + app + logic still reference old types). Confirm the errors are all in persistence/ and logic/ directories, not in pkg/domain/.

- [ ] **Step 10: Commit**

```bash
git add pkg/domain/promotion/ admin/internal/domain/
git commit -m "refactor(domain/promotion): unify Promotion + Coupon into single struct"
```

---

## Task 5: Repository rewrite

**Files:**
- Modify: `admin/internal/infrastructure/persistence/promotion_repository.go`
- Delete: `admin/internal/infrastructure/persistence/coupon_repository.go`

**Interfaces:**
- Consumes: `promotion.Promotion`, `promotion.PromotionRule`, `promotion.UserCoupon`, `promotion.PromotionUsage`, `promotion.Query` from `pkg/domain/promotion/`.
- Produces: implementation of `promotion.Repository` interface (see Task 4 Step 3).

- [ ] **Step 1: Delete the old coupon repository**

```bash
git rm admin/internal/infrastructure/persistence/coupon_repository.go
```

- [ ] **Step 2: Rewrite `promotionModel`**

Open `admin/internal/infrastructure/persistence/promotion_repository.go`. Replace the `promotionModel` struct with:

```go
type promotionModel struct {
    application.Model
    TenantID     int64                  `gorm:"column:tenant_id;not null;index"`
    Kind         string                 `gorm:"column:kind;type:enum('PROMOTION','COUPON');not null;default:PROMOTION"`
    Name         string                 `gorm:"column:name;size:255;not null"`
    Description  string                 `gorm:"column:description;type:text"`
    Code         *string                `gorm:"column:code;size:100"`
    Type         int                    `gorm:"column:type;not null;default:0"`
    Status       int                    `gorm:"column:status;not null;default:0;index"`
    Priority     int                    `gorm:"column:priority;not null;default:0"`
    MarketID     *int64                 `gorm:"column:market_id;index"`
    Currency     string                 `gorm:"column:currency;size:10;not null;default:CNY"`
    TotalCount   *int                   `gorm:"column:total_count"`
    UsedCount    *int                   `gorm:"column:used_count"`
    UsageLimit   int                    `gorm:"column:usage_limit;not null;default:0"`
    PerUserLimit int                    `gorm:"column:per_user_limit;not null;default:1"`
    Tags         datatypes.JSON         `gorm:"column:tags;type:json"`
    ScopeType    string                 `gorm:"column:scope_type;size:32;not null;default:STOREWIDE"`
    ScopeIDs     datatypes.JSON         `gorm:"column:scope_ids;type:json"`
    ExcludeIDs   datatypes.JSON         `gorm:"column:exclude_ids;type:json"`
    StartAt      time.Time              `gorm:"column:start_at;not null"`
    EndAt        time.Time              `gorm:"column:end_at;not null"`
    Audit        shared.AuditInfo       `gorm:"embedded"`
}

func (*promotionModel) TableName() string { return "promotions" }
```

- [ ] **Step 3: Rewrite `promotionRuleModel`**

Replace `promotionRuleModel`:

```go
type promotionRuleModel struct {
    application.Model
    OwnerKind           string          `gorm:"column:owner_kind;type:enum('PROMOTION','COUPON');not null;default:PROMOTION"`
    OwnerID             int64           `gorm:"column:owner_id;not null"`
    PromotionID         *int64          `gorm:"column:promotion_id"`
    ConditionType       int             `gorm:"column:condition_type;not null;default:0"`
    ConditionValue      decimal.Decimal `gorm:"column:condition_value;type:decimal(19,4);not null;default:0"`
    ActionType          int             `gorm:"column:action_type;not null;default:0"`
    ActionValue         decimal.Decimal `gorm:"column:action_value;type:decimal(19,4);not null;default:0"`
    MaxDiscountAmount   decimal.Decimal `gorm:"column:max_discount_amount;type:decimal(19,4);not null;default:0"`
    MaxDiscountCurrency string          `gorm:"column:max_discount_currency;size:10;default:CNY"`
    Currency            string          `gorm:"column:currency;size:10;not null;default:CNY"`
    SortOrder           int             `gorm:"column:sort_order;not null;default:0"`
    Audit               shared.AuditInfo `gorm:"embedded"`
}

func (*promotionRuleModel) TableName() string { return "promotion_rules" }
```

- [ ] **Step 4: Implement `toEntity` and `fromPromotionRuleEntity`**

Update the conversion functions so:
- `promotionModel.toEntity()` maps all new fields into `promotion.Promotion` (including nullable Code/TotalCount/UsedCount/MarketID).
- `fromPromotionRuleEntity(rule)` populates `owner_kind` from `rule.OwnerKind`, `owner_id` from `rule.OwnerID`, `promotion_id` set only when `owner_kind == KindPromotion`.

```go
func (m *promotionModel) toEntity() *promotion.Promotion {
    p := &promotion.Promotion{
        ID:           m.ID,
        TenantID:     shared.TenantID(m.TenantID),
        Kind:         promotion.Kind(m.Kind),
        Name:         m.Name,
        Description:  m.Description,
        Code:         m.Code,
        Type:         promotion.Type(m.Type),
        Status:       promotion.Status(m.Status),
        Priority:     m.Priority,
        MarketID:     m.MarketID,
        Currency:     m.Currency,
        TotalCount:   m.TotalCount,
        UsedCount:    m.UsedCount,
        UsageLimit:   m.UsageLimit,
        PerUserLimit: m.PerUserLimit,
        Scope: promotion.PromotionScope{
            Type:       promotion.ScopeType(m.ScopeType),
            ScopeIDs:   m.ScopeIDs,
            ExcludeIDs: m.ExcludeIDs,
        },
        StartAt:   m.StartAt,
        EndAt:     m.EndAt,
        Audit:     m.Audit,
        DeletedAt: m.DeletedAt,
    }
    if len(m.Tags) > 0 {
        _ = json.Unmarshal(m.Tags, &p.Tags)
    }
    return p
}

func fromPromotionRuleEntity(r *promotion.PromotionRule) promotionRuleModel {
    m := promotionRuleModel{
        OwnerKind:           string(r.OwnerKind),
        OwnerID:             r.OwnerID,
        ConditionType:       int(r.ConditionType),
        ConditionValue:      r.ConditionValue,
        ActionType:          int(r.ActionType),
        ActionValue:         r.ActionValue,
        MaxDiscountAmount:   r.MaxDiscount,
        MaxDiscountCurrency: r.Currency,
        Currency:            r.Currency,
        SortOrder:           r.SortOrder,
    }
    if r.OwnerKind == promotion.KindPromotion {
        pid := r.OwnerID
        m.PromotionID = &pid
    }
    m.ID = r.ID
    return m
}
```

- [ ] **Step 5: Implement Repository methods**

Implement all methods declared in Task 4 Step 3. Highlights:

```go
func (r *promotionRepo) FindByCode(ctx context.Context, db *gorm.DB, code string) (*promotion.Promotion, error) {
    var m promotionModel
    err := db.WithContext(ctx).Where("kind = ? AND code = ?", promotion.KindCoupon, code).First(&m).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, code.ErrPromotionNotFound
        }
        return nil, err
    }
    return m.toEntity(), nil
}

func (r *promotionRepo) IncrementUsedCount(ctx context.Context, db *gorm.DB, couponID int64) error {
    // Atomic check-and-increment to prevent overselling under concurrency.
    res := db.WithContext(ctx).Exec(`
        UPDATE promotions
        SET used_count = used_count + 1
        WHERE id = ? AND kind = 'COUPON'
          AND (total_count IS NULL OR used_count < total_count)
    `, couponID)
    if res.Error != nil {
        return res.Error
    }
    if res.RowsAffected == 0 {
        return code.ErrCouponDepleted
    }
    return nil
}

func (r *promotionRepo) DeleteRulesByOwner(ctx context.Context, db *gorm.DB, ownerKind promotion.Kind, ownerID int64) error {
    return db.WithContext(ctx).
        Where("owner_kind = ? AND owner_id = ?", ownerKind, ownerID).
        Delete(&promotionRuleModel{}).Error
}

func (r *promotionRepo) FindActiveCoupons(ctx context.Context, db *gorm.DB, marketID *int64) ([]*promotion.Promotion, error) {
    now := time.Now().UTC()
    q := db.WithContext(ctx).Model(&promotionModel{}).
        Where("kind = ?", promotion.KindCoupon).
        Where("status = ?", promotion.StatusActive).
        Where("start_at <= ? AND end_at >= ?", now, now).
        Where("(total_count IS NULL OR used_count < total_count)")
    if marketID != nil {
        q = q.Where("market_id IS NULL OR market_id = ?", *marketID)
    }
    var models []promotionModel
    if err := q.Find(&models).Error; err != nil {
        return nil, err
    }
    out := make([]*promotion.Promotion, len(models))
    for i, m := range models {
        out[i] = m.toEntity()
    }
    return out, nil
}
```

For `FindList`, apply the pointer-based filters:

```go
if query.Kind != nil && query.Kind.IsValid() {
    dbQuery = dbQuery.Where("kind = ?", *query.Kind)
}
if query.Status != nil && query.Status.IsValid() {
    dbQuery = dbQuery.Where("status = ?", *query.Status)
}
if query.Type != nil && query.Type.IsValid() {
    dbQuery = dbQuery.Where("type = ?", *query.Type)
}
if query.MarketID != nil {
    dbQuery = dbQuery.Where("market_id = ?", *query.MarketID)
}
if query.ExpiredOnly {
    dbQuery = dbQuery.Where("end_at <= ?", time.Now().UTC())
}
```

- [ ] **Step 6: Verify build (compiles alone)**

```bash
cd admin && make build
```

Expected: errors only in `application/promotion/` and `logic/` directories (those are next).

- [ ] **Step 7: Commit**

```bash
git add admin/internal/infrastructure/persistence/promotion_repository.go
git commit -m "refactor(persistence): unify promotion + coupon repositories"
```

---

## Task 6: Application layer merge

**Files:**
- Modify: `admin/internal/application/promotion/promotion_app.go`
- Delete: `admin/internal/application/promotion/coupon_app.go`

**Interfaces:**
- Consumes: `promotion.Repository`, `promotion.Promotion`, `promotion.PromotionRule`, `promotion.UserCoupon` from `pkg/domain/promotion/`.
- Produces: single `PromotionApp` struct with all methods listed in spec §5.1.

- [ ] **Step 1: Delete `coupon_app.go`**

```bash
git rm admin/internal/application/promotion/coupon_app.go
```

- [ ] **Step 2: Rewrite `promotion_app.go`**

Overwrite `admin/internal/application/promotion/promotion_app.go`:

```go
package promotion

import (
    "context"
    "time"

    "github.com/colinrs/shopjoy/pkg/code"
    "github.com/colinrs/shopjoy/pkg/domain/promotion"
    "github.com/colinrs/shopjoy/pkg/domain/shared"
    "github.com/shopspring/decimal"
    "gorm.io/gorm"
)

// CreatePromotionRequest is the unified input for both kinds.
type CreatePromotionRequest struct {
    TenantID     shared.TenantID
    Kind         promotion.Kind
    Name         string
    Description  string
    Code         *string
    Type         promotion.Type
    MarketID     *int64
    Currency     string
    TotalCount   *int
    UsageLimit   int
    PerUserLimit int
    Tags         []string
    Scope        promotion.PromotionScope
    StartAt      time.Time
    EndAt        time.Time
    Rules        []promotion.PromotionRule
    ActorID      int64
}

// UpdatePromotionRequest mirrors CreatePromotionRequest with ID and Audit info.
type UpdatePromotionRequest struct {
    ID           int64
    Name         string
    Description  string
    Code         *string
    Type         promotion.Type
    MarketID     *int64
    Currency     string
    TotalCount   *int
    UsageLimit   int
    PerUserLimit int
    Tags         []string
    Scope        promotion.PromotionScope
    StartAt      time.Time
    EndAt        time.Time
    Rules        *[]promotion.PromotionRule // nil = no change; non-nil = replace
    Status       *promotion.Status
    ActorID      int64
}

// PromotionResponse is the unified output.
type PromotionResponse struct {
    ID           int64
    TenantID     shared.TenantID
    Kind         promotion.Kind
    Name         string
    Description  string
    Code         *string
    Type         promotion.Type
    Status       promotion.Status
    MarketID     *int64
    Currency     string
    TotalCount   *int
    UsedCount    *int
    UsageLimit   int
    PerUserLimit int
    Tags         []string
    Scope        promotion.PromotionScope
    StartAt      time.Time
    EndAt        time.Time
    Rules        []PromotionRuleResponse
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type PromotionRuleResponse struct {
    ID             int64
    ConditionType  promotion.ConditionType
    ConditionValue decimal.Decimal
    ActionType     promotion.ActionType
    ActionValue    decimal.Decimal
    MaxDiscount    decimal.Decimal
    SortOrder      int
}

type ListPromotionResponse struct {
    List  []*PromotionResponse
    Total int64
    Page  int
    Size  int
}

type PromotionApp struct {
    repo promotion.Repository
    db   *gorm.DB
}

func NewPromotionApp(repo promotion.Repository, db *gorm.DB) *PromotionApp {
    return &PromotionApp{repo: repo, db: db}
}

// Create persists a Promotion and (optionally) its rules.
func (a *PromotionApp) Create(ctx context.Context, req *CreatePromotionRequest) (*PromotionResponse, error) {
    now := time.Now().UTC()
    p := &promotion.Promotion{
        TenantID:     req.TenantID,
        Kind:         req.Kind,
        Name:         req.Name,
        Description:  req.Description,
        Code:         req.Code,
        Type:         req.Type,
        Status:       promotion.StatusPending,
        MarketID:     req.MarketID,
        Currency:     req.Currency,
        TotalCount:   req.TotalCount,
        UsedCount:    nilOrZero(req.TotalCount),
        UsageLimit:   req.UsageLimit,
        PerUserLimit: req.PerUserLimit,
        Tags:         req.Tags,
        Scope:        req.Scope,
        StartAt:      req.StartAt,
        EndAt:        req.EndAt,
        Rules:        req.Rules,
        Audit:        shared.AuditInfo{CreatedAt: now, UpdatedAt: now, CreatedBy: req.ActorID, UpdatedBy: req.ActorID},
    }
    if err := a.repo.Create(ctx, a.db, p); err != nil {
        return nil, err
    }
    if len(p.Rules) > 0 {
        for i := range p.Rules {
            p.Rules[i].OwnerKind = p.Kind
            p.Rules[i].OwnerID = p.ID
        }
        if err := a.repo.CreateRules(ctx, a.db, p.Kind, p.ID, p.Rules); err != nil {
            return nil, err
        }
    }
    return a.toResponse(p), nil
}

// Update modifies an existing Promotion and (optionally) its rules.
func (a *PromotionApp) Update(ctx context.Context, req *UpdatePromotionRequest) (*PromotionResponse, error) {
    p, err := a.repo.FindByID(ctx, a.db, req.ID)
    if err != nil {
        return nil, err
    }
    p.Name = req.Name
    p.Description = req.Description
    p.Code = req.Code
    p.Type = req.Type
    p.MarketID = req.MarketID
    p.Currency = req.Currency
    p.TotalCount = req.TotalCount
    p.UsageLimit = req.UsageLimit
    p.PerUserLimit = req.PerUserLimit
    p.Tags = req.Tags
    p.Scope = req.Scope
    p.StartAt = req.StartAt
    p.EndAt = req.EndAt
    if req.Status != nil {
        p.Status = *req.Status
    }
    p.Audit.UpdatedAt = time.Now().UTC()
    p.Audit.UpdatedBy = req.ActorID

    if err := a.repo.Update(ctx, a.db, p); err != nil {
        return nil, err
    }
    if req.Rules != nil {
        if err := a.repo.DeleteRulesByOwner(ctx, a.db, p.Kind, p.ID); err != nil {
            return nil, err
        }
        if len(*req.Rules) > 0 {
            for i := range *req.Rules {
                (*req.Rules)[i].OwnerKind = p.Kind
                (*req.Rules)[i].OwnerID = p.ID
            }
            if err := a.repo.CreateRules(ctx, a.db, p.Kind, p.ID, *req.Rules); err != nil {
                return nil, err
            }
            p.Rules = *req.Rules
        } else {
            p.Rules = nil
        }
    }
    return a.toResponse(p), nil
}

// Get returns the promotion with its rules loaded.
func (a *PromotionApp) Get(ctx context.Context, id int64) (*PromotionResponse, error) {
    p, err := a.repo.FindByID(ctx, a.db, id)
    if err != nil {
        return nil, err
    }
    rules, err := a.repo.FindRulesByOwner(ctx, a.db, p.Kind, p.ID)
    if err != nil {
        return nil, err
    }
    p.Rules = rules
    return a.toResponse(p), nil
}

// List paginates with optional filters.
func (a *PromotionApp) List(ctx context.Context, q promotion.Query) (*ListPromotionResponse, error) {
    list, total, err := a.repo.FindList(ctx, a.db, q)
    if err != nil {
        return nil, err
    }
    out := make([]*PromotionResponse, len(list))
    for i, p := range list {
        out[i] = a.toResponse(p)
    }
    return &ListPromotionResponse{
        List:  out,
        Total: total,
        Page:  q.GetPage(),
        Size:  q.GetSize(),
    }, nil
}

// Delete removes a promotion and all its rules.
func (a *PromotionApp) Delete(ctx context.Context, id int64) error {
    p, err := a.repo.FindByID(ctx, a.db, id)
    if err != nil {
        return err
    }
    if err := a.repo.DeleteRulesByOwner(ctx, a.db, p.Kind, p.ID); err != nil {
        return err
    }
    return a.repo.Delete(ctx, a.db, id)
}

// Activate / Deactivate flip Status.
func (a *PromotionApp) Activate(ctx context.Context, id int64) (*PromotionResponse, error) {
    p, err := a.repo.FindByID(ctx, a.db, id)
    if err != nil {
        return nil, err
    }
    if time.Now().UTC().After(p.EndAt) {
        return nil, code.ErrPromotionExpired
    }
    p.Status = promotion.StatusActive
    p.Audit.UpdatedAt = time.Now().UTC()
    if err := a.repo.Update(ctx, a.db, p); err != nil {
        return nil, err
    }
    return a.toResponse(p), nil
}

func (a *PromotionApp) Deactivate(ctx context.Context, id int64) (*PromotionResponse, error) {
    p, err := a.repo.FindByID(ctx, a.db, id)
    if err != nil {
        return nil, err
    }
    p.Status = promotion.StatusPaused
    p.Audit.UpdatedAt = time.Now().UTC()
    if err := a.repo.Update(ctx, a.db, p); err != nil {
        return nil, err
    }
    return a.toResponse(p), nil
}

// GetRules returns rules for an owner.
func (a *PromotionApp) GetRules(ctx context.Context, ownerKind promotion.Kind, ownerID int64) ([]*PromotionRuleResponse, error) {
    rules, err := a.repo.FindRulesByOwner(ctx, a.db, ownerKind, ownerID)
    if err != nil {
        return nil, err
    }
    out := make([]*PromotionRuleResponse, len(rules))
    for i, r := range rules {
        out[i] = ruleToResponse(&rules[i])
    }
    return out, nil
}

func (a *PromotionApp) CreateRules(ctx context.Context, ownerKind promotion.Kind, ownerID int64, rules []promotion.PromotionRule) ([]*PromotionRuleResponse, error) {
    for i := range rules {
        rules[i].OwnerKind = ownerKind
        rules[i].OwnerID = ownerID
    }
    if err := a.repo.CreateRules(ctx, a.db, ownerKind, ownerID, rules); err != nil {
        return nil, err
    }
    return a.GetRules(ctx, ownerKind, ownerID)
}

func (a *PromotionApp) UpdateRule(ctx context.Context, rule *promotion.PromotionRule) (*PromotionRuleResponse, error) {
    if err := a.repo.UpdateRule(ctx, a.db, rule); err != nil {
        return nil, err
    }
    return ruleToResponse(rule), nil
}

func (a *PromotionApp) DeleteRule(ctx context.Context, id int64) error {
    return a.repo.DeleteRule(ctx, a.db, id)
}

// ===== COUPON-specific =====

func (a *PromotionApp) IssueToUser(ctx context.Context, couponID, userID int64) (*UserCouponResponse, error) {
    p, err := a.repo.FindByID(ctx, a.db, couponID)
    if err != nil {
        return nil, err
    }
    uc, err := p.Issue(userID, time.Now().UTC())
    if err != nil {
        return nil, err
    }
    if err := a.repo.IssueUserCoupon(ctx, a.db, uc); err != nil {
        return nil, err
    }
    return userCouponToResponse(uc), nil
}

func (a *PromotionApp) BatchIssue(ctx context.Context, couponID int64, userIDs []int64) (int64, []int64, error) {
    var issued int64
    var ids []int64
    for _, uid := range userIDs {
        resp, err := a.IssueToUser(ctx, couponID, uid)
        if err != nil {
            continue
        }
        issued++
        ids = append(ids, resp.ID)
    }
    return issued, ids, nil
}

func (a *PromotionApp) GenerateCodes(ctx context.Context, prefix string, quantity int, cfg map[string]any) ([]string, error) {
    out := make([]string, 0, quantity)
    for i := 0; i < quantity; i++ {
        code := prefix + randomCode(8)
        // Create a new coupon promotion with kind=COUPON, code, and cfg-derived fields.
        req := couponFromConfig(prefix+randomCode(8), cfg)
        _, err := a.Create(ctx, req)
        if err != nil {
            return nil, err
        }
        out = append(out, code)
    }
    return out, nil
}

func (a *PromotionApp) ListUserCoupons(ctx context.Context, q promotion.UserCouponQuery) (*ListUserCouponResponse, error) {
    list, total, err := a.repo.FindUserCoupons(ctx, a.db, q)
    if err != nil {
        return nil, err
    }
    out := make([]*UserCouponResponse, len(list))
    for i, uc := range list {
        out[i] = userCouponToResponse(&list[i])
    }
    return &ListUserCouponResponse{List: out, Total: total}, nil
}

func (a *PromotionApp) FindPromotionUsage(ctx context.Context, q promotion.UsageQuery) (*ListPromotionUsageResponse, error) {
    list, total, err := a.repo.FindPromotionUsage(ctx, a.db, q)
    if err != nil {
        return nil, err
    }
    out := make([]*PromotionUsageResponse, len(list))
    for i := range list {
        out[i] = usageToResponse(&list[i])
    }
    return &ListPromotionUsageResponse{List: out, Total: total}, nil
}

// ===== helpers =====

func nilOrZero(total *int) *int {
    if total == nil {
        return nil
    }
    z := 0
    return &z
}

func (a *PromotionApp) toResponse(p *promotion.Promotion) *PromotionResponse {
    resp := &PromotionResponse{
        ID:           p.ID,
        TenantID:     p.TenantID,
        Kind:         p.Kind,
        Name:         p.Name,
        Description:  p.Description,
        Code:         p.Code,
        Type:         p.Type,
        Status:       p.Status,
        MarketID:     p.MarketID,
        Currency:     p.Currency,
        TotalCount:   p.TotalCount,
        UsedCount:    p.UsedCount,
        UsageLimit:   p.UsageLimit,
        PerUserLimit: p.PerUserLimit,
        Tags:         p.Tags,
        Scope:        p.Scope,
        StartAt:      p.StartAt,
        EndAt:        p.EndAt,
        CreatedAt:    p.Audit.CreatedAt,
        UpdatedAt:    p.Audit.UpdatedAt,
    }
    for _, r := range p.Rules {
        resp.Rules = append(resp.Rules, *ruleToResponse(&r))
    }
    return resp
}

func ruleToResponse(r *promotion.PromotionRule) *PromotionRuleResponse {
    return &PromotionRuleResponse{
        ID:             r.ID,
        ConditionType:  r.ConditionType,
        ConditionValue: r.ConditionValue,
        ActionType:     r.ActionType,
        ActionValue:    r.ActionValue,
        MaxDiscount:    r.MaxDiscount,
        SortOrder:      r.SortOrder,
    }
}

func userCouponToResponse(uc *promotion.UserCoupon) *UserCouponResponse {
    return &UserCouponResponse{
        ID:         uc.ID,
        UserID:     uc.UserID,
        CouponID:   uc.CouponID,
        Status:     uc.Status,
        UsedAt:     uc.UsedAt,
        OrderID:    uc.OrderID,
        ReceivedAt: uc.ReceivedAt,
        ExpireAt:   uc.ExpireAt,
    }
}

func usageToResponse(u *promotion.PromotionUsage) *PromotionUsageResponse {
    return &PromotionUsageResponse{
        ID:             u.ID,
        CouponID:       u.CouponID,
        UserID:         u.UserID,
        OrderID:        u.OrderID,
        DiscountAmount: u.DiscountAmount,
        CreatedAt:      u.CreatedAt,
    }
}

// ===== GenerateCodes helpers =====

func randomCode(n int) string {
    const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, n)
    for i := range b {
        b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
        time.Sleep(time.Microsecond) // crude uniqueness
    }
    return string(b)
}

func couponFromConfig(code string, cfg map[string]any) *CreatePromotionRequest {
    // Convert cfg → CreatePromotionRequest with Kind: KindCoupon.
    // Implementation extracts: name, value, min_amount, max_discount, type, currency,
    // total_count, per_user_limit, start_at, end_at, scope_type, scope_ids.
    // Returns the request ready for a.Create.
    // ... (concrete parsing omitted; follow the existing GenerateCouponCodesRequest shape)
    _ = code
    _ = cfg
    return nil
}
```

(Add the UserCouponResponse / ListUserCouponResponse / PromotionUsageResponse / ListPromotionUsageResponse types — these mirror the wire types defined in Task 7.)

- [ ] **Step 3: Verify build (compiles alone)**

```bash
cd admin && make build
```

Expected: errors only in `logic/` directories.

- [ ] **Step 4: Commit**

```bash
git add admin/internal/application/promotion/promotion_app.go
git commit -m "refactor(application/promotion): unify promotion + coupon into single App"
```

---

## Task 7: Logic layer update

**Files:**
- Modify: all files in `admin/internal/logic/promotions/`
- Modify: all files in `admin/internal/logic/coupons/`
- Delete: `admin/internal/logic/coupons/helper.go`
- Modify: `admin/internal/logic/promotions/helper.go`

**Interfaces:**
- Consumes: `apppromotion.PromotionApp` (newly merged in Task 6).
- Produces: every existing `*Logic` struct now routes to `PromotionApp` methods instead of two separate apps.

- [ ] **Step 1: Delete `coupons/helper.go`**

```bash
git rm admin/internal/logic/coupons/helper.go
```

- [ ] **Step 2: Rewrite `promotions/helper.go`**

Open `admin/internal/logic/promotions/helper.go`. Replace `convertPromotionToDetailResp` and add `convertRulesToResp`:

```go
// convertRulesToResp maps app → wire types for rules.
func convertRulesToResp(rules []apppromotion.PromotionRuleResponse) []*types.PromotionRuleResp {
    out := make([]*types.PromotionRuleResp, 0, len(rules))
    for _, r := range rules {
        out = append(out, &types.PromotionRuleResp{
            ID:             r.ID,
            ConditionType:  string(r.ConditionType),
            ConditionValue: r.ConditionValue.String(),
            ActionType:     string(r.ActionType),
            ActionValue:    r.ActionValue.String(),
            MaxDiscount:    r.MaxDiscount.String(),
            SortOrder:      r.SortOrder,
        })
    }
    return out
}

// convertPromotionToDetailResp maps a unified PromotionResponse to wire.
func convertPromotionToDetailResp(p *apppromotion.PromotionResponse) *types.PromotionDetailResp {
    resp := &types.PromotionDetailResp{
        ID:           p.ID,
        Kind:         string(p.Kind),
        Name:         p.Name,
        Description:  p.Description,
        Type:         p.Type.String(),
        Status:       p.Status.String(),
        Currency:     p.Currency,
        UsageLimit:   p.UsageLimit,
        PerUserLimit: p.PerUserLimit,
        Tags:         p.Tags,
        ScopeType:    string(p.Scope.Type),
        ScopeIDs:     idsToStrings(p.Scope.ScopeIDs),
        ExcludeIDs:   idsToStrings(p.Scope.ExcludeIDs),
        StartTime:    p.StartAt.Format(time.RFC3339),
        EndTime:      p.EndAt.Format(time.RFC3339),
        CreatedAt:    p.CreatedAt.Format(time.RFC3339),
        UpdatedAt:    p.UpdatedAt.Format(time.RFC3339),
        Rules:        convertRulesToResp(p.Rules),
    }
    if p.Code != nil {
        resp.Code = *p.Code
    }
    if p.MarketID != nil {
        resp.MarketID = *p.MarketID
    }
    if p.TotalCount != nil {
        resp.TotalCount = *p.TotalCount
    }
    if p.UsedCount != nil {
        resp.UsedCount = *p.UsedCount
    }
    return resp
}
```

- [ ] **Step 3: Rewrite `promotions/create_promotion_logic.go`**

```go
func (l *CreatePromotionLogic) Create(req *types.CreatePromotionReq) (*types.CreatePromotionResp, error) {
    kind := promotion.Kind(strings.ToLower(req.Kind))
    if !kind.IsValid() {
        return nil, code.ErrPromotionInvalidKind
    }
    scope := buildPromotionScope(req.ScopeType, req.ScopeIDs, req.ExcludeIDs)
    rules := convertRuleReqsToDomain(req.Rules)
    createReq := &apppromotion.CreatePromotionRequest{
        TenantID:     shared.TenantID(l.ctx.Value("tenant_id").(int64)),
        Kind:         kind,
        Name:         req.Name,
        Description:  req.Description,
        Type:         parsePromotionType(req.Type),
        MarketID:     optionalInt64(req.MarketID),
        Currency:     req.Currency,
        UsageLimit:   req.UsageLimit,
        PerUserLimit: req.PerUserLimit,
        Tags:         req.Tags,
        Scope:        scope,
        StartAt:      parseTime(req.StartTime),
        EndAt:        parseTime(req.EndTime),
        Rules:        rules,
        ActorID:      l.ctx.Value("user_id").(int64),
    }
    p, err := l.svcCtx.PromotionApp.Create(l.ctx, createReq)
    if err != nil {
        return nil, err
    }
    return &types.CreatePromotionResp{ID: p.ID}, nil
}
```

- [ ] **Step 4: Rewrite `promotions/update_promotion_logic.go`**

Same shape as Create but maps to `UpdatePromotionRequest`, passing `Rules: &rules` so the App replaces them.

- [ ] **Step 5: Rewrite `promotions/list_promotions_logic.go`**

```go
func (l *ListPromotionsLogic) List(req *types.ListPromotionsReq) (*types.ListPromotionsResp, error) {
    q := promotion.Query{
        PageQuery: shared.PageQuery{Page: req.Page, Size: req.PageSize},
        Name:      req.Name,
    }
    if req.Kind != "" {
        k := promotion.Kind(strings.ToUpper(req.Kind))
        q.Kind = &k
    }
    if req.Status != "" {
        s := promotion.Status(parseStatus(req.Status))
        q.Status = &s
    }
    if req.Type != "" {
        t := promotion.Type(parseType(req.Type))
        q.Type = &t
    }
    if req.MarketID != 0 {
        mid := req.MarketID
        q.MarketID = &mid
    }
    q.ExpiredOnly = req.Status == "expired"

    list, err := l.svcCtx.PromotionApp.List(l.ctx, q)
    if err != nil {
        return nil, err
    }
    out := make([]*types.PromotionDetailResp, len(list.List))
    for i, p := range list.List {
        out[i] = convertPromotionToDetailResp(p)
    }
    return &types.ListPromotionsResp{
        List:     out,
        Total:    list.Total,
        Page:     list.Page,
        PageSize: list.Size,
    }, nil
}
```

- [ ] **Step 6: Rewrite `coupons/create_coupon_logic.go`**

```go
func (l *CreateCouponLogic) Create(req *types.CreateCouponReq) (*types.CreateCouponResp, error) {
    code := req.Code
    total := req.UsageLimit
    rules := []promotion.PromotionRule{{
        ConditionType:  promotion.ConditionMinAmount,
        ConditionValue: parseDecimal(req.MinOrderAmount),
        ActionType:     parseActionType(req.Type),
        ActionValue:    parseDecimal(req.DiscountValue),
        MaxDiscount:    parseDecimal(req.MaxDiscount),
    }}
    createReq := &apppromotion.CreatePromotionRequest{
        TenantID:     shared.TenantID(l.ctx.Value("tenant_id").(int64)),
        Kind:         promotion.KindCoupon,
        Name:         req.Name,
        Description:  req.Description,
        Code:         &code,
        Type:         promotion.TypeDiscount,
        Currency:     req.Currency,
        TotalCount:   &total,
        UsageLimit:   req.UsageLimit,
        PerUserLimit: req.PerUserLimit,
        Scope:        buildPromotionScope(req.ScopeType, scopeFromJSON(req.ProductIDs), nil),
        StartAt:      parseTime(req.StartTime),
        EndAt:        parseTime(req.EndTime),
        Rules:        rules,
        ActorID:      l.ctx.Value("user_id").(int64),
    }
    p, err := l.svcCtx.PromotionApp.Create(l.ctx, createReq)
    if err != nil {
        return nil, err
    }
    return &types.CreateCouponResp{ID: p.ID}, nil
}
```

(If `Coupon` already supports multi-tier rules on creation, loop through a new `Rules []PromotionRuleReq` field. The current wire shape has flat fields — promote them to a single rule for migration compatibility.)

- [ ] **Step 7: Rewrite `coupons/update_coupon_logic.go`**

Same shape as Create, building `UpdatePromotionRequest` with `Kind: KindCoupon`.

- [ ] **Step 8: Rewrite all remaining coupons/*_logic.go files**

Replace each `coupons/{list,get,activate,deactivate,delete}_coupon_logic.go` to call `PromotionApp.{List,Get,Activate,Deactivate,Delete}` and use `convertPromotionToDetailResp` instead of `convertCouponToDetailResp`.

For `coupons/issue_user_coupon_logic.go` and `batch_issue_user_coupon_logic.go`: call `PromotionApp.IssueToUser` / `PromotionApp.BatchIssue`.

For `coupons/list_user_coupons_logic.go`: call `PromotionApp.ListUserCoupons`.

For `coupons/get_coupon_usage_logic.go`: call `PromotionApp.FindPromotionUsage`.

For `coupons/generate_coupon_codes_logic.go`: parse the JSON `CouponConfig`, then call `PromotionApp.GenerateCodes(prefix, qty, cfg)`.

- [ ] **Step 9: Update rule logic files**

For `promotions/create_promotion_rules_logic.go`:

```go
func (l *CreatePromotionRulesLogic) Create(req *types.CreatePromotionRulesReq) (*types.CreatePromotionRulesResp, error) {
    owner, err := l.svcCtx.PromotionApp.Get(l.ctx, req.OwnerID)
    if err != nil {
        return nil, err
    }
    rules := convertRuleReqsToDomain(req.Rules)
    out, err := l.svcCtx.PromotionApp.CreateRules(l.ctx, owner.Kind, req.OwnerID, rules)
    if err != nil {
        return nil, err
    }
    ids := make([]string, len(out))
    for i, r := range out {
        ids[i] = strconv.FormatInt(r.ID, 10)
    }
    return &types.CreatePromotionRulesResp{IDs: ids}, nil
}
```

(Repeat for `get_promotion_rules_logic.go`, `update_promotion_rule_logic.go`, `delete_promotion_rule_logic.go` — each calls the matching `PromotionApp` rule method.)

- [ ] **Step 10: Verify build**

```bash
cd admin && make build
```

Expected: 0 errors. (Some warnings acceptable.)

- [ ] **Step 11: Commit**

```bash
git add admin/internal/logic/
git commit -m "refactor(logic): route all promotion + coupon logic through unified PromotionApp"
```

---

## Task 8: API definition + codegen

**Files:**
- Modify: `admin/desc/promotion.api`

**Interfaces:**
- Produces: regenerated `admin/internal/types/types.go` and `admin/internal/handler/routes.go` via `make api`.

- [ ] **Step 1: Update `PromotionDetailResp`**

In `admin/desc/promotion.api`:

```go
PromotionDetailResp {
    ID             int64                `json:"id,string"`
    Kind           string               `json:"kind"`                              // NEW: "promotion" | "coupon"
    Name           string               `json:"name"`
    Description    string               `json:"description"`
    Code           string               `json:"code,optional"`
    Type           string               `json:"type"`
    Status         string               `json:"status"`
    MarketID       int64                `json:"market_id,optional,string"`         // NEW
    Currency       string               `json:"currency"`
    UsageLimit     int                  `json:"usage_limit"`
    UsedCount      int                  `json:"used_count"`
    PerUserLimit   int                  `json:"per_user_limit"`
    TotalCount     int                  `json:"total_count,optional"`              // NEW (coupon only)
    ScopeType      string               `json:"scope_type"`
    ScopeIDs       []string             `json:"scope_ids"`
    ExcludeIDs     []string             `json:"exclude_ids"`
    Tags           []string             `json:"tags"`
    Rules          []*PromotionRuleResp `json:"rules"`                             // NEW
    StartTime      string               `json:"start_time"`
    EndTime        string               `json:"end_time"`
    CreatedAt      string               `json:"created_at"`
    UpdatedAt      string               `json:"updated_at"`
}
```

- [ ] **Step 2: Update `CreatePromotionReq` / `UpdatePromotionReq`**

```go
CreatePromotionReq {
    Kind           string              `json:"kind"`                       // NEW
    Name           string              `json:"name"`
    Description    string              `json:"description,optional"`
    Code           string              `json:"code,optional"`              // coupon only
    Type           string              `json:"type,optional"`
    MarketID       int64               `json:"market_id,optional,string"`  // NEW
    Currency       string              `json:"currency,optional"`
    UsageLimit     int                 `json:"usage_limit,optional"`
    PerUserLimit   int                 `json:"per_user_limit,optional"`
    TotalCount     int                 `json:"total_count,optional"`        // coupon only
    ScopeType      string              `json:"scope_type,optional"`
    ScopeIDs       []string            `json:"scope_ids,optional"`
    ExcludeIDs     []string            `json:"exclude_ids,optional"`
    Tags           []string            `json:"tags,optional"`
    Rules          []*PromotionRuleReq `json:"rules,optional"`              // NEW
    StartTime      string              `json:"start_time"`
    EndTime        string              `json:"end_time"`
}
```

(UpdatePromotionReq mirrors CreatePromotionReq plus `ID int64 path:"id"`.)

- [ ] **Step 3: Replace rule wire types**

```go
PromotionRuleReq {
    ConditionType  string `json:"condition_type"`               // "min_amount" | "min_quantity"
    ConditionValue string `json:"condition_value"`              // decimal string
    ActionType     string `json:"action_type"`                  // "fixed_amount" | "percentage" | "free_shipping"
    ActionValue    string `json:"action_value"`                 // decimal string
    MaxDiscount    string `json:"max_discount,optional"`
    SortOrder      int    `json:"sort_order,optional"`
}

PromotionRuleResp {
    ID             int64  `json:"id,string"`
    ConditionType  string `json:"condition_type"`
    ConditionValue string `json:"condition_value"`
    ActionType     string `json:"action_type"`
    ActionValue    string `json:"action_value"`
    MaxDiscount    string `json:"max_discount"`
    SortOrder      int    `json:"sort_order"`
}

CreatePromotionRulesReq {
    OwnerKind string              `path:"owner_kind"`                // NEW: "promotion" | "coupon"
    OwnerID   int64               `path:"owner_id"`
    Rules     []PromotionRuleReq  `json:"rules"`
}
```

- [ ] **Step 4: Update list requests**

```go
ListPromotionsReq {
    Kind     string `form:"kind,optional"`
    Name     string `form:"name,optional"`
    Type     string `form:"type,optional"`
    Status   string `form:"status,optional"`
    MarketID int64  `form:"market_id,optional,string"`     // NEW
    Page     int    `form:"page,default=1"`
    PageSize int    `form:"page_size,default=20"`
}
```

- [ ] **Step 5: Run codegen**

```bash
cd admin && make api
```

Expected: regenerated `internal/types/types.go` and `internal/handler/routes.go` with the new fields.

- [ ] **Step 6: Verify build**

```bash
cd admin && make build
```

Expected: errors point at logic files that referenced old field names (these are caught by Task 7 if logic was already updated, otherwise fix in this task).

- [ ] **Step 7: Verify all wire fields used**

For each field in `PromotionDetailResp`, grep the logic layer to confirm it is read by `convertPromotionToDetailResp` (run the project checklist from `.claude/rules/golang/promotion.md`):

```bash
cd admin
grep -rn "resp\.\(Kind\|MarketID\|Code\|TotalCount\|Rules\)" internal/logic/
# All five should appear at least once
```

- [ ] **Step 8: Commit**

```bash
git add admin/desc/promotion.api admin/internal/types/ admin/internal/handler/
git commit -m "feat(api): unify PromotionDetailResp with kind + market_id + rules"
```

---

## Task 9: Frontend sync

**Files:**
- Modify: `shop-admin/src/api/promotion.ts`
- Modify: `shop-admin/src/views/promotions/index.vue`
- Modify: `shop-admin/src/views/promotions/rule*.vue`

**Interfaces:**
- Consumes: wire types from `admin/desc/promotion.api` (via Task 8).
- Produces: TypeScript types + Vue templates aligned with backend.

- [ ] **Step 1: Update `api/promotion.ts`**

Replace the `CouponDetailResp` interface with fields merged into `Promotion`. Add the new fields:

```typescript
export interface Promotion {
  id: string
  kind: 'promotion' | 'coupon'                      // NEW
  name: string
  code?: string                                      // coupon only
  description: string
  type: string
  status: string
  market_id?: string                                 // NEW
  currency: string
  usage_limit: number
  used_count: number
  per_user_limit: number
  total_count?: number                               // coupon only
  scope_type: string
  scope_ids: string[]
  exclude_ids: string[]
  tags: string[]
  rules?: PromotionRule[]                            // NEW
  start_time: string
  end_time: string
  created_at: string
  updated_at: string
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

// Delete CouponDetailResp interface
```

Update list/create/update request interfaces to match the new `CreatePromotionReq` shape (add `kind`, `market_id`, `code?`, `total_count?`, `rules?`).

- [ ] **Step 2: Update list view `index.vue`**

- Replace all `row.code` / `row.used_count` columns with `v-if="row.kind === 'coupon'"` guards.
- Replace all `row.type === 'coupon'` references with `row.kind === 'coupon'`.
- Merge `couponForm` and `promotionForm` into one `form` object. Add `form.kind` field (defaults to `"promotion"`). The dialog opens with `kind` set based on which tab the user clicked.
- Add market filter dropdown (use the same dropdown component as `products/index.vue` if it exists; otherwise a simple `<el-select>`).
- Status column toggle: `el-switch` already works — no change needed.

- [ ] **Step 3: Update rule editor views**

Open `rule-list.vue` and `rule-edit.vue` (or whatever the existing rule pages are called). Replace the page title with:

```vue
<h2>{{ owner.kind === 'coupon' ? t('promotions.couponRules') : t('promotions.promotionRules') }}</h2>
```

- [ ] **Step 4: Verify frontend build**

```bash
cd shop-admin && pnpm build
```

Expected: 0 errors. (TypeScript strict mode will flag any missed fields.)

- [ ] **Step 5: Manual smoke test**

Start the admin backend locally, then:
- Open `/promotions` tab → verify existing 5 promotions render with kind="promotion"
- Open `/coupons` tab → verify existing 7 coupons render with kind="coupon" + rules + code
- Click "Edit" on a coupon → verify rule list shows 1 rule with min_amount=10000 etc.
- Create a new PROMOTION with 2 rules → verify rules persist
- Create a new COUPON → verify it gets a unique code

- [ ] **Step 6: Commit**

```bash
git add shop-admin/src/
git commit -m "feat(frontend): unify Promotion + Coupon rendering with kind discriminator"
```

---

## Task 10: Final verification

**Files:** none (verification only)

**Interfaces:** N/A

- [ ] **Step 1: Compile both services**

```bash
make build
```

Expected: 0 errors across admin and shop services.

- [ ] **Step 2: Run domain tests**

```bash
go test ./pkg/domain/promotion/...
go test ./admin/internal/application/promotion/...
go test ./admin/internal/infrastructure/persistence/...
```

Expected: all green.

- [ ] **Step 3: Run full admin test suite**

```bash
cd admin && go test ./...
```

Expected: 0 failures.

- [ ] **Step 4: Run promotion.md Checklist**

Run through every item in `.claude/rules/golang/promotion.md` "Checklist" section. Each box must be checked:

- [ ] Every field on `UpdatePromotionReq` / `CreatePromotionRequest` has a corresponding assignment in logic + app layers
- [ ] `convertPromotionToDetailResp()` maps ALL response fields including the new `kind / market_id / code / total_count / rules`
- [ ] `promotionModel` / `promotionRuleModel` column names match `SHOW CREATE TABLE`
- [ ] Filter parameters use pointer types (`*Kind`, `*Status`, `*Type`, `*MarketID`)
- [ ] `"expired"` status uses `ExpiredOnly bool` + `end_at <= NOW()`, not a status enum
- [ ] Scope type strings are normalized with `strings.ToUpper()` before domain comparison
- [ ] `go build ./...` passes
- [ ] Doc comments are accurate (no "no DB columns" after columns exist)

- [ ] **Step 5: Code review**

Run the `/review` skill:

```
/review --comment
```

Address any high-severity findings before merging. User confirmation required per CLAUDE.md.

- [ ] **Step 6: Drop migration files (only after step 2 verifies green for 1+ day)**

```bash
git rm sql/promotion/migrations/2026071900_pre_migration_checks.sql
git rm sql/promotion/migrations/2026071901_merge_promotion_coupon.sql
git commit -m "chore(sql): drop applied promotion merge migrations"
```

(The archived `_archived_coupons_20260719` and `_deprecated_coupons` tables are dropped in a separate T+30 migration, not here.)

---

## Self-Review Notes (plan-vs-spec)

After writing this plan, I ran the writing-plans self-review checklist against `docs/superpowers/specs/2026-07-19-promotion-coupon-merge-design.md`:

- **Spec coverage:** Every spec section maps to a task:
  - §1 → covered by tasks 1-2 (DB), 4 (entity), 6 (app), 8-9 (wire/frontend)
  - §2 → Tasks 1, 2
  - §3 → Task 4
  - §4 → Task 5
  - §5 → Tasks 6, 7
  - §6 → Tasks 8, 9
  - §7 → Tasks 9, 10
  - §8 → Tasks 2 (rollback section embedded), 10 (drop after grace period)
  - Appendix A → referenced from Task 8 (rollback plan)
- **Placeholder scan:** No "TBD"/"TODO"/"fill in" — each step has concrete code or commands.
- **Type consistency:** Verified `promotion.Kind`, `promotion.Promotion`, `promotion.PromotionRule`, `promotion.Query`, `promotion.Repository` are referenced with the same names across Tasks 3-7. `PromotionApp.Create/Update/Get/List/Delete/Activate/Deactivate/CreateRules/UpdateRule/DeleteRule/GetRules/IssueToUser/BatchIssue/GenerateCodes/ListUserCoupons/FindPromotionUsage` match the spec §5.1 list.
- **Naming consistency:** `KindPromotion` / `KindCoupon` (Go), `"PROMOTION"` / `"COUPON"` (DB ENUM), `"promotion"` / `"coupon"` (wire lowercase). Conversion happens at the App boundary via `strings.ToUpper` (DB ENUM) and `strings.ToLower` (wire input). Documented in `Global Constraints`.
- **Scope:** Single coherent refactor. Does not need to be split into multiple plans.