<template>
  <el-dialog
    v-model="visible"
    title="Create Shipment"
    width="600px"
    :close-on-click-modal="false"
  >
    <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
      <!-- Order Info -->
      <div class="order-info-section">
        <div class="section-header">
          <el-icon><Document /></el-icon>
          <span>Order Information</span>
        </div>
        <el-descriptions :column="2" border size="small">
          <el-descriptions-item label="Order No.">{{ shipment?.order_no }}</el-descriptions-item>
          <el-descriptions-item label="Status">
            <status-tag :status="shipment?.status" :type-map="statusTypeMap" />
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <!-- Logistics Info -->
      <div class="form-section">
        <div class="section-header">
          <el-icon><Van /></el-icon>
          <span>Logistics Information</span>
        </div>

        <el-form-item label="Carrier" prop="carrier_code">
          <el-select
            v-model="form.carrier_code"
            placeholder="Select carrier"
            style="width: 100%"
            @change="handleCarrierChange"
          >
            <el-option
              v-for="carrier in carriers"
              :key="carrier.code"
              :label="carrier.name"
              :value="carrier.code"
            />
          </el-select>
        </el-form-item>

        <el-form-item v-if="form.carrier_code === 'OTHER'" label="Custom Carrier" prop="carrier">
          <el-input v-model="form.carrier" placeholder="Enter carrier name" />
        </el-form-item>

        <el-form-item label="Tracking No." prop="tracking_no">
          <el-input v-model="form.tracking_no" placeholder="Enter tracking number">
            <template #prefix>
              <el-icon><Tickets /></el-icon>
            </template>
          </el-input>
        </el-form-item>

        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="Shipping Cost">
              <el-input-number
                v-model="form.shipping_cost"
                :min="0"
                :precision="2"
                style="width: 100%"
                placeholder="Optional"
              >
                <template #prefix>{{ form.currency }}</template>
              </el-input-number>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Weight (kg)">
              <el-input-number
                v-model="form.weight"
                :min="0"
                :precision="2"
                style="width: 100%"
                placeholder="Optional"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </div>

      <!-- Items Selection -->
      <div class="items-section">
        <div class="section-header">
          <el-icon><Goods /></el-icon>
          <span>Items to Ship</span>
          <el-tag size="small" type="info">{{ selectedItems.length }} selected</el-tag>
        </div>

        <div class="items-list">
          <div
            v-for="item in availableItems"
            :key="item.id"
            class="item-row"
            :class="{ selected: isItemSelected(item.id) }"
            @click="toggleItem(item)"
          >
            <el-checkbox :model-value="isItemSelected(item.id)" />
            <el-image :src="item.image" class="item-image" fit="cover">
              <template #error>
                <div class="image-placeholder">
                  <el-icon><Picture /></el-icon>
                </div>
              </template>
            </el-image>
            <div class="item-info">
              <p class="item-name">{{ item.product_name }}</p>
              <p class="item-sku">{{ item.sku_name }}</p>
            </div>
            <div class="item-quantity">
              <el-input-number
                v-if="isItemSelected(item.id)"
                :model-value="getItemQuantity(item.id)"
                @update:model-value="(val: number) => setItemQuantity(item.id, val)"
                :min="1"
                :max="item.quantity"
                size="small"
                style="width: 100px"
                @click.stop
              />
              <span v-else class="quantity-text">x {{ item.quantity }}</span>
            </div>
          </div>
        </div>

        <div class="select-actions">
          <el-button size="small" @click="selectAllItems">Select All</el-button>
          <el-button size="small" @click="clearItems">Clear</el-button>
        </div>
      </div>

      <!-- Remark -->
      <el-form-item label="Remark" class="remark-item">
        <el-input
          v-model="form.remark"
          type="textarea"
          :rows="3"
          placeholder="Optional notes"
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="visible = false">Cancel</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          Confirm Shipment
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Document, Van, Goods, Tickets, Picture } from '@element-plus/icons-vue'
import StatusTag from '@/components/common/StatusTag.vue'
import { createShipment, type Shipment, type Carrier } from '@/api/fulfillment'

