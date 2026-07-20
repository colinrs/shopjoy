# Promotion Activation Toggle Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Move promotion activate/deactivate from the action column to an `el-switch` toggle in the status column so the UX matches the Coupons tab.

**Architecture:** Single-file Vue template + script change in `shop-admin/src/views/promotions/index.vue`. Reuse the existing `handleToggleCouponStatus` pattern (optimistic update + per-row loading state + rollback on failure). Backend is already complete — no Go changes, no API changes, no locale changes.

**Tech Stack:** Vue 3 + Element Plus + TypeScript, vue-i18n.

## Global Constraints

- **Spec location:** `docs/superpowers/specs/2026-07-18-promotion-activation-toggle-design.md` — read fully before each task.
- **Build command:** Always `cd shop-admin && make build` (or `make build` from repo root, which builds both). Never `vite build` / `tsc` directly.
- **Frontend rules:** Frontend enums must match backend `.api` files exactly. `any` types forbidden. Money/price strings at API boundary.
- **Error handling:** Use `useErrorHandler` composable + i18n keys; do not surface raw `Error.message` to users.
- **i18n:** Reuse existing keys `promotions.activatedStatus`, `promotions.unactivated`, `promotions.activateSuccess`, `promotions.deactivateSuccess`, `promotions.activatePromotionFailed`, `promotions.deactivatePromotionFailed`. Do not add new keys.
- **Commit style:** Conventional commits. End with `Co-Authored-By: Claude <noreply@anthropic.com>`.

---

## File Structure

### Frontend — modify (single file)
- `shop-admin/src/views/promotions/index.vue`
  - **Template (~line 479-525):** Status column `el-tag` → `el-switch` with `el-tag` fallback for `ended`. Remove two link buttons from actions column.
  - **Script:** Add `promotionToggleLoading` reactive and `handleTogglePromotionStatus` function near the existing coupon toggle block (around line 1630-1656). Remove the now-unused `handleActivatePromotion` / `handleDeactivatePromotion` functions (lines 1658-1676).

### No other files change
- `shop-admin/src/api/promotion.ts` — `activatePromotion` / `deactivatePromotion` already exist.
- `admin/...` — backend untouched.
- `shop-admin/src/locales/{zh,en}.json` — no new keys.

---

## Task 1: Add `handleTogglePromotionStatus` and `promotionToggleLoading`

**Files:**
- Modify: `shop-admin/src/views/promotions/index.vue` (script block, after the existing `couponToggleLoading` declaration around line 1630)

**Interfaces:**
- Consumes: `activatePromotion(id)`, `deactivatePromotion(id)` from `@/api/promotion`; `useErrorHandler().handleError`; `ElMessage`; `t` from i18n; existing `loadPromotions()`.
- Produces: function `handleTogglePromotionStatus(row: Promotion, nextActive: boolean): Promise<void>`; reactive `promotionToggleLoading: Record<string, boolean>`.

- [ ] **Step 1: Locate the insertion point**

Open `shop-admin/src/views/promotions/index.vue` and find the line:

```typescript
const couponToggleLoading = reactive<Record<string, boolean>>({})
```

It is at approximately line 1630. The new code goes immediately **after** the existing `handleToggleCouponStatus` function and **before** the existing `handleActivatePromotion` function.

- [ ] **Step 2: Add `promotionToggleLoading` reactive + `handleTogglePromotionStatus`**

Insert the following block **immediately after** `handleToggleCouponStatus` and **before** `handleActivatePromotion`. Do **NOT** delete `handleActivatePromotion` or `handleDeactivatePromotion` in this task — their template callers still exist and are removed in Task 3:

