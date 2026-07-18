# Promotion usage_limit / per_user_limit / tags Persistence — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Persist `usage_limit`, `per_user_limit`, and `tags` on promotions so admins can edit them in the form and have the values survive a save→re-fetch cycle.

**Architecture:** Add three columns to the `promotions` table, mirror them through the existing four-layer stack (entity → repository → application → logic). The wire types already carry these fields from earlier work — only the persistence + the application-side wiring are new.

**Tech Stack:** Go 1.x, go-zero (`httpx.Parse`, `logx`), GORM v2, MySQL 8 (JSON column type), go-playground/validator (already imported via `pkg/httpy/parse.go`).

**Source spec:** `docs/domains/promotion/2026-07-18-promotion-usage-limit-tags-design.md` (commit `16488ce`).

## Global Constraints

- All time fields use `time.Time`; storage column `TIMESTAMP`. GORM handles conversion.
- Money/price rules: not directly relevant (this ticket is int + JSON).
- Errors: always use `pkg/code` defined errors. Two new ones: `ErrPromotionUsageLimitInvalid` (code 80016) and `ErrPromotionPerUserLimitInvalid` (code 80017).
- HTTP status from `code.Err`: 400 for invalid input.
- DB names: lowercase + underscores.
- Default `not null` with explicit `DEFAULT`.
- After any code changes: `cd admin && go build ./...` must succeed before commit.
- Two-round review per project rules; `/review` before commit on each task.

## File Structure (locked-in by this plan)

| Path | Responsibility |
|---|---|
| `sql/promotion/schema.sql` | Single source of truth for promotion tables. Add 3 columns. |
| `pkg/domain/promotion/entity.go` | `Promotion` aggregate adds 3 typed fields. |
| `pkg/code/code.go` | Two new error codes in the Promotion Module (80xxx) range. |
| `admin/internal/infrastructure/persistence/promotion_repository.go` | `promotionModel` mirrors schema; `toEntity`/`fromPromotionEntity` map JSON column ↔ `[]string`; `Update` map includes new columns. |
| `admin/internal/application/promotion/promotion_app.go` | `Create/UpdatePromotionRequest` add 3 fields; `PromotionResponse` adds 3; `CreatePromotion`/`UpdatePromotion`/`toPromotionResponse` wire them. |
| `admin/internal/logic/promotions/update_promotion_logic.go` | Pass-through 3 fields into `UpdatePromotionRequest`. |
| `admin/internal/logic/promotions/create_promotion_logic.go` | Pass-through 3 fields into `CreatePromotionRequest`. |

No frontend changes (per spec §3.4, §10). No wire-type changes (per spec §6).

---

### Task 1: Schema column additions

**Files:**
- Modify: `sql/promotion/schema.sql` — add 3 columns via `ALTER TABLE` block at the bottom of the file.

- [ ] **Step 1: Open `sql/promotion/schema.sql` and append the ALTER block**

Append after the existing `CREATE TABLE promotions (...)` block:

```sql
-- 2026-07-18: add usage_limit / per_user_limit / tags (commit 16488ce)
-- These columns back the wire-level fields that were previously
-- dropped between the API edge and the DB layer. Defaults chosen so
-- existing rows behave as "unlimited" without backfill.
ALTER TABLE promotions
    ADD COLUMN usage_limit    INT     NOT NULL DEFAULT 0
        COMMENT '0 = unlimited',
    ADD COLUMN per_user_limit INT     NOT NULL DEFAULT 1
        COMMENT 'per-user cap; 0 = unlimited',
    ADD COLUMN tags           JSON    NULL
        COMMENT 'free-form labels, stored as JSON array of strings';
```

- [ ] **Step 2: Apply the migration to the local dev DB**

Run via mycli:
```bash
mycli -h 192.168.0.100 -P 3306 -u root -p123456 shopjoy -e "
ALTER TABLE promotions
    ADD COLUMN usage_limit    INT     NOT NULL DEFAULT 0,
    ADD COLUMN per_user_limit INT     NOT NULL DEFAULT 1,
    ADD COLUMN tags           JSON    NULL;"
```
Expected: "0 row(s) affected" (DDL).

- [ ] **Step 3: Verify the columns exist**

```bash
mycli -h 192.168.0.100 -P 3306 -u root -p123456 shopjoy -e "DESCRIBE promotions;"
```
Expected: rows for `usage_limit`, `per_user_limit`, `tags` appear.

