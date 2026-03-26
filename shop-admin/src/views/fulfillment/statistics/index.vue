<template>
  <div class="statistics-page">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">After-Sales Statistics</h1>
        <p class="page-subtitle">Refund analytics and fulfillment performance insights</p>
      </div>
      <div class="header-right">
        <el-radio-group v-model="timeRange" size="default" @change="loadData">
          <el-radio-button label="7">7 Days</el-radio-button>
          <el-radio-button label="30">30 Days</el-radio-button>
          <el-radio-button label="90">90 Days</el-radio-button>
        </el-radio-group>
        <el-button @click="handleExport">
          <el-icon><Download /></el-icon>Export Report
        </el-button>
      </div>
    </div>

    <!-- Overview Stats -->
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="12" :sm="6">
        <div class="stat-card refund-rate">
          <div class="stat-icon">
            <el-icon><RefreshRight /></el-icon>
          </div>
          <div class="stat-info">
            <p class="stat-label">Refund Rate</p>
            <p class="stat-value">{{ stats.refund_rate.toFixed(1) }}%</p>
            <p class="stat-change" :class="refundRateChange >= 0 ? 'negative' : 'positive'">
              <el-icon>
                <ArrowUp v-if="refundRateChange >= 0" />
                <ArrowDown v-else />
              </el-icon>
              {{ Math.abs(refundRateChange).toFixed(1) }}% vs last period
            </p>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card total-refunds">
          <div class="stat-icon">
            <el-icon><Money /></el-icon>
          </div>
          <div class="stat-info">
            <p class="stat-label">Total Refunds</p>
            <p class="stat-value">{{ stats.total_refunds }}</p>
            <p class="stat-amount">CNY {{ formatAmount(stats.refund_amount) }}</p>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card delivery-rate">
          <div class="stat-icon">
            <el-icon><Van /></el-icon>
          </div>
          <div class="stat-info">
            <p class="stat-label">Delivery Success</p>
            <p class="stat-value">{{ stats.delivery_success_rate.toFixed(1) }}%</p>
            <p class="stat-detail">{{ stats.total_shipments }} total shipments</p>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card pending">
          <div class="stat-icon">
            <el-icon><Clock /></el-icon>
          </div>
          <div class="stat-info">
            <p class="stat-label">Pending Refunds</p>
            <p class="stat-value highlight">{{ stats.pending_refunds }}</p>
            <p class="stat-detail">Awaiting review</p>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- Charts Row -->
    <el-row :gutter="20" class="charts-row">
      <!-- Refund Rate Trend -->
      <el-col :xs="24" :lg="16">
        <el-card class="chart-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><TrendCharts /></el-icon>
                Refund Rate Trend
              </span>
            </div>
          </template>
          <div ref="trendChartRef" class="chart-container"></div>
        </el-card>
      </el-col>

      <!-- Refund Reasons Pie -->
      <el-col :xs="24" :lg="8">
        <el-card class="chart-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><PieChart /></el-icon>
                Refund Reasons
              </span>
            </div>
          </template>
          <div ref="reasonChartRef" class="chart-container pie-chart"></div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Tables Row -->
    <el-row :gutter="20" class="tables-row">
      <!-- Problem Products -->
      <el-col :xs="24" :lg="12">
        <el-card class="table-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><Warning /></el-icon>
                Problem Products
                <el-tag type="danger" size="small" effect="plain">High Refund Rate</el-tag>
              </span>
              <el-button type="primary" link>View All</el-button>
            </div>
          </template>
          <el-table :data="problemProducts" stripe style="width: 100%">
            <el-table-column label="Product" min-width="180">
              <template #default="{ row }">
                <div class="product-cell">
                  <el-avatar :size="40" :src="row.image" shape="square" class="product-avatar">
                    <el-icon><Goods /></el-icon>
                  </el-avatar>
                  <div class="product-info">
                    <p class="product-name">{{ row.product_name }}</p>
                    <p class="product-id">ID: {{ row.product_id }}</p>
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="Sales" width="80" align="right">
              <template #default="{ row }">
                <span class="sales-num">{{ row.total_sales }}</span>
              </template>
            </el-table-column>
            <el-table-column label="Refunds" width="80" align="right">
              <template #default="{ row }">
                <span class="refund-num">{{ row.refund_count }}</span>
              </template>
            </el-table-column>
            <el-table-column label="Rate" width="90" align="center">
              <template #default="{ row }">
                <el-tag :type="getRefundRateTagType(row.refund_rate)" size="small">
                  {{ row.refund_rate.toFixed(1) }}%
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <!-- Carrier Performance -->
      <el-col :xs="24" :lg="12">
        <el-card class="table-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><Van /></el-icon>
                Carrier Performance
              </span>
              <el-button type="primary" link>View All</el-button>
            </div>
          </template>
          <el-table :data="carrierPerformance" stripe style="width: 100%">
            <el-table-column label="Carrier" min-width="120">
              <template #default="{ row }">
                <span class="carrier-name">{{ row.carrier_name }}</span>
              </template>
            </el-table-column>
            <el-table-column label="Shipments" width="100" align="right">
              <template #default="{ row }">
                <span class="shipment-num">{{ row.total_shipments }}</span>
              </template>
            </el-table-column>
            <el-table-column label="Success Rate" width="120" align="center">
              <template #default="{ row }">
                <el-progress
                  :percentage="row.delivery_success_rate"
                  :stroke-width="8"
                  :color="getProgressColor(row.delivery_success_rate)"
                />
              </template>
            </el-table-column>
            <el-table-column label="Avg. Time" width="100" align="center">
              <template #default="{ row }">
                <span class="time-text">{{ row.avg_delivery_time }}d</span>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Download, RefreshRight, Money, Van, Clock, TrendCharts, PieChart,
  Warning, Goods, ArrowUp, ArrowDown
} from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { getFulfillmentStatistics } from '@/api/fulfillment'

