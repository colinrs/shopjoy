# Storefront Theme Preview Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace placeholder thumbnails on `/storefront/themes` with a config-driven CSS mini-mock for each theme card.

**Architecture:** A new `ThemePreviewCard.vue` component renders a fixed-shape mini storefront themed by CSS variables (`--theme-primary`, `--theme-secondary`, `--theme-font`, `--theme-radius`). The component reads `theme.default_config`. The backend `ListThemes` currently drops `DefaultConfig` from its DTO — fix that, extend the frontend `ThemeItem` type, then drop the component into the card loop.

**Tech Stack:** Vue 3 (`<script setup>`) · TypeScript · Element Plus · Go (GORM) · go-zero API definition

## Global Constraints

- Backend project rules: `pkg/code/code.go` for errors, `decimal.Decimal` for money, `time.Time` for timestamps (already satisfied — no monetary/time fields added).
- API definitions live in `admin/desc/storefront.api`; run `cd admin && make api` after editing.
- Frontend API types live in `shop-admin/src/api/storefront.ts`; keep in sync after backend changes.
- Build verification: `cd shop-admin && pnpm build` (frontend), `cd admin && make build` (backend).
- Component uses `scoped` style only. No global CSS additions.
- Plan file: `docs/superpowers/plans/2026-07-19-storefront-theme-preview.md`
- Spec file: `docs/superpowers/specs/2026-07-19-storefront-theme-preview-design.md`

---

## File Structure

| File | Responsibility | Action |
|------|----------------|--------|
| `admin/internal/domain/storefront/entity.go` | Change `Theme.DefaultConfig` type to `json.RawMessage` | Edit |
| `admin/internal/infrastructure/persistence/theme_repository.go` | Stop pre-parsing `DefaultConfig` into struct; pass bytes through | Edit |
| `admin/internal/application/storefront/service.go` | Add `defaultConfigToDTO` helper; populate `DefaultConfig` in `ListThemes`; fix existing `GetCurrentTheme` caller | Edit |
| `admin/desc/storefront.api` | (No change — DTO already has `default_config` field) | — |
| `admin/internal/types/types.go` | (No change — fields already generated) | — |
| `shop-admin/src/api/storefront.ts` | Add `default_config?: ThemeConfigDTO` to `ThemeItem` | Edit |
| `shop-admin/src/views/storefront/themes/components/ThemePreviewCard.vue` | New: themed mini-mock renderer | Create |
| `shop-admin/src/views/storefront/themes/index.vue` | Swap `<el-image>` for new component in card loop | Edit |

---

## Task 1: Backend — populate `DefaultConfig` in `ListThemes`

**Files:**
- Modify: `admin/internal/domain/storefront/entity.go:40-54` (change `DefaultConfig` field type)
- Modify: `admin/internal/infrastructure/persistence/theme_repository.go:54-69` (skip pre-parsing for `DefaultConfig`)
- Modify: `admin/internal/application/storefront/service.go:255-258` (update existing caller to use new helper)
- Modify: `admin/internal/application/storefront/service.go:187-213` (`ListThemes`)
- Modify: `admin/internal/application/storefront/service.go:424` area (add helper)

**Interfaces:**
- Consumes: `*storefront.Theme` (entity); `Theme.DefaultConfig` will become `json.RawMessage` carrying the original column bytes
- Produces: `*ThemeDTO.DefaultConfig *ThemeConfigDTO` populated with `PrimaryColor`, `SecondaryColor`, `FontFamily`, `ButtonStyle`

**Context the implementer needs:**
The seed SQL (`sql/storefront/schema.sql:234-242`) inserts `default_config` as **flat DTO format**:
```json
{"primary_color":"#3B82F6","secondary_color":"#1E40AF","font_family":"inter","button_style":"rounded"}
```
But `ThemeConfig` is deserialized into nested fields (`Colors["primary"]`, `Fonts["heading"]`, `Components["button_style"]`). Once GORM unmarshals flat JSON into `ThemeConfig`, the flat keys are **lost** — the struct only knows `Colors/Fonts/Layout/Components`. So `t.DefaultConfig.Colors["primary"]` returns `""` for seed data.