- [ ] **Step 4: Commit**

```bash
git add sql/promotion/schema.sql
git commit -m "sql(promotion): add usage_limit, per_user_limit, tags columns"
```

---

### Task 2: Domain entity fields

**Files:**
- Modify: `pkg/domain/promotion/entity.go` — `Promotion` struct.

- [ ] **Step 1: Locate the `Promotion` struct**

Open `pkg/domain/promotion/entity.go`. Find the `Promotion` struct (around line 234 in the current file). It currently ends with:
```go
type Promotion struct {
    // ... existing fields ...
    Audit       shared.AuditInfo `json:"audit"`
    DeletedAt   *time.Time       `json:"deleted_at,omitempty"`
}
```

- [ ] **Step 2: Add three fields after `Scope`**

Insert before `Audit`:
```go
    UsageLimit   int      `json:"usage_limit"`
    PerUserLimit int      `json:"per_user_limit"`
    Tags         []string `json:"tags,omitempty" gorm:"type:json"`
```

Resulting struct:
```go
type Promotion struct {
    ID          int64            `json:"id"`
    TenantID    shared.TenantID  `json:"tenant_id"`
    Name        string           `json:"name"`
    Description string           `json:"description"`
    Type        Type             `json:"type"`
    Status      Status           `json:"status"`
    Priority    int              `json:"priority"`
    StartAt     time.Time        `json:"start_at"`
    EndAt       time.Time        `json:"end_at"`
    Scope       PromotionScope   `json:"scope"`
    UsageLimit   int             `json:"usage_limit"`
    PerUserLimit int             `json:"per_user_limit"`
    Tags         []string        `json:"tags,omitempty" gorm:"type:json"`
    Currency    string           `json:"currency"`
    Audit       shared.AuditInfo `json:"audit"`
    DeletedAt   *time.Time       `json:"deleted_at,omitempty"`
}
```

- [ ] **Step 3: Verify the package still compiles**

Run: `cd /Users/dengyichuan/workspace/go/src/shopjoy && go build ./pkg/domain/promotion/...`
Expected: no output, exit 0.

- [ ] **Step 4: Commit**

```bash
git add pkg/domain/promotion/entity.go
git commit -m "domain(promotion): add UsageLimit, PerUserLimit, Tags fields to Promotion"
```

---

### Task 3: Persistence layer

**Files:**
- Modify: `admin/internal/infrastructure/persistence/promotion_repository.go` — `promotionModel`, `toEntity`, `fromPromotionEntity`, `Update` map.

**Interfaces:**
- Consumes: `pkgpromotion.Promotion` struct with 3 new fields (from Task 2).
- Produces: DB rows with new columns populated; `*promotion.Promotion` returned from `FindByID`/`FindList` carries the 3 fields.

- [ ] **Step 1: Read the current `promotionModel` struct (around lines 24-43)**

Open `admin/internal/infrastructure/persistence/promotion_repository.go`. Locate the `promotionModel` definition. The `Update` method (around line 169-191) uses an explicit column list.

- [ ] **Step 2: Add 3 columns to `promotionModel`**

Insert after `UpdatedBy`:
```go
    UsageLimit    int    `gorm:"column:usage_limit;not null;default:0"`
    PerUserLimit  int    `gorm:"column:per_user_limit;not null;default:1"`
    Tags          string `gorm:"column:tags;type:json"` // JSON-encoded []string, empty/nullsafe
```

(Notes:
- `int`/`int`/`string` for the model columns. We use Go's `int` for the two limits and `string` for the JSON-encoded array to match the project's existing `ScopeIDs string + json.Marshal` pattern from the coupon repo — keeps the dependency surface flat and avoids pulling `gorm.io/datatypes` into this codebase.)

- [ ] **Step 3: Extend `toEntity()` to populate the three new fields**

Inside `toEntity()`, after the `Scope:` literal that builds `PromotionScope`, add:
```go
        UsageLimit:   int(m.UsageLimit),
        PerUserLimit: int(m.PerUserLimit),
        Tags:         decodeJSONStringSlice(m.Tags),
```

- [ ] **Step 4: Extend `fromPromotionEntity()` to populate the three new fields**

Inside `fromPromotionEntity()`, after `ScopeType:`:
```go
        UsageLimit:   int(p.UsageLimit),
        PerUserLimit: int(p.PerUserLimit),
        Tags:         encodeJSONStringSlice(p.Tags),
```

