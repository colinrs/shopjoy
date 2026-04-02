<template>
  <div class="skeleton-page">
    <!-- Header skeleton -->
    <el-card
      v-if="showHeader"
      class="skeleton-page-header"
      shadow="never"
    >
      <div class="skeleton-header-content">
        <div class="skeleton-header-left">
          <div class="skeleton-page-title" />
          <div class="skeleton-page-subtitle" />
        </div>
        <div
          v-if="showHeaderAction"
          class="skeleton-header-action"
        />
      </div>
    </el-card>

    <!-- Stats skeleton -->
    <div
      v-if="showStats"
      class="skeleton-stats-grid"
    >
      <el-card
        v-for="i in statsCount"
        :key="i"
        class="skeleton-stat-card"
        shadow="never"
      >
        <div class="skeleton-stat">
          <div class="skeleton-stat-icon" />
          <div class="skeleton-stat-info">
            <div class="skeleton-stat-value" />
            <div class="skeleton-stat-label" />
          </div>
        </div>
      </el-card>
    </div>

    <!-- Filter bar skeleton -->
    <el-card
      v-if="showFilter"
      class="skeleton-filter-card"
      shadow="never"
    >
      <div class="skeleton-filter-bar">
        <div class="skeleton-filter-left">
          <div class="skeleton-search" />
          <div class="skeleton-select" />
          <div class="skeleton-select" />
        </div>
        <div class="skeleton-filter-action" />
      </div>
    </el-card>

    <!-- Table skeleton -->
    <el-card
      v-if="showTable"
      class="skeleton-table-card"
      shadow="never"
    >
      <TableSkeleton
        :rows="tableRows"
        :columns="tableColumns"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import TableSkeleton from './TableSkeleton.vue'

defineProps({
  showHeader: {
    type: Boolean,
    default: true
  },
  showHeaderAction: {
    type: Boolean,
    default: true
  },
  showStats: {
    type: Boolean,
    default: false
  },
  statsCount: {
    type: Number,
    default: 4
  },
  showFilter: {
    type: Boolean,
    default: true
  },
  showTable: {
    type: Boolean,
    default: true
  },
  tableRows: {
    type: Number,
    default: 5
  },
  tableColumns: {
    type: Number,
    default: 6
  }
})
</script>

<style scoped>
.skeleton-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.skeleton-page-header {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.skeleton-header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.skeleton-header-left {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.skeleton-page-title {
  width: 160px;
  height: 24px;
  background: linear-gradient(90deg, #E5E7EB 25%, #F3F4F6 50%, #E5E7EB 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 6px;
}

[data-theme="dark"] .skeleton-page-title {
  background: linear-gradient(90deg, #312E81 25%, #3730A3 50%, #312E81 75%);
  background-size: 200% 100%;
}

.skeleton-page-subtitle {
  width: 240px;
  height: 14px;
  background: linear-gradient(90deg, #E5E7EB 25%, #F3F4F6 50%, #E5E7EB 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 4px;
}

[data-theme="dark"] .skeleton-page-subtitle {
  background: linear-gradient(90deg, #312E81 25%, #3730A3 50%, #312E81 75%);
  background-size: 200% 100%;
}

.skeleton-header-action {
  width: 100px;
  height: 36px;
  background: linear-gradient(90deg, #E5E7EB 25%, #F3F4F6 50%, #E5E7EB 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 18px;
}

[data-theme="dark"] .skeleton-header-action {
  background: linear-gradient(90deg, #312E81 25%, #3730A3 50%, #312E81 75%);
  background-size: 200% 100%;
}

/* Stats Grid */
.skeleton-stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
}

.skeleton-stat-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.skeleton-stat {
  display: flex;
  align-items: center;
  gap: 16px;
}

.skeleton-stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  background: linear-gradient(90deg, #E5E7EB 25%, #F3F4F6 50%, #E5E7EB 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}

[data-theme="dark"] .skeleton-stat-icon {
  background: linear-gradient(90deg, #312E81 25%, #3730A3 50%, #312E81 75%);
  background-size: 200% 100%;
}

.skeleton-stat-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.skeleton-stat-value {
  width: 80px;
  height: 28px;
  background: linear-gradient(90deg, #E5E7EB 25%, #F3F4F6 50%, #E5E7EB 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 6px;
}

[data-theme="dark"] .skeleton-stat-value {
  background: linear-gradient(90deg, #312E81 25%, #3730A3 50%, #312E81 75%);
  background-size: 200% 100%;
}

.skeleton-stat-label {
  width: 60px;
  height: 14px;
  background: linear-gradient(90deg, #E5E7EB 25%, #F3F4F6 50%, #E5E7EB 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 4px;
}

[data-theme="dark"] .skeleton-stat-label {
  background: linear-gradient(90deg, #312E81 25%, #3730A3 50%, #312E81 75%);
  background-size: 200% 100%;
}

/* Filter Bar */
.skeleton-filter-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.skeleton-filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.skeleton-filter-left {
  display: flex;
  gap: 12px;
}

.skeleton-search {
  width: 260px;
  height: 36px;
  background: linear-gradient(90deg, #E5E7EB 25%, #F3F4F6 50%, #E5E7EB 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 18px;
}

[data-theme="dark"] .skeleton-search {
  background: linear-gradient(90deg, #312E81 25%, #3730A3 50%, #312E81 75%);
  background-size: 200% 100%;
}

.skeleton-select {
  width: 140px;
  height: 36px;
  background: linear-gradient(90deg, #E5E7EB 25%, #F3F4F6 50%, #E5E7EB 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 8px;
}

[data-theme="dark"] .skeleton-select {
  background: linear-gradient(90deg, #312E81 25%, #3730A3 50%, #312E81 75%);
  background-size: 200% 100%;
}

.skeleton-filter-action {
  width: 100px;
  height: 36px;
  background: linear-gradient(90deg, #E5E7EB 25%, #F3F4F6 50%, #E5E7EB 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 18px;
}

[data-theme="dark"] .skeleton-filter-action {
  background: linear-gradient(90deg, #312E81 25%, #3730A3 50%, #312E81 75%);
  background-size: 200% 100%;
}

/* Table Card */
.skeleton-table-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

@keyframes shimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}

/* Responsive */
@media (max-width: 768px) {
  .skeleton-header-content {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .skeleton-filter-bar {
    flex-direction: column;
    gap: 12px;
  }

  .skeleton-filter-left {
    flex-wrap: wrap;
  }

  .skeleton-search {
    width: 100%;
  }
}
</style>