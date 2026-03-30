<template>
  <el-dialog
    v-model="visible"
    :title="$t('orders.batchShipOrders')"
    width="700px"
    :close-on-click-modal="false"
    destroy-on-close
  >
    <div class="batch-ship-content">
      <!-- Selected Orders -->
      <div class="selected-section">
        <div class="section-header">
          <el-icon><List /></el-icon>
          <span>{{ $t('orders.selectedOrders') }}</span>
          <el-tag size="small" type="primary">{{ $t('orders.itemsCount', { count: orders.length }) }}</el-tag>
        </div>

        <div class="orders-list">
          <div v-for="order in orders" :key="order.order_id" class="order-item">
            <div class="order-info">
              <span class="order-no">{{ order.order_no }}</span>
              <span class="order-buyer">{{ order.user_name }} - {{ order.user_phone }}</span>
            </div>
            <div class="order-items">
              <span class="items-count">{{ $t('orders.itemsCount', { count: order.item_count }) }}</span>
              <span class="items-amount">{{ order.currency }} {{ formatAmount(order.pay_amount) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Logistics Form -->
      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px" class="logistics-form">
        <div class="section-header">
          <el-icon><Van /></el-icon>
          <span>{{ $t('orders.logisticsSettings') }}</span>
        </div>

        <el-form-item :label="$t('orders.carrier')" prop="carrier_code">
          <el-select
            v-model="form.carrier_code"
            :placeholder="$t('orders.selectCarrier')"
            style="width: 100%"
            :loading="loadingCarriers"
          >
            <el-option
              v-for="carrier in carriers"
              :key="carrier.code"
              :label="carrier.name"
              :value="carrier.code"
            />
          </el-select>
        </el-form-item>

        <el-form-item :label="$t('orders.trackingNoStart')" prop="tracking_no_start">
          <div class="tracking-input-group">
            <el-input
              v-model="form.tracking_no_start"
              :placeholder="$t('orders.enterTrackingNoStart')"
            >
              <template #prefix>
                <el-icon><Tickets /></el-icon>
              </template>
            </el-input>
            <div class="tracking-hint">
              <el-icon><InfoFilled /></el-icon>
              <span>{{ $t('orders.trackingAutoIncrement') }}</span>
            </div>
          </div>
        </el-form-item>

        <el-form-item :label="$t('orders.preview')">
          <div class="tracking-preview">
            <div v-for="(preview, index) in trackingPreviews" :key="index" class="preview-item">
              <span class="preview-order">{{ preview.order_no }}</span>
              <span class="preview-tracking">{{ preview.tracking_no }}</span>
            </div>
            <p v-if="trackingPreviews.length < orders.length" class="more-preview">
              {{ $t('orders.moreOrders', { count: orders.length - trackingPreviews.length }) }}
            </p>
          </div>
        </el-form-item>
      </el-form>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="visible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          {{ $t('orders.confirmBatchShip') }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { List, Van, Tickets, InfoFilled } from '@element-plus/icons-vue'
import { batchShipOrders, getCarrierList, type Order, type Carrier } from '@/api/order'
import { t } from '@/plugins/i18n'

const props = defineProps<{
  modelValue: boolean
  orders: Order[]
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
const loadingCarriers = ref(false)
const carriers = ref<Carrier[]>([])

const form = reactive({
  carrier_code: '',
  tracking_no_start: ''
})

const rules = {
  carrier_code: [{ required: true, message: t('orders.selectCarrier'), trigger: 'change' }],
  tracking_no_start: [{ required: true, message: t('orders.enterTrackingNoStart'), trigger: 'blur' }]
}

// Tracking previews
const trackingPreviews = computed(() => {
  if (!form.tracking_no_start) return []

  return props.orders.slice(0, 5).map((order, index) => {
    const baseTracking = form.tracking_no_start
    const numericPart = baseTracking.match(/\d+$/)?.[0] || ''
    const prefix = baseTracking.replace(/\d+$/, '')

    if (numericPart) {
      const num = parseInt(numericPart, 10) + index
      const paddedNum = num.toString().padStart(numericPart.length, '0')
      return {
        order_no: order.order_no,
        tracking_no: prefix + paddedNum
      }
    }

    return {
      order_no: order.order_no,
      tracking_no: baseTracking + (index > 0 ? `-${index}` : '')
    }
  })
})

// Format amount
const formatAmount = (amount: string) => {
  return parseFloat(amount).toFixed(2)
}

// Load carriers
const loadCarriers = async () => {
  loadingCarriers.value = true
  try {
    const res = await getCarrierList()
    carriers.value = res.filter(c => c.is_active)
  } catch (error) {
    console.error('Failed to load carriers:', error)
    ElMessage.error(t('orders.loadCarriersFailed'))
  } finally {
    loadingCarriers.value = false
  }
}

// Handle submit
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (!valid) return

    submitting.value = true
    try {
      await batchShipOrders({
        order_ids: props.orders.map(o => o.order_id),
        carrier_code: form.carrier_code,
        tracking_no_start: form.tracking_no_start
      })
      ElMessage.success(t('orders.batchShipSuccess', { count: props.orders.length }))
      emit('success')
      visible.value = false
    } catch (error: any) {
      ElMessage.error(error?.message || t('orders.batchShipFailed'))
    } finally {
      submitting.value = false
    }
  })
}

// Watch for visibility changes
watch(visible, (val) => {
  if (val) {
    loadCarriers()
    form.carrier_code = ''
    form.tracking_no_start = ''
  }
})
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

/* Orders List */
.orders-list {
  max-height: 200px;
  overflow-y: auto;
  border: 1px solid #E5E7EB;
  border-radius: 12px;
  padding: 8px;
}

.order-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  border-radius: 8px;
  transition: background 0.2s ease;
}

.order-item:hover {
  background: #F5F3FF;
}

.order-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.order-no {
  font-family: 'Fira Code', monospace;
  font-size: 14px;
  color: #6366F1;
  font-weight: 500;
}

.order-buyer {
  font-size: 12px;
  color: #6B7280;
}

.order-items {
  text-align: right;
}

.items-count {
  font-size: 13px;
  color: #1E1B4B;
  font-weight: 500;
}

.items-amount {
  display: block;
  font-size: 13px;
  color: #EF4444;
  font-weight: 600;
  margin-top: 4px;
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

.preview-order {
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