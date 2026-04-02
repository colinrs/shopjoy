<template>
  <div class="refunds-page">
    <!-- Statistics Cards -->
    <el-row
      :gutter="16"
      class="stats-row"
    >
      <el-col
        :xs="12"
        :sm="6"
      >
        <div
          class="stat-item pending"
          @click="handleStatusFilter('0')"
        >
          <p class="stat-number">
            {{ stats.pending }}
          </p>
          <p class="stat-label">
            {{ $t('fulfillment.pendingRefundStatus') }}
          </p>
        </div>
      </el-col>
      <el-col
        :xs="12"
        :sm="6"
      >
        <div
          class="stat-item approved"
          @click="handleStatusFilter('1')"
        >
          <p class="stat-number">
            {{ stats.approved }}
          </p>
          <p class="stat-label">
            {{ $t('fulfillment.approvedRefundStatus') }}
          </p>
        </div>
      </el-col>
      <el-col
        :xs="12"
        :sm="6"
      >
        <div
          class="stat-item rejected"
          @click="handleStatusFilter('2')"
        >
          <p class="stat-number">
            {{ stats.rejected }}
          </p>
          <p class="stat-label">
            {{ $t('fulfillment.rejectedRefundStatus') }}
          </p>
        </div>
      </el-col>
      <el-col
        :xs="12"
        :sm="6"
      >
        <div
          class="stat-item completed"
          @click="handleStatusFilter('3')"
        >
          <p class="stat-number">
            {{ stats.completed }}
          </p>
          <p class="stat-label">
            {{ $t('fulfillment.completedRefundStatus') }}
          </p>
        </div>
      </el-col>
    </el-row>

    <!-- Filter Bar -->
    <el-card
      class="filter-card"
      shadow="never"
    >
      <div class="filter-bar">
        <div class="filter-left">
          <el-input
            v-model="searchQuery"
            :placeholder="$t('fulfillment.searchRefundNo')"
            class="search-input"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select
            v-model="statusFilter"
            :placeholder="$t('fulfillment.refundStatus')"
            clearable
            class="filter-select"
          >
            <el-option
              :label="$t('common.all')"
              value=""
            />
            <el-option
              :label="$t('fulfillment.pendingRefundStatus')"
              value="0"
            />
            <el-option
              :label="$t('fulfillment.approvedRefundStatus')"
              value="1"
            />
            <el-option
              :label="$t('fulfillment.rejectedRefundStatus')"
              value="2"
            />
            <el-option
              :label="$t('fulfillment.completedRefundStatus')"
              value="3"
            />
            <el-option
              :label="$t('fulfillment.cancelledRefundStatus')"
              value="4"
            />
          </el-select>
          <el-select
            v-model="reasonFilter"
            :placeholder="$t('fulfillment.reasonType')"
            clearable
            class="filter-select"
          >
            <el-option
              :label="$t('common.all')"
              value=""
            />
            <el-option
              v-for="reason in refundReasons"
              :key="reason.code"
              :label="reason.name"
              :value="reason.code"
            />
          </el-select>
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="to"
            :start-placeholder="$t('fulfillment.startDate')"
            :end-placeholder="$t('fulfillment.endDate')"
            class="date-picker"
            value-format="YYYY-MM-DD"
          />
        </div>
        <div class="filter-right">
          <el-button @click="handleExport">
            <el-icon><Download /></el-icon>{{ $t('common.export') }}
          </el-button>
          <el-button
            type="primary"
            @click="handleRefresh"
          >
            <el-icon><Refresh /></el-icon>{{ $t('common.refresh') }}
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Pending Alert -->
    <transition name="slide-down">
      <el-alert
        v-if="stats.pending > 0"
        :title="$t('fulfillment.pendingRefundRequests', { n: stats.pending })"
        type="warning"
        :closable="false"
        show-icon
        class="pending-alert"
      >
        <template #action>
          <el-button
            type="warning"
            size="small"
            @click="handleStatusFilter('0')"
          >
            {{ $t('fulfillment.reviewNow') }}
          </el-button>
        </template>
      </el-alert>
    </transition>

    <!-- Refunds Table -->
    <el-card
      class="table-card"
      shadow="never"
    >
      <EmptyState
        v-if="refundList.length === 0 && !loading"
        :title="$t('fulfillment.noRefunds')"
        :description="$t('fulfillment.noRefundsDesc')"
      />
      <el-table
        v-else
        v-loading="loading"
        :data="refundList"
        stripe
      >
        <el-table-column
          prop="refund_no"
          :label="$t('fulfillment.refundNo')"
          min-width="150"
        >
          <template #default="{ row }">
            <div class="refund-no-cell">
              <span class="refund-no">{{ row.refund_no }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          prop="order_no"
          :label="$t('fulfillment.orderNo')"
          min-width="150"
        >
          <template #default="{ row }">
            <el-button
              type="primary"
              link
              size="small"
              @click="viewOrder(row.order_id)"
            >
              {{ row.order_no }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('fulfillment.buyer')"
          min-width="140"
        >
          <template #default="{ row }">
            <div class="buyer-info">
              <p class="buyer-name">
                {{ row.user_name }}
              </p>
              <p class="buyer-phone">
                {{ row.user_phone }}
              </p>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('fulfillment.refundAmount')"
          width="120"
          align="right"
        >
          <template #default="{ row }">
            <div class="amount-cell">
              <p class="refund-amount">
                {{ row.currency }} {{ row.amount }}
              </p>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('fulfillment.reason')"
          min-width="150"
        >
          <template #default="{ row }">
            <div class="reason-cell">
              <el-tag
                size="small"
                effect="plain"
              >
                {{ getReasonName(row.reason_type) }}
              </el-tag>
              <p
                v-if="row.reason"
                class="reason-detail"
              >
                {{ row.reason }}
              </p>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          prop="status"
          :label="$t('common.status')"
          width="100"
          align="center"
        >
          <template #default="{ row }">
            <status-tag
              :status="row.status"
              :type-map="statusTypeMap"
            />
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('fulfillment.appliedAt')"
          width="160"
        >
          <template #default="{ row }">
            <span class="time-text">{{ row.created_at }}</span>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('common.actions')"
          width="180"
          fixed="right"
        >
          <template #default="{ row }">
            <el-button
              type="primary"
              link
              size="small"
              @click="viewDetail(row)"
            >
              {{ $t('fulfillment.detailsAction') }}
            </el-button>
            <el-button
              v-if="row.status === '0'"
              type="success"
              link
              size="small"
              @click="quickApprove(row)"
            >
              {{ $t('fulfillment.approveAction') }}
            </el-button>
            <el-button
              v-if="row.status === '0'"
              type="danger"
              link
              size="small"
              @click="openRejectDialog(row)"
            >
              {{ $t('fulfillment.rejectAction') }}
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

    <!-- Reject Dialog -->
    <el-dialog
      v-model="rejectDialogVisible"
      :title="$t('fulfillment.rejectRefund')"
      width="500px"
    >
      <el-form
        ref="rejectFormRef"
        :model="rejectForm"
        :rules="rejectRules"
        label-width="100px"
      >
        <el-form-item
          :label="$t('fulfillment.rejectReason')"
          prop="reject_reason"
        >
          <el-input
            v-model="rejectForm.reject_reason"
            type="textarea"
            :rows="4"
            :placeholder="$t('fulfillment.pleaseEnterRejectionReason')"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rejectDialogVisible = false">
          {{ $t('common.cancel') }}
        </el-button>
        <el-button
          type="danger"
          :loading="rejecting"
          @click="confirmReject"
        >
          {{ $t('fulfillment.confirmRejectAction') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Download, Refresh } from '@element-plus/icons-vue'
import StatusTag from '@/components/common/StatusTag.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import { t } from '@/plugins/i18n'
import {
  getRefundList,
  approveRefund,
  rejectRefund,
  getRefundReasonList,
  exportRefundsUrl,
  type Refund,
  type RefundReason,
  type RefundListParams
} from '@/api/fulfillment'
import { downloadFile } from '@/utils/download'
import { useErrorHandler } from '@/composables/useErrorHandler'

const router = useRouter()
const { handleError } = useErrorHandler()

// State
const loading = ref(false)
const searchQuery = ref('')
const statusFilter = ref<string | ''>('')
const reasonFilter = ref('')
const dateRange = ref<[string, string] | null>(null)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const rejectDialogVisible = ref(false)
const rejecting = ref(false)
const currentRefund = ref<Refund | null>(null)
const rejectFormRef = ref()

const rejectForm = reactive({
  reject_reason: ''
})

const rejectRules = {
  reject_reason: [{ required: true, message: t('fulfillment.pleaseEnterRejectionReason'), trigger: 'blur' }]
}

const refundReasons = ref<RefundReason[]>([])

const stats = ref({
  pending: 0,
  approved: 0,
  rejected: 0,
  completed: 0
})

const statusTypeMap: Record<string, { type: 'warning' | 'success' | 'danger' | 'primary' | 'info', text: string }> = {
  '0': { type: 'warning', text: 'Pending' },
  '1': { type: 'success', text: 'Approved' },
  '2': { type: 'danger', text: 'Rejected' },
  '3': { type: 'primary', text: 'Completed' },
  '4': { type: 'info', text: 'Cancelled' }
}

// Refund list
const refundList = ref<Refund[]>([])

// Methods
const loadRefundReasons = async () => {
  try {
    const res = await getRefundReasonList()
    refundReasons.value = res
  } catch (error) {
    ElMessage.error(t('fulfillment.loadRefundReasonsFailed'))
  }
}

const loadStats = async () => {
  try {
    // Fetch counts for each status via API (parallel calls for efficiency)
    const [pendingRes, approvedRes, rejectedRes, completedRes] = await Promise.all([
      getRefundList({ status: '0', page_size: 1 }),
      getRefundList({ status: '1', page_size: 1 }),
      getRefundList({ status: '2', page_size: 1 }),
      getRefundList({ status: '3', page_size: 1 })
    ])
    stats.value = {
      pending: pendingRes.total,
      approved: approvedRes.total,
      rejected: rejectedRes.total,
      completed: completedRes.total
    }
  } catch (error) {
    ElMessage.error(t('fulfillment.loadStatisticsFailed'))
  }
}

const loadData = async () => {
  loading.value = true
  try {
    const params: RefundListParams = {
      page: currentPage.value,
      page_size: pageSize.value,
      reason_type: reasonFilter.value || undefined,
      start_time: dateRange.value?.[0],
      end_time: dateRange.value?.[1]
    }
    if (statusFilter.value !== '') {
      params.status = statusFilter.value as import('@/api/fulfillment').RefundStatus
    }
    const res = await getRefundList(params)
    refundList.value = res.list
    total.value = res.total
  } catch (error) {
    ElMessage.error(t('fulfillment.loadRefundsFailed'))
  } finally {
    loading.value = false
  }
}

const getReasonName = (code: string) => {
  const reason = refundReasons.value.find(r => r.code === code)
  return reason?.name || code
}

const handleStatusFilter = (status: string) => {
  statusFilter.value = status
  currentPage.value = 1
  loadData()
}

const handleSearch = () => {
  currentPage.value = 1
  loadData()
}

const handleExport = async () => {
  try {
    const { url, params } = exportRefundsUrl({
      order_no: searchQuery.value || undefined,
      refund_no: searchQuery.value || undefined,
      status: statusFilter.value !== '' ? statusFilter.value : undefined,
      reason_type: reasonFilter.value || undefined,
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
  loadData()
}

const viewOrder = (orderId: number) => {
  router.push(`/orders?id=${orderId}`)
}

const viewDetail = (row: Refund) => {
  router.push(`/fulfillment/refunds/${row.id}`)
}

const quickApprove = async (row: Refund) => {
  try {
    await ElMessageBox.confirm(
      t('fulfillment.approveRefundConfirm', { currency: row.currency, amount: row.amount }),
      t('fulfillment.approveRefund'),
      { type: 'success' }
    )
    await approveRefund(row.id)
    ElMessage.success(t('fulfillment.refundApproved'))
    loadData()
    stats.value.pending--
    stats.value.approved++
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(t('fulfillment.refundApprovedFailed'))
    }
  }
}

const openRejectDialog = (row: Refund) => {
  currentRefund.value = row
  rejectForm.reject_reason = ''
  rejectDialogVisible.value = true
}

const confirmReject = async () => {
  if (!rejectFormRef.value) return

  await rejectFormRef.value.validate(async (valid: boolean) => {
    if (!valid) return

    rejecting.value = true
    try {
      await rejectRefund(currentRefund.value!.id, rejectForm.reject_reason)
      ElMessage.success(t('fulfillment.refundRejected'))
      rejectDialogVisible.value = false
      loadData()
      stats.value.pending--
      stats.value.rejected++
    } catch (error) {
      ElMessage.error(t('fulfillment.refundRejectedFailed'))
    } finally {
      rejecting.value = false
    }
  })
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  loadData()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  loadData()
}

onMounted(() => {
  loadRefundReasons()
  loadStats()
  loadData()
})
</script>

<style scoped>
.refunds-page {
  padding: 0;
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

.stat-item.pending {
  border-left: 4px solid #F59E0B;
}

.stat-item.approved {
  border-left: 4px solid #10B981;
}

.stat-item.rejected {
  border-left: 4px solid #EF4444;
}

.stat-item.completed {
  border-left: 4px solid #6366F1;
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
  width: 220px;
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

/* Pending Alert */
.pending-alert {
  margin-bottom: 20px;
  border-radius: 12px;
}

/* Table */
.table-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

:deep(.el-table__row:hover > td) {
  background-color: #F5F3FF !important;
}

/* Refund No */
.refund-no {
  font-family: 'Fira Code', monospace;
  font-size: 13px;
  color: #6366F1;
  font-weight: 500;
}

/* Buyer Info */
.buyer-info {
  line-height: 1.5;
}

.buyer-name {
  font-weight: 500;
  color: #1E1B4B;
  margin: 0;
}

.buyer-phone {
  font-size: 12px;
  color: #6B7280;
  margin: 4px 0 0 0;
  font-family: 'Fira Code', monospace;
}

/* Amount */
.refund-amount {
  font-weight: 600;
  color: #EF4444;
  margin: 0;
  font-family: 'Fira Sans', sans-serif;
}

/* Reason */
.reason-cell {
  line-height: 1.5;
}

.reason-detail {
  font-size: 12px;
  color: #6B7280;
  margin: 4px 0 0 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 150px;
}

/* Time Text */
.time-text {
  font-size: 13px;
  color: #6B7280;
  font-family: 'Fira Code', monospace;
}

/* Pagination */
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding-top: 20px;
  border-top: 1px solid #F3F4F6;
  margin-top: 20px;
}

/* Transitions */
.slide-down-enter-active,
.slide-down-leave-active {
  transition: all 0.3s ease;
}

.slide-down-enter-from,
.slide-down-leave-to {
  opacity: 0;
  transform: translateY(-10px);
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