- [ ] **Step 5: Add two local helpers at the top of the file (after imports)**

Above `// promotionModel represents ...`:

```go
// decodeJSONStringSlice parses a MySQL JSON column value into a
// string slice, returning nil on missing / malformed / empty
// payloads. The wire/UI layer treats nil and [] as equivalent.
func decodeJSONStringSlice(s string) []string {
    if strings.TrimSpace(s) == "" {
        return nil
    }
    var out []string
    if err := json.Unmarshal([]byte(s), &out); err != nil {
        return nil
    }
    return out
}

// encodeJSONStringSlice renders a string slice as a MySQL JSON
// column value. nil slices become NULL ("null") so the row
// distinguishes "no tags" from "empty array" by absence.
func encodeJSONStringSlice(in []string) string {
    if len(in) == 0 {
        return "null"
    }
    b, err := json.Marshal(in)
    if err != nil {
        return "null"
    }
    return string(b)
}
```

- [ ] **Step 6: Add `strings` to the import block**

Edit the import block at the top of the file:
```go
import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "strings"
    "time"
    // ... rest unchanged
)
```

- [ ] **Step 7: Extend `Update` map to include the 3 columns**

Find the `Updates(map[string]any{...})` block in the `Update` method. After `ScopeIDs:` insert:
```go
            "usage_limit":    model.UsageLimit,
            "per_user_limit": model.PerUserLimit,
            "tags":           model.Tags,
```

- [ ] **Step 8: Verify the package compiles**

Run: `cd /Users/dengyichuan/workspace/go/src/shopjoy/admin && go build ./...`
Expected: no output, exit 0.

- [ ] **Step 9: Commit**

```bash
git add admin/internal/infrastructure/persistence/promotion_repository.go
git commit -m "persistence(promotion): store usage_limit / per_user_limit / tags"
```

---

### Task 4: Application layer request/response types + persistence

**Files:**
- Modify: `admin/internal/application/promotion/promotion_app.go` — `CreatePromotionRequest`, `UpdatePromotionRequest`, `PromotionResponse`, `CreatePromotion`, `UpdatePromotion`, `toPromotionResponse`.

**Interfaces:**
- Consumes: `pkgpromotion.Promotion` (with the 3 fields from Task 2).
- Produces: callers receive `UpdatePromotionRequest`/`CreatePromotionRequest`/`PromotionResponse` carrying the 3 fields.

- [ ] **Step 1: Locate the three structs**

In `admin/internal/application/promotion/promotion_app.go`, find:
- `CreatePromotionRequest` (line ~17)
- `UpdatePromotionRequest` (line ~50)
- `PromotionResponse` (line ~61)

- [ ] **Step 2: Extend `CreatePromotionRequest`**

Insert before the closing brace:
```go
    UsageLimit   int
    PerUserLimit int
    Tags         []string
```

- [ ] **Step 3: Extend `UpdatePromotionRequest`**

Insert before the closing brace:
```go
    UsageLimit   int
    PerUserLimit int
    Tags         []string
```

- [ ] **Step 4: Extend `PromotionResponse`**

Insert before the closing brace:
```go
    UsageLimit   int      `json:"usage_limit"`
    PerUserLimit int      `json:"per_user_limit"`
    Tags         []string `json:"tags"`
```

- [ ] **Step 5: Wire `CreatePromotion` to write the three fields onto the entity**

Inside `promotionApp.CreatePromotion`'s transaction, after the line `Scope: req.Scope,`:
```go
            UsageLimit:   req.UsageLimit,
            PerUserLimit: req.PerUserLimit,
            Tags:         req.Tags,
```

- [ ] **Step 6: Wire `UpdatePromotion` to write the three fields onto the entity**

Inside `promotionApp.UpdatePromotion`, after the line `p.Scope = req.Scope`:
```go
    p.UsageLimit = req.UsageLimit
    p.PerUserLimit = req.PerUserLimit
    p.Tags = req.Tags
```

- [ ] **Step 7: Add input validation for the two limit fields**

Inside `UpdatePromotion`, before the `if !req.Type.IsValid()` line (right after the existing `if p.Status == pkgpromotion.StatusActive` block), insert:
```go
    if req.UsageLimit < 0 {
        return nil, code.ErrPromotionUsageLimitInvalid
    }
    if req.PerUserLimit < 0 {
        return nil, code.ErrPromotionPerUserLimitInvalid
    }
```

