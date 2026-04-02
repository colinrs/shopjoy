<template>
  <div class="transaction-detail-page">
    <!-- Page Header -->
    <el-card
      class="header-card"
      shadow="never"
    >
      <div class="page-header">
        <div class="header-left">
          <el-button
            link
            @click="handleBack"
          >
            <el-icon><ArrowLeft /></el-icon>
            {{ $t('payments.backToList') }}
          </el-button>
          <el-divider direction="vertical" />
          <h2 class="page-title">
            {{ $t('payments.transactionDetail') }}
            <el-tag
              v-if="transaction"
              :type="getStatusType(transaction.status)"
              size="small"
            >
              {{ transaction.status_text }}
            </el-tag>
          </h2>
        </div>
        <div class="header-right">
          <el-button
            :loading="loading"
            @click="handleRefresh"
          >
            <el-icon><Refresh /></el-icon>
            {{ $t('common.refresh') }}
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Loading State -->
    <el-skeleton
      v-if="loading"
      :rows="10"
      animated
    />

    <!-- Transaction Details -->
    <template v-else-if="transaction">
      <!-- Basic Information -->
      <el-card
        class="detail-card"
        shadow="never"
      >
        <template #header>
          <div class="card-header">
            <span class="card-title">{{ $t('payments.basicInfo') }}</span>
          </div>
        </template>
        <el-descriptions
          :column="2"
          border
        >
          <el-descriptions-item :label="$t('payments.transactionId')">
            <div class="id-cell">
              <span class="id-text">{{ transaction.transaction_id }}</span>
              <el-button
                link
                size="small"
                @click="copyTransactionId"
              >
                <el-icon><CopyDocument /></el-icon>
              </el-button>
            </div>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('payments.internalId')">
            {{ transaction.id }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('payments.orderNo')">
            <el-button
              type="primary"
              link
              size="small"
              @click="viewOrder"
            >
              {{ transaction.order_no }}
            </el-button>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('payments.paymentMethod')">
            <el-tag
              :type="getPaymentMethodTagType(transaction.payment_method)"
              effect="plain"
              size="small"
            >
              {{ transaction.payment_method_text }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('payments.channelTransactionId')">
            <span class="channel-id">{{ transaction.channel_transaction_id || '-' }}</span>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('payments.status')">
            <status-tag
              :status="transaction.status"
              :type-map="statusTypeMap"
            />
          </el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- Amount Information -->
      <el-card
        class="detail-card"
        shadow="never"
      >
        <template #header>
          <div class="card-header">
            <span class="card-title">{{ $t('payments.amountInfo') }}</span>
          </div>
        </template>
        <el-descriptions
          :column="3"
          border
        >
          <el-descriptions-item :label="$t('payments.amount')">
            <span class="amount-value">{{ transaction.currency }} {{ formatAmount(transaction.amount) }}</span>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('payments.transactionFee')">
            <span class="fee-value">{{ transaction.currency }} {{ formatAmount(transaction.transaction_fee) }}</span>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('payments.currency')">
            {{ transaction.currency }}
          </el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- Timeline -->
      <el-card
        class="detail-card"
        shadow="never"
      >
        <template #header>
          <div class="card-header">
            <span class="card-title">{{ $t('payments.timeline') }}</span>
          </div>
        </template>
        <el-timeline>
          <el-timeline-item
            :timestamp="formatTime(transaction.created_at)"
            :type="transaction.status === 0 ? 'primary' : 'success'"
            :hollow="transaction.status === 0"
          >
            <p class="timeline-title">
              {{ $t('payments.transactionCreated') }}
            </p>
            <p class="timeline-detail">
              {{ transaction.transaction_id }}
            </p>
          </el-timeline-item>
          <el-timeline-item
            v-if="transaction.paid_at"
            timestamp="formatTime(transaction.paid_at)"
            type="success"
          >
            <p class="timeline-title">
              {{ $t('payments.paymentSucceeded') }}
            </p>
            <p class="timeline-detail">
              {{ transaction.channel_transaction_id }}
            </p>
          </el-timeline-item>
          <el-timeline-item
            v-if="transaction.status === 2"
            :timestamp="formatTime(transaction.created_at)"
            type="danger"
          >
            <p class="timeline-title">
              {{ $t('payments.paymentFailed') }}
            </p>
          </el-timeline-item>
        </el-timeline>
      </el-card>
    </template>

    <!-- Error State -->
    <el-card
      v-else-if="error"
      class="error-card"
      shadow="never"
    >
      <el-empty :description="$t('payments.loadFailed')">
        <el-button
          type="primary"
          @click="handleRefresh"
        >
          {{ $t('common.refresh') }}
        </el-button>
      </el-empty>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Refresh, CopyDocument } from '@element-plus/icons-vue'
import StatusTag from '@/components/common/StatusTag.vue'
import { getPaymentStatusType } from '@/utils/status'
import { getTransactionDetail, type Transaction } from '@/api/payment'
import { t } from '@/plugins/i18n'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const error = ref(false)
const transaction = ref<Transaction | null>(null)

const statusTypeMap = {
  0: { type: 'warning' as const, text: t('payments.pending') },
  1: { type: 'success' as const, text: t('payments.success') },
  2: { type: 'danger' as const, text: t('payments.failed') }
}

const formatAmount = (amount: string) => {
  const num = parseFloat(amount)
  if (isNaN(num)) return '0.00'
  return num.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const formatTime = (timeStr: string | null | undefined) => {
  if (!timeStr) return '-'
  try {
    const date = new Date(timeStr)
    return date.toLocaleString()
  } catch {
    return timeStr
  }
}

const getStatusType = getPaymentStatusType

const getPaymentMethodTagType = (method: string) => {
  const types: Record<string, string> = {
    'stripe_card': 'primary',
    'stripe_alipay': 'success',
    'stripe_wechat': 'warning'
  }
  return types[method] || 'info'
}

const copyTransactionId = () => {
  if (transaction.value?.transaction_id) {
    navigator.clipboard.writeText(transaction.value.transaction_id)
    ElMessage.success(t('payments.transactionCopied'))
  }
}

const viewOrder = () => {
  if (transaction.value?.order_id) {
    router.push(`/orders?id=${transaction.value.order_id}`)
  }
}

const handleBack = () => {
  router.push('/payments')
}

const handleRefresh = () => {
  loadData()
}

const loadData = async () => {
  const id = route.params.id
  if (!id) return

  loading.value = true
  error.value = false

  try {
    const res = await getTransactionDetail(Number(id))
    transaction.value = res
  } catch (err) {
    console.error('Failed to load transaction detail:', err)
    error.value = true
    transaction.value = null
    ElMessage.error(t('payments.loadFailed'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.transaction-detail-page {
  padding: 0;
}

/* Header Card */
.header-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241,  0.06);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.page-title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #1E1B4B;
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-right {
  display: flex;
  gap: 12px;
}

/* Detail Cards */
.detail-card {
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
}

/* ID Cell */
.id-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.id-text {
  font-family: 'Fira Code', monospace;
  font-size: 13px;
  color: #6366F1;
}

.channel-id {
  font-family: 'Fira Code', monospace;
  font-size: 13px;
  color: #6B7280;
}

/* Amount Values */
.amount-value {
  font-weight: 600;
  color: #1E1B4B;
  font-size: 16px;
}

.fee-value {
  color: #6B7280;
}

/* Timeline */
.timeline-title {
  margin: 0;
  font-size: 14px;
  font-weight: 500;
  color: #1E1B4B;
}

.timeline-detail {
  margin: 4px 0 0 0;
  font-size: 13px;
  color: #6B7280;
  font-family: 'Fira Code', monospace;
}

/* Error Card */
.error-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

/* Responsive */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .header-left {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
