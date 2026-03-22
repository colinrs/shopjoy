<template>
  <el-card class="header-card" shadow="never">
    <div class="header-bar">
      <div class="header-left">
        <h2 class="page-title">{{ title }}</h2>
        <p v-if="description" class="page-desc">{{ description }}</p>
      </div>
      <div class="header-right">
        <slot name="actions"></slot>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  description?: string
}>()
</script>

<style scoped>
.header-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

[data-theme="dark"] .header-card {
  background: var(--color-bg-card);
  border-color: var(--color-border);
}

.header-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  flex: 1;
}

.page-title {
  margin: 0 0 8px 0;
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-primary, #1E1B4B);
}

.page-desc {
  margin: 0;
  font-size: 14px;
  color: var(--color-text-muted, #6B7280);
}

.header-right {
  display: flex;
  gap: 12px;
}

/* Responsive */
@media (max-width: 768px) {
  .header-bar {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .header-right {
    width: 100%;
  }

  .header-right :deep(.el-button) {
    width: 100%;
  }
}
</style>