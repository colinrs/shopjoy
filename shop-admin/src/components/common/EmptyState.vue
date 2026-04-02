<template>
  <div class="empty-state">
    <div class="empty-icon">
      <el-icon :size="64">
        <component :is="icon" />
      </el-icon>
    </div>
    <h3 class="empty-title">
      {{ displayTitle }}
    </h3>
    <p
      v-if="description"
      class="empty-description"
    >
      {{ description }}
    </p>
    <div
      v-if="$slots.action"
      class="empty-action"
    >
      <slot name="action" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { FolderOpened } from '@element-plus/icons-vue'
import { computed, type Component } from 'vue'
import { t } from '@/plugins/i18n'

const props = withDefaults(defineProps<{
  title?: string
  description?: string
  icon?: Component
}>(), {
  title: '',
  description: '',
  icon: FolderOpened
})

const displayTitle = computed(() => props.title || t('common.noData'))
</script>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
}

.empty-icon {
  color: var(--color-text-light, #9CA3AF);
  margin-bottom: 20px;
  opacity: 0.6;
}

[data-theme="dark"] .empty-icon {
  color: var(--color-text-muted, #6B7280);
}

.empty-title {
  font-size: 16px;
  font-weight: 500;
  color: var(--color-text-secondary, #4B5563);
  margin: 0 0 8px 0;
}

[data-theme="dark"] .empty-title {
  color: var(--color-text-secondary, #C7D2FE);
}

.empty-description {
  font-size: 14px;
  color: var(--color-text-muted, #6B7280);
  margin: 0 0 20px 0;
  max-width: 300px;
}

[data-theme="dark"] .empty-description {
  color: var(--color-text-muted, #A5B4FC);
}

.empty-action {
  margin-top: 8px;
}
</style>