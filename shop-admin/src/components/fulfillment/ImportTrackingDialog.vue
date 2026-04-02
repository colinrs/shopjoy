<template>
  <el-dialog
    v-model="visible"
    :title="$t('fulfillment.importTracking')"
    width="600px"
    :close-on-click-modal="false"
  >
    <div class="import-content">
      <!-- Template Download -->
      <div class="template-section">
        <el-button
          type="primary"
          link
          @click="downloadTemplate"
        >
          <el-icon><Download /></el-icon>
          {{ $t('fulfillment.importTemplate') }}
        </el-button>
      </div>

      <!-- Upload Area -->
      <el-upload
        ref="uploadRef"
        class="csv-uploader"
        drag
        :auto-upload="false"
        :limit="1"
        :on-change="handleFileChange"
        :on-remove="handleFileRemove"
        :file-list="fileList"
        accept=".csv"
      >
        <el-icon class="upload-icon">
          <UploadFilled />
        </el-icon>
        <div class="upload-text">
          {{ $t('fulfillment.importTrackingTip') }}
        </div>
        <template #tip>
          <div class="upload-tip">
            CSV: shipment_id,carrier_code,tracking_no,weight
          </div>
        </template>
      </el-upload>

      <!-- Preview Table -->
      <div
        v-if="parsedData.length > 0"
        class="preview-section"
      >
        <div class="section-header">
          <el-icon><List /></el-icon>
          <span>{{ $t('common.preview') || 'Preview' }} ({{ parsedData.length }} items)</span>
        </div>
        <el-table
          :data="parsedData"
          size="small"
          max-height="200"
          border
        >
          <el-table-column
            prop="shipment_id"
            label="shipment_id"
            width="100"
          />
          <el-table-column
            prop="carrier_code"
            label="carrier_code"
            width="100"
          />
          <el-table-column
            prop="tracking_no"
            label="tracking_no"
            min-width="150"
          />
          <el-table-column
            prop="weight"
            label="weight"
            width="80"
          />
        </el-table>
      </div>

      <!-- Error Display -->
      <div
        v-if="errors.length > 0"
        class="error-section"
      >
        <el-alert
          type="error"
          :closable="false"
        >
          <template #title>
            {{ $t('fulfillment.importFailed', { failed: errors.length }) }}
          </template>
          <div class="error-list">
            <div
              v-for="(err, idx) in errors.slice(0, 5)"
              :key="idx"
              class="error-item"
            >
              Shipment ID {{ err.shipment_id }}: {{ err.message }}
            </div>
            <div
              v-if="errors.length > 5"
              class="error-more"
            >
              ... and {{ errors.length - 5 }} more errors
            </div>
          </div>
        </el-alert>
      </div>

      <!-- Summary -->
      <div
        v-if="parsedData.length > 0"
        class="summary-section"
      >
        <el-alert
          type="info"
          :closable="false"
        >
          <template #title>
            {{ $t('fulfillment.importSuccess', { success: parsedData.length }) }}
          </template>
        </el-alert>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="visible = false">
          {{ $t('common.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="submitting"
          :disabled="parsedData.length === 0"
          @click="handleSubmit"
        >
          {{ $t('common.confirm') }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage, type UploadFile } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { UploadFilled, Download, List } from '@element-plus/icons-vue'
import { batchUpdateTracking, type BatchUpdateTrackingRequest } from '@/api/fulfillment'

const { t } = useI18n()

const props = defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  success: []
}>()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const uploadRef = ref()
const fileList = ref<UploadFile[]>([])
const submitting = ref(false)
const parsedData = ref<BatchUpdateTrackingRequest[]>([])
const errors = ref<{ shipment_id: number; message: string }[]>([])

// CSV Template download
const downloadTemplate = () => {
  const template = 'shipment_id,carrier_code,tracking_no,weight\n123456,SF,SF1234567890,1.5\n123457,YTO,YTO9876543210,2.0'
  const blob = new Blob([template], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = 'tracking_import_template.csv'
  link.click()
  URL.revokeObjectURL(link.href)
}

// Parse CSV content
const parseCSV = (content: string): BatchUpdateTrackingRequest[] => {
  const lines = content.split('\n').filter(line => line.trim())
  if (lines.length < 2) {
    return []
  }

  const data: BatchUpdateTrackingRequest[] = []
  errors.value = []

  // Skip header
  for (let i = 1; i < lines.length; i++) {
    const parts = lines[i].split(',').map(p => p.trim())
    if (parts.length >= 3) {
      const shipmentId = parseInt(parts[0], 10)
      if (isNaN(shipmentId)) {
        errors.value.push({ shipment_id: 0, message: `Invalid shipment_id: ${parts[0]}` })
        continue
      }
      data.push({
        shipment_ids: [shipmentId],
        carrier_code: parts[1],
        tracking_no: parts[2],
        weight: parts[3] || undefined
      })
    }
  }
  return data
}

// Handle file selection
const handleFileChange = (file: UploadFile) => {
  errors.value = []
  parsedData.value = []

  if (!file.raw) {
    ElMessage.error(t('fulfillment.importEmptyError'))
    return
  }

  const reader = new FileReader()
  reader.onload = (e) => {
    const content = e.target?.result as string
    if (!content) {
      ElMessage.error(t('fulfillment.importEmptyError'))
      return
    }
    parsedData.value = parseCSV(content)
    if (parsedData.value.length === 0 && errors.value.length === 0) {
      ElMessage.error(t('fulfillment.importEmptyError'))
    }
  }
  reader.readAsText(file.raw)
}

// Handle file removal
const handleFileRemove = () => {
  parsedData.value = []
  errors.value = []
}

// Handle submit - process each row individually
const handleSubmit = async () => {
  if (parsedData.value.length === 0) {
    ElMessage.error(t('fulfillment.importEmptyError'))
    return
  }

  submitting.value = true
  let successCount = 0
  let failedCount = 0
  const failedErrors: { shipment_id: number; message: string }[] = []

  try {
    // Process each shipment individually
    for (const item of parsedData.value) {
      try {
        const result = await batchUpdateTracking(item)
        if (result.failed && result.failed.length > 0) {
          failedCount++
          result.failed.forEach(f => {
            failedErrors.push({ shipment_id: f.shipment_id, message: f.message })
          })
        } else {
          successCount++
        }
      } catch (error) {
        failedCount++
        failedErrors.push({ shipment_id: item.shipment_ids[0], message: 'Network error' })
      }
    }

    if (failedCount > 0) {
      errors.value = failedErrors
      ElMessage.error(t('fulfillment.importFailed', { failed: failedCount }))
    }

    if (successCount > 0) {
      ElMessage.success(t('fulfillment.importSuccess', { success: successCount }))
      emit('success')
      visible.value = false
      // Reset state
      fileList.value = []
      parsedData.value = []
      errors.value = []
    }
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.import-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.template-section {
  display: flex;
  align-items: center;
}

.csv-uploader {
  width: 100%;
}

.upload-icon {
  font-size: 48px;
  color: #6366F1;
  margin-bottom: 16px;
}

.upload-text {
  font-size: 14px;
  color: #6B7280;
}

.upload-tip {
  font-size: 12px;
  color: #9CA3AF;
  margin-top: 8px;
}

.preview-section {
  margin-top: 8px;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
  font-weight: 600;
  color: #1E1B4B;
  font-size: 14px;
}

.section-header .el-icon {
  color: #6366F1;
}

.error-section {
  margin-top: 8px;
}

.error-list {
  max-height: 120px;
  overflow-y: auto;
}

.error-item {
  font-size: 12px;
  padding: 4px 0;
  font-family: 'Fira Code', monospace;
}

.error-more {
  font-size: 12px;
  color: #9CA3AF;
  margin-top: 4px;
}

.summary-section {
  margin-top: 8px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