```typescript
const promotionToggleLoading = reactive<Record<string, boolean>>({})

const handleTogglePromotionStatus = async (row: Promotion, nextActive: boolean) => {
  const wasActive = row.status === 'active'
  // 乐观翻转，避免网络慢时开关回弹
  row.status = nextActive ? 'active' : 'paused'
  promotionToggleLoading[row.id] = true
  try {
    if (nextActive) {
      await activatePromotion(row.id)
      ElMessage.success(t('promotions.activateSuccess'))
    } else {
      await deactivatePromotion(row.id)
      ElMessage.success(t('promotions.deactivateSuccess'))
    }
    loadPromotions()
  } catch (error) {
    // 失败回滚，保证 UI 与服务端状态一致
    row.status = wasActive ? 'active' : 'paused'
    handleError(
      error,
      nextActive ? t('promotions.activatePromotionFailed') : t('promotions.deactivatePromotionFailed')
    )
  } finally {
    promotionToggleLoading[row.id] = false
  }
}
```

- [ ] **Step 3: Verify the file still type-checks**

Run from the repo root:

```bash
cd shop-admin && make build 2>&1 | tail -30
```

Expected: build succeeds. The template buttons that call the old `handleActivatePromotion` / `handleDeactivatePromotion` are still in place (they get removed in Task 3), so those functions must also stay until Task 3 — that's why this task only **adds** the new function, it does not delete the old ones.

- [ ] **Step 4: Commit**

```bash
git add shop-admin/src/views/promotions/index.vue
git commit -m "feat(promotions): add handleTogglePromotionStatus for status switch

Mirrors handleToggleCouponStatus: optimistic update + per-row loading +
rollback on failure. The new function is not yet wired up to the UI
(template change + deletion of the now-redundant handleActivatePromotion
and handleDeactivatePromotion follow in subsequent commits).

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 2: Switch status column to `el-switch`

**Files:**
- Modify: `shop-admin/src/views/promotions/index.vue` (template, status column at lines 479-493; also widen column from 100 to 120; also remove the `// @ts-expect-error` directive above `handleTogglePromotionStatus` in the script block, around line 1659)

> **Carry-over from Task 1:** The new `handleTogglePromotionStatus` function declared in Task 1 is currently suppressed with a `// @ts-expect-error` directive because `tsconfig.json` has `noUnusedLocals: true` and the function is not yet called by the template. Once you wire the function into the template in this task, the directive becomes "unused" and vue-tsc will fail with `Unused '@ts-expect-error' directive`. **You must delete the entire `// @ts-expect-error` comment line** as part of this task (keep the `const handleTogglePromotionStatus = ...` declaration that follows). The comment is at approximately line 1659 and looks like:
>
> ```typescript
> // @ts-expect-error - Wired to template in subsequent commit; intentionally unused here.
> const handleTogglePromotionStatus = async (row: Promotion, nextActive: boolean) => {
> ```
>
> Delete the comment line. Do not replace it with anything.

- [ ] **Step 1: Locate the status column**

Find the existing promotions table status column (around line 479-493). It currently looks like:

```vue
<el-table-column
  :label="$t('promotions.statusColumn')"
  width="100"
  align="center"
>
  <template #default="{ row }">
    <el-tag
      :type="getPromoStatusType(row.status)"
      effect="light"
      size="small"
    >
      {{ getPromoStatusText(row.status) }}
    </el-tag>
  </template>
</el-table-column>
```

- [ ] **Step 2: Remove the `// @ts-expect-error` directive above `handleTogglePromotionStatus`**

The new function is about to be wired up by Step 3 of this task, so the temporary suppression must go. Find the line in the script block (around line 1659):

```typescript
// @ts-expect-error - Wired to template in subsequent commit; intentionally unused here.
const handleTogglePromotionStatus = async (row: Promotion, nextActive: boolean) => {
```

Delete the `// @ts-expect-error ...` line entirely. Leave the `const handleTogglePromotionStatus = ...` declaration untouched.

- [ ] **Step 3: Replace the column with the switch + tag-fallback version**

Change `width="100"` to `width="120"`, and replace the inner template body with:

```vue
<el-table-column
  :label="$t('promotions.statusColumn')"
  width="120"
  align="center"
>
  <template #default="{ row }">
    <el-switch
      v-if="row.status === 'active' || row.status === 'pending' || row.status === 'paused'"
      :model-value="row.status === 'active'"
      :loading="promotionToggleLoading[row.id] === true"
      :active-text="$t('promotions.activatedStatus')"
      :inactive-text="$t('promotions.unactivated')"
      inline-prompt
      @change="(val: boolean) => handleTogglePromotionStatus(row, val)"
    />
    <el-tag
      v-else
      :type="getPromoStatusType(row.status)"
      effect="light"
      size="small"
    >
      {{ getPromoStatusText(row.status) }}
    </el-tag>
  </template>
</el-table-column>
```

- [ ] **Step 4: Build to confirm template type-checks**

```bash
cd shop-admin && make build 2>&1 | tail -30
```

Expected: build succeeds. If the build fails with a "Property 'handleTogglePromotionStatus' does not exist" type error, return to Task 1 and confirm the function was added inside `<script setup lang="ts">`. If the build fails with i18n key errors, double-check that `promotions.activatedStatus` and `promotions.unactivated` exist in both `shop-admin/src/locales/zh.json` and `shop-admin/src/locales/en.json`.

- [ ] **Step 5: Commit**

```bash
git add shop-admin/src/views/promotions/index.vue
git commit -m "feat(promotions): switch status column to el-switch toggle

Replaces the static el-tag with an inline-prompt el-switch for rows in
pending/paused/active status. ended (and any other non-toggleable)
status still renders a tag.

Wider column (100 -> 120) to fit the inline prompt text. Switch i18n
keys reuse the Coupons tab values for consistent UX language.

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 3: Remove redundant "启用" / "停用" buttons from actions column

**Files:**
- Modify: `shop-admin/src/views/promotions/index.vue` (template, actions column at lines 509-525)

- [ ] **Step 1: Locate the action column block**

Find the promotions table actions column (around line 494-535). It currently has these two link buttons (after "编辑"):

```vue
<el-button
  v-if="row.status === 'pending' || row.status === 'paused'"
  type="success"
  link
  size="small"
  @click="handleActivatePromotion(row)"
>
  {{ $t('promotions.activate') }}
</el-button>
<el-button
  v-if="row.status === 'active'"
  type="warning"
  link
  size="small"
  @click="handleDeactivatePromotion(row)"
>
  {{ $t('promotions.deactivate') }}
</el-button>
```

- [ ] **Step 2: Delete the two buttons**

Delete both `el-button` blocks (8 + 8 = 16 lines including blanks) entirely. Leave "编辑" and "删除" untouched.

After deletion, the actions column should be just:

```vue
<el-table-column
  :label="$t('promotions.actionsColumn')"
  width="200"
  fixed="right"
>
  <template #default="{ row }">
    <el-button
      type="primary"
      link
      size="small"
      @click="handleEditPromotion(row)"
    >
      {{ $t('promotions.edit') }}
    </el-button>
    <el-button
      v-if="row.status !== 'active'"
      type="danger"
      link
      size="small"
      @click="handleDeletePromotion(row)"
    >
      {{ $t('promotions.delete') }}
    </el-button>
  </template>
</el-table-column>
```

- [ ] **Step 3: Widen the actions column**

Change `width="200"` to `width="180"` — the column only contains "编辑" and "删除" now, so 180 is enough.

- [ ] **Step 4: Delete the now-unused `handleActivatePromotion` and `handleDeactivatePromotion` script functions**

The two button blocks you just deleted in Step 2 were the only callers of these functions, so they are dead code. Delete them (they sit at approximately lines 1658-1676):

```typescript
const handleActivatePromotion = async (row: Promotion) => {
  try {
    await activatePromotion(row.id)
    ElMessage.success(t('promotions.activateSuccess'))
    loadPromotions()
  } catch (error) {
    handleError(error, t('promotions.activatePromotionFailed'))
  }
}