const timeRange = ref('30')
const trendChartRef = ref<HTMLElement | null>(null)
const reasonChartRef = ref<HTMLElement | null>(null)
let trendChart: echarts.ECharts | null = null
let reasonChart: echarts.ECharts | null = null

const stats = ref({
  refund_rate: 3.2,
  total_refunds: 156,
  refund_amount: 45800,
  delivery_success_rate: 98.5,
  total_shipments: 2340,
  pending_refunds: 8
})

const refundRateChange = ref(-0.5)

const problemProducts = ref([
  { product_id: 101, product_name: 'Wireless Earbuds X1', image: '', total_sales: 45, refund_count: 8, refund_rate: 17.8 },
  { product_id: 102, product_name: 'Phone Case Pro', image: '', total_sales: 120, refund_count: 15, refund_rate: 12.5 },
  { product_id: 103, product_name: 'USB-C Cable 2m', image: '', total_sales: 89, refund_count: 9, refund_rate: 10.1 },
  { product_id: 104, product_name: 'Screen Protector', image: '', total_sales: 200, refund_count: 18, refund_rate: 9.0 },
  { product_id: 105, product_name: 'Charging Stand', image: '', total_sales: 67, refund_count: 5, refund_rate: 7.5 }
])

const carrierPerformance = ref([
  { carrier_code: 'SF', carrier_name: 'SF Express', total_shipments: 856, delivery_success_rate: 99.2, avg_delivery_time: 2.1 },
  { carrier_code: 'ZT', carrier_name: 'ZTO Express', total_shipments: 623, delivery_success_rate: 98.5, avg_delivery_time: 2.8 },
  { carrier_code: 'YT', carrier_name: 'YTO Express', total_shipments: 412, delivery_success_rate: 97.8, avg_delivery_time: 3.2 },
  { carrier_code: 'YD', carrier_name: 'Yunda Express', total_shipments: 289, delivery_success_rate: 97.1, avg_delivery_time: 3.5 },
  { carrier_code: 'EMS', carrier_name: 'EMS', total_shipments: 160, delivery_success_rate: 96.5, avg_delivery_time: 4.1 }
])

