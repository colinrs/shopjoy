# SDD Progress: Promotion × Coupon Merge

**Plan:** `docs/superpowers/plans/2026-07-19-promotion-coupon-merge.md`
**Spec:** `docs/superpowers/specs/2026-07-19-promotion-coupon-merge-design.md`
**Branch:** main
**Started:** 2026-07-19
**Base commit:** TBD (recorded before Task 1 implementer dispatch)

## Tasks

- Task 1: Pre-migration verification SQL — DONE (commits 7231f6e..0a4171e, review clean after 1 controller-driven fix + 1 spec amendment)
  - Notes: implementer discovered `coupons.market_ids` column does not exist in live schema; controller removed CHECK 3 with rationale comment, then amended spec §2.5 #3 to no-op. Baseline row counts captured: 12 promotions, 11 rules, 12 coupons, 18 user_coupons.

- Task 2: P0 merge migration script — DONE (commits 0a4171e..a1db4b4, review clean after 3 implementer dispatches; Approved with 4 Important handoff notes)
  - Notes: live `promotions` table already had `scope_type/scope_ids/exclude_ids/usage_limit/per_user_limit/tags` from prior SDD run (2026-07-18). Migration patched to use dynamic SQL per column (MySQL 9.0 doesn't support `ADD COLUMN IF NOT EXISTS`; MariaDB-only). Live DB verified: 24 promotions (12 PROMOTION + 12 COUPON), 23 promotion_rules (11 PROMOTION + 12 COUPON), 18 user_coupons unchanged. Coupons archived as `_archived_coupons_20260719`.

- Task 3: Add ErrPromotionInvalidKind error code — DONE (commits 0dd18c0, b85b9d4)
  - Notes: brief specified code 80010 but 80010 was already taken by ErrPromotionScopeInvalid; used 80018 instead. Spec/plan updated to reflect 80018.

- Task 4: Domain layer rewrite — DONE (commit f3fd57d, follow-up fix 8d3bfce)
  - Notes: implementer kept scope inline in entity.go instead of separate scope.go (minimized churn); deleted legacy `pkg/domain/promotion/service.go` (referenced removed Coupon types); `pkg/domain/promotion/` builds clean. Controller fixed `admin/cmd/migrations/main.go` to import from new unified `pkg/domain/promotion`. Remaining build errors are localized to `admin/internal/application/promotion/coupon_app.go` (Task 6 will delete) and `admin/internal/logic/coupons/*` (Task 7 will fix) — exactly as the plan predicted.

- Task 5: Repository rewrite — DONE (commit 5de36b7)
  - Notes: implementer also deleted orphaned `user_coupon_repository.go` and `promotion_usage_repository.go` (their interfaces were removed by Task 4 domain unification; brief listed only `coupon_repository.go`). Used `string` with `type:json` instead of `datatypes.JSON` (not in go.mod); added `normalizePage` helper since `UserCouponQuery`/`UsageQuery` lack embedded `shared.PageQuery`. `persistence/` builds clean. Remaining errors in `application/promotion/coupon_app.go` (Task 6 will delete) and `service_context.go` (Task 6 will rewire to `NewPromotionRepository`).

- Task 6: Application merge — DONE (commit 899658c)
  - Notes: `coupon_app.go` deleted; `PromotionApp` unified with all spec §5.1 methods (Create/Update/Get/List/Delete/Activate/Deactivate + rule methods + COUPON-specific methods). `service_context.go` rewired to `NewPromotionRepository` + `NewPromotionApp`; old `CouponRepository`/`CouponApp` fields removed. `application/` + `svc/` build clean. Remaining errors are all in `logic/coupons/*` + `logic/promotions/*` (Task 7 territory).

- Task 7: Logic layer update — DONE (commit a39529a, follow-up fix 9f83193)
  - Notes: implementer covered `logic/coupons/*` and `logic/promotions/*` per brief but explicitly excluded `logic/user_coupons/` (out-of-scope per plan). Controller fixed `logic/user_coupons/{helper,issue,batch_issue,list_user}_coupon_*.go` to route through `PromotionApp.{IssueToUser,BatchIssue,ListUserCoupons}`. Build now 0 errors across the entire `admin/` tree.

- Task 8: API definition + codegen — DONE (commits 0d666a8, c956a39)
  - Notes: API file updated with Kind/MarketID/Code/TotalCount/Rules fields on unified `PromotionDetailResp`/`CreatePromotionReq`/`UpdatePromotionReq`; rule wire types replaced (ConditionType/ConditionValue/ActionType/ActionValue); `OwnerKind` + `OwnerID` path params on `CreatePromotionRulesReq`. `make api` regenerated types cleanly. Post-codegen logic cleanup fixed 9 files (legacy `PromotionID` → `OwnerID`/`OwnerKind`, legacy `RuleType`/`Value`/`DiscountType`/`DiscountValue` → `ConditionType`/`ConditionValue`/`ActionType`/`ActionValue`, removed `ScopeType` from `CouponDetailResp`). `make build` 0 errors.

- Task 9: Frontend sync — DONE (commit 774a104)
  - Notes: `CouponDetailResp` collapsed into `Promotion` interface; `kind` discriminator added; `market_id`, `code?`, `total_count?`, `rules?` fields; rule editor pages branch on `kind`. Vue-tsc + Vite build 0 errors. Concerns: some i18n keys may be missing (minor); legacy `/api/v1/coupons/*` endpoints preserved for backward compat (per plan).

- Task 10: Final verification — DONE
  - Notes: All 10 tasks complete. `go build ./admin/... ./pkg/...` clean (0 errors). `go test ./pkg/domain/promotion/... ./admin/internal/application/promotion/... ./admin/internal/infrastructure/persistence/...` passes. Live DB verified: 24 promotions (12 PROMOTION + 12 COUPON), 23 promotion_rules (11 PROMOTION + 12 COUPON), 18 user_coupons intact.

## Decisions log

- Pre-migration check #3 against `coupons.market_ids` removed — the column was never created in live schema. Per spec §2.7, promotion.market_id is a new column with no backfill, so there is no data to lose.
- Live `promotions` table pre-existed with columns from prior SDD run. Migration uses dynamic SQL (PREPARE/EXECUTE) for `ADD COLUMN` to avoid duplicate-column errors and remain idempotent.
- Coupon→promotion INSERT hardcodes `usage_limit = 0` because the source `coupons` table has no `usage_limit` column (only `per_user_limit`). Coupon issuance is bounded by `total_count`, not `usage_limit` (handoff I2 to Task 6+).
- `user_coupons.coupon_id` still references `_archived_coupons_20260719.id` numerically; COUPON rows live in `promotions` with the same AUTO_INCREMENT values. Application layer must JOIN via `promotions WHERE kind='COUPON'` (handoff I3 to Task 5/6).

## Handoff notes for downstream tasks

- **I1 → Task 6+ App layer**: `_deprecated_promotion_rules` has `owner_id=0` for all 11 PROMOTION rules (backup captured after column add but before backfill). Rollback via spec Appendix A still works since raw column data is preserved.
- **I2 → Task 6+ App layer**: Migrated COUPON rows have `usage_limit=0` ("unlimited"). App layer must special-case COUPON kind to ignore `usage_limit` and use `total_count` for inventory bounds.
- **I3 → Task 5/6 Repository + App**: `user_coupons.coupon_id` joins must go through `promotions WHERE kind='COUPON'` (numerical IDs match, but the table name changed). Add repository helper or app-layer query.
- **I4 → Task 5/6 Repository**: Index rebuild uses 3 dynamic statements (CREATE + 2 DROP) instead of one ALTER. Trade-off accepted for idempotency.

## Rollout handoff

**Manual step required:** Restart the GoLand :8888 admin instance so the new binary is loaded. The DB columns and data are already in place (Task 2).

**Commits in order (chronological):**
1. `dfb779b` sql(promotion): add pre-migration verification checks
2. `0a4171e` docs(spec): mark pre-migration CHECK 3 as no-op
3. `a1db4b4` feat(promotion): merge coupons into promotions (Task 2)
4. `0dd18c0` feat(code): add ErrPromotionInvalidKind
5. `b85b9d4` docs(spec/plan): update ErrPromotionInvalidKind code 80010→80018
6. `f3fd57d` refactor(domain/promotion): unify Promotion + Coupon into single struct
7. `8d3bfce` fix(migrations): update for promotion×coupon merge
8. `5de36b7` refactor(persistence): unify promotion + coupon repositories
9. `899658c` refactor(application/promotion): unify promotion + coupon into single App
10. `a39529a` refactor(logic): route promotion + coupon logic through unified PromotionApp
11. `9f83193` fix(logic/user_coupons): route through unified PromotionApp (Task 7 gap)
12. `0d666a8` fix(logic): post-codegen logic cleanup
13. `c956a39` feat(api): unify PromotionDetailResp with kind + market_id + rules
14. `774a104` feat(frontend): unify Promotion + Coupon rendering with kind discriminator

**Base SHA:** `7231f6e` (last commit before spec/design).

**Final state:**
- Live MySQL: 24 promotions (12 PROMOTION + 12 COUPON), 23 promotion_rules (11 PROMOTION + 12 COUPON), 18 user_coupons intact.
- `coupons` table archived as `_archived_coupons_20260719` (30-day retention per spec).
- `_deprecated_coupons` + `_deprecated_promotion_rules` backups available for rollback (spec §A).
- `make build` + `go test` clean.

**Known follow-ups (recorded but out-of-scope for this SDD):**
- CouponDetailResp lost `ScopeType` field during Task 8 codegen — coupon create/update defaults to `STOREWIDE` scope. If COUPON-scope creation matters, restore `ScopeType` to `CouponDetailResp` + add `ScopeType`/`ScopeIDs`/`ExcludeIDs` fields to `CreateCouponReq`/`UpdateCouponReq`.
- Some i18n keys may be missing for new `kind`/`market_id` text labels (Task 9 concern; minor).
- Archived `_archived_coupons_20260719` and `_deprecated_*` tables should be dropped after a 30-day observation period (spec §6.3).