<template>
  <el-dialog
    :model-value="visible"
    :title="isEdit ? $t('points.editRule') : $t('points.createRule')"
    width="700px"
    destroy-on-close
    @update:model-value="$emit('update:visible', $event)"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="120px"
    >
      <!-- Basic Info -->
      <el-form-item
        :label="$t('points.ruleName')"
        prop="name"
      >
        <el-input
          v-model="form.name"
          :placeholder="$t('points.orderPoints')"
          maxlength="100"
        />
      </el-form-item>

      <el-form-item :label="$t('points.description')">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="2"
          :placeholder="$t('points.ruleDescPlaceholder')"
        />
      </el-form-item>

      <!-- Scenario -->
      <el-form-item
        :label="$t('points.scenario')"
        prop="scenario"
      >
        <el-radio-group v-model="form.scenario">
          <el-radio value="ORDER_PAYMENT">
            {{ $t('points.orderPayment') }}
          </el-radio>
          <el-radio value="SIGN_IN">
            {{ $t('points.signIn') }}
          </el-radio>
          <el-radio value="PRODUCT_REVIEW">
            {{ $t('points.productReview') }}
          </el-radio>
          <el-radio value="FIRST_ORDER">
            {{ $t('points.firstOrder') }}
          </el-radio>
        </el-radio-group>
      </el-form-item>

      <!-- Calculation Type -->
      <el-form-item
        :label="$t('points.calculationType')"
        prop="calculation_type"
      >
        <el-radio-group v-model="form.calculation_type">
          <el-radio value="FIXED">
            {{ $t('points.fixed') }}
          </el-radio>
          <el-radio value="RATIO">
            {{ $t('points.ratio') }}
          </el-radio>
          <el-radio value="TIERED">
            {{ $t('points.tiered') }}
          </el-radio>
        </el-radio-group>
      </el-form-item>

      <!-- Fixed Points -->
      <el-form-item
        v-if="form.calculation_type === 'FIXED'"
        :label="$t('points.fixedPoints')"
        prop="fixed_points"
      >
        <el-input-number
          v-model="form.fixed_points"
          :min="1"
          :max="99999"
          style="width: 200px"
        />
        <span class="form-hint">{{ $t('points.orderPointsReward') }}</span>
      </el-form-item>

      <!-- Ratio -->
      <el-form-item
        v-if="form.calculation_type === 'RATIO'"
        :label="$t('points.ratioPoints')"
        prop="ratio"
      >
        <el-input-number
          v-model="form.ratio_num"
          :min="0.1"
          :max="100"
          :precision="2"
          :step="0.5"
          style="width: 200px"
        />
        <span class="form-hint">{{ $t('points.pointsPerDollar') }}</span>
      </el-form-item>

      <!-- Tiered Config -->
      <el-form-item
        v-if="form.calculation_type === 'TIERED'"
        :label="$t('points.tieredConfig')"
      >
        <TieredConfig v-model="form.tiers" />
      </el-form-item>

      <!-- Condition Type -->
      <el-form-item :label="$t('points.triggerCondition')">
        <el-select
          v-model="form.condition_type"
          :placeholder="$t('points.condition')"
        >
          <el-option
            :label="$t('points.unconditional')"
            value="NONE"
          />
          <el-option
            :label="$t('points.newUser')"
            value="NEW_USER"
          />
          <el-option
            :label="$t('points.firstPurchase')"
            value="FIRST_ORDER"
          />
          <el-option
            :label="$t('points.specificProducts')"
            value="SPECIFIC_PRODUCTS"
          />
          <el-option
            :label="$t('points.minAmount')"
            value="MIN_AMOUNT"
          />
        </el-select>
      </el-form-item>

      <!-- Condition Value for MIN_AMOUNT -->
      <el-form-item
        v-if="form.condition_type === 'MIN_AMOUNT'"
        :label="$t('points.minAmount')"
      >
        <el-input-number
          v-model="form.min_amount"
          :min="1"
          :precision="2"
          style="width: 200px"
        />
        <span class="form-hint">{{ $t('points.amountDollar') }}</span>
      </el-form-item>

      <!-- Expiration -->
      <el-form-item :label="$t('points.pointsExpiration')">
        <el-input-number
          v-model="form.expiration_months"
          :min="0"
          :max="36"
          style="width: 200px"
        />
        <span class="form-hint">{{ $t('points.zeroNeverExpires') }}</span>
      </el-form-item>

      <!-- Time Range -->
      <el-form-item :label="$t('points.validPeriod')">
        <el-date-picker
          v-model="form.dateRange"
          type="datetimerange"
          :start-placeholder="$t('points.startTime')"
          :end-placeholder="$t('points.endTime')"
          value-format="YYYY-MM-DDTHH:mm:ss[Z]"
          style="width: 100%"
        />
      </el-form-item>

      <!-- Priority -->
      <el-form-item :label="$t('points.priority')">
        <el-input-number
          v-model="form.priority"
          :min="0"
          :max="100"
          style="width: 200px"
        />
        <span class="form-hint">{{ $t('points.greaterPriority') }}</span>
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="$emit('update:visible', false)">
        {{ $t('common.cancel') }}
      </el-button>
      <el-button
        :loading="loading"
        @click="handleSaveDraft"
      >
        {{ $t('points.saveDraft') }}
      </el-button>
      <el-button
        type="primary"
        :loading="loading"
        @click="handleSaveAndActivate"
      >
        {{ $t('points.saveAndActivate') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import TieredConfig from './TieredConfig.vue'
import type { EarnRule, TierConfig, CreateEarnRuleParams } from '@/api/points'
import { t } from '@/plugins/i18n'

const props = defineProps<{
  visible: boolean
  rule: EarnRule | null
  loading?: boolean
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  submit: [data: CreateEarnRuleParams]
}>()

const formRef = ref<FormInstance>()

const isEdit = computed(() => !!props.rule)

const form = reactive({
  name: '',
  description: '',
  scenario: 'ORDER_PAYMENT' as string,
  calculation_type: 'FIXED' as string,
  fixed_points: 10,
  ratio_num: 1,
  tiers: [{ threshold: 10000, ratio: '1.0' }, { threshold: null, ratio: '1.5' }] as TierConfig[],
  condition_type: 'NONE' as string,
  min_amount: 0,
  expiration_months: 12,
  dateRange: [] as string[],
  priority: 0
})

const rules: FormRules = {
  name: [{ required: true, message: t('points.ruleName'), trigger: 'blur' }],
  scenario: [{ required: true, message: t('points.scenario'), trigger: 'change' }],
  calculation_type: [{ required: true, message: t('points.calculationType'), trigger: 'change' }]
}

// Watch rule prop to populate form
watch(() => props.rule, (rule) => {
  if (rule) {
    form.name = rule.name
    form.description = rule.description || ''
    form.scenario = rule.scenario
    form.calculation_type = rule.calculation_type
    form.fixed_points = rule.fixed_points || 10
    form.ratio_num = parseFloat(rule.ratio) || 1
    form.tiers = rule.tiers || [{ threshold: 10000, ratio: '1.0' }, { threshold: null, ratio: '1.5' }]
    form.condition_type = rule.condition_type || 'NONE'
    form.min_amount = rule.condition_value?.min_amount || 0
    form.expiration_months = rule.expiration_months
    form.dateRange = rule.start_at && rule.end_at ? [rule.start_at, rule.end_at] : []
    form.priority = rule.priority
  } else {
    // Reset form
    form.name = ''
    form.description = ''
    form.scenario = 'ORDER_PAYMENT'
    form.calculation_type = 'FIXED'
    form.fixed_points = 10
    form.ratio_num = 1
    form.tiers = [{ threshold: 10000, ratio: '1.0' }, { threshold: null, ratio: '1.5' }]
    form.condition_type = 'NONE'
    form.min_amount = 0
    form.expiration_months = 12
    form.dateRange = []
    form.priority = 0
  }
}, { immediate: true })

const getFormData = (status: 'draft' | 'active'): CreateEarnRuleParams => {
  const data: CreateEarnRuleParams = {
    name: form.name,
    description: form.description,
    scenario: form.scenario,
    calculation_type: form.calculation_type,
    condition_type: form.condition_type,
    expiration_months: form.expiration_months,
    priority: form.priority,
    status
  }

  if (form.calculation_type === 'FIXED') {
    data.fixed_points = form.fixed_points
  } else if (form.calculation_type === 'RATIO') {
    data.ratio = String(form.ratio_num)
  } else if (form.calculation_type === 'TIERED') {
    data.tiers = form.tiers
  }

  if (form.condition_type === 'MIN_AMOUNT' && form.min_amount > 0) {
    data.condition_value = { min_amount: form.min_amount }
  }

  if (form.dateRange && form.dateRange.length === 2) {
    data.start_at = form.dateRange[0]
    data.end_at = form.dateRange[1]
  }

  return data
}

const handleSaveDraft = async () => {
  if (!formRef.value) return
  await formRef.value.validate((valid) => {
    if (valid) {
      emit('submit', getFormData('draft'))
    }
  })
}

const handleSaveAndActivate = async () => {
  if (!formRef.value) return
  await formRef.value.validate((valid) => {
    if (valid) {
      emit('submit', getFormData('active'))
    }
  })
}
</script>

<style scoped>
.form-hint {
  margin-left: 12px;
  font-size: 12px;
  color: #9CA3AF;
}

:deep(.el-dialog) {
  border-radius: 16px;
}

:deep(.el-dialog__header) {
  border-bottom: 1px solid #F3F4F6;
  padding: 16px 20px;
}

:deep(.el-dialog__title) {
  font-weight: 600;
  color: #1E1B4B;
}

:deep(.el-dialog__body) {
  padding: 20px;
  max-height: 60vh;
  overflow-y: auto;
}

:deep(.el-dialog__footer) {
  border-top: 1px solid #F3F4F6;
  padding: 16px 20px;
}

:deep(.el-radio-group) {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

:deep(.el-radio) {
  margin-right: 0;
}
</style>