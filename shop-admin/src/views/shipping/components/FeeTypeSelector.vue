<template>
  <div class="fee-type-selector">
    <!-- Fee Type Radio -->
    <el-radio-group v-model="localFeeType" class="fee-type-group" @change="handleFeeTypeChange">
      <el-radio-button value="fixed">
        <el-icon><Coin /></el-icon>
        {{ $t('shipping.fixed') }}
      </el-radio-button>
      <el-radio-button value="by_count">
        <el-icon><Box /></el-icon>
        {{ $t('shipping.byCount') }}
      </el-radio-button>
      <el-radio-button value="by_weight">
        <el-icon><Odometer /></el-icon>
        {{ $t('shipping.byWeight') }}
      </el-radio-button>
      <el-radio-button value="free">
        <el-icon><Present /></el-icon>
        {{ $t('shipping.free') }}
      </el-radio-button>
    </el-radio-group>

    <!-- Fee Configuration -->
    <div class="fee-config" v-if="localFeeType !== 'free'">
      <!-- Fixed Fee -->
      <template v-if="localFeeType === 'fixed'">
        <el-form-item :label="$t('shipping.shippingFee')">
          <el-input-number
            v-model="localForm.first_fee"
            :min="0"
            :precision="2"
            style="width: 200px"
          >
            <template #prefix>¥</template>
          </el-input-number>
        </el-form-item>
      </template>

      <!-- Per Item -->
      <template v-if="localFeeType === 'by_count'">
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item :label="$t('shipping.firstItemCount')">
              <el-input-number
                v-model="localForm.first_unit"
                :min="1"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('shipping.firstItemFee')">
              <el-input-number
                v-model="localForm.first_fee"
                :min="0"
                :precision="2"
                style="width: 100%"
              >
                <template #prefix>¥</template>
              </el-input-number>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item :label="$t('shipping.additionalItemCount')">
              <el-input-number
                v-model="localForm.additional_unit"
                :min="1"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('shipping.additionalItemFee')">
              <el-input-number
                v-model="localForm.additional_fee"
                :min="0"
                :precision="2"
                style="width: 100%"
              >
                <template #prefix>¥</template>
              </el-input-number>
            </el-form-item>
          </el-col>
        </el-row>
        <div class="fee-preview">
          <span class="preview-text">
            {{ $t('shipping.billingExampleCount', { first: localForm.first_unit, fee: localForm.first_fee, add: localForm.additional_unit, addFee: localForm.additional_fee }) }}
          </span>
        </div>
      </template>

      <!-- By Weight -->
      <template v-if="localFeeType === 'by_weight'">
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item :label="$t('shipping.firstWeightGrams')">
              <el-input-number
                v-model="localForm.first_unit"
                :min="1"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('shipping.firstWeightFee')">
              <el-input-number
                v-model="localForm.first_fee"
                :min="0"
                :precision="2"
                style="width: 100%"
              >
                <template #prefix>¥</template>
              </el-input-number>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item :label="$t('shipping.additionalWeightGrams')">
              <el-input-number
                v-model="localForm.additional_unit"
                :min="1"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('shipping.additionalWeightFee')">
              <el-input-number
                v-model="localForm.additional_fee"
                :min="0"
                :precision="2"
                style="width: 100%"
              >
                <template #prefix>¥</template>
              </el-input-number>
            </el-form-item>
          </el-col>
        </el-row>
        <div class="fee-preview">
          <span class="preview-text">
            {{ $t('shipping.billingExampleWeight', { first: localForm.first_unit, fee: localForm.first_fee, add: localForm.additional_unit, addFee: localForm.additional_fee }) }}
          </span>
        </div>
      </template>
    </div>

    <!-- Free Shipping Note -->
    <el-alert
      v-if="localFeeType === 'free'"
      type="success"
      :closable="false"
      show-icon
    >
      {{ $t('shipping.freeShippingNote') }}
    </el-alert>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { Coin, Box, Odometer, Present } from '@element-plus/icons-vue'
import type { CreateZoneRequest } from '@/api/shipping'

const props = defineProps<{
  modelValue: Partial<CreateZoneRequest>
}>()

const emit = defineEmits<{
  'update:modelValue': [value: Partial<CreateZoneRequest>]
}>()

// Local state
const localFeeType = ref(props.modelValue.fee_type || 'fixed')
const localForm = ref({
  first_unit: props.modelValue.first_unit || 1,
  first_fee: props.modelValue.first_fee || '0',
  additional_unit: props.modelValue.additional_unit || 1,
  additional_fee: props.modelValue.additional_fee || '0'
})

// Methods
const handleFeeTypeChange = (_value: string) => {
  emitUpdate()
}

const emitUpdate = () => {
  emit('update:modelValue', {
    ...props.modelValue,
    fee_type: localFeeType.value as any,
    first_unit: localForm.value.first_unit,
    first_fee: String(localForm.value.first_fee),
    additional_unit: localForm.value.additional_unit,
    additional_fee: String(localForm.value.additional_fee)
  })
}

// Watch for external changes
watch(() => props.modelValue, (newVal) => {
  localFeeType.value = newVal.fee_type || 'fixed'
  localForm.value = {
    first_unit: newVal.first_unit || 1,
    first_fee: newVal.first_fee || '0',
    additional_unit: newVal.additional_unit || 1,
    additional_fee: newVal.additional_fee || '0'
  }
}, { immediate: true })

// Watch local changes
watch([localFeeType, localForm], () => {
  emitUpdate()
}, { deep: true })
</script>

<style scoped>
.fee-type-selector {
  width: 100%;
}

.fee-type-group {
  display: flex;
  margin-bottom: 20px;
}

.fee-type-group :deep(.el-radio-button) {
  flex: 1;
}

.fee-type-group :deep(.el-radio-button__inner) {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 12px 16px;
  border-radius: 0;
}

.fee-type-group :deep(.el-radio-button:first-child .el-radio-button__inner) {
  border-radius: 8px 0 0 8px;
}

.fee-type-group :deep(.el-radio-button:last-child .el-radio-button__inner) {
  border-radius: 0 8px 8px 0;
}

.fee-type-group :deep(.el-radio-button__original-radio:checked + .el-radio-button__inner) {
  background-color: #6366F1;
  border-color: #6366F1;
}

.fee-config {
  padding: 16px;
  background: #F9FAFB;
  border-radius: 12px;
  margin-top: 16px;
}

.fee-preview {
  margin-top: 16px;
  padding: 12px 16px;
  background: rgba(99, 102, 241, 0.05);
  border-radius: 8px;
  border-left: 3px solid #6366F1;
}

.preview-text {
  font-size: 13px;
  color: #374151;
}

:deep(.el-alert--success) {
  margin-top: 16px;
  border-radius: 12px;
}

/* Responsive */
@media (max-width: 768px) {
  .fee-type-group {
    flex-direction: column;
  }

  .fee-type-group :deep(.el-radio-button__inner) {
    border-radius: 0;
  }

  .fee-type-group :deep(.el-radio-button:first-child .el-radio-button__inner) {
    border-radius: 8px 8px 0 0;
  }

  .fee-type-group :deep(.el-radio-button:last-child .el-radio-button__inner) {
    border-radius: 0 0 8px 8px;
  }
}
</style>
