<template>
  <el-dialog
    :model-value="visible"
    @update:model-value="$emit('update:visible', $event)"
    :title="isEdit ? '编辑规则' : '创建规则'"
    width="600px"
    destroy-on-close
  >
    <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
      <!-- Basic Info -->
      <el-form-item label="规则名称" prop="name">
        <el-input v-model="form.name" placeholder="例如: $10优惠券兑换" maxlength="100" />
      </el-form-item>

      <el-form-item label="描述">
        <el-input v-model="form.description" type="textarea" :rows="2" placeholder="规则描述" />
      </el-form-item>

      <!-- Coupon Selection -->
      <el-form-item label="关联优惠券" prop="coupon_id">
        <CouponSelector v-model="form.coupon_id" @change="handleCouponChange" />
      </el-form-item>

      <!-- Points Required -->
      <el-form-item label="所需积分" prop="points_required">
        <el-input-number v-model="form.points_required" :min="1" :max="999999" style="width: 200px" />
        <span class="form-hint">兑换所需积分数量</span>
      </el-form-item>

      <!-- Stock -->
      <el-form-item label="兑换总量">
        <el-input-number v-model="form.total_stock" :min="0" :max="999999" style="width: 200px" />
        <span class="form-hint">0 表示不限制</span>
      </el-form-item>

      <!-- Per User Limit -->
      <el-form-item label="每人限兑">
        <el-input-number v-model="form.per_user_limit" :min="0" :max="100" style="width: 200px" />
        <span class="form-hint">0 表示不限制</span>
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
    </el-form>

    <template #footer>
      <el-button @click="$emit('update:visible', false)">取消</el-button>
      <el-button type="primary" @click="handleSave" :loading="loading">
        保存
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import CouponSelector from './CouponSelector.vue'
import type { RedeemRule, CreateRedeemRuleParams } from '@/api/points'

const props = defineProps<{
  visible: boolean
  rule: RedeemRule | null
  loading?: boolean
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  submit: [data: CreateRedeemRuleParams]
}>()

const formRef = ref<FormInstance>()

const isEdit = computed(() => !!props.rule)

const form = reactive({
  name: '',
  description: '',
  coupon_id: 0,
  coupon_name: '',
  points_required: 500,
  total_stock: 100,
  per_user_limit: 5,
  dateRange: [] as string[]
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入规则名称', trigger: 'blur' }],
  coupon_id: [{ required: true, message: '请选择优惠券', trigger: 'change' }],
  points_required: [{ required: true, message: '请输入所需积分', trigger: 'blur' }]
}

// Watch rule prop to populate form
watch(() => props.rule, (rule) => {
  if (rule) {
    form.name = rule.name
    form.description = rule.description || ''
    form.coupon_id = rule.coupon_id
    form.coupon_name = rule.coupon_name
    form.points_required = rule.points_required
    form.total_stock = rule.total_stock
    form.per_user_limit = rule.per_user_limit
    form.dateRange = rule.start_at && rule.end_at ? [rule.start_at, rule.end_at] : []
  } else {
    // Reset form
    form.name = ''
    form.description = ''
    form.coupon_id = 0
    form.coupon_name = ''
    form.points_required = 500
    form.total_stock = 100
    form.per_user_limit = 5
    form.dateRange = []
  }
}, { immediate: true })

const handleCouponChange = (coupon: { id: number; name: string }) => {
  form.coupon_id = coupon.id
  form.coupon_name = coupon.name
}

const handleSave = async () => {
  if (!formRef.value) return

  await formRef.value.validate((valid) => {
    if (valid) {
      const data: CreateRedeemRuleParams = {
        name: form.name,
        description: form.description,
        coupon_id: form.coupon_id,
        points_required: form.points_required,
        total_stock: form.total_stock,
        per_user_limit: form.per_user_limit,
        status: isEdit.value ? 1 : 1 // active by default
      }

      if (form.dateRange && form.dateRange.length === 2) {
        data.start_at = form.dateRange[0]
        data.end_at = form.dateRange[1]
      }

      emit('submit', data)
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
</style>