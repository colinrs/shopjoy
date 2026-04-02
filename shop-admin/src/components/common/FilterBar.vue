<template>
  <el-card
    class="filter-card"
    shadow="never"
  >
    <div class="filter-bar">
      <div class="filter-left">
        <slot name="filters" />
      </div>
      <div class="filter-right">
        <slot name="actions" />
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
// Common filter bar component for list pages
</script>

<style scoped>
.filter-card {
  margin-bottom: 16px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

[data-theme="dark"] .filter-card {
  background: var(--color-bg-card);
  border-color: var(--color-border);
}

.filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.filter-left {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  align-items: center;
}

.filter-right {
  display: flex;
  gap: 12px;
}

/* Common filter element styles */
.filter-left :deep(.search-input) {
  width: 260px;
}

.filter-left :deep(.search-input .el-input__wrapper) {
  border-radius: 12px;
}

.filter-left :deep(.filter-select) {
  width: 140px;
}

.filter-left :deep(.filter-select .el-input__wrapper) {
  border-radius: 12px;
}

/* Responsive */
@media (max-width: 768px) {
  .filter-bar {
    flex-direction: column;
    align-items: stretch;
  }

  .filter-left {
    flex-direction: column;
  }

  .filter-left :deep(.search-input),
  .filter-left :deep(.filter-select) {
    width: 100%;
  }

  .filter-right {
    justify-content: flex-end;
  }
}
</style>