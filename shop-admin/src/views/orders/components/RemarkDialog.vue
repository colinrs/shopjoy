<template>
  <el-dialog
    v-model="visible"
    :title="$t('orders.editRemarkTitle')"
    width="500px"
    :close-on-click-modal="false"
    destroy-on-close
  >
    <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
      <!-- Current Remark Display -->
      <div class="current-remark-section">
        <p class="section-label">{{ $t('orders.currentRemark') }}</p>
        <div class="current-remark-content">
          <p v-if="currentRemark">{{ currentRemark }}</p>
          <p v-else class="no-remark">{{ $t('orders.noRemark') }}</p>
        </div>
      </div>

      <!-- New Remark Input -->
      <el-form-item :label="$t('orders.newRemark')" prop="remark">
        <el-input
          v-model="form.remark"
          type="textarea"
          :rows="4"
          :placeholder="$t('orders.quickTags')"
          maxlength="500"
          show-word-limit
        />
      </el-form-item>

      <!-- Quick Tags -->
      <div class="quick-tags-section">
        <p class="section-label">{{ $t('orders.quickTags') }}</p>
        <div class="quick-tags">
          <el-tag
            v-for="tag in quickTags"
            :key="tag"
            class="quick-tag"
            @click="appendTag(tag)"
          >
            {{ tag }}
          </el-tag>
        </div>
      </div>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="visible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          {{ $t('orders.saveRemark') }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { updateOrderRemark } from '@/api/order'
import { t } from '@/plugins/i18n'

const props = defineProps<{
  modelValue: boolean
  orderId: string
  currentRemark: string
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
  remark: ''
})

const rules = {
  remark: [
    { max: 500, message: t('orders.remarkTooLong'), trigger: 'blur' }
  ]
}

const quickTags = computed(() => [
  t('orders.vipCustomer'),
  t('orders.returnCustomer'),
  t('orders.largeOrder'),
  t('orders.urgent'),
  t('orders.giftPackage'),
  t('orders.contactBeforeShip')
])

const appendTag = (tag: string) => {
  if (form.remark) {
    form.remark += ', ' + tag
  } else {
    form.remark = tag
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (!valid) return

    submitting.value = true
    try {
      await updateOrderRemark(props.orderId, {
        remark: form.remark
      })
      ElMessage.success(t('orders.remarkSaveSuccess'))
      emit('success')
      visible.value = false
    } catch (error: any) {
      ElMessage.error(error?.message || t('orders.remarkSaveFailed'))
    } finally {
      submitting.value = false
    }
  })
}

// Reset form when dialog opens
watch(visible, (val) => {
  if (val) {
    form.remark = props.currentRemark || ''
  }
})
</script>

<style scoped>
.current-remark-section,
.quick-tags-section {
  margin-bottom: 20px;
}

.section-label {
  font-size: 13px;
  color: #6B7280;
  margin: 0 0 8px 0;
}

.current-remark-content {
  background: #F9FAFB;
  border-radius: 12px;
  padding: 12px 16px;
}

.current-remark-content p {
  margin: 0;
  color: #4B5563;
}

.no-remark {
  color: #9CA3AF;
  font-style: italic;
}

.quick-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.quick-tag {
  cursor: pointer;
  transition: all 0.2s ease;
}

.quick-tag:hover {
  background: #EEF2FF;
  border-color: #6366F1;
  color: #6366F1;
}

/* Dialog Footer */
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>