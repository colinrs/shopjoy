<template>
  <div class="account-detail-page">
    <!-- Back Button -->
    <div class="page-header">
      <el-button link @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
        {{ $t('points.backToAccountList') }}
      </el-button>
      <h2 class="page-title">{{ $t('points.userPointsAccount', { userId: account?.user_id }) }}</h2>
    </div>

    <!-- Account Summary Card -->
    <el-card class="summary-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span class="card-title">
            <el-icon><Wallet /></el-icon>
            {{ $t('points.accountOverview') }}
          </span>
          <el-button type="primary" @click="openAdjustDialog">
            <el-icon><Edit /></el-icon>
            {{ $t('points.adjustPoints') }}
          </el-button>
        </div>
      </template>

      <el-row :gutter="16">
        <el-col :xs="12" :sm="6">
          <div class="summary-item">
            <div class="summary-icon balance">
              <el-icon><Star /></el-icon>
            </div>
            <div class="summary-info">
              <p class="summary-label">{{ $t('points.availableBalance') }}</p>
              <p class="summary-value">{{ account?.balance?.toLocaleString() || 0 }}</p>
            </div>
          </div>
        </el-col>
        <el-col :xs="12" :sm="6">
          <div class="summary-item">
            <div class="summary-icon frozen">
              <el-icon><Lock /></el-icon>
            </div>
            <div class="summary-info">
              <p class="summary-label">{{ $t('points.frozenPoints') }}</p>
              <p class="summary-value">{{ account?.frozen_balance?.toLocaleString() || 0 }}</p>
            </div>
          </div>
        </el-col>
        <el-col :xs="12" :sm="6">
          <div class="summary-item">
            <div class="summary-icon earned">
              <el-icon><TrendCharts /></el-icon>
            </div>
            <div class="summary-info">
              <p class="summary-label">{{ $t('points.cumulativeEarned') }}</p>
              <p class="summary-value">{{ account?.total_earned?.toLocaleString() || 0 }}</p>
            </div>
          </div>
        </el-col>
        <el-col :xs="12" :sm="6">
          <div class="summary-item">
            <div class="summary-icon redeemed">
              <el-icon><Present /></el-icon>
            </div>
            <div class="summary-info">
              <p class="summary-label">{{ $t('points.cumulativeRedeemed') }}</p>
              <p class="summary-value">{{ account?.total_redeemed?.toLocaleString() || 0 }}</p>
            </div>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <!-- Transaction History -->
    <el-card class="transactions-card" shadow="never">
      <template #header>
        <span class="card-title">
          <el-icon><List /></el-icon>
          {{ $t('points.transactionRecords') }}
        </span>
      </template>

      <TransactionTable
        :transactions="transactions"
        :loading="transactionLoading"
        :show-filters="true"
        :show-pagination="true"
        :total="transactionTotal"
        :page="transactionPage"
        :page-size="transactionPageSize"
        @filter="handleFilter"
        @page-change="handlePageChange"
      />
    </el-card>

    <!-- Manual Adjust Dialog -->
    <ManualAdjustDialog
      v-model:visible="adjustDialogVisible"
      :account="account"
      :loading="adjustLoading"
      @submit="handleAdjust"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Wallet, Edit, Star, Lock, TrendCharts, Present, List } from '@element-plus/icons-vue'
import TransactionTable from '../components/TransactionTable.vue'
import ManualAdjustDialog from '../components/ManualAdjustDialog.vue'
import {
  getPointsAccount,
  getAccountTransactions,
  adjustPoints,
  type PointsAccount,
  type PointsTransaction,
  type ListTransactionsParams
} from '@/api/points'
import { useErrorHandler } from '@/composables/useErrorHandler'

const { t } = useI18n()
const { handleError } = useErrorHandler()

const route = useRoute()
const router = useRouter()

// State
const account = ref<PointsAccount | null>(null)
const loading = ref(false)

const transactions = ref<PointsTransaction[]>([])
const transactionLoading = ref(false)
const transactionTotal = ref(0)
const transactionPage = ref(1)
const transactionPageSize = ref(10)

const adjustDialogVisible = ref(false)
const adjustLoading = ref(false)

// Load functions
const loadAccount = async () => {
  const id = parseInt(route.params.id as string)
  if (!id) return

  loading.value = true
  try {
    const res = await getPointsAccount(id)
    account.value = res
  } catch (error) {
    handleError(error, t('points.loadAccountDetailFailed'))
  } finally {
    loading.value = false
  }
}

const loadTransactions = async (filters: { type?: string; startDate?: string; endDate?: string } = {}) => {
  if (!account.value) return

  transactionLoading.value = true
  try {
    const params: ListTransactionsParams = {
      page: transactionPage.value,
      page_size: transactionPageSize.value
    }
    if (filters.type) params.type = filters.type
    if (filters.startDate) params.start_time = filters.startDate
    if (filters.endDate) params.end_time = filters.endDate

    const res = await getAccountTransactions(account.value.id, params)
    transactions.value = res.list || []
    transactionTotal.value = res.total || 0
  } catch (error) {
    handleError(error, t('points.loadTransactionRecordsFailed'))
  } finally {
    transactionLoading.value = false
  }
}

// Handlers
const goBack = () => {
  router.push('/points/accounts')
}

const openAdjustDialog = () => {
  adjustDialogVisible.value = true
}

const handleFilter = (filters: { type: string; startDate: string; endDate: string }) => {
  transactionPage.value = 1
  loadTransactions(filters)
}

const handlePageChange = (page: number, pageSize: number) => {
  transactionPage.value = page
  transactionPageSize.value = pageSize
  loadTransactions()
}

const handleAdjust = async (data: { adjustment_type: 'ADD' | 'DEDUCT'; points: number; reason: string }) => {
  if (!account.value) return

  adjustLoading.value = true
  try {
    await adjustPoints(account.value.id, data)
    ElMessage.success(t('points.adjustSuccess'))
    adjustDialogVisible.value = false
    // Reload account and transactions
    await loadAccount()
    await loadTransactions()
  } catch (error) {
    handleError(error, t('points.adjustFailed'))
  } finally {
    adjustLoading.value = false
  }
}

// Initialize
onMounted(() => {
  loadAccount()
  loadTransactions()
})
</script>

<style scoped>
.account-detail-page {
  padding: 0;
}

.page-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: #1E1B4B;
  margin: 0;
}

/* Summary Card */
.summary-card {
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

.summary-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: #F9FAFB;
  border-radius: 12px;
}

.summary-icon {
  width: 48px;
  height: 48px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.summary-icon.balance {
  background: linear-gradient(135deg, #F59E0B 0%, #FBBF24 100%);
}

.summary-icon.frozen {
  background: linear-gradient(135deg, #6B7280 0%, #9CA3AF 100%);
}

.summary-icon.earned {
  background: linear-gradient(135deg, #10B981 0%, #34D399 100%);
}

.summary-icon.redeemed {
  background: linear-gradient(135deg, #3B82F6 0%, #60A5FA 100%);
}

.summary-info {
  flex: 1;
}

.summary-label {
  font-size: 12px;
  color: #6B7280;
  margin: 0 0 4px 0;
}

.summary-value {
  font-size: 24px;
  font-weight: 700;
  color: #1E1B4B;
  margin: 0;
  font-family: 'Fira Sans', sans-serif;
}

/* Transactions Card */
.transactions-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

/* Responsive */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .summary-item {
    margin-bottom: 12px;
  }

  .summary-value {
    font-size: 20px;
  }
}
</style>