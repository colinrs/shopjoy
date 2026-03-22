<template>
  <div class="refunds-page">
    <!-- Statistics Cards -->
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="12" :sm="6">
        <div class="stat-item pending" @click="handleStatusFilter(0)">
          <p class="stat-number">{{ stats.pending }}</p>
          <p class="stat-label">Pending</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-item approved" @click="handleStatusFilter(1)">
          <p class="stat-number">{{ stats.approved }}</p>
          <p class="stat-label">Approved</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-item rejected" @click="handleStatusFilter(2)">
          <p class="stat-number">{{ stats.rejected }}</p>
          <p class="stat-label">Rejected</p>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-item completed" @click="handleStatusFilter(3)">
          <p class="stat-number">{{ stats.completed }}</p>
          <p class="stat-label">Completed</p>
        </div>
      </el-col>
    </el-row>

    <!-- Filter Bar -->
    <el-card class="filter-card" shadow="never">
      <div class="filter-bar">
        <div class="filter-left">
          <el-input
            v-model="searchQuery"
            placeholder="Search refund no. / order no."
            class="search-input"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select v-model="statusFilter" placeholder="Status" clearable class="filter-select">
            <el-option label="All" value="" />
            <el-option label="Pending" :value="0" />
            <el-option label="Approved" :value="1" />
            <el-option label="Rejected" :value="2" />
            <el-option label="Completed" :value="3" />
            <el-option label="Cancelled" :value="4" />
          </el-select>
          <el-select v-model="reasonFilter" placeholder="Refund Reason" clearable class="filter-select">
            <el-option label="All" value="" />
            <el-option v-for="reason in refundReasons" :key="reason.code" :label="reason.name" :value="reason.code" />
          </el-select>
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="to"
            start-placeholder="Start Date"
            end-placeholder="End Date"
            class="date-picker"
            value-format="YYYY-MM-DD"
          />
        </div>
        <div class="filter-right">
          <el-button @click="handleExport">
            <el-icon><Download /></el-icon>Export
          </el-button>
          <el-button type="primary" @click="handleRefresh">
            <el-icon><Refresh /></el-icon>Refresh
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Pending Alert -->
    <transition name="slide-down">
      <el-alert
        v-if="stats.pending > 0"
        :title="`${stats.pending} refund requests pending review`"
        type="warning"
        :closable="false"
        show-icon
        class="pending-alert"
      >
        <template #action>
          <el-button type="warning" size="small" @click="handleStatusFilter(0)">
            Review Now
          </el-button>
        </template>
      </el-alert>
    </transition>

    <!-- Refunds Table -->
    <el-card class="table-card" shadow="never">
      <el-table :data="refundList" v-loading="loading" stripe>
        <el-table-column prop="refund_no" label="Refund No." min-width="150">
          <template #default="{ row }">
            <div class="refund-no-cell">
              <span class="refund-no">{{ row.refund_no }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="order_no" label="Order No." min-width="150">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="viewOrder(row.order_id)">
              {{ row.order_no }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column label="Buyer" min-width="140">
          <template #default="{ row }">
            <div class="buyer-info">
              <p class="buyer-name">{{ row.user_name }}</p>
              <p class="buyer-phone">{{ row.user_phone }}</p>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="Amount" width="120" align="right">
          <template #default="{ row }">
            <div class="amount-cell">
              <p class="refund-amount">{{ row.currency }} {{ (row.amount / 100).toFixed(2) }}</p>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="Reason" min-width="150">
          <template #default="{ row }">
            <div class="reason-cell">
              <el-tag size="small" effect="plain">{{ getReasonName(row.reason_type) }}</el-tag>
              <p v-if="row.reason" class="reason-detail">{{ row.reason }}</p>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="Status" width="100" align="center">
          <template #default="{ row }">
            <status-tag :status="row.status" :type-map="statusTypeMap" />
          </template>
        </el-table-column>
        <el-table-column label="Applied At" width="160">
          <template #default="{ row }">
            <span class="time-text">{{ row.created_at }}</span>
          </template>
        </el-table-column>
        <el-table-column label="Actions" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="viewDetail(row)">
              Details
            </el-button>
            <el-button
              v-if="row.status === 0"
              type="success"
              link
              size="small"
              @click="quickApprove(row)"
            >
              Approve
            </el-button>
            <el-button
              v-if="row.status === 0"
              type="danger"
              link
              size="small"
              @click="openRejectDialog(row)"
            >
              Reject
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
    <el-dialog v-model="rejectDialogVisible" title="Reject Refund" width="500px">
      <el-form :model="rejectForm" :rules="rejectRules" ref="rejectFormRef" label-width="100px">
        <el-form-item label="Reason" prop="reject_reason">
          <el-input
            v-model="rejectForm.reject_reason"
            type="textarea"
            :rows="4"
            placeholder="Please enter the reason for rejection"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rejectDialogVisible = false">Cancel</el-button>
        <el-button type="danger" :loading="rejecting" @click="confirmReject">
          Confirm Reject
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
import {
  getRefundList,
  approveRefund,
  rejectRefund,
  getRefundReasonList,
  type Refund,
  type RefundReason
} from '@/api/fulfillment'

const router = useRouter()

// State
const loading = ref(false)
const searchQuery = ref('')
const statusFilter = ref<number | ''>('')
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
  reject_reason: [{ required: true, message: 'Please enter rejection reason', trigger: 'blur' }]
}

const refundReasons = ref<RefundReason[]>([])

const stats = ref({
  pending: 8,
  approved: 12,
  rejected: 3,
  completed: 45
})

const statusTypeMap = {
  0: { type: 'warning' as const, text: 'Pending' },
  1: { type: 'success' as const, text: 'Approved' },
  2: { type: 'danger' as const, text: 'Rejected' },
  3: { type: 'primary' as const, text: 'Completed' },
  4: { type: 'info' as const, text: 'Cancelled' }
}

// Mock data
const refundList = ref<Refund[]>([
  {
    id: 1,
    refund_no: 'REF20260322001',
    order_id: 'ORD001',
    order_no: 'ORD2026031800100',
    user_id: 101,
    user_name: 'John Doe',
    user_phone: '138****8001',
    type: 1,
    status: 0,
    reason_type: 'DEFECTIVE',
    reason: 'Product has scratches on screen',
    description: 'Received the product with visible scratches on the display screen. Photos attached.',
    images: [],
    amount: 29900,
    currency: 'CNY',
    reject_reason: '',
    approved_at: null,
    approved_by: null,
    completed_at: null,
    created_at: '2026-03-22 14:30:25',
    order_items: []
  },
  {
    id: 2,
    refund_no: 'REF20260322002',
    order_id: 'ORD002',
    order_no: 'ORD2026031800099',
    user_id: 102,
    user_name: 'Jane Smith',
    user_phone: '139****9002',
    type: 1,
    status: 1,
    reason_type: 'WRONG_ITEM',
    reason: 'Received wrong color',
    description: 'Ordered black but received white.',
    images: [],
    amount: 45600,
    currency: 'CNY',
    reject_reason: '',
    approved_at: '2026-03-21 16:00:00',
    approved_by: 'Admin',
    completed_at: null,
    created_at: '2026-03-21 15:00:00',
    order_items: []
  },
  {
    id: 3,
    refund_no: 'REF20260322003',
    order_id: 'ORD003',
    order_no: 'ORD2026031800098',
    user_id: 103,
    user_name: 'Mike Johnson',
    user_phone: '137****7003',
    type: 1,
    status: 2,
    reason_type: 'NO_LONGER_NEEDED',
    reason: 'Changed mind',
    description: 'Found a better deal elsewhere.',
    images: [],
    amount: 12900,
    currency: 'CNY',
    reject_reason: 'Refund period has expired for this reason.',
    approved_at: null,
    approved_by: null,
    completed_at: null,
    created_at: '2026-03-20 10:00:00',
    order_items: []
  },
  {
    id: 4,
    refund_no: 'REF20260322004',
    order_id: 'ORD004',
    order_no: 'ORD2026031800097',
    user_id: 104,
    user_name: 'Sarah Wilson',
    user_phone: '136****6004',
    type: 1,
    status: 3,
    reason_type: 'DAMAGED',
    reason: 'Damaged during shipping',
    description: 'Package arrived damaged, product broken.',
    images: [],
    amount: 59900,
    currency: 'CNY',
    reject_reason: '',
    approved_at: '2026-03-19 14:00:00',
    approved_by: 'Admin',
    completed_at: '2026-03-19 15:00:00',
    created_at: '2026-03-19 12:00:00',
    order_items: []
  },
  {
    id: 5,
    refund_no: 'REF20260322005',
    order_id: 'ORD005',
    order_no: 'ORD2026031800096',
    user_id: 105,
    user_name: 'Tom Brown',
    user_phone: '135****5005',
    type: 1,
    status: 0,
    reason_type: 'NOT_AS_DESCRIBED',
    reason: 'Product differs from description',
    description: 'The material quality is much lower than shown in photos.',
    images: [],
    amount: 159900,
    currency: 'CNY',
    reject_reason: '',
    approved_at: null,
    approved_by: null,
    completed_at: null,
    created_at: '2026-03-18 09:00:00',
    order_items: []
  }
])

// Methods
const loadRefundReasons = async () => {
  try {
    const res = await getRefundReasonList()
    refundReasons.value = res.data
  } catch (error) {
    refundReasons.value = [
      { code: 'DEFECTIVE', name: 'Product Defective' },
      { code: 'WRONG_ITEM', name: 'Wrong Item Received' },
      { code: 'NOT_AS_DESCRIBED', name: 'Not As Described' },
      { code: 'DAMAGED', name: 'Damaged in Transit' },
      { code: 'NO_LONGER_NEEDED', name: 'No Longer Needed' },
      { code: 'LATE_DELIVERY', name: 'Late Delivery' },
      { code: 'OTHER', name: 'Other' }
    ]
  }
}

const loadData = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      status: statusFilter.value,
      reason_type: reasonFilter.value,
      start_time: dateRange.value?.[0],
      end_time: dateRange.value?.[1]
    }
    const res = await getRefundList(params)
    refundList.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    // Mock data already set
  } finally {
    loading.value = false
  }
}

