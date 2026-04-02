<template>
  <div class="payments-page">
    <!-- Payment Stats Card -->
    <PaymentStatsCard :stats="paymentStats" :period="selectedPeriod" />

    <!-- Period Selector -->
    <el-card class="period-card" shadow="never">
      <div class="period-bar">
        <span class="period-label">{{ $t('payments.period') }}:</span>
        <el-radio-group v-model="selectedPeriod" size="small" @change="handlePeriodChange">
          <el-radio-button value="7d">{{ $t('payments.period7d') }}</el-radio-button>
          <el-radio-button value="30d">{{ $t('payments.period30d') }}</el-radio-button>
          <el-radio-button value="90d">{{ $t('payments.period90d') }}</el-radio-button>
        </el-radio-group>
      </div>
    </el-card>

    <!-- Statistics Cards -->
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="12" :sm="8">
        <div class="stat-item success" @click="handleStatusFilter(1)">
          <p class="stat-number">{{ transactionStats.success }}</p>
          <p class="stat-label">{{ $t('payments.success') }}</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="8">
        <div class="stat-item pending" @click="handleStatusFilter(0)">
          <p class="stat-number">{{ transactionStats.pending }}</p>
          <p class="stat-label">{{ $t('payments.pending') }}</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="8">
        <div class="stat-item failed" @click="handleStatusFilter(2)">
          <p class="stat-number">{{ transactionStats.failed }}</p>
          <p class="stat-label">{{ $t('payments.failed') }}</p>
        </div>
      </el-col>
    </el-row>

    <!-- Filter Bar -->
    <el-card class="filter-card" shadow="never">
      <div class="filter-bar">
        <div class="filter-left">
          <el-input
            v-model="searchQuery"
            :placeholder="$t('payments.searchPlaceholder')"
            class="search-input"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select v-model="statusFilter" :placeholder="$t('payments.status')" clearable class="filter-select">
            <el-option :label="$t('payments.all')" value="" />
            <el-option :label="$t('payments.pending')" :value="0" />
            <el-option :label="$t('payments.success')" :value="1" />
            <el-option :label="$t('payments.failed')" :value="2" />
          </el-select>
          <el-select v-model="paymentMethodFilter" :placeholder="$t('payments.paymentMethod')" clearable class="filter-select">
            <el-option :label="$t('payments.all')" value="" />
            <el-option :label="$t('payments.stripeCard')" value="stripe_card" />
            <el-option :label="$t('payments.stripeAlipay')" value="stripe_alipay" />
            <el-option :label="$t('payments.stripeWechat')" value="stripe_wechat" />
          </el-select>
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="to"
            :start-placeholder="$t('payments.startDate')"
            :end-placeholder="$t('payments.endDate')"
            class="date-picker"
            value-format="YYYY-MM-DD"
          />
        </div>
        <div class="filter-right">
          <el-button @click="handleExport">
            <el-icon><Download /></el-icon>
            {{ $t('common.export') }}
          </el-button>
          <el-button type="primary" @click="handleRefresh">
            <el-icon><Refresh /></el-icon>
            {{ $t('common.refresh') }}
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Transactions Table -->
    <el-card class="table-card" shadow="never">
      <el-table :data="transactionList" v-loading="loading" stripe>
        <el-table-column prop="transaction_id" :label="$t('payments.transactionId')" min-width="180">
          <template #default="{ row }">
            <div class="transaction-id-cell">
              <span class="transaction-id" :title="row.transaction_id">
                {{ truncateId(row.transaction_id) }}
              </span>
              <el-button link size="small" @click="copyTransactionId(row.transaction_id)">
                <el-icon><CopyDocument /></el-icon>
              </el-button>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="order_no" :label="$t('payments.orderNo')" min-width="160">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="viewOrder(row.order_id)">
              {{ row.order_no }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column :label="$t('payments.paymentMethod')" width="140" align="center">
          <template #default="{ row }">
            <el-tag :type="getPaymentMethodTagType(row.payment_method)" effect="plain" size="small">
              {{ row.payment_method_text }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('payments.amount')" width="140" align="right">
          <template #default="{ row }">
            <div class="amount-cell">
              <p class="transaction-amount">{{ row.currency }} {{ formatAmount(row.amount) }}</p>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="$t('payments.fee')" width="100" align="right">
          <template #default="{ row }">
            <span class="fee-amount">{{ row.currency }} {{ formatAmount(row.transaction_fee) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('payments.status')" width="100" align="center">
          <template #default="{ row }">
            <status-tag :status="row.status" :type-map="statusTypeMap" />
          </template>
        </el-table-column>
        <el-table-column :label="$t('payments.createdAt')" width="160">
          <template #default="{ row }">
            <span class="time-text">{{ row.created_at }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('payments.paidAt')" width="160">
          <template #default="{ row }">
            <span v-if="row.paid_at" class="time-text">{{ row.paid_at }}</span>
            <span v-else class="no-data">-</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('common.actions')" width="100" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="viewDetail(row)">
              {{ $t('payments.details') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- Pagination -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Search, Download, Refresh, CopyDocument } from '@element-plus/icons-vue'
import StatusTag from '@/components/common/StatusTag.vue'
import PaymentStatsCard from './components/PaymentStatsCard.vue'
import {
  getPaymentStats,
  getTransactionList,
  exportPaymentTransactionsUrl,
  type PaymentStats,
  type Transaction
} from '@/api/payment'
import { t } from '@/plugins/i18n'
import { downloadFile } from '@/utils/download'
import { useErrorHandler } from '@/composables/useErrorHandler'

const router = useRouter()
const { handleError } = useErrorHandler()

// State
const loading = ref(false)
const statsLoading = ref(false)
const searchQuery = ref('')
const statusFilter = ref<number | ''>('')
const paymentMethodFilter = ref('')
const dateRange = ref<[string, string] | null>(null)
const selectedPeriod = ref<'7d' | '30d' | '90d'>('7d')
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const paymentStats = ref<PaymentStats>({
  today_received: '',
  today_growth: '',
  period_received: '',
  refund_amount: '',
  refund_rate: '',
  currency: 'USD',
  channel_distribution: []
})

const transactionStats = ref({
  success: 0,
  pending: 0,
  failed: 0
})

const statusTypeMap = {
  0: { type: 'warning' as const, text: t('payments.pending') },
  1: { type: 'success' as const, text: t('payments.success') },
  2: { type: 'danger' as const, text: t('payments.failed') }
}

const transactionList = ref<Transaction[]>([])

// Methods
const formatAmount = (amount: string) => {
  const num = parseFloat(amount)
  if (isNaN(num)) return '0.00'
  return num.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const truncateId = (id: string) => {
  if (!id || id.length <= 16) return id
  return `${id.slice(0, 8)}...${id.slice(-4)}`
}

const getPaymentMethodTagType = (method: string) => {
  const types: Record<string, string> = {
    'stripe_card': 'primary',
    'stripe_alipay': 'success',
    'stripe_wechat': 'warning'
  }
  return types[method] || 'info'
}

const copyTransactionId = (id: string) => {
  navigator.clipboard.writeText(id)
  ElMessage.success(t('payments.transactionCopied'))
}

const handleStatusFilter = (status: number) => {
  statusFilter.value = status
  currentPage.value = 1
  loadData()
}

const handlePeriodChange = () => {
  loadStats()
}

const handleSearch = () => {
  currentPage.value = 1
  loadData()
}

const handleExport = async () => {
  try {
    const { url, params } = exportPaymentTransactionsUrl({
      order_no: searchQuery.value || undefined,
      transaction_id: searchQuery.value || undefined,
      status: statusFilter.value !== '' ? statusFilter.value : undefined,
      payment_method: paymentMethodFilter.value || undefined,
      start_time: dateRange.value?.[0],
      end_time: dateRange.value?.[1]
    })

    await downloadFile(url, params)
  } catch (error) {
    handleError(error)
    // Error message is handled by downloadFile utility
  }
}

const handleRefresh = () => {
  loadStats()
  loadData()
}

const viewOrder = (orderId: string) => {
  router.push(`/orders?id=${orderId}`)
}

const viewDetail = (row: Transaction) => {
  router.push(`/payments/transactions/${row.id}`)
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  loadData()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  loadData()
}

const loadStats = async () => {
  statsLoading.value = true
  try {
    const res = await getPaymentStats(selectedPeriod.value)
    paymentStats.value = res
  } catch (error) {
    handleError(error, t('payments.loadFailed'))
    // Reset to empty state on error
    paymentStats.value = {
      today_received: '',
      today_growth: '',
      period_received: '',
      refund_amount: '',
      refund_rate: '',
      currency: 'USD',
      channel_distribution: []
    }
  } finally {
    statsLoading.value = false
  }
}

const loadData = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      order_no: searchQuery.value || undefined,
      payment_method: paymentMethodFilter.value || undefined,
      status: statusFilter.value !== '' ? statusFilter.value : undefined,
      start_time: dateRange.value?.[0],
      end_time: dateRange.value?.[1]
    }
    const res = await getTransactionList(params as any)
    transactionList.value = res.list
    total.value = res.total
    transactionStats.value = res.stats
  } catch (error) {
    handleError(error, t('payments.loadFailed'))
    transactionList.value = []
    total.value = 0
    // Do not use mock data on error - keep the previous stats or show empty
    transactionStats.value = { success: 0, pending: 0, failed: 0 }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadStats()
  loadData()
})
</script>

<style scoped>
.payments-page {
  padding: 0;
}

/* Period Card */
.period-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.period-bar {
  display: flex;
  align-items: center;
  gap: 12px;
}

.period-label {
  font-size: 14px;
  font-weight: 500;
  color: #6B7280;
}

/* Stats Row */
.stats-row {
  margin-bottom: 20px;
}

.stat-item {
  background: #fff;
  border-radius: 16px;
  padding: 24px;
  text-align: center;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  border: 1px solid rgba(99, 102, 241, 0.06);
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;
}

.stat-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px -4px rgba(99, 102, 241, 0.12);
}

.stat-item.success {
  border-left: 4px solid #10B981;
}

.stat-item.pending {
  border-left: 4px solid #F59E0B;
}

.stat-item.failed {
  border-left: 4px solid #EF4444;
}

.stat-number {
  font-size: 32px;
  font-weight: 700;
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0 0 6px 0;
  font-family: 'Fira Sans', sans-serif;
}

.stat-label {
  font-size: 14px;
  color: #6B7280;
  margin: 0;
}

/* Filter Bar */
.filter-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
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
}

.search-input {
  width: 240px;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.filter-select {
  width: 140px;
}

.filter-select :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.date-picker {
  width: 260px;
}

.date-picker :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.filter-right {
  display: flex;
  gap: 12px;
}

/* Table */
.table-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

:deep(.el-table__row:hover > td) {
  background-color: #F5F3FF !important;
}

.transaction-id-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.transaction-id {
  font-family: 'Fira Code', monospace;
  font-size: 13px;
  color: #6366F1;
}

.transaction-amount {
  font-weight: 600;
  color: #1E1B4B;
  margin: 0;
  font-family: 'Fira Sans', sans-serif;
}

.fee-amount {
  color: #6B7280;
  font-size: 13px;
}

.time-text {
  font-size: 13px;
  color: #6B7280;
  font-family: 'Fira Code', monospace;
}

.no-data {
  color: #9CA3AF;
  font-size: 13px;
}

/* Pagination */
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding-top: 20px;
  border-top: 1px solid #F3F4F6;
  margin-top: 20px;
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

  .search-input,
  .filter-select,
  .date-picker {
    width: 100%;
  }

  .stat-item {
    border-radius: 14px;
    padding: 20px;
  }

  .stat-number {
    font-size: 28px;
  }
}
</style>