const props = defineProps<{
  modelValue: boolean
  shipment: Shipment | null
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
const selectedItems = ref<number[]>([])

const form = reactive({
  carrier_code: '',
  carrier: '',
  tracking_no: '',
  shipping_cost: '0',
  weight: '0',
  currency: 'CNY',
  remark: '',
  items: [] as { order_item_id: number; quantity: number }[]
})

const rules = {
  carrier_code: [{ required: true, message: 'Please select a carrier', trigger: 'change' }],
  carrier: [{ required: true, message: 'Please enter carrier name', trigger: 'blur' }],
  tracking_no: [{ required: true, message: 'Please enter tracking number', trigger: 'blur' }]
}

const statusTypeMap = {
  0: { type: 'warning' as const, text: 'Pending' },
  1: { type: 'primary' as const, text: 'Shipped' },
  2: { type: 'info' as const, text: 'In Transit' },
  3: { type: 'success' as const, text: 'Delivered' },
  4: { type: 'danger' as const, text: 'Failed' }
}

const availableItems = computed(() => {
  return props.shipment?.items || []
})

const handleCarrierChange = (code: string) => {
  if (code !== 'OTHER') {
    form.carrier = ''
  }
}

const isItemSelected = (id: number) => {
  return selectedItems.value.includes(id)
}

const getItemQuantity = (id: number) => {
  const item = form.items.find(i => i.order_item_id === id)
  return item?.quantity || 1
}

const setItemQuantity = (id: number, quantity: number) => {
  const item = form.items.find(i => i.order_item_id === id)
  if (item) {
    item.quantity = quantity
  }
}

const toggleItem = (item: any) => {
  const index = selectedItems.value.indexOf(item.id)
  if (index === -1) {
    selectedItems.value.push(item.id)
    form.items.push({ order_item_id: item.id, quantity: item.quantity })
  } else {
    selectedItems.value.splice(index, 1)
    form.items = form.items.filter(i => i.order_item_id !== item.id)
  }
}

const selectAllItems = () => {
  selectedItems.value = availableItems.value.map(item => item.id)
  form.items = availableItems.value.map(item => ({
    order_item_id: item.id,
    quantity: item.quantity
  }))
}

const clearItems = () => {
  selectedItems.value = []
  form.items = []
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (!valid) return

    if (form.items.length === 0) {
      ElMessage.warning('Please select at least one item')
      return
    }

    submitting.value = true
    try {
      await createShipment({
        order_id: props.shipment!.order_id,
        carrier_code: form.carrier_code,
        carrier: form.carrier,
        tracking_no: form.tracking_no,
        shipping_cost: String(Math.round(parseFloat(form.shipping_cost) * 100)),
        shipping_currency: form.currency,
        weight: form.weight,
        remark: form.remark,
        items: form.items
      })
      emit('success')
      visible.value = false
    } catch (error) {
      ElMessage.error('Failed to create shipment')
    } finally {
      submitting.value = false
    }
  })
}

// Reset form when dialog opens
watch(visible, (val) => {
  if (val && props.shipment) {
    form.carrier_code = ''
    form.carrier = ''
    form.tracking_no = ''
    form.shipping_cost = '0'
    form.weight = '0'
    form.remark = ''
    selectedItems.value = props.shipment.items.map(item => item.id)
    form.items = props.shipment.items.map(item => ({
      order_item_id: item.id,
      quantity: item.quantity
    }))
  }
})
</script>

<style scoped>
.order-info-section,
.form-section,
.items-section {
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

/* Descriptions */
:deep(.el-descriptions__label) {
  font-weight: 500;
  color: #6B7280;
  background: #F9FAFB;
}

/* Items List */
.items-list {
  max-height: 280px;
  overflow-y: auto;
  border: 1px solid #E5E7EB;
  border-radius: 12px;
  padding: 8px;
}

.item-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.item-row:hover {
  background: #F5F3FF;
}

.item-row.selected {
  background: #EEF2FF;
  border: 1px solid rgba(99, 102, 241, 0.2);
}

.item-image {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  flex-shrink: 0;
}

.image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #F5F3FF 0%, #EEF2FF 100%);
  color: #6366F1;
}

.item-info {
  flex: 1;
}

.item-name {
  font-weight: 500;
  color: #1E1B4B;
  margin: 0 0 4px 0;
  font-size: 14px;
}

.item-sku {
  font-size: 12px;
  color: #6B7280;
  margin: 0;
}

.item-quantity {
  flex-shrink: 0;
}

.quantity-text {
  color: #6366F1;
  font-weight: 600;
}

.select-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 12px;
}

/* Remark */
.remark-item {
  margin-top: 16px;
}

/* Dialog Footer */
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>