The seed also inserts `config` in nested format which is what the existing helper expects — that's why `GetCurrentTheme` works when `shop.ThemeConfig` was previously set via `UpdateThemeConfig`.

The fix: keep `DefaultConfig` as **raw bytes** (`json.RawMessage`) so we can attempt both parse orders at the service layer.

- [ ] **Step 1: Change the entity field type**

In `admin/internal/domain/storefront/entity.go`, add `"encoding/json"` to the imports (line 3 currently imports only `context` and `time`). Then change the `DefaultConfig` field on `Theme` (line 50):

```go
DefaultConfig json.RawMessage `gorm:"column:default_config;type:text"`
```

(The `Config` and `ConfigSchema` fields stay as their current struct types — only `DefaultConfig` needs raw bytes for the dual-format case.)

- [ ] **Step 2: Stop pre-parsing `DefaultConfig` in the repository**

In `admin/internal/infrastructure/persistence/theme_repository.go`, replace lines 54-57:

```go
var defaultConfig storefront.ThemeConfig
if m.DefaultConfig != "" {
    _ = json.Unmarshal([]byte(m.DefaultConfig), &defaultConfig)
}
```

with:

```go
var defaultConfig json.RawMessage
if m.DefaultConfig != "" {
    defaultConfig = json.RawMessage(m.DefaultConfig)
}
```

And update the assignment on line 69 to use `defaultConfig` (now `json.RawMessage`) directly — it already matches the new entity type, no code change there.

Also update the `fromThemeEntity` function (line 79): the existing code does `defaultConfig, _ := json.Marshal(t.DefaultConfig)`. With `RawMessage`, you must use `[]byte(t.DefaultConfig)` directly to avoid the double-marshal (`RawMessage` is already bytes; `Marshal` on it is a no-op but emits the raw JSON verbatim, which still works but is wasteful). Change line 79 to:

```go
defaultConfig := []byte(t.DefaultConfig)
```

Then line 104 stays the same: `DefaultConfig: string(defaultConfig),`.

- [ ] **Step 3: Add the tolerant helper**

In `admin/internal/application/storefront/service.go`, insert after the existing `themeConfigToDTO` function (after line 455):

```go
// defaultConfigToDTO parses theme.DefaultConfig JSON which may be in either
// nested entity format ({"colors":{"primary":"#xxx"},"fonts":{"heading":"..."},"components":{"button_style":"..."}})
// or flat DTO format ({"primary_color":"#xxx","font_family":"...","button_style":"..."}).
// The latter is what the seed SQL inserts. Returns nil if neither format yields data.
func defaultConfigToDTO(raw json.RawMessage) *ThemeConfigDTO {
	if len(raw) == 0 {
		return nil
	}

	// Try nested form first
	var nested storefront.ThemeConfig
	if err := json.Unmarshal(raw, &nested); err == nil {
		if dto := themeConfigToDTO(&nested); dto != nil &&
			(dto.PrimaryColor != "" || dto.SecondaryColor != "" ||
				dto.FontFamily != "" || dto.ButtonStyle != "") {
			return dto
		}
	}

	// Fall back to flat form
	var flat ThemeConfigDTO
	if err := json.Unmarshal(raw, &flat); err != nil {
		return nil
	}
	if flat.PrimaryColor == "" && flat.SecondaryColor == "" &&
		flat.FontFamily == "" && flat.ButtonStyle == "" {
		return nil
	}
	return &flat
}
```

- [ ] **Step 4: Update the existing caller in `GetCurrentTheme`**

In `admin/internal/application/storefront/service.go` line 257, change:

```go
config = themeConfigToDTO(&theme.DefaultConfig)
```

to:

```go
config = defaultConfigToDTO(theme.DefaultConfig)
```

(`theme.DefaultConfig` is now `json.RawMessage`; `defaultConfigToDTO` accepts that type.)

- [ ] **Step 5: Use the helper in `ListThemes`**

Edit `admin/internal/application/storefront/service.go` lines 199-211. Replace the DTO construction with:

```go
dtos := make([]*ThemeDTO, len(themes))
for i, t := range themes {
	dtos[i] = &ThemeDTO{
		ID:            t.ID,
		Code:          t.Code,
		Name:          t.Name,
		Description:   t.Description,
		PreviewImage:  t.PreviewImage,
		Thumbnail:     t.Thumbnail,
		IsPreset:      t.IsPreset,
		IsCurrent:     t.ID == currentThemeID,
		DefaultConfig: defaultConfigToDTO(t.DefaultConfig),
	}
}
```

- [ ] **Step 6: Build backend**

Run: `cd admin && make build`
Expected: build succeeds. (Behavior is unchanged for callers that don't read `DefaultConfig`.)

- [ ] **Step 7: Commit**

```bash
git add admin/internal/domain/storefront/entity.go \
        admin/internal/infrastructure/persistence/theme_repository.go \
        admin/internal/application/storefront/service.go
git commit -m "feat(storefront): populate DefaultConfig in ListThemes response

Change Theme.DefaultConfig to json.RawMessage so ListThemes can attempt
both nested entity format and flat DTO format (the seed uses flat keys).
Add defaultConfigToDTO helper that tries nested first, falls back to flat.
Frontend theme card preview depends on this to show distinct colors per theme."
```

---

## Task 2: Frontend — extend `ThemeItem` type

**Files:**
- Modify: `shop-admin/src/api/storefront.ts:5-14` (`ThemeItem` interface)

**Interfaces:**
- Consumes: `ThemeConfigDTO` (already defined at `storefront.ts:20-25`)
- Produces: `ThemeItem.default_config?: ThemeConfigDTO`

- [ ] **Step 1: Add the field**

In `shop-admin/src/api/storefront.ts`, replace the `ThemeItem` interface:

```ts
export interface ThemeItem {
  id: string
  code: string
  name: string
  description: string
  preview_image: string
  thumbnail: string
  is_preset: boolean
  is_current: boolean
  default_config?: ThemeConfigDTO
}
```

- [ ] **Step 2: Verify TypeScript compiles**

Run: `cd shop-admin && pnpm build`
Expected: build succeeds. No new errors. (The field is optional, so existing call sites are unaffected.)

- [ ] **Step 3: Commit**

```bash
git add shop-admin/src/api/storefront.ts
git commit -m "feat(storefront-ui): expose default_config on ThemeItem type

Mirrors the backend ThemeDTO field so the theme card preview can read
per-theme colors, fonts, and button style."
```

---

## Task 3: Frontend — create `ThemePreviewCard.vue`

**Files:**
- Create: `shop-admin/src/views/storefront/themes/components/ThemePreviewCard.vue`

**Interfaces:**
- Consumes: `theme: ThemeItem` (from `@/api/storefront`)
- Produces: a 100% × 100% sized mini-mock themed by inline `--theme-*` CSS variables

- [ ] **Step 1: Create the components directory and file**

Create `shop-admin/src/views/storefront/themes/components/ThemePreviewCard.vue` with this exact content:

```vue
<template>
  <div
    class="theme-preview-card"
    :style="rootStyle"
  >
    <!-- Header bar -->
    <div class="tpc-header">
      <span class="tpc-dots">
        <span class="tpc-dot" />
        <span class="tpc-dot" />
        <span class="tpc-dot" />
      </span>
      <span class="tpc-shop-name">Shop</span>
      <span class="tpc-icons">🔍 🛒</span>
    </div>

    <!-- Body -->
    <div class="tpc-body">
      <div class="tpc-section-title">
        Featured Collection
      </div>
      <div class="tpc-section-rule" />

      <div class="tpc-grid">
        <div
          v-for="i in 3"
          :key="i"
          class="tpc-tile"
        >
          <div class="tpc-tile-img" />
          <div class="tpc-tile-price">
            ${{ 20 + i * 10 }}
          </div>
        </div>
      </div>

      <div class="tpc-cta-wrap">
        <component
          :is="buttonStyle === 'underline' ? 'span' : 'button'"
          class="tpc-cta"
          :class="{ 'tpc-cta-link': buttonStyle === 'underline' }"
          type="button"
        >
          Shop Now
        </component>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { ThemeItem } from '@/api/storefront'

const props = defineProps<{
  theme: ThemeItem
}>()

const FONT_MAP: Record<string, string> = {
  inter:        'Inter, system-ui, sans-serif',
  roboto:       'Roboto, system-ui, sans-serif',
  opensans:     '"Open Sans", system-ui, sans-serif',
  poppins:      'Poppins, system-ui, sans-serif',
  montserrat:   'Montserrat, system-ui, sans-serif',
  helvetica:    '"Helvetica Neue", Helvetica, Arial, sans-serif',
  dmsans:       '"DM Sans", system-ui, sans-serif',
  nunito:       'Nunito, system-ui, sans-serif',
  merriweather: 'Merriweather, Georgia, serif',
  lora:         'Lora, Georgia, serif',
  notosans:     '"Noto Sans", system-ui, sans-serif',
}

const BUTTON_RADIUS: Record<string, string> = {
  rounded:   '8px',
  pill:      '999px',
  square:    '0',
  underline: '0',
}

const cfg = computed(() => props.theme.default_config ?? null)

const rootStyle = computed(() => ({
  '--theme-primary':   cfg.value?.primary_color   || '#3B82F6',
  '--theme-secondary': cfg.value?.secondary_color || '#1E40AF',
  '--theme-font':      FONT_MAP[cfg.value?.font_family || ''] || 'Inter, system-ui, sans-serif',
  '--theme-radius':    BUTTON_RADIUS[cfg.value?.button_style || ''] || '8px',
}))

const buttonStyle = computed(() => cfg.value?.button_style || 'rounded')
</script>

<style scoped>
.theme-preview-card {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  font-family: var(--theme-font);
  background: #fff;
  overflow: hidden;
}

.tpc-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 10px;
  background: var(--theme-secondary);
  color: #fff;
  font-size: 10px;
  flex-shrink: 0;
}

.tpc-dots {
  display: flex;
  gap: 3px;
}

.tpc-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.55);
}

.tpc-shop-name {
  flex: 1;
  font-weight: 600;
}

.tpc-icons {
  font-size: 9px;
  letter-spacing: 2px;
}

.tpc-body {
  flex: 1;
  padding: 8px 10px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  background: #fff;
}

.tpc-section-title {
  font-size: 11px;
  font-weight: 700;
  color: var(--theme-primary);
}

.tpc-section-rule {
  width: 24px;
  height: 2px;
  background: var(--theme-primary);
  border-radius: 1px;
}

.tpc-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 6px;
  flex: 1;
}

.tpc-tile {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.tpc-tile-img {
  flex: 1;
  background: linear-gradient(135deg, #F3F4F6 0%, #E5E7EB 100%);
  border-radius: 3px;
}

.tpc-tile-price {
  font-size: 9px;
  color: #6B7280;
  font-weight: 600;
}

.tpc-cta-wrap {
  display: flex;
  justify-content: center;
  margin-top: auto;
  padding-top: 6px;
}

.tpc-cta {
  background: var(--theme-primary);
  color: #fff;
  border: none;
  padding: 4px 14px;
  border-radius: var(--theme-radius);
  font-family: inherit;
  font-size: 10px;
  font-weight: 600;
  cursor: default;
}

.tpc-cta-link {
  background: transparent;
  color: var(--theme-primary);
  text-decoration: underline;
  padding: 4px 6px;
}
</style>
```

- [ ] **Step 2: Verify TypeScript compiles**

Run: `cd shop-admin && pnpm build`
Expected: build succeeds.

- [ ] **Step 3: Commit**

```bash
git add shop-admin/src/views/storefront/themes/components/ThemePreviewCard.vue
git commit -m "feat(storefront-ui): add ThemePreviewCard mock component

CSS-variable driven mini storefront. Header bar tinted by secondary,
section title and CTA by primary, font from FONT_MAP, button shape from
button_style (rounded/pill/square/underline). Falls back to defaults
when theme.default_config is missing."
```

---

## Task 4: Frontend — wire the component into the themes page

**Files:**
- Modify: `shop-admin/src/views/storefront/themes/index.vue`

**Interfaces:**
- Consumes: `ThemePreviewCard` from Task 3; `ThemeItem[]` already in `themes` ref (line 334)
- Produces: each theme card's `.theme-thumbnail` renders `<ThemePreviewCard>` instead of `<el-image>`

- [ ] **Step 1: Add the import**

In `shop-admin/src/views/storefront/themes/index.vue`, after the existing `import PageHeader from '@/components/common/PageHeader.vue'` line (~316), add:

```ts
import ThemePreviewCard from './components/ThemePreviewCard.vue'
```

- [ ] **Step 2: Replace the `<el-image>` block in the card loop**

In the same file, replace the `.theme-thumbnail` div content (current lines 167-203) with:

```vue
<div class="theme-thumbnail">
  <ThemePreviewCard :theme="theme" />
  <div class="theme-overlay">
    <el-button
      type="primary"
      size="small"
      @click.stop="previewTheme(theme)"
    >
      {{ $t('storefront.preview') }}
    </el-button>
  </div>
  <div
    v-if="theme.is_current"
    class="current-badge"
  >
    <el-icon><Check /></el-icon>
    {{ $t('storefront.inUse') }}
  </div>
  <div
    v-if="theme.is_preset"
    class="preset-badge"
  >
    {{ $t('storefront.officialTheme') }}
  </div>
</div>
```

What changed: `<el-image :src="theme.thumbnail">` and its `<template #error>` block are gone; `<ThemePreviewCard :theme="theme" />` is in their place. The overlay, current-badge, and preset-badge are unchanged.

- [ ] **Step 3: Verify TypeScript compiles**

Run: `cd shop-admin && pnpm build`
Expected: build succeeds.

- [ ] **Step 4: Commit**

```bash
git add shop-admin/src/views/storefront/themes/index.vue
git commit -m "feat(storefront-ui): use ThemePreviewCard in theme cards

Replaces non-functional <el-image> + gradient placeholder with the
config-driven mock. Drop the dead #error template in the card loop."
```

---

## Task 5: End-to-end verification

- [ ] **Step 1: Backend build clean**

Run: `cd admin && make build`
Expected: no errors.

- [ ] **Step 2: Frontend build clean**

Run: `cd shop-admin && pnpm build`
Expected: no errors.

- [ ] **Step 3: Manual visual check**

In a browser, open `http://localhost:3000/storefront/themes`. Verify the 5 preset themes render distinct, themed previews:
- **Classic** → blue header (deep blue), blue section title, rounded blue CTA, Inter font
- **Modern** → emerald header (deep emerald), emerald section title, pill emerald CTA, Poppins font
- **Minimal** → black header, gray section title, **underlined** text-link CTA, Helvetica font
- **Bold** → violet header (deep violet), violet section title, rounded violet CTA, DM Sans font
- **Nature** → emerald header (deep emerald), emerald section title, rounded emerald CTA, Merriweather font

- [ ] **Step 4: Layout integrity**

In the same browser view, verify:
- `current-badge` ("在用") still appears on the active theme card
- `preset-badge` ("官方主题") still appears on preset cards
- Hover overlay (preview button) still appears on mouseover
- Card heights are unchanged (~160 px thumbnail)

- [ ] **Step 5: No console errors**

Open DevTools console. Reload. Verify no 404s for image src, no Vue warnings about missing props, no TypeErrors.

- [ ] **Step 6: Custom-theme fallback (if applicable)**

If there are custom themes (uploaded by merchant) in the DB, their `default_config` may be null. Verify those cards render with the fallback blue palette rather than breaking — `default_config?.primary_color || '#3B82F6'` handles this.

- [ ] **Step 7: Final commit (if any drift)**

If any commit was needed during verification (shouldn't be — Tasks 1-4 should leave the tree clean), commit with:

```bash
git commit --allow-empty -m "chore: verified storefront theme preview end-to-end"
```

---

## Out-of-scope follow-ups (recorded for future tasks)

- **Preview dialog**: `index.vue:260-307` still uses `<el-image :src="previewThemeData.preview_image">`. Out of scope for this iteration; could adopt `ThemePreviewCard` (with a `large` size variant) later.
- **Live config edit preview**: the "current theme" card's `configForm` could drive a live preview by passing a `configOverride` prop to `ThemePreviewCard`. CSS-variable architecture makes this trivial.
- **Real screenshot assets**: not pursued. Seeded CDN URLs stay as `preview_image`/`thumbnail` placeholders.