<template>
  <div class="payment-stats">
    <!-- Amount Stats Row -->
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="24" :sm="8">
        <div class="stat-item today">
          <div class="stat-header">
            <el-icon><Calendar /></el-icon>
            <span class="stat-label">Today Received</span>
          </div>
          <p class="stat-amount">{{ currency }} {{ formatAmount(stats.today_received) }}</p>
          <div class="stat-trend" :class="{ up: isPositiveGrowth(stats.today_growth) }">
            <el-icon><TrendCharts /></el-icon>
            <span>{{ formatGrowth(stats.today_growth) }}%</span>
            <span class="trend-label">vs yesterday</span>
          </div>
        </div>
      </el-col>
      <el-col :xs="24" :sm="8">
        <div class="stat-item period">
          <div class="stat-header">
            <el-icon><Timer /></el-icon>
            <span class="stat-label">{{ periodLabel }} Received</span>
          </div>
          <p class="stat-amount">{{ currency }} {{ formatAmount(stats.period_received) }}</p>
        </div>
      </el-col>
      <el-col :xs="24" :sm="8">
        <div class="stat-item refund">
          <div class="stat-header">
            <el-icon><RefreshLeft /></el-icon>
            <span class="stat-label">Refund Amount</span>
          </div>
          <p class="stat-amount refund-amount">{{ currency }} {{ formatAmount(stats.refund_amount) }}</p>
          <div class="stat-trend">
            <span class="refund-rate">Rate: {{ stats.refund_rate }}%</span>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- Channel Distribution Card -->
    <el-card class="channel-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span class="card-title">
            <el-icon><PieChart /></el-icon>
            Payment Channel Distribution
          </span>
        </div>
      </template>
      <el-row :gutter="24">
        <el-col :xs="24" :md="12">
          <!-- Progress Bars -->
          <div class="channel-list">
            <div
              v-for="channel in stats.channel_distribution"
              :key="channel.name"
              class="channel-item"
            >
              <div class="channel-info">
                <span class="channel-name">{{ channel.name }}</span>
                <span class="channel-percent">{{ channel.percent }}%</span>
              </div>
              <el-progress
                :percentage="parseFloat(channel.percent)"
                :color="channel.color"
                :stroke-width="8"
                :show-text="false"
              />
              <div class="channel-amount">
                {{ currency }} {{ formatAmount(channel.amount) }}
                <span class="channel-count">{{ channel.count }} transactions</span>
              </div>
            </div>
          </div>
        </el-col>
        <el-col :xs="24" :md="12">
          <!-- Chart Placeholder -->
          <div class="chart-container">
            <div class="chart-placeholder">
              <el-icon :size="48"><PieChart /></el-icon>
              <p>Channel Distribution</p>
            </div>
          </div>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Calendar, Timer, RefreshLeft, TrendCharts, PieChart } from '@element-plus/icons-vue'
import type { PaymentStats } from '@/api/payment'

const props = withDefaults(defineProps<{
  stats: PaymentStats
  period?: '7d' | '30d' | '90d'
}>(), {
  period: '7d'
})

const currency = computed(() => props.stats.currency || 'USD')

const periodLabel = computed(() => {
  const labels: Record<string, string> = {
    '7d': '7-Day',
    '30d': '30-Day',
    '90d': '90-Day'
  }
  return labels[props.period] || '7-Day'
})

const formatAmount = (amount: string) => {
  const num = parseFloat(amount)
  if (isNaN(num)) return '0.00'
  return num.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const formatGrowth = (growth: string) => {
  const num = parseFloat(growth)
  if (isNaN(num)) return '0.0'
  return (num >= 0 ? '+' : '') + num.toFixed(1)
}

const isPositiveGrowth = (growth: string) => {
  const num = parseFloat(growth)
  return !isNaN(num) && num >= 0
}
</script>

<style scoped>
.payment-stats {
  margin-bottom: 20px;
}

.stats-row {
  margin-bottom: 16px;
}

.stat-item {
  background: #fff;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  border: 1px solid rgba(99, 102, 241, 0.06);
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  margin-bottom: 12px;
}

.stat-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px -4px rgba(99, 102, 241, 0.12);
}

.stat-item.today {
  border-left: 4px solid #6366F1;
}

.stat-item.period {
  border-left: 4px solid #10B981;
}

.stat-item.refund {
  border-left: 4px solid #EF4444;
}

.stat-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.stat-header .el-icon {
  color: #6366F1;
  font-size: 18px;
}

.stat-label {
  font-size: 14px;
  color: #6B7280;
}

.stat-amount {
  font-size: 32px;
  font-weight: 700;
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0 0 8px 0;
  font-family: 'Fira Sans', sans-serif;
}

.stat-amount.refund-amount {
  background: linear-gradient(135deg, #EF4444 0%, #F87171 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.stat-trend {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #6B7280;
}

.stat-trend.up {
  color: #10B981;
}

.stat-trend .el-icon {
  font-size: 14px;
}

.trend-label {
  color: #9CA3AF;
  margin-left: 4px;
}

.refund-rate {
  color: #6B7280;
  font-weight: 500;
}

/* Channel Card */
.channel-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #1E1B4B;
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-title .el-icon {
  color: #6366F1;
}

.channel-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.channel-item {
  padding: 16px;
  background: #F9FAFB;
  border-radius: 12px;
  transition: background 0.2s ease;
}

.channel-item:hover {
  background: #F5F3FF;
}

.channel-info {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.channel-name {
  font-weight: 500;
  color: #1E1B4B;
}

.channel-percent {
  font-weight: 600;
  color: #6366F1;
  font-family: 'Fira Sans', sans-serif;
}

.channel-amount {
  margin-top: 8px;
  font-size: 13px;
  color: #6B7280;
  display: flex;
  justify-content: space-between;
}

.channel-count {
  color: #9CA3AF;
}

.chart-container {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 200px;
  background: #F9FAFB;
  border-radius: 12px;
}

.chart-placeholder {
  text-align: center;
  color: #9CA3AF;
}

.chart-placeholder p {
  margin: 8px 0 0 0;
  font-size: 13px;
}

/* Responsive */
@media (max-width: 768px) {
  .stat-amount {
    font-size: 24px;
  }

  .stat-item {
    margin-bottom: 12px;
  }
}
</style>