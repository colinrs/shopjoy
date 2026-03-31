<template>
  <el-dialog
    v-model="visible"
    title="Batch Ship"
    width="700px"
    :close-on-click-modal="false"
  >
    <div class="batch-ship-content">
      <!-- Selected Shipments -->
      <div class="selected-section">
        <div class="section-header">
          <el-icon><List /></el-icon>
          <span>Selected Shipments</span>
          <el-tag size="small" type="primary">{{ shipments.length }} items</el-tag>
        </div>

        <div class="shipments-list">
          <div v-for="shipment in shipments" :key="shipment.id" class="shipment-item">
            <div class="shipment-info">
              <span class="shipment-no">{{ shipment.shipment_no }}</span>
              <span class="order-no">Order: {{ shipment.order_no }}</span>
            </div>
            <div class="shipment-items">
              <span class="items-count">{{ shipment.items.length }} items</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Logistics Form -->
      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px" class="logistics-form">
        <div class="section-header">
          <el-icon><Van /></el-icon>
          <span>Logistics Settings</span>
        </div>

        <el-form-item label="Carrier" prop="carrier_code">
          <el-select
            v-model="form.carrier_code"
            placeholder="Select carrier"
            style="width: 100%"
          >
            <el-option
              v-for="carrier in carriers"
              :key="carrier.code"
              :label="carrier.name"
              :value="carrier.code"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="Tracking No." prop="tracking_no_start">
          <div class="tracking-input-group">
            <el-input
              v-model="form.tracking_no_start"
              placeholder="Enter starting tracking number"
            >
              <template #prefix>
                <el-icon><Tickets /></el-icon>
              </template>
            </el-input>
            <div class="tracking-hint">
              <el-icon><InfoFilled /></el-icon>
              <span>Tracking numbers will auto-increment for each shipment</span>
            </div>
          </div>
        </el-form-item>

        <el-form-item label="Preview">
          <div class="tracking-preview">
            <div v-for="(preview, index) in trackingPreviews" :key="index" class="preview-item">
              <span class="preview-shipment">{{ preview.shipment_no }}</span>
              <span class="preview-tracking">{{ preview.tracking_no }}</span>
            </div>
            <p v-if="trackingPreviews.length < shipments.length" class="more-preview">
              ... and {{ shipments.length - trackingPreviews.length }} more
            </p>
          </div>
        </el-form-item>
      </el-form>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="visible = false">Cancel</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          Confirm Batch Ship
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { List, Van, Tickets, InfoFilled } from '@element-plus/icons-vue'
import { batchCreateShipments, type Shipment, type Carrier } from '@/api/fulfillment'

const props = defineProps<{
  modelValue: boolean
  shipments: Shipment[]
  carriers: Carrier[]
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
  carrier_code: '',
  tracking_no_start: ''
})

const rules = {
  carrier_code: [{ required: true, message: 'Please select a carrier', trigger: 'change' }],
  tracking_no_start: [{ required: true, message: 'Please enter starting tracking number', trigger: 'blur' }]
}

const trackingPreviews = computed(() => {
  if (!form.tracking_no_start) return []

  return props.shipments.slice(0, 5).map((shipment, index) => {
    const baseTracking = form.tracking_no_start
    const numericPart = baseTracking.match(/\d+$/)?.[0] || ''
    const prefix = baseTracking.replace(/\d+$/, '')

    if (numericPart) {
      const num = parseInt(numericPart, 10) + index
      const paddedNum = num.toString().padStart(numericPart.length, '0')
      return {
        shipment_no: shipment.shipment_no,
        tracking_no: prefix + paddedNum
      }
    }

    return {
      shipment_no: shipment.shipment_no,
      tracking_no: baseTracking + (index > 0 ? `-${index}` : '')
    }
  })
})

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (!valid) return

    submitting.value = true
    try {
      // Generate tracking numbers for each shipment
      const shipments = props.shipments.map((shipment, index) => {
        let tracking_no = form.tracking_no_start
        if (form.tracking_no_start) {
          const numericPart = form.tracking_no_start.match(/\d+$/)?.[0] || ''
          const prefix = form.tracking_no_start.replace(/\d+$/, '')

          if (numericPart) {
            const num = parseInt(numericPart, 10) + index
            const paddedNum = num.toString().padStart(numericPart.length, '0')
            tracking_no = prefix + paddedNum
          } else if (index > 0) {
            tracking_no = `${form.tracking_no_start}-${index}`
          }
        }

        return {
          order_id: shipment.order_id,
          tracking_no
        }
      })

      await batchCreateShipments({
        carrier_code: form.carrier_code,
        shipments
      })
      ElMessage.success(`Successfully shipped ${props.shipments.length} orders`)
      emit('success')
      visible.value = false
    } catch (error) {
      ElMessage.error('Failed to batch ship')
    } finally {
      submitting.value = false
    }
  })
}
</script>

<style scoped>
.batch-ship-content {
  max-height: 60vh;
  overflow-y: auto;
}

.selected-section,
.logistics-form {
  margin-bottom: 24px;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  font-weight: 600;
  color: #1E1B4B;
  font-size: 15px;
}

.section-header .el-icon {
  color: #6366F1;
}

/* Shipments List */
.shipments-list {
  max-height: 200px;
  overflow-y: auto;
  border: 1px solid #E5E7EB;
  border-radius: 12px;
  padding: 8px;
}

.shipment-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  border-radius: 8px;
  transition: background 0.2s ease;
}

.shipment-item:hover {
  background: #F5F3FF;
}

.shipment-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.shipment-no {
  font-family: 'Fira Code', monospace;
  font-size: 14px;
  color: #6366F1;
  font-weight: 500;
}

.order-no {
  font-size: 12px;
  color: #6B7280;
}

.shipment-items {
  text-align: right;
}

.items-count {
  font-size: 13px;
  color: #1E1B4B;
  font-weight: 500;
}

/* Tracking Input */
.tracking-input-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 100%;
}

.tracking-hint {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #6B7280;
}

.tracking-hint .el-icon {
  color: #6366F1;
}

/* Tracking Preview */
.tracking-preview {
  background: #F9FAFB;
  border-radius: 12px;
  padding: 16px;
  width: 100%;
}

.preview-item {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid #E5E7EB;
}

.preview-item:last-of-type {
  border-bottom: none;
}

.preview-shipment {
  font-family: 'Fira Code', monospace;
  font-size: 13px;
  color: #6B7280;
}

.preview-tracking {
  font-family: 'Fira Code', monospace;
  font-size: 13px;
  color: #6366F1;
  font-weight: 500;
}

.more-preview {
  font-size: 12px;
  color: #9CA3AF;
  text-align: center;
  margin: 8px 0 0 0;
}

/* Dialog Footer */
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>