<template>
  <div class="refund-detail-page">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-left">
        <el-button link @click="goBack">
          <el-icon><ArrowLeft /></el-icon>
          Back
        </el-button>
        <div class="title-section">
          <h1 class="page-title">Refund Details</h1>
          <p class="refund-no">{{ refund?.refund_no }}</p>
        </div>
      </div>
      <div class="header-right">
        <status-tag :status="refund?.status" :type-map="statusTypeMap" size="large" />
      </div>
    </div>

    <el-row :gutter="20">
      <!-- Left Column -->
      <el-col :xs="24" :lg="16">
        <!-- Refund Info Card -->
        <el-card class="info-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><Money /></el-icon>
                Refund Information
              </span>
            </div>
          </template>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="Refund No.">
              <span class="value-text">{{ refund?.refund_no }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="Order No.">
              <el-button type="primary" link @click="viewOrder">
                {{ refund?.order_no }}
              </el-button>
            </el-descriptions-item>
            <el-descriptions-item label="Refund Amount">
              <span class="refund-amount-value">
                {{ refund?.currency }} {{ refund?.amount || '0.00' }}
              </span>
            </el-descriptions-item>
            <el-descriptions-item label="Refund Type">
              <el-tag size="small">Full Refund</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="Status">
              <status-tag :status="refund?.status" :type-map="statusTypeMap" />
            </el-descriptions-item>
            <el-descriptions-item label="Applied At">
              {{ refund?.created_at }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>

        <!-- Reason Card -->
        <el-card class="info-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><Warning /></el-icon>
                Refund Reason
              </span>
            </div>
          </template>
          <div class="reason-section">
            <div class="reason-type">
              <el-tag effect="plain" size="large">
                {{ getReasonName(refund?.reason_type) }}
              </el-tag>
            </div>
            <div class="reason-summary">
              <p class="summary-label">Summary</p>
              <p class="summary-text">{{ refund?.reason }}</p>
            </div>
            <div v-if="refund?.description" class="reason-detail">
              <p class="detail-label">Description</p>
              <p class="detail-text">{{ refund.description }}</p>
            </div>
          </div>

          <!-- Evidence Images -->
          <div v-if="refund?.images && refund.images.length > 0" class="evidence-section">
            <p class="evidence-label">Evidence Images</p>
            <div class="image-gallery">
              <el-image
                v-for="(img, index) in refund.images"
                :key="index"
                :src="img"
                :preview-src-list="refund.images"
                :initial-index="index"
                fit="cover"
                class="evidence-image"
              >
                <template #error>
                  <div class="image-placeholder">
                    <el-icon><Picture /></el-icon>
                  </div>
                </template>
              </el-image>
            </div>
          </div>
        </el-card>

        <!-- Order Items Card -->
        <el-card class="info-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><Goods /></el-icon>
                Order Items
              </span>
              <span class="item-count">Total {{ refund?.order_items?.length || 0 }} items</span>
            </div>
          </template>
          <div class="items-list">
            <div v-for="item in refund?.order_items" :key="item.id" class="item-row">
              <el-image :src="item.image" class="item-image" fit="cover">
                <template #error>
                  <div class="image-placeholder">
                    <el-icon><Picture /></el-icon>
                  </div>
                </template>
              </el-image>
              <div class="item-info">
                <p class="item-name">{{ item.product_name }}</p>
                <p class="item-sku">SKU: {{ item.sku_name }}</p>
              </div>
              <div class="item-price">
                <p class="unit-price">{{ refund?.currency }} {{ item.price }}</p>
                <p class="quantity">x {{ item.quantity }}</p>
              </div>
            </div>
          </div>
        </el-card>

        <!-- Rejection Info -->
        <el-card v-if="refund?.status === 2 && refund.reject_reason" class="info-card rejection-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><CircleClose /></el-icon>
                Rejection Reason
              </span>
            </div>
          </template>
          <div class="rejection-content">
            <p class="rejection-text">{{ refund.reject_reason }}</p>
            <p class="rejected-by">Rejected by: {{ refund.approved_by }}</p>
          </div>
        </el-card>
      </el-col>

      <!-- Right Column -->
      <el-col :xs="24" :lg="8">
        <!-- Buyer Info Card -->
        <el-card class="info-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><User /></el-icon>
                Buyer Information
              </span>
            </div>
          </template>
          <div class="buyer-section">
            <el-avatar :size="48" class="buyer-avatar">
              {{ refund?.user_name?.charAt(0) }}
            </el-avatar>
            <div class="buyer-details">
              <p class="buyer-name">{{ refund?.user_name }}</p>
              <p class="buyer-contact">
                <el-icon><Phone /></el-icon>
                {{ refund?.user_phone }}
              </p>
            </div>
          </div>
        </el-card>

        <!-- Timeline Card -->
        <el-card class="timeline-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><Clock /></el-icon>
                Status Timeline
              </span>
            </div>
          </template>
          <el-timeline>
            <el-timeline-item
              v-for="(event, index) in timeline"
              :key="index"
              :type="event.type"
              :timestamp="event.time"
              :hollow="!event.active"
              :class="{ 'is-active': event.active }"
            >
              <div class="timeline-content">
                <p class="timeline-title">{{ event.title }}</p>
                <p v-if="event.description" class="timeline-desc">{{ event.description }}</p>
              </div>
            </el-timeline-item>
          </el-timeline>
        </el-card>

        <!-- Actions Card -->
        <el-card v-if="refund?.status === 0" class="actions-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span class="card-title">
                <el-icon><Operation /></el-icon>
                Actions
              </span>
            </div>
          </template>
          <div class="action-buttons">
            <el-button type="success" class="action-btn" @click="handleApprove">
              <el-icon><CircleCheck /></el-icon>
              Approve Refund
            </el-button>
            <el-button type="danger" class="action-btn" @click="openRejectDialog">
              <el-icon><CircleClose /></el-icon>
              Reject Refund
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Reject Dialog -->
    <el-dialog v-model="rejectDialogVisible" title="Reject Refund" width="500px">
      <el-form :model="rejectForm" :rules="rejectRules" ref="rejectFormRef" label-width="100px">
        <el-alert
          type="warning"
          :closable="false"
          style="margin-bottom: 16px"
        >
          Please provide a clear reason for rejection. This will be shown to the buyer.
        </el-alert>
        <el-form-item label="Reason" prop="reject_reason">
          <el-input
            v-model="rejectForm.reject_reason"
            type="textarea"
            :rows="4"
            placeholder="Enter the reason for rejection"
            maxlength="500"
            show-word-limit
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
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  ArrowLeft, Money, Warning, Goods, User, Clock, Operation,
  CircleCheck, CircleClose, Phone, Picture
} from '@element-plus/icons-vue'
import StatusTag from '@/components/common/StatusTag.vue'
import {
  getRefundDetail,
  approveRefund,
  rejectRefund,
  getRefundReasonList,
  type Refund,
  type RefundReason
} from '@/api/fulfillment'

const route = useRoute()
const router = useRouter()

const refund = ref<Refund | null>(null)
const refundReasons = ref<RefundReason[]>([])
const rejectDialogVisible = ref(false)
const rejecting = ref(false)
const rejectFormRef = ref()

const rejectForm = reactive({
  reject_reason: ''
})

const rejectRules = {
  reject_reason: [{ required: true, message: 'Please enter rejection reason', trigger: 'blur' }]
}

const statusTypeMap = {
  0: { type: 'warning' as const, text: 'Pending' },
  1: { type: 'success' as const, text: 'Approved' },
  2: { type: 'danger' as const, text: 'Rejected' },
  3: { type: 'primary' as const, text: 'Completed' },
  4: { type: 'info' as const, text: 'Cancelled' }
}

const timeline = computed(() => {
  if (!refund.value) return []

  type TimelineEvent = {
    title: string
    time: string
    type: 'primary' | 'success' | 'danger' | 'warning' | 'info'
    active: boolean
    description?: string
  }

  const events: TimelineEvent[] = [
    {
      title: 'Refund Requested',
      time: refund.value.created_at,
      type: 'primary',
      active: true,
      description: `Reason: ${getReasonName(refund.value.reason_type)}`
    }
  ]

  if (refund.value.status >= 1 && refund.value.approved_at) {
    events.push({
      title: 'Approved',
      time: refund.value.approved_at,
      type: 'success',
      active: refund.value.status >= 1,
      description: refund.value.approved_by ? `By: ${refund.value.approved_by}` : ''
    })
  }

  if (refund.value.status === 2) {
    events.push({
      title: 'Rejected',
      time: refund.value.approved_at || refund.value.created_at,
      type: 'danger',
      active: true,
      description: refund.value.reject_reason
    })
  }

  if (refund.value.status === 3 && refund.value.completed_at) {
    events.push({
      title: 'Refund Completed',
      time: refund.value.completed_at,
      type: 'primary',
      active: true,
      description: 'Payment has been refunded to buyer'
    })
  }

  if (refund.value.status === 4) {
    events.push({
      title: 'Cancelled',
      time: refund.value.completed_at || refund.value.created_at,
      type: 'info',
      active: true,
      description: 'Request cancelled by buyer'
    })
  }

  return events
})

const loadRefund = async () => {
  const id = route.params.id
  try {
    const res = await getRefundDetail(Number(id))
    refund.value = res
  } catch (error) {
    // Mock data
    refund.value = {
      id: Number(id),
      refund_no: 'REF20260322001',
      order_id: 'ORD001',
      order_no: 'ORD2026031800100',
      user_id: 101,
      user_name: 'John Doe',
      user_phone: '138****8001',
      type: 1,
      type_text: '全额退款',
      status: 0,
      status_text: '待处理',
      reason_type: 'DEFECTIVE',
      reason: 'Product has scratches on screen',
      description: 'Received the product with visible scratches on the display screen. The scratches are clearly visible and affect the user experience. I have attached photos showing the damage.',
      images: [
        'https://picsum.photos/400/300?random=1',
        'https://picsum.photos/400/300?random=2'
      ],
      amount: "299.00",
      currency: 'CNY',
      order_amount: "299.00",
      reject_reason: '',
      approved_at: null,
      approved_by: null,
      completed_at: null,
      created_at: '2026-03-22T14:30:25Z',
      updated_at: '2026-03-22T14:30:25Z',
      order_items: [
        {
          id: 1,
          product_id: 1,
          product_name: 'Wireless Bluetooth Earphones Pro',
          sku_name: 'Black - Premium Edition',
          image: '',
          quantity: 1,
          price: "299.00"
        }
      ]
    }
  }
}

const loadRefundReasons = async () => {
  try {
    const res = await getRefundReasonList()
    refundReasons.value = res
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

const getReasonName = (code?: string) => {
  if (!code) return ''
  const reason = refundReasons.value.find(r => r.code === code)
  return reason?.name || code
}

const goBack = () => {
  router.back()
}

const viewOrder = () => {
  if (refund.value) {
    router.push(`/orders?id=${refund.value.order_id}`)
  }
}

const handleApprove = async () => {
  try {
    await ElMessageBox.confirm(
      `Approve refund of ${refund.value?.currency} ${refund.value?.amount || '0.00'}?`,
      'Approve Refund',
      { type: 'success' }
    )
    await approveRefund(refund.value!.id)
    ElMessage.success('Refund approved successfully')
    loadRefund()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('Failed to approve refund')
    }
  }
}

const openRejectDialog = () => {
  rejectForm.reject_reason = ''
  rejectDialogVisible.value = true
}

const confirmReject = async () => {
  if (!rejectFormRef.value) return

  await rejectFormRef.value.validate(async (valid: boolean) => {
    if (!valid) return

    rejecting.value = true
    try {
      await rejectRefund(refund.value!.id, rejectForm.reject_reason)
      ElMessage.success('Refund rejected')
      rejectDialogVisible.value = false
      loadRefund()
    } catch (error) {
      ElMessage.error('Failed to reject refund')
    } finally {
      rejecting.value = false
    }
  })
}

onMounted(() => {
  loadRefund()
  loadRefundReasons()
})
</script>

<style scoped>
.refund-detail-page {
  padding: 0;
}

/* Page Header */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  padding-bottom: 20px;
  border-bottom: 1px solid #E5E7EB;
}

.header-left {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.title-section {
  display: flex;
  align-items: baseline;
  gap: 16px;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: #1E1B4B;
  margin: 0;
}

.refund-no {
  font-family: 'Fira Code', monospace;
  font-size: 14px;
  color: #6366F1;
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* Cards */
.info-card,
.timeline-card,
.actions-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
  margin-bottom: 20px;
}

.info-card :deep(.el-card__header),
.timeline-card :deep(.el-card__header),
.actions-card :deep(.el-card__header) {
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

.item-count {
  font-size: 13px;
  color: #6B7280;
}

/* Descriptions */
:deep(.el-descriptions__label) {
  font-weight: 500;
  color: #6B7280;
  background: #F9FAFB;
}

:deep(.el-descriptions__content) {
  color: #1E1B4B;
}

.value-text {
  font-family: 'Fira Code', monospace;
  color: #6366F1;
}

.refund-amount-value {
  font-weight: 700;
  color: #EF4444;
  font-size: 18px;
}

/* Reason Section */
.reason-section {
  padding: 8px 0;
}

.reason-type {
  margin-bottom: 16px;
}

.reason-summary,
.reason-detail {
  margin-bottom: 12px;
}

.summary-label,
.detail-label {
  font-size: 13px;
  font-weight: 500;
  color: #6B7280;
  margin: 0 0 8px 0;
}

.summary-text,
.detail-text {
  font-size: 14px;
  color: #1E1B4B;
  margin: 0;
  line-height: 1.6;
}

.detail-text {
  background: #F9FAFB;
  padding: 12px 16px;
  border-radius: 8px;
}

/* Evidence Section */
.evidence-section {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #F3F4F6;
}

.evidence-label {
  font-size: 13px;
  font-weight: 500;
  color: #6B7280;
  margin: 0 0 12px 0;
}

.image-gallery {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.evidence-image {
  width: 120px;
  height: 90px;
  border-radius: 8px;
  cursor: pointer;
}

.image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #F5F3FF 0%, #EEF2FF 100%);
  color: #6366F1;
}

/* Items List */
.items-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.item-row {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px;
  background: #F9FAFB;
  border-radius: 12px;
  transition: all 0.2s ease;
}

.item-row:hover {
  background: #F5F3FF;
}

.item-image {
  width: 60px;
  height: 60px;
  border-radius: 10px;
  flex-shrink: 0;
}

.item-info {
  flex: 1;
}

.item-name {
  font-weight: 500;
  color: #1E1B4B;
  margin: 0 0 4px 0;
}

.item-sku {
  font-size: 13px;
  color: #6B7280;
  margin: 0;
}

.item-price {
  text-align: right;
}

.unit-price {
  font-weight: 600;
  color: #1E1B4B;
  margin: 0;
}

.quantity {
  font-size: 13px;
  color: #6B7280;
  margin: 4px 0 0 0;
}

/* Rejection Card */
.rejection-card {
  border-left: 4px solid #EF4444;
}

.rejection-content {
  padding: 8px 0;
}

.rejection-text {
  font-size: 14px;
  color: #EF4444;
  margin: 0 0 8px 0;
  line-height: 1.6;
}

.rejected-by {
  font-size: 12px;
  color: #6B7280;
  margin: 0;
}

/* Buyer Section */
.buyer-section {
  display: flex;
  align-items: center;
  gap: 16px;
}

.buyer-avatar {
  background: linear-gradient(135deg, #6366F1 0%, #818CF8 100%);
  color: white;
  font-weight: 600;
  font-size: 18px;
}

.buyer-details {
  flex: 1;
}

.buyer-name {
  font-weight: 600;
  color: #1E1B4B;
  margin: 0 0 4px 0;
}

.buyer-contact {
  font-size: 13px;
  color: #6B7280;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 4px;
}

/* Timeline */
.timeline-content {
  padding-left: 4px;
}

.timeline-title {
  font-weight: 500;
  color: #1E1B4B;
  margin: 0;
}

.timeline-desc {
  font-size: 13px;
  color: #6B7280;
  margin: 4px 0 0 0;
}

:deep(.el-timeline-item.is-active .el-timeline-item__node) {
  background-color: #6366F1;
}

:deep(.el-timeline-item.is-active .el-timeline-item__tail) {
  border-left-color: #6366F1;
}

/* Actions */
.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.action-btn {
  width: 100%;
  height: 44px;
  justify-content: flex-start;
  padding: 0 20px;
}

/* Responsive */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
  }

  .header-right {
    width: 100%;
    justify-content: space-between;
  }

  .info-card,
  .timeline-card,
  .actions-card {
    border-radius: 14px;
  }
}
</style>