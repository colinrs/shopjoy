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
            <p class="stat-label">今日订单</p>
            <p class="stat-value">{{ stats.todayOrders }}</p>
            <p class="stat-change positive">
              <el-icon><ArrowUp />+12.5%</el-icon> 较昨日
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
            <p class="stat-label">今日销售额</p>
            <p class="stat-value">¥{{ formatNumber(stats.todaySales) }}</p>
            <p class="stat-change positive">
              <el-icon><ArrowUp />+8.2%</el-icon> 较昨日
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
            <p class="stat-label">商品总数</p>
            <p class="stat-value">{{ stats.totalProducts }}</p>
            <p class="stat-change">
              <span class="neutral">持平</span> 较昨日
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
            <p class="stat-label">用户总数</p>
            <p class="stat-value">{{ stats.totalUsers }}</p>
            <p class="stat-change positive">
              <el-icon><ArrowUp />+23</el-icon> 新增用户
            </p>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- Charts Row -->
    <el-row :gutter="20" class="charts-row">
      <el-col :xs="24" :lg="16">
        <el-card class="chart-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">销售趋势</span>
              <el-radio-group v-model="timeRange" size="small">
                <el-radio-button label="week">本周</el-radio-button>
                <el-radio-button label="month">本月</el-radio-button>
                <el-radio-button label="year">全年</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div ref="salesChart" class="chart-container"></div>
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="8">
        <el-card class="chart-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">订单状态分布</span>
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
                待处理订单
              </span>
              <el-button type="primary" link>查看全部</el-button>
            </div>
          </template>
          <el-table :data="pendingOrders" stripe style="width: 100%" v-loading="loading">
            <el-table-column prop="order_no" label="订单号" min-width="120">
              <template #default="{ row }">
                <span class="order-no">{{ row.order_no }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="pay_amount" label="金额" width="100">
              <template #default="{ row }">
                <span class="amount">¥{{ (row.pay_amount / 100).toFixed(2) }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="90">
              <template #default="{ row }">
                <el-tag :type="getOrderStatusType(row.status)" size="small">
                  {{ getOrderStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="80" align="center">
              <template #default>
                <el-button type="primary" link size="small">处理</el-button>
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
                热销商品 TOP5
              </span>
              <el-button type="primary" link>查看全部</el-button>
            </div>
          </template>
          <el-table :data="hotProducts" stripe style="width: 100%" v-loading="loading">
            <el-table-column type="index" label="排名" width="60" align="center">
              <template #default="{ $index }">
                <div class="rank" :class="{ 'top3': $index < 3 }">
                  {{ $index + 1 }}
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="name" label="商品名称" min-width="150">
              <template #default="{ row }">
                <div class="product-name">
                  <el-avatar :size="32" :src="row.image" shape="square" class="product-avatar">
                    <el-icon><Goods /></el-icon>
                  </el-avatar>
                  <span>{{ row.name }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="sales" label="销量" width="100" align="right">
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
                最近活动
              </span>
            </div>
          </template>
          <el-timeline>
            <el-timeline-item
              v-for="(activity, index) in recentActivities"
              :key="index"
              :type="activity.type"
              :timestamp="activity.time"
            >
              {{ activity.content }}
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import * as echarts from 'echarts'
import { 
  ShoppingCart, Money, Goods, User, ArrowUp, 
  Bell, TrendCharts, Timer 
} from '@element-plus/icons-vue'

const timeRange = ref('week')
const loading = ref(false)
const salesChart = ref<HTMLElement | null>(null)
const orderChart = ref<HTMLElement | null>(null)
let salesChartInstance: echarts.ECharts | null = null
let orderChartInstance: echarts.ECharts | null = null

const stats = ref({
  todayOrders: 128,
  todaySales: 15860,
  totalProducts: 486,
  totalUsers: 3256
})

const pendingOrders = ref([
  { order_no: 'ORD20240318001', pay_amount: 29900, status: 1 },
  { order_no: 'ORD20240318002', pay_amount: 15800, status: 2 },
  { order_no: 'ORD20240318003', pay_amount: 8990, status: 1 },
  { order_no: 'ORD20240318004', pay_amount: 45600, status: 2 },
  { order_no: 'ORD20240318005', pay_amount: 12300, status: 1 }
])

const hotProducts = ref([
  { name: '无线蓝牙耳机 Pro', sales: 256, image: '' },
  { name: '智能手表 Series 7', sales: 198, image: '' },
  { name: '便携充电宝 20000mAh', sales: 167, image: '' },
  { name: '机械键盘 RGB', sales: 134, image: '' },
  { name: '4K 高清显示器', sales: 98, image: '' }
])

const recentActivities = ref([
  { content: '新订单 #ORD20240318010 已创建', time: '10分钟前', type: 'primary' },
  { content: '用户 张三 完成了支付', time: '30分钟前', type: 'success' },
  { content: '商品 "无线蓝牙耳机 Pro" 库存不足', time: '1小时前', type: 'warning' },
  { content: '订单 #ORD20240318005 已发货', time: '2小时前', type: 'success' },
  { content: '系统备份完成', time: '3小时前', type: 'info' }
])

const formatNumber = (num: number) => {
  return num.toLocaleString()
}

const getOrderStatusType = (status: number) => {
  const types: Record<number, string> = {
    0: 'warning',
    1: 'success',
    2: 'info',
    3: 'primary',
    4: 'success',
    5: 'danger'
  }
  return types[status] || 'info'
}

const getOrderStatusText = (status: number) => {
  const texts: Record<number, string> = {
    0: '待支付',
    1: '已支付',
    2: '待发货',
    3: '已发货',
    4: '已完成',
    5: '已取消'
  }
  return texts[status] || '未知'
}

const initSalesChart = () => {
  if (!salesChart.value) return

  salesChartInstance = echarts.init(salesChart.value)
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'cross' }
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
      data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
      axisLine: { lineStyle: { color: '#E5E7EB' } },
      axisLabel: { color: '#6B7280' }
    },
    yAxis: {
      type: 'value',
      axisLine: { show: false },
      axisTick: { show: false },
      splitLine: { lineStyle: { color: '#F3F4F6' } },
      axisLabel: { color: '#6B7280' }
    },
    series: [
      {
        name: '销售额',
        type: 'line',
        smooth: true,
        symbol: 'circle',
        symbolSize: 8,
        sampling: 'average',
        itemStyle: { color: '#6366F1' },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(99, 102, 241, 0.25)' },
            { offset: 1, color: 'rgba(99, 102, 241, 0.01)' }
          ])
        },
        data: [8200, 9320, 9010, 14340, 12900, 15300, 15860]
      }
    ]
  }
  salesChartInstance.setOption(option)
}

const initOrderChart = () => {
  if (!orderChart.value) return

  orderChartInstance = echarts.init(orderChart.value)
  const option = {
    tooltip: {
      trigger: 'item'
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
        data: [
          { value: 1048, name: '已完成', itemStyle: { color: '#10B981' } },
          { value: 735, name: '待发货', itemStyle: { color: '#F59E0B' } },
          { value: 580, name: '已发货', itemStyle: { color: '#6366F1' } },
          { value: 234, name: '待支付', itemStyle: { color: '#6B7280' } },
          { value: 148, name: '已取消', itemStyle: { color: '#EF4444' } }
        ]
      }
    ]
  }
  orderChartInstance.setOption(option)
}

const handleResize = () => {
  salesChartInstance?.resize()
  orderChartInstance?.resize()
}

onMounted(() => {
  initSalesChart()
  initOrderChart()
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
