# SDD Progress: Shipping Internationalization (Phase 1 / Batch 1)

**Plan:** `docs/superpowers/plans/2026-07-20-shipping-internationalization.md`
**Phase Scope:** Phase 1 only (Tasks 1.1 ~ 1.11)
**Branch:** main
**Started:** 2026-07-20
**Base commit:** `6a3dce8` (HEAD before starting)

## Lessons carried over from storefront-theme-preview SDD

- **Silent field loss in auto-generated logic layer**: previous SDD caught a bug where logic-layer constructors omitted `DefaultConfig: t.DefaultConfig` even though wire types had the field. Compile + unit tests passed. Only end-to-end browser verification caught it. **Risk for this SDD**: Task 1.7 and 1.8 add many fields (Taxable/TaxRate/RemoteSurcharge/...); reviewer must explicitly verify each wire field has corresponding assignment in logic.
- Pre-existing dirty files in working tree: `shop-admin/src/views/settings/markets/index.vue`, `shop-admin/src/views/shop/index.vue` — out of scope, do not touch.

## Tasks

- Task 1.1: DONE (commit e5aa5d3, review approved — 0 Critical, 2 Important, 5 Minor)
  - Important follow-ups scheduled for Task 1.7: cover `SetDefault`/`CanDelete` with new MarketID layout
  - Minor #2 (missing trailing newline in entity_test.go) — fixed inline 2026-07-20
  - Minor #3-7: non-blocking, recorded below for context
- Task 1.2: DONE (commit 74b8db4 + fix 2aef22e, re-review APPROVED)
  - 9 fields added to ShippingZone: Currency/NameI18n/Taxable/TaxRate/TaxIncluded/IossApplicable/RemoteSurcharge/RemoteZipPatterns/FuelSurchargePct/VolumetricDivisor
  - NameI18n + RemoteZipPatterns use `json.RawMessage` placeholder; Task 1.3 will swap to typed value objects
  - Fix commit 2aef22e strengthened test assertions (bytes.Equal + TaxIncluded logic)
- Task 1.3: DONE (commit 13aefd7, review approved — 0 Critical/Important, 3 Minor)
  - StringI18n + StringArray value objects created with Value/Scan/Get/Contains
  - Get() hardened with empty-locale guard (documented deviation from brief — non-defect)
  - Minor #1 (missing trailing newlines) — fixed inline post-review commit (tracked separately)
- Task 1.4: DONE (commit 05437b3 + fix c3a960f, review APPROVED with documented concerns)
  - Migration applied to live DB: 11 new columns on shipping_zones, 4 on shipping_templates
  - schema.sql merged with new columns in struct-aligned order
  - Concern #3 (live DB missing idx_warehouse_id) — fixed via separate migration c3a960f
  - Known follow-ups (non-blocking): Concern #2 (idempotency pattern), Concern #1 (DB column physical order vs struct order)
- Task 1.5: DONE (commit ef8f39e, review approved — 0 Critical/Important, 2 Minor)
  - FindDefaultByMarket + FindListByMarket added with market_id fallback to market_id=0
  - UnsetAllDefaultByMarket added as sibling (kept original UnsetAllDefault for backward compat)
  - 8 sqlmock unit tests cover: find default, market fallback, not-found, list filtering, pagination
  - Trailing newline fixed inline
