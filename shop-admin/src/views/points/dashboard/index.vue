<template>
  <div class="points-dashboard">
    <!-- Stats Cards -->
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="12" :sm="6">
        <PointsStatsCard
          :title="$t('points.issuedPoints')"
          :value="stats.total_issued"
          :icon="TrendCharts"
          icon-color="primary"
        />
      </el-col>
      <el-col :xs="12" :sm="6">
        <PointsStatsCard
          :title="$t('points.redeemedPoints')"
          :value="stats.total_redeemed"
          :icon="Present"
          icon-color="success"
        />
      </el-col>
      <el-col :xs="12" :sm="6">
        <PointsStatsCard
          :title="$t('points.activeBalance')"
          :value="stats.outstanding_balance"
          :icon="Star"
          icon-color="warning"
        />
      </el-col>
      <el-col :xs="12" :sm="6">
        <PointsStatsCard
          :title="$t('points.activeUsers')"
          :value="stats.active_users"
          :icon="User"
          icon-color="primary"
        />
      </el-col>
    </el-row>

    <!-- Trend Chart Section -->
    <el-card class="trend-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span class="card-title">
            <el-icon><DataLine /></el-icon>
            {{ $t('points.pointsTrend') }}
          </span>
          <div class="card-actions">
            <el-radio-group v-model="trendPeriod" size="small" @change="loadTrendData">
              <el-radio-button value="7d">{{ $t('points.sevenDays') }}</el-radio-button>
              <el-radio-button value="30d">{{ $t('points.thirtyDays') }}</el-radio-button>
            </el-radio-group>
          </div>
        </div>
      </template>
      <div class="chart-container">
        <div v-if="trendLoading" class="chart-loading">
          <el-icon class="is-loading" :size="32"><Loading /></el-icon>
        </div>
        <div v-else-if="trendData.length === 0" class="chart-empty">
          <el-icon :size="48"><DataLine /></el-icon>
          <p>{{ $t('points.noTrendData') }}</p>
        </div>
        <div v-else class="chart-wrapper">
          <PointsTrendChart :data="trendData" />
        </div>
      </div>
    </el-card>

    <!-- Two Column Section -->
    <el-row :gutter="16" class="bottom-section">
      <!-- Expiring Soon -->
      <el-col :xs="24" :lg="12">
        <el-card class="expiring-card" shadow="never">
          <template #header>
            <span class="card-title">
              <el-icon><Timer /></el-icon>
              {{ $t('points.expiringSoon') }}
            </span>
          </template>
          <div v-if="expiringLoading" class="loading-state">
            <el-icon class="is-loading"><Loading /></el-icon>
          </div>
          <div v-else-if="expiringList.length === 0" class="empty-state">
            <el-icon :size="32"><CircleCheck /></el-icon>
            <p>{{ $t('points.noExpiringPoints') }}</p>
          </div>
          <div v-else class="expiring-list">
            <div
              v-for="item in expiringList"
              :key="item.date"
              class="expiring-item"
            >
              <div class="expiring-date">
                <span class="date">{{ formatDate(item.date) }}</span>
                <span class="days">{{ getDaysUntil(item.date) }} {{ $t('points.daysRemaining') }}</span>
              </div>
              <div class="expiring-points">
                <span class="points">{{ item.points.toLocaleString() }}</span>
                <span class="users">{{ item.user_count }} {{ $t('points.users') }}</span>
              </div>
            </div>
            <div class="expiring-total">
              <span>{{ $t('points.totalLabel') }}:</span>
              <span class="total-points">{{ totalExpiringPoints.toLocaleString() }} {{ $t('points.pointsLabel') }}</span>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- Top Users -->
      <el-col :xs="24" :lg="12">
        <el-card class="top-users-card" shadow="never">
          <template #header>
            <span class="card-title">
              <el-icon><Trophy /></el-icon>
              {{ $t('points.ranking') }}
            </span>
          </template>
          <div v-if="topUsersLoading" class="loading-state">
            <el-icon class="is-loading"><Loading /></el-icon>
          </div>
          <div v-else-if="topUsers.length === 0" class="empty-state">
            <el-icon :size="32"><User /></el-icon>
            <p>{{ $t('points.noUserData') }}</p>
          </div>
          <div v-else class="top-users-list">
            <div
              v-for="(user, index) in topUsers"
              :key="user.user_id"
              class="top-user-item"
            >
              <div class="rank" :class="{ top3: index < 3 }">
                {{ index + 1 }}
              </div>
              <div class="user-info">
                <span class="user-id">U{{ user.user_id }}</span>
                <span v-if="user.user_email" class="user-email">{{ user.user_email }}</span>
              </div>
              <div class="user-points">
                {{ user.points_earned.toLocaleString() }} {{ $t('points.pointsLabel') }}
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  TrendCharts, Present, Star, User, DataLine, Timer, Trophy,
  CircleCheck, Loading
} from '@element-plus/icons-vue'
import PointsStatsCard from '../components/PointsStatsCard.vue'
import PointsTrendChart from '../components/PointsTrendChart.vue'
import {
  getPointsStats,
  getPointsTrend,
  getTopUsers,
  getExpiringPoints,
  type PointsStats,
  type TrendDataPoint,
  type TopUser,
  type ExpiringPoints
} from '@/api/points'
import { useErrorHandler } from '@/composables/useErrorHandler'