Inside `CreatePromotion`, after the line `if !req.Type.IsValid() {`, insert:
```go
    if req.UsageLimit < 0 {
        return nil, code.ErrPromotionUsageLimitInvalid
    }
    if req.PerUserLimit < 0 {
        return nil, code.ErrPromotionPerUserLimitInvalid
    }
```

- [ ] **Step 8: Tag sanitization helper at the top of the file**

Insert near the top (after imports):
```go
// sanitizeTags enforces a max length of 64 chars per entry and
// drops empty entries. Returns nil for empty result so persistence
// can write SQL NULL consistently.
func sanitizeTags(in []string) []string {
    if len(in) == 0 {
        return nil
    }
    out := make([]string, 0, len(in))
    for _, t := range in {
        t = strings.TrimSpace(t)
        if t == "" {
            continue
        }
        if len(t) > 64 {
            t = t[:64]
        }
        out = append(out, t)
    }
    if len(out) == 0 {
        return nil
    }
    return out
}
```

- [ ] **Step 9: Apply tag sanitization at the request boundary**

In both `CreatePromotion` and `UpdatePromotion`, before constructing the entity, call `req.Tags = sanitizeTags(req.Tags)`. Add as the very first line inside each method body.

- [ ] **Step 10: Wire `toPromotionResponse` to read the three fields off the entity**

In `toPromotionResponse` (around line 562), extend the returned struct:
```go
    UsageLimit:   p.UsageLimit,
    PerUserLimit: p.PerUserLimit,
    Tags:         p.Tags,
```

- [ ] **Step 11: Add `strings` to imports**

Ensure `"strings"` is in the import block.

- [ ] **Step 12: Verify the package compiles**

Run: `cd /Users/dengyichuan/workspace/go/src/shopjoy/admin && go build ./...`
Expected: no output, exit 0.

- [ ] **Step 13: Commit**

```bash
git add admin/internal/application/promotion/promotion_app.go
git commit -m "app(promotion): persist usage_limit / per_user_limit / tags on create+update"
```

---

### Task 5: Error codes

**Files:**
- Modify: `pkg/code/code.go` — add `ErrPromotionUsageLimitInvalid`, `ErrPromotionPerUserLimitInvalid`.

- [ ] **Step 1: Locate the Promotion block**

In `pkg/code/code.go`, find the Promotion Module block comment (`// ErrPromotionNotFound ... Promotion Module (80xxx) ...`).

- [ ] **Step 2: Insert two new codes after `ErrPromotionInvalidTimeRange`**

```go
    ErrPromotionUsageLimitInvalid   = &Err{HTTPCode: http.StatusBadRequest, Code: 80016, Msg: "promotion usage_limit must be >= 0"}
    ErrPromotionPerUserLimitInvalid = &Err{HTTPCode: http.StatusBadRequest, Code: 80017, Msg: "promotion per_user_limit must be >= 0"}
```

- [ ] **Step 3: Verify the package compiles**

Run: `cd /Users/dengyichuan/workspace/go/src/shopjoy && go build ./...`
Expected: no output, exit 0.

- [ ] **Step 4: Commit**

```bash
git add pkg/code/code.go
git commit -m "code(promotion): ErrPromotionUsageLimitInvalid, ErrPromotionPerUserLimitInvalid"
```

---

### Task 6: Logic layer pass-through

**Files:**
- Modify: `admin/internal/logic/promotions/update_promotion_logic.go` — pass 3 fields into `UpdatePromotionRequest`.
- Modify: `admin/internal/logic/promotions/create_promotion_logic.go` — pass 3 fields into `CreatePromotionRequest`.

- [ ] **Step 1: Open `update_promotion_logic.go`**

Locate the `updateReq := apppromotion.UpdatePromotionRequest{...}` literal (after Task 1 / Phase 1 already extended it).

- [ ] **Step 2: Extend the Update literal**

After the closing brace of `Scope: scope,`:
```go
        UsageLimit:   req.UsageLimit,
        PerUserLimit: req.PerUserLimit,
        Tags:         req.Tags,
```

- [ ] **Step 3: Open `create_promotion_logic.go`**

Locate the `createReq := apppromotion.CreatePromotionRequest{...}` literal.

- [ ] **Step 4: Extend the Create literal**

After the closing brace of `Scope: buildPromotionScope(req.ScopeType, req.ProductIDs, req.CategoryIDs),`:
```go
        UsageLimit:   req.UsageLimit,
        PerUserLimit: req.PerUserLimit,
        Tags:         req.Tags,
```

