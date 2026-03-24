<template>
  <div class="transactions-page">
    <!-- Filter Bar -->
    <el-card class="filter-card" shadow="never">
      <div class="filter-bar">
        <div class="filter-left">
          <el-input
            v-model="userIdInput"
            placeholder="用户ID"
            clearable
            class="filter-input"
            @keyup.enter="loadTransactions"
          >
            <template #prefix>
              <el-icon><User /></el-icon>
            </template>
          </el-input>
          <el-select v-model="searchParams.type" placeholder="交易类型" clearable class="filter-select" @change="loadTransactions">
            <el-option label="全部" value="" />
            <el-option label="获得" value="EARN" />
            <el-option label="兑换" value="REDEEM" />
            <el-option label="调整" value="ADJUST" />
            <el-option label="过期" value="EXPIRE" />
            <el-option label="冻结" value="FREEZE" />
            <el-option label="解冻" value="UNFREEZE" />
          </el-select>
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            class="date-picker"
            @change="handleDateChange"
          />
        </div>
        <el-button @click="handleExport">
          <el-icon><Download /></el-icon>
          导出CSV
        </el-button>
      </div>
    </el-card>

    <!-- Stats Row -->
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="12" :sm="8">
        <div class="stat-item earned">
          <p class="stat-value">{{ transactionStats.total_earned.toLocaleString() }}</p>
          <p class="stat-label">累计获得</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="8">
        <div class="stat-item redeemed">
          <p class="stat-value">{{ transactionStats.total_redeemed.toLocaleString() }}</p>
          <p class="stat-label">累计兑换</p>
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
import { ElMessage } from 'element-plus'
import { User, Download } from '@element-plus/icons-vue'
import TransactionTable from '../components/TransactionTable.vue'
import { getTransactions, type PointsTransaction, type ListTransactionsParams } from '@/api/points'

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
    transactionList.value = res.data.list || []
    total.value = res.data.total || 0
    transactionStats.value = res.data.stats
  } catch (error) {
    console.error('Failed to load transactions:', error)
    // Mock data
    transactionList.value = [
      {
        id: 1001,
        user_id: 12345,
        account_id: 1,
        points: 520,
        balance_after: 5520,
        type: 'EARN',
        reference_type: 'ORDER',
        reference_id: 'ORD20260324001',
        description: '订单支付积分奖励',
        expires_at: '2027-03-24T00:00:00Z',
        created_at: '2026-03-24T10:30:00Z'
      },
      {
        id: 1002,
        user_id: 12346,
        account_id: 2,
        points: -500,
        balance_after: 3000,
        type: 'REDEEM',
        reference_type: 'REDEEM_RULE',
        reference_id: '1',
        description: '兑换 $10 优惠券',
        expires_at: null,
        created_at: '2026-03-24T09:15:00Z'
      },
      {
        id: 1003,
        user_id: 12345,
        account_id: 1,
        points: 5,
        balance_after: 5505,
        type: 'EARN',
        reference_type: 'SIGN_IN',
        reference_id: '',
        description: '每日签到奖励',
        expires_at: '2026-09-24T00:00:00Z',
        created_at: '2026-03-24T08:00:00Z'
      },
      {
        id: 1004,
        user_id: 12347,
        account_id: 3,
        points: 100,
        balance_after: 1300,
        type: 'ADJUST',
        reference_type: 'MANUAL',
        reference_id: '',
        description: '客服补偿 - 订单发货延迟',
        expires_at: '2027-03-23T00:00:00Z',
        created_at: '2026-03-23T16:00:00Z'
      },
      {
        id: 1005,
        user_id: 12348,
        account_id: 4,
        points: -200,
        balance_after: 500,
        type: 'EXPIRE',
        reference_type: 'SYSTEM',
        reference_id: '',
        description: '积分过期',
        expires_at: null,
        created_at: '2026-03-23T00:00:00Z'
      }
    ]
    total.value = 5
    transactionStats.value = { total_earned: 625, total_redeemed: 700 }
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

const handleExport = () => {
  ElMessage.success('导出功能开发中')
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