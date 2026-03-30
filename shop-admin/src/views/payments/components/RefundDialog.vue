<template>
  <el-dialog
    v-model="visible"
    :title="$t('payments.initiateRefund')"
    width="560px"
    :close-on-click-modal="false"
    destroy-on-close
  >
    <el-steps :active="currentStep" simple class="refund-steps">
      <el-step :title="$t('payments.selectType')" />
      <el-step :title="$t('payments.confirm')" />
    </el-steps>

    <!-- Step 1: Type & Reason -->
    <div v-show="currentStep === 0" class="step-content">
      <!-- Order Info -->
      <div class="order-info-banner">
        <el-icon><Document /></el-icon>
        <div class="order-info-content">
          <p class="order-no">{{ $t('payments.order') }}: {{ order?.order_no }}</p>
          <p class="order-amount">{{ $t('payments.paid') }}: {{ order?.currency }} {{ formatAmount(maxRefundable) }}</p>
        </div>
      </div>

      <!-- Refund Type -->
      <div class="refund-type-section">
        <p class="section-label">{{ $t('payments.refundType') }}</p>
        <el-radio-group v-model="form.refundType" class="refund-type-group">
          <el-radio-button value="full">
            <div class="type-option">
              <el-icon><FullScreen /></el-icon>
              <span class="type-title">{{ $t('payments.fullRefund') }}</span>
              <span class="type-desc">{{ $t('payments.refundFullAmount') }}</span>
            </div>
          </el-radio-button>
          <el-radio-button value="partial">
            <div class="type-option">
              <el-icon><Compass /></el-icon>
              <span class="type-title">{{ $t('payments.partialRefund') }}</span>
              <span class="type-desc">{{ $t('payments.refundPartOfAmount') }}</span>
            </div>
          </el-radio-button>
        </el-radio-group>
      </div>

      <!-- Partial Amount Input -->
      <transition name="slide-down">
        <div v-if="form.refundType === 'partial'" class="amount-input-section">
          <p class="section-label">{{ $t('payments.refundAmount') }}</p>
          <div class="amount-input-wrapper">
            <span class="currency-label">{{ order?.currency }}</span>
            <el-input-number
              v-model="form.refundAmount"
              :min="0.01"
              :max="maxRefundable / 100"
              :precision="2"
              :controls="false"
              class="amount-input"
            />
          </div>
          <div class="amount-hint">
            <span>{{ $t('payments.maxRefundable') }}: {{ order?.currency }} {{ formatAmount(maxRefundable) }}</span>
            <el-progress
              :percentage="refundPercentage"
              :stroke-width="4"
              :show-text="false"
              class="amount-progress"
            />
          </div>
        </div>
      </transition>

      <!-- Reason -->
      <div class="reason-section">
        <p class="section-label">{{ $t('payments.refundReason') }}</p>
        <el-select v-model="form.reasonType" :placeholder="$t('payments.selectRefundReason')" style="width: 100%">
          <el-option
            v-for="reason in refundReasons"
            :key="reason.code"
            :label="reason.name"
            :value="reason.code"
          />
        </el-select>
      </div>

      <!-- Notes -->
      <div class="notes-section">
        <p class="section-label">{{ $t('payments.additionalNotes') }}</p>
        <el-input
          v-model="form.notes"
          type="textarea"
          :rows="3"
          :placeholder="$t('payments.enterAdditionalInfo')"
          maxlength="500"
          show-word-limit
        />
      </div>
    </div>

    <!-- Step 2: Confirmation -->
    <div v-show="currentStep === 1" class="step-content">
      <el-alert type="warning" :closable="false" class="confirm-alert">
        <template #title>
          {{ $t('payments.confirmRefundDetails') }}
        </template>
      </el-alert>

      <div class="confirm-summary">
        <el-descriptions :column="1" border>
          <el-descriptions-item :label="$t('payments.orderNo')">
            <span class="value-text">{{ order?.order_no }}</span>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('payments.refundType')">
            <el-tag :type="form.refundType === 'full' ? 'primary' : 'warning'" size="small">
              {{ form.refundType === 'full' ? $t('payments.fullRefund') : $t('payments.partialRefund') }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('payments.refundAmount')">
            <span class="confirm-amount">
              {{ order?.currency }} {{ formatAmount(actualRefundAmount) }}
            </span>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('payments.reason')">
            {{ getReasonName(form.reasonType) }}
          </el-descriptions-item>
          <el-descriptions-item v-if="form.notes" :label="$t('payments.notes')">
            {{ form.notes }}
          </el-descriptions-item>
        </el-descriptions>

        <div v-if="form.refundType === 'partial'" class="remaining-info">
          <el-icon><InfoFilled /></el-icon>
          <span>{{ $t('payments.remainingOrderAmount') }}: {{ order?.currency }} {{ formatAmount(remainingAmount) }}</span>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button v-if="currentStep > 0" @click="prevStep">
          <el-icon><ArrowLeft /></el-icon>
          {{ $t('payments.back') }}
        </el-button>
        <el-button @click="handleCancel">{{ $t('common.cancel') }}</el-button>
        <el-button
          v-if="currentStep === 0"
          type="primary"
          :disabled="!canProceed"
          @click="nextStep"
        >
          {{ $t('payments.nextStep') }}
          <el-icon><ArrowRight /></el-icon>
        </el-button>
        <el-button
          v-else
          type="danger"
          :loading="submitting"
          @click="confirmRefund"
        >
          <el-icon><Check /></el-icon>
          {{ $t('payments.confirmRefund') }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Document,
  FullScreen,
  Compass,
  InfoFilled,
  ArrowLeft,
  ArrowRight,
  Check
} from '@element-plus/icons-vue'
import { initiateRefund, type OrderPayment } from '@/api/payment'
import { t } from '@/plugins/i18n'

// Generate UUID for idempotency
const generateIdempotencyKey = () => {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    const r = Math.random() * 16 | 0
    const v = c === 'x' ? r : (r & 0x3 | 0x8)
    return v.toString(16)
  })
}