- [ ] **Step 5: Verify the package compiles**

Run: `cd /Users/dengyichuan/workspace/go/src/shopjoy/admin && go build ./...`
Expected: no output, exit 0.

- [ ] **Step 6: Commit**

```bash
git add admin/internal/logic/promotions/update_promotion_logic.go admin/internal/logic/promotions/create_promotion_logic.go
git commit -m "logic(promotion): pass usage_limit / per_user_limit / tags to app layer"
```

---

### Task 7: End-to-end verification

**Files:** none (just runs and tests).

- [ ] **Step 1: Build the admin binary**

Run: `cd /Users/dengyichuan/workspace/go/src/shopjoy && make build`
Expected: `Go code formatted` and exit 0.

- [ ] **Step 2: Start a fresh test instance on port 8893**

Run:
```bash
cat > /tmp/admin-api-verify.yaml <<'EOF'
Name: admin-api-verify
Host: 0.0.0.0
Port: 8893
Verbose: false
Log:
  Stat: false
  Level: error
JWT:
  Secret: "ncjancjabjkacam"
  AccessExpiry: 36000000
  RefreshExpiry: 72000000
StripeWebhookSecret: "aaaa"
MySQL:
  Host: 192.168.0.100
  Port: 3306
  UserName: root
  Password: "123456"
  Database: shopjoy
  MaxIdleConn: 10
  MaxOpenConn: 100
  ConnMaxLifeTime: 3600
EOF
/Users/dengyichuan/workspace/go/src/shopjoy/admin/bin/admin -f /tmp/admin-api-verify.yaml > /tmp/admin-verify.log 2>&1 &
sleep 3
```
Expected: process listening on port 8893 (`lsof -nP -iTCP:8893 -sTCP:LISTEN` shows one entry).

- [ ] **Step 3: Login and stash a token**

```bash
TOKEN=$(curl -s -X POST http://127.0.0.1:8893/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"account":"superadmin@shopjoy.com","password":"password123"}' \
    | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['access_token'])")
```
Expected: long JWT token echoed.

- [ ] **Step 4: Verify create path persists the three fields**

```bash
RESP=$(curl -s -X POST http://127.0.0.1:8893/api/v1/promotions \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
        "name":"TASK-7",
        "description":"verify",
        "type":"discount",
        "start_time":"2026-08-01T00:00:00Z",
        "end_time":"2027-12-31T00:00:00Z",
        "usage_limit":100,
        "per_user_limit":5,
        "tags":["精选","双11"]
    }')
NEW_ID=$(echo "$RESP" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['id'])")
echo "  -> created $NEW_ID"
```
Expected: `data.id` echoed.

Then:
```bash
mycli -h 192.168.0.100 -P 3306 -u root -p123456 shopjoy -e "
SELECT id,name,usage_limit,per_user_limit,tags
FROM promotions WHERE id=$NEW_ID;"
```
Expected: row shows `usage_limit=100`, `per_user_limit=5`, `tags=["精选","双11"]`.

- [ ] **Step 5: Verify update path persists the three fields**

```bash
curl -s -X PUT "http://127.0.0.1:8893/api/v1/promotions/$NEW_ID" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
        "id":"'"$NEW_ID"'",
        "name":"TASK-7",
        "description":"verify",
        "type":"discount",
        "start_time":"2026-08-01T00:00:00Z",
        "end_time":"2027-12-31T00:00:00Z",
        "usage_limit":200,
        "per_user_limit":2,
        "tags":[""]
    }' > /dev/null

mycli -h 192.168.0.100 -P 3306 -u root -p123456 shopjoy -e "
SELECT usage_limit,per_user_limit,tags FROM promotions WHERE id=$NEW_ID;"
```
Expected: `usage_limit=200`, `per_user_limit=2`, `tags=NULL` (because `[""]` sanitizes to empty slice, which encode as SQL `null`).

- [ ] **Step 6: Verify GET round-trips the three fields**

```bash
curl -s -H "Authorization: Bearer $TOKEN" \
    "http://127.0.0.1:8893/api/v1/promotions?page=1&page_size=50" \
    | python3 -c "
import sys,json
d=json.load(sys.stdin)['data']
hit=[p for p in d['list'] if p['id']=='$NEW_ID'][0]
print('usage_limit=',hit.get('usage_limit'))
print('per_user_limit=',hit.get('per_user_limit'))
print('tags=',hit.get('tags'))"
```
Expected: `usage_limit=200`, `per_user_limit=2`, `tags=None` or `[]` (Python view of `null`).

