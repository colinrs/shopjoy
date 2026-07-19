# Storefront Theme Preview — Design Spec

**Date:** 2026-07-19
**Scope:** `shop-admin/src/views/storefront/themes/` only
**Status:** Approved (brainstorming) — pending spec review

## Problem

The merchant-facing Themes page (`/storefront/themes`) currently renders every
theme card as an `<el-image :src="theme.thumbnail">`. The seeded data points to
non-existent CDN URLs (`https://cdn.shopjoy.com/themes/*-thumb.png`), so every
card falls back to a generic gradient + Picture icon placeholder. Merchants
cannot tell the themes apart visually — defeating the purpose of the page.

The theme's actual visual identity is defined in
`theme.default_config.{primary_color, secondary_color, font_family,
button_style}` and `theme.description`. We render from that data instead of
trying to fetch screenshots that don't exist.

## Goals

- Each theme card on `/storefront/themes` shows a distinct, themed mini
  storefront rendered purely from `ThemeItem.default_config`.
- No backend changes, no API changes, no new dependencies.
- The "look" lives in one reusable component so other surfaces (current-theme
  card, preview dialog) can adopt it later without copy-paste.

## Non-Goals

- Replacing the preview dialog's image (the large dialog at the bottom of
  `index.vue`). Dialog keeps existing behavior in this iteration.
- Real storefront rendering, iframe embedding, or actual screenshot assets.
- Live-update on current-config edits — this preview is for *selecting* a
  theme, not editing it.

## Solution

### New component

**File:** `shop-admin/src/views/storefront/themes/components/ThemePreviewCard.vue`

**Props:**
```ts
defineProps<{
  theme: ThemeItem
}>()
```
`ThemeItem` is imported from `@/api/storefront` (the existing API type).

**Renders** a fixed-shape mini storefront, ~280×160 px, themed via CSS variables
on the root `<div>`. Same shape for every theme — only colors / font / button
shape differ.

### Mock anatomy

```
┌─────────────────────────────────────┐
│ [● ● ●]   Shop          🔍  🛒      │  ← header bar (secondary bg)
├─────────────────────────────────────┤
│   Featured Collection               │  ← primary-color title
│   ─────                             │
│   ┌────┐ ┌────┐ ┌────┐              │  ← 3 product tiles
│   │ □  │ │ □  │ │ □  │              │     (neutral gray fill)
│   └────┘ └────┘ └────┘              │
│   $29  $39  $49                     │
│         [  Shop Now  ]              │  ← CTA, shape per button_style
└─────────────────────────────────────┘
```

Visual chrome:
- Header bar: `background: var(--theme-secondary)`, light text.
- "Featured Collection" title: `color: var(--theme-primary)`, font family from
  `--theme-font`.
- Product tiles: gray rectangles with a small price label below each.
- CTA button: `background: var(--theme-primary)`, color white, radius from
  `--theme-radius`. If `button_style === 'underline'`, render as a text link
  with `text-decoration: underline` instead of a filled button.

### CSS variables — mapping

The component's root `<div>` sets these via inline `:style` from
`theme.default_config`:

| CSS var              | Source                              | Fallback     |
| -------------------- | ----------------------------------- | ------------ |
| `--theme-primary`    | `default_config.primary_color`      | `#3B82F6`    |
| `--theme-secondary`  | `default_config.secondary_color`    | `#1E40AF`    |
| `--theme-font`       | `FONT_MAP[default_config.font_family]` | `Inter, system-ui, sans-serif` |
| `--theme-radius`     | mapped from `default_config.button_style` | `8px`  |
| `--theme-btn-style`  | raw `default_config.button_style` (string) | `"rounded"` |

**Radius / button-style map:**
```ts
const BUTTON_MAP: Record<string, string> = {
  rounded: '8px',
  pill:    '999px',
  square:  '0',
  underline: '0',  // CTA rendered as text — see template
}
```

**Font map** (kept inside the component, not in API data):
```ts
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
```
Unknown `font_family` values fall back to `Inter, system-ui, sans-serif`.

### Integration in `themes/index.vue`

Replace the `<el-image :src="theme.thumbnail">` block (and its `#error`
template) inside `.theme-thumbnail` with:

```vue
<div class="theme-thumbnail">
  <ThemePreviewCard :theme="theme" />
  <div class="theme-overlay">...</div>           <!-- unchanged -->
  <div v-if="theme.is_current" class="current-badge">...</div>
  <div v-if="theme.is_preset"  class="preset-badge">...</div>
</div>
```

`.theme-thumbnail` keeps its `height: 160px; overflow: hidden;` — the
component fills the container 100%.

Add the import alongside the existing ones:
```ts
import ThemePreviewCard from './components/ThemePreviewCard.vue'
```

The dead `<template #error>` block inside the card's `el-image` is removed.
The `image-placeholder` gradient stays in the preview dialog (out of scope).

### Files changed

| File | Action | Approx. size |
|------|--------|-------------:|
| `shop-admin/src/views/storefront/themes/components/ThemePreviewCard.vue` | Create | ~150 lines |
| `shop-admin/src/views/storefront/themes/index.vue`                  | Edit   | -30 / +5 lines |

No backend, no API, no dependency changes.

## Verification

- [ ] `cd shop-admin && pnpm build` passes
- [ ] `/storefront/themes` shows 5 visually distinct cards (blue classic,
      emerald modern, black minimal, violet bold, emerald nature)
- [ ] `current-badge`, `preset-badge`, and hover-overlay still render
      correctly on top of the new mock
- [ ] No console errors from missing image `src`
- [ ] The gradient placeholder no longer appears on any card
- [ ] Custom themes (no `default_config` or unknown fields) render with
      fallback values, not a broken layout

## Open risks / future

- If we later want the preview dialog to also use this component, the
  extracted component makes that a one-line swap.
- If we want live-updating preview when editing `configForm`, the CSS-variable
  approach makes that easy — pass a `configOverride` prop and bind its values
  ahead of `default_config`.
- Real screenshot assets remain a future option, not replaced by this work.