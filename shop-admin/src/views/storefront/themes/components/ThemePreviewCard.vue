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