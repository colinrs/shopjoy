<template>
  <div class="transactions-page">
    <!-- Filter Bar -->
    <el-card class="filter-card" shadow="never">
      <div class="filter-bar">
        <div class="filter-left">
          <el-input
            v-model="userIdInput"
            :placeholder="$t('points.searchPlaceholder')"
            clearable
            class="filter-input"
            @keyup.enter="loadTransactions"
          >
            <template #prefix>
              <el-icon><User /></el-icon>
            </template>
          </el-input>
          <el-select v-model="searchParams.type" :placeholder="$t('points.filterType')" clearable class="filter-select" @change="loadTransactions">
            <el-option :label="$t('points.all')" value="" />
            <el-option :label="$t('points.earn')" value="EARN" />
            <el-option :label="$t('points.redeem')" value="REDEEM" />
            <el-option :label="$t('points.adjust')" value="ADJUST" />
            <el-option :label="$t('points.expire')" value="EXPIRE" />
            <el-option :label="$t('points.freeze')" value="FREEZE" />
            <el-option :label="$t('points.unfreeze')" value="UNFREEZE" />
          </el-select>
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="-"
            :start-placeholder="$t('points.startDate')"
            :end-placeholder="$t('points.endDate')"
            value-format="YYYY-MM-DD"
            class="date-picker"
            @change="handleDateChange"
          />
        </div>
        <el-button @click="handleExport">
          <el-icon><Download /></el-icon>
          {{ $t('points.exportCSV') }}
        </el-button>
      </div>
    </el-card>

    <!-- Stats Row -->
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="12" :sm="8">
        <div class="stat-item earned">
          <p class="stat-value">{{ transactionStats.total_earned.toLocaleString() }}</p>
          <p class="stat-label">{{ $t('points.totalEarned') }}</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="8">
        <div class="stat-item redeemed">
          <p class="stat-value">{{ transactionStats.total_redeemed.toLocaleString() }}</p>
          <p class="stat-label">{{ $t('points.totalRedeemed') }}</p>
        </div>
      </el-col>
    </el-row>

    <!-- Transactions Table -->
    <el-card class="table-card" shadow="never">
      <TransactionTable
        :transactions="transactionList"
        :loading="loading"
        :show-user="true"
        :show-expires="true"
        :show-pagination="true"
        :total="total"
        :page="searchParams.page"
        :page-size="searchParams.page_size"
        @page-change="handlePageChange"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { User, Download } from '@element-plus/icons-vue'
import TransactionTable from '../components/TransactionTable.vue'
import { getTransactions, exportPointsTransactionsUrl, type PointsTransaction, type ListTransactionsParams } from '@/api/points'
import { t } from '@/plugins/i18n'
import { downloadFile } from '@/utils/download'
import { useErrorHandler } from '@/composables/useErrorHandler'

const { handleError } = useErrorHandler()

// State
const loading = ref(false)
const dateRange = ref<[string, string] | null>(null)
const userIdInput = ref('')

const transactionList = ref<PointsTransaction[]>([])
const total = ref(0)

const transactionStats = ref({
  total_earned: 0,
  total_redeemed: 0
})

const searchParams = reactive<ListTransactionsParams>({
  page: 1,
  page_size: 20,
  type: '',
  start_time: '',
  end_time: ''
})

// Load functions
const loadTransactions = async () => {
  loading.value = true
  try {
    const params: ListTransactionsParams = {
      page: searchParams.page,
      page_size: searchParams.page_size
    }
    // Parse user_id from string input to number
    if (userIdInput.value && /^\d+$/.test(userIdInput.value.trim())) {
      params.user_id = parseInt(userIdInput.value.trim(), 10)
    }
    if (searchParams.type) params.type = searchParams.type
    if (searchParams.start_time) params.start_time = searchParams.start_time
    if (searchParams.end_time) params.end_time = searchParams.end_time

    const res = await getTransactions(params)
    transactionList.value = res.list || []
    total.value = res.total || 0
    transactionStats.value = res.stats
  } catch (error) {
    handleError(error, t('points.loadTransactionsFailed'))
  } finally {
    loading.value = false
  }
}

// Handlers
const handleDateChange = () => {
  if (dateRange.value) {
    searchParams.start_time = dateRange.value[0]
    searchParams.end_time = dateRange.value[1]
  } else {
    searchParams.start_time = ''
    searchParams.end_time = ''
  }
  loadTransactions()
}

const handlePageChange = (page: number, pageSize: number) => {
  searchParams.page = page
  searchParams.page_size = pageSize
  loadTransactions()
}

const handleExport = async () => {
  try {
    // Parse user_id from string input to number
    const userId = userIdInput.value && /^\d+$/.test(userIdInput.value.trim())
      ? parseInt(userIdInput.value.trim(), 10)
      : undefined

    const { url, params } = exportPointsTransactionsUrl({
      user_id: userId,
      type: searchParams.type || undefined,
      start_time: dateRange.value?.[0],
      end_time: dateRange.value?.[1]
    })

    await downloadFile(url, params)
  } catch (error) {
    handleError(error)
    // Error message is handled by downloadFile utility
  }
}

// Initialize
onMounted(() => {
  loadTransactions()
})
</script>

<style scoped>
.transactions-page {
  padding: 0;
}

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

.filter-input {
  width: 140px;
}

.filter-input :deep(.el-input__wrapper) {
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
}

.stat-item.earned {
  border-left: 4px solid #10B981;
}

.stat-item.redeemed {
  border-left: 4px solid #3B82F6;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  margin: 0 0 6px 0;
  font-family: 'Fira Sans', sans-serif;
}

.stat-item.earned .stat-value {
  color: #10B981;
}

.stat-item.redeemed .stat-value {
  color: #3B82F6;
}

.stat-label {
  font-size: 14px;
  color: #6B7280;
  margin: 0;
}

.table-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
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

  .filter-input,
  .filter-select,
  .date-picker {
    width: 100%;
  }
}
</style>