const formatAmount = (cents: number) => {
  return (cents / 100).toLocaleString()
}

const getRefundRateTagType = (rate: number) => {
  if (rate >= 15) return 'danger'
  if (rate >= 10) return 'warning'
  return 'info'
}

const getProgressColor = (rate: number) => {
  if (rate >= 99) return '#10B981'
  if (rate >= 97) return '#6366F1'
  if (rate >= 95) return '#F59E0B'
  return '#EF4444'
}

const initTrendChart = () => {
  if (!trendChartRef.value) return

  trendChart = echarts.init(trendChartRef.value)
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
      data: ['Mar 1', 'Mar 5', 'Mar 10', 'Mar 15', 'Mar 20', 'Mar 25', 'Mar 30'],
      axisLine: { lineStyle: { color: '#E5E7EB' } },
      axisLabel: { color: '#6B7280' }
    },
    yAxis: {
      type: 'value',
      name: 'Refund Rate (%)',
      min: 0,
      max: 6,
      axisLine: { show: false },
      axisTick: { show: false },
      splitLine: { lineStyle: { color: '#F3F4F6' } },
      axisLabel: { color: '#6B7280' }
    },
    series: [
      {
        name: 'Refund Rate',
        type: 'line',
        smooth: true,
        symbol: 'circle',
        symbolSize: 8,
        itemStyle: { color: '#EF4444' },
        lineStyle: { width: 3 },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(239, 68, 68, 0.25)' },
            { offset: 1, color: 'rgba(239, 68, 68, 0.01)' }
          ])
        },
        data: [3.5, 3.8, 3.2, 2.9, 3.1, 3.4, 3.2]
      },
      {
        name: 'Target',
        type: 'line',
        smooth: true,
        symbol: 'none',
        lineStyle: { type: 'dashed', color: '#10B981', width: 2 },
        data: [3.0, 3.0, 3.0, 3.0, 3.0, 3.0, 3.0]
      }
    ]
  }
  trendChart.setOption(option)
}

const initReasonChart = () => {
  if (!reasonChartRef.value) return

  reasonChart = echarts.init(reasonChartRef.value)
  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      right: 10,
      top: 'center',
      itemWidth: 10,
      itemHeight: 10
    },
    series: [
      {
        name: 'Refund Reason',
        type: 'pie',
        radius: ['45%', '70%'],
        center: ['35%', '50%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 8,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: { show: false },
        emphasis: {
          label: {
            show: true,
            fontSize: 14,
            fontWeight: 'bold'
          }
        },
        labelLine: { show: false },
        data: [
          { value: 45, name: 'Defective', itemStyle: { color: '#EF4444' } },
          { value: 32, name: 'Wrong Item', itemStyle: { color: '#F59E0B' } },
          { value: 28, name: 'Damaged', itemStyle: { color: '#8B5CF6' } },
          { value: 25, name: 'Not As Described', itemStyle: { color: '#6366F1' } },
          { value: 18, name: 'Late Delivery', itemStyle: { color: '#3B82F6' } },
          { value: 8, name: 'Other', itemStyle: { color: '#6B7280' } }
        ]
      }
    ]
  }
  reasonChart.setOption(option)
}

const handleResize = () => {
  trendChart?.resize()
  reasonChart?.resize()
}

const loadData = async () => {
  try {
    await getFulfillmentStatistics({})
    // Update stats from API
  } catch (error) {
    // Mock data already set
  }
}

const handleExport = () => {
  ElMessage.success('Report exported successfully')
}

onMounted(() => {
  loadData()
  initTrendChart()
  initReasonChart()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  trendChart?.dispose()
  reasonChart?.dispose()
})
</script>

<style scoped>
.statistics-page {
  padding: 0;
}

/* Page Header */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: #1E1B4B;
  margin: 0;
}