const props = defineProps<{
  modelValue: boolean
  order: OrderPayment | null
  maxRefundable: number  // In cents
  orderId: number
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}>()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const currentStep = ref(0)
const submitting = ref(false)

const form = reactive({
  idempotencyKey: generateIdempotencyKey(),
  refundType: 'full' as 'full' | 'partial',
  refundAmount: 0,
  reasonType: '',
  notes: ''
})

const refundReasons = [
  { code: 'DEFECTIVE', name: 'Product Defective' },
  { code: 'WRONG_ITEM', name: 'Wrong Item Received' },
  { code: 'NOT_AS_DESCRIBED', name: 'Not As Described' },
  { code: 'DAMAGED', name: 'Damaged in Transit' },
  { code: 'NO_LONGER_NEEDED', name: 'No Longer Needed' },
  { code: 'LATE_DELIVERY', name: 'Late Delivery' },
  { code: 'OTHER', name: 'Other' }
]

// Computed
const actualRefundAmount = computed(() => {
  if (form.refundType === 'full') {
    return props.maxRefundable
  }
  return Math.round(form.refundAmount * 100)
})

const remainingAmount = computed(() => {
  return props.maxRefundable - actualRefundAmount.value
})

const refundPercentage = computed(() => {
  if (props.maxRefundable === 0) return 0
  return Math.round((actualRefundAmount.value / props.maxRefundable) * 100)
})

const canProceed = computed(() => {
  if (!form.reasonType) return false
  if (form.refundType === 'partial' && form.refundAmount <= 0) return false
  return true
})

// Methods
const formatAmount = (amount: number) => {
  return (amount / 100).toFixed(2)
}

const getReasonName = (code: string) => {
  const reason = refundReasons.find(r => r.code === code)
  return reason?.name || code
}