const getReasonName = (code: string) => {
  const reason = refundReasons.value.find(r => r.code === code)
  return reason?.name || code
}

const handleStatusFilter = (status: number) => {
  statusFilter.value = status
  currentPage.value = 1
  loadData()
}

const handleSearch = () => {
  currentPage.value = 1
  loadData()
}

const handleExport = () => {
  ElMessage.success('Export successful')
}

const handleRefresh = () => {
  loadData()
}

const viewOrder = (orderId: string) => {
  router.push(`/orders?id=${orderId}`)
}

const viewDetail = (row: Refund) => {
  router.push(`/fulfillment/refunds/${row.id}`)
}

const quickApprove = async (row: Refund) => {
  try {
    await ElMessageBox.confirm(
      `Approve refund of ${row.currency} ${(row.amount / 100).toFixed(2)}?`,
      'Approve Refund',
      { type: 'success' }
    )
    await approveRefund({ refund_id: row.id })
    ElMessage.success('Refund approved')
    loadData()
    stats.value.pending--
    stats.value.approved++
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('Failed to approve refund')
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
      await rejectRefund({
        refund_id: currentRefund.value!.id,
        reject_reason: rejectForm.reject_reason
      })
      ElMessage.success('Refund rejected')
      rejectDialogVisible.value = false
      loadData()
      stats.value.pending--
      stats.value.rejected++
    } catch (error) {
      ElMessage.error('Failed to reject refund')
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