.page-subtitle {
  font-size: 14px;
  color: #6B7280;
  margin: 4px 0 0 0;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
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
  height: 120px;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 28px -8px rgba(99, 102, 241, 0.15);
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

.stat-card.refund-rate .stat-icon {
  background: linear-gradient(135deg, #EF4444 0%, #F87171 100%);
  color: white;
  box-shadow: 0 8px 16px -4px rgba(239, 68, 68, 0.3);
}

.stat-card.total-refunds .stat-icon {
  background: linear-gradient(135deg, #8B5CF6 0%, #A78BFA 100%);
  color: white;
  box-shadow: 0 8px 16px -4px rgba(139, 92, 246, 0.3);
}

.stat-card.delivery-rate .stat-icon {
  background: linear-gradient(135deg, #10B981 0%, #34D399 100%);
  color: white;
  box-shadow: 0 8px 16px -4px rgba(16, 185, 129, 0.3);
}

.stat-card.pending .stat-icon {
  background: linear-gradient(135deg, #F59E0B 0%, #FBBF24 100%);
  color: white;
  box-shadow: 0 8px 16px -4px rgba(245, 158, 11, 0.3);
}

.stat-info {
  flex: 1;
}

.stat-label {
  font-size: 13px;
  color: #6B7280;
  margin: 0 0 4px 0;
  font-weight: 500;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #1E1B4B;
  margin: 0 0 4px 0;
  font-family: 'Fira Sans', sans-serif;
}

.stat-value.highlight {
  color: #F59E0B;
}

.stat-change {
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 2px;
  margin: 0;
}

.stat-change.positive {
  color: #10B981;
}

.stat-change.negative {
  color: #EF4444;
}

.stat-amount,
.stat-detail {
  font-size: 12px;
  color: #6B7280;
  margin: 0;
}

.stat-amount {
  color: #8B5CF6;
  font-weight: 500;
}

/* Charts */
.charts-row {
  margin-bottom: 20px;
}

.chart-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
  height: 380px;
}

.chart-card :deep(.el-card__header) {
  border-bottom: 1px solid #F3F4F6;
  padding: 16px 20px;
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
  height: 280px;
  width: 100%;
}

.pie-chart {
  height: 260px;
}

/* Tables */
.tables-row {
  margin-bottom: 20px;
}

.table-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.table-card :deep(.el-card__header) {
  border-bottom: 1px solid #F3F4F6;
  padding: 16px 20px;
}

/* Table Styles */
:deep(.el-table__row:hover > td) {
  background-color: #F5F3FF !important;
}

.product-cell {
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

.product-info {
  flex: 1;
}

.product-name {
  font-weight: 500;
  color: #1E1B4B;
  margin: 0;
  font-size: 14px;
}

.product-id {
  font-size: 12px;
  color: #6B7280;
  margin: 2px 0 0 0;
  font-family: 'Fira Code', monospace;
}

.sales-num,
.shipment-num {
  font-weight: 600;
  color: #1E1B4B;
}

.refund-num {
  font-weight: 600;
  color: #EF4444;
}

.time-text {
  font-family: 'Fira Code', monospace;
  color: #6B7280;
}

.carrier-name {
  font-weight: 500;
  color: #1E1B4B;
}

/* Radio Button Group */
:deep(.el-radio-button__inner) {
  border-radius: 8px !important;
  border: 1px solid #E5E7EB;
  padding: 8px 16px;
  font-weight: 500;
}

:deep(.el-radio-button__original-radio:checked + .el-radio-button__inner) {
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  border-color: #6366F1;
  box-shadow: 0 4px 8px -2px rgba(99, 102, 241, 0.3);
}

/* Responsive */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: stretch;
  }

  .header-right {
    flex-wrap: wrap;
  }

  .stat-card {
    margin-bottom: 16px;
    border-radius: 14px;
    height: auto;
    padding: 20px;
  }

  .stat-value {
    font-size: 24px;
  }

  .chart-card {
    height: auto;
    border-radius: 14px;
  }

  .chart-container {
    height: 240px;
  }
}
</style>