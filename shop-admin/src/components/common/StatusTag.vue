<template>
  <el-tag
    :type="tagType"
    :size="size"
    :effect="effect"
    class="status-tag"
  >
    <slot>{{ displayText }}</slot>
  </el-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

type StatusType = 'success' | 'warning' | 'danger' | 'info' | 'primary'

const { t } = useI18n()

const props = withDefaults(defineProps<{
  status?: string | number
  typeMap?: Record<string | number, { type: StatusType; text: string }>
  size?: 'large' | 'default' | 'small'
  effect?: 'dark' | 'light' | 'plain'
}>(), {
  size: 'small',
  effect: 'light',
  typeMap: () => ({
    // Default status mappings with translation keys
    'success': { type: 'success', text: 'status.success' },
    'warning': { type: 'warning', text: 'status.warning' },
    'danger': { type: 'danger', text: 'status.danger' },
    'error': { type: 'danger', text: 'status.danger' },
    'info': { type: 'info', text: 'status.info' },
    'pending': { type: 'warning', text: 'status.pending' },
    'processing': { type: 'primary', text: 'status.processing' },
    'completed': { type: 'success', text: 'status.completed' },
    'failed': { type: 'danger', text: 'status.failed' },
    'active': { type: 'success', text: 'status.active' },
    'inactive': { type: 'info', text: 'status.inactive' },
    'on_sale': { type: 'success', text: 'status.onSale' },
    'off_sale': { type: 'info', text: 'status.offSale' },
    'draft': { type: 'warning', text: 'status.draft' },
    1: { type: 'success', text: 'status.normal' },
    0: { type: 'info', text: 'status.disabled' },
    2: { type: 'danger', text: 'status.abnormal' }
  })
})

const tagType = computed(() => {
  const mapping = props.typeMap[props.status as string | number]
  return mapping?.type || 'info'
})

const displayText = computed(() => {
  const mapping = props.typeMap[props.status as string | number]
  const text = mapping?.text || String(props.status)
  // Check if text looks like a translation key (contains dot notation)
  if (text.includes('.')) {
    return t(text)
  }
  return text
})
</script>

<style scoped>
.status-tag {
  font-weight: 500;
}

/* Primary tag styling */
:deep(.el-tag--primary) {
  background-color: rgba(99, 102, 241, 0.1);
  border-color: rgba(99, 102, 241, 0.2);
  color: var(--color-primary, #6366F1);
}

[data-theme="dark"] :deep(.el-tag--primary) {
  background-color: rgba(99, 102, 241, 0.2);
  border-color: rgba(99, 102, 241, 0.3);
  color: var(--color-primary-light, #818CF8);
}

/* Success tag styling */
:deep(.el-tag--success) {
  background-color: rgba(16, 185, 129, 0.1);
  border-color: rgba(16, 185, 129, 0.2);
  color: var(--color-cta, #10B981);
}

[data-theme="dark"] :deep(.el-tag--success) {
  background-color: rgba(16, 185, 129, 0.2);
  border-color: rgba(16, 185, 129, 0.3);
  color: var(--color-cta-light, #34D399);
}

/* Warning tag styling */
:deep(.el-tag--warning) {
  background-color: rgba(245, 158, 11, 0.1);
  border-color: rgba(245, 158, 11, 0.2);
  color: var(--color-warning, #F59E0B);
}

[data-theme="dark"] :deep(.el-tag--warning) {
  background-color: rgba(245, 158, 11, 0.2);
  border-color: rgba(245, 158, 11, 0.3);
  color: #FBBF24;
}

/* Danger tag styling */
:deep(.el-tag--danger) {
  background-color: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.2);
  color: var(--color-danger, #EF4444);
}

[data-theme="dark"] :deep(.el-tag--danger) {
  background-color: rgba(239, 68, 68, 0.2);
  border-color: rgba(239, 68, 68, 0.3);
  color: #F87171;
}

/* Info tag styling */
:deep(.el-tag--info) {
  background-color: rgba(107, 114, 128, 0.1);
  border-color: rgba(107, 114, 128, 0.2);
  color: var(--color-text-muted, #6B7280);
}

[data-theme="dark"] :deep(.el-tag--info) {
  background-color: rgba(107, 114, 128, 0.2);
  border-color: rgba(107, 114, 128, 0.3);
  color: var(--color-text-secondary, #9CA3AF);
}
</style>