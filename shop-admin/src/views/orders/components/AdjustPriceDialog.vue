<template>
  <el-dialog
    v-model="visible"
    title="Adjust Order Price"
    width="500px"
    :close-on-click-modal="false"
    destroy-on-close
  >
    <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
      <!-- Current Amount Display -->
      <div class="amount-display">
        <div class="amount-row">
          <span class="amount-label">Current Pay Amount</span>
          <span class="amount-value">
            <span class="currency">{{ currency }}</span>
            {{ formatAmount(currentAmount) }}
          </span>
        </div>
      </div>

      <!-- Adjustment Type -->
      <el-form-item label="Adjustment Type">
        <el-radio-group v-model="form.adjustType">
          <el-radio-button value="decrease">
            <el-icon><Minus /></el-icon>
            Decrease Price
          </el-radio-button>
          <el-radio-button value="increase">
            <el-icon><Plus /></el-icon>
            Increase Price
          </el-radio-button>
        </el-radio-group>
      </el-form-item>

      <!-- Adjustment Amount -->
      <el-form-item label="Amount" prop="adjustAmount">
        <div class="amount-input-wrapper">
          <span class="currency-label">{{ currency }}</span>
          <el-input-number
            v-model="form.adjustAmount"
            :min="0.01"
            :max="maxAdjustAmount / 100"
            :precision="2"
            :controls="false"
            class="amount-input"
            placeholder="Enter amount"
          />
        </div>
        <div v-if="form.adjustType === 'decrease'" class="amount-hint">
          <el-icon><InfoFilled /></el-icon>
          <span>Maximum: {{ currency }} {{ formatAmount(currentAmount) }}</span>
        </div>
      </el-form-item>

      <!-- Preview -->
      <div class="preview-section">
        <div class="preview-row">
          <span class="preview-label">Original Amount</span>
          <span class="preview-value">{{ currency }} {{ formatAmount(currentAmount) }}</span>
        </div>
        <div class="preview-row">
          <span class="preview-label">
            {{ form.adjustType === 'decrease' ? 'Decrease' : 'Increase' }}
          </span>
          <span class="preview-value" :class="form.adjustType">
            {{ form.adjustType === 'decrease' ? '-' : '+' }}{{ currency }}
            {{ form.adjustAmount ? form.adjustAmount.toFixed(2) : '0.00' }}
          </span>
        </div>
        <el-divider />
        <div class="preview-row total">
          <span class="preview-label">New Pay Amount</span>
          <span class="preview-value new-amount">{{ currency }} {{ formatAmount(newAmount) }}</span>
        </div>
      </div>

      <!-- Reason -->
      <el-form-item label="Reason" prop="reason">
        <el-input
          v-model="form.reason"
          type="textarea"
          :rows="3"
          placeholder="Enter the reason for price adjustment (required)"
          maxlength="200"
          show-word-limit
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="visible = false">Cancel</el-button>
        <el-button
          type="primary"
          :loading="submitting"
          :disabled="!canSubmit"
          @click="handleSubmit"
        >
          Confirm Adjustment
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Minus, Plus, InfoFilled } from '@element-plus/icons-vue'
import { adjustOrderPrice } from '@/api/order'

const props = defineProps<{
  modelValue: boolean
  orderId: string
  currentAmount: string
  currency: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  success: []
}>()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const formRef = ref()
const submitting = ref(false)

const form = reactive({
  adjustType: 'decrease' as 'decrease' | 'increase',
  adjustAmount: 0,
  reason: ''
})

const rules = {
  adjustAmount: [
    { required: true, message: 'Please enter amount', trigger: 'blur' },
    { type: 'number', min: 0.01, message: 'Amount must be greater than 0', trigger: 'blur' }
  ],
  reason: [
    { required: true, message: 'Please enter reason', trigger: 'blur' },
    { min: 5, max: 200, message: 'Reason must be 5-200 characters', trigger: 'blur' }
  ]
}

// Max adjustment amount (current amount for decrease)
const maxAdjustAmount = computed(() => {
  if (form.adjustType === 'decrease') {
    return parseFloat(props.currentAmount) * 100
  }
  return 999999999 // Large number for increase
})

// New amount after adjustment
const newAmount = computed(() => {
  const current = parseFloat(props.currentAmount) || 0
  const adjust = form.adjustAmount || 0
  if (form.adjustType === 'decrease') {
    return Math.max(0, current - adjust)
  }
  return current + adjust
})

// Can submit
const canSubmit = computed(() => {
  return form.adjustAmount > 0 && form.reason.length >= 5
})

// Format amount
const formatAmount = (amount: string | number) => {
  const num = typeof amount === 'string' ? parseFloat(amount) : amount
  return num.toFixed(2)
}

// Handle submit
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (!valid) return

    submitting.value = true
    try {
      const adjustAmountStr = form.adjustType === 'decrease'
        ? `-${form.adjustAmount.toFixed(2)}`
        : form.adjustAmount.toFixed(2)

      await adjustOrderPrice(props.orderId, {
        adjust_amount: adjustAmountStr,
        reason: form.reason
      })
      ElMessage.success('Price adjusted successfully')
      emit('success')
      visible.value = false
    } catch (error: any) {
      ElMessage.error(error?.message || 'Failed to adjust price')
    } finally {
      submitting.value = false
    }
  })
}

// Reset form when dialog opens
watch(visible, (val) => {
  if (val) {
    form.adjustType = 'decrease'
    form.adjustAmount = 0
    form.reason = ''
  }
})
</script>

<style scoped>
.amount-display {
  background: linear-gradient(135deg, #F5F3FF 0%, #EEF2FF 100%);
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 20px;
}

.amount-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.amount-label {
  font-weight: 500;
  color: #4B5563;
}

.amount-value {
  font-size: 24px;
  font-weight: 700;
  color: #1E1B4B;
  font-family: 'Fira Sans', sans-serif;
}

.currency {
  font-size: 14px;
  color: #6B7280;
  margin-right: 4px;
}

/* Radio buttons */
:deep(.el-radio-button__inner) {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 16px;
}

/* Amount Input */
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
  font-size: 20px;
  font-weight: 600;
  text-align: right;
  background: transparent;
  border: none;
  box-shadow: none;
}

.amount-hint {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 8px;
  font-size: 12px;
  color: #6B7280;
}

.amount-hint .el-icon {
  color: #6366F1;
}

/* Preview Section */
.preview-section {
  background: #F9FAFB;
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 20px;
}

.preview-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
}

.preview-label {
  color: #6B7280;
  font-size: 14px;
}

.preview-value {
  font-weight: 500;
  color: #1E1B4B;
}

.preview-value.decrease {
  color: #10B981;
}

.preview-value.increase {
  color: #EF4444;
}

.preview-value.new-amount {
  font-size: 18px;
  font-weight: 700;
  color: #EF4444;
}

.preview-row.total {
  padding-top: 12px;
}

.preview-row.total .preview-label {
  font-weight: 600;
  color: #1E1B4B;
}

/* Dialog Footer */
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>