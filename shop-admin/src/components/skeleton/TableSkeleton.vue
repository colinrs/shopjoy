<template>
  <div class="skeleton-table">
    <div class="skeleton-table-header">
      <div
        v-for="i in columns"
        :key="i"
        class="skeleton-cell skeleton-header"
      />
    </div>
    <div
      v-for="row in rows"
      :key="row"
      class="skeleton-table-row"
    >
      <div
        v-for="col in columns"
        :key="col"
        class="skeleton-cell"
        :class="getCellClass(col)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps({
  rows: {
    type: Number,
    default: 5
  },
  columns: {
    type: Number,
    default: 6
  }
})

const getCellClass = (col: number) => {
  // First column is usually wider (name/title)
  if (col === 1) return 'skeleton-cell--lg'
  // Last column is usually actions (narrower)
  if (col === 6) return 'skeleton-cell--sm'
  return 'skeleton-cell--md'
}
</script>

<style scoped>
.skeleton-table {
  width: 100%;
}

.skeleton-table-header {
  display: flex;
  gap: 16px;
  padding: 16px;
  background: var(--color-bg-base, #F5F3FF);
  border-radius: 12px 12px 0 0;
}

.skeleton-table-row {
  display: flex;
  gap: 16px;
  padding: 16px;
  border-bottom: 1px solid var(--color-border-light, #F3F4F6);
}

.skeleton-table-row:last-child {
  border-bottom: none;
}

.skeleton-cell {
  background: linear-gradient(90deg, #E5E7EB 25%, #F3F4F6 50%, #E5E7EB 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 6px;
  height: 20px;
}

[data-theme="dark"] .skeleton-cell {
  background: linear-gradient(90deg, #312E81 25%, #3730A3 50%, #312E81 75%);
  background-size: 200% 100%;
}

.skeleton-header {
  height: 16px;
  width: 80px;
}

.skeleton-cell--sm {
  width: 60px;
  height: 32px;
  border-radius: 16px;
}

.skeleton-cell--md {
  flex: 1;
  min-width: 80px;
}

.skeleton-cell--lg {
  flex: 2;
  min-width: 180px;
}

@keyframes shimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}
</style>