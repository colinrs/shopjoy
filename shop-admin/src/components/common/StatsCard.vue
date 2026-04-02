<template>
  <div class="stats-card" :class="[`stats-card--${color}`, { 'stats-card--clickable': clickable }]" @click="handleClick">
    <div class="stats-icon" v-if="icon || $slots.icon">
      <slot name="icon">
        <el-icon v-if="icon" :size="32">
          <component :is="icon" />
        </el-icon>
      </slot>
    </div>
    <div class="stats-info">
      <p class="stats-value">
        <slot name="value">{{ value }}</slot>
      </p>
      <p class="stats-label">
        <slot name="label">{{ title }}</slot>
      </p>
      <p class="stats-trend" v-if="trend || $slots.trend">
        <slot name="trend">{{ trend }}</slot>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Component } from 'vue'

interface Props {
  title?: string
  value?: string | number
  icon?: Component | string
  trend?: string
  color?: 'default' | 'primary' | 'success' | 'warning' | 'danger' | 'info'
  clickable?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  title: '',
  value: 0,
  color: 'default',
  clickable: false
})

const emit = defineEmits<{
  click: [event: MouseEvent]
}>()

const handleClick = (event: MouseEvent) => {
  if (props.clickable) {
    emit('click', event)
  }
}
</script>

<style scoped>
.stats-card {
  background: #fff;
  border-radius: 16px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  border: 1px solid rgba(99, 102, 241, 0.06);
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.stats-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px -4px rgba(99, 102, 241, 0.12);
}

.stats-card--clickable {
  cursor: pointer;
}

.stats-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #F3F4F6;
  color: #6B7280;
  flex-shrink: 0;
}

.stats-card--primary .stats-icon {
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  color: white;
}

.stats-card--success .stats-icon {
  background: linear-gradient(135deg, #10B981 0%, #34D399 100%);
  color: white;
}

.stats-card--warning .stats-icon {
  background: linear-gradient(135deg, #F59E0B 0%, #FBBF24 100%);
  color: white;
}

.stats-card--danger .stats-icon {
  background: linear-gradient(135deg, #EF4444 0%, #F87171 100%);
  color: white;
}

.stats-card--info .stats-icon {
  background: linear-gradient(135deg, #3B82F6 0%, #60A5FA 100%);
  color: white;
}

.stats-info {
  flex: 1;
  min-width: 0;
}

.stats-value {
  font-size: 24px;
  font-weight: 700;
  color: #1E1B4B;
  margin: 0 0 4px 0;
  font-family: 'Fira Sans', sans-serif;
}

.stats-label {
  font-size: 13px;
  color: #6B7280;
  margin: 0;
}

.stats-trend {
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 4px;
  margin: 4px 0 0 0;
  color: #6B7280;
}

.stats-trend.positive {
  color: #10B981;
}

.stats-trend.negative {
  color: #EF4444;
}

@media (max-width: 768px) {
  .stats-card {
    border-radius: 14px;
    padding: 16px;
  }

  .stats-icon {
    width: 48px;
    height: 48px;
  }

  .stats-value {
    font-size: 20px;
  }
}
</style>
