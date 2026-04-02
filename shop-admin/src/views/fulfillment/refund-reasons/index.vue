<template>
  <div class="refund-reasons-page">
    <!-- Header -->
    <el-card
      class="header-card"
      shadow="never"
    >
      <div class="header-bar">
        <div class="header-left">
          <h2>{{ $t('refundReasons.title') }}</h2>
        </div>
        <div class="header-right">
          <el-button
            type="primary"
            @click="handleAdd"
          >
            <el-icon><Plus /></el-icon>
            {{ $t('refundReasons.addReason') }}
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Reason List -->
    <el-card
      class="list-card"
      shadow="never"
    >
      <el-table
        v-loading="loading"
        :data="reasons"
        row-key="id"
        border
        stripe
      >
        <el-table-column
          prop="sort"
          :label="$t('refundReasons.sort')"
          width="80"
          align="center"
        />
        <el-table-column
          prop="code"
          :label="$t('refundReasons.code')"
          min-width="150"
        />
        <el-table-column
          prop="name"
          :label="$t('refundReasons.name')"
          min-width="200"
        />
        <el-table-column
          :label="$t('refundReasons.status')"
          width="100"
          align="center"
        >
          <template #default="{ row }">
            <el-switch
              v-model="row.is_active"
              :active-value="true"
              :inactive-value="false"
              @change="(val: boolean) => handleStatusChange(row, val)"
            />
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('common.actions')"
          width="180"
          align="center"
        >
          <template #default="{ row }">
            <el-button
              type="primary"
              link
              size="small"
              @click="handleEdit(row)"
            >
              {{ $t('common.edit') }}
            </el-button>
            <el-button
              type="danger"
              link
              size="small"
              @click="handleDelete(row)"
            >
              {{ $t('common.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty
        v-if="!loading && reasons.length === 0"
        :description="$t('common.noData')"
      />
    </el-card>

    <!-- Add/Edit Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? $t('refundReasons.editReason') : $t('refundReasons.addReason')"
      width="450px"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="reasonForm"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item
          :label="$t('refundReasons.code')"
          prop="code"
        >
          <el-input
            v-model="reasonForm.code"
            :placeholder="$t('refundReasons.enterCode')"
          />
        </el-form-item>
        <el-form-item
          :label="$t('refundReasons.name')"
          prop="name"
        >
          <el-input
            v-model="reasonForm.name"
            :placeholder="$t('refundReasons.enterName')"
          />
        </el-form-item>
        <el-form-item :label="$t('refundReasons.sort')">
          <el-input-number
            v-model="reasonForm.sort"
            :min="0"
            :max="9999"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">
          {{ $t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="saveLoading"
          @click="handleSave"
        >
          {{ $t('common.save') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import {
  getRefundReasonList,
  type RefundReason
} from '@/api/fulfillment'
import { useErrorHandler } from '@/composables/useErrorHandler'

const { t } = useI18n()
const { handleError } = useErrorHandler()

const loading = ref(false)
const saveLoading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const reasons = ref<RefundReason[]>([])
const formRef = ref()

const reasonForm = reactive({
  id: 0,
  code: '',
  name: '',
  sort: 0
})

const formRules = {
  code: [{ required: true, message: t('refundReasons.enterCode'), trigger: 'blur' }],
  name: [{ required: true, message: t('refundReasons.enterName'), trigger: 'blur' }]
}

// Load refund reasons
const loadReasons = async () => {
  loading.value = true
  try {
    const data = await getRefundReasonList()
    reasons.value = data || []
  } catch (error) {
    handleError(error, t('refundReasons.loadFailed'))
  } finally {
    loading.value = false
  }
}

// Handle add
const handleAdd = () => {
  isEdit.value = false
  Object.assign(reasonForm, {
    id: 0,
    code: '',
    name: '',
    sort: reasons.value.length > 0 ? Math.max(...reasons.value.map(r => r.sort)) + 1 : 0
  })
  dialogVisible.value = true
}

// Handle edit
const handleEdit = (row: RefundReason) => {
  isEdit.value = true
  Object.assign(reasonForm, {
    id: row.id,
    code: row.code,
    name: row.name,
    sort: row.sort
  })
  dialogVisible.value = true
}

// Handle save
const handleSave = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      saveLoading.value = true
      try {
        // Note: Backend only has GET endpoint for refund reasons
        // If create/update/delete are needed, they would need to be added to the backend API
        // For now, we show a message that this is read-only
        ElMessage.info(t('refundReasons.readOnlyMessage'))
        dialogVisible.value = false
        // Reload to ensure data is fresh
        loadReasons()
      } catch (error) {
        handleError(error, t('refundReasons.saveFailed'))
      } finally {
        saveLoading.value = false
      }
    }
  })
}

// Handle delete
const handleDelete = async (row: RefundReason) => {
  try {
    await ElMessageBox.confirm(
      t('refundReasons.confirmDelete', { name: row.name }),
      t('refundReasons.deleteConfirm'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )
    // Note: Backend only has GET endpoint for refund reasons
    // If delete is needed, it would need to be added to the backend API
    ElMessage.info(t('refundReasons.readOnlyMessage'))
  } catch (error) {
    // User cancelled
  }
}

// Handle status change
const handleStatusChange = async (row: RefundReason, isActive: boolean) => {
  try {
    // Note: Backend only has GET endpoint for refund reasons
    // If status update is needed, it would need to be added to the backend API
    ElMessage.info(t('refundReasons.readOnlyMessage'))
    // Revert the change since backend doesn't support it
    row.is_active = !isActive
  } catch (error) {
    handleError(error, t('refundReasons.updateStatusFailed'))
    // Revert the change
    row.is_active = !isActive
  }
}

onMounted(() => {
  loadReasons()
})
</script>

<style scoped>
.refund-reasons-page {
  padding: 0;
}

.header-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

.header-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #1E1B4B;
}

.list-card {
  border-radius: 16px;
  border: 1px solid rgba(99, 102, 241, 0.06);
}

:deep(.el-table__row:hover > td) {
  background-color: #F5F3FF !important;
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

:deep(.el-dialog__footer) {
  border-top: 1px solid #F3F4F6;
  padding: 16px 20px;
}
</style>
