<template>
  <el-dialog
    :model-value="visible"
    @update:model-value="$emit('update:visible', $event)"
    :title="isEdit ? '编辑规则' : '创建规则'"
    width="700px"
    destroy-on-close
  >
    <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
      <!-- Basic Info -->
      <el-form-item label="规则名称" prop="name">
        <el-input v-model="form.name" placeholder="例如: 订单积分奖励" maxlength="100" />
      </el-form-item>

      <el-form-item label="描述">
        <el-input v-model="form.description" type="textarea" :rows="2" placeholder="规则描述" />
      </el-form-item>

      <!-- Scenario -->
      <el-form-item label="获取场景" prop="scenario">
        <el-radio-group v-model="form.scenario">
          <el-radio value="ORDER_PAYMENT">订单支付</el-radio>
          <el-radio value="SIGN_IN">每日签到</el-radio>
          <el-radio value="PRODUCT_REVIEW">商品评价</el-radio>
          <el-radio value="FIRST_ORDER">首单奖励</el-radio>
        </el-radio-group>
      </el-form-item>

      <!-- Calculation Type -->
      <el-form-item label="计算类型" prop="calculation_type">
        <el-radio-group v-model="form.calculation_type">
          <el-radio value="FIXED">固定积分</el-radio>
          <el-radio value="RATIO">比例计算</el-radio>
          <el-radio value="TIERED">阶梯计算</el-radio>
        </el-radio-group>
      </el-form-item>

      <!-- Fixed Points -->
      <el-form-item v-if="form.calculation_type === 'FIXED'" label="固定积分" prop="fixed_points">
        <el-input-number v-model="form.fixed_points" :min="1" :max="99999" style="width: 200px" />
        <span class="form-hint">每次获得固定积分数量</span>
      </el-form-item>

      <!-- Ratio -->
      <el-form-item v-if="form.calculation_type === 'RATIO'" label="积分比例" prop="ratio">
        <el-input-number
          v-model="form.ratio_num"
          :min="0.1"
          :max="100"
          :precision="2"
          :step="0.5"
          style="width: 200px"
        />
        <span class="form-hint">每 $1 可获得积分数量</span>
      </el-form-item>

      <!-- Tiered Config -->
      <el-form-item v-if="form.calculation_type === 'TIERED'" label="阶梯配置">
        <TieredConfig v-model="form.tiers" />
      </el-form-item>

      <!-- Condition Type -->
      <el-form-item label="触发条件">
        <el-select v-model="form.condition_type" placeholder="选择条件">
          <el-option label="无条件" value="NONE" />
          <el-option label="新用户" value="NEW_USER" />
          <el-option label="首单" value="FIRST_ORDER" />
          <el-option label="指定商品" value="SPECIFIC_PRODUCTS" />
          <el-option label="最低金额" value="MIN_AMOUNT" />
        </el-select>
      </el-form-item>

      <!-- Condition Value for MIN_AMOUNT -->
      <el-form-item v-if="form.condition_type === 'MIN_AMOUNT'" label="最低金额">
        <el-input-number v-model="form.min_amount" :min="1" :precision="2" style="width: 200px" />
        <span class="form-hint">订单最低金额（美元）</span>
      </el-form-item>

      <!-- Expiration -->
      <el-form-item label="积分有效期">
        <el-input-number v-model="form.expiration_months" :min="0" :max="36" style="width: 200px" />
        <span class="form-hint">0 表示永不过期</span>
      </el-form-item>

      <!-- Time Range -->
      <el-form-item label="有效期">
        <el-date-picker
          v-model="form.dateRange"
          type="datetimerange"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          value-format="YYYY-MM-DDTHH:mm:ss[Z]"
          style="width: 100%"
        />
      </el-form-item>

      <!-- Priority -->
      <el-form-item label="优先级">
        <el-input-number v-model="form.priority" :min="0" :max="100" style="width: 200px" />
        <span class="form-hint">数值越大优先级越高</span>
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="$emit('update:visible', false)">取消</el-button>
      <el-button @click="handleSaveDraft" :loading="loading">保存草稿</el-button>
      <el-button type="primary" @click="handleSaveAndActivate" :loading="loading">
        保存并激活
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import TieredConfig from './TieredConfig.vue'
import type { EarnRule, TierConfig, CreateEarnRuleParams } from '@/api/points'

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
  name: [{ required: true, message: '请输入规则名称', trigger: 'blur' }],
  scenario: [{ required: true, message: '请选择获取场景', trigger: 'change' }],
  calculation_type: [{ required: true, message: '请选择计算类型', trigger: 'change' }]
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

const getFormData = (status: number): CreateEarnRuleParams => {
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
      emit('submit', getFormData(0)) // status = 0 (draft)
    }
  })
}

const handleSaveAndActivate = async () => {
  if (!formRef.value) return
  await formRef.value.validate((valid) => {
    if (valid) {
      emit('submit', getFormData(1)) // status = 1 (active)
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