const nextStep = () => {
  if (canProceed.value) {
    currentStep.value = 1
  }
}

const prevStep = () => {
  currentStep.value = 0
}

const handleCancel = () => {
  visible.value = false
}

const confirmRefund = async () => {
  submitting.value = true
  try {
    const amountStr = (actualRefundAmount.value / 100).toFixed(2)
    await initiateRefund(props.orderId, {
      idempotency_key: form.idempotencyKey,
      amount: amountStr,
      reason_type: form.reasonType,
      reason: form.notes || undefined
    })
    ElMessage.success(t('payments.refundInitiatedSuccessfully'))
    emit('success')
    visible.value = false
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.msg || t('payments.failedToInitiateRefund'))
  } finally {
    submitting.value = false
  }
}

// Reset form when dialog opens
watch(visible, (val) => {
  if (val) {
    form.idempotencyKey = generateIdempotencyKey()
    form.refundType = 'full'
    form.refundAmount = 0
    form.reasonType = ''
    form.notes = ''
    currentStep.value = 0
  }
})
</script>

<style scoped>
.refund-steps {
  margin-bottom: 24px;
}

.step-content {
  min-height: 300px;
}

.order-info-banner {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: linear-gradient(135deg, #F5F3FF 0%, #EEF2FF 100%);
  border-radius: 12px;
  margin-bottom: 20px;
}

.order-info-banner .el-icon {
  font-size: 24px;
  color: #6366F1;
}

.order-info-content {
  flex: 1;
}

.order-no {
  font-weight: 600;
  color: #1E1B4B;
  margin: 0;
  font-family: 'Fira Code', monospace;
  font-size: 14px;
}

.order-amount {
  font-size: 13px;
  color: #6B7280;
  margin: 4px 0 0 0;
}

.section-label {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  margin: 0 0 8px 0;
}

.refund-type-section {
  margin-bottom: 20px;
}

.refund-type-group {
  display: flex;
  gap: 12px;
  width: 100%;
}

.refund-type-group :deep(.el-radio-button) {
  flex: 1;
}

.refund-type-group :deep(.el-radio-button__inner) {
  width: 100%;
  border-radius: 12px !important;
  border: 1px solid #E5E7EB;
  padding: 16px;
}

.refund-type-group :deep(.el-radio-button.is-active .el-radio-button__inner) {
  border-color: #6366F1;
  background: #F5F3FF;
}

.type-option {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.type-option .el-icon {
  font-size: 24px;
  color: #6366F1;
}

.type-title {
  font-weight: 600;
  color: #1E1B4B;
}

.type-desc {
  font-size: 12px;
  color: #6B7280;
}

.amount-input-section {
  margin-bottom: 20px;
}

.amount-input-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #F9FAFB;
  border: 1px solid #E5E7EB;
  border-radius: 12px;
  padding: 12px 16px;
}

.currency-label {
  font-weight: 600;
  color: #6366F1;
  font-size: 16px;
}

.amount-input {
  flex: 1;
}

.amount-input :deep(.el-input__inner) {
  font-size: 24px;
  font-weight: 700;
  text-align: right;
  background: transparent;
  border: none;
  box-shadow: none;
}

.amount-hint {
  margin-top: 8px;
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 13px;
  color: #6B7280;
}

.amount-progress {
  flex: 1;
}

.reason-section,
.notes-section {
  margin-bottom: 20px;
}

/* Step 2 */
.confirm-alert {
  margin-bottom: 16px;
}

.confirm-summary {
  background: #F9FAFB;
  border-radius: 12px;
  padding: 16px;
}

.confirm-amount {
  font-size: 18px;
  font-weight: 700;
  color: #EF4444;
  font-family: 'Fira Sans', sans-serif;
}

.remaining-info {
  margin-top: 12px;
  padding: 12px;
  background: #FEF3C7;
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #92400E;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
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
@media (max-width: 576px) {
  .refund-type-group {
    flex-direction: column;
  }
}
</style>