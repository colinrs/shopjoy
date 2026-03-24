<template>
  <el-dialog
    :model-value="visible"
    @update:model-value="$emit('update:visible', $event)"
    title="调整积分"
    width="500px"
    destroy-on-close
  >
    <div class="adjust-dialog">
      <!-- Current Balance -->
      <div class="current-balance">
        <span class="label">当前余额:</span>
        <span class="value">{{ account?.balance?.toLocaleString() || 0 }} 积分</span>
      </div>

      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <!-- Adjustment Type -->
        <el-form-item label="调整类型" prop="adjustment_type">
          <el-radio-group v-model="form.adjustment_type">
            <el-radio value="ADD">
              <span class="radio-label add">增加积分</span>
            </el-radio>
            <el-radio value="DEDUCT">
              <span class="radio-label deduct">扣减积分</span>
            </el-radio>
          </el-radio-group>
        </el-form-item>

        <!-- Points Amount -->
        <el-form-item label="积分数量" prop="points">
          <el-input-number
            v-model="form.points"
            :min="1"
            :max="form.adjustment_type === 'DEDUCT' ? (account?.balance || 0) : 999999"
            :step="10"
            style="width: 100%"
          />
        </el-form-item>

        <!-- Reason -->
        <el-form-item label="原因" prop="reason">
          <el-input
            v-model="form.reason"
            type="textarea"
            :rows="3"
            placeholder="请输入调整原因（必填）"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>

        <!-- Preview -->
        <div class="preview-box">
          <div class="preview-label">调整后余额:</div>
          <div class="preview-value" :class="{ negative: previewBalance < 0 }">
            {{ previewBalance.toLocaleString() }} 积分
          </div>
        </div>
      </el-form>
    </div>

    <template #footer>
      <el-button @click="$emit('update:visible', false)">取消</el-button>
      <el-button type="primary" @click="handleSubmit" :loading="loading">
        确认调整
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import type { PointsAccount } from '@/api/points'

const props = defineProps<{
  visible: boolean
  account: PointsAccount | null
  loading?: boolean
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  submit: [data: { adjustment_type: 'ADD' | 'DEDUCT'; points: number; reason: string }]
}>()

const formRef = ref<FormInstance>()

const form = reactive({
  adjustment_type: 'ADD' as 'ADD' | 'DEDUCT',
  points: 100,
  reason: ''
})

const rules: FormRules = {
  adjustment_type: [
    { required: true, message: '请选择调整类型', trigger: 'change' }
  ],
  points: [
    { required: true, message: '请输入积分数量', trigger: 'blur' },
    { type: 'number', min: 1, message: '积分数量必须大于0', trigger: 'blur' }
  ],
  reason: [
    { required: true, message: '请输入调整原因', trigger: 'blur' },
    { min: 5, message: '原因至少5个字符', trigger: 'blur' }
  ]
}

const previewBalance = computed(() => {
  const current = props.account?.balance || 0
  const points = form.points || 0
  if (form.adjustment_type === 'ADD') {
    return current + points
  } else {
    return current - points
  }
})

// Reset form when dialog opens
watch(() => props.visible, (val) => {
  if (val) {
    form.adjustment_type = 'ADD'
    form.points = 100
    form.reason = ''
  }
})

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate((valid) => {
    if (valid) {
      emit('submit', {
        adjustment_type: form.adjustment_type,
        points: form.points,
        reason: form.reason
      })
    }
  })
}
</script>

<style scoped>
.adjust-dialog {
  padding: 0;
}

.current-balance {
  background: #F9FAFB;
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.current-balance .label {
  font-size: 14px;
  color: #6B7280;
}

.current-balance .value {
  font-size: 20px;
  font-weight: 700;
  color: #6366F1;
  font-family: 'Fira Sans', sans-serif;
}

.radio-label.add {
  color: #10B981;
}

.radio-label.deduct {
  color: #EF4444;
}

.preview-box {
  background: linear-gradient(135deg, #F5F3FF 0%, #EEF2FF 100%);
  border-radius: 12px;
  padding: 16px;
  margin-top: 16px;
  text-align: center;
}

.preview-label {
  font-size: 13px;
  color: #6B7280;
  margin-bottom: 4px;
}

.preview-value {
  font-size: 24px;
  font-weight: 700;
  color: #6366F1;
  font-family: 'Fira Sans', sans-serif;
}

.preview-value.negative {
  color: #EF4444;
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
}

:deep(.el-dialog__footer) {
  border-top: 1px solid #F3F4F6;
  padding: 16px 20px;
}
</style>