- Task 1.6: DONE (commit 4e03862, review approved — 0 Critical/Important, 6 Minor)
  - .api fully rewritten per brief; types.go regenerated (+278/-148)
  - Build failures isolated to logic layer (3 errors) — gate for Task 1.7/1.8 as planned
  - ⚠️ Forward notes:
    - Task 1.7 must also fix `calculate_shipping_fee_logic.go` (CalculatorItem.ProductID/SKUID string flip + FindMappingByTarget call path)
    - Task 1.8 implementer must compile `internal/logic/shipping_zones/...` directly — umbrella build masks zone errors
    - Task 1.9 must remove `product_count`/`category_count` references in frontend (removed from ShippingTemplateListItem)
    - Task 1.7 must apply CNY default in logic (wire layer doesn't default Currency)
- Task 1.7: DONE (commit 11e35ba + fix c6b6366, review approved — 0 Critical, 2 Important, 4 Minor)
  - 7 logic files updated; MarketID/Currency/CarrierCode/WarehouseID all wired through (audit confirmed no silent drops)
  - 5 silent field drops discovered and fixed: Get/Update/SetDefault responses, List removed deprecated fields, Calculator CNY hardcode
  - Calculator: ProductID/SKUID string→int64 conversion with ErrSharedInvalidParam on failure
  - Error codes: 230006/230007 used (project actual range, brief was wrong with 120xxx)
  - Fix commit added TemplateName to calc response + trailing newlines; `parsed` variable confirmed not dead code (used by findTemplateForItems)
- Task 1.8: DONE (commit 2cf3e44, review approved — 0 Critical/Important, 5 Minor)
  - FeeTypeByVolume constant + IsValidV2() + Validate() by_volume check
  - CreateShippingZoneLogic: all 21 wire fields + 2 derived (TenantID, MarketID) = 23 entity fields
  - UpdateShippingZoneLogic: all 20 wire fields optional-handled
  - **buildZoneDetails now populates all 24 ShippingZoneDetail fields** (reviewer corrected brief's "22" — actual is 24)
  - i18n_helper.go created with exported ToStringI18n/FromStringI18n
  - Agent hit 429 rate-limit on response but committed work successfully
- Task 1.9: DONE (commit 72a88be + fix 0cb7734, review approved — 0 Critical, 2 Important out-of-scope, 2 Minor)
  - api/shipping.ts fully rewrote types to match backend (24 ShippingZoneDetail fields, NameI18nEntry, CalculatorAddress with country_code, etc.)
  - Removed product_count/category_count (no longer in wire)
  - Fix commit migrated 2 view files: calculator (added market selector for market_id + country_code) and TemplateCard (removed deprecated stat blocks)
  - `pnpm build` green
- Task 1.10: DONE (commit 601df90, review APPROVED)
  - useCurrency composable + vitest tests (3/3 passing)
  - Supports 12 currencies with correct decimals and position (VND after, JPY/KRW/IDR/VND no decimals)
- Task 1.11: DONE (commit 6631e22, review APPROVED — 0 Critical/Important, 3 Minor)
  - Added `ErrShippingCalcMarketRequired` (Code=230308, HTTPCode=400)
  - Added MarketID required validation at top of CalculateShippingFee
  - Note: CNY hardcode was already removed in Task 1.7; this task focused on MarketID guard
  - Doc follow-up pending: error-codes.md row for 230308

## Final state

**Final whole-branch review verdict: NOT READY TO MERGE**

Per-task reviews all passed (per Task reviews), but final whole-branch review found **4 Critical** + **6 Important** issues that block safe merge. The Phase 1 schema/field additions were done correctly, but several **calculation paths, security checks, and UI forms are incomplete** despite the wire types claiming the features exist.

### Critical (block merge)

1. **C1 — `by_volume` 计费未实现**: `ShippingZone.CalculateFee()` 的 switch 没有处理 `FeeTypeByVolume` —— 所有 by_volume 区域都会退化为只收 FirstFee，体积变化不影响运费。P1-9 实质未实现。
2. **C2 — CalculateShippingFeeResp 静默丢字段**: `Tax/Total/PriceIncludesTax/CarrierCode/EstimatedDays` 在 wire 中声明，但 calculator logic 只填 `ShippingFee`。国际订单会显示错误的总价/税费/含税标识。
3. **C3 — `FindDefaultByMarket`/`FindListByMarket`/`UnsetAllDefaultByMarket` 缺 `tenant_id` 过滤**: 跨租户数据访问风险 —— 同一 market_id 下不同租户的模板可被读取/批量取消默认状态。**安全级别问题**。
4. **C4 — zone currency vs template currency 不一致未拒绝**: zone 创建时若传入与 template 不同 currency，逻辑不报错；响应可能用 template.Currency 显示 zone 配置的另一种货币金额。容易产生"标价币种 ≠ 计费币种"bug。

### Important (should fix before merge)

- Calculator 视图缺市场下拉/国家选择 UI（market_id/country_code 永远为空）
- ZoneConfigForm 未添加 P1 字段（tax/surcharge/i18n/dimensions 等），管理后台无法创建/编辑 P1 功能
- Calculator 强制要求 city_code，国际地址通常只有 country_code/postal_code → 走不下去
- bool 字段（Taxable/TaxIncluded/IossApplicable）用 value 而非 pointer，Update 无法显式设为 false
- TaxRate 解析失败时静默归零
- 模板列表 N+1 查询（每个模板单独 CountZonesByTemplateID）

### Recommended path forward

Stop here. Create **Batch 1.5 = Phase 1 Patch** to fix C1-C4 + form/UI gaps before Phase 2 (国家代码 + 区域匹配). Otherwise Phase 2 will compound on top of broken calculation logic and tenant-isolation bugs.

## Commits in chronological order

1. `e5aa5d3` feat(shipping): add MarketID and Currency fields to ShippingTemplate
2. `74b8db4` feat(shipping): add Currency/Taxable/i18n/surcharge fields to ShippingZone
3. `2aef22e` fix(shipping): strengthen ShippingZone tests for i18n/surcharge fields
4. `13aefd7` feat(shipping): add StringI18n and StringArray value objects
5. `75e4ef4` fix(shipping): add trailing newlines to value_objects files
6. `05437b3` feat(shipping): add internationalization columns (currency/market/tax/i18n/surcharge)
7. `c3a960f` fix(shipping): add missing idx_warehouse_id to live DB (sync with schema.sql)
8. `ef8f39e` feat(shipping): add market-aware template lookup with fallback
9. `33c08d5` fix(shipping): add trailing newline to shipping_repository_test.go
10. `4e03862` feat(shipping): rewrite API types for internationalization (build broken, awaits 1.7/1.8)
11. `11e35ba` fix(shipping-templates): wire MarketID/Currency/CarrierCode/WarehouseID through CRUD logic
12. `c6b6366` fix(shipping-templates): add TemplateName to calc response + trailing newlines
13. `2cf3e44` feat(shipping-zones): wire 9 P1 fields through CRUD logic + by_volume support
14. `72a88be` feat(shipping): sync TypeScript types with backend
15. `0cb7734` fix(shipping-ui): migrate views to internationalization API types
16. `601df90` feat(shipping): add useCurrency composable with 12-currency support
17. `6631e22` feat(shipping-calc): require market_id + add ErrShippingCalcMarketRequired
18. `f91f16a` fix(shipping): add tenant_id filter to market-aware repo methods (C3)
19. `1261657` fix(shipping): implement by_volume billing with volumetric weight
20. `d187ba8` fix(shipping-calc): populate Tax/Total/CarrierCode in response
21. `48875a9` fix(shipping): enforce currency + tenant consistency on zone create
22. `ea51e7a` fix(shipping): pointer bools for zone updates + batch zone count + strict tax rate
23. `a30c7fd` shipping UI changes (ZoneConfigForm P1 + Calculator market/country selectors) — committed under misleading promotion-mislabeled commit message; content is correct
24. `ea07bb0` docs: add 230107 invalid tax rate and 230308 market_id required

## Decisions log

- **Branch: main** — user chose to commit directly to main (explicit override of skill default). Reason: project is currently single-developer, no parallel feature work in flight.
- **Phase 1 final verdict: NOT READY TO MERGE** → **Patch Batch 1.5 created and executed** to fix Critical gaps. **NOW READY TO MERGE** subject to known follow-ups (see Notes).

## Notes

- Task 1.1 review (Minor #4): composite index `idx_market_default` will be created as `(market_id, is_default)` because that's struct field order. Matches "find default within a market" use case. Task 1.4 must keep column order aligned.
- Task 1.1 review (Important #1): SetDefault/UnsetDefault/CanDelete not exercised by Task 1.1 tests. Schedule coverage in Task 1.7 when logic layer is wired.
- `application.Model` provides `SetDefault` / `UnsetDefault` / `CanDelete` — implementer should check whether `CanDelete` logic needs MarketID-aware update (e.g., per-market default).