<template>
  <el-tag :type="tagType" :size="size" :effect="effect" class="status-tag">
    <slot>{{ displayText }}</slot>
  </el-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue'

type StatusType = 'success' | 'warning' | 'danger' | 'info' | 'primary'

const props = withDefaults(defineProps<{
  status?: string | number
  typeMap?: Record<string | number, { type: StatusType; text: string }>
  size?: 'large' | 'default' | 'small'
  effect?: 'dark' | 'light' | 'plain'
}>(), {
  size: 'small',
  effect: 'light',
  typeMap: () => ({
    // Default status mappings
    'success': { type: 'success', text: '成功' },
    'warning': { type: 'warning', text: '警告' },
    'danger': { type: 'danger', text: '危险' },
    'error': { type: 'danger', text: '错误' },
    'info': { type: 'info', text: '信息' },
    'pending': { type: 'warning', text: '待处理' },
    'processing': { type: 'primary', text: '处理中' },
    'completed': { type: 'success', text: '已完成' },
    'failed': { type: 'danger', text: '失败' },
    'active': { type: 'success', text: '启用' },
    'inactive': { type: 'info', text: '禁用' },
    'on_sale': { type: 'success', text: '在售' },
    'off_sale': { type: 'info', text: '下架' },
    'draft': { type: 'warning', text: '草稿' },
    1: { type: 'success', text: '正常' },
    0: { type: 'info', text: '禁用' },
    2: { type: 'danger', text: '异常' }
  })
})

const tagType = computed(() => {
  const mapping = props.typeMap[props.status as string | number]
  return mapping?.type || 'info'
})

const displayText = computed(() => {
  const mapping = props.typeMap[props.status as string | number]
  return mapping?.text || String(props.status)
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