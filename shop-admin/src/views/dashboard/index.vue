<template>
  <div class="dashboard">
    <!-- Stats Cards -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="24" :sm="12" :lg="6">
        <div class="stat-card orders">
          <div class="stat-icon">
            <el-icon><ShoppingCart /></el-icon>
          </div>
          <div class="stat-info">
            <p class="stat-label">{{ $t('dashboard.todayOrders') }}</p>
            <p class="stat-value">{{ stats.today_orders }}</p>
            <p class="stat-change" :class="{ positive: isGrowthPositive(stats.today_growth), negative: !isGrowthPositive(stats.today_growth) }">
              <el-icon><ArrowUp v-if="isGrowthPositive(stats.today_growth)" /><ArrowDown v-else /></el-icon>
              {{ stats.today_growth }} {{ $t('dashboard.vsLastPeriod') }}
            </p>
          </div>
        </div>
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <div class="stat-card revenue">
          <div class="stat-icon">
            <el-icon><Money /></el-icon>
          </div>
          <div class="stat-info">
            <p class="stat-label">{{ $t('dashboard.todaySales') }}</p>
            <p class="stat-value">¥{{ formatNumber(stats.today_sales) }}</p>
            <p class="stat-change">
              <span class="neutral">{{ $t('dashboard.yesterdaySales') }} ¥{{ formatNumber(stats.yesterday_sales) }}</span>
            </p>
          </div>
        </div>
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <div class="stat-card products">
          <div class="stat-icon">
            <el-icon><Goods /></el-icon>
          </div>
          <div class="stat-info">
            <p class="stat-label">{{ $t('dashboard.totalProducts') }}</p>
            <p class="stat-value">{{ stats.total_products }}</p>
            <p class="stat-change">
              <span class="neutral">{{ $t('dashboard.activeProducts') }}</span>
            </p>
          </div>
        </div>
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <div class="stat-card users">
          <div class="stat-icon">
            <el-icon><User /></el-icon>
          </div>
          <div class="stat-info">
            <p class="stat-label">{{ $t('dashboard.totalUsers') }}</p>
            <p class="stat-value">{{ stats.total_users }}</p>
            <p class="stat-change positive">
              <el-icon><ArrowUp /></el-icon>+{{ stats.new_users_today }} {{ $t('dashboard.newUsersToday') }}
            </p>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- Charts Row -->
    <el-row :gutter="20" class="charts-row">
      <el-col :xs="24" :lg="16">
        <el-card class="chart-card" v-loading="loading.charts">
          <template #header>
            <div class="card-header">
              <span class="card-title">{{ $t('dashboard.salesTrend') }}</span>
              <el-radio-group v-model="timeRange" size="small" @change="fetchSalesTrend">
                <el-radio-button label="week">{{ $t('dashboard.thisWeek') }}</el-radio-button>
                <el-radio-button label="month">{{ $t('dashboard.thisMonth') }}</el-radio-button>
                <el-radio-button label="year">{{ $t('dashboard.thisYear') }}</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div ref="salesChart" class="chart-container"></div>
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="8">
        <el-card class="chart-card" v-loading="loading.charts">
          <template #header>
            <div class="card-header">
              <span class="card-title">{{ $t('dashboard.orderStatusDistribution') }}</span>
            </div>
          </template>
          <div ref="orderChart" class="chart-container pie-chart"></div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Tables Row -->
    <el-row :gutter="20" class="tables-row">
      <el-col :xs="24" :lg="12">
        <el-card class="table-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><Bell /></el-icon>
                {{ $t('dashboard.pendingOrders') }}
                <el-badge :value="pendingOrdersTotal" :max="99" v-if="pendingOrdersTotal > 0" />
              </span>
              <el-button type="primary" link @click="goToOrders">{{ $t('common.viewAll') }}</el-button>
            </div>
          </template>
          <el-table :data="pendingOrders" stripe style="width: 100%" v-loading="loading.pending">
            <el-table-column prop="order_no" :label="$t('orders.orderNo')" min-width="120">
              <template #default="{ row }">
                <span class="order-no">{{ row.order_no }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="pay_amount" :label="$t('orders.amount')" width="100">
              <template #default="{ row }">
                <span class="amount">¥{{ row.pay_amount }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="status_text" :label="$t('common.status')" width="90">
              <template #default="{ row }">
                <el-tag :type="getOrderStatusType(row.status)" size="small">
                  {{ row.status_text }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('common.actions')" width="80" align="center">
              <template #default="{ row }">
                <el-button type="primary" link size="small" @click="goToOrderDetail(row.id)">{{ $t('dashboard.process') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="12">
        <el-card class="table-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><TrendCharts /></el-icon>
                {{ $t('dashboard.hotProductsTop5') }}
              </span>
              <el-button type="primary" link @click="goToProducts">{{ $t('common.viewAll') }}</el-button>
            </div>
          </template>
          <el-table :data="hotProducts" stripe style="width: 100%" v-loading="loading.products">
            <el-table-column type="index" :label="$t('dashboard.rank')" width="60" align="center">
              <template #default="{ $index }">
                <div class="rank" :class="{ 'top3': $index < 3 }">
                  {{ $index + 1 }}
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="product_name" :label="$t('dashboard.productName')" min-width="150">
              <template #default="{ row }">
                <div class="product-name">
                  <el-avatar :size="32" :src="row.image" shape="square" class="product-avatar">
                    <el-icon><Goods /></el-icon>
                  </el-avatar>
                  <span>{{ row.product_name }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="sales" :label="$t('dashboard.sales')" width="100" align="right">
              <template #default="{ row }">
                <span class="sales-num">{{ row.sales }}</span>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <!-- Activity Row -->
    <el-row :gutter="20" class="activity-row">
      <el-col :xs="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><Timer /></el-icon>
                {{ $t('dashboard.recentActivity') }}
              </span>
            </div>
          </template>
          <el-timeline v-if="recentActivities.length > 0">
            <el-timeline-item
              v-for="(activity, index) in recentActivities"
              :key="index"
              :type="getActivityType(activity.type)"
              :timestamp="formatActivityTime(activity.time)"
            >
              {{ activity.content }}
            </el-timeline-item>
          </el-timeline>
          <el-empty v-else :description="$t('dashboard.noActivityRecords')" :image-size="80" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import * as echarts from 'echarts'
import {
  ShoppingCart, Money, Goods, User, ArrowUp, ArrowDown,
  Bell, TrendCharts, Timer
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { t } from '@/plugins/i18n'
import {
  getDashboardOverview,
  getSalesTrend,
  getOrderStatusDistribution,
  getTopProducts,
  getPendingOrders,
  getRecentActivities,
  type DashboardOverview,
  type SalesTrendData,
  type OrderStatusItem,
  type TopProductItem,
  type PendingOrderItem,
  type ActivityItem
} from '@/api/dashboard'

const router = useRouter()
const timeRange = ref<'week' | 'month' | 'year'>('week')
const loading = ref({
  overview: false,
  charts: false,
  pending: false,
  products: false
})

const salesChart = ref<HTMLElement | null>(null)
const orderChart = ref<HTMLElement | null>(null)
let salesChartInstance: echarts.ECharts | null = null
let orderChartInstance: echarts.ECharts | null = null

// Stats data
const stats = ref<DashboardOverview>({
  today_orders: 0,
  today_sales: '0.00',
  today_growth: '0%',
  yesterday_sales: '0.00',
  total_products: 0,
  total_users: 0,
  new_users_today: 0,
  currency: 'CNY'
})

// Pending orders
const pendingOrders = ref<PendingOrderItem[]>([])
const pendingOrdersTotal = ref(0)

// Hot products
const hotProducts = ref<TopProductItem[]>([])

// Recent activities
const recentActivities = ref<ActivityItem[]>([])

// Sales trend data
const salesTrendData = ref<SalesTrendData[]>([])

// Order status data
const orderStatusData = ref<OrderStatusItem[]>([])

const formatNumber = (num: string | number) => {
  const value = typeof num === 'string' ? parseFloat(num) : num
  return value.toLocaleString()
}

const isGrowthPositive = (growth: string) => {
  return growth.startsWith('+') || growth === '0%' || !growth.startsWith('-')
}

const getOrderStatusType = (status: string) => {
  const types: Record<string, string> = {
    'pending_payment': 'warning',
    'paid': 'primary',
    'shipped': 'info',
    'delivered': 'success',
    'cancelled': 'danger',
    'refunded': 'danger'
  }
  return types[status] || 'info'
}

const getActivityType = (type: string) => {
  const types: Record<string, string> = {
    'order_created': 'primary',
    'payment_received': 'success',
    'product_low_stock': 'warning',
    'order_shipped': 'info'
  }
  return types[type] || 'primary'
}

const formatActivityTime = (time: string) => {
  const date = new Date(time)
  const now = new Date()
  const diff = now.getTime() - date.getTime()

  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (minutes < 60) return t('dashboard.minutesAgo', { n: minutes })
  if (hours < 24) return t('dashboard.hoursAgo', { n: hours })
  if (days < 7) return t('dashboard.daysAgo', { n: days })
  return date.toLocaleDateString()
}

// Fetch dashboard overview
const fetchOverview = async () => {
  loading.value.overview = true
  try {
    stats.value = await getDashboardOverview()
  } catch (error) {
    console.error('Failed to fetch overview:', error)
    ElMessage.error(t('dashboard.loadOverviewFailed'))
  } finally {
    loading.value.overview = false
  }
}

// Fetch sales trend
const fetchSalesTrend = async () => {
  loading.value.charts = true
  try {
    const data = await getSalesTrend({ period: timeRange.value })
    salesTrendData.value = data.data
    updateSalesChart()
  } catch (error) {
    console.error('Failed to fetch sales trend:', error)
    ElMessage.error(t('dashboard.loadSalesTrendFailed'))
  } finally {
    loading.value.charts = false
  }
}

// Fetch order status distribution
const fetchOrderStatus = async () => {
  try {
    const data = await getOrderStatusDistribution()
    orderStatusData.value = data.list
    updateOrderChart()
  } catch (error) {
    console.error('Failed to fetch order status:', error)
    ElMessage.error(t('dashboard.loadOrderStatusFailed'))
  }
}

// Fetch pending orders
const fetchPendingOrders = async () => {
  loading.value.pending = true
  try {
    const data = await getPendingOrders({ limit: 5 })
    pendingOrders.value = data.list
    pendingOrdersTotal.value = data.total
  } catch (error) {
    console.error('Failed to fetch pending orders:', error)
    ElMessage.error(t('dashboard.loadPendingOrdersFailed'))
  } finally {
    loading.value.pending = false
  }
}

// Fetch top products
const fetchTopProducts = async () => {
  loading.value.products = true
  try {
    const data = await getTopProducts({ limit: 5, period: 'week' })
    hotProducts.value = data.list
  } catch (error) {
    console.error('Failed to fetch top products:', error)
    ElMessage.error(t('dashboard.loadHotProductsFailed'))
  } finally {
    loading.value.products = false
  }
}

// Fetch recent activities
const fetchRecentActivities = async () => {
  try {
    const data = await getRecentActivities({ limit: 10 })
    recentActivities.value = data.list
  } catch (error) {
    console.error('Failed to fetch recent activities:', error)
    ElMessage.error(t('dashboard.loadRecentActivityFailed'))
  }
}

// Initialize sales chart
const initSalesChart = () => {
  if (!salesChart.value) return
  salesChartInstance = echarts.init(salesChart.value)
}

// Update sales chart with data
const updateSalesChart = () => {
  if (!salesChartInstance) return

  const dates = salesTrendData.value.map(d => d.date.slice(5)) // MM-DD format
  const sales = salesTrendData.value.map(d => parseFloat(d.sales))

  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'cross' },
      formatter: (params: any) => {
        const point = params[0]
        return `${point.axisValue}<br/>销售额: ¥${formatNumber(point.value)}`
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: dates,
      axisLine: { lineStyle: { color: '#E5E7EB' } },
      axisLabel: { color: '#6B7280' }
    },
    yAxis: {
      type: 'value',
      axisLine: { show: false },
      axisTick: { show: false },
      splitLine: { lineStyle: { color: '#F3F4F6' } },
      axisLabel: {
        color: '#6B7280',
        formatter: (value: number) => {
          if (value >= 10000) return `${(value / 10000).toFixed(1)}万`
          return value.toString()
        }
      }
    },
    series: [
      {
        name: '销售额',
        type: 'line',
        smooth: true,
        symbol: 'circle',
        symbolSize: 6,
        sampling: 'average',
        itemStyle: { color: '#6366F1' },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(99, 102, 241, 0.25)' },
            { offset: 1, color: 'rgba(99, 102, 241, 0.01)' }
          ])
        },
        data: sales
      }
    ]
  }
  salesChartInstance.setOption(option)
}

// Initialize order chart
const initOrderChart = () => {
  if (!orderChart.value) return
  orderChartInstance = echarts.init(orderChart.value)
}

// Update order chart with data
const updateOrderChart = () => {
  if (!orderChartInstance) return

  const pieData = orderStatusData.value.map(item => ({
    value: item.count,
    name: item.status_text,
    itemStyle: { color: item.color }
  }))

  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} ({d}%)'
    },
    legend: {
      bottom: '5%',
      left: 'center'
    },
    series: [
      {
        name: '订单状态',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 16,
            fontWeight: 'bold'
          }
        },
        labelLine: { show: false },
        data: pieData
      }
    ]
  }
  orderChartInstance.setOption(option)
}