- [ ] **Step 7: Verify negative value rejected**

```bash
curl -s -X PUT "http://127.0.0.1:8893/api/v1/promotions/$NEW_ID" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
        "id":"'"$NEW_ID"'",
        "name":"TASK-7",
        "type":"discount",
        "start_time":"2026-08-01T00:00:00Z",
        "end_time":"2027-12-31T00:00:00Z",
        "usage_limit":-1
    }' | python3 -c "import sys,json; print(json.load(sys.stdin).get('code'))"
```
Expected: prints `80016` (`ErrPromotionUsageLimitInvalid`).

- [ ] **Step 8: Verify existing rows after migration**

```bash
mycli -h 192.168.0.100 -P 3306 -u root -p123456 shopjoy -e "
SELECT id,name,usage_limit,per_user_limit,tags FROM promotions WHERE id < 10;"
```
Expected: every pre-existing row shows `usage_limit=0`, `per_user_limit=1`, `tags=NULL`.

- [ ] **Step 9: Stop the verify instance**

```bash
pkill -f "admin-api-verify"
sleep 1
lsof -nP -iTCP:8893 -sTCP:LISTEN 2>/dev/null && echo "still up" || echo "verify instance stopped"
```
Expected: `verify instance stopped`.

- [ ] **Step 10: Document & commit**

No code changes in this task — verification only. If a Task 1–6 implementation slipped and revealed an actual bug, fix and commit as a follow-up. Otherwise add the verification output to your final status report (not a separate commit).

---

### Task 8: Rollout handoff

**Files:** none (operational).

- [ ] **Step 1: Inform the user that the long-running GoLand `:8888` instance needs a restart**

The dev server under GoLand is running with the old binary. The three new columns are already in the DB (Task 1 step 2) — once the binary is restarted, all four layers will read/write the new fields.

- [ ] **Step 2: Stop here**

Do NOT auto-restart the user's IDE process — destructive and easy to recover but explicit consent is the project's rule. Document the manual step in your final summary and exit.

---

## Self-Review

After writing this plan I checked:

**Spec coverage:**
- §1 Schema → Task 1 ✅
- §2 Domain → Task 2 ✅
- §3 Persistence model + mapping + Update map → Task 3 ✅
- §4 Application request/response + Create/Update/`toPromotionResponse` → Task 4 ✅
- §5 Logic pass-through → Task 6 ✅
- §6 "no wire type changes" → noted in Global Constraints; not a task ✅
- §7 Error codes → Task 5 ✅
- §8 Verification matrix → Task 7 ✅
- §9 Rollout → Task 8 ✅
- §10 Out of scope → Task 4 Step 9 (lenient tag handling) and Task 8 Step 1 (no scope-side enforcement) ✅
- §11 Files-touched → all 6 paths covered in tasks above ✅

**No placeholders:** No "TBD", "TODO", "implement later", "fill in". Code blocks contain real content for every change.

**Type consistency:**
- `Promotion.UsageLimit int`, `Promotion.PerUserLimit int`, `Promotion.Tags []string` — defined once in Task 2, used consistently in Tasks 3, 4, 6.
- `UsageLimit int`, `PerUserLimit int`, `Tags []string` field names match across `CreatePromotionRequest`, `UpdatePromotionRequest`, `PromotionResponse`, and the persistence layer's `promotionModel`.
- `promotionModel.UsageLimit int` and `.PerUserLimit int` are then assigned via `int(p.UsageLimit)` and `int(m.UsageLimit)` casts — matching the existing `Scope promotion.ScopeType(m.ScopeType)` pattern in the same file.

**Notes:**
- Task 3 deliberately uses `string` (not `datatypes.JSON`) for the JSON column to mirror the existing `ScopeIDs string` style and avoid pulling in a new module dependency. The `gorm:"type:json"` tag is what makes MySQL treat it as a JSON column.
- Task 4 Step 7 + Step 9 calls `req.Tags = sanitizeTags(req.Tags)` BEFORE constructing the entity, so the response and DB both see the cleaned list.
- Task 7 Step 5 sanitization strips `[""]` to empty → DB NULL. Documented; treat this as expected behavior, not a regression.