const { t } = useI18n()
const { handleError } = useErrorHandler()

// Stats
const stats = ref<PointsStats>({
  total_issued: 0,
  total_redeemed: 0,
  total_expired: 0,
  outstanding_balance: 0,
  redemption_rate: '0',
  active_users: 0,
  period_start: '',
  period_end: ''
})

// Trend
const trendPeriod = ref('7d')
const trendData = ref<TrendDataPoint[]>([])
const trendLoading = ref(false)

// Top Users
const topUsers = ref<TopUser[]>([])
const topUsersLoading = ref(false)

// Expiring
const expiringList = ref<ExpiringPoints[]>([])
const expiringLoading = ref(false)

const totalExpiringPoints = computed(() => {
  return expiringList.value.reduce((sum, item) => sum + item.points, 0)
})

// Methods
const loadStats = async () => {
  try {
    const res = await getPointsStats(trendPeriod.value)
    stats.value = res
  } catch (error) {
    handleError(error, t('points.loadStatsFailed'))
  }
}

const loadTrendData = async () => {
  trendLoading.value = true
  try {
    const res = await getPointsTrend(trendPeriod.value)
    trendData.value = res.data
  } catch (error) {
    handleError(error, t('points.loadTrendFailed'))
  } finally {
    trendLoading.value = false
  }
}

const loadTopUsers = async () => {
  topUsersLoading.value = true
  try {
    const res = await getTopUsers(trendPeriod.value, 10)
    topUsers.value = res.list
  } catch (error) {
    handleError(error, t('points.loadRankingFailed'))
  } finally {
    topUsersLoading.value = false
  }
}

const loadExpiring = async () => {
  expiringLoading.value = true
  try {
    const res = await getExpiringPoints(30)
    expiringList.value = res.list
  } catch (error) {
    handleError(error, t('points.loadExpiringFailed'))
  } finally {
    expiringLoading.value = false
  }
}

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', { month: 'long', day: 'numeric' })
}

const getDaysUntil = (dateStr: string) => {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = date.getTime() - now.getTime()
  return Math.ceil(diff / (1000 * 60 * 60 * 24))
}

// Initialize
onMounted(() => {
  loadStats()
  loadTrendData()
  loadTopUsers()
  loadExpiring()
})
</script>

<style scoped>
.points-dashboard {
  padding: 0;
}

.stats-row {
  margin-bottom: 20px;
}

/* Trend Card */
.trend-card {
  margin-bottom: 20px;
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

.chart-container {
  min-height: 280px;
}

.chart-loading,
.chart-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 280px;
  color: #9CA3AF;
}

.chart-placeholder {
  padding: 16px 0;
}

.chart-wrapper {
  padding: 16px 0;
}

/* Bottom Section */
.bottom-section {
  margin-bottom: 20px;
}

.expiring-card,
.top-users-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
  margin-bottom: 16px;
}

.loading-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
  color: #9CA3AF;
}

.empty-state p {
  margin: 8px 0 0 0;
  font-size: 13px;
}

/* Expiring List */
.expiring-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.expiring-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #F9FAFB;
  border-radius: 12px;
  transition: background 0.2s ease;
}

.expiring-item:hover {
  background: #F5F3FF;
}

.expiring-date {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.expiring-date .date {
  font-weight: 500;
  color: #1E1B4B;
}

.expiring-date .days {
  font-size: 12px;
  color: #F59E0B;
}

.expiring-points {
  text-align: right;
}

.expiring-points .points {
  display: block;
  font-weight: 600;
  color: #EF4444;
  font-family: 'Fira Sans', sans-serif;
}

.expiring-points .users {
  font-size: 12px;
  color: #9CA3AF;
}

.expiring-total {
  display: flex;
  justify-content: space-between;
  padding: 16px;
  background: linear-gradient(135deg, #FEF2F2 0%, #FEE2E2 100%);
  border-radius: 12px;
  margin-top: 8px;
  font-size: 14px;
}

.expiring-total .total-points {
  font-weight: 700;
  color: #EF4444;
  font-family: 'Fira Sans', sans-serif;
}

/* Top Users List */
.top-users-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.top-user-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: #F9FAFB;
  border-radius: 12px;
  transition: background 0.2s ease;
}

.top-user-item:hover {
  background: #F5F3FF;
}

.rank {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 13px;
  background: #E5E7EB;
  color: #6B7280;
}

.rank.top3 {
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  color: white;
}

.user-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.user-info .user-id {
  font-weight: 500;
  color: #1E1B4B;
  font-family: 'Fira Code', monospace;
}

.user-info .user-email {
  font-size: 12px;
  color: #9CA3AF;
}

.user-points {
  font-weight: 600;
  color: #6366F1;
  font-family: 'Fira Sans', sans-serif;
}

/* Responsive */
@media (max-width: 768px) {
  .chart-wrapper {
    padding: 8px 0;
  }
}
</style>