const handleResize = () => {
  salesChartInstance?.resize()
  orderChartInstance?.resize()
}

// Navigation handlers
const goToOrders = () => {
  router.push('/orders')
}

const goToProducts = () => {
  router.push('/products')
}

const goToOrderDetail = (orderId: number | string) => {
  router.push(`/orders?id=${orderId}`)
}

// Fetch all dashboard data
const fetchAllData = async () => {
  await Promise.all([
    fetchOverview(),
    fetchSalesTrend(),
    fetchOrderStatus(),
    fetchPendingOrders(),
    fetchTopProducts(),
    fetchRecentActivities()
  ])
}

onMounted(() => {
  initSalesChart()
  initOrderChart()
  fetchAllData()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  salesChartInstance?.dispose()
  orderChartInstance?.dispose()
})
</script>

<style scoped>
.dashboard {
  padding: 0;
}

/* Stats Cards */
.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  background: #fff;
  border-radius: 16px;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 28px -8px rgba(99, 102, 241, 0.15);
  border-color: rgba(99, 102, 241, 0.15);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 26px;
}

.stat-card.orders .stat-icon {
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  color: white;
  box-shadow: 0 8px 16px -4px rgba(99, 102, 241, 0.3);
}

.stat-card.revenue .stat-icon {
  background: linear-gradient(135deg, #F59E0B 0%, #FBBF24 100%);
  color: white;
  box-shadow: 0 8px 16px -4px rgba(245, 158, 11, 0.3);
}

.stat-card.products .stat-icon {
  background: linear-gradient(135deg, #3B82F6 0%, #60A5FA 100%);
  color: white;
  box-shadow: 0 8px 16px -4px rgba(59, 130, 246, 0.3);
}

.stat-card.users .stat-icon {
  background: linear-gradient(135deg, #8B5CF6 0%, #A78BFA 100%);
  color: white;
  box-shadow: 0 8px 16px -4px rgba(139, 92, 246, 0.3);
}

.stat-info {
  flex: 1;
}

.stat-label {
  font-size: 14px;
  color: #6B7280;
  margin: 0 0 4px 0;
  font-weight: 500;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #1E1B4B;
  margin: 0 0 6px 0;
  font-family: 'Fira Sans', sans-serif;
  letter-spacing: -0.5px;
}

.stat-change {
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 4px;
  margin: 0;
}

.stat-change.positive {
  color: #10B981;
}

.stat-change.negative {
  color: #EF4444;
}

.stat-change .neutral {
  color: #6B7280;
}

/* Charts */
.charts-row {
  margin-bottom: 20px;
}

.chart-card {
  height: 400px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.chart-card :deep(.el-card__header) {
  border-bottom: 1px solid #F3F4F6;
  padding: 16px 20px;
}

.chart-container {
  height: 320px;
  width: 100%;
}

.pie-chart {
  height: 280px;
}

/* Tables */
.tables-row {
  margin-bottom: 20px;
}

.table-card {
  height: 380px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.table-card :deep(.el-card__header) {
  border-bottom: 1px solid #F3F4F6;
  padding: 16px 20px;
}

/* Card Header */
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
  color: var(--color-primary);
}

/* Table Styles */
.order-no {
  font-family: 'Fira Code', monospace;
  font-size: 13px;
  color: #6366F1;
  font-weight: 500;
}

.amount {
  font-weight: 600;
  color: #1E1B4B;
}

.rank {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  background: #F5F3FF;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  color: #6B7280;
  margin: 0 auto;
  font-size: 13px;
}

.rank.top3 {
  background: linear-gradient(135deg, #F59E0B 0%, #FBBF24 100%);
  color: white;
  box-shadow: 0 4px 8px -2px rgba(245, 158, 11, 0.3);
}

.product-name {
  display: flex;
  align-items: center;
  gap: 12px;
}

.product-avatar {
  background: linear-gradient(135deg, #EEF2FF 0%, #E0E7FF 100%);
  border-radius: 8px;
}

.product-avatar .el-icon {
  color: #6366F1;
}

.sales-num {
  font-weight: 600;
  color: #10B981;
  font-family: 'Fira Code', monospace;
}

/* Timeline */
:deep(.el-timeline-item__node) {
  background-color: #6366F1;
}

:deep(.el-timeline-item__tail) {
  border-left-color: #E5E7EB;
}

:deep(.el-timeline-item__timestamp) {
  font-family: 'Fira Code', monospace;
  font-size: 12px;
}

/* Radio Button Group */
:deep(.el-radio-button__inner) {
  border-radius: 8px !important;
  border: 1px solid #E5E7EB;
  padding: 6px 14px;
  font-weight: 500;
  transition: all 0.2s ease;
}

:deep(.el-radio-button__original-radio:checked + .el-radio-button__inner) {
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  border-color: #6366F1;
  box-shadow: 0 4px 8px -2px rgba(99, 102, 241, 0.3);
}

/* Table row hover */
:deep(.el-table__row:hover > td) {
  background-color: #F5F3FF !important;
}

/* Responsive */
@media (max-width: 768px) {
  .stat-card {
    margin-bottom: 16px;
    border-radius: 14px;
  }

  .stat-value {
    font-size: 24px;
  }

  .stat-icon {
    width: 48px;
    height: 48px;
    font-size: 22px;
    border-radius: 12px;
  }

  .chart-card,
  .table-card {
    border-radius: 14px;
  }
}
</style>