<template>
  <div
    class="stat-card"
    :class="[iconColor, { clickable: clickable }]"
    @click="handleClick"
  >
    <div
      class="stat-icon"
      :class="iconColor"
    >
      <el-icon :size="24">
        <component :is="icon" />
      </el-icon>
    </div>
    <div class="stat-info">
      <p class="stat-label">
        {{ title }}
      </p>
      <p class="stat-value">
        {{ formatValue(value) }}
      </p>
      <div
        v-if="trend"
        class="stat-trend"
        :class="{ up: trend.value >= 0, down: trend.value < 0 }"
      >
        <el-icon v-if="trend.value >= 0">
          <CaretTop />
        </el-icon>
        <el-icon v-else>
          <CaretBottom />
        </el-icon>
        <span>{{ Math.abs(trend.value) }}%</span>
        <span class="trend-label">{{ trend.label }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { CaretTop, CaretBottom } from '@element-plus/icons-vue'
import type { Component } from 'vue'

interface Trend {
  value: number
  label: string
}

const props = withDefaults(defineProps<{
  title: string
  value: number | string
  icon: Component
  iconColor?: 'primary' | 'success' | 'warning' | 'danger'
  trend?: Trend
  clickable?: boolean
}>(), {
  iconColor: 'primary',
  clickable: false
})

const emit = defineEmits<{
  click: []
}>()

const formatValue = (val: number | string) => {
  if (typeof val === 'string') return val
  return val.toLocaleString('en-US')
}

const handleClick = () => {
  if (props.clickable) {
    emit('click')
  }
}
</script>

<style scoped>
.stat-card {
  background: #fff;
  border-radius: 16px;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  border: 1px solid rgba(99, 102, 241, 0.06);
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.stat-card.clickable {
  cursor: pointer;
}

.stat-card.clickable:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px -4px rgba(99, 102, 241, 0.12);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-icon.primary {
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  color: white;
}

.stat-icon.success {
  background: linear-gradient(135deg, #10B981 0%, #34D399 100%);
  color: white;
}

.stat-icon.warning {
  background: linear-gradient(135deg, #F59E0B 0%, #FBBF24 100%);
  color: white;
}

.stat-icon.danger {
  background: linear-gradient(135deg, #EF4444 0%, #F87171 100%);
  color: white;
}

.stat-info {
  flex: 1;
  min-width: 0;
}

.stat-label {
  font-size: 13px;
  color: #6B7280;
  margin: 0 0 4px 0;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #1E1B4B;
  margin: 0 0 4px 0;
  font-family: 'Fira Sans', sans-serif;
}

.stat-trend {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #6B7280;
}

.stat-trend.up {
  color: #10B981;
}

.stat-trend.down {
  color: #EF4444;
}

.trend-label {
  color: #9CA3AF;
  margin-left: 4px;
}

/* Border left indicator */
.stat-card.primary {
  border-left: 4px solid #6366F1;
}

.stat-card.success {
  border-left: 4px solid #10B981;
}

.stat-card.warning {
  border-left: 4px solid #F59E0B;
}

.stat-card.danger {
  border-left: 4px solid #EF4444;
}

/* Responsive */
@media (max-width: 768px) {
  .stat-card {
    padding: 20px;
  }

  .stat-value {
    font-size: 24px;
  }

  .stat-icon {
    width: 48px;
    height: 48px;
  }
}
</style>