# Promotion usage_limit / per_user_limit / tags — Persistence Design

- **Date:** 2026-07-18
- **Status:** Approved (brainstorming complete, awaiting implementation)
- **Scope:** Wire-level fields exist on promotion create/update payloads but are silently dropped between the logic and persistence layers. Coupon table already has these columns; bring the promotions table to parity.
- **Out of scope:** enforcement at order/checkout time, per-user `promotion_usage` table, frontend tag-edit UI.

## 1. Schema

Add three columns to `sql/promotion/schema.sql`. The project consolidates migrations into `schema.sql` (no per-migration files) — see `.claude/rules/document/README.md`.

```sql
ALTER TABLE promotions
    ADD COLUMN usage_limit    INT     NOT NULL DEFAULT 0
        COMMENT '0 = unlimited',
    ADD COLUMN per_user_limit INT     NOT NULL DEFAULT 1
        COMMENT 'per-user cap; 0 = unlimited',
    ADD COLUMN tags           JSON    NULL
        COMMENT 'free-form labels';
```

Defaults:
- `usage_limit = 0` → unlimited
- `per_user_limit = 1` → one-per-user
- `tags = NULL` → no tags

Existing rows pick up the defaults, which match the "unlimited" semantics that the current code incidentally has (because the field is missing). No data backfill needed.

## 2. Domain entity

`pkg/domain/promotion/entity.go` extends the `Promotion` struct:

```go
type Promotion struct {
    // ... existing fields ...
    UsageLimit    int             `json:"usage_limit"`
    PerUserLimit  int             `json:"per_user_limit"`
    Tags          []string        `json:"tags,omitempty" gorm:"type:json"`
    // ... rest unchanged ...
}
```

The `gorm:"type:json"` tag is a hint; the actual MySQL-side persistence uses GORM's standard JSON column handling on the persistence-side `promotionModel`. The struct tag here is informational.

## 3. Persistence

`admin/internal/infrastructure/persistence/promotion_repository.go`:

```go
import "gorm.io/datatypes" // or fall back to string + json marshaller if absent

type promotionModel struct {
    // ... existing columns ...
    UsageLimit    int               `gorm:"column:usage_limit;not null;default:0"`
    PerUserLimit  int               `gorm:"column:per_user_limit;not null;default:1"`
    Tags          datatypes.JSON    `gorm:"column:tags;type:json"`
    // ... rest unchanged ...
}
```

`toEntity()` / `fromPromotionEntity()` mirror `coupons_repository.go`'s `ScopeIDs string + json.Marshal/Unmarshal` approach if `gorm.io/datatypes` is not on the module path. The persistence layer is the single place that does the JSON string ↔ `[]string` translation.

The `Update(ctx, ...)` method's `Updates(map[string]any{...})` block must include the new columns (`"usage_limit"`, `"per_user_limit"`, `"tags"`) so writes actually land.

## 4. Application layer

`admin/internal/application/promotion/promotion_app.go`:

### Request types

```go
type CreatePromotionRequest struct {
    // ... existing ...
    UsageLimit    int
    PerUserLimit  int
    Tags          []string
}

type UpdatePromotionRequest struct {
    // ... existing ...
    UsageLimit    int
    PerUserLimit  int
    Tags          []string
}
```

### Response type

```go
type PromotionResponse struct {
    // ... existing ...
    UsageLimit    int        `json:"usage_limit"`
    PerUserLimit  int        `json:"per_user_limit"`
    Tags          []string   `json:"tags"`
}
```

### Persistence wiring

- `promotionApp.CreatePromotion`: assign `p.UsageLimit = req.UsageLimit`, etc., when constructing `pkgpromotion.Promotion`.
- `promotionApp.UpdatePromotion`: assign the same three fields onto `p`. The existing `if p.Status == StatusActive { return ErrPromotionCannotDelete }` guard is preserved.
- `toPromotionResponse`: read the new fields off `p` and put them on the response.

## 5. Logic layer

`admin/internal/logic/promotions/update_promotion_logic.go` and `create_promotion_logic.go`:

When constructing `apppromotion.UpdatePromotionRequest` / `CreatePromotionRequest`, copy the three fields off the wire request:

```go
updateReq := apppromotion.UpdatePromotionRequest{
    // ... existing ...
    UsageLimit:   req.UsageLimit,
    PerUserLimit: req.PerUserLimit,
    Tags:         req.Tags,
}
```

These changes mirror what Phase 1 / 2 already did for `type` and `scope`. No new helper functions are needed.

## 6. Wire types

`admin/internal/types/types.go`:
- `UpdatePromotionReq`, `CreatePromotionReq`, `PromotionDetailResp` — **no changes needed**. They already carry `UsageLimit` / `PerUserLimit` / `Tags` from earlier refactors; the gap was purely on the persistence side.

## 7. Error handling

Add to `pkg/code/code.go` (Promotion Module 80xxx):

```go
ErrPromotionUsageLimitInvalid    = &Err{HTTPCode: http.StatusBadRequest, Code: 80016, Msg: "promotion usage_limit must be >= 0"}
ErrPromotionPerUserLimitInvalid  = &Err{HTTPCode: http.StatusBadRequest, Code: 80017, Msg: "promotion per_user_limit must be >= 0"}
```

Reuse for both Create and Update paths — they're identical validation rules.

`tags` validation: lenient. Strip empty strings and truncate entries over 64 chars before persistence. Do not 400 on bad tag input.

## 8. End-to-end verification matrix

Run against a fresh admin binary on a non-conflicting port (same pattern as previous phases).

| Case | Expectation |
|---|---|
| `POST /promotions` with `usage_limit=100`, `per_user_limit=5`, `tags=["精选","双11"]` | All three columns written; `GET` returns same values |
| `PUT /promotions/:id` updates `usage_limit=200` | DB row shows `usage_limit=200` |
| `tags=[]` on save | DB `tags=NULL` (Go empty slice → JSON `null`) |
| Omit `usage_limit` from payload | Default 0; no error |
| Existing 7 promotion rows pre-migration | After ALTER, behave as "unlimited"; `GET` returns defaults |
| Lists `GET /promotions` | Each item carries the three new keys; old clients ignore them |

## 9. Migration & rollout

1. `make build` from project root (per `.claude/CLAUDE.md`).
2. Schema migration: `ALTER TABLE promotions ...` applied manually to the dev DB; mirrors `sql/promotion/schema.sql`.
3. Restart the long-running `admin` process (currently `:8888` from GoLand) so it picks up the new binary.
4. Frontend hot-reloads; no frontend changes in this ticket (`tags` has no UI yet, but the data path is plumbed end-to-end so a future UI ticket can use it).

## 10. Out of scope (deferred tickets)

| Item | Why deferred |
|---|---|
| `promotion_usage(promotion_id, user_id, used_count)` table + checkout enforcement of `usage_limit`/`per_user_limit` | Different layer (redemption flow); needs its own spec. Brainstorming confirmed user wants this as a separate ticket. |
| `used_count` auto-update on order | Same — order-side concern. |
| Frontend tags input UI | No current UX; not in user's ask. |
| Inline validation annotations on schema columns (`@`, `gorm:"comment:..."` exact syntax) | Implementation detail picked during the write step. |

## 11. Files touched

- `sql/promotion/schema.sql` — schema columns
- `pkg/domain/promotion/entity.go` — Promotion fields + `IsValid`-style checks if any (none currently)
- `pkg/code/code.go` — two new error codes
- `admin/internal/infrastructure/persistence/promotion_repository.go` — model + `toEntity`/`fromPromotionEntity` + Update map
- `admin/internal/application/promotion/promotion_app.go` — request/response types + Create/Update/`toPromotionResponse`
- `admin/internal/logic/promotions/update_promotion_logic.go` + `create_promotion_logic.go` — pass-through

No frontend changes. No `types/types.go` changes. No new error codes beyond the two above.
