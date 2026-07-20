# SDD Progress: Storefront Theme Preview

**Plan:** `docs/superpowers/plans/2026-07-19-storefront-theme-preview.md`
**Spec:** `docs/superpowers/specs/2026-07-19-storefront-theme-preview-design.md`
**Branch:** main
**Started:** 2026-07-19
**Base commit:** `f5441eb` (the plan commit)

## Tasks

- Task 1: DONE (commits f5441eb..848e3e5, review clean — 0 Critical/Important, 3 Minor observations only)
- Task 2: DONE (commit cd5d3e0, review clean — 0 Critical/Important, 1 Minor observation only)
- Task 3: DONE (commit 2fdb3c7, review clean — 0 Critical/Important, 4 Minor observations only; verbatim transcription confirmed)
- Task 4: DONE (commit 1871886, review clean — 0 Critical, 0 Important, 0 Minor)
- Task 5: DONE (no commits; machine-verifiable steps 1, 2, 7 passed; bundle grep confirms ThemePreviewCard compiled in; Steps 3-6 marked PENDING HUMAN — visual/layout/console/DB-fallback)

## Final state

- Branch is 7 commits ahead of `f5441eb` (the plan commit): 848e3e5, cd5d3e0, 2fdb3c7, 1871886, 4c5e347, 99ab112, 16a0c0c
- Backend `admin/`: `make build` clean; binary rebuilt and restarted after each fix
- Frontend `shop-admin/`: `pnpm build` clean; `merriweather` FONT_MAP key found in `dist/assets/index-*.js`
- **Human verification: PASSED** (user confirmed "正常了") — 5 preset themes render distinct, dialog renders themed mock
- Admin service running on `:8888` with the final binary

## Commits in chronological order

1. `848e3e5` feat(storefront): populate DefaultConfig in ListThemes response
2. `cd5d3e0` feat(storefront-ui): expose default_config on ThemeItem type
3. `2fdb3c7` feat(storefront-ui): add ThemePreviewCard mock component
4. `1871886` feat(storefront-ui): use ThemePreviewCard in theme cards
5. `4c5e347` fix(storefront): drop redundant double-marshal in themeRepo.Update (final-review fix)
6. `99ab112` feat(storefront-ui): use ThemePreviewCard in preview dialog (out-of-scope add-on after user feedback)
7. `16a0c0c` fix(storefront): wire DefaultConfig through logic layer (post-verification bug)

## Post-verification bug captured

After the SDD ran clean (all per-task reviewers approved, final review approved), browser verification revealed themes still rendered in the fallback palette. Root cause was a **silent-field-loss in the auto-generated logic layer** (`internal/logic/themes/list_themes_logic.go` and `get_current_theme_logic.go`): the wire-type constructors omitted `DefaultConfig: t.DefaultConfig`, even though the service returned it and the wire struct (`types.ThemeItem`) had the field. The bug class is exactly what `CLAUDE.md` flags for the promotion module as the #1 risk.

Detection required manual browser verification — no test would have caught this because the missing field doesn't produce a build error, log warning, or panic. Per-task reviewers trusted the implementer reports (which only confirmed `make build` and `pnpm build`) and the final review caught the type-mismatch but missed the logic-layer field drop.

**Lesson for future SDDs in this project:** even after final-review approval, the controller should make one verification pass with actual API responses end-to-end before declaring done. The `defaultConfigToWire` helper is small enough to be reviewed by inspection; the logic-layer omission is the kind of thing only end-to-end data flow reveals.

## Known follow-ups (recorded but out-of-scope for this SDD)

- Live config edit preview: pass `configOverride` prop to `ThemePreviewCard`; CSS-var architecture makes it easy.
- Real screenshot assets: not pursued.
- Custom-theme (non-preset) themes with NULL `default_config` render in the fallback blue palette — this is correct per spec ("missing config → fallback"), but the 4 Chinese custom themes the user has look identical to each other. If they want those to differentiate, the seed (`sql/storefront/schema.sql:228`) needs default_config values inserted, or the frontend could offer a config-edit UI.
- Minor observations from per-task reviews (e.g., `aria-hidden` on decorative `.tpc-dot` spans, `type="button"` on the `<span>` underline variant, missing trailing newline in component file) — left as-is, none are defects.

## Decisions log

_(filled as work proceeds)_

## Notes

- Working tree has pre-existing dirty files (`shop-admin/src/views/settings/markets/index.vue`, `shop-admin/src/views/shop/index.vue`) — out of scope, do not touch.
- `.gitignore` was modified to drop `docs/superpowers` — no `-f` needed for spec/plan commits going forward.
- Seed SQL (`sql/storefront/schema.sql:234-242`) inserts `default_config` as flat DTO format; helper must handle both nested and flat formats.