const handleDeactivatePromotion = async (row: Promotion) => {
  try {
    await deactivatePromotion(row.id)
    ElMessage.success(t('promotions.deactivateSuccess'))
    loadPromotions()
  } catch (error) {
    handleError(error, t('promotions.deactivatePromotionFailed'))
  }
}
```

> This is committed together with the template-button deletion (Steps 2-3) because they form one logical change: "the only callers of these functions are gone, so the functions are dead code."

- [ ] **Step 5: Build**

```bash
cd shop-admin && make build 2>&1 | tail -30
```

Expected: build succeeds. If the build complains about unused i18n keys (`promotions.activate`, `promotions.deactivate`), that is **expected** and can be ignored in this task — leave the keys in place because they may be referenced elsewhere or in future UI. (Vue-i18n does not error on unused keys; the build should pass cleanly.)

- [ ] **Step 6: Commit**

```bash
git add shop-admin/src/views/promotions/index.vue
git commit -m "refactor(promotions): drop redundant activate/deactivate buttons

State-column switch is now the single, discoverable toggle. The two
link buttons in the actions column would have duplicated the same
affordance and created inconsistent UI surfaces vs. the Coupons tab.

Also deletes the now-unused handleActivatePromotion and
handleDeactivatePromotion script functions (their only callers were
the buttons just removed) and narrows the actions column 200 -> 180
since it now holds only edit + delete.

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 4: Final build + manual smoke test

**Files:** none (verification only)

- [ ] **Step 1: Build the frontend**

```bash
cd shop-admin && make build 2>&1 | tail -30
```

Expected: build succeeds with no errors. (Warnings about unused i18n keys are acceptable; errors are not.)

- [ ] **Step 2: Lint pass**

```bash
cd shop-admin && make lint 2>&1 | tail -30
```

Expected: no new errors introduced by the three commits.

- [ ] **Step 3: Manual smoke test**

Start the admin dev server:

```bash
cd shop-admin && make dev
```

Open `http://localhost:3000/promotions` in a browser and verify against the spec acceptance criteria:

| Status of test row | Expected |
|--------------------|----------|
| `active` (进行中) | Switch ON, click turns it OFF, row re-renders as `paused` (已暂停) after reload |
| `pending` (待开始) | Switch OFF, click turns it ON, row re-renders as `active` (进行中) after reload |
| `paused` (已暂停) | Switch OFF, click turns it ON, row re-renders as `active` after reload |
| `ended` (已结束) | Tag only, no switch, not clickable |

Verify also:
- Switches show inline prompt text "已激活" / "未激活" (matches Coupons tab).
- During the in-flight request, the switch shows a loading spinner and is non-interactive.
- If the backend returns an error (simulate by stopping the dev server mid-click), the switch snaps back to its previous position and an error toast appears.
- The "操作" column now shows only "编辑" and "删除" — no "启用" / "停用".

- [ ] **Step 4: Confirm clean diff**

```bash
git log --oneline -5
git diff HEAD~3 HEAD --stat
```

Expected: three new commits on top of `main`, touching only `shop-admin/src/views/promotions/index.vue`.

- [ ] **Step 5: Push (only if user has explicitly asked)**

Skip this step unless the user requests a push. The plan stops at "feature branch ready for review."

---

## Self-Review

**1. Spec coverage:**

| Spec section | Plan task |
|--------------|-----------|
| §3.1 Status column → switch | Task 2 |
| §3.2 Remove action column buttons | Task 3 |
| §3.3 `handleTogglePromotionStatus` function | Task 1 |
| §3.4 Status semantics | Task 2 (rendered) + Task 1 (handler logic) |
| §4 Edge cases (rollback, ended, loading) | Task 1 (try/catch/finally), Task 2 (v-if on ended) |
| §5 i18n keys | All three tasks reuse existing keys; no task adds new ones |
| §6 Acceptance criteria | Task 4 step 3 manual verification |
| §7 File list | Tasks 1-3 touch only `views/promotions/index.vue` |

All spec items have a corresponding task. ✓

**2. Placeholder scan:** No "TBD" / "TODO" / "implement later" patterns. No "add appropriate error handling" — try/catch is shown in full. No "similar to Task N" — code is repeated in each task for the engineer's benefit.

**3. Type consistency:** `handleTogglePromotionStatus` is the only handler introduced; referenced consistently in Task 1 (definition) and Task 2 (template binding). `promotionToggleLoading` likewise. `Promotion` type already imported at the top